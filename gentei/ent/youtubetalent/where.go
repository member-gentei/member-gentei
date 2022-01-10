// Code generated by entc, DO NOT EDIT.

package youtubetalent

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// ChannelName applies equality check predicate on the "channel_name" field. It's identical to ChannelNameEQ.
func ChannelName(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldChannelName), v))
	})
}

// ThumbnailURL applies equality check predicate on the "thumbnail_url" field. It's identical to ThumbnailURLEQ.
func ThumbnailURL(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldThumbnailURL), v))
	})
}

// LastUpdated applies equality check predicate on the "last_updated" field. It's identical to LastUpdatedEQ.
func LastUpdated(v time.Time) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLastUpdated), v))
	})
}

// ChannelNameEQ applies the EQ predicate on the "channel_name" field.
func ChannelNameEQ(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldChannelName), v))
	})
}

// ChannelNameNEQ applies the NEQ predicate on the "channel_name" field.
func ChannelNameNEQ(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldChannelName), v))
	})
}

// ChannelNameIn applies the In predicate on the "channel_name" field.
func ChannelNameIn(vs ...string) predicate.YouTubeTalent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldChannelName), v...))
	})
}

// ChannelNameNotIn applies the NotIn predicate on the "channel_name" field.
func ChannelNameNotIn(vs ...string) predicate.YouTubeTalent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldChannelName), v...))
	})
}

// ChannelNameGT applies the GT predicate on the "channel_name" field.
func ChannelNameGT(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldChannelName), v))
	})
}

// ChannelNameGTE applies the GTE predicate on the "channel_name" field.
func ChannelNameGTE(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldChannelName), v))
	})
}

// ChannelNameLT applies the LT predicate on the "channel_name" field.
func ChannelNameLT(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldChannelName), v))
	})
}

// ChannelNameLTE applies the LTE predicate on the "channel_name" field.
func ChannelNameLTE(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldChannelName), v))
	})
}

// ChannelNameContains applies the Contains predicate on the "channel_name" field.
func ChannelNameContains(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldChannelName), v))
	})
}

// ChannelNameHasPrefix applies the HasPrefix predicate on the "channel_name" field.
func ChannelNameHasPrefix(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldChannelName), v))
	})
}

// ChannelNameHasSuffix applies the HasSuffix predicate on the "channel_name" field.
func ChannelNameHasSuffix(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldChannelName), v))
	})
}

// ChannelNameEqualFold applies the EqualFold predicate on the "channel_name" field.
func ChannelNameEqualFold(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldChannelName), v))
	})
}

// ChannelNameContainsFold applies the ContainsFold predicate on the "channel_name" field.
func ChannelNameContainsFold(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldChannelName), v))
	})
}

// ThumbnailURLEQ applies the EQ predicate on the "thumbnail_url" field.
func ThumbnailURLEQ(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLNEQ applies the NEQ predicate on the "thumbnail_url" field.
func ThumbnailURLNEQ(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLIn applies the In predicate on the "thumbnail_url" field.
func ThumbnailURLIn(vs ...string) predicate.YouTubeTalent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldThumbnailURL), v...))
	})
}

// ThumbnailURLNotIn applies the NotIn predicate on the "thumbnail_url" field.
func ThumbnailURLNotIn(vs ...string) predicate.YouTubeTalent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldThumbnailURL), v...))
	})
}

// ThumbnailURLGT applies the GT predicate on the "thumbnail_url" field.
func ThumbnailURLGT(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLGTE applies the GTE predicate on the "thumbnail_url" field.
func ThumbnailURLGTE(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLLT applies the LT predicate on the "thumbnail_url" field.
func ThumbnailURLLT(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLLTE applies the LTE predicate on the "thumbnail_url" field.
func ThumbnailURLLTE(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLContains applies the Contains predicate on the "thumbnail_url" field.
func ThumbnailURLContains(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLHasPrefix applies the HasPrefix predicate on the "thumbnail_url" field.
func ThumbnailURLHasPrefix(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLHasSuffix applies the HasSuffix predicate on the "thumbnail_url" field.
func ThumbnailURLHasSuffix(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLEqualFold applies the EqualFold predicate on the "thumbnail_url" field.
func ThumbnailURLEqualFold(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldThumbnailURL), v))
	})
}

// ThumbnailURLContainsFold applies the ContainsFold predicate on the "thumbnail_url" field.
func ThumbnailURLContainsFold(v string) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldThumbnailURL), v))
	})
}

// LastUpdatedEQ applies the EQ predicate on the "last_updated" field.
func LastUpdatedEQ(v time.Time) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldLastUpdated), v))
	})
}

// LastUpdatedNEQ applies the NEQ predicate on the "last_updated" field.
func LastUpdatedNEQ(v time.Time) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldLastUpdated), v))
	})
}

// LastUpdatedIn applies the In predicate on the "last_updated" field.
func LastUpdatedIn(vs ...time.Time) predicate.YouTubeTalent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldLastUpdated), v...))
	})
}

// LastUpdatedNotIn applies the NotIn predicate on the "last_updated" field.
func LastUpdatedNotIn(vs ...time.Time) predicate.YouTubeTalent {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldLastUpdated), v...))
	})
}

// LastUpdatedGT applies the GT predicate on the "last_updated" field.
func LastUpdatedGT(v time.Time) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldLastUpdated), v))
	})
}

// LastUpdatedGTE applies the GTE predicate on the "last_updated" field.
func LastUpdatedGTE(v time.Time) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldLastUpdated), v))
	})
}

// LastUpdatedLT applies the LT predicate on the "last_updated" field.
func LastUpdatedLT(v time.Time) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldLastUpdated), v))
	})
}

// LastUpdatedLTE applies the LTE predicate on the "last_updated" field.
func LastUpdatedLTE(v time.Time) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldLastUpdated), v))
	})
}

// HasGuilds applies the HasEdge predicate on the "guilds" edge.
func HasGuilds() predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(GuildsTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, GuildsTable, GuildsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGuildsWith applies the HasEdge predicate on the "guilds" edge with a given conditions (other predicates).
func HasGuildsWith(preds ...predicate.Guild) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(GuildsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, GuildsTable, GuildsPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.YouTubeTalent) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.YouTubeTalent) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.YouTubeTalent) predicate.YouTubeTalent {
	return predicate.YouTubeTalent(func(s *sql.Selector) {
		p(s.Not())
	})
}