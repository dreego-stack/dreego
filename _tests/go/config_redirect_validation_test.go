package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestGenerateRejectsInvalidRedirect(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/dreego.config.json":      `{"redirects":[{"from":"/old","to":"relative","status":301}]}`,
		"www/app/routes/+page.dreego": `<body><p>home</p></body>`,
	})
}

func TestGenerateRejectsRedirectCycle(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":      `{"redirects":[{"from":"/a","to":"/b","status":301},{"from":"/b","to":"/a","status":301}]}`,
		"www/app/routes/+page.dreego": `<body><p>home</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate to reject a redirect cycle, got: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "cycle") && !strings.Contains(strings.ToLower(out), "loop") {
		t.Fatalf("expected a cycle diagnostic, got: %s", out)
	}
}

func TestGenerateRejectsInvalidRedirectStatus(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/dreego.config.json":      `{"redirects":[{"from":"/old","to":"/new","status":200}]}`,
		"www/app/routes/+page.dreego": `<body><p>home</p></body>`,
	})
}

func TestGenerateRejectsInvalidRewrite(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/dreego.config.json":      `{"rewrites":[{"from":"/old","to":"relative"}]}`,
		"www/app/routes/+page.dreego": `<body><p>home</p></body>`,
	})
}
