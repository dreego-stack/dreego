package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugDuplicateRouteFlatVsIndex(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego":   `<body><p>get</p></body>`,
		"www/routes/index.dreego": `<body><p>index</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for duplicate route, got success: %s", out)
	}
	if !strings.Contains(out, "www/routes/get.dreego") {
		t.Fatalf("error must name the first source path, got: %s", out)
	}
	if !strings.Contains(out, "www/routes/index.dreego") {
		t.Fatalf("error must name the second source path, got: %s", out)
	}
}

func TestBugDuplicateRouteMethodAttrVsFile(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<server method="post">msg := "posted"</server>
<body><p>{{ msg }}</p></body>`,
		"www/routes/post.dreego": `<body><p>post</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for duplicate POST route, got success: %s", out)
	}
	if !strings.Contains(out, "POST") {
		t.Fatalf("error must name the conflicting method, got: %s", out)
	}
	if !strings.Contains(out, "www/routes/get.dreego") || !strings.Contains(out, "www/routes/post.dreego") {
		t.Fatalf("error must name both source paths, got: %s", out)
	}
}

func TestBugDuplicateRouteFormWithoutHandler(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body>
<form g-action="Missing" method="post">
    <input name="x">
    <button>OK</button>
</form>
</body>`,
		"www/routes/post.dreego": `<body><p>post</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("form without handler must not claim POST: %v\n%s", err, out)
	}
}
