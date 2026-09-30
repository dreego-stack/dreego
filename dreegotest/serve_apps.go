package dreegotest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// MultiApp is a locally built multi-app server: one website root with several
// app subdirectories, each served on its own port. Apps maps an app package
// name (the app directory name) to its client.
type MultiApp struct {
	Apps map[string]*Client
}

// App returns the client for one app package, failing when it does not exist.
func (m *MultiApp) App(t *testing.T, name string) *Client {
	t.Helper()
	c, ok := m.Apps[name]
	if !ok {
		t.Fatalf("no app %q; available: %v", name, appNames(m.Apps))
	}
	return c
}

func appNames(apps map[string]*Client) []string {
	names := make([]string, 0, len(apps))
	for name := range apps {
		names = append(names, name)
	}
	return names
}

// ServeApps builds a project with one website root and one app per entry in
// apps (a relative app directory path, e.g. "www"), starts the binary, and
// returns a client per app. Each app is served on its own port. files are the
// usual relative dreego file map; the caller provides a main that starts every
// app. This is the multi-app counterpart to Serve.
func ServeApps(t *testing.T, files map[string]string, apps ...string) *MultiApp {
	t.Helper()
	dir := t.TempDir()
	repoRoot, err := RepoRoot()
	if err != nil {
		t.Fatalf("ServeApps: %v", err)
	}
	goMod := testModuleFile(repoRoot, true)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("ServeApps: write go.mod: %v", err)
	}
	copyModuleSum(t, dir, repoRoot)

	imports := ""
	newApps := ""
	last := ""
	for i, app := range apps {
		env := fmt.Sprintf("APP_%d_ADDR", i)
		imports += fmt.Sprintf("\tapp%d %q\n", i, "t/"+app)
		newApps += fmt.Sprintf("\ta%d := dreego.New(app%d.App)\n", i, i)
		if i < len(apps)-1 {
			newApps += fmt.Sprintf("\tgo ssr.Listen(a%d, os.Getenv(%q))\n", i, env)
		} else {
			last = fmt.Sprintf("a%d", i)
		}
	}
	mainGo := fmt.Sprintf(`package main

import (
	"os"

	ssr "github.com/dreego-stack/dreego/adapter/ssr"
	dreego "github.com/dreego-stack/dreego/core"
%s)

func main() {
%s	if err := ssr.Listen(%s, os.Getenv("APP_%d_ADDR")); err != nil {
		os.Exit(1)
	}
}
`, imports, newApps, last, len(apps)-1)
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatalf("ServeApps: write main.go: %v", err)
	}
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatalf("ServeApps: mkdir: %v", err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatalf("ServeApps: write %s: %v", path, err)
		}
	}
	ensureConfig(t, dir, files)
	if _, err := RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("ServeApps: generate: %v", err)
	}
	bin := filepath.Join(dir, "server")
	build := exec.Command("go", "build", "-mod=mod", "-o", bin, ".")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("ServeApps: go build: %v\n%s", err, out)
	}

	ports := make([]int, len(apps))
	env := os.Environ()
	for i := range apps {
		ports[i] = FreePort(t)
		env = append(env, fmt.Sprintf("APP_%d_ADDR=127.0.0.1:%d", i, ports[i]))
	}
	proc := exec.Command(bin)
	proc.Dir = dir
	proc.Env = env
	if err := proc.Start(); err != nil {
		t.Fatalf("ServeApps: start: %v", err)
	}
	t.Cleanup(func() {
		proc.Process.Kill()
		proc.Wait()
	})

	result := &MultiApp{Apps: map[string]*Client{}}
	for i, app := range apps {
		WaitForPort(t, ports[i])
		name := filepath.Base(app)
		result.Apps[name] = &Client{base: fmt.Sprintf("http://127.0.0.1:%d", ports[i]), jar: newJar()}
	}
	return result
}
