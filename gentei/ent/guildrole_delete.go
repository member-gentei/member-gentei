// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
)

// GuildRoleDelete is the builder for deleting a GuildRole entity.
type GuildRoleDelete struct {
	config
	hooks    []Hook
	mutation *GuildRoleMutation
}

// Where appends a list predicates to the GuildRoleDelete builder.
func (grd *GuildRoleDelete) Where(ps ...predicate.GuildRole) *GuildRoleDelete {
	grd.mutation.Where(ps...)
	return grd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (grd *GuildRoleDelete) Exec(ctx context.Context) (int, error) {
	return withHooks[int, GuildRoleMutation](ctx, grd.sqlExec, grd.mutation, grd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (grd *GuildRoleDelete) ExecX(ctx context.Context) int {
	n, err := grd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (grd *GuildRoleDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(guildrole.Table, sqlgraph.NewFieldSpec(guildrole.FieldID, field.TypeUint64))
	if ps := grd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, grd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	grd.mutation.done = true
	return affected, err
}

// GuildRoleDeleteOne is the builder for deleting a single GuildRole entity.
type GuildRoleDeleteOne struct {
	grd *GuildRoleDelete
}

// Where appends a list predicates to the GuildRoleDelete builder.
func (grdo *GuildRoleDeleteOne) Where(ps ...predicate.GuildRole) *GuildRoleDeleteOne {
	grdo.grd.mutation.Where(ps...)
	return grdo
}

// Exec executes the deletion query.
func (grdo *GuildRoleDeleteOne) Exec(ctx context.Context) error {
	n, err := grdo.grd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{guildrole.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (grdo *GuildRoleDeleteOne) ExecX(ctx context.Context) {
	if err := grdo.Exec(ctx); err != nil {
		panic(err)
	}
}
