package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// YouTubeTalent holds the schema definition for the YouTubeTalent entity.
type YouTubeTalent struct {
	ent.Schema
}

// Fields of the YouTubeTalent.
func (YouTubeTalent) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().
			Comment("YouTube channel ID"),
		field.String("channel_name").
			Comment("YouTube channel name"),
		field.String("thumbnail_url").
			Comment("URL of the talent's YouTube thumbnail"),
		field.String("membership_video_id").
			Optional().
			Comment("ID of a members-only video"),
		field.Time("last_membership_video_id_miss").
			Optional().
			Comment("Last time membership_video_id returned no results"),
		field.Time("last_updated").
			Default(time.Now).
			Comment("Last time data was fetched"),
		field.Time("disabled").
			Optional().
			Comment("When refresh/membership checks were disabled. Set to zero value to re-enable."),
	}
}

// Edges of the YouTubeTalent.
func (YouTubeTalent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("guilds", Guild.Type),
		edge.To("roles", GuildRole.Type),
		edge.From("memberships", UserMembership.Type).
			Ref("youtube_talent"),
	}
}
