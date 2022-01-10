package membership

import (
	"context"

	"github.com/member-gentei/member-gentei/gentei/ent"
)

type CheckForUserOptions struct {
	// Specify ChannelIDs to restruct checks to these channels.
	ChannelIDs []string
}

func CheckForUser(ctx context.Context, db *ent.Client, userID uint64, options *CheckForUserOptions) error {
	return nil
}
