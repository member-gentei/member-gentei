// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
	"golang.org/x/oauth2"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// Where appends a list predicates to the UserUpdate builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetFullName sets the "full_name" field.
func (uu *UserUpdate) SetFullName(s string) *UserUpdate {
	uu.mutation.SetFullName(s)
	return uu
}

// SetAvatarHash sets the "avatar_hash" field.
func (uu *UserUpdate) SetAvatarHash(s string) *UserUpdate {
	uu.mutation.SetAvatarHash(s)
	return uu
}

// SetLastCheck sets the "last_check" field.
func (uu *UserUpdate) SetLastCheck(t time.Time) *UserUpdate {
	uu.mutation.SetLastCheck(t)
	return uu
}

// SetNillableLastCheck sets the "last_check" field if the given value is not nil.
func (uu *UserUpdate) SetNillableLastCheck(t *time.Time) *UserUpdate {
	if t != nil {
		uu.SetLastCheck(*t)
	}
	return uu
}

// SetYoutubeID sets the "youtube_id" field.
func (uu *UserUpdate) SetYoutubeID(s string) *UserUpdate {
	uu.mutation.SetYoutubeID(s)
	return uu
}

// SetNillableYoutubeID sets the "youtube_id" field if the given value is not nil.
func (uu *UserUpdate) SetNillableYoutubeID(s *string) *UserUpdate {
	if s != nil {
		uu.SetYoutubeID(*s)
	}
	return uu
}

// ClearYoutubeID clears the value of the "youtube_id" field.
func (uu *UserUpdate) ClearYoutubeID() *UserUpdate {
	uu.mutation.ClearYoutubeID()
	return uu
}

// SetYoutubeToken sets the "youtube_token" field.
func (uu *UserUpdate) SetYoutubeToken(o *oauth2.Token) *UserUpdate {
	uu.mutation.SetYoutubeToken(o)
	return uu
}

// ClearYoutubeToken clears the value of the "youtube_token" field.
func (uu *UserUpdate) ClearYoutubeToken() *UserUpdate {
	uu.mutation.ClearYoutubeToken()
	return uu
}

// SetDiscordToken sets the "discord_token" field.
func (uu *UserUpdate) SetDiscordToken(o *oauth2.Token) *UserUpdate {
	uu.mutation.SetDiscordToken(o)
	return uu
}

// ClearDiscordToken clears the value of the "discord_token" field.
func (uu *UserUpdate) ClearDiscordToken() *UserUpdate {
	uu.mutation.ClearDiscordToken()
	return uu
}

// AddGuildIDs adds the "guilds" edge to the Guild entity by IDs.
func (uu *UserUpdate) AddGuildIDs(ids ...uint64) *UserUpdate {
	uu.mutation.AddGuildIDs(ids...)
	return uu
}

// AddGuilds adds the "guilds" edges to the Guild entity.
func (uu *UserUpdate) AddGuilds(g ...*Guild) *UserUpdate {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uu.AddGuildIDs(ids...)
}

// AddGuildsAdminIDs adds the "guilds_admin" edge to the Guild entity by IDs.
func (uu *UserUpdate) AddGuildsAdminIDs(ids ...uint64) *UserUpdate {
	uu.mutation.AddGuildsAdminIDs(ids...)
	return uu
}

// AddGuildsAdmin adds the "guilds_admin" edges to the Guild entity.
func (uu *UserUpdate) AddGuildsAdmin(g ...*Guild) *UserUpdate {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uu.AddGuildsAdminIDs(ids...)
}

// AddMembershipIDs adds the "memberships" edge to the UserMembership entity by IDs.
func (uu *UserUpdate) AddMembershipIDs(ids ...int) *UserUpdate {
	uu.mutation.AddMembershipIDs(ids...)
	return uu
}

// AddMemberships adds the "memberships" edges to the UserMembership entity.
func (uu *UserUpdate) AddMemberships(u ...*UserMembership) *UserUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uu.AddMembershipIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uu *UserUpdate) Mutation() *UserMutation {
	return uu.mutation
}

// ClearGuilds clears all "guilds" edges to the Guild entity.
func (uu *UserUpdate) ClearGuilds() *UserUpdate {
	uu.mutation.ClearGuilds()
	return uu
}

// RemoveGuildIDs removes the "guilds" edge to Guild entities by IDs.
func (uu *UserUpdate) RemoveGuildIDs(ids ...uint64) *UserUpdate {
	uu.mutation.RemoveGuildIDs(ids...)
	return uu
}

