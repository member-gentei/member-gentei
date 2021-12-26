// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/member-gentei/member-gentei/gentei/ent/guildrole"
	"github.com/member-gentei/member-gentei/gentei/ent/schema"
	"github.com/member-gentei/member-gentei/gentei/ent/user"
	"github.com/member-gentei/member-gentei/gentei/ent/youtubetalent"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	guildFields := schema.Guild{}.Fields()
	_ = guildFields
	guildroleFields := schema.GuildRole{}.Fields()
	_ = guildroleFields
	// guildroleDescLastUpdated is the schema descriptor for last_updated field.
	guildroleDescLastUpdated := guildroleFields[2].Descriptor()
	// guildrole.DefaultLastUpdated holds the default value on creation for the last_updated field.
	guildrole.DefaultLastUpdated = guildroleDescLastUpdated.Default.(func() time.Time)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescFullName is the schema descriptor for full_name field.
	userDescFullName := userFields[1].Descriptor()
	// user.FullNameValidator is a validator for the "full_name" field. It is called by the builders before save.
	user.FullNameValidator = userDescFullName.Validators[0].(func(string) error)
	// userDescLastCheck is the schema descriptor for last_check field.
	userDescLastCheck := userFields[3].Descriptor()
	// user.DefaultLastCheck holds the default value on creation for the last_check field.
	user.DefaultLastCheck = userDescLastCheck.Default.(func() time.Time)
	youtubetalentFields := schema.YouTubeTalent{}.Fields()
	_ = youtubetalentFields
	// youtubetalentDescLastUpdated is the schema descriptor for last_updated field.
	youtubetalentDescLastUpdated := youtubetalentFields[3].Descriptor()
	// youtubetalent.DefaultLastUpdated holds the default value on creation for the last_updated field.
	youtubetalent.DefaultLastUpdated = youtubetalentDescLastUpdated.Default.(func() time.Time)
}
