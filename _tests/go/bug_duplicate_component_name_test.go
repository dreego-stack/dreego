package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugDuplicateComponentNameFailsGenerate(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/a/Card.dreego": "Component Card ()\n<div><p>A</p></div>",
		"www/components/b/Card.dreego": "Component Card ()\n<div><p>B</p></div>",
		"www/routes/get.dreego":        "<div><@Card/></div>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted duplicate component names:\n%s", out)
	}
	if !strings.Contains(out, "duplicate component Card") {
		t.Fatalf("unexpected diagnostic:\n%s", out)
	}
}
