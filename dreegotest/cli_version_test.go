package dreegotest

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestLatestTagIgnoresModuleTags(t *testing.T) {
	repo := t.TempDir()
	runGit := func(arguments ...string) {
		t.Helper()
		command := exec.Command("git", arguments...)
		command.Dir = repo
		command.Env = append(os.Environ(), "GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com", "GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com")
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", arguments, err, output)
		}
	}
	runGit("init", "-q")
	if err := os.WriteFile(filepath.Join(repo, "file"), []byte("test"), 0o600); err != nil {
		t.Fatal(err)
	}
	runGit("add", "file")
	runGit("commit", "-qm", "initial")
	runGit("tag", "v0.8.0")
	runGit("tag", "adapter/ssr/v0.8.0")

	if got := latestTag(repo); got != "v0.8.0" {
		t.Fatalf("latestTag = %q, want root tag v0.8.0", got)
	}
}
