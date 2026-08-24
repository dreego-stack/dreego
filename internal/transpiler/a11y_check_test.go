package transpiler

import (
	"strings"
	"testing"
)

func TestA11yCheckImageWithoutAlt(t *testing.T) {
	src := "<body>\n    <img src=\"/logo.png\">\n</body>\n"
	f := parseFile(t, src)
	setNodeSource(f.Body.Nodes, "dreego/routes/get.dreego", 0)
	d := a11yCheck(f.Body.Nodes)
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
	src := "<body>\n<form>\n    <input name=\"email\" type=\"email\">\n</form>\n</body>\n"
	f := parseFile(t, src)
	setNodeSource(f.Body.Nodes, "dreego/routes/get.dreego", 0)
	d := a11yCheck(f.Body.Nodes)
	if len(d) == 0 {
		t.Fatal("expected an accessibility diagnostic for <input> without label")
	}
	if !strings.Contains(d[0].String(), "label") {
		t.Fatalf("diagnostic must mention label, got %q", d[0].String())
	}
}

func TestA11yCheckExplicitAltAndLabel(t *testing.T) {
	src := "<body>\n    <img src=\"/logo.png\" alt=\"Dreego logo\">\n    <label for=\"email\">Email</label>\n    <input id=\"email\" name=\"email\" type=\"email\">\n</body>\n"
	f := parseFile(t, src)
	setNodeSource(f.Body.Nodes, "dreego/routes/get.dreego", 0)
	d := a11yCheck(f.Body.Nodes)
	if len(d) != 0 {
		t.Fatalf("expected no diagnostics for accessible markup, got %q", d)
	}
}

func TestA11yCheckLabelForMatchesID(t *testing.T) {
	src := "<body>\n    <label for=\"email\">Email</label>\n    <input id=\"email\" name=\"email\">\n</body>\n"
	f := parseFile(t, src)
	setNodeSource(f.Body.Nodes, "dreego/routes/get.dreego", 0)
	if d := a11yCheck(f.Body.Nodes); len(d) != 0 {
		t.Fatalf("label[for] must count as an association, got %q", d)
	}
}

func TestA11yCheckFormGetsDiagnostics(t *testing.T) {
	src := "<body>\n<form>\n    <input name=\"email\" type=\"email\">\n</form>\n</body>\n"
	f := parseFile(t, src)
	setNodeSource(f.Body.Nodes, "dreego/routes/get.dreego", 0)
	d := a11yDiagnostics(f.Body.Nodes)
	if len(d) != 1 {
		t.Fatalf("expected exactly one diagnostic, got %q", d)
	}
	if !strings.Contains(d[0], "input") {
		t.Fatalf("diagnostic must name input, got %q", d[0])
	}
}

func TestA11yCheckIsCaseInsensitive(t *testing.T) {
	src := `<body><IMG SRC="/logo.png" ALT="Logo"><LABEL FOR="email">Email</LABEL><INPUT ID="email" TYPE="EMAIL"></body>`
	f := parseFile(t, src)
	setNodeSource(f.Body.Nodes, "route.dreego", 0)
	if d := a11yCheck(f.Body.Nodes); len(d) != 0 {
		t.Fatalf("uppercase markup produced diagnostics: %q", d)
	}
}

func TestA11yCheckRecognizesWrappedLabel(t *testing.T) {
	src := `<body><label>Email <input name="email" type="email"></label></body>`
	f := parseFile(t, src)
	setNodeSource(f.Body.Nodes, "route.dreego", 0)
	if d := a11yCheck(f.Body.Nodes); len(d) != 0 {
		t.Fatalf("wrapped label produced diagnostics: %q", d)
	}
}
