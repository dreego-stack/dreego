package core

import (
	"strings"
	"testing"
)

// compGen NodeEach with an else branch must wrap the range in
// "if len(items) > 0 { … } else { … }" and emit the EachLoop helper.
func TestCompGenEachWithElse(t *testing.T) {
	n := TemplateNode{
		Type:  NodeEach,
		Item:  "item",
		Items: "items",
		Children: []TemplateNode{
			{Type: NodeText, Content: "X"},
		},
		ElseChildren: []TemplateNode{
			{Type: NodeText, Content: "empty"},
		},
	}
	out, err := genTemplateNodeComp(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"if len(items) > 0", "for i, item := range items", "dreego.EachLoop{", "} else {", "`empty`"} {
		if !strings.Contains(out, want) {
			t.Errorf("compGen each-with-else missing %q, got:\n%s", want, out)
		}
	}
}

// compGen NodeIf with a chain of else-if nodes must emit "else if" rather
// than a single plain else.
func TestCompGenIfElseIfChain(t *testing.T) {
	n := TemplateNode{
		Type: NodeIf,
		Cond: "a",
		Children: []TemplateNode{
			{Type: NodeText, Content: "A"},
		},
		ElseChildren: []TemplateNode{
			{
				Type: NodeIf,
				Cond: "b",
				Children: []TemplateNode{
					{Type: NodeText, Content: "B"},
				},
			},
		},
	}
	out, err := genTemplateNodeComp(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"if a {", "} else if b {", "`A`", "`B`"} {
		if !strings.Contains(out, want) {
			t.Errorf("compGen else-if chain missing %q, got:\n%s", want, out)
		}
	}
	if strings.Contains(out, "} else {\n") {
		t.Errorf("compGen else-if chain must not emit a trailing plain else, got:\n%s", out)
	}
}

// compGen NodeSlot with a name must read the named slot from the context.
func TestCompGenSlotNamed(t *testing.T) {
	n := TemplateNode{Type: NodeSlot, Content: "header"}
	out, err := genTemplateNodeComp(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `ctx.Get("slot_header")`) {
		t.Errorf("compGen named slot must read ctx.Get(\"slot_header\"), got:\n%s", out)
	}
}

// compGen NodeSlot without a name must read the default slot.
func TestCompGenSlotDefault(t *testing.T) {
	n := TemplateNode{Type: NodeSlot}
	out, err := genTemplateNodeComp(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `ctx.Get("slot")`) {
		t.Errorf("compGen default slot must read ctx.Get(\"slot\"), got:\n%s", out)
	}
}

// compGen NodeVerbatim must emit the raw content without escaping.
func TestCompGenVerbatim(t *testing.T) {
	n := TemplateNode{Type: NodeVerbatim, Content: "<b>raw</b>"}
	out, err := genTemplateNodeComp(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "`<b>raw</b>`") {
		t.Errorf("compGen verbatim must emit literal content, got:\n%s", out)
	}
	if strings.Contains(out, "html.EscapeString") {
		t.Errorf("compGen verbatim must NOT escape, got:\n%s", out)
	}
}

// compGen NodeExpression with raw + upper filters: upper must wrap the value,
// raw must skip html.EscapeString.
func TestCompGenFilterRawUpper(t *testing.T) {
	n := TemplateNode{
		Type:    NodeExpression,
		Content: "name",
		Filters: []string{"raw", "upper"},
	}
	out, err := genTemplateNodeComp(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "strings.ToUpper") {
		t.Errorf("compGen upper filter must wrap in strings.ToUpper, got:\n%s", out)
	}
	if strings.Contains(out, "html.EscapeString") {
		t.Errorf("compGen raw filter must skip html.EscapeString, got:\n%s", out)
	}
}

// compGen NodeExpression without raw must escape by default.
func TestCompGenExpressionEscapesByDefault(t *testing.T) {
	n := TemplateNode{Type: NodeExpression, Content: "name"}
	out, err := genTemplateNodeComp(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "html.EscapeString") {
		t.Errorf("compGen expression must escape by default, got:\n%s", out)
	}
}

// genComponentCall for a non-self-close node must return a literal <@Tag>
// placeholder, not invoke the component.
func TestGenComponentCallNonSelfClose(t *testing.T) {
	n := TemplateNode{
		Type:      NodeComponentCall,
		Tag:       "Nav.Link",
		Attrs:     `href="/x"`,
		SelfClose: false,
	}
	out, err := genComponentCall(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "<@Nav.Link>") {
		t.Errorf("non-self-close genComponentCall must emit <@Nav.Link> placeholder, got:\n%s", out)
	}
	if strings.Contains(out, ".Render(ctx)") {
		t.Errorf("non-self-close genComponentCall must NOT invoke the component, got:\n%s", out)
	}
}
