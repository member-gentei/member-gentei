// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uint64) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uint64) predicate.User {
	return predicate.User(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uint64) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uint64) predicate.User {
	return predicate.User(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uint64) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uint64) predicate.User {
	return predicate.User(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uint64) predicate.User {
	return predicate.User(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uint64) predicate.User {
	return predicate.User(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uint64) predicate.User {
	return predicate.User(sql.FieldLTE(FieldID, id))
}

// FullName applies equality check predicate on the "full_name" field. It's identical to FullNameEQ.
func FullName(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldFullName, v))
}

// AvatarHash applies equality check predicate on the "avatar_hash" field. It's identical to AvatarHashEQ.
func AvatarHash(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldAvatarHash, v))
}

// LastCheck applies equality check predicate on the "last_check" field. It's identical to LastCheckEQ.
func LastCheck(v time.Time) predicate.User {
	return predicate.User(sql.FieldEQ(FieldLastCheck, v))
}

// YoutubeID applies equality check predicate on the "youtube_id" field. It's identical to YoutubeIDEQ.
func YoutubeID(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldYoutubeID, v))
}

// FullNameEQ applies the EQ predicate on the "full_name" field.
func FullNameEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldFullName, v))
}

// FullNameNEQ applies the NEQ predicate on the "full_name" field.
func FullNameNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldFullName, v))
}

// FullNameIn applies the In predicate on the "full_name" field.
func FullNameIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldFullName, vs...))
}

// FullNameNotIn applies the NotIn predicate on the "full_name" field.
func FullNameNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldFullName, vs...))
}

// FullNameGT applies the GT predicate on the "full_name" field.
func FullNameGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldFullName, v))
}

// FullNameGTE applies the GTE predicate on the "full_name" field.
func FullNameGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldFullName, v))
}

// FullNameLT applies the LT predicate on the "full_name" field.
func FullNameLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldFullName, v))
}

// FullNameLTE applies the LTE predicate on the "full_name" field.
func FullNameLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldFullName, v))
}

// FullNameContains applies the Contains predicate on the "full_name" field.
func FullNameContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldFullName, v))
}

// FullNameHasPrefix applies the HasPrefix predicate on the "full_name" field.
func FullNameHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldFullName, v))
}

// FullNameHasSuffix applies the HasSuffix predicate on the "full_name" field.
func FullNameHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldFullName, v))
}

// FullNameEqualFold applies the EqualFold predicate on the "full_name" field.
func FullNameEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldFullName, v))
}

// FullNameContainsFold applies the ContainsFold predicate on the "full_name" field.
func FullNameContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldFullName, v))
}

// AvatarHashEQ applies the EQ predicate on the "avatar_hash" field.
func AvatarHashEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldAvatarHash, v))
}

// AvatarHashNEQ applies the NEQ predicate on the "avatar_hash" field.
func AvatarHashNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldAvatarHash, v))
}

// AvatarHashIn applies the In predicate on the "avatar_hash" field.
func AvatarHashIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldAvatarHash, vs...))
}

// AvatarHashNotIn applies the NotIn predicate on the "avatar_hash" field.
func AvatarHashNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldAvatarHash, vs...))
}

// AvatarHashGT applies the GT predicate on the "avatar_hash" field.
func AvatarHashGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldAvatarHash, v))
}

// AvatarHashGTE applies the GTE predicate on the "avatar_hash" field.
func AvatarHashGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldAvatarHash, v))
}

// AvatarHashLT applies the LT predicate on the "avatar_hash" field.
func AvatarHashLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldAvatarHash, v))
}

// AvatarHashLTE applies the LTE predicate on the "avatar_hash" field.
func AvatarHashLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldAvatarHash, v))
}

// AvatarHashContains applies the Contains predicate on the "avatar_hash" field.
func AvatarHashContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldAvatarHash, v))
}

// AvatarHashHasPrefix applies the HasPrefix predicate on the "avatar_hash" field.
func AvatarHashHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldAvatarHash, v))
}

// AvatarHashHasSuffix applies the HasSuffix predicate on the "avatar_hash" field.
func AvatarHashHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldAvatarHash, v))
}

// AvatarHashEqualFold applies the EqualFold predicate on the "avatar_hash" field.
func AvatarHashEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldAvatarHash, v))
}

// AvatarHashContainsFold applies the ContainsFold predicate on the "avatar_hash" field.
func AvatarHashContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldAvatarHash, v))
}

// LastCheckEQ applies the EQ predicate on the "last_check" field.
func LastCheckEQ(v time.Time) predicate.User {
	return predicate.User(sql.FieldEQ(FieldLastCheck, v))
}

// LastCheckNEQ applies the NEQ predicate on the "last_check" field.
func LastCheckNEQ(v time.Time) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldLastCheck, v))
}

// LastCheckIn applies the In predicate on the "last_check" field.
func LastCheckIn(vs ...time.Time) predicate.User {
	return predicate.User(sql.FieldIn(FieldLastCheck, vs...))
}

// LastCheckNotIn applies the NotIn predicate on the "last_check" field.
func LastCheckNotIn(vs ...time.Time) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldLastCheck, vs...))
}

// LastCheckGT applies the GT predicate on the "last_check" field.
func LastCheckGT(v time.Time) predicate.User {
	return predicate.User(sql.FieldGT(FieldLastCheck, v))
}

// LastCheckGTE applies the GTE predicate on the "last_check" field.
func LastCheckGTE(v time.Time) predicate.User {
	return predicate.User(sql.FieldGTE(FieldLastCheck, v))
}

// LastCheckLT applies the LT predicate on the "last_check" field.
func LastCheckLT(v time.Time) predicate.User {
	return predicate.User(sql.FieldLT(FieldLastCheck, v))
}

