package dreegotest

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	transpiler "github.com/dreego-stack/dreego/internal/transpiler"
)

// Generate transpiles a .dreego source string to generated Go code using the
// transpiler pipeline directly (ParseHeader → Lex → Parse → GenerateMethodHandler).
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

// MustCompile asserts that a .dreego source transpiles without error.
func MustCompile(t *testing.T, src string) {
	t.Helper()
	if _, err := generate(src); err != nil {
		t.Fatalf("MustCompile: %v", err)
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

// GenerateComponent transpiles a .dreego component source to generated Go code
// using the transpiler pipeline directly (ParseHeader → Lex → Parse → GenerateComponent).
func GenerateComponent(t *testing.T, src string) string {
	t.Helper()
	out, err := generateComponent(src)
	if err != nil {
		t.Fatalf("GenerateComponent: %v", err)
	}
	return out
}

// MustCompileComponent asserts that a .dreego component source transpiles.
func MustCompileComponent(t *testing.T, src string) {
	t.Helper()
	if _, err := generateComponent(src); err != nil {
		t.Fatalf("MustCompileComponent: %v", err)
	}
}

func generate(src string) (string, error) {
	_, imports, body := transpiler.ParseHeader(src)
	tokens, err := transpiler.Lex(body)
	if err != nil {
		return "", err
	}
	p := transpiler.NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return "", err
	}
	file.Imports = imports
	file.SourceContent = src
	if len(file.Go) == 0 {
		file.Go = []transpiler.GoSection{{Method: "GET"}}
	}
	for i := range file.Go {
		if !file.Go[i].MethodExplicit {
			file.Go[i].Method = "GET"
		}
	}
	h := sha256.Sum256([]byte(src))
	scopeHash := hex.EncodeToString(h[:])[:12]
	gen := transpiler.NewGenerator()
	out, _, err := transpiler.GenerateMethodHandler(gen, file, nil, "routes", "index", "/{$}", scopeHash)
	return out, err
}

func generateComponent(src string) (string, error) {
	comp, _, body := transpiler.ParseHeader(src)
	if comp == nil || comp.Name == "" {
		return "", nil
	}
	tokens, err := transpiler.Lex(body)
	if err != nil {
		return "", err
	}
	p := transpiler.NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return "", err
	}
	file.Component = comp
	file.SourceContent = src
	if len(file.Go) == 0 {
		file.Go = []transpiler.GoSection{{Method: ""}}
	}
	h := sha256.Sum256([]byte(src))
	scopeHash := hex.EncodeToString(h[:])[:12]
	gen := transpiler.NewGenerator()
	return transpiler.GenerateComponent(gen, file, scopeHash)
}
