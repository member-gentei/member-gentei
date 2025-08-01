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
	"github.com/member-gentei/member-gentei/gentei/ent/schema"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
)

// GuildCreate is the builder for creating a Guild entity.
type GuildCreate struct {
	config
	mutation *GuildMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetName sets the "name" field.
func (gc *GuildCreate) SetName(s string) *GuildCreate {
	gc.mutation.SetName(s)
	return gc
}

// SetIconHash sets the "icon_hash" field.
func (gc *GuildCreate) SetIconHash(s string) *GuildCreate {
	gc.mutation.SetIconHash(s)
	return gc
}

// SetNillableIconHash sets the "icon_hash" field if the given value is not nil.
func (gc *GuildCreate) SetNillableIconHash(s *string) *GuildCreate {
	if s != nil {
		gc.SetIconHash(*s)
	}
	return gc
}

// SetAuditChannel sets the "audit_channel" field.
func (gc *GuildCreate) SetAuditChannel(u uint64) *GuildCreate {
	gc.mutation.SetAuditChannel(u)
	return gc
}

// SetNillableAuditChannel sets the "audit_channel" field if the given value is not nil.
func (gc *GuildCreate) SetNillableAuditChannel(u *uint64) *GuildCreate {
	if u != nil {
		gc.SetAuditChannel(*u)
	}
	return gc
}

// SetLanguage sets the "language" field.
func (gc *GuildCreate) SetLanguage(gu guild.Language) *GuildCreate {
	gc.mutation.SetLanguage(gu)
	return gc
}

// SetNillableLanguage sets the "language" field if the given value is not nil.
func (gc *GuildCreate) SetNillableLanguage(gu *guild.Language) *GuildCreate {
	if gu != nil {
		gc.SetLanguage(*gu)
	}
	return gc
}

// SetSettings sets the "settings" field.
func (gc *GuildCreate) SetSettings(ss *schema.GuildSettings) *GuildCreate {
	gc.mutation.SetSettings(ss)
	return gc
}

// SetID sets the "id" field.
func (gc *GuildCreate) SetID(u uint64) *GuildCreate {
	gc.mutation.SetID(u)
	return gc
}

// AddMemberIDs adds the "members" edge to the User entity by IDs.
func (gc *GuildCreate) AddMemberIDs(ids ...uint64) *GuildCreate {
	gc.mutation.AddMemberIDs(ids...)
	return gc
}

// AddMembers adds the "members" edges to the User entity.
func (gc *GuildCreate) AddMembers(u ...*User) *GuildCreate {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return gc.AddMemberIDs(ids...)
}

// AddAdminIDs adds the "admins" edge to the User entity by IDs.
func (gc *GuildCreate) AddAdminIDs(ids ...uint64) *GuildCreate {
	gc.mutation.AddAdminIDs(ids...)
	return gc
}

// AddAdmins adds the "admins" edges to the User entity.
func (gc *GuildCreate) AddAdmins(u ...*User) *GuildCreate {
	ids := make([]uint64, len(u))
	for i := range u {
		ids[i] = u[i].ID
	}
	return gc.AddAdminIDs(ids...)
}

// AddRoleIDs adds the "roles" edge to the GuildRole entity by IDs.
func (gc *GuildCreate) AddRoleIDs(ids ...uint64) *GuildCreate {
	gc.mutation.AddRoleIDs(ids...)
	return gc
}

// AddRoles adds the "roles" edges to the GuildRole entity.
func (gc *GuildCreate) AddRoles(g ...*GuildRole) *GuildCreate {
	ids := make([]uint64, len(g))
	for i := range g {
		ids[i] = g[i].ID
	}
	return gc.AddRoleIDs(ids...)
}

// AddYoutubeTalentIDs adds the "youtube_talents" edge to the YouTubeTalent entity by IDs.
func (gc *GuildCreate) AddYoutubeTalentIDs(ids ...string) *GuildCreate {
	gc.mutation.AddYoutubeTalentIDs(ids...)
	return gc
}

// AddYoutubeTalents adds the "youtube_talents" edges to the YouTubeTalent entity.
func (gc *GuildCreate) AddYoutubeTalents(y ...*YouTubeTalent) *GuildCreate {
	ids := make([]string, len(y))
	for i := range y {
		ids[i] = y[i].ID
	}
	return gc.AddYoutubeTalentIDs(ids...)
}

// Mutation returns the GuildMutation object of the builder.
func (gc *GuildCreate) Mutation() *GuildMutation {
	return gc.mutation
}

// Save creates the Guild in the database.
func (gc *GuildCreate) Save(ctx context.Context) (*Guild, error) {
	gc.defaults()
	return withHooks(ctx, gc.sqlSave, gc.mutation, gc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (gc *GuildCreate) SaveX(ctx context.Context) *Guild {
	v, err := gc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gc *GuildCreate) Exec(ctx context.Context) error {
	_, err := gc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gc *GuildCreate) ExecX(ctx context.Context) {
	if err := gc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gc *GuildCreate) defaults() {
	if _, ok := gc.mutation.Language(); !ok {
		v := guild.DefaultLanguage
		gc.mutation.SetLanguage(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (gc *GuildCreate) check() error {
	if _, ok := gc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Guild.name"`)}
	}
	if _, ok := gc.mutation.Language(); !ok {
		return &ValidationError{Name: "language", err: errors.New(`ent: missing required field "Guild.language"`)}
	}
	if v, ok := gc.mutation.Language(); ok {
		if err := guild.LanguageValidator(v); err != nil {
			return &ValidationError{Name: "language", err: fmt.Errorf(`ent: validator failed for field "Guild.language": %w`, err)}
		}
	}
	return nil
}

func (gc *GuildCreate) sqlSave(ctx context.Context) (*Guild, error) {
	if err := gc.check(); err != nil {
		return nil, err
	}
	_node, _spec := gc.createSpec()
	if err := sqlgraph.CreateNode(ctx, gc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != _node.ID {
		id := _spec.ID.Value.(int64)
		_node.ID = uint64(id)
	}
	gc.mutation.id = &_node.ID
	gc.mutation.done = true
	return _node, nil
}

func (gc *GuildCreate) createSpec() (*Guild, *sqlgraph.CreateSpec) {
	var (
		_node = &Guild{config: gc.config}
		_spec = sqlgraph.NewCreateSpec(guild.Table, sqlgraph.NewFieldSpec(guild.FieldID, field.TypeUint64))
	)
	_spec.OnConflict = gc.conflict
	if id, ok := gc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = id
	}
	if value, ok := gc.mutation.Name(); ok {
		_spec.SetField(guild.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := gc.mutation.IconHash(); ok {
		_spec.SetField(guild.FieldIconHash, field.TypeString, value)
		_node.IconHash = value
	}
	if value, ok := gc.mutation.AuditChannel(); ok {
		_spec.SetField(guild.FieldAuditChannel, field.TypeUint64, value)
		_node.AuditChannel = value
	}
	if value, ok := gc.mutation.Language(); ok {
		_spec.SetField(guild.FieldLanguage, field.TypeEnum, value)
		_node.Language = value
	}
	if value, ok := gc.mutation.Settings(); ok {
		_spec.SetField(guild.FieldSettings, field.TypeJSON, value)
		_node.Settings = value
	}
	if nodes := gc.mutation.MembersIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := gc.mutation.AdminsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := gc.mutation.RolesIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := gc.mutation.YoutubeTalentsIDs(); len(nodes) > 0 {
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
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Guild.Create().
//		SetName(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.GuildUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (gc *GuildCreate) OnConflict(opts ...sql.ConflictOption) *GuildUpsertOne {
	gc.conflict = opts
	return &GuildUpsertOne{
		create: gc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Guild.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (gc *GuildCreate) OnConflictColumns(columns ...string) *GuildUpsertOne {
	gc.conflict = append(gc.conflict, sql.ConflictColumns(columns...))
	return &GuildUpsertOne{
		create: gc,
	}
}

type (
	// GuildUpsertOne is the builder for "upsert"-ing
	//  one Guild node.
	GuildUpsertOne struct {
		create *GuildCreate
	}

	// GuildUpsert is the "OnConflict" setter.
	GuildUpsert struct {
		*sql.UpdateSet
	}
)

// SetName sets the "name" field.
func (u *GuildUpsert) SetName(v string) *GuildUpsert {
	u.Set(guild.FieldName, v)
	return u
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *GuildUpsert) UpdateName() *GuildUpsert {
	u.SetExcluded(guild.FieldName)
	return u
}

// SetIconHash sets the "icon_hash" field.
func (u *GuildUpsert) SetIconHash(v string) *GuildUpsert {
	u.Set(guild.FieldIconHash, v)
	return u
}

// UpdateIconHash sets the "icon_hash" field to the value that was provided on create.
func (u *GuildUpsert) UpdateIconHash() *GuildUpsert {
	u.SetExcluded(guild.FieldIconHash)
	return u
}

// ClearIconHash clears the value of the "icon_hash" field.
func (u *GuildUpsert) ClearIconHash() *GuildUpsert {
	u.SetNull(guild.FieldIconHash)
	return u
}

// SetAuditChannel sets the "audit_channel" field.
func (u *GuildUpsert) SetAuditChannel(v uint64) *GuildUpsert {
	u.Set(guild.FieldAuditChannel, v)
	return u
}

// UpdateAuditChannel sets the "audit_channel" field to the value that was provided on create.
func (u *GuildUpsert) UpdateAuditChannel() *GuildUpsert {
	u.SetExcluded(guild.FieldAuditChannel)
	return u
}

// AddAuditChannel adds v to the "audit_channel" field.
func (u *GuildUpsert) AddAuditChannel(v uint64) *GuildUpsert {
	u.Add(guild.FieldAuditChannel, v)
	return u
}

// ClearAuditChannel clears the value of the "audit_channel" field.
func (u *GuildUpsert) ClearAuditChannel() *GuildUpsert {
	u.SetNull(guild.FieldAuditChannel)
	return u
}

// SetLanguage sets the "language" field.
func (u *GuildUpsert) SetLanguage(v guild.Language) *GuildUpsert {
	u.Set(guild.FieldLanguage, v)
	return u
}

// UpdateLanguage sets the "language" field to the value that was provided on create.
func (u *GuildUpsert) UpdateLanguage() *GuildUpsert {
	u.SetExcluded(guild.FieldLanguage)
	return u
}

// SetSettings sets the "settings" field.
func (u *GuildUpsert) SetSettings(v *schema.GuildSettings) *GuildUpsert {
	u.Set(guild.FieldSettings, v)
	return u
}

// UpdateSettings sets the "settings" field to the value that was provided on create.
func (u *GuildUpsert) UpdateSettings() *GuildUpsert {
	u.SetExcluded(guild.FieldSettings)
	return u
}

// ClearSettings clears the value of the "settings" field.
func (u *GuildUpsert) ClearSettings() *GuildUpsert {
	u.SetNull(guild.FieldSettings)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create except the ID field.
// Using this option is equivalent to using:
//
//	client.Guild.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(guild.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *GuildUpsertOne) UpdateNewValues() *GuildUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		if _, exists := u.create.mutation.ID(); exists {
			s.SetIgnore(guild.FieldID)
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Guild.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *GuildUpsertOne) Ignore() *GuildUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *GuildUpsertOne) DoNothing() *GuildUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the GuildCreate.OnConflict
// documentation for more info.
func (u *GuildUpsertOne) Update(set func(*GuildUpsert)) *GuildUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&GuildUpsert{UpdateSet: update})
	}))
	return u
}

// SetName sets the "name" field.
func (u *GuildUpsertOne) SetName(v string) *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *GuildUpsertOne) UpdateName() *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateName()
	})
}

