package discord

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
)

const expectedOutput = `Memberships confirmed! You will be granted roles corresponding to the following channels:
` + "◦ `uwu`" + `
` + "◦ `what's`" + `
` + "◦ `this`"

func TestMultiMembersConfirmedTemplate(t *testing.T) {
	// tests that the template does newlines correctly.
	var (
		buf        bytes.Buffer
		mmTemplate = template.Must(template.New("multimember").Parse(multiMembersConfirmed))
	)
	err := mmTemplate.Execute(&buf, map[string]interface{}{
		"titles": []string{"uwu", "what's", "this"},
	})
	if err != nil {
		t.Fatalf("error executing template: %s", err)
	}
	if diff := cmp.Diff(expectedOutput, buf.String()); diff != "" {
		t.Fatalf("unexpected template output: %s", diff)
	}
}
