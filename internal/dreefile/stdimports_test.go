package dreefile

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func TestRegisterGoImportsAcceptsStdlibAndModule(t *testing.T) {
	gen := NewGenerator()
	gen.Module = "example.com/app"
	gen.Requires = map[string]string{"statuna/auth": "v1.2.3", "github.com/dreego-stack/dreego-ui": "v0.4.0"}
	imports := []ir.GoImport{
		{Path: "os"},
		{Path: "strings"},
		{Alias: "myauth", Path: "statuna/auth"},
		{Path: "example.com/app/internal/foo"},
		{Path: "github.com/dreego-stack/dreego-ui/components"},
	}
	if err := registerGoImports(gen, "routes", "www/routes/+page.dreego", imports); err != nil {
		t.Fatalf("registerGoImports: %v", err)
	}
	if len(gen.GoImports["routes"]) != len(imports) {
		t.Fatalf("expected %d imports, got %+v", len(imports), gen.GoImports["routes"])
	}
}

func TestRegisterGoImportsRejectsMissingModule(t *testing.T) {
	gen := NewGenerator()
	gen.Module = "example.com/app"
	gen.Requires = map[string]string{"statuna/auth": "v1.2.3"}
	err := registerGoImports(gen, "routes", "www/routes/+page.dreego", []ir.GoImport{{Path: "statuna/missing"}})
	if err == nil {
		t.Fatal("expected an error for an unresolvable import")
	}
	for _, want := range []string{"statuna/missing", "not in go.mod", "go get"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("diagnostic must contain %q, got: %v", want, err)
		}
	}
}

func TestStdImportsForEmitsAlias(t *testing.T) {
	gen := NewGenerator()
	gen.Module = "example.com/app"
	gen.Requires = map[string]string{"statuna/auth": "v1.2.3"}
	if err := registerGoImports(gen, "routes", "www/routes/+page.dreego", []ir.GoImport{{Alias: "myauth", Path: "statuna/auth"}}); err != nil {
		t.Fatalf("registerGoImports: %v", err)
	}
	out, err := stdImportsFor(gen, "routes", "")
	if err != nil {
		t.Fatalf("stdImportsFor: %v", err)
	}
	if !strings.Contains(out, `myauth "statuna/auth"`) {
		t.Fatalf("expected aliased import, got:\n%s", out)
	}
}

func TestStdImportsForBarePathUsesBaseName(t *testing.T) {
	gen := NewGenerator()
	gen.Module = "example.com/app"
	gen.Requires = map[string]string{"statuna/auth": "v1.2.3"}
	if err := registerGoImports(gen, "routes", "www/routes/+page.dreego", []ir.GoImport{{Path: "statuna/auth"}}); err != nil {
		t.Fatalf("registerGoImports: %v", err)
	}
	out, err := stdImportsFor(gen, "routes", "")
	if err != nil {
		t.Fatalf("stdImportsFor: %v", err)
	}
	if !strings.Contains(out, `"statuna/auth"`) || strings.Contains(out, "auth ") {
		t.Fatalf("expected bare import path, got:\n%s", out)
	}
}

func TestStdImportsForRejectsBaseNameCollision(t *testing.T) {
	gen := NewGenerator()
	gen.Module = "example.com/app"
	gen.Requires = map[string]string{"a/auth": "v1", "b/auth": "v1"}
	if err := registerGoImports(gen, "routes", "www/routes/+page.dreego", []ir.GoImport{{Path: "a/auth"}, {Path: "b/auth"}}); err != nil {
		t.Fatalf("registerGoImports: %v", err)
	}
	_, err := stdImportsFor(gen, "routes", "")
	if err == nil {
		t.Fatal("expected a base-name collision error")
	}
	for _, want := range []string{"a/auth", "b/auth", "auth"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("collision diagnostic must contain %q, got: %v", want, err)
		}
	}
}

func TestStdImportsForKeepsAutoDetected(t *testing.T) {
	gen := NewGenerator()
	out, err := stdImportsFor(gen, "routes", `v := strings.ToUpper("x"); _ = fmt.Sprintf("%v", v)`)
	if err != nil {
		t.Fatalf("stdImportsFor: %v", err)
	}
	for _, want := range []string{`"strings"`, `"fmt"`} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %q in import block, got:\n%s", want, out)
		}
	}
}

func TestStdImportsForDetectsOnlyWholeIdentifiers(t *testing.T) {
	gen := NewGenerator()
	out, err := stdImportsFor(gen, "routes", `v := mystrings.ToUpper("x"); _ = myfmt.Sprintf("%v", v)`)
	if err != nil {
		t.Fatalf("stdImportsFor: %v", err)
	}
	if strings.Contains(out, `"strings"`) || strings.Contains(out, `"fmt"`) {
		t.Fatalf("aliased identifiers must not trigger auto-detection, got:\n%s", out)
	}
}

func TestStdImportsForDeduplicatesDeclaredAndDetected(t *testing.T) {
	gen := NewGenerator()
	if err := registerGoImports(gen, "routes", "www/routes/+page.dreego", []ir.GoImport{{Path: "strings"}}); err != nil {
		t.Fatalf("registerGoImports: %v", err)
	}
	out, err := stdImportsFor(gen, "routes", `v := strings.ToUpper("x")`)
	if err != nil {
		t.Fatalf("stdImportsFor: %v", err)
	}
	if strings.Count(out, `"strings"`) != 1 {
		t.Fatalf("expected exactly one strings import, got:\n%s", out)
	}
}

func TestIsStdlibImport(t *testing.T) {
	for _, path := range []string{"strings", "net/http", "fmt", "os", "encoding/json"} {
		if !isStdlibImport(path) {
			t.Errorf("expected %q to be stdlib", path)
		}
	}
	for _, path := range []string{"github.com/evil/pkg", "statuna/auth", "unsafe/x"} {
		if isStdlibImport(path) {
			t.Errorf("expected %q not to be stdlib", path)
		}
	}
}
