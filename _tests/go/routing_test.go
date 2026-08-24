package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestRouting404Page(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRouting500Page(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingCatchall(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingDeepNesting(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/a/b/c/d/get.dreego": `<body><p>deep</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingDeleteMethod(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/delete.dreego": `<body><p>delete works</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingDynamic(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingGetMethod(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingGroups(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingMultiSegment(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/a/get.dreego": `<server>a:=c.Param("a")</server>
<body><p>{{ a }}</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingNestedRoutes(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/about/get.dreego":       `<body><p>about page</p></body>`,
		"www/routes/users/about/get.dreego": `<body><p>users about page</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingOptional(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingPostMethod(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingPutMethod(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/put.dreego": `<body><p>put works</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingServemuxCache(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body><p>hello</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}
