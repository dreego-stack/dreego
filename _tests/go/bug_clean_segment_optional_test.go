package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugCleanSegmentOptional(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/[[opt]]/get.dreego": `<body>optional</body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected generate failure for optional segment, got success: %s", out)
	}
	if !strings.Contains(out, "[[opt]]") {
		t.Fatalf("error must name the optional segment, got: %s", out)
	}
	if !strings.Contains(out, "www/routes/[[opt]]") {
		t.Fatalf("error must include the source path, got: %s", out)
	}
}
