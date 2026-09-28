package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugNewPreservesFullModulePath(t *testing.T) {
	t.Parallel()
	parent := t.TempDir()
	out, err := dreegotest.RunCLIIn(t, parent, dreegotest.CLIBin(t), "new", "github.com/me/myapp")
	if err != nil {
		t.Fatalf("new failed: %v\n%s", err, out)
	}
	data, err := os.ReadFile(filepath.Join(parent, "github.com", "me", "myapp", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "module github.com/me/myapp\n") {
		t.Fatalf("full module path was not preserved:\n%s", data)
	}
}
