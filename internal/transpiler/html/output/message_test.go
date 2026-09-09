package output

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/codegen"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func TestGenMessageNode(t *testing.T) {
	state := codegen.NewState()
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
	if !strings.Contains(code, `dreego.Message(c, "cart.items", dreego.MessageArg{Name: "count", Value: len(items)})`) {
		t.Fatalf("generated code = %s", code)
	}
	use := state.MessageUses["cart.items"]
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
