package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func csrfFormRoute(form string) string {
	return "<server>\n    type F struct{}\n    func Save(c dreego.Context, f F) error { return nil }\n</server>\n<body>\n" + form + "\n</body>"
}

func TestCSRFFormLintWarnsWithoutToken(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": csrfFormRoute(`<form g-action="Save" method="post">
    <input name="name">
    <button>OK</button>
</form>`),
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("generate must succeed with the warning: %v\n%s", err, out)
	}
	if !strings.Contains(out, "warning:") || !strings.Contains(out, "csrf_token") {
		t.Fatalf("expected a csrf_token warning, got:\n%s", out)
	}
	if !strings.Contains(out, "403") {
		t.Fatalf("warning must state the 403 outcome, got:\n%s", out)
	}
	if !strings.Contains(out, "www/routes/+page.dreego") || !strings.Contains(out, "Fix:") {
		t.Fatalf("warning must carry the source location and fix, got:\n%s", out)
	}
}

func TestCSRFFormLintSilentWithToken(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": csrfFormRoute(`<form g-action="Save" method="post">
    <input type="hidden" name="csrf_token" value="{{ c.CSRFToken() }}">
    <input name="name">
    <button>OK</button>
</form>`),
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if strings.Contains(out, "csrf_token") {
		t.Fatalf("a form with a token must not warn, got:\n%s", out)
	}
}
