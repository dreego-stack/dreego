package ts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProcessUsesNativeCompilerOutput(t *testing.T) {
	t.Setenv("DREEGO_TYPESCRIPT_COMPILER", os.Args[0])
	t.Setenv("DREEGO_TYPESCRIPT_COMPILER_ARGS", "-test.run=TestTypeScriptCompilerHelper --")
	t.Setenv("DREEGO_TYPESCRIPT_TEST_HELPER", "emit")

	artifact, err := Process(`const count: number = 2;`, "", "component.dreego", 3)
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Code != `const count = 2;` {
		t.Fatalf("JavaScript = %q", artifact.Code)
	}
}

func TestProcessReturnsCompilerDiagnostic(t *testing.T) {
	t.Setenv("DREEGO_TYPESCRIPT_COMPILER", os.Args[0])
	t.Setenv("DREEGO_TYPESCRIPT_COMPILER_ARGS", "-test.run=TestTypeScriptCompilerHelper --")
	t.Setenv("DREEGO_TYPESCRIPT_TEST_HELPER", "fail")

	_, err := Process(`const count: number = "wrong";`, "", "component.dreego", 8)
	if err == nil || !strings.Contains(err.Error(), "component.dreego(8,7)") || !strings.Contains(err.Error(), "TS2322") {
		t.Fatalf("error = %v", err)
	}
}

func TestProcessWithInstalledNativeCompiler(t *testing.T) {
	if _, err := compilerPath(); err != nil {
		t.Skip(err)
	}
	t.Setenv("DREEGO_TYPESCRIPT_COMPILER_ARGS", "")

	artifact, err := Process(`const count: number = 2;`, "", "route.dreego", 1)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(artifact.Code, ": number") || !strings.Contains(artifact.Code, "const count = 2") {
		t.Fatalf("JavaScript = %q", artifact.Code)
	}
}

func TestTypeScriptCompilerHelper(t *testing.T) {
	mode := os.Getenv("DREEGO_TYPESCRIPT_TEST_HELPER")
	if mode == "" {
		return
	}
	if mode == "fail" {
		_, _ = os.Stderr.WriteString("input.ts(1,7): error TS2322: Type mismatch\n")
		os.Exit(1)
	}
	args := os.Args
	var input, outDir string
	for i, arg := range args {
		if arg == "--outDir" && i+1 < len(args) {
			outDir = args[i+1]
		}
		if strings.HasSuffix(arg, "input.ts") {
			input = arg
		}
	}
	source, err := os.ReadFile(input)
	if err != nil {
		os.Exit(2)
	}
	javascript := strings.ReplaceAll(string(source), ": number", "")
	if err := os.WriteFile(filepath.Join(outDir, "input.js"), []byte(javascript), 0600); err != nil {
		os.Exit(2)
	}
	os.Exit(0)
}