// SetIconHash sets the "icon_hash" field.
func (u *GuildUpsertOne) SetIconHash(v string) *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.SetIconHash(v)
	})
}

// UpdateIconHash sets the "icon_hash" field to the value that was provided on create.
func (u *GuildUpsertOne) UpdateIconHash() *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateIconHash()
	})
}

// ClearIconHash clears the value of the "icon_hash" field.
func (u *GuildUpsertOne) ClearIconHash() *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.ClearIconHash()
	})
}

// SetAuditChannel sets the "audit_channel" field.
func (u *GuildUpsertOne) SetAuditChannel(v uint64) *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.SetAuditChannel(v)
	})
}

// AddAuditChannel adds v to the "audit_channel" field.
func (u *GuildUpsertOne) AddAuditChannel(v uint64) *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.AddAuditChannel(v)
	})
}

// UpdateAuditChannel sets the "audit_channel" field to the value that was provided on create.
func (u *GuildUpsertOne) UpdateAuditChannel() *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateAuditChannel()
	})
}

// ClearAuditChannel clears the value of the "audit_channel" field.
func (u *GuildUpsertOne) ClearAuditChannel() *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.ClearAuditChannel()
	})
}

// SetLanguage sets the "language" field.
func (u *GuildUpsertOne) SetLanguage(v guild.Language) *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.SetLanguage(v)
	})
}