// RemoveGuilds removes "guilds" edges to Guild entities.
func (uu *UserUpdate) RemoveGuilds(g ...*Guild) *UserUpdate {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uu.RemoveGuildIDs(ids...)
}

// ClearGuildsAdmin clears all "guilds_admin" edges to the Guild entity.
func (uu *UserUpdate) ClearGuildsAdmin() *UserUpdate {
	uu.mutation.ClearGuildsAdmin()
	return uu
}

// RemoveGuildsAdminIDs removes the "guilds_admin" edge to Guild entities by IDs.
func (uu *UserUpdate) RemoveGuildsAdminIDs(ids ...uint64) *UserUpdate {
	uu.mutation.RemoveGuildsAdminIDs(ids...)
	return uu
}

// RemoveGuildsAdmin removes "guilds_admin" edges to Guild entities.
func (uu *UserUpdate) RemoveGuildsAdmin(g ...*Guild) *UserUpdate {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uu.RemoveGuildsAdminIDs(ids...)
}

// ClearMemberships clears all "memberships" edges to the UserMembership entity.
func (uu *UserUpdate) ClearMemberships() *UserUpdate {
	uu.mutation.ClearMemberships()
	return uu
}

// RemoveMembershipIDs removes the "memberships" edge to UserMembership entities by IDs.
func (uu *UserUpdate) RemoveMembershipIDs(ids ...int) *UserUpdate {
	uu.mutation.RemoveMembershipIDs(ids...)
	return uu
}

// RemoveMemberships removes "memberships" edges to UserMembership entities.
func (uu *UserUpdate) RemoveMemberships(u ...*UserMembership) *UserUpdate {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uu.RemoveMembershipIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(uu.hooks) == 0 {
		if err = uu.check(); err != nil {
			return 0, err
		}
		affected, err = uu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = uu.check(); err != nil {
				return 0, err
			}
			uu.mutation = mutation
			affected, err = uu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(uu.hooks) - 1; i >= 0; i-- {
			if uu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = uu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uu *UserUpdate) check() error {
	if v, ok := uu.mutation.FullName(); ok {
		if err := user.FullNameValidator(v); err != nil {
			return &ValidationError{Name: "full_name", err: fmt.Errorf(`ent: validator failed for field "User.full_name": %w`, err)}
		}
	}
	return nil
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: user.FieldID,
			},
		},
	}
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.FullName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldFullName,
		})
	}
	if value, ok := uu.mutation.AvatarHash(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldAvatarHash,
		})
	}
	if value, ok := uu.mutation.LastCheck(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldLastCheck,
		})
	}
	if value, ok := uu.mutation.YoutubeID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldYoutubeID,
		})
	}
	if uu.mutation.YoutubeIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldYoutubeID,
		})
	}
	if value, ok := uu.mutation.YoutubeToken(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: user.FieldYoutubeToken,
		})
	}
	if uu.mutation.YoutubeTokenCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: user.FieldYoutubeToken,
		})
	}
	if value, ok := uu.mutation.DiscordToken(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: user.FieldDiscordToken,
		})
	}
	if uu.mutation.DiscordTokenCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: user.FieldDiscordToken,
		})
	}
	if uu.mutation.GuildsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsTable,
			Columns: user.GuildsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.RemovedGuildsIDs(); len(nodes) > 0 && !uu.mutation.GuildsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsTable,
			Columns: user.GuildsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.GuildsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsTable,
			Columns: user.GuildsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if uu.mutation.GuildsAdminCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsAdminTable,
			Columns: user.GuildsAdminPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.RemovedGuildsAdminIDs(); len(nodes) > 0 && !uu.mutation.GuildsAdminCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsAdminTable,
			Columns: user.GuildsAdminPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.GuildsAdminIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsAdminTable,
			Columns: user.GuildsAdminPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if uu.mutation.MembershipsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.MembershipsTable,
			Columns: []string{user.MembershipsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usermembership.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.RemovedMembershipsIDs(); len(nodes) > 0 && !uu.mutation.MembershipsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.MembershipsTable,
			Columns: []string{user.MembershipsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usermembership.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.MembershipsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.MembershipsTable,
			Columns: []string{user.MembershipsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usermembership.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UserMutation
}

// SetFullName sets the "full_name" field.
func (uuo *UserUpdateOne) SetFullName(s string) *UserUpdateOne {
	uuo.mutation.SetFullName(s)
	return uuo
}

// SetAvatarHash sets the "avatar_hash" field.
func (uuo *UserUpdateOne) SetAvatarHash(s string) *UserUpdateOne {
	uuo.mutation.SetAvatarHash(s)
	return uuo
}

// SetLastCheck sets the "last_check" field.
func (uuo *UserUpdateOne) SetLastCheck(t time.Time) *UserUpdateOne {
	uuo.mutation.SetLastCheck(t)
	return uuo
}

// SetNillableLastCheck sets the "last_check" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableLastCheck(t *time.Time) *UserUpdateOne {
	if t != nil {
		uuo.SetLastCheck(*t)
	}
	return uuo
}

// SetYoutubeID sets the "youtube_id" field.
func (uuo *UserUpdateOne) SetYoutubeID(s string) *UserUpdateOne {
	uuo.mutation.SetYoutubeID(s)
	return uuo
}

// SetNillableYoutubeID sets the "youtube_id" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableYoutubeID(s *string) *UserUpdateOne {
	if s != nil {
		uuo.SetYoutubeID(*s)
	}
	return uuo
}

