package membership

import (
	"context"

	"github.com/member-gentei/member-gentei/gentei/ent"
	"github.com/member-gentei/member-gentei/gentei/ent/hook"
)

// ChangeHandler is implemented by things that react to membership changes.
type ChangeHandler interface {
	GainedMembership(userMembershipID int)
	LostMembership(userMembershipID int)
	SetChangeReason(reason string)
}

// HookMembershipChanges uses ent hooks to listen to relevant changes to UserMembership ents.
func HookMembershipChanges(db *ent.Client, handler ChangeHandler) {
	db.UserMembership.Use(
		hook.On(
			func(next ent.Mutator) ent.Mutator {
				return createPostUserMembershipFunc(next, handler)
			},
			ent.OpCreate|ent.OpUpdate|ent.OpUpdateOne,
		),
	)
}

func createPostUserMembershipFunc(next ent.Mutator, handler ChangeHandler) hook.UserMembershipFunc {
	return func(ctx context.Context, m *ent.UserMembershipMutation) (ent.Value, error) {
		v, err := next.Mutate(ctx, m)
		if err != nil {
			return v, err
		}
		var (
			affectedIDs []int
			op          = m.Op()
		)
		if op.Is(ent.OpCreate) || op.Is(ent.OpUpdateOne) {
			umID, _ := m.ID()
			affectedIDs = []int{umID}
		} else if op.Is(ent.OpUpdate) {
			affectedIDs, _ = m.IDs(ctx)
		}
		for _, userMembershipID := range affectedIDs {
			// if created with a FailCount of 0, this is an added/new membership.
			if op.Is(ent.OpCreate) {
				if failCount, included := m.FailCount(); included && failCount == 0 {
					handler.GainedMembership(userMembershipID)
				}
				// creates only have 1 ID, so just break out here.
				break
			}
			var (
				failCount, failCountIsSet   = m.FailCount()
				_, failCountIncremented     = m.AddedFailCount()
				firstFailed, firstFailedSet = m.FirstFailed()
			)
			if failCountIsSet && firstFailedSet && failCount == 0 && firstFailed.IsZero() {
				// if FailCount is set to 0 and FirstFailed is zeroed, this is a regained membership.
				handler.GainedMembership(userMembershipID)
			} else if failCountIncremented && firstFailedSet && !firstFailed.IsZero() {
				// if FailCount is incremented and FirstFailed is set, this is a newly lost membership.
				handler.LostMembership(userMembershipID)
			}

		}
		return v, err
	}
}
