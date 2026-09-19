package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

var internalModulePattern = regexp.MustCompile(`(?m)^[ \t]*(?:require[ \t]+|replace[ \t]+)?(github\.com/dreego-stack/dreego(?:/[^\s]+)?)[ \t]+([^\s]+)`)

var semanticVersionPattern = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

func internalModuleVersions(body string) map[string]string {
	versions := make(map[string]string)
	for _, match := range internalModulePattern.FindAllStringSubmatch(body, -1) {
		versions[match[1]] = match[2]
	}
	return versions
}

func coordinatedVersion(bodies map[string]string) (string, error) {
	version := ""
	for path, body := range bodies {
		for module, found := range internalModuleVersions(body) {
			if !semanticVersionPattern.MatchString(found) {
				return "", fmt.Errorf("%s: %s uses invalid version %q", path, module, found)
			}
			if version == "" {
				version = found
				continue
			}
			if found != version {
				return "", fmt.Errorf("%s: %s uses %s, coordinated version is %s", path, module, found, version)
			}
		}
	}
	if version == "" {
		return "", fmt.Errorf("no internal module requirement found")
	}
	return version, nil
}

func TestModuleBoundaries(t *testing.T) {
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	modules := map[string]string{
		"go.mod":               "module github.com/dreego-stack/dreego\n",
		"core/go.mod":          "module github.com/dreego-stack/dreego/core\n",
		"adapter/ssr/go.mod":   "module github.com/dreego-stack/dreego/adapter/ssr\n",
		"adapter/wails/go.mod": "module github.com/dreego-stack/dreego/adapter/wails\n",
		"dreegotest/go.mod":    "module github.com/dreego-stack/dreego/dreegotest\n",
		"cmd/dreego/go.mod":    "module github.com/dreego-stack/dreego/cmd/dreego\n",
	}
	for path, declaration := range modules {
		contents, err := os.ReadFile(filepath.Join(repoRoot, path))
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		if !strings.Contains(string(contents), declaration) {
			t.Errorf("%s does not declare %q", path, strings.TrimSpace(declaration))
		}
		if strings.Contains(string(contents), "replace github.com/dreego-stack/dreego") {
			t.Errorf("published module %s contains a repository-local replacement", path)
		}
	}
	requirements := map[string][]string{
		"core/go.mod":          {"github.com/dreego-stack/dreego"},
		"adapter/ssr/go.mod":   {"github.com/dreego-stack/dreego", "github.com/dreego-stack/dreego/core"},
		"adapter/wails/go.mod": {"github.com/dreego-stack/dreego/core"},
		"dreegotest/go.mod":    {"github.com/dreego-stack/dreego", "github.com/dreego-stack/dreego/core"},
		"cmd/dreego/go.mod":    {"github.com/dreego-stack/dreego"},
	}
	bodies := make(map[string]string, len(requirements)+1)
	for path, expected := range requirements {
		contents, err := os.ReadFile(filepath.Join(repoRoot, path))
		if err != nil {
			t.Fatal(err)
		}
		bodies[path] = string(contents)
		versions := internalModuleVersions(bodies[path])
		for _, requirement := range expected {
			if _, ok := versions[requirement]; !ok {
				t.Errorf("%s is missing coordinated requirement %q", path, requirement)
			}
		}
	}
	workspace, err := os.ReadFile(filepath.Join(repoRoot, "go.work"))
	if err != nil {
		t.Fatal(err)
	}
	bodies["go.work"] = string(workspace)
	if version, err := coordinatedVersion(bodies); err != nil {
		t.Errorf("coordinated module version mismatch: %v", err)
	} else if !semanticVersionPattern.MatchString(version) {
		t.Errorf("coordinated module version %q is not a semantic version", version)
	}
	if info, err := os.Stat(filepath.Join(repoRoot, "_docs")); err != nil || !info.IsDir() {
		t.Error("central documentation directory _docs is missing")
	}
	for _, path := range []string{
		"core/_docs",
		"adapter/ssr/_docs",
		"adapter/wails/_docs",
		"dreegotest/_docs",
		"cmd/dreego/_docs",
	} {
		if _, err := os.Stat(filepath.Join(repoRoot, path)); err == nil {
			t.Errorf("module documentation directory %s must live centrally in _docs", path)
		}
	}
	for _, path := range []string{"core/ssr", "cli/dreego", "target"} {
		if legacyPathExists(filepath.Join(repoRoot, path)) {
			t.Errorf("removed v0.7 path %s still exists", path)
		}
	}
	oldImports := []string{
		"github.com/dreego-stack/dreego/core/ssr",
		"github.com/dreego-stack/dreego/cli/dreego",
	}
	if err := filepath.WalkDir(repoRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
			if walkErr == nil && entry.IsDir() && (entry.Name() == ".worktrees" || entry.Name() == ".git") {
				return filepath.SkipDir
			}
			return walkErr
		}
		if filepath.Base(path) == "module_boundaries_test.go" {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, oldImport := range oldImports {
			if strings.Contains(string(contents), oldImport) {
				t.Errorf("removed import %q remains in %s", oldImport, path)
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestCoordinatedVersionRejectsMismatch(t *testing.T) {
	if _, err := coordinatedVersion(map[string]string{
		"core/go.mod":        "module github.com/dreego-stack/dreego/core\n\nrequire github.com/dreego-stack/dreego v0.9.0\n",
		"adapter/ssr/go.mod": "module github.com/dreego-stack/dreego/adapter/ssr\n\nrequire github.com/dreego-stack/dreego/core v0.8.0\n",
	}); err == nil {
		t.Fatal("coordinatedVersion must reject mixed v0.9.0 and v0.8.0 requirements")
	}
}

func TestCoordinatedVersionRejectsInvalidSemanticVersion(t *testing.T) {
	if _, err := coordinatedVersion(map[string]string{
		"core/go.mod": "module github.com/dreego-stack/dreego/core\n\nrequire github.com/dreego-stack/dreego v0.9\n",
	}); err == nil {
		t.Fatal("coordinatedVersion must reject non-semantic module versions")
	}
}

func TestCoordinatedVersionAcceptsSingleVersion(t *testing.T) {
	version, err := coordinatedVersion(map[string]string{
		"core/go.mod":        "module github.com/dreego-stack/dreego/core\n\nrequire github.com/dreego-stack/dreego v0.9.0\n",
		"adapter/ssr/go.mod": "module github.com/dreego-stack/dreego/adapter/ssr\n\nrequire github.com/dreego-stack/dreego/core v0.9.0\n",
		"go.work":            "replace github.com/dreego-stack/dreego v0.9.0 => .\n\nreplace github.com/dreego-stack/dreego/core v0.9.0 => ./core\n",
	})
	if err != nil {
		t.Fatalf("coordinatedVersion rejected a consistent version set: %v", err)
	}
	if version != "v0.9.0" {
		t.Fatalf("coordinatedVersion = %q, want v0.9.0", version)
	}
}

func legacyPathExists(path string) bool {
	entries, err := os.ReadDir(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		return true
	}
	for _, entry := range entries {
		if entry.Name() != ".DS_Store" {
			return true
		}
	}
	return false
}

func TestWailsAdapterKeepsApplicationOwnershipExplicit(t *testing.T) {
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	adapterModule, err := os.ReadFile(filepath.Join(repoRoot, "adapter/wails/go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(adapterModule), "github.com/wailsapp/wails") {
		t.Fatal("adapter/wails must not depend on Wails")
	}
	mainSource, err := os.ReadFile(filepath.Join(repoRoot, "demo/demo-wailsv3/main.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, ownership := range []string{"application.New", "application.NewService", "Window.NewWithOptions", "wailsApp.Run"} {
		if !strings.Contains(string(mainSource), ownership) {
			t.Errorf("demo main.go does not visibly own %s", ownership)
		}
	}
}
