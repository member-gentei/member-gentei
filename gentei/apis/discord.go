package apis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"libs.altipla.consulting/tokensource"
)

func GetRefreshingDiscordTokenSource(ctx context.Context, db *ent.Client, config *oauth2.Config, userID uint64) (oauth2.TokenSource, error) {
	u, err := db.User.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	logger := log.With().Str("userID", strconv.FormatUint(userID, 10)).Logger()
	notify := tokensource.NewNotifyHook(ctx, config, u.DiscordToken, func(token *oauth2.Token) error {
		logger.Debug().Msg("Discord token for user refreshed")
		err := db.User.UpdateOneID(userID).
			SetDiscordToken(token).
			Exec(ctx)
		if err != nil {
			return fmt.Errorf("error saving refreshed Discord token: %w", err)
		}
		return err
	})
	return notify, nil
}
