package dreegotest

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

// Generate transpiles a .dreego source string to generated Go code using the
// core pipeline directly (ParseHeader → Lex → Parse → GenerateMethodHandler).
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

func generate(src string) (string, error) {
	_, imports, body := dreego.ParseHeader(src)
	tokens, err := dreego.Lex(body)
	if err != nil {
		return "", err
	}
	p := dreego.NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return "", err
	}
	file.Imports = imports
	if len(file.Go) == 0 {
		file.Go = []dreego.GoSection{{Method: "GET"}}
	}
	for i := range file.Go {
		file.Go[i].Method = "GET"
	}
	h := sha256.Sum256([]byte(src))
	scopeHash := hex.EncodeToString(h[:])[:12]
	return dreego.GenerateMethodHandler(file, nil, "routes", "index", "/{$}", scopeHash)
}
