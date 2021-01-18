package common

import (
	"time"

	"cloud.google.com/go/firestore"
)

// These are names for Firestore collections. For preventing dumb typos!
const (
	UsersCollection         = "users"
	ChannelCollection       = "channels"
	ChannelCheckCollection  = "check"
	ChannelMemberCollection = "members"
	DiscordGuildCollection  = "guilds"
	PrivateCollection       = "private"
	DMTemplateCollection    = "dm-templates"
)

// DiscordIdentity is saved to Firestore in the `DiscordIdentityCollection` collection.
// fs: users/{UserID}
type DiscordIdentity struct {
	UserID            string
	Username          string
	Discriminator     string
	YoutubeChannelID  string                   // the user's linked Youtube Channel ID
	CandidateChannels []*firestore.DocumentRef // possibly relevant YouTube channels
	Memberships       []*firestore.DocumentRef // verified memberships
}

// Channel defines a YouTube channel whose membership we might check.
// fs: channels/{slug}
type Channel struct {
	ChannelID,
	ChannelTitle,
	Thumbnail string
}

// ChannelCheck defines channel membership check criteria.
// fs: channels/{slug}/check/check
type ChannelCheck struct {
	VideoID string
}

// ChannelMember defines channel membership.
// fs: channels/{slug}/members/{UserID}
type ChannelMember struct {
	DiscordID string
	Timestamp time.Time
}

// DiscordGuild defines loaded Discord guild associations.
// fs: guilds/{guildID]
type DiscordGuild struct {
	ID                    string
	Name                  string
	LastMembershipRefresh time.Time `json:",omitempty"`
	AuditLogChannelID     string
	BCP47                 string            `json:",omitempty"`
	MembershipRoles       map[string]string // channelSlug -> roleID. Please never leave RoleID empty.
}

// DMTemplate is a message template that can be sent out (en masse).
// fs: dm-templates/[name]
type DMTemplate struct {
	Name string
	Body string
}

// DMTemplateData is the data passed in a text/template.Execute() call with a DMTemplate.
type DMTemplateData struct {
	User  *DiscordIdentity
	Extra map[string]interface{}
}
