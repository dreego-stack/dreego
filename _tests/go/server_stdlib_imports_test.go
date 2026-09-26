package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestServerStdlibImportSyncMutexCompiles(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": `GOIMPORT { sync }
<server>
var mu sync.Mutex
mu.Lock()
value := "locked"
mu.Unlock()
</server>
<body><p>{{ value }}</p></body>`,
	})
	routes := gen["www/routes/dree.go"]
	if !strings.Contains(routes, `"sync"`) {
		t.Fatalf("generated import block must contain \"sync\", got:\n%s", routes)
	}
	if !strings.Contains(routes, "sync.Mutex") {
		t.Fatalf("server code must be emitted verbatim, got:\n%s", routes)
	}
}

func TestServerStdlibImportMultiplePackages(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": `GOIMPORT { sync, encoding/json }
<server>
var mu sync.Mutex
mu.Lock()
payload, _ := json.Marshal(map[string]string{"state": "ok"})
mu.Unlock()
value := string(payload)
</server>
<body><p>{{ value }}</p></body>`,
	})
	routes := gen["www/routes/dree.go"]
	for _, want := range []string{`"sync"`, `"encoding/json"`} {
		if !strings.Contains(routes, want) {
			t.Fatalf("generated import block missing %s, got:\n%s", want, routes)
		}
	}
}

func TestComponentStdlibImportCompiles(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/components/Counter.dreego": `DREEFILE component ()
GOIMPORT { sync }
<server>
var mu sync.Mutex
mu.Lock()
label := "counter"
mu.Unlock()
</server>
<body><span>{{ label }}</span></body>`,
		"www/routes/+page.dreego": `<body><@Counter/></body>`,
	})
	components := gen["www/components/dree.go"]
	if !strings.Contains(components, `"sync"`) {
		t.Fatalf("generated component import block must contain \"sync\", got:\n%s", components)
	}
}

func TestGoImportRejectsModuleNotInGoMod(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": `GOIMPORT { statuna/auth }
<server>
value := "x"
</server>
<body><p>{{ value }}</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted a module that is not in go.mod:\n%s", out)
	}
	for _, want := range []string{"statuna/auth", "not in go.mod", "go get"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestGoImportStdlibWithoutAllowlist(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": `GOIMPORT { os }
<server>
value := os.Getenv("HOME")
if value == "" { value = "none" }
</server>
<body><p>{{ value }}</p></body>`,
	})
	routes := gen["www/routes/dree.go"]
	if !strings.Contains(routes, `"os"`) {
		t.Fatalf("generated import block must contain \"os\", got:\n%s", routes)
	}
}
