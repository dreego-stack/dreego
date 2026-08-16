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
	out, err := genTemplateNodeComp(NewGenerator(), n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"if len(items) > 0", "for i, item := range items", "dreego.EachLoop{", "} else {", "`empty`"} {
		if !strings.Contains(out, want) {
			t.Errorf("compGen each-with-else missing %q, got:\n%s", want, out)
		}
	}
}

// compGen: $loop. inside a {#if} cond nested in an {#each} body must be
// substituted too — the whole generated child code (including NodeIf conds)
// is rewritten, not only direct expressions.
func TestCompGenEachSubstitutesLoopInIfCond(t *testing.T) {
	n := TemplateNode{
		Type:  NodeEach,
		Item:  "item",
		Items: "items",
		Children: []TemplateNode{
			{
				Type: NodeIf,
				Cond: "!$loop.Last",
				Children: []TemplateNode{
					{Type: NodeText, Content: ", "},
				},
			},
			{Type: NodeExpression, Content: "item.Name"},
		},
	}
	out, err := genTemplateNodeComp(NewGenerator(), n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "if !loop.Last {") {
		t.Errorf("compGen $loop. in {#if} cond must be substituted to loop., got:\n%s", out)
	}
	if strings.Contains(out, "$loop.") {
		t.Errorf("compGen raw $loop. must not remain, got:\n%s", out)
	}
}

// compGen full pipeline (lex → parse → codegen): {#if !$loop.Last} inside
// {#each} must produce "if !loop.Last {".
func TestCompGenEachLoopInIfCondFullParse(t *testing.T) {
	input := `<div>{#each items as item}<span>{#if !$loop.Last}, {/if}{{ item }}</span>{/each}</div>`
	tokens, err := Lex(input)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := genTemplateNodeComp(NewGenerator(), file.Template.Nodes[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "if !loop.Last {") {
		t.Errorf("compGen full pipeline must substitute $loop. in {#if} cond, got:\n%s", out)
	}
	if strings.Contains(out, "$loop.") {
		t.Errorf("compGen raw $loop. must not remain, got:\n%s", out)
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
	out, err := genTemplateNodeComp(NewGenerator(), n)
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

// compGen NodeIf with mixed else children (NodeIf + NodeText) must emit a
// plain "else" block with recursive processing of both children, not an
// else-if chain.
func TestCompGenIfElseMixedChildren(t *testing.T) {
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
			{Type: NodeText, Content: "fallback"},
		},
	}
	out, err := genTemplateNodeComp(NewGenerator(), n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"if a {", "} else {", "if b {", "`A`", "`B`", "`fallback`"} {
		if !strings.Contains(out, want) {
			t.Errorf("compGen mixed else missing %q, got:\n%s", want, out)
		}
	}
	if strings.Contains(out, "} else if b {") {
		t.Errorf("compGen mixed else must not emit an else-if chain, got:\n%s", out)
	}
}

// compGen NodeSlot with a name must read the named slot from the context.
func TestCompGenSlotNamed(t *testing.T) {
	n := TemplateNode{Type: NodeSlot, Content: "header"}
	out, err := genTemplateNodeComp(NewGenerator(), n)
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
	out, err := genTemplateNodeComp(NewGenerator(), n)
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
	out, err := genTemplateNodeComp(NewGenerator(), n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "`<b>raw</b>`") {
		t.Errorf("compGen verbatim must emit literal content, got:\n%s", out)
	}
	if strings.Contains(out, "dreego.SafeText") {
		t.Errorf("compGen verbatim must NOT escape, got:\n%s", out)
	}
}

// compGen NodeExpression with raw + upper filters: upper must wrap the value,
// raw must skip dreego.SafeText.
func TestCompGenFilterRawUpper(t *testing.T) {
	n := TemplateNode{
		Type:    NodeExpression,
		Content: "name",
		Filters: []string{"raw", "upper"},
	}
	out, err := genTemplateNodeComp(NewGenerator(), n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "strings.ToUpper") {
		t.Errorf("compGen upper filter must wrap in strings.ToUpper, got:\n%s", out)
	}
	if strings.Contains(out, "dreego.SafeText") {
		t.Errorf("compGen raw filter must skip dreego.SafeText, got:\n%s", out)
	}
}

// compGen NodeExpression without raw must escape by default.
func TestCompGenExpressionEscapesByDefault(t *testing.T) {
	n := TemplateNode{Type: NodeExpression, Content: "name"}
	out, err := genTemplateNodeComp(NewGenerator(), n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "dreego.SafeText") {
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
	out, err := (&compGen{gen: NewGenerator()}).genComponentCall(n)
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
