package transpiler

import (
	"strings"
	"testing"
)

func TestAllowedStdlibImportList(t *testing.T) {
	for _, path := range []string{"strings", "net/http", "fmt", "sync", "encoding/json", "time", "errors", "strconv"} {
		if !allowedStdlibImport(path) {
			t.Errorf("expected %q to be allow-listed", path)
		}
	}
	for _, path := range []string{"os", "os/exec", "unsafe", "reflect", "syscall", "net", "plugin", "github.com/evil/pkg"} {
		if allowedStdlibImport(path) {
			t.Errorf("expected %q to be rejected", path)
		}
	}
}

func TestRegisterGoImportsRejectsUnknown(t *testing.T) {
	gen := NewGenerator()
	err := registerGoImports(gen, "routes", "www/routes/+page.dreego", []string{"os"})
	if err == nil {
		t.Fatal("expected an error for a non-allow-listed import")
	}
	if !strings.Contains(err.Error(), "os") || !strings.Contains(err.Error(), "www/routes/+page.dreego") {
		t.Fatalf("diagnostic must name the path and the file, got: %v", err)
	}
	if !strings.Contains(err.Error(), "strings") {
		t.Fatalf("diagnostic must list supported packages, got: %v", err)
	}
}

func TestStdImportsForIncludesDeclared(t *testing.T) {
	gen := NewGenerator()
	if err := registerGoImports(gen, "routes", "www/routes/+page.dreego", []string{"sync", "encoding/json"}); err != nil {
		t.Fatalf("registerGoImports: %v", err)
	}
	out := stdImportsFor(gen, "routes", "")
	for _, want := range []string{`"sync"`, `"encoding/json"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in import block, got:\n%s", want, out)
		}
	}
}

func TestStdImportsForKeepsAutoDetected(t *testing.T) {
	gen := NewGenerator()
	out := stdImportsFor(gen, "routes", `v := strings.ToUpper("x"); _ = fmt.Sprintf("%v", v)`)
	for _, want := range []string{`"strings"`, `"fmt"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in import block, got:\n%s", want, out)
		}
	}
}

func TestStdImportsForDeduplicatesDeclaredAndDetected(t *testing.T) {
	gen := NewGenerator()
	if err := registerGoImports(gen, "routes", "www/routes/+page.dreego", []string{"strings"}); err != nil {
		t.Fatalf("registerGoImports: %v", err)
	}
	out := stdImportsFor(gen, "routes", `v := strings.ToUpper("x")`)
	if strings.Count(out, `"strings"`) != 1 {
		t.Fatalf("expected exactly one strings import, got:\n%s", out)
	}
}
