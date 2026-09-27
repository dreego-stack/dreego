package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugDuplicateComponentNameFailsGenerate(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/a/Card.dreego": "DREEFILE component ()\n<body><p>A</p></body>",
		"www/components/b/Card.dreego": "DREEFILE component ()\n<body><p>B</p></body>",
		"www/routes/+page.dreego":      "<body><@Card/></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted duplicate component names:\n%s", out)
	}
	if !strings.Contains(out, "duplicate component Card") {
		t.Fatalf("unexpected diagnostic:\n%s", out)
	}
}
