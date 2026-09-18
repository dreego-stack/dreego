package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestLayoutChainMissingTargetFailsGenerate(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/layouts/default.dreego":             "<body><html><body><nav>Base</nav>{#slot}</body></html></body>",
		"www/routes/admin/layouts/layout.dreego": "DREEFILE layout\n\nLAYOUT \"www/layouts/missing.dreego\"\n\n<body><html><body><nav>Admin</nav>{#slot}</body></html></body>",
		"www/routes/admin/+page.dreego":          "<body><p>Admin page</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted a missing LAYOUT target:\n%s", out)
	}
	for _, want := range []string{"layout.dreego", "www/layouts/missing.dreego", "3:8"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestLayoutChainCycleFailsGenerate(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/layouts/default.dreego":             "DREEFILE layout\n\nLAYOUT \"www/routes/admin/layouts/layout.dreego\"\n\n<body><html><body>{#slot}</body></html></body>",
		"www/routes/admin/layouts/layout.dreego": "DREEFILE layout\n\nLAYOUT \"www/layouts/default.dreego\"\n\n<body><html><body><nav>Admin</nav>{#slot}</body></html></body>",
		"www/routes/admin/+page.dreego":          "<body><p>Admin page</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted a layout cycle:\n%s", out)
	}
	for _, want := range []string{"cycle", "default.dreego", "layout.dreego"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestLayoutChainValidRendersOutermost(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/layouts/default.dreego":             "<body><html><body><nav>Base</nav>{#slot}</body></html></body>",
		"www/routes/admin/layouts/layout.dreego": "DREEFILE layout\n\nLAYOUT \"www/layouts/default.dreego\"\n\n<body><html><body><nav>Admin</nav>{#slot}</body></html></body>",
		"www/routes/admin/+page.dreego":          "<body><p>Admin page</p></body>",
	})
	code, body := c.Get(t, "/admin")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	for _, want := range []string{"<nav>Base</nav>", "Admin page"} {
		if !strings.Contains(body, want) {
			t.Fatalf("response missing %q, got: %s", want, body)
		}
	}
}