// ClearYoutubeID clears the value of the "youtube_id" field.
func (uuo *UserUpdateOne) ClearYoutubeID() *UserUpdateOne {
	uuo.mutation.ClearYoutubeID()
	return uuo
}

// SetYoutubeToken sets the "youtube_token" field.
func (uuo *UserUpdateOne) SetYoutubeToken(o *oauth2.Token) *UserUpdateOne {
	uuo.mutation.SetYoutubeToken(o)
	return uuo
}

// ClearYoutubeToken clears the value of the "youtube_token" field.
func (uuo *UserUpdateOne) ClearYoutubeToken() *UserUpdateOne {
	uuo.mutation.ClearYoutubeToken()
	return uuo
}

// SetDiscordToken sets the "discord_token" field.
func (uuo *UserUpdateOne) SetDiscordToken(o *oauth2.Token) *UserUpdateOne {
	uuo.mutation.SetDiscordToken(o)
	return uuo
}

// ClearDiscordToken clears the value of the "discord_token" field.
func (uuo *UserUpdateOne) ClearDiscordToken() *UserUpdateOne {
	uuo.mutation.ClearDiscordToken()
	return uuo
}

// AddGuildIDs adds the "guilds" edge to the Guild entity by IDs.
func (uuo *UserUpdateOne) AddGuildIDs(ids ...uint64) *UserUpdateOne {
	uuo.mutation.AddGuildIDs(ids...)
	return uuo
}

// AddGuilds adds the "guilds" edges to the Guild entity.
func (uuo *UserUpdateOne) AddGuilds(g ...*Guild) *UserUpdateOne {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uuo.AddGuildIDs(ids...)
}

// AddGuildsAdminIDs adds the "guilds_admin" edge to the Guild entity by IDs.
func (uuo *UserUpdateOne) AddGuildsAdminIDs(ids ...uint64) *UserUpdateOne {
	uuo.mutation.AddGuildsAdminIDs(ids...)
	return uuo
}

// AddGuildsAdmin adds the "guilds_admin" edges to the Guild entity.
func (uuo *UserUpdateOne) AddGuildsAdmin(g ...*Guild) *UserUpdateOne {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uuo.AddGuildsAdminIDs(ids...)
}

// AddMembershipIDs adds the "memberships" edge to the UserMembership entity by IDs.
func (uuo *UserUpdateOne) AddMembershipIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.AddMembershipIDs(ids...)
	return uuo
}

// AddMemberships adds the "memberships" edges to the UserMembership entity.
func (uuo *UserUpdateOne) AddMemberships(u ...*UserMembership) *UserUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uuo.AddMembershipIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uuo *UserUpdateOne) Mutation() *UserMutation {
	return uuo.mutation
}

// ClearGuilds clears all "guilds" edges to the Guild entity.
func (uuo *UserUpdateOne) ClearGuilds() *UserUpdateOne {
	uuo.mutation.ClearGuilds()
	return uuo
}

// RemoveGuildIDs removes the "guilds" edge to Guild entities by IDs.
func (uuo *UserUpdateOne) RemoveGuildIDs(ids ...uint64) *UserUpdateOne {
	uuo.mutation.RemoveGuildIDs(ids...)
	return uuo
}

// RemoveGuilds removes "guilds" edges to Guild entities.
func (uuo *UserUpdateOne) RemoveGuilds(g ...*Guild) *UserUpdateOne {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uuo.RemoveGuildIDs(ids...)
}

// ClearGuildsAdmin clears all "guilds_admin" edges to the Guild entity.
func (uuo *UserUpdateOne) ClearGuildsAdmin() *UserUpdateOne {
	uuo.mutation.ClearGuildsAdmin()
	return uuo
}

