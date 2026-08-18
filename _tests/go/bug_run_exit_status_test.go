package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugRunReturnsChildExitStatus(t *testing.T) {
	dir := dreegotest.ProjectDir(t, nil)
	mainGo := "package main\n\nimport \"os\"\n\nfunc main() { os.Exit(7) }\n"
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := dreegotest.RunCLI(t, dir, "run")
	if err == nil {
		t.Fatalf("run succeeded after child exit 7:\n%s", out)
	}
	if !strings.Contains(out, "server exited") {
		t.Fatalf("run output lacks child failure:\n%s", out)
	}
}
