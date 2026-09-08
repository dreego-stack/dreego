package lua

import (
	"strings"
	"testing"
)

func TestCompileDirectStatements(t *testing.T) {
	artifact, err := Compile(`
local count = 2
count = count + 1
document.title = "Count " .. count
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"let count = 2;",
		"count = globalThis.dreegoLua.add(count, 1);",
		`document.title = globalThis.dreegoLua.concat("Count ", count);`,
	} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
	if artifact.Runtime != "add,concat" {
		t.Fatalf("runtime = %q", artifact.Runtime)
	}
}

func TestCompilePreservesLuaTruthiness(t *testing.T) {
	artifact, err := Compile(`
local value = 0
if value then
  print("zero is true")
elseif not false then
  print("fallback")
else
  print("never")
end
local selected = value or "fallback"
`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"globalThis.dreegoLua.truthy(value)",
		"!globalThis.dreegoLua.truthy(false)",
		"globalThis.dreegoLua.or(value, () => \"fallback\")",
		"globalThis.dreegoLua.print(\"zero is true\")",
	} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
	if artifact.Runtime != "or,print,truthy" {
		t.Fatalf("runtime = %q", artifact.Runtime)
	}
}

func TestCompileRejectsUnsupportedTable(t *testing.T) {
	_, err := Compile(`local values = {1, 2, 3}`)
	if err == nil || !strings.Contains(err.Error(), "tables are not supported") {
		t.Fatalf("error = %v", err)
	}
}

func TestRuntimeIncludesOnlyRequiredHelpers(t *testing.T) {
	artifact, err := Compile(`print("hello")`)
	if err != nil {
		t.Fatal(err)
	}
	runtime := Bundle(artifact.Runtime)
	if !strings.Contains(runtime, "console.log") || strings.Contains(runtime, "toString") {
		t.Fatalf("runtime bundle = %s", runtime)
	}
}

func TestDirectLuaNeedsNoRuntime(t *testing.T) {
	artifact, err := Compile(`document.title = "Dreego"`)
	if err != nil {
		t.Fatal(err)
	}
	if artifact.Runtime != "" || Bundle(artifact.Runtime) != "" {
		t.Fatalf("unexpected runtime %q", artifact.Runtime)
	}
}

func TestCompileManglesJavaScriptReservedNames(t *testing.T) {
	artifact, err := Compile(`local class = 1
local _lua_class = 2
class = class + _lua_class`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(artifact.Code, "let $lua_class = 1;") ||
		!strings.Contains(artifact.Code, "let _lua_class = 2;") ||
		!strings.Contains(artifact.Code, "$lua_class = globalThis.dreegoLua.add($lua_class, _lua_class);") {
		t.Fatalf("JavaScript = %s", artifact.Code)
	}
}

func TestCompileRuntimeNamespaceCannotBeShadowed(t *testing.T) {
	artifact, err := Compile(`local dreegoLua = "application value"
local total = 1 + 2`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(artifact.Code, "globalThis.dreegoLua.add(1, 2)") {
		t.Fatalf("JavaScript runtime namespace can be shadowed:\n%s", artifact.Code)
	}
}

func TestCompileNumbersAndRejectsBareExpressions(t *testing.T) {
	artifact, err := Compile(`local small = 1.5e-2`)
	if err != nil || !strings.Contains(artifact.Code, "let small = 1.5e-2;") {
		t.Fatalf("artifact = %#v, error = %v", artifact, err)
	}
	if _, err := Compile(`1 + 2`); err == nil || !strings.Contains(err.Error(), "only function calls") {
		t.Fatalf("bare expression error = %v", err)
	}
	if _, err := Compile(`local broken = 1e`); err == nil || !strings.Contains(err.Error(), "invalid number") {
		t.Fatalf("invalid number error = %v", err)
	}
}

func TestCompileLocalFunctionsAndBrowserCallbacks(t *testing.T) {
	artifact, err := Compile(`local button = document:querySelector("#count")
local function update(label)
  button.textContent = label
end
button:addEventListener("click", function(event)
  update(event.type)
end)`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`document.querySelector("#count")`,
		"let update = (label) => {",
		"button.textContent = label;",
		`button.addEventListener("click", (event) => {`,
		"update(event.type);",
	} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
}

func TestCompileAllowsMultilineCalls(t *testing.T) {
	artifact, err := Compile("print(\n  \"hello\",\n  \"world\"\n)")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(artifact.Code, `globalThis.dreegoLua.print("hello", "world")`) {
		t.Fatalf("JavaScript = %s", artifact.Code)
	}
}

func TestCompileRejectsServerOnlyLibraries(t *testing.T) {
	for _, source := range []string{`require("module")`, `load("code")`, `dofile("file")`} {
		if _, err := Compile(source); err == nil || !strings.Contains(err.Error(), "not available") {
			t.Fatalf("Compile(%q) error = %v", source, err)
		}
	}
}

func TestCompileCommonLuaSyntax(t *testing.T) {
	artifact, err := Compile(`-- browser behavior
local greeting = 'hello'
local equal = 2 == 2
local different = greeting ~= "world";
local remainder = -3 % 2`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`let greeting = "hello";`, "let equal = (2 === 2);", `let different = (greeting !== "world");`, "globalThis.dreegoLua.mod(globalThis.dreegoLua.neg(3), 2)"} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
}

func TestCompileIsolatesEveryLuaBlock(t *testing.T) {
	artifact, err := Compile(`local status = "ready"`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(artifact.Code, "(() => {\n") || !strings.HasSuffix(artifact.Code, "\n})();") {
		t.Fatalf("JavaScript is not isolated:\n%s", artifact.Code)
	}
}

func TestRuntimeRejectsModuloByZero(t *testing.T) {
	runtime := Bundle("mod")
	if !strings.Contains(runtime, "right === 0") || !strings.Contains(runtime, "Lua modulo by zero") {
		t.Fatalf("modulo helper does not reject zero: %s", runtime)
	}
}

func TestRuntimeTreatsJavaScriptUndefinedAsNil(t *testing.T) {
	runtime := Bundle("truthy")
	if !strings.Contains(runtime, "value !== undefined") {
		t.Fatalf("truthiness helper does not normalize undefined: %s", runtime)
	}
}
