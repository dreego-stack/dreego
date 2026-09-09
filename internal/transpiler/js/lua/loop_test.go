package lua

import (
	"strings"
	"testing"
)

func TestCompileWhileLoop(t *testing.T) {
	artifact, err := Compile(`local count = 0
while count < 3 do
  count = count + 1
end`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"while (globalThis.dreegoLua.truthy(",
		"globalThis.dreegoLua.lt(count, 3)",
		"count = globalThis.dreegoLua.add(count, 1);",
	} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
}

func TestCompileNumericForLoop(t *testing.T) {
	artifact, err := Compile(`local total = 0
for index = 3, 1, -1 do
  total = total + index
end`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"for (let index of globalThis.dreegoLua.numericFor(3, 1, globalThis.dreegoLua.neg(1)))",
		"total = globalThis.dreegoLua.add(total, index);",
	} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
}

func TestCompileNumericForDefaultsToOne(t *testing.T) {
	artifact, err := Compile(`for index = 1, 3 do print(index) end`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(artifact.Code, "globalThis.dreegoLua.numericFor(1, 3, 1)") {
		t.Fatalf("JavaScript = %s", artifact.Code)
	}
}

func TestCompileBreakRequiresLoop(t *testing.T) {
	if _, err := Compile(`break`); err == nil || !strings.Contains(err.Error(), "break is only valid inside a loop") {
		t.Fatalf("error = %v", err)
	}
	if _, err := Compile(`while true do break end`); err != nil {
		t.Fatal(err)
	}
}

func TestNumericForRuntimeRejectsInvalidBoundsAndZeroStep(t *testing.T) {
	runtime := Bundle("numericFor")
	for _, want := range []string{"Lua numeric for requires finite numbers", "Lua numeric for step cannot be zero"} {
		if !strings.Contains(runtime, want) {
			t.Fatalf("runtime missing %q: %s", want, runtime)
		}
	}
}