// LastCheckLTE applies the LTE predicate on the "last_check" field.
func LastCheckLTE(v time.Time) predicate.User {
	return predicate.User(sql.FieldLTE(FieldLastCheck, v))
}

// YoutubeIDEQ applies the EQ predicate on the "youtube_id" field.
func YoutubeIDEQ(v string) predicate.User {
	return predicate.User(sql.FieldEQ(FieldYoutubeID, v))
}

// YoutubeIDNEQ applies the NEQ predicate on the "youtube_id" field.
func YoutubeIDNEQ(v string) predicate.User {
	return predicate.User(sql.FieldNEQ(FieldYoutubeID, v))
}

// YoutubeIDIn applies the In predicate on the "youtube_id" field.
func YoutubeIDIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldIn(FieldYoutubeID, vs...))
}

// YoutubeIDNotIn applies the NotIn predicate on the "youtube_id" field.
func YoutubeIDNotIn(vs ...string) predicate.User {
	return predicate.User(sql.FieldNotIn(FieldYoutubeID, vs...))
}

// YoutubeIDGT applies the GT predicate on the "youtube_id" field.
func YoutubeIDGT(v string) predicate.User {
	return predicate.User(sql.FieldGT(FieldYoutubeID, v))
}

// YoutubeIDGTE applies the GTE predicate on the "youtube_id" field.
func YoutubeIDGTE(v string) predicate.User {
	return predicate.User(sql.FieldGTE(FieldYoutubeID, v))
}

// YoutubeIDLT applies the LT predicate on the "youtube_id" field.
func YoutubeIDLT(v string) predicate.User {
	return predicate.User(sql.FieldLT(FieldYoutubeID, v))
}

// YoutubeIDLTE applies the LTE predicate on the "youtube_id" field.
func YoutubeIDLTE(v string) predicate.User {
	return predicate.User(sql.FieldLTE(FieldYoutubeID, v))
}

// YoutubeIDContains applies the Contains predicate on the "youtube_id" field.
func YoutubeIDContains(v string) predicate.User {
	return predicate.User(sql.FieldContains(FieldYoutubeID, v))
}

// YoutubeIDHasPrefix applies the HasPrefix predicate on the "youtube_id" field.
func YoutubeIDHasPrefix(v string) predicate.User {
	return predicate.User(sql.FieldHasPrefix(FieldYoutubeID, v))
}

// YoutubeIDHasSuffix applies the HasSuffix predicate on the "youtube_id" field.
func YoutubeIDHasSuffix(v string) predicate.User {
	return predicate.User(sql.FieldHasSuffix(FieldYoutubeID, v))
}

// YoutubeIDIsNil applies the IsNil predicate on the "youtube_id" field.
func YoutubeIDIsNil() predicate.User {
	return predicate.User(sql.FieldIsNull(FieldYoutubeID))
}

// YoutubeIDNotNil applies the NotNil predicate on the "youtube_id" field.
func YoutubeIDNotNil() predicate.User {
	return predicate.User(sql.FieldNotNull(FieldYoutubeID))
}

// YoutubeIDEqualFold applies the EqualFold predicate on the "youtube_id" field.
func YoutubeIDEqualFold(v string) predicate.User {
	return predicate.User(sql.FieldEqualFold(FieldYoutubeID, v))
}

// YoutubeIDContainsFold applies the ContainsFold predicate on the "youtube_id" field.
func YoutubeIDContainsFold(v string) predicate.User {
	return predicate.User(sql.FieldContainsFold(FieldYoutubeID, v))
}

// YoutubeTokenIsNil applies the IsNil predicate on the "youtube_token" field.
func YoutubeTokenIsNil() predicate.User {
	return predicate.User(sql.FieldIsNull(FieldYoutubeToken))
}

// YoutubeTokenNotNil applies the NotNil predicate on the "youtube_token" field.
func YoutubeTokenNotNil() predicate.User {
	return predicate.User(sql.FieldNotNull(FieldYoutubeToken))
}

// DiscordTokenIsNil applies the IsNil predicate on the "discord_token" field.
func DiscordTokenIsNil() predicate.User {
	return predicate.User(sql.FieldIsNull(FieldDiscordToken))
}

// DiscordTokenNotNil applies the NotNil predicate on the "discord_token" field.
func DiscordTokenNotNil() predicate.User {
	return predicate.User(sql.FieldNotNull(FieldDiscordToken))
}

// HasGuilds applies the HasEdge predicate on the "guilds" edge.
func HasGuilds() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, GuildsTable, GuildsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGuildsWith applies the HasEdge predicate on the "guilds" edge with a given conditions (other predicates).
func HasGuildsWith(preds ...predicate.Guild) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(GuildsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, GuildsTable, GuildsPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasGuildsAdmin applies the HasEdge predicate on the "guilds_admin" edge.
func HasGuildsAdmin() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, GuildsAdminTable, GuildsAdminPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGuildsAdminWith applies the HasEdge predicate on the "guilds_admin" edge with a given conditions (other predicates).
func HasGuildsAdminWith(preds ...predicate.Guild) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(GuildsAdminInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, GuildsAdminTable, GuildsAdminPrimaryKey...),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasMemberships applies the HasEdge predicate on the "memberships" edge.
func HasMemberships() predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MembershipsTable, MembershipsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMembershipsWith applies the HasEdge predicate on the "memberships" edge with a given conditions (other predicates).
func HasMembershipsWith(preds ...predicate.UserMembership) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.To(MembershipsInverseTable, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MembershipsTable, MembershipsColumn),
		)
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.User) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.User) predicate.User {
	return predicate.User(func(s *sql.Selector) {
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
func Not(p predicate.User) predicate.User {
	return predicate.User(func(s *sql.Selector) {
		p(s.Not())
	})
}
