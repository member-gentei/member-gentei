// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

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
	var (
		err      error
		affected int
	)
	if len(grd.hooks) == 0 {
		affected, err = grd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*GuildRoleMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			grd.mutation = mutation
			affected, err = grd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(grd.hooks) - 1; i >= 0; i-- {
			if grd.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = grd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, grd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
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
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: guildrole.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUint64,
				Column: guildrole.FieldID,
			},
		},
	}
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
	return affected, err
}

// GuildRoleDeleteOne is the builder for deleting a single GuildRole entity.
type GuildRoleDeleteOne struct {
	grd *GuildRoleDelete
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
	grdo.grd.ExecX(ctx)
}
