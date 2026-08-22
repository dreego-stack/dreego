package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestCLIGenerateCheckTracksMethodRouteChanges(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": `<go>message := "get"</go><div>{{ message }}</div><go method="post">message := "post"</go><div method="post">{{ message }}</div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate", "--check"); err != nil || !strings.Contains(out, "up-to-date") {
		t.Fatalf("initial check: %v\n%s", err, out)
	}
	path := filepath.Join(dir, "www/routes/+page.dreego")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.Replace(string(data), `message := "post"`, `message := "changed"`, 1)), 0644); err != nil {
		t.Fatal(err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate", "--check"); err == nil || !strings.Contains(strings.ToLower(out), "out of date") {
		t.Fatalf("expected stale method route, err=%v output=%s", err, out)
	}
}