// UpdateLanguage sets the "language" field to the value that was provided on create.
func (u *GuildUpsertOne) UpdateLanguage() *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateLanguage()
	})
}

// SetSettings sets the "settings" field.
func (u *GuildUpsertOne) SetSettings(v *schema.GuildSettings) *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.SetSettings(v)
	})
}

// UpdateSettings sets the "settings" field to the value that was provided on create.
func (u *GuildUpsertOne) UpdateSettings() *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateSettings()
	})
}

// ClearSettings clears the value of the "settings" field.
func (u *GuildUpsertOne) ClearSettings() *GuildUpsertOne {
	return u.Update(func(s *GuildUpsert) {
		s.ClearSettings()
	})
}

// Exec executes the query.
func (u *GuildUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for GuildCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *GuildUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *GuildUpsertOne) ID(ctx context.Context) (id uint64, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *GuildUpsertOne) IDX(ctx context.Context) uint64 {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// GuildCreateBulk is the builder for creating many Guild entities in bulk.
type GuildCreateBulk struct {
	config
	err      error
	builders []*GuildCreate
	conflict []sql.ConflictOption
}

// Save creates the Guild entities in the database.
func (gcb *GuildCreateBulk) Save(ctx context.Context) ([]*Guild, error) {
	if gcb.err != nil {
		return nil, gcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(gcb.builders))
	nodes := make([]*Guild, len(gcb.builders))
	mutators := make([]Mutator, len(gcb.builders))
	for i := range gcb.builders {
		func(i int, root context.Context) {
			builder := gcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*GuildMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, gcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = gcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, gcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil && nodes[i].ID == 0 {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = uint64(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, gcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (gcb *GuildCreateBulk) SaveX(ctx context.Context) []*Guild {
	v, err := gcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (gcb *GuildCreateBulk) Exec(ctx context.Context) error {
	_, err := gcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gcb *GuildCreateBulk) ExecX(ctx context.Context) {
	if err := gcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.Guild.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.GuildUpsert) {
//			SetName(v+v).
//		}).
//		Exec(ctx)
func (gcb *GuildCreateBulk) OnConflict(opts ...sql.ConflictOption) *GuildUpsertBulk {
	gcb.conflict = opts
	return &GuildUpsertBulk{
		create: gcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.Guild.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (gcb *GuildCreateBulk) OnConflictColumns(columns ...string) *GuildUpsertBulk {
	gcb.conflict = append(gcb.conflict, sql.ConflictColumns(columns...))
	return &GuildUpsertBulk{
		create: gcb,
	}
}

// GuildUpsertBulk is the builder for "upsert"-ing
// a bulk of Guild nodes.
type GuildUpsertBulk struct {
	create *GuildCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.Guild.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//			sql.ResolveWith(func(u *sql.UpdateSet) {
//				u.SetIgnore(guild.FieldID)
//			}),
//		).
//		Exec(ctx)
func (u *GuildUpsertBulk) UpdateNewValues() *GuildUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(s *sql.UpdateSet) {
		for _, b := range u.create.builders {
			if _, exists := b.mutation.ID(); exists {
				s.SetIgnore(guild.FieldID)
			}
		}
	}))
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.Guild.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *GuildUpsertBulk) Ignore() *GuildUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *GuildUpsertBulk) DoNothing() *GuildUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the GuildCreateBulk.OnConflict
// documentation for more info.
func (u *GuildUpsertBulk) Update(set func(*GuildUpsert)) *GuildUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&GuildUpsert{UpdateSet: update})
	}))
	return u
}

// SetName sets the "name" field.
func (u *GuildUpsertBulk) SetName(v string) *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.SetName(v)
	})
}

