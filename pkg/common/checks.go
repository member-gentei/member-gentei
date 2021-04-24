package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrDiscordTokenInvalid denotes that an OAuth token has expired or been revoked
	ErrDiscordTokenInvalid = errors.New("Discord token expired or revoked")

	// ErrDiscordTokenNotFound denotes that an OAuth token has been removed
	ErrDiscordTokenNotFound = errors.New("Discord token removed")

	// ErrYouTubeTokenInvalid denotes that an OAuth token has expired or been revoked
	ErrYouTubeTokenInvalid = errors.New("YouTube token expired or revoked")

	// ErrYouTubeInvalidGrant denotes an ephemeral invalid_grant error that can mean anything, really
	ErrYouTubeInvalidGrant = errors.New("generic invalid_grant error")

	// ErrYouTubeUserGone
	ErrYouTubeUserUnavailable = errors.New("YouTube user account closed or suspended")
)

// EnforceMembershipsOptions is the multiselect-y options struct for EnforceMemberships
type EnforceMembershipsOptions struct {
	ReloadDiscordGuilds       bool
	OnlyChannelSlug           string
	RemoveInvalidDiscordToken bool // removes users with permanently invalid (revoked, de-scoped) tokens
	RemoveInvalidYouTubeToken bool // removes users with permanently invalid (revoked, de-scoped) tokens
	Apply                     bool // apply changes
	UserIDs                   []string

	// only refresh memberships that were validated before this time
	RefreshBefore time.Time

	// amount of worker threads to use (default is 1)
	NumWorkers uint
}

// EnforceMembershipsResult contains metrics useful for monitoring/debugging/fun.
type EnforceMembershipsResult struct {
	UserCount uint

	// Number of users that have disconnected/removed their YouTube or Discord accounts outside of the main website UI.
	UsersDisconnected uint

	// Number of membership lapses during enforcement
	MembershipsLapsed uint

	// Number of memberships added during enforcement
	MembershipsAdded uint

	// Number of memberships re-confirmed during enforcement
	MembershipsReconfirmed uint
}

// EnforceMemberships checks all users' memberships against candidate channels.
func EnforceMemberships(ctx context.Context, fs *firestore.Client, options *EnforceMembershipsOptions) (result EnforceMembershipsResult, err error) {
	var query firestore.Query
	if len(options.UserIDs) > 0 {
		query = fs.Collection(UsersCollection).Where("UserID", "in", options.UserIDs)
	} else {
		if options.OnlyChannelSlug != "" {
			query = fs.Collection(UsersCollection).Where("CandidateChannels", "array-contains", fs.Collection(ChannelCollection).Doc(options.OnlyChannelSlug))
		} else {
			query = fs.Collection(UsersCollection).Query
		}
		if !options.RefreshBefore.IsZero() {
			query = query.Where("LastRefreshed", "<", options.RefreshBefore)
		}
		if options.ReloadDiscordGuilds {
			query = query.Select()
		} else {
			query = query.Select("CandidateChannels")
		}
	}
	// cache so that we don't have to perform a lot of expensive array-in queries
	slug2MemberVideos, err := getMemberVideoIDs(ctx, fs)
	if err != nil {
		log.Err(err).Msg("error getting member video IDs")
		return
	}
	// we should be able to slowly paginate through userIDs, but Firestore returns an internal error more often than not when we paginate this query.
	// load it all into RAM!
	startTime := time.Now()
	docs, err := query.Documents(ctx).GetAll()
	if err != nil {
		log.Err(err).Msg("error getting all user IDs")
		return
	}
	result.UserCount = uint(len(docs))
	log.Debug().
		Str("duration", time.Since(startTime).String()).
		Uint("count", result.UserCount).
		Msg("loaded user IDs")
	var (
		numWorkers   = int(math.Max(1, float64(options.NumWorkers)))
		wg           = &sync.WaitGroup{}
		workerWG     = &sync.WaitGroup{}
		docsChan     = make(chan *firestore.DocumentSnapshot, numWorkers)
		resultChan   = make(chan enforceMembershipsWorkerResult, numWorkers)
		cctx, cancel = context.WithCancel(ctx)
	)
	defer cancel()
	// start workers
	for i := 0; i < numWorkers; i++ {
		workerWG.Add(1)
		go enforceMembershipsWorker(
			cctx, fs, options, slug2MemberVideos,
			docsChan, resultChan, workerWG,
		)
	}
	// when all of the workers are done, close the channel
	wg.Add(3)
	go func() {
		defer wg.Done()
		log.Debug().Msg("waiting for workers to complete")
		workerWG.Wait()
		close(resultChan)
		log.Debug().Msg("workers are done")
	}()
	// start doc producer
	go func() {
		defer wg.Done()
		for _, doc := range docs {
			docsChan <- doc
		}
		close(docsChan)
	}()
	// start results consumer
	go func() {
		defer wg.Done()
		for workerResult := range resultChan {
			if workerResult.err != nil {
				log.Err(err).Msg("recieved error from worker, cancelling context")
				cancel()
				// deplete the error queue
				go func() {
					for result := range resultChan {
						if result.err != nil {
							log.Err(err).Msg("post-cancellation")
						}
					}
				}()
				return
			}
			// aggregate
			result.UsersDisconnected += workerResult.UsersDisconnected
			result.MembershipsLapsed += workerResult.MembershipsLapsed
			result.MembershipsAdded += workerResult.MembershipsAdded
			result.MembershipsReconfirmed += workerResult.MembershipsReconfirmed
		}
	}()
	// docs producer + workers done + worker result consumer
	wg.Wait()
	return
}

