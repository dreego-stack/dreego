package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugDuplicateRoutePageVsIndex(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": `<body><p>get</p></body>`,
		"www/routes/index.dreego": `<body><p>index</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for duplicate route, got success: %s", out)
	}
	if !strings.Contains(out, "www/routes/+page.dreego") {
		t.Fatalf("error must name the first source path, got: %s", out)
	}
	if !strings.Contains(out, "www/routes/index.dreego") {
		t.Fatalf("error must name the second source path, got: %s", out)
	}
}

func TestNamedFileDoesNotClaimDirectoryRoute(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego": `<server method="post">msg := "posted"</server>
<body method="post"><p>{{ msg }}</p></body>`,
		"www/routes/profile.dreego": `<body><p>profile</p></body>`,
	})
	code, body, _ := c.Request(t, "POST", "/", "", nil)
	if code != 200 || !strings.Contains(body, "posted") {
		t.Fatalf("POST / = %d %q, want directory route", code, body)
	}
	code, body = c.Get(t, "/profile")
	if code != 200 || !strings.Contains(body, "profile") {
		t.Fatalf("GET /profile = %d %q, want named route", code, body)
	}
}

func TestBugDuplicateRouteFormWithoutHandler(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": `<body>
<form g-action="Missing" method="post">
    <input name="x">
    <button>OK</button>
</form>
</body>`,
		"www/routes/profile.dreego": `<body><p>profile</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("form without handler must not claim POST: %v\n%s", err, out)
	}
}