// UpdateName sets the "name" field to the value that was provided on create.
func (u *GuildUpsertBulk) UpdateName() *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateName()
	})
}

// SetIconHash sets the "icon_hash" field.
func (u *GuildUpsertBulk) SetIconHash(v string) *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.SetIconHash(v)
	})
}

// UpdateIconHash sets the "icon_hash" field to the value that was provided on create.
func (u *GuildUpsertBulk) UpdateIconHash() *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateIconHash()
	})
}

// ClearIconHash clears the value of the "icon_hash" field.
func (u *GuildUpsertBulk) ClearIconHash() *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.ClearIconHash()
	})
}

// SetAuditChannel sets the "audit_channel" field.
func (u *GuildUpsertBulk) SetAuditChannel(v uint64) *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.SetAuditChannel(v)
	})
}

// AddAuditChannel adds v to the "audit_channel" field.
func (u *GuildUpsertBulk) AddAuditChannel(v uint64) *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.AddAuditChannel(v)
	})
}

// UpdateAuditChannel sets the "audit_channel" field to the value that was provided on create.
func (u *GuildUpsertBulk) UpdateAuditChannel() *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateAuditChannel()
	})
}

// ClearAuditChannel clears the value of the "audit_channel" field.
func (u *GuildUpsertBulk) ClearAuditChannel() *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.ClearAuditChannel()
	})
}

// SetLanguage sets the "language" field.
func (u *GuildUpsertBulk) SetLanguage(v guild.Language) *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.SetLanguage(v)
	})
}

// UpdateLanguage sets the "language" field to the value that was provided on create.
func (u *GuildUpsertBulk) UpdateLanguage() *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateLanguage()
	})
}

// SetSettings sets the "settings" field.
func (u *GuildUpsertBulk) SetSettings(v *schema.GuildSettings) *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.SetSettings(v)
	})
}

// UpdateSettings sets the "settings" field to the value that was provided on create.
func (u *GuildUpsertBulk) UpdateSettings() *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.UpdateSettings()
	})
}

// ClearSettings clears the value of the "settings" field.
func (u *GuildUpsertBulk) ClearSettings() *GuildUpsertBulk {
	return u.Update(func(s *GuildUpsert) {
		s.ClearSettings()
	})
}

// Exec executes the query.
func (u *GuildUpsertBulk) Exec(ctx context.Context) error {
	if u.create.err != nil {
		return u.create.err
	}
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the GuildCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for GuildCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *GuildUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