type enforceMembershipsWorkerResult struct {
	EnforceMembershipsResult
	err error
}

func enforceMembershipsWorker(
	ctx context.Context, fs *firestore.Client, options *EnforceMembershipsOptions,
	slug2MemberVideos map[string]string,
	docs <-chan *firestore.DocumentSnapshot, resultChan chan<- enforceMembershipsWorkerResult, wg *sync.WaitGroup,
) {
	defer wg.Done()
	for doc := range docs {
		select {
		case <-ctx.Done():
			log.Warn().Msg("worker aborting")
			return
		default:
		}
		// acquire candidate YouTube channels (via a Discord refresh or otherwise)
		var (
			result            enforceMembershipsWorkerResult
			candidateChannels []*firestore.DocumentRef
			userID            = doc.Ref.ID
			logger            = log.With().Str("userID", userID).Logger()
			err               error
		)
		if options.ReloadDiscordGuilds {
			candidateChannels, err = ReloadDiscordGuilds(ctx, fs, userID)
			if errors.Is(err, ErrDiscordTokenInvalid) || errors.Is(err, ErrDiscordTokenNotFound) {
				if options.RemoveInvalidDiscordToken {
					logger.Warn().Err(err).Msg("Discord token invalid, deleting user")
					err = DeleteUser(ctx, fs, userID)
					if err != nil {
						logger.Err(err).Msg("error deleting user")
						result.err = err
						resultChan <- result
						return
					}
					result.UsersDisconnected++
				} else {
					logger.Warn().Err(err).Msg("Discord token invalid, skipping user")
					err = nil
				}
				resultChan <- result
				continue
			} else if err != nil {
				logger.Err(err).Msg("error reloading Discord guilds for user")
				result.err = err
				resultChan <- result
				return
			}
		} else {
			var user DiscordIdentity
			err := doc.DataTo(&user)
			if err != nil {
				logger.Err(err).Msg("error unmarshalling user")
				result.err = err
				resultChan <- result
				return
			}
			candidateChannels = user.CandidateChannels
		}
		// check memberships
		var (
			verifiedMemberships = map[string]time.Time{}
			skipUser            bool
		)
		for _, candidateRef := range candidateChannels {
			checkLogger := logger.With().Str("channelSlug", candidateRef.ID).Logger()
			checkOpts := &CheckChannelMembershipOptions{
				UserID:                   userID,
				ChannelMembershipVideoID: slug2MemberVideos[candidateRef.ID],
				Logger:                   &checkLogger,
			}
			if checkOpts.ChannelMembershipVideoID == "" {
				checkOpts.ChannelSlug = candidateRef.ID
				logger.Warn().Str("channelSlug", checkOpts.ChannelSlug).Msg("could not find membership video ID for candidate channel")
			}
			isMember, err := CheckChannelMembership(ctx, fs, checkOpts)
			if errors.Is(err, ErrYouTubeTokenInvalid) || status.Code(err) == codes.NotFound {
				if options.RemoveInvalidYouTubeToken {
					logger.Warn().Err(err).Msg("YouTube token invalid for user, removing token and memberships")
					err = RevokeYouTubeAccess(ctx, fs, userID)
					if err != nil {
						logger.Err(err).Msg("error revoking YouTube access for user")
						result.err = err
						resultChan <- result
						return
					}
					result.UsersDisconnected++
				} else {
					logger.Info().Err(err).Msg("YouTube token invalid for user, skipping")
					err = nil
				}
				skipUser = true
				break
			} else if errors.Is(err, ErrYouTubeInvalidGrant) {
				logger.Warn().Err(err).Msg("mystery invalid_grant, need to retry user's checks later")
				skipUser = true
				break
			} else if err != nil {
				logger.Err(err).Msg("unhandled error while checking membership for user")
				result.err = err
				resultChan <- result
				return
			}
			if isMember {
				verifiedMemberships[candidateRef.ID] = time.Now().In(time.UTC)
			}
		}
		if !options.Apply {
			logger.Info().Interface("memberships", verifiedMemberships).Msg("verified memberships")
		}
		if skipUser {
			_, err = fs.Collection(UsersCollection).Doc(userID).Update(ctx, []firestore.Update{{
				Path:  "LastRefreshed",
				Value: time.Now().In(time.UTC),
			}})
			result.err = err
			resultChan <- result
			continue
		}
		if options.Apply {
			logger.Debug().Interface("memberships", verifiedMemberships).Msg("setting memberships")
			// query for existing
			var (
				toDelete            []*firestore.DocumentRef
				selects             []*firestore.DocumentSnapshot
				existingMemberships = map[string]struct{}{}
			)
			selects, err := fs.CollectionGroup("members").Where("DiscordID", "==", userID).Select().Documents(ctx).GetAll()
			if err != nil {
				logger.Err(err).Msg("error querying for members CollectionGroup")
				result.err = err
				resultChan <- result
				return
			}
			for _, selected := range selects {
				slug := selected.Ref.Parent.Parent.ID
				if _, found := verifiedMemberships[slug]; !found {
					logger.Info().Str("slug", slug).Msg("membership lapsed for channel")
					toDelete = append(toDelete, selected.Ref)
					result.MembershipsLapsed++
				}
				existingMemberships[slug] = struct{}{}
			}
			// delete stale
			for _, stale := range toDelete {
				_, err = stale.Delete(ctx)
				if err != nil {
					logger.Err(err).Msg("error deleting stale membership")
					result.err = err
					resultChan <- result
					return
				}
			}
			// upsert new
			for channelSlug, ts := range verifiedMemberships {
				docRef := fs.Collection("channels").Doc(channelSlug).Collection("members").Doc(userID)
				_, err = docRef.Set(ctx, map[string]interface{}{
					"DiscordID": userID,
					"Timestamp": ts,
				})
				if err != nil {
					logger.Err(err).Msg("error adding member verification")
					result.err = err
					resultChan <- result
					return
				}
				if _, exists := existingMemberships[channelSlug]; exists {
					result.MembershipsReconfirmed++
				} else {
					result.MembershipsAdded++
				}
			}
			// update user
			membershipDocRefs := make([]*firestore.DocumentRef, 0, len(verifiedMemberships))
			for slug := range verifiedMemberships {
				membershipDocRefs = append(membershipDocRefs, fs.Collection("channels").Doc(slug))
			}
			_, err = fs.Collection("users").Doc(userID).Update(ctx,
				[]firestore.Update{
					{
						Path:  "Memberships",
						Value: membershipDocRefs,
					},
					{
						Path:  "LastRefreshed",
						Value: time.Now().In(time.UTC),
					},
				},
			)
			if err != nil {
				logger.Err(err).Msg("error updating user object memberships")
				result.err = err
				resultChan <- result
				return
			}
		}
		resultChan <- result
	}
}

