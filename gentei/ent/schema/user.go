package schema

import (
	"regexp"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"golang.org/x/oauth2"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").Unique().
			Comment("Discord user snowflake"),
		field.String("full_name").
			Match(regexp.MustCompile(`^.+?#\d{4}$`)).
			Comment("Username + discriminator"),
		field.String("avatar_hash"),
		field.Time("last_check").
			Default(func() time.Time {
				return time.Time{}
			}).
			Comment("Timestamp of last membership check"),
		field.String("youtube_id").
			Unique().Nillable().
			Optional().
			Comment("user's YouTube channel ID"),
		field.JSON("youtube_token", &oauth2.Token{}).
			Optional(),
		field.JSON("discord_token", &oauth2.Token{}).
			Optional(), // TODO: remove Optional() on next PR
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guilds", Guild.Type).
			Ref("members").
			Comment("Guild that this user has joined"),
		edge.From("guilds_admin", Guild.Type).
			Ref("admins"),
		edge.To("memberships", UserMembership.Type),
	}
}
