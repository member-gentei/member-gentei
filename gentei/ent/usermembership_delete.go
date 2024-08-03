// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/member-gentei/member-gentei/gentei/ent/predicate"
	"github.com/member-gentei/member-gentei/gentei/ent/usermembership"
)

// UserMembershipDelete is the builder for deleting a UserMembership entity.
type UserMembershipDelete struct {
	config
	hooks    []Hook
	mutation *UserMembershipMutation
}

// Where appends a list predicates to the UserMembershipDelete builder.
func (umd *UserMembershipDelete) Where(ps ...predicate.UserMembership) *UserMembershipDelete {
	umd.mutation.Where(ps...)
	return umd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (umd *UserMembershipDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, umd.sqlExec, umd.mutation, umd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (umd *UserMembershipDelete) ExecX(ctx context.Context) int {
	n, err := umd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (umd *UserMembershipDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(usermembership.Table, sqlgraph.NewFieldSpec(usermembership.FieldID, field.TypeInt))
	if ps := umd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, umd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	umd.mutation.done = true
	return affected, err
}

// UserMembershipDeleteOne is the builder for deleting a single UserMembership entity.
type UserMembershipDeleteOne struct {
	umd *UserMembershipDelete
}

// Where appends a list predicates to the UserMembershipDelete builder.
func (umdo *UserMembershipDeleteOne) Where(ps ...predicate.UserMembership) *UserMembershipDeleteOne {
	umdo.umd.mutation.Where(ps...)
	return umdo
}

// Exec executes the deletion query.
func (umdo *UserMembershipDeleteOne) Exec(ctx context.Context) error {
	n, err := umdo.umd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{usermembership.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (umdo *UserMembershipDeleteOne) ExecX(ctx context.Context) {
	if err := umdo.Exec(ctx); err != nil {
		panic(err)
	}
}
