package output

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/codegen"
)

func TestMessageInAttributeUsesAttributeContext(t *testing.T) {
	state := codegen.NewState()
	state.MessageArguments["account.url"] = map[string]string{}
	state.MessageArguments["account.title"] = map[string]string{"name": "string"}
	code, _, err := compTextSection(state, `<a href="[[ account.url ]]" title="[[ account.title name=user.Name ]]">`, false, "c")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(code, `dreego.SafeURL(dreego.Message(c, "account.url"))`) {
		t.Fatalf("URL message is not URL-safe: %s", code)
	}
	if !strings.Contains(code, `dreego.SafeAttr(dreego.Message(c, "account.title", dreego.StringMessageArg("name", user.Name)))`) {
		t.Fatalf("title message is not attribute-safe: %s", code)
	}
	if len(state.MessageUses) != 2 {
		t.Fatalf("message uses = %+v", state.MessageUses)
	}
}
