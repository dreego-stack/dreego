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
		"www/routes/a/b/c/d/get.dreego": `<div><p>deep</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingDeleteMethod(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/delete.dreego": `<div><p>delete works</p></div>`,
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
		"www/routes/a/get.dreego": `<go>a:=c.Param("a")</go>
<div><p>{{ a }}</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingNestedRoutes(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/about/get.dreego":       `<div><p>about page</p></div>`,
		"www/routes/users/about/get.dreego": `<div><p>users about page</p></div>`,
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
		"www/routes/put.dreego": `<div><p>put works</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestRoutingServemuxCache(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<div><p>hello</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}
