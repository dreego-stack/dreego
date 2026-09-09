package output

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/codegen"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func TestGenMessageNode(t *testing.T) {
	state := codegen.NewState()
	state.MessageArguments["cart.items"] = map[string]string{"count": "number"}
	node := ir.TemplateNode{
		Type:       ir.NodeMessage,
		MessageKey: "cart.items",
		MessageArgs: []ir.MessageArgument{
			{Name: "count", Expression: "len(items)"},
		},
	}
	code, err := GenTemplateNode(state, node, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(code, `dreego.Message(c, "cart.items", dreego.NumberMessageArg("count", len(items)))`) {
		t.Fatalf("generated code = %s", code)
	}
	if len(state.MessageUses) != 1 {
		t.Fatalf("message uses = %+v", state.MessageUses)
	}
	use := state.MessageUses[0]
	if use.Key != "cart.items" {
		t.Fatalf("message key = %q", use.Key)
	}
	if len(use.Arguments) != 1 || use.Arguments[0] != "count" {
		t.Fatalf("message use = %+v", use)
	}
}

func TestGenComponentMessageNode(t *testing.T) {
	state := codegen.NewState()
	code, err := GenTemplateNodeComp(state, ir.TemplateNode{Type: ir.NodeMessage, MessageKey: "home.title"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(code, `dreego.Message(ctx, "home.title")`) {
		t.Fatalf("generated code = %s", code)
	}
}

func TestMessageNodeIsRejectedInsideRawElement(t *testing.T) {
	state := codegen.NewState()
	inSection := true
	_, err := GenTemplateNodeToState(state, ir.TemplateNode{Type: ir.NodeMessage, MessageKey: "unsafe"}, 0, "b", &inSection)
	if err == nil || !strings.Contains(err.Error(), "not allowed inside script or style") {
		t.Fatalf("unexpected error: %v", err)
	}
}
