package lua

import (
	"strings"
	"testing"
)

func TestCompileTableLiteralIndexAndLength(t *testing.T) {
	artifact, err := Compile(`local values = {"first", "second", label = "ready", [4] = "fourth"}
local first = values[1]
local label = values["label"]
local size = #values`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`globalThis.dreegoLua.table([[1, "first"], [2, "second"], ["label", "ready"], [4, "fourth"]])`,
		"globalThis.dreegoLua.get(values, 1)",
		`globalThis.dreegoLua.get(values, "label")`,
		"globalThis.dreegoLua.length(values)",
	} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
}

func TestCompileTableAssignmentUsesNilDeletion(t *testing.T) {
	artifact, err := Compile(`local values = {"first"; "second"}
values[1] = "updated"
values[2] = nil
values.label = nil`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`globalThis.dreegoLua.set(values, 1, "updated");`,
		"globalThis.dreegoLua.set(values, 2, null);",
		"values.label = null;",
	} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
}

func TestTableRuntimeSupportsDottedFieldsAndNilDeletion(t *testing.T) {
	runtime := Bundle("table")
	for _, want := range []string{"new Proxy", "entries.get", "entries.delete"} {
		if !strings.Contains(runtime, want) {
			t.Fatalf("runtime missing %q: %s", want, runtime)
		}
	}
}

func TestCompileTableRejectsNilAndNaNKeysAtRuntime(t *testing.T) {
	runtime := Bundle("table")
	for _, want := range []string{"Lua table index is nil", "Lua table index is NaN"} {
		if !strings.Contains(runtime, want) {
			t.Fatalf("runtime missing %q: %s", want, runtime)
		}
	}
}

func TestCompileTableLengthUsesContiguousSequence(t *testing.T) {
	runtime := Bundle("length")
	if !strings.Contains(runtime, "while (entries.has(length + 1))") {
		t.Fatalf("runtime does not define contiguous one-based length: %s", runtime)
	}
}
