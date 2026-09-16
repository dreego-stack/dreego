package tests

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestWailsAdapterLeastPrivilegeDependencies(t *testing.T) {
	t.Parallel()
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	command := exec.Command("go", "list", "-deps", "-f", "{{if not .Standard}}{{.ImportPath}}{{end}}", ".")
	command.Dir = filepath.Join(repoRoot, "adapter", "wails")
	out, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("go list -deps ./adapter/wails: %v\n%s", err, out)
	}
	for _, line := range strings.Split(string(out), "\n") {
		dep := strings.TrimSpace(line)
		if dep == "" {
			continue
		}
		if strings.HasPrefix(dep, "golang.org/x/") || strings.HasPrefix(dep, "github.com/dreego-stack/dreego") {
			continue
		}
		t.Errorf("adapter/wails links unapproved dependency %q; it must stay free of privileged Wails APIs", dep)
	}
}
