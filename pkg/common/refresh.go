package common

import (
	"context"
	"fmt"

	"google.golang.org/api/youtube/v3"

	"cloud.google.com/go/firestore"
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

// RefreshYouTubeChannels updates transient YouTube channel data.
func RefreshYouTubeChannels(ctx context.Context, fs *firestore.Client, yt *youtube.Service) error {
	docs, err := fs.Collection(ChannelCollection).Documents(ctx).GetAll()
	if err != nil {
		log.Err(err).Msg("error creating Firestore client")
		return err
	}
	for _, doc := range docs {
		var (
			channel Channel
			logger  = log.With().Str("slug", doc.Ref.ID).Logger()
		)
		err = doc.DataTo(&channel)
		if err != nil {
			logger.Err(err).Msg("error unmarshalling doc to Channel")
			return err
		}
		clr, err := yt.Channels.List([]string{"snippet"}).Id(channel.ChannelID).Do()
		if err != nil {
			logger.Err(err).Msg("error calling channels.list")
			return err
		}
		var (
			updates []firestore.Update
			snippet = clr.Items[0].Snippet
		)
		if snippet.Title != channel.ChannelTitle {
			updates = append(updates, firestore.Update{
				Path:  "ChannelTitle",
				Value: snippet.Title,
			})
		}
		if snippet.Thumbnails.High.Url != channel.Thumbnail {
			updates = append(updates, firestore.Update{
				Path:  "Thumbnail",
				Value: snippet.Thumbnails.High.Url,
			})
		}
		if len(updates) > 0 {
			_, err = doc.Ref.Update(ctx, updates)
			if err != nil {
				logger.Err(err).Msg("error setting new channel data")
				return err
			}
			logger.Info().Interface("updates", updates).Msg("updated channel data")
		}
	}
	return nil
}

// RefreshDiscordGuilds updates transient Discord guild data.
func RefreshDiscordGuilds(ctx context.Context, fs *firestore.Client, dg *discordgo.Session) error {
	if dg == nil {
		return fmt.Errorf("must specify a discordgo.Session")
	}
	docs, err := fs.Collection(DiscordGuildCollection).Documents(ctx).GetAll()
	if err != nil {
		log.Err(err).Msg("error getting Discord guilds from firestore")
		return err
	}
	for _, doc := range docs {
		var (
			logger = log.With().Str("guildID", doc.Ref.ID).Logger()
			guild  DiscordGuild
		)
		err = doc.DataTo(&guild)
		if err != nil {
			logger.Err(err).Msg("error unmarshalling Guild")
			return err
		}
		if guild.APIOnly {
			continue
		}
		dgGuild, err := dg.Guild(doc.Ref.ID)
		if err != nil {
			logger.Err(err).Msg("error getting guild info")
			return err
		}
		var updates []firestore.Update
		if guild.Name != dgGuild.Name {
			updates = append(updates, firestore.Update{
				Path:  "Name",
				Value: dgGuild.Name,
			})
		}
		if len(updates) > 0 {
			_, err = doc.Ref.Update(ctx, updates)
			if err != nil {
				logger.Err(err).Msg("error updating guild info")
				return err
			}
			logger.Info().Interface("updates", updates).Msg("updated guild info")
		}
	}
	return nil
}
