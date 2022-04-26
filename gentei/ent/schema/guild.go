package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// GuildSettings contains customization options for a server.
type GuildSettings struct {
	MaxStale time.Duration
}

// Guild holds the schema definition for the Guild entity.
type Guild struct {
	ent.Schema
}

// Fields of the Guild.
func (Guild) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Unique().
			Comment("Discord guild ID"),
		field.String("name").
			Comment("Discord guild name"),
		field.String("icon_hash").
			Optional().
			Comment("Discord guild icon hash"),
		field.Uint64("audit_channel").
			Optional().
			Unique().
			Comment("Audit log channel ID"),
		field.Enum("language").
			Values("en-US").
			Default("en-US").
			Comment("IETF BCP 47 language tag"),
		field.JSON("admin_snowflakes", []uint64{}).
			Comment("Discord snowflakes of users and groups that can modify server settings. The first snowflake is always the server owner."),
		field.JSON("moderator_snowflakes", []uint64{}).
			Optional().
			Comment("Discord snowflakes of users and groups that can read server settings"),
		field.JSON("settings", &GuildSettings{}).
			Optional(),
	}
}

// Edges of the Guild.
func (Guild) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("members", User.Type),
		edge.To("admins", User.Type),
		edge.To("roles", GuildRole.Type),
		edge.From("youtube_talents", YouTubeTalent.Type).
			Ref("guilds"),
	}
}
