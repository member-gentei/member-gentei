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
	AdministrativeRoles   []string
	LastMembershipRefresh time.Time
	MemberRoleID          string
	AuditLogChannelID     string
	Channel               *firestore.DocumentRef
}
