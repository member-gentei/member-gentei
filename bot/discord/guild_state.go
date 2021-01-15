package discord

import (
	"github.com/member-gentei/member-gentei/pkg/common"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type guildLoadState int

const (
	guildFirstEncounter guildLoadState = iota
	guildWaitingForAssociationData
	guildWaitingForCreateEvent
	guildLoaded
)

type guildState struct {
	Doc       common.DiscordGuild
	LoadState guildLoadState
	localizer *i18n.Localizer

	noFancyReply bool // whether we can use message replies instead of @user in this guild
}

// GetMembershipInfo retrieves channelSlug -> membership role information.
func (g guildState) GetMembershipInfo() map[string]string {
	if g.Doc.Channel == nil || g.Doc.MemberRoleID == "" || g.Doc.Channel.ID == "" {
		return nil
	}
	return map[string]string{g.Doc.Channel.ID: g.Doc.MemberRoleID}
}
