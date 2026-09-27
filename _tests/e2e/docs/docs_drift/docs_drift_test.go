package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

var staleTestCounts = []string{
	"147 Integration Tests",
	"Implemented | 40",
	"Planned (v0.0.7) | 30",
	"Named Slots (v0.0.8) | 4",
}

var globalAppAPIs = []*regexp.Regexp{
	regexp.MustCompile(`dreego\.Listen\(`),
	regexp.MustCompile(`dreego\.Register\(`),
	regexp.MustCompile(`dreego\.SetSessionStore\(`),
	regexp.MustCompile(`dreego\.SetCSRF\(`),
	regexp.MustCompile(`dreego\.SetLogging\(`),
	regexp.MustCompile(`dreego\.SetCSP\(`),
	regexp.MustCompile(`dreego\.SetErrorHandler\(`),
	regexp.MustCompile(`dreego\.SetReady\(`),
}

var blobLink = regexp.MustCompile(`https://github\.com/dreego-stack/dreego/blob/main/([^)\s]+)`)

func docsFiles(t *testing.T) []string {
	t.Helper()
	root, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	docs, err := filepath.Glob(filepath.Join(root, "_docs", "*.md"))
	if err != nil {
		t.Fatalf("glob _docs: %v", err)
	}
	return append(docs, filepath.Join(root, "README.md"))
}

func TestDocsNoStaleTestCounts(t *testing.T) {
	t.Parallel()
	for _, f := range docsFiles(t) {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		text := string(data)
		for _, pat := range staleTestCounts {
			if strings.Contains(text, pat) {
				t.Errorf("%s still claims %q — test counts drift; describe the layout instead", f, pat)
			}
		}
		if strings.Contains(text, "⬜") {
			t.Errorf("%s still lists planned tests with ⬜ — remove planned rows", f)
		}
	}
}

func TestDocsExamplesUseExplicitApp(t *testing.T) {
	t.Parallel()
	for _, f := range docsFiles(t) {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		for _, re := range globalAppAPIs {
			if m := re.FindString(string(data)); m != "" {
				t.Errorf("%s uses the removed global API call %q — use the explicit App (app := dreego.New(); gen.Register(app))", f, m)
			}
		}
	}
}

func TestDocsLinksResolve(t *testing.T) {
	t.Parallel()
	root, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	for _, f := range docsFiles(t) {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		for _, m := range blobLink.FindAllStringSubmatch(string(data), -1) {
			rel := m[1]
			if i := strings.IndexAny(rel, "#?"); i >= 0 {
				rel = rel[:i]
			}
			target := filepath.Join(root, filepath.FromSlash(rel))
			if _, err := os.Stat(target); err != nil {
				t.Errorf("%s links to %s which does not exist in the repo", f, m[0])
			}
		}
	}
}

func TestConfigDocWebsiteRoot(t *testing.T) {
	t.Parallel()
	root, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "_docs", "config.md"))
	if err != nil {
		t.Fatalf("read config.md: %v", err)
	}
	text := string(data)
	if strings.Contains(text, "in the project root") {
		t.Error("config.md still says dreego.config.json lives in the project root; it lives in the website root (www/ by default)")
	}
	if !strings.Contains(text, "website root") {
		t.Error("config.md must name the website root as the location of dreego.config.json")
	}
}

func TestDocsInstallPathUsesCmdDreego(t *testing.T) {
	t.Parallel()
	for _, f := range docsFiles(t) {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		text := string(data)
		if strings.Contains(text, "go install github.com/dreego-stack/dreego@") {
			t.Errorf("%s documents the module-root install path; the CLI install path is github.com/dreego-stack/dreego/cmd/dreego@…", f)
		}
	}
}
