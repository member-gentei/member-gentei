// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/member-gentei/member-gentei/gentei/ent/guild"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/schema"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
)

// GuildUpdate is the builder for updating Guild entities.
type GuildUpdate struct {
	config
	hooks    []Hook
	mutation *GuildMutation
}

// Where appends a list predicates to the GuildUpdate builder.
func (gu *GuildUpdate) Where(ps ...predicate.Guild) *GuildUpdate {
	gu.mutation.Where(ps...)
	return gu
}

// SetName sets the "name" field.
func (gu *GuildUpdate) SetName(s string) *GuildUpdate {
	gu.mutation.SetName(s)
	return gu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (gu *GuildUpdate) SetNillableName(s *string) *GuildUpdate {
	if s != nil {
		gu.SetName(*s)
	}
	return gu
}

// SetIconHash sets the "icon_hash" field.
func (gu *GuildUpdate) SetIconHash(s string) *GuildUpdate {
	gu.mutation.SetIconHash(s)
	return gu
}

// SetNillableIconHash sets the "icon_hash" field if the given value is not nil.
func (gu *GuildUpdate) SetNillableIconHash(s *string) *GuildUpdate {
	if s != nil {
		gu.SetIconHash(*s)
	}
	return gu
}

// ClearIconHash clears the value of the "icon_hash" field.
func (gu *GuildUpdate) ClearIconHash() *GuildUpdate {
	gu.mutation.ClearIconHash()
	return gu
}

// SetAuditChannel sets the "audit_channel" field.
func (gu *GuildUpdate) SetAuditChannel(u uint64) *GuildUpdate {
	gu.mutation.ResetAuditChannel()
	gu.mutation.SetAuditChannel(u)
	return gu
}

// SetNillableAuditChannel sets the "audit_channel" field if the given value is not nil.
func (gu *GuildUpdate) SetNillableAuditChannel(u *uint64) *GuildUpdate {
	if u != nil {
		gu.SetAuditChannel(*u)
	}
	return gu
}

// AddAuditChannel adds u to the "audit_channel" field.
func (gu *GuildUpdate) AddAuditChannel(u int64) *GuildUpdate {
	gu.mutation.AddAuditChannel(u)
	return gu
}

// ClearAuditChannel clears the value of the "audit_channel" field.
func (gu *GuildUpdate) ClearAuditChannel() *GuildUpdate {
	gu.mutation.ClearAuditChannel()
	return gu
}

// SetLanguage sets the "language" field.
func (gu *GuildUpdate) SetLanguage(value guild.Language) *GuildUpdate {
	gu.mutation.SetLanguage(value)
	return gu
}

// SetNillableLanguage sets the "language" field if the given value is not nil.
func (gu *GuildUpdate) SetNillableLanguage(value *guild.Language) *GuildUpdate {
	if value != nil {
		gu.SetLanguage(*value)
	}
	return gu
}

// SetSettings sets the "settings" field.
func (gu *GuildUpdate) SetSettings(ss *schema.GuildSettings) *GuildUpdate {
	gu.mutation.SetSettings(ss)
	return gu
}

// ClearSettings clears the value of the "settings" field.
func (gu *GuildUpdate) ClearSettings() *GuildUpdate {
	gu.mutation.ClearSettings()
	return gu
}

// AddMemberIDs adds the "members" edge to the User entity by IDs.
func (gu *GuildUpdate) AddMemberIDs(ids ...uint64) *GuildUpdate {
	gu.mutation.AddMemberIDs(ids...)
	return gu
}

// AddMembers adds the "members" edges to the User entity.
func (gu *GuildUpdate) AddMembers(u ...*User) *GuildUpdate {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return gu.AddMemberIDs(ids...)
}

// AddAdminIDs adds the "admins" edge to the User entity by IDs.
func (gu *GuildUpdate) AddAdminIDs(ids ...uint64) *GuildUpdate {
	gu.mutation.AddAdminIDs(ids...)
	return gu
}

// AddAdmins adds the "admins" edges to the User entity.
func (gu *GuildUpdate) AddAdmins(u ...*User) *GuildUpdate {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return gu.AddAdminIDs(ids...)
}

// AddRoleIDs adds the "roles" edge to the GuildRole entity by IDs.
func (gu *GuildUpdate) AddRoleIDs(ids ...uint64) *GuildUpdate {
	gu.mutation.AddRoleIDs(ids...)
	return gu
}

// AddRoles adds the "roles" edges to the GuildRole entity.
func (gu *GuildUpdate) AddRoles(g ...*GuildRole) *GuildUpdate {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return gu.AddRoleIDs(ids...)
}

// AddYoutubeTalentIDs adds the "youtube_talents" edge to the YouTubeTalent entity by IDs.
func (gu *GuildUpdate) AddYoutubeTalentIDs(ids ...string) *GuildUpdate {
	gu.mutation.AddYoutubeTalentIDs(ids...)
	return gu
}

// AddYoutubeTalents adds the "youtube_talents" edges to the YouTubeTalent entity.
func (gu *GuildUpdate) AddYoutubeTalents(y ...*YouTubeTalent) *GuildUpdate {
	ids := make([]string, len(y))
	for i := range y {
		ids[i] = y[i].ID
	}
	return gu.AddYoutubeTalentIDs(ids...)
}

// Mutation returns the GuildMutation object of the builder.
func (gu *GuildUpdate) Mutation() *GuildMutation {
	return gu.mutation
}

// ClearMembers clears all "members" edges to the User entity.
func (gu *GuildUpdate) ClearMembers() *GuildUpdate {
	gu.mutation.ClearMembers()
	return gu
}

// RemoveMemberIDs removes the "members" edge to User entities by IDs.
func (gu *GuildUpdate) RemoveMemberIDs(ids ...uint64) *GuildUpdate {
	gu.mutation.RemoveMemberIDs(ids...)
	return gu
}

// RemoveMembers removes "members" edges to User entities.
func (gu *GuildUpdate) RemoveMembers(u ...*User) *GuildUpdate {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return gu.RemoveMemberIDs(ids...)
}

// ClearAdmins clears all "admins" edges to the User entity.
func (gu *GuildUpdate) ClearAdmins() *GuildUpdate {
	gu.mutation.ClearAdmins()
	return gu
}

// RemoveAdminIDs removes the "admins" edge to User entities by IDs.
func (gu *GuildUpdate) RemoveAdminIDs(ids ...uint64) *GuildUpdate {
	gu.mutation.RemoveAdminIDs(ids...)
	return gu
}

// RemoveAdmins removes "admins" edges to User entities.
func (gu *GuildUpdate) RemoveAdmins(u ...*User) *GuildUpdate {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return gu.RemoveAdminIDs(ids...)
}

// ClearRoles clears all "roles" edges to the GuildRole entity.
func (gu *GuildUpdate) ClearRoles() *GuildUpdate {
	gu.mutation.ClearRoles()
	return gu
}

// RemoveRoleIDs removes the "roles" edge to GuildRole entities by IDs.
func (gu *GuildUpdate) RemoveRoleIDs(ids ...uint64) *GuildUpdate {
	gu.mutation.RemoveRoleIDs(ids...)
	return gu
}

// RemoveRoles removes "roles" edges to GuildRole entities.
func (gu *GuildUpdate) RemoveRoles(g ...*GuildRole) *GuildUpdate {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return gu.RemoveRoleIDs(ids...)
}

// ClearYoutubeTalents clears all "youtube_talents" edges to the YouTubeTalent entity.
func (gu *GuildUpdate) ClearYoutubeTalents() *GuildUpdate {
	gu.mutation.ClearYoutubeTalents()
	return gu
}

// RemoveYoutubeTalentIDs removes the "youtube_talents" edge to YouTubeTalent entities by IDs.
func (gu *GuildUpdate) RemoveYoutubeTalentIDs(ids ...string) *GuildUpdate {
	gu.mutation.RemoveYoutubeTalentIDs(ids...)
	return gu
}

// RemoveYoutubeTalents removes "youtube_talents" edges to YouTubeTalent entities.
func (gu *GuildUpdate) RemoveYoutubeTalents(y ...*YouTubeTalent) *GuildUpdate {
	ids := make([]string, len(y))
	for i := range y {
		ids[i] = y[i].ID
	}
	return gu.RemoveYoutubeTalentIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (gu *GuildUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, gu.sqlSave, gu.mutation, gu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gu *GuildUpdate) SaveX(ctx context.Context) int {
	affected, err := gu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (gu *GuildUpdate) Exec(ctx context.Context) error {
	_, err := gu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gu *GuildUpdate) ExecX(ctx context.Context) {
	if err := gu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (gu *GuildUpdate) check() error {
	if v, ok := gu.mutation.Language(); ok {
		if err := guild.LanguageValidator(v); err != nil {
			return &ValidationError{Name: "language", err: fmt.Errorf(`ent: validator failed for field "Guild.language": %w`, err)}
		}
	}
	return nil
}

func (gu *GuildUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := gu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(guild.Table, guild.Columns, sqlgraph.NewFieldSpec(guild.FieldID, field.TypeUint64))
	if ps := gu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gu.mutation.Name(); ok {
		_spec.SetField(guild.FieldName, field.TypeString, value)
	}
	if value, ok := gu.mutation.IconHash(); ok {
		_spec.SetField(guild.FieldIconHash, field.TypeString, value)
	}
	if gu.mutation.IconHashCleared() {
		_spec.ClearField(guild.FieldIconHash, field.TypeString)
	}
	if value, ok := gu.mutation.AuditChannel(); ok {
		_spec.SetField(guild.FieldAuditChannel, field.TypeUint64, value)
	}
	if value, ok := gu.mutation.AddedAuditChannel(); ok {
		_spec.AddField(guild.FieldAuditChannel, field.TypeUint64, value)
	}
	if gu.mutation.AuditChannelCleared() {
		_spec.ClearField(guild.FieldAuditChannel, field.TypeUint64)
	}
	if value, ok := gu.mutation.Language(); ok {
		_spec.SetField(guild.FieldLanguage, field.TypeEnum, value)
	}
	if value, ok := gu.mutation.Settings(); ok {
		_spec.SetField(guild.FieldSettings, field.TypeJSON, value)
	}
	if gu.mutation.SettingsCleared() {
		_spec.ClearField(guild.FieldSettings, field.TypeJSON)
	}
	if gu.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.MembersTable,
			Columns: guild.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gu.mutation.RemovedMembersIDs(); len(nodes) > 0 && !gu.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.MembersTable,
			Columns: guild.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gu.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.MembersTable,
			Columns: guild.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if gu.mutation.AdminsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.AdminsTable,
			Columns: guild.AdminsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gu.mutation.RemovedAdminsIDs(); len(nodes) > 0 && !gu.mutation.AdminsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.AdminsTable,
			Columns: guild.AdminsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gu.mutation.AdminsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.AdminsTable,
			Columns: guild.AdminsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if gu.mutation.RolesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   guild.RolesTable,
			Columns: []string{guild.RolesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guildrole.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gu.mutation.RemovedRolesIDs(); len(nodes) > 0 && !gu.mutation.RolesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   guild.RolesTable,
			Columns: []string{guild.RolesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guildrole.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gu.mutation.RolesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   guild.RolesTable,
			Columns: []string{guild.RolesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guildrole.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if gu.mutation.YoutubeTalentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   guild.YoutubeTalentsTable,
			Columns: guild.YoutubeTalentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(youtubetalent.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gu.mutation.RemovedYoutubeTalentsIDs(); len(nodes) > 0 && !gu.mutation.YoutubeTalentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   guild.YoutubeTalentsTable,
			Columns: guild.YoutubeTalentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(youtubetalent.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gu.mutation.YoutubeTalentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   guild.YoutubeTalentsTable,
			Columns: guild.YoutubeTalentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(youtubetalent.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, gu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{guild.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	gu.mutation.done = true
	return n, nil
}

// GuildUpdateOne is the builder for updating a single Guild entity.
type GuildUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *GuildMutation
}

// SetName sets the "name" field.
func (guo *GuildUpdateOne) SetName(s string) *GuildUpdateOne {
	guo.mutation.SetName(s)
	return guo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (guo *GuildUpdateOne) SetNillableName(s *string) *GuildUpdateOne {
	if s != nil {
		guo.SetName(*s)
	}
	return guo
}

// SetIconHash sets the "icon_hash" field.
func (guo *GuildUpdateOne) SetIconHash(s string) *GuildUpdateOne {
	guo.mutation.SetIconHash(s)
	return guo
}

// SetNillableIconHash sets the "icon_hash" field if the given value is not nil.
func (guo *GuildUpdateOne) SetNillableIconHash(s *string) *GuildUpdateOne {
	if s != nil {
		guo.SetIconHash(*s)
	}
	return guo
}

// ClearIconHash clears the value of the "icon_hash" field.
func (guo *GuildUpdateOne) ClearIconHash() *GuildUpdateOne {
	guo.mutation.ClearIconHash()
	return guo
}

// SetAuditChannel sets the "audit_channel" field.
func (guo *GuildUpdateOne) SetAuditChannel(u uint64) *GuildUpdateOne {
	guo.mutation.ResetAuditChannel()
	guo.mutation.SetAuditChannel(u)
	return guo
}

// SetNillableAuditChannel sets the "audit_channel" field if the given value is not nil.
func (guo *GuildUpdateOne) SetNillableAuditChannel(u *uint64) *GuildUpdateOne {
	if u != nil {
		guo.SetAuditChannel(*u)
	}
	return guo
}

// AddAuditChannel adds u to the "audit_channel" field.
func (guo *GuildUpdateOne) AddAuditChannel(u int64) *GuildUpdateOne {
	guo.mutation.AddAuditChannel(u)
	return guo
}

// ClearAuditChannel clears the value of the "audit_channel" field.
func (guo *GuildUpdateOne) ClearAuditChannel() *GuildUpdateOne {
	guo.mutation.ClearAuditChannel()
	return guo
}

// SetLanguage sets the "language" field.
func (guo *GuildUpdateOne) SetLanguage(gu guild.Language) *GuildUpdateOne {
	guo.mutation.SetLanguage(gu)
	return guo
}

// SetNillableLanguage sets the "language" field if the given value is not nil.
func (guo *GuildUpdateOne) SetNillableLanguage(gu *guild.Language) *GuildUpdateOne {
	if gu != nil {
		guo.SetLanguage(*gu)
	}
	return guo
}

// SetSettings sets the "settings" field.
func (guo *GuildUpdateOne) SetSettings(ss *schema.GuildSettings) *GuildUpdateOne {
	guo.mutation.SetSettings(ss)
	return guo
}

// ClearSettings clears the value of the "settings" field.
func (guo *GuildUpdateOne) ClearSettings() *GuildUpdateOne {
	guo.mutation.ClearSettings()
	return guo
}

// AddMemberIDs adds the "members" edge to the User entity by IDs.
func (guo *GuildUpdateOne) AddMemberIDs(ids ...uint64) *GuildUpdateOne {
	guo.mutation.AddMemberIDs(ids...)
	return guo
}

// AddMembers adds the "members" edges to the User entity.
func (guo *GuildUpdateOne) AddMembers(u ...*User) *GuildUpdateOne {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return guo.AddMemberIDs(ids...)
}

// AddAdminIDs adds the "admins" edge to the User entity by IDs.
func (guo *GuildUpdateOne) AddAdminIDs(ids ...uint64) *GuildUpdateOne {
	guo.mutation.AddAdminIDs(ids...)
	return guo
}

// AddAdmins adds the "admins" edges to the User entity.
func (guo *GuildUpdateOne) AddAdmins(u ...*User) *GuildUpdateOne {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return guo.AddAdminIDs(ids...)
}

// AddRoleIDs adds the "roles" edge to the GuildRole entity by IDs.
func (guo *GuildUpdateOne) AddRoleIDs(ids ...uint64) *GuildUpdateOne {
	guo.mutation.AddRoleIDs(ids...)
	return guo
}

// AddRoles adds the "roles" edges to the GuildRole entity.
func (guo *GuildUpdateOne) AddRoles(g ...*GuildRole) *GuildUpdateOne {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return guo.AddRoleIDs(ids...)
}

// AddYoutubeTalentIDs adds the "youtube_talents" edge to the YouTubeTalent entity by IDs.
func (guo *GuildUpdateOne) AddYoutubeTalentIDs(ids ...string) *GuildUpdateOne {
	guo.mutation.AddYoutubeTalentIDs(ids...)
	return guo
}

// AddYoutubeTalents adds the "youtube_talents" edges to the YouTubeTalent entity.
func (guo *GuildUpdateOne) AddYoutubeTalents(y ...*YouTubeTalent) *GuildUpdateOne {
	ids := make([]string, len(y))
	for i := range y {
		ids[i] = y[i].ID
	}
	return guo.AddYoutubeTalentIDs(ids...)
}

// Mutation returns the GuildMutation object of the builder.
func (guo *GuildUpdateOne) Mutation() *GuildMutation {
	return guo.mutation
}

// ClearMembers clears all "members" edges to the User entity.
func (guo *GuildUpdateOne) ClearMembers() *GuildUpdateOne {
	guo.mutation.ClearMembers()
	return guo
}

// RemoveMemberIDs removes the "members" edge to User entities by IDs.
func (guo *GuildUpdateOne) RemoveMemberIDs(ids ...uint64) *GuildUpdateOne {
	guo.mutation.RemoveMemberIDs(ids...)
	return guo
}

// RemoveMembers removes "members" edges to User entities.
func (guo *GuildUpdateOne) RemoveMembers(u ...*User) *GuildUpdateOne {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return guo.RemoveMemberIDs(ids...)
}

// ClearAdmins clears all "admins" edges to the User entity.
func (guo *GuildUpdateOne) ClearAdmins() *GuildUpdateOne {
	guo.mutation.ClearAdmins()
	return guo
}

// RemoveAdminIDs removes the "admins" edge to User entities by IDs.
func (guo *GuildUpdateOne) RemoveAdminIDs(ids ...uint64) *GuildUpdateOne {
	guo.mutation.RemoveAdminIDs(ids...)
	return guo
}

// RemoveAdmins removes "admins" edges to User entities.
func (guo *GuildUpdateOne) RemoveAdmins(u ...*User) *GuildUpdateOne {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return guo.RemoveAdminIDs(ids...)
}

// ClearRoles clears all "roles" edges to the GuildRole entity.
func (guo *GuildUpdateOne) ClearRoles() *GuildUpdateOne {
	guo.mutation.ClearRoles()
	return guo
}

// RemoveRoleIDs removes the "roles" edge to GuildRole entities by IDs.
func (guo *GuildUpdateOne) RemoveRoleIDs(ids ...uint64) *GuildUpdateOne {
	guo.mutation.RemoveRoleIDs(ids...)
	return guo
}

// RemoveRoles removes "roles" edges to GuildRole entities.
func (guo *GuildUpdateOne) RemoveRoles(g ...*GuildRole) *GuildUpdateOne {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return guo.RemoveRoleIDs(ids...)
}

// ClearYoutubeTalents clears all "youtube_talents" edges to the YouTubeTalent entity.
func (guo *GuildUpdateOne) ClearYoutubeTalents() *GuildUpdateOne {
	guo.mutation.ClearYoutubeTalents()
	return guo
}

// RemoveYoutubeTalentIDs removes the "youtube_talents" edge to YouTubeTalent entities by IDs.
func (guo *GuildUpdateOne) RemoveYoutubeTalentIDs(ids ...string) *GuildUpdateOne {
	guo.mutation.RemoveYoutubeTalentIDs(ids...)
	return guo
}

// RemoveYoutubeTalents removes "youtube_talents" edges to YouTubeTalent entities.
func (guo *GuildUpdateOne) RemoveYoutubeTalents(y ...*YouTubeTalent) *GuildUpdateOne {
	ids := make([]string, len(y))
	for i := range y {
		ids[i] = y[i].ID
	}
	return guo.RemoveYoutubeTalentIDs(ids...)
}

// Where appends a list predicates to the GuildUpdate builder.
func (guo *GuildUpdateOne) Where(ps ...predicate.Guild) *GuildUpdateOne {
	guo.mutation.Where(ps...)
	return guo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (guo *GuildUpdateOne) Select(field string, fields ...string) *GuildUpdateOne {
	guo.fields = append([]string{field}, fields...)
	return guo
}

// Save executes the query and returns the updated Guild entity.
func (guo *GuildUpdateOne) Save(ctx context.Context) (*Guild, error) {
	return withHooks(ctx, guo.sqlSave, guo.mutation, guo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (guo *GuildUpdateOne) SaveX(ctx context.Context) *Guild {
	node, err := guo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (guo *GuildUpdateOne) Exec(ctx context.Context) error {
	_, err := guo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (guo *GuildUpdateOne) ExecX(ctx context.Context) {
	if err := guo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (guo *GuildUpdateOne) check() error {
	if v, ok := guo.mutation.Language(); ok {
		if err := guild.LanguageValidator(v); err != nil {
			return &ValidationError{Name: "language", err: fmt.Errorf(`ent: validator failed for field "Guild.language": %w`, err)}
		}
	}
	return nil
}

func (guo *GuildUpdateOne) sqlSave(ctx context.Context) (_node *Guild, err error) {
	if err := guo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(guild.Table, guild.Columns, sqlgraph.NewFieldSpec(guild.FieldID, field.TypeUint64))
	id, ok := guo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Guild.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := guo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, guild.FieldID)
		for _, f := range fields {
			if !guild.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != guild.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := guo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := guo.mutation.Name(); ok {
		_spec.SetField(guild.FieldName, field.TypeString, value)
	}
	if value, ok := guo.mutation.IconHash(); ok {
		_spec.SetField(guild.FieldIconHash, field.TypeString, value)
	}
	if guo.mutation.IconHashCleared() {
		_spec.ClearField(guild.FieldIconHash, field.TypeString)
	}
	if value, ok := guo.mutation.AuditChannel(); ok {
		_spec.SetField(guild.FieldAuditChannel, field.TypeUint64, value)
	}
	if value, ok := guo.mutation.AddedAuditChannel(); ok {
		_spec.AddField(guild.FieldAuditChannel, field.TypeUint64, value)
	}
	if guo.mutation.AuditChannelCleared() {
		_spec.ClearField(guild.FieldAuditChannel, field.TypeUint64)
	}
	if value, ok := guo.mutation.Language(); ok {
		_spec.SetField(guild.FieldLanguage, field.TypeEnum, value)
	}
	if value, ok := guo.mutation.Settings(); ok {
		_spec.SetField(guild.FieldSettings, field.TypeJSON, value)
	}
	if guo.mutation.SettingsCleared() {
		_spec.ClearField(guild.FieldSettings, field.TypeJSON)
	}
	if guo.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.MembersTable,
			Columns: guild.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := guo.mutation.RemovedMembersIDs(); len(nodes) > 0 && !guo.mutation.MembersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.MembersTable,
			Columns: guild.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := guo.mutation.MembersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.MembersTable,
			Columns: guild.MembersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if guo.mutation.AdminsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.AdminsTable,
			Columns: guild.AdminsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := guo.mutation.RemovedAdminsIDs(); len(nodes) > 0 && !guo.mutation.AdminsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.AdminsTable,
			Columns: guild.AdminsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := guo.mutation.AdminsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   guild.AdminsTable,
			Columns: guild.AdminsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if guo.mutation.RolesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   guild.RolesTable,
			Columns: []string{guild.RolesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guildrole.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := guo.mutation.RemovedRolesIDs(); len(nodes) > 0 && !guo.mutation.RolesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   guild.RolesTable,
			Columns: []string{guild.RolesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guildrole.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := guo.mutation.RolesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   guild.RolesTable,
			Columns: []string{guild.RolesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(guildrole.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if guo.mutation.YoutubeTalentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   guild.YoutubeTalentsTable,
			Columns: guild.YoutubeTalentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(youtubetalent.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := guo.mutation.RemovedYoutubeTalentsIDs(); len(nodes) > 0 && !guo.mutation.YoutubeTalentsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   guild.YoutubeTalentsTable,
			Columns: guild.YoutubeTalentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(youtubetalent.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := guo.mutation.YoutubeTalentsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: true,
			Table:   guild.YoutubeTalentsTable,
			Columns: guild.YoutubeTalentsPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(youtubetalent.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Guild{config: guo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, guo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{guild.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	guo.mutation.done = true
	return _node, nil
}
