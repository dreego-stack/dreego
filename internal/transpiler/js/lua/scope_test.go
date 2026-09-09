package lua

import (
	"strings"
	"testing"
)

func TestCompileLocalPrintShadowsBuiltin(t *testing.T) {
	artifact, err := Compile(`local print = function(value)
  document.title = value
end
print("local")`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(artifact.Code, "globalThis.dreegoLua.print") {
		t.Fatalf("local print was compiled as the runtime builtin:\n%s", artifact.Code)
	}
	if !strings.Contains(artifact.Code, `print("local");`) {
		t.Fatalf("local print call missing:\n%s", artifact.Code)
	}
	if artifact.Runtime != "" {
		t.Fatalf("runtime = %q, want no linked builtin", artifact.Runtime)
	}
}

func TestCompileFunctionParameterShadowsPrint(t *testing.T) {
	artifact, err := Compile(`local function call(print)
  print("parameter")
end`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(artifact.Code, "globalThis.dreegoLua.print") {
		t.Fatalf("print parameter was compiled as the runtime builtin:\n%s", artifact.Code)
	}
}

func TestCompileLocalFunctionPrintIsRecursive(t *testing.T) {
	artifact, err := Compile(`local function print(value)
  if value then
    print(false)
  end
end`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(artifact.Code, "globalThis.dreegoLua.print") || !strings.Contains(artifact.Code, "print(false);") {
		t.Fatalf("local function does not reference itself:\n%s", artifact.Code)
	}
}

func TestCompilePrintValueLinksBuiltin(t *testing.T) {
	artifact, err := Compile(`local logger = print
logger("value")`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(artifact.Code, "let logger = globalThis.dreegoLua.print;") || artifact.Runtime != "print" {
		t.Fatalf("JavaScript = %s, runtime = %q", artifact.Code, artifact.Runtime)
	}
}

func TestCompileBlockLocalPrintDoesNotEscape(t *testing.T) {
	artifact, err := Compile(`if true then
  local print = function(value)
    document.title = value
  end
  print("inside")
end
print("outside")`)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(artifact.Code, "globalThis.dreegoLua.print") != 1 {
		t.Fatalf("block scope was not preserved:\n%s", artifact.Code)
	}
}

func TestCompileLocalRequireShadowsForbiddenGlobal(t *testing.T) {
	artifact, err := Compile(`local require = function(value)
  document.title = value
end
require("local")`)
	if err != nil {
		t.Fatalf("local require was rejected: %v", err)
	}
	if !strings.Contains(artifact.Code, `require("local");`) {
		t.Fatalf("local require call missing:\n%s", artifact.Code)
	}
}
