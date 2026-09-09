package head

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/codegen"
)

func TestGenWithMessages(t *testing.T) {
	state := codegen.NewState()
	state.MessageArguments["page.title"] = map[string]string{"name": "string"}
	code, err := GenWithMessages(state, `<title>[[ page.title name=user.Name ]]</title>`, "b", "c")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(code, `dreego.SafeText(dreego.Message(c, "page.title", dreego.StringMessageArg("name", user.Name)))`) {
		t.Fatalf("generated code = %s", code)
	}
	if len(state.MessageUses) != 1 || state.MessageUses[0].Key != "page.title" {
		t.Fatalf("message uses = %+v", state.MessageUses)
	}
}
