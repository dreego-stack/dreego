package dreegotest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"

	dreefile "github.com/dreego-stack/dreego/internal/dreefile"
)

// Generate transpiles a .dreego source string to generated Go code using the
// transpiler pipeline directly (ParseFileHeaderStrict → Lex → Parse → GenerateMethodHandler).
// It replaces shell tests that run `dreego generate` and grep the output.
func Generate(t *testing.T, src string) string {
	t.Helper()
	out, err := generate(src)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	return out
}

// MustGenerate is like Generate but returns the generated code instead of
// failing the test. Use it when the caller wants to assert on the error.
func MustGenerate(t *testing.T, src string) (string, error) {
	t.Helper()
	return generate(src)
}

// MustCompile asserts that a .dreego source transpiles to non-empty Go code.
func MustCompile(t *testing.T, src string) {
	t.Helper()
	out, err := generate(src)
	if err != nil {
		t.Fatalf("MustCompile: %v", err)
	}
	if out == "" {
		t.Fatal("MustCompile: generated empty output")
	}
}

// MustFail asserts that a .dreego source produces a transpile error.
func MustFail(t *testing.T, src string) {
	t.Helper()
	if _, err := generate(src); err == nil {
		t.Fatal("MustFail: expected error, got none")
	}
}

// MustFailWith asserts that a .dreego source fails with an error containing
// the given substring.
func MustFailWith(t *testing.T, src, want string) {
	t.Helper()
	_, err := generate(src)
	if err == nil {
		t.Fatalf("MustFailWith: expected error containing %q, got none", want)
	}
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("MustFailWith: error %q does not contain %q", err.Error(), want)
	}
}

// MustContain asserts that the generated code contains the given substring.
func MustContain(t *testing.T, out, want string) {
	t.Helper()
	if !strings.Contains(out, want) {
		t.Fatalf("MustContain: generated code missing %q\n---\n%s", want, out)
	}
}

// MustNotContain asserts that the generated code does not contain the substring.
func MustNotContain(t *testing.T, out, want string) {
	t.Helper()
	if strings.Contains(out, want) {
		t.Fatalf("MustNotContain: generated code contains %q\n---\n%s", want, out)
	}
}

// GenerateComponent transpiles a DREEFILE component source to generated Go
// code. The component name comes from the file name in the real pipeline, so
// the caller must supply it explicitly; an empty or unusable name is an error.
func GenerateComponent(t *testing.T, name, src string) string {
	t.Helper()
	out, err := generateComponent(name, src)
	if err != nil {
		t.Fatalf("GenerateComponent: %v", err)
	}
	return out
}

// MustCompileComponent asserts that a DREEFILE component source transpiles to
// non-empty Go code under the given component name.
func MustCompileComponent(t *testing.T, name, src string) {
	t.Helper()
	out, err := generateComponent(name, src)
	if err != nil {
		t.Fatalf("MustCompileComponent: %v", err)
	}
	if out == "" {
		t.Fatal("MustCompileComponent: generated empty output")
	}
}

func generate(src string) (string, error) {
	header, body, err := dreefile.ParseFileHeaderStrict(src)
	if err != nil {
		return "", err
	}
	tokens, err := dreefile.Lex(body)
	if err != nil {
		return "", err
	}
	p := dreefile.NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return "", err
	}
	file.Imports = header.Imports
	file.Kind = header.Kind
	file.Layout = header.Layout
	file.GoImports = header.GoImports
	file.SourceContent = src
	if len(file.Server) == 0 {
		file.Server = []dreefile.ServerSection{{Method: "GET"}}
	}
	for i := range file.Server {
		if !file.Server[i].MethodExplicit {
			file.Server[i].Method = "GET"
		}
	}
	h := sha256.Sum256([]byte(src))
	scopeHash := hex.EncodeToString(h[:])[:12]
	gen := dreefile.NewGenerator()
	out, _, err := dreefile.GenerateMethodHandler(gen, file, nil, "routes", "index", "/{$}", scopeHash)
	return out, err
}

func generateComponent(name, src string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("component name is required: the DREEFILE component name comes from the file name")
	}
	if !dreefile.IsExportedGoIdentifier(name) {
		return "", fmt.Errorf("invalid component name %q: must be an exported Go identifier such as Card or ProductCard", name)
	}
	header, body, err := dreefile.ParseFileHeaderStrict(src)
	if err != nil {
		return "", err
	}
	if !header.IsComponent() {
		return "", fmt.Errorf("source is not a DREEFILE component")
	}
	tokens, err := dreefile.Lex(body)
	if err != nil {
		return "", err
	}
	file, err := dreefile.NewParserConcatServer(tokens).Parse()
	if err != nil {
		return "", err
	}
	comp := &dreefile.ComponentDef{Name: name, Props: header.Props, Slots: header.Slots}
	file.Component = comp
	file.Imports = header.Imports
	file.Kind = header.Kind
	file.Layout = header.Layout
	file.GoImports = header.GoImports
	file.SourceContent = src
	if len(file.Server) == 0 {
		file.Server = []dreefile.ServerSection{{Method: ""}}
	}
	h := sha256.Sum256([]byte(src))
	scopeHash := hex.EncodeToString(h[:])[:12]
	gen := dreefile.NewGenerator()
	return dreefile.GenerateComponent(gen, file, scopeHash)
}
