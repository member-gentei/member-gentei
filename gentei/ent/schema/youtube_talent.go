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
		field.Time("last_updated").
			Default(time.Now).
			Comment("Last time data was fetched"),
	}
}

// Edges of the YouTubeTalent.
func (YouTubeTalent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("guilds", Guild.Type),
	}
}
