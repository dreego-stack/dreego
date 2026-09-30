package tests

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestMultiAppTwoPorts is the black-box gate for the "one website root, many
// apps" model: a website root with two apps and a shared root component builds
// through the CLI and serves each app on its own port.
func TestMultiAppTwoPorts(t *testing.T) {
	t.Parallel()
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	dir := t.TempDir()

	files := map[string]string{
		"www/dreego.config.json":           `{"logging":{"enabled":false}}`,
		"www/components/Badge.dreego":      "DREEFILE component (label string)\n<body><span class=\"badge\">{{ label }}</span></body>",
		"www/blog/components/Local.dreego": "DREEFILE component (label string)\n<body><em class=\"local\">{{ label }}</em></body>",
		"www/app/routes/+page.dreego":      "COMPONENT \"www/components\" IMPORT { Badge }\n<body><h1>www home</h1><@Badge label=\"shared\"/></body>",
		"www/blog/routes/+page.dreego":     "COMPONENT \"www/components\" IMPORT { Badge }\nCOMPONENT \"www/blog/components\" IMPORT { Local }\n<body><h1>blog home</h1><@Badge label=\"shared\"/><@Local label=\"own\"/></body>",
	}
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	gomod := "module t\n\ngo 1.27\n\nrequire github.com/dreego-stack/dreego/core v0.0.0\n\n" +
		"replace github.com/dreego-stack/dreego => " + repoRoot + "\n" +
		"replace github.com/dreego-stack/dreego/core => " + filepath.Join(repoRoot, "core") + "\n" +
		"replace github.com/dreego-stack/dreego/adapter/ssr => " + filepath.Join(repoRoot, "adapter", "ssr") + "\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644); err != nil {
		t.Fatal(err)
	}
	sum, err := os.ReadFile(filepath.Join(repoRoot, "go.sum"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.sum"), sum, 0644); err != nil {
		t.Fatal(err)
	}
	mainGo := `package main

import (
	"fmt"
	"os"

	ssr "github.com/dreego-stack/dreego/adapter/ssr"
	dreego "github.com/dreego-stack/dreego/core"
	"t/www/app"
	"t/www/blog"
)

func main() {
	wwwApp := dreego.New(app.App)
	blogApp := dreego.New(blog.App)
	go ssr.Listen(wwwApp, os.Getenv("WWW_ADDR"))
	if err := ssr.Listen(blogApp, os.Getenv("BLOG_ADDR")); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatal(err)
	}

	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	bin := filepath.Join(dir, "server")
	build := exec.Command("go", "build", "-mod=mod", "-o", bin, ".")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	wwwPort := dreegotest.FreePort(t)
	blogPort := dreegotest.FreePort(t)
	proc := exec.Command(bin)
	proc.Dir = dir
	proc.Env = append(os.Environ(),
		"WWW_ADDR="+fmt.Sprintf("127.0.0.1:%d", wwwPort),
		"BLOG_ADDR="+fmt.Sprintf("127.0.0.1:%d", blogPort),
	)
	if err := proc.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	t.Cleanup(func() {
		proc.Process.Kill()
		proc.Wait()
	})
	dreegotest.WaitForPort(t, wwwPort)
	dreegotest.WaitForPort(t, blogPort)

	for _, tc := range []struct {
		port  int
		want  string
		local bool
	}{
		{wwwPort, "www home", false},
		{blogPort, "blog home", true},
	} {
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", tc.port))
		if err != nil {
			t.Fatalf("GET :%d: %v", tc.port, err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("GET :%d = %d, want 200", tc.port, resp.StatusCode)
		}
		if !strings.Contains(string(body), tc.want) {
			t.Fatalf("GET :%d missing %q: %s", tc.port, tc.want, body)
		}
		if !strings.Contains(string(body), `class="badge"`) {
			t.Fatalf("GET :%d missing shared component output: %s", tc.port, body)
		}
		hasLocal := strings.Contains(string(body), `class="local"`)
		if hasLocal != tc.local {
			t.Fatalf("GET :%d app-local component present=%v, want %v: %s", tc.port, hasLocal, tc.local, body)
		}
	}
}
