// Package templates contains Discord bot message templates that are somewhere near complicated.
package templates

import (
	_ "embed"

	"github.com/cbroglie/mustache"
)

//go:embed role_permission_failure.md
var RoleAlreadyMapped string

type RoleAlreadyMappedContext struct {
	ChannelID,
	ChannelName,
	RoleMention string
}

//go:embed role_applied.md
var RoleApplied string

type RoleAppliedContext struct {
	ChannelName,
	RoleMention string
}

//go:embed role_permission_failure.md
var RolePermissionFailure string

type RolePermissionFailureContext struct {
	Action,
	RoleMention string
}

func MustRender(data string, context ...interface{}) string {
	rendered, err := mustache.Render(data, context...)
	if err != nil {
		panic(err)
	}
	return rendered
}

// plaintext messages

//go:embed user_deleted.md
var PlaintextUserDeleted string
