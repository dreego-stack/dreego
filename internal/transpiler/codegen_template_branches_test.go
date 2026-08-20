package transpiler

import (
	"strings"
	"testing"
)

func TestGenTemplateNodeIfElseIfChain(t *testing.T) {
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
	gen := NewGenerator()
	gen.registerDef("Button", &ComponentDef{Name: "Button", Props: []Prop{{Name: "label", Type: "string"}}})
	result, err := genTemplateNode(gen, n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"if a", "else if b", "`A`", "`B`"} {
		if !strings.Contains(result, want) {
			t.Errorf("else-if chain missing %q, got:\n%s", want, result)
		}
	}
}

func TestGenTemplateNodeEachWithElse(t *testing.T) {
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
	gen := NewGenerator()
	gen.registerDef("Card", &ComponentDef{Name: "Card", Slots: []string{"header"}, HasDefaultSlot: true})
	result, err := genTemplateNode(gen, n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"if len(items) > 0", "else", "`empty`"} {
		if !strings.Contains(result, want) {
			t.Errorf("each-with-else missing %q, got:\n%s", want, result)
		}
	}
}

func TestGenTemplateNodeEachLoopVar(t *testing.T) {
	n := TemplateNode{
		Type:  NodeEach,
		Item:  "item",
		Items: "items",
		Children: []TemplateNode{
			{Type: NodeExpression, Content: "$loop.Index"},
		},
	}
	result, err := genTemplateNode(NewGenerator(), n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "loop.Index") {
		t.Errorf("$loop. must be substituted to loop., got:\n%s", result)
	}
	if strings.Contains(result, "$loop.") {
		t.Errorf("raw $loop. must not remain, got:\n%s", result)
	}
}

// $loop. inside a {#if} cond nested in an {#each} body must be substituted
// too — the whole generated child code (including NodeIf conds) is rewritten,
// not only direct expressions. {#if !$loop.Last} must generate
// "if !loop.Last {" or the emitted Go does not compile.
func TestGenTemplateNodeEachSubstitutesLoopInIfCond(t *testing.T) {
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
	result, err := genTemplateNode(NewGenerator(), n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "if !loop.Last {") {
		t.Errorf("$loop. in {#if} cond must be substituted to loop., got:\n%s", result)
	}
	if strings.Contains(result, "$loop.") {
		t.Errorf("raw $loop. must not remain, got:\n%s", result)
	}
}

// $loop. in an {#else if} cond inside an {#each} body must be substituted too.
func TestGenTemplateNodeEachSubstitutesLoopInElseIfCond(t *testing.T) {
	n := TemplateNode{
		Type:  NodeEach,
		Item:  "item",
		Items: "items",
		Children: []TemplateNode{
			{
				Type: NodeIf,
				Cond: "$loop.First",
				Children: []TemplateNode{
					{Type: NodeText, Content: "A"},
				},
				ElseChildren: []TemplateNode{
					{
						Type:     NodeIf,
						Cond:     "$loop.Last",
						Children: []TemplateNode{{Type: NodeText, Content: "B"}},
					},
				},
			},
		},
	}
	result, err := genTemplateNode(NewGenerator(), n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"if loop.First {", "} else if loop.Last {"} {
		if !strings.Contains(result, want) {
			t.Errorf("$loop. in else-if cond must be substituted, missing %q, got:\n%s", want, result)
		}
	}
	if strings.Contains(result, "$loop.") {
		t.Errorf("raw $loop. must not remain, got:\n%s", result)
	}
}

// Full pipeline (lex → parse → codegen): {#if !$loop.Last} inside {#each}
// must produce "if !loop.Last {" — the exact feedback.md scenario.
func TestGenTemplateNodeEachLoopInIfCondFullParse(t *testing.T) {
	input := `<div>{#each items as item}<span>{#if !$loop.Last}, {/if}{{ item }}</span>{/each}</div>`
	tokens, err := Lex(input)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := genTemplateNode(NewGenerator(), file.Template.Nodes[0], 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "if !loop.Last {") {
		t.Errorf("full pipeline must substitute $loop. in {#if} cond, got:\n%s", out)
	}
	if strings.Contains(out, "$loop.") {
		t.Errorf("raw $loop. must not remain, got:\n%s", out)
	}
}

func TestGenTemplateNodeSlotNamed(t *testing.T) {
	n := TemplateNode{
		Type:    NodeSlot,
		Content: "name",
	}
	result, err := genTemplateNode(NewGenerator(), n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, `c.Get("slot_name")`) {
		t.Errorf("named slot must read c.Get(\"slot_name\"), got:\n%s", result)
	}
}

func TestGenTemplateNodeSlotDefault(t *testing.T) {
	n := TemplateNode{
		Type: NodeSlot,
	}
	result, err := genTemplateNode(NewGenerator(), n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, `c.Get("slot")`) {
		t.Errorf("default slot must read c.Get(\"slot\"), got:\n%s", result)
	}
}

func TestGenTemplateNodeSlotWithChildren(t *testing.T) {
	n := TemplateNode{
		Type:    NodeSlot,
		Content: "name",
		Children: []TemplateNode{
			{Type: NodeText, Content: "fallback"},
		},
	}
	result, err := genTemplateNode(NewGenerator(), n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, `c.Set("slot_name"`) {
		t.Errorf("slot with children must set c.Set(\"slot_name\"), got:\n%s", result)
	}
	if !strings.Contains(result, "`fallback`") {
		t.Errorf("slot fallback children must be emitted, got:\n%s", result)
	}
}

func TestGenTemplateNodeComponentCallSelfClose(t *testing.T) {
	n := TemplateNode{
		Type:      NodeComponentCall,
		Tag:       "Button",
		Attrs:     `label="Hi"`,
		SelfClose: true,
	}
	gen := NewGenerator()
	gen.registerDef("Button", &ComponentDef{Name: "Button", Props: []Prop{{Name: "label", Type: "string"}}})
	result, err := genTemplateNode(gen, n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Button(") {
		t.Errorf("self-close component must call Button(...), got:\n%s", result)
	}
	if !strings.Contains(result, ".Render(c)") {
		t.Errorf("self-close component must call .Render(c), got:\n%s", result)
	}
}

func TestGenTemplateNodeComponentCallWithSlot(t *testing.T) {
	n := TemplateNode{
		Type: NodeComponentCall,
		Tag:  "Card",
		Children: []TemplateNode{
			{
				Type:    NodeSlot,
				Content: "header",
				Children: []TemplateNode{
					{Type: NodeText, Content: "H"},
				},
			},
			{Type: NodeText, Content: "body"},
		},
	}
	gen := NewGenerator()
	gen.registerDef("Card", &ComponentDef{Name: "Card", Slots: []string{"header"}, HasDefaultSlot: true})
	result, err := genTemplateNode(gen, n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{`c.Set("slot_header"`, `c.Set("slot"`, "Card(", ".Render(c)"} {
		if !strings.Contains(result, want) {
			t.Errorf("component-with-slot missing %q, got:\n%s", want, result)
		}
	}
}

func TestGenTemplateNodeVerbatim(t *testing.T) {
	n := TemplateNode{
		Type:    NodeVerbatim,
		Content: "<b>raw</b>",
	}
	result, err := genTemplateNode(NewGenerator(), n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "`<b>raw</b>`") {
		t.Errorf("verbatim must emit literal content, got:\n%s", result)
	}
	if strings.Contains(result, "dreego.SafeText") {
		t.Errorf("verbatim must NOT escape, got:\n%s", result)
	}
}
