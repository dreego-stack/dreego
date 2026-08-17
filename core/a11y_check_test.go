package core

import (
	"strings"
	"testing"
)

func TestA11yCheckImageWithoutAlt(t *testing.T) {
	src := "<div>\n    <img src=\"/logo.png\">\n</div>\n"
	f := parseFile(t, src)
	setNodeSource(f.Template.Nodes, "dreego/routes/get.dreego", 0)
	d := a11yCheck(f.Template.Nodes)
	if len(d) == 0 {
		t.Fatal("expected an accessibility diagnostic for <img> without alt")
	}
	if !strings.Contains(d[0].String(), "dreego/routes/get.dreego:2:5") {
		t.Fatalf("diagnostic must carry file:line:col, got %q", d[0].String())
	}
	if !strings.Contains(d[0].String(), "img") || !strings.Contains(d[0].String(), "alt") {
		t.Fatalf("diagnostic must name img and alt, got %q", d[0].String())
	}
	if !strings.Contains(d[0].String(), "Fix:") {
		t.Fatalf("diagnostic must end with a next action, got %q", d[0].String())
	}
}

func TestA11yCheckInputWithoutLabel(t *testing.T) {
	src := "<div>\n<form>\n    <input name=\"email\" type=\"email\">\n</form>\n</div>\n"
	f := parseFile(t, src)
	setNodeSource(f.Template.Nodes, "dreego/routes/get.dreego", 0)
	d := a11yCheck(f.Template.Nodes)
	if len(d) == 0 {
		t.Fatal("expected an accessibility diagnostic for <input> without label")
	}
	if !strings.Contains(d[0].String(), "label") {
		t.Fatalf("diagnostic must mention label, got %q", d[0].String())
	}
}

func TestA11yCheckExplicitAltAndLabel(t *testing.T) {
	src := "<div>\n    <img src=\"/logo.png\" alt=\"Dreego logo\">\n    <label for=\"email\">Email</label>\n    <input id=\"email\" name=\"email\" type=\"email\">\n</div>\n"
	f := parseFile(t, src)
	setNodeSource(f.Template.Nodes, "dreego/routes/get.dreego", 0)
	d := a11yCheck(f.Template.Nodes)
	if len(d) != 0 {
		t.Fatalf("expected no diagnostics for accessible markup, got %q", d)
	}
}

func TestA11yCheckLabelForMatchesID(t *testing.T) {
	src := "<div>\n    <label for=\"email\">Email</label>\n    <input id=\"email\" name=\"email\">\n</div>\n"
	f := parseFile(t, src)
	setNodeSource(f.Template.Nodes, "dreego/routes/get.dreego", 0)
	if d := a11yCheck(f.Template.Nodes); len(d) != 0 {
		t.Fatalf("label[for] must count as an association, got %q", d)
	}
}

func TestA11yCheckFormGetsDiagnostics(t *testing.T) {
	src := "<div>\n<form>\n    <input name=\"email\" type=\"email\">\n</form>\n</div>\n"
	f := parseFile(t, src)
	setNodeSource(f.Template.Nodes, "dreego/routes/get.dreego", 0)
	d := a11yDiagnostics(f.Template.Nodes)
	if len(d) != 1 {
		t.Fatalf("expected exactly one diagnostic, got %q", d)
	}
	if !strings.Contains(d[0], "input") {
		t.Fatalf("diagnostic must name input, got %q", d[0])
	}
}
