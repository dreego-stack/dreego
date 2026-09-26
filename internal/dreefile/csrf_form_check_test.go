package dreefile

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestParseRouteFileEmitsCSRFWarningToStderr(t *testing.T) {
	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w
	defer func() { os.Stderr = orig }()

	_, _, perr := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<server>\n    type F struct{}\n    func Save(c dreego.Context, f F) error { return nil }\n</server>\n<body>\n<form g-action=\"Save\" method=\"post\">\n    <input name=\"name\">\n</form>\n</body>\n"))
	_ = w.Close()
	os.Stderr = orig
	if perr != nil {
		t.Fatalf("parseRouteFile must not fail on a missing csrf field: %v", perr)
	}
	out, _ := io.ReadAll(r)
	if !strings.Contains(string(out), "warning:") || !strings.Contains(string(out), "csrf_token") {
		t.Fatalf("generate must warn on stderr and keep exit 0, got: %q", string(out))
	}
}

func TestCSRFFormDiagnosticWarnsWithoutToken(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<server>\n    type F struct{}\n    func Save(c dreego.Context, f F) error { return nil }\n</server>\n<body>\n<form g-action=\"Save\" method=\"post\">\n    <input name=\"name\">\n</form>\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	diags := csrfFormDiagnostics(file.Body.Nodes, "www/routes/+page.dreego")
	if len(diags) != 1 {
		t.Fatalf("expected one csrf diagnostic, got %d: %q", len(diags), diags)
	}
	if !strings.Contains(diags[0], "csrf_token") {
		t.Fatalf("diagnostic must name csrf_token, got: %s", diags[0])
	}
	if !strings.Contains(diags[0], "403") {
		t.Fatalf("diagnostic must state the 403 rejection, got: %s", diags[0])
	}
	if !strings.Contains(diags[0], "www/routes/+page.dreego:") {
		t.Fatalf("diagnostic must carry the file location, got: %s", diags[0])
	}
	if !strings.Contains(diags[0], "Fix:") {
		t.Fatalf("diagnostic must name the next action, got: %s", diags[0])
	}
}

func TestCSRFFormDiagnosticSilentWithHiddenInput(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<body>\n<form g-action=\"Save\" method=\"post\">\n    <input type=\"hidden\" name=\"csrf_token\" value=\"{{ c.CSRFToken() }}\">\n    <input name=\"name\">\n</form>\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if diags := csrfFormDiagnostics(file.Body.Nodes, "www/routes/+page.dreego"); len(diags) != 0 {
		t.Fatalf("a form with csrf_token must not warn, got: %q", diags)
	}
}

func TestCSRFFormDiagnosticAcceptsCSRFInput(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<body>\n<form g-action=\"Save\" method=\"post\">\n    {{ c.CSRFInput() }}\n    <input name=\"name\">\n</form>\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if diags := csrfFormDiagnostics(file.Body.Nodes, "www/routes/+page.dreego"); len(diags) != 0 {
		t.Fatalf("c.CSRFInput() must satisfy the lint, got: %q", diags)
	}
}

func TestCSRFFormDiagnosticSilentWithoutGAction(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<body>\n<form method=\"post\" action=\"/x\">\n    <input name=\"name\">\n</form>\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if diags := csrfFormDiagnostics(file.Body.Nodes, "www/routes/+page.dreego"); len(diags) != 0 {
		t.Fatalf("a form without g-action must not warn, got: %q", diags)
	}
}

func TestCSRFFormDiagnosticReportsOnlyMissingForm(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<body>\n<form g-action=\"First\" method=\"post\">\n    <input name=\"a\">\n</form>\n<form g-action=\"Second\" method=\"post\">\n    <input type=\"hidden\" name=\"csrf_token\" value=\"{{ c.CSRFToken() }}\">\n</form>\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	diags := csrfFormDiagnostics(file.Body.Nodes, "www/routes/+page.dreego")
	if len(diags) != 1 {
		t.Fatalf("expected exactly one diagnostic for the missing form, got %d: %q", len(diags), diags)
	}
	if !strings.Contains(diags[0], `g-action="First"`) {
		t.Fatalf("diagnostic must name the offending action, got: %s", diags[0])
	}
}
