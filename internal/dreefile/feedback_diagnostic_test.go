package dreefile

import (
	"strings"
	"testing"
)

func TestBodyAttrDiagnosticWarnsOnRouteBodyAttributes(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<body x-data=\"app()\" x-init=\"init()\">\n<p>hi</p>\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, ok := bodyAttrDiagnostic(file, "www/routes/+page.dreego", 0)
	if !ok {
		t.Fatal("attributes on the route <body> tag must warn")
	}
	if !strings.Contains(out, "x-data") || !strings.Contains(out, "not applied") {
		t.Fatalf("diagnostic must name the dropped attributes, got: %s", out)
	}
}

func TestBodyAttrDiagnosticSilentWithoutAttributes(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<body lang=\"md\">\n# hi\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, ok := bodyAttrDiagnostic(file, "www/routes/+page.dreego", 0); ok {
		t.Fatal("section directives must not warn")
	}
}

func TestAlpineCSPDiagnosticWarnsOnAlpineDirectives(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<body>\n<div x-data=\"{n: 0}\" x-text=\"n\"></div>\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, ok := alpineCSPDiagnostic(file.Body.Nodes, "www/routes/+page.dreego")
	if !ok {
		t.Fatal("Alpine directives must warn about the default CSP")
	}
	if !strings.Contains(out, "unsafe-eval") {
		t.Fatalf("diagnostic must mention unsafe-eval, got: %s", out)
	}
}

func TestAlpineCSPDiagnosticSilentForPlainHTML(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/+page.dreego",
		[]byte("<body>\n<p class=\"hint\">plain</p>\n</body>\n"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, ok := alpineCSPDiagnostic(file.Body.Nodes, "www/routes/+page.dreego"); ok {
		t.Fatal("plain HTML must not warn about CSP")
	}
}
