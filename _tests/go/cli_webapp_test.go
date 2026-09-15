package tests

import (
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func webAppTemplateFiles(t *testing.T) map[string]string {
	t.Helper()
	dir := t.TempDir()
	if out, err := dreegotest.RunCLI(t, dir, "init", ".", "-t", "web-app"); err != nil {
		t.Fatalf("init -t web-app: %v\n%s", err, out)
	}
	root := filepath.Join(dir, "www")
	files := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return walkErr
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = string(data)
		return nil
	})
	if err != nil {
		t.Fatalf("walk web-app scaffold: %v", err)
	}
	return files
}

func TestCLIWebAppTemplateServes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, webAppTemplateFiles(t))

	for _, tc := range []struct {
		path  string
		title string
		nav   string
	}{
		{path: "/", title: "Notes", nav: `<a href="/" aria-current="page" class="active">Home</a>`},
		{path: "/dashboard", title: "Dashboard", nav: `<a href="/dashboard" aria-current="page" class="active">Dashboard</a>`},
	} {
		code, body := c.Get(t, tc.path)
		if code != 200 {
			t.Fatalf("GET %s = %d, want 200", tc.path, code)
		}
		if got := strings.Count(body, "<title"); got != 1 {
			t.Fatalf("GET %s must render exactly one <title>, got %d\n---\n%s", tc.path, got, body)
		}
		if !strings.Contains(body, "<title>"+tc.title+" ") {
			t.Fatalf("GET %s missing its own title %q\n---\n%s", tc.path, tc.title, body)
		}
		if !strings.Contains(body, tc.nav) {
			t.Fatalf("GET %s missing active nav link %q\n---\n%s", tc.path, tc.nav, body)
		}
		if got := strings.Count(body, `aria-current="page"`); got != 1 {
			t.Fatalf("GET %s must mark exactly one nav link current, got %d\n---\n%s", tc.path, got, body)
		}
	}

	_, dashboard := c.Get(t, "/dashboard")
	if !strings.Contains(dashboard, "<table>") || !strings.Contains(dashboard, "<caption>Tasks</caption>") {
		t.Fatalf("dashboard missing its captioned table\n---\n%s", dashboard)
	}
	if !strings.Contains(dashboard, `<th scope="col">Task</th>`) {
		t.Fatalf("dashboard table header missing scope\n---\n%s", dashboard)
	}

	_, home := c.Get(t, "/")
	if strings.Contains(home, `name="csrf_token"`) {
		t.Fatalf("home must not ship an inert csrf_token field\n---\n%s", home)
	}

	code, _, headers := c.Request(t, "POST", "/", "title=First+note", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 303 {
		t.Fatalf("POST / = %d, want 303 redirect\n---\n%s", code, headers.Get("Location"))
	}
	if loc := headers.Get("Location"); loc != "/" {
		t.Fatalf("POST / redirect = %q, want /", loc)
	}
	_, home = c.Get(t, "/")
	if !strings.Contains(home, "First note") {
		t.Fatalf("saved note missing after POST\n---\n%s", home)
	}

	code, body, _ := c.Request(t, "POST", "/", "title=", map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 200 {
		t.Fatalf("POST / with empty title = %d, want 200 re-render", code)
	}
	if !strings.Contains(body, "is required") {
		t.Fatalf("empty title must render the validation error\n---\n%s", body)
	}
	if !strings.Contains(body, `class="error"`) {
		t.Fatalf("validation error must be marked up\n---\n%s", body)
	}
}

const webAppRaceTest = `package routes_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
	"t/www/routes"
)

func TestGeneratedNotesConcurrent(t *testing.T) {
	app := dreego.New()
	if err := routes.Register(app); err != nil {
		t.Fatal(err)
	}
	handler := app.Handler()
	var wg sync.WaitGroup
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(fmt.Sprintf("title=note-%d", i)))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusSeeOther {
				t.Errorf("POST / = %d, want 303", rec.Code)
			}
		}(i)
	}
	for i := 0; i < 25; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK {
				t.Errorf("GET / = %d, want 200", rec.Code)
			}
		}()
	}
	wg.Wait()
}
`

func TestCLIWebAppNotesConcurrentRace(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, webAppTemplateFiles(t))
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	testFile := filepath.Join(dir, "www", "routes", "notes_race_test.go")
	if err := os.WriteFile(testFile, []byte(webAppRaceTest), 0644); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("go", "test", "-mod=mod", "-race", "./www/routes", "-run", "TestGeneratedNotesConcurrent", "-count=1")
	command.Dir = dir
	command.Env = append(os.Environ(), "CGO_ENABLED=1")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated -race concurrency test failed: %v\n%s", err, output)
	}
}
