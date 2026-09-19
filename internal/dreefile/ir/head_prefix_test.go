package ir

import "testing"

func TestHasHeadDedupeTagCaseInsensitive(t *testing.T) {
	cases := map[string]bool{
		`<title>Page</title>`:                   true,
		`<TITLE>Page</TITLE>`:                   true,
		`<Title>Page</Title>`:                   true,
		`<meta name="description" content="d">`: true,
		`<meta NAME="description" content="d">`: true,
		`<meta name='description' content='d'>`: true,
		`<meta name='DESCRIPTION' content='d'>`: true,
		`<meta name="viewport" content="w">`:    false,
		`<p>hello</p>`:                          false,
	}
	for in, want := range cases {
		if got := HasHeadDedupeTag(in); got != want {
			t.Errorf("HasHeadDedupeTag(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestBodyHeadConflictLiteralPlaceholder(t *testing.T) {
	f := &File{Body: &BodySection{Nodes: []TemplateNode{
		{Type: NodeText, Content: "<html><head><title>Site</title>"},
		{Type: NodeText, Content: HeadPlaceholder + "</head><body>"},
		{Type: NodeSlot},
	}}}
	if _, _, found := f.BodyHeadConflict(); found {
		t.Fatal("literal placeholder must not report a conflict")
	}
}

func TestBodyHeadConflictExpressionBeforePlaceholder(t *testing.T) {
	f := &File{Body: &BodySection{Nodes: []TemplateNode{
		{Type: NodeText, Content: "<html><head>"},
		{Type: NodeExpression, Content: "title"},
		{Type: NodeText, Content: HeadPlaceholder + "</head>"},
	}}}
	if _, _, found := f.BodyHeadConflict(); !found {
		t.Fatal("expression before the placeholder must report a conflict")
	}
}

func TestBodyHeadConflictMessageBeforePlaceholder(t *testing.T) {
	f := &File{Body: &BodySection{Nodes: []TemplateNode{
		{Type: NodeText, Content: "<html><head>"},
		{Type: NodeMessage, Content: "site.title"},
		{Type: NodeText, Content: HeadPlaceholder + "</head>"},
	}}}
	if _, _, found := f.BodyHeadConflict(); !found {
		t.Fatal("message before the placeholder must report a conflict")
	}
}

func TestBodyHeadConflictNoPlaceholder(t *testing.T) {
	f := &File{Body: &BodySection{Nodes: []TemplateNode{
		{Type: NodeText, Content: "<html><body>hi</body></html>"},
	}}}
	if _, _, found := f.BodyHeadConflict(); found {
		t.Fatal("missing placeholder must not report a conflict")
	}
}