// RemoveGuildsAdminIDs removes the "guilds_admin" edge to Guild entities by IDs.
func (uuo *UserUpdateOne) RemoveGuildsAdminIDs(ids ...uint64) *UserUpdateOne {
	uuo.mutation.RemoveGuildsAdminIDs(ids...)
	return uuo
}

// RemoveGuildsAdmin removes "guilds_admin" edges to Guild entities.
func (uuo *UserUpdateOne) RemoveGuildsAdmin(g ...*Guild) *UserUpdateOne {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return uuo.RemoveGuildsAdminIDs(ids...)
}

// ClearMemberships clears all "memberships" edges to the UserMembership entity.
func (uuo *UserUpdateOne) ClearMemberships() *UserUpdateOne {
	uuo.mutation.ClearMemberships()
	return uuo
}

// RemoveMembershipIDs removes the "memberships" edge to UserMembership entities by IDs.
func (uuo *UserUpdateOne) RemoveMembershipIDs(ids ...int) *UserUpdateOne {
	uuo.mutation.RemoveMembershipIDs(ids...)
	return uuo
}

// RemoveMemberships removes "memberships" edges to UserMembership entities.
func (uuo *UserUpdateOne) RemoveMemberships(u ...*UserMembership) *UserUpdateOne {
	ids := make([]int, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return uuo.RemoveMembershipIDs(ids...)
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UserUpdateOne) Select(field string, fields ...string) *UserUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated User entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	var (
		err  error
		node *User
	)
	if len(uuo.hooks) == 0 {
		if err = uuo.check(); err != nil {
			return nil, err
		}
		node, err = uuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*UserMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = uuo.check(); err != nil {
				return nil, err
			}
			uuo.mutation = mutation
			node, err = uuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(uuo.hooks) - 1; i >= 0; i-- {
			if uuo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = uuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, uuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (uuo *UserUpdateOne) check() error {
	if v, ok := uuo.mutation.FullName(); ok {
		if err := user.FullNameValidator(v); err != nil {
			return &ValidationError{Name: "full_name", err: fmt.Errorf(`ent: validator failed for field "User.full_name": %w`, err)}
		}
	}
	return nil
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (_node *User, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   user.Table,
			Columns: user.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: user.FieldID,
			},
		},
	}
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "User.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, user.FieldID)
		for _, f := range fields {
			if !user.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != user.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.FullName(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldFullName,
		})
	}
	if value, ok := uuo.mutation.AvatarHash(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldAvatarHash,
		})
	}
	if value, ok := uuo.mutation.LastCheck(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: user.FieldLastCheck,
		})
	}
	if value, ok := uuo.mutation.YoutubeID(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: user.FieldYoutubeID,
		})
	}
	if uuo.mutation.YoutubeIDCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Column: user.FieldYoutubeID,
		})
	}
	if value, ok := uuo.mutation.YoutubeToken(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: user.FieldYoutubeToken,
		})
	}
	if uuo.mutation.YoutubeTokenCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: user.FieldYoutubeToken,
		})
	}
	if value, ok := uuo.mutation.DiscordToken(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Value:  value,
			Column: user.FieldDiscordToken,
		})
	}
	if uuo.mutation.DiscordTokenCleared() {
		_spec.Fields.Clear = append(_spec.Fields.Clear, &sqlgraph.FieldSpec{
			Type:   field.TypeJSON,
			Column: user.FieldDiscordToken,
		})
	}
	if uuo.mutation.GuildsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsTable,
			Columns: user.GuildsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.RemovedGuildsIDs(); len(nodes) > 0 && !uuo.mutation.GuildsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsTable,
			Columns: user.GuildsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.GuildsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsTable,
			Columns: user.GuildsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if uuo.mutation.GuildsAdminCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsAdminTable,
			Columns: user.GuildsAdminPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.RemovedGuildsAdminIDs(); len(nodes) > 0 && !uuo.mutation.GuildsAdminCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsAdminTable,
			Columns: user.GuildsAdminPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.GuildsAdminIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   user.GuildsAdminTable,
			Columns: user.GuildsAdminPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUint64,
					Column: guild.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if uuo.mutation.MembershipsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.MembershipsTable,
			Columns: []string{user.MembershipsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usermembership.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.RemovedMembershipsIDs(); len(nodes) > 0 && !uuo.mutation.MembershipsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.MembershipsTable,
			Columns: []string{user.MembershipsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usermembership.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.MembershipsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.MembershipsTable,
			Columns: []string{user.MembershipsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: usermembership.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &User{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}
