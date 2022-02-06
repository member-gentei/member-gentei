package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserMembership holds the schema definition for the UserMembership entity.
type UserMembership struct {
	ent.Schema
}

// Fields of the UserMembership.
func (UserMembership) Fields() []ent.Field {
	return []ent.Field{
		field.Time("first_failed").Optional(),
		field.Time("last_verified"),
		field.Int("fail_count").Default(0),
	}
}

// Edges of the UserMembership.
func (UserMembership) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("memberships").
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("youtube_talent", YouTubeTalent.Type).
			Required().
			Unique().
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
		edge.To("roles", GuildRole.Type).
			Annotations(entsql.Annotation{
				OnDelete: entsql.Cascade,
			}),
	}
}
