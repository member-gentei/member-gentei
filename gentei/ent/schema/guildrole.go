package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
			Required().Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.From("user_memberships", UserMembership.Type).
			Ref("roles"),
		edge.From("talent", YouTubeTalent.Type).
			Ref("roles").Unique(),
	}
}

func (GuildRole) Indexes() []ent.Index {
	return []ent.Index{
		// enforces one role mapping per talent
		index.Edges("guild", "talent").StorageKey("guild_talent").Unique(),
	}
}
