package dreegotest

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func initTaggedRepo(t *testing.T, tags ...string) string {
	t.Helper()
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
	for _, tag := range tags {
		runGit("tag", tag)
	}
	return repo
}

func TestLatestTagIgnoresModuleTags(t *testing.T) {
	repo := initTaggedRepo(t, "v0.8.0", "adapter/ssr/v0.8.0")

	if got := latestTag(repo); got != "v0.8.0" {
		t.Fatalf("latestTag = %q, want root tag v0.8.0", got)
	}
}

func TestLatestTagPrefersRepoOverEnv(t *testing.T) {
	t.Setenv("DREEGO_VERSION", "v0.9.0")
	repo := initTaggedRepo(t, "v0.8.0", "adapter/ssr/v0.8.0")

	if got := latestTag(repo); got != "v0.8.0" {
		t.Fatalf("latestTag = %q, want repository tag v0.8.0 over DREEGO_VERSION", got)
	}
}

func TestLatestTagFallsBackToEnv(t *testing.T) {
	t.Setenv("DREEGO_VERSION", "v0.9.0")
	repo := initTaggedRepo(t)

	if got := latestTag(repo); got != "v0.9.0" {
		t.Fatalf("latestTag = %q, want fallback DREEGO_VERSION v0.9.0", got)
	}
}
