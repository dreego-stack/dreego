package core

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
	result, err := genTemplateNode(n, 0)
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
	result, err := genTemplateNode(n, 0)
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
	result, err := genTemplateNode(n, 0)
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

func TestGenTemplateNodeSlotNamed(t *testing.T) {
	n := TemplateNode{
		Type:    NodeSlot,
		Content: "name",
	}
	result, err := genTemplateNode(n, 0)
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
	result, err := genTemplateNode(n, 0)
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
	result, err := genTemplateNode(n, 0)
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
	result, err := genTemplateNode(n, 0)
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
	result, err := genTemplateNode(n, 0)
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
	result, err := genTemplateNode(n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "`<b>raw</b>`") {
		t.Errorf("verbatim must emit literal content, got:\n%s", result)
	}
	if strings.Contains(result, "html.EscapeString") {
		t.Errorf("verbatim must NOT escape, got:\n%s", result)
	}
}

func TestGenTemplateNodeFilterRawUpper(t *testing.T) {
	n := TemplateNode{
		Type:    NodeExpression,
		Content: "name",
		Filters: []string{"raw", "upper"},
	}
	result, err := genTemplateNode(n, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "strings.ToUpper") {
		t.Errorf("upper filter must wrap in strings.ToUpper, got:\n%s", result)
	}
	if strings.Contains(result, "html.EscapeString") {
		t.Errorf("raw filter must skip html.EscapeString, got:\n%s", result)
	}
}
