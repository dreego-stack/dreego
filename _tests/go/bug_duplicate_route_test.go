package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugDuplicateRouteFlatVsIndex(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego":   `<div><p>get</p></div>`,
		"dreego/routes/index.dreego": `<div><p>index</p></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for duplicate route, got success: %s", out)
	}
	if !strings.Contains(out, "dreego/routes/get.dreego") {
		t.Fatalf("error must name the first source path, got: %s", out)
	}
	if !strings.Contains(out, "dreego/routes/index.dreego") {
		t.Fatalf("error must name the second source path, got: %s", out)
	}
}

func TestBugDuplicateRouteMethodAttrVsFile(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<go method="post">msg := "posted"</go>
<div><p>{{ msg }}</p></div>`,
		"dreego/routes/post.dreego": `<div><p>post</p></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for duplicate POST route, got success: %s", out)
	}
	if !strings.Contains(out, "POST") {
		t.Fatalf("error must name the conflicting method, got: %s", out)
	}
	if !strings.Contains(out, "dreego/routes/get.dreego") || !strings.Contains(out, "dreego/routes/post.dreego") {
		t.Fatalf("error must name both source paths, got: %s", out)
	}
}

func TestBugDuplicateRouteFormWithoutHandler(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div>
<form g-action="Missing" method="post">
    <input name="x">
    <button>OK</button>
</form>
</div>`,
		"dreego/routes/post.dreego": `<div><p>post</p></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("form without handler must not claim POST: %v\n%s", err, out)
	}
}
