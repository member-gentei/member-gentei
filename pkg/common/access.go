package common

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

// RevokeYouTubeAccess removes any and all YouTube stuff from a user ID.
func RevokeYouTubeAccess(ctx context.Context, fs *firestore.Client, userID string) error {
	var (
		logger        = log.With().Str("userID", userID).Logger()
		userDocRef    = fs.Collection(UsersCollection).Doc(userID)
		youtubeDocRef = userDocRef.Collection("private").Doc("youtube")
	)
	// get the refresh token to delete. If it doesn't exist, move on to deleting memberships anyway.
	var token oauth2.Token
	youtubeDoc, err := youtubeDocRef.Get(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("error getting user's stored YouTube token, skipping revocation")
	} else {
		err = youtubeDoc.DataTo(&token)
		if err != nil {
			logger.Err(err).Msg("error unmarshalling user's stored YouTube token")
			return err
		}
		revokeYoutubeToken(token.RefreshToken, logger)
		// delete YouTube token
		_, err = youtubeDocRef.Delete(ctx)
		if err != nil {
			logger.Err(err).Msg("error deleting user's YouTube token")
			return err
		}
	}
	// delete memberships
	_, err = userDocRef.Update(ctx, []firestore.Update{
		{
			Path:  "YoutubeChannelID",
			Value: "",
		},
		{
			Path:  "Memberships",
			Value: firestore.Delete,
		},
	})
	if err != nil {
		logger.Err(err).Msg("error removing user fields")
		return err
	}
	snaps, err := fs.CollectionGroup(ChannelMemberCollection).
		Where("DiscordID", "==", userID).Select().Documents(ctx).GetAll()
	if err != nil {
		logger.Err(err).Msg("error getting user memberships")
		return err
	}
	for _, snap := range snaps {
		_, err = snap.Ref.Delete(ctx)
		if err != nil {
			logger.Err(err).Msg("error deleting user membership doc")
			return err
		}
	}
	return nil
}

// DeleteUser removes all traces of a user.
func DeleteUser(ctx context.Context, fs *firestore.Client, userID string) (err error) {
	// revoke YouTube access, which supports not having access to revoke + deletes all memberships
	err = RevokeYouTubeAccess(ctx, fs, userID)
	if err != nil {
		log.Err(err).Msg("error removing YouTube token and memberships")
		return
	}
	// delete user object and its collections (its tokens)
	userDoc := fs.Collection(UsersCollection).Doc(userID)
	iter := userDoc.Collections(ctx)
	collectionRefs, err := iter.GetAll()
	if err != nil {
		log.Err(err).Msg("error getting user collections")
		return
	}
	for _, collectionRef := range collectionRefs {
		// n.b. DocumentRefs includes missing documents
		docRefs, err := collectionRef.DocumentRefs(ctx).GetAll()
		if err != nil {
			log.Err(err).Str("collectionRef", collectionRef.Path).
				Msg("error getting subcollection docref")
			return err
		}
		for _, docRef := range docRefs {
			_, err = docRef.Delete(ctx)
			if err != nil {
				log.Err(err).Str("collectionRef", collectionRef.Path).
					Msg("error deleting subcollection doc")
			}
			return err
		}
	}
	_, err = userDoc.Delete(ctx)
	if err != nil {
		log.Err(err).Msg("error deleting user object")
		return
	}
	return
}
