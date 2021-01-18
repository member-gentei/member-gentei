package discord

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/member-gentei/member-gentei/bot/discord/lang"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const expectedOutput = `Memberships confirmed! You will be granted roles corresponding to the following channels:
` + "◦ `uwu`" + `
` + "◦ `what's`" + `
` + "◦ `this`"

func TestMultiMembershipConfirmedReply(t *testing.T) {
	// tests that the template does newlines correctly.
	var (
		localizer = makeLocalizer(lang.NewBundle(), "")
	)
	output := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "MembershipConfirmedReply",
			One:   "Membership confirmed! You will be added as a member shortly.",
			Other: multiMembersConfirmed,
		},
		TemplateData: map[string]interface{}{
			"titles": []string{"uwu", "what's", "this"},
		},
		PluralCount: 3,
	})
	if diff := cmp.Diff(expectedOutput, output); diff != "" {
		t.Fatalf("unexpected template output: %s", diff)
	}
}
