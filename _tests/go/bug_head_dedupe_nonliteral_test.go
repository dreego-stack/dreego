package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugHeadDedupeNonLiteralPrefixDiagnostic(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/layouts/default.dreego": `<body>
<html>
<head>
<title>Site</title>
{{ layoutTitle }}
{#head}
</head>
<body><main>{#slot}</main></body>
</html>
</body>`,
		"www/routes/+page.dreego": `<head><title>Page</title></head>
<body><h1>Page</h1></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("generate must still succeed while warning: %v\n%s", err, out)
	}
	if !strings.Contains(out, "warning:") {
		t.Fatalf("expected a generate-time dedupe warning, got:\n%s", out)
	}
	if !strings.Contains(out, "dedupe") {
		t.Fatalf("warning must name the disabled dedupe, got:\n%s", out)
	}
	if !strings.Contains(out, "www/layouts/default.dreego") {
		t.Fatalf("warning must carry the layout source location, got:\n%s", out)
	}
	if !strings.Contains(out, "Fix:") {
		t.Fatalf("warning must name the next action, got:\n%s", out)
	}
}

func TestBugHeadDedupeLiteralPrefixNoDiagnostic(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/layouts/default.dreego": `<body>
<html>
<head>
<title>Site</title>
{#head}
</head>
<body><main>{#slot}</main></body>
</html>
</body>`,
		"www/routes/+page.dreego": `<head><title>Page</title></head>
<body><h1>Page</h1></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if strings.Contains(out, "dedupe is disabled") {
		t.Fatalf("literal layout head must not warn, got:\n%s", out)
	}
}
