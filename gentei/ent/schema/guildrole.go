package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// GuildRole holds the schema definition for the GuildRole entity.
type GuildRole struct {
	ent.Schema
}

// Fields of the GuildRole.
func (GuildRole) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			Unique().Immutable().
			Comment("Discord snowflake for this role"),
		field.String("name").
			Comment("Human name for this role"),
		field.Time("last_updated").Default(time.Now).
			Comment("When the name was last synchronized for this role"),
	}
}

// Edges of the GuildRole.
func (GuildRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("guild", Guild.Type).
			Ref("roles").
			Required().Unique(),
		edge.From("user_memberships", UserMembership.Type).
			Ref("roles"),
		edge.From("talent", YouTubeTalent.Type).
			Ref("roles").Unique(),
	}
}