// CheckChannelMembershipOptions is the multiselect-y options struct for CheckChannelMembership
type CheckChannelMembershipOptions struct {
	UserID      string
	UserService *youtube.Service

	ChannelSlug              string
	ChannelMembershipVideoID string

	Logger *zerolog.Logger // optional logging context
}

const checkRetries = 4

// CheckChannelMembership checks a user's membership against a channel.
func CheckChannelMembership(
	ctx context.Context, fs *firestore.Client,
	options *CheckChannelMembershipOptions,
) (isMember bool, err error) {
	var (
		userService        = options.UserService
		memberCheckVideoID = options.ChannelMembershipVideoID
		logger             zerolog.Logger
	)
	if options.Logger == nil {
		logger = log.Logger
	} else {
		logger = *options.Logger
	}
	// acquire a membership check video
	if memberCheckVideoID == "" {
		if options.ChannelSlug == "" {
			return false, fmt.Errorf("must specify ChannelSlug or ChannelMembershipVideoID")
		}
		var snap *firestore.DocumentSnapshot
		snap, err = fs.Collection(ChannelCollection).Doc(options.ChannelSlug).
			Collection(ChannelCheckCollection).Doc(ChannelCheckCollection).Get(ctx)
		if err != nil {
			logger.Err(err).Str("channelSlug", options.ChannelSlug).Msg("error getting membership check video")
			return
		}
		var checkDoc ChannelCheck
		err = snap.DataTo(&checkDoc)
		if err != nil {
			logger.Err(err).Str("channelSlug", options.ChannelSlug).Msg("error unmarshalling membership check video doc")
			return
		}
		memberCheckVideoID = checkDoc.VideoID
	}
	// acquire a youtube.Service in the user's context
	if userService == nil {
		if options.UserID == "" {
			return false, fmt.Errorf("must specify UserID or UserService")
		}
		userService, err = GetYouTubeService(ctx, fs, options.UserID)
		if err != nil {
			logger.Err(err).Str("userID", options.UserID).Msg("error getting YouTube service for user")
			return
		}
	}
	// perform The Membership Check
	var ctr *youtube.CommentThreadListResponse
	for i := 0; i < checkRetries; i++ {
		ctr, err = userService.CommentThreads.List([]string{"id"}).VideoId(memberCheckVideoID).Do()
		if err != nil {
			errString := err.Error()
			if strings.HasSuffix(errString, "commentsDisabled") {
				err = nil
				return
			} else if strings.Contains(errString, "Token has been expired or revoked.") {
				logger.Warn().Err(err).Send()
				err = ErrYouTubeTokenInvalid
				return
			} else if strings.Contains(errString, "Request had invalid authentication credentials") {
				logger.Warn().Err(err).Send()
				err = ErrYouTubeTokenInvalid
			} else if strings.Contains(errString, `"error": "invalid_grant",`) {
				logger.Warn().Err(err).Send()
				err = ErrYouTubeInvalidGrant
			} else if strings.Contains(errString, "Invalid \\\"invalid_grant\\\" in request.") {
				logger.Warn().Err(err).Send()
				err = ErrYouTubeInvalidGrant
				return
			} else if strings.HasSuffix(errString, "authenticatedUserAccountClosed") || strings.HasSuffix(errString, "authenticatedUserAccountSuspended") {
				logger.Warn().Err(err).Send()
				err = ErrYouTubeTokenInvalid
				return
			} else if strings.Contains(errString, "processingFailure") || strings.HasSuffix(errString, "videoNotFound") {
				// retry with linear backoff
				logger.Warn().Int("try", i+1).Err(err).Msg("membership check attempt failed")
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
			logger.Err(err).Msg("error getting comment threads for video")
			return
		}
		break
	}
	if err == nil {
		isMember = true
		logger.Info().
			Int("commentThreads", len(ctr.Items)).
			Str("memberCheckVideoID", memberCheckVideoID).
			Msg("confirmed membership")
	}
	return
}

const discordMeGuildsURL = "https://discord.com/api/users/@me/guilds"

// ReloadDiscordGuilds reloads Discord guilds for a user and returns the new candidate
// channels. Results usually piped into CheckChannelMembership().
func ReloadDiscordGuilds(
	ctx context.Context, fs *firestore.Client, userID string,
) (candidateChannels []*firestore.DocumentRef, err error) {
	httpClient, err := GetDiscordHTTPClient(ctx, fs, userID)
	if status.Code(err) == codes.NotFound {
		err = ErrDiscordTokenNotFound
		return
	} else if err != nil {
		log.Err(err).Msg("error getting Discord client for user")
		return
	}
	// load guilds
	response, err := httpClient.Get(discordMeGuildsURL)
	if err != nil {
		// oauth2 client complains about the refresh token via this GET. Annoyingly,
		// the http client mangles it real bad and we can't cast the error conventionally!
		if rErr, ok := scavengeRetrieveError(response, err); ok {
			var errResponse struct {
				Error            string
				ErrorDescription string `json:"error_description"`
			}
			// if this fails to unmarshal, we return the error as-is anyway
			json.Unmarshal(rErr.Body, &errResponse)
			if errResponse.ErrorDescription == `Invalid "refresh_token" in request.` {
				err = ErrDiscordTokenInvalid
			} else if errResponse.Error == "invalid_grant" {
				err = ErrDiscordTokenInvalid
			}
			log.Debug().Str("userID", userID).Interface("retrieveError", errResponse).Send()
		}
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode == http.StatusUnauthorized {
		err = ErrDiscordTokenInvalid
		log.Warn().Bytes("body", body).
			Msg("401 getting discord guilds for user")
		return
	}
	if response.StatusCode != http.StatusOK {
		var jsonErr struct {
			Code    int
			Message string
		}
		json.Unmarshal(body, &jsonErr)
		switch jsonErr.Code {
		case 50014, 50025:
			err = ErrDiscordTokenInvalid
			return
		}
		log.Warn().Int("code", response.StatusCode).Bytes("body", body).
			Msg("non-200 status getting discord guilds for user")
		return
	}
	var guildMemberships []struct {
		ID string
	}
	err = json.Unmarshal(body, &guildMemberships)
	if err != nil {
		return
	}
	if len(guildMemberships) == 0 {
		return
	}
	// load heavily cached guildID -> channel map
	guildsToChannels, err := getCachedGuildsToChannels(ctx, fs)
	if err != nil {
		log.Err(err).Msg("error getting guild to channel mapping")
		return
	}
	candidateMap := map[string]*firestore.DocumentRef{}
	for _, datum := range guildMemberships {
		if channelRefs := guildsToChannels[datum.ID]; channelRefs != nil {
			for _, channelRef := range channelRefs {
				candidateMap[channelRef.ID] = channelRef
			}
		}
	}
	for _, channelRef := range candidateMap {
		candidateChannels = append(candidateChannels, channelRef)
	}
	// sort by docID
	sort.Slice(candidateChannels, func(i, j int) bool {
		return candidateChannels[i].ID < candidateChannels[j].ID
	})
	// write to user object
	_, err = fs.Collection(UsersCollection).Doc(userID).Update(ctx, []firestore.Update{
		{
			Path:  "CandidateChannels",
			Value: candidateChannels,
		},
	})
	return
}

// SetUserMemberships sets user memberships in the appropriate places across Firestore.
func SetUserMemberships(
	ctx context.Context, fs *firestore.Client, userID string,
	verifiedMemberships map[string]time.Time,
) error {
	var (
		toDelete            []*firestore.DocumentRef
		existingMemberships = map[string]struct{}{}
	)
	selects, err := fs.CollectionGroup(ChannelMemberCollection).Where("DiscordID", "==", userID).Select().Documents(ctx).GetAll()
	if err != nil {
		log.Err(err).Msg("error querying for members CollectionGroup")
		return err
	}
	for _, selected := range selects {
		slug := selected.Ref.Parent.Parent.ID
		if _, found := verifiedMemberships[slug]; !found {
			log.Info().Str("slug", slug).Msg("membership lapsed for channel")
			toDelete = append(toDelete, selected.Ref)
		} else {
			existingMemberships[slug] = struct{}{}
		}
	}
	// delete stale
	for _, stale := range toDelete {
		_, err = stale.Delete(ctx)
		if err != nil {
			log.Err(err).Msg("error deleting stale membership")
			return err
		}
	}
	// upsert new
	for channelSlug, ts := range verifiedMemberships {
		if _, isExisting := existingMemberships[channelSlug]; !isExisting {
			log.Info().Str("slug", channelSlug).Str("userID", userID).Msg("adding new membership")
		}
		docRef := fs.Collection(ChannelCollection).Doc(channelSlug).Collection(ChannelMemberCollection).Doc(userID)
		_, err = docRef.Set(ctx, map[string]interface{}{
			"DiscordID": userID,
			"Timestamp": ts,
		})
		if err != nil {
			log.Err(err).Msg("error adding member verification")
			return err
		}
	}
	// update user
	membershipDocRefs := make([]*firestore.DocumentRef, 0, len(verifiedMemberships))
	for slug := range verifiedMemberships {
		membershipDocRefs = append(membershipDocRefs, fs.Collection(ChannelCollection).Doc(slug))
	}
	_, err = fs.Collection(UsersCollection).Doc(userID).Update(ctx, []firestore.Update{{
		Path:  "Memberships",
		Value: membershipDocRefs,
	}})
	if err != nil {
		log.Err(err).Msg("error updating user object memberships")
		return err
	}
	return nil
}

var cachedGuildsToChannels = make(map[string][]*firestore.DocumentRef)

func getCachedGuildsToChannels(ctx context.Context, fs *firestore.Client) (map[string][]*firestore.DocumentRef, error) {
	if len(cachedGuildsToChannels) > 0 {
		return cachedGuildsToChannels, nil
	}
	snaps, err := fs.Collection(DiscordGuildCollection).Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	chCollection := fs.Collection(ChannelCollection)
	for _, snap := range snaps {
		var guild DiscordGuild
		err = snap.DataTo(&guild)
		if err != nil {
			log.Err(err).Msg("error unmarshalling DiscordGuild")
			return nil, err
		}
		channelRefs := make([]*firestore.DocumentRef, 0, len(guild.MembershipRoles))
		for channelSlug := range guild.MembershipRoles {
			channelRefs = append(channelRefs, chCollection.Doc(channelSlug))
		}
		cachedGuildsToChannels[guild.ID] = channelRefs
	}
	return cachedGuildsToChannels, nil
}
