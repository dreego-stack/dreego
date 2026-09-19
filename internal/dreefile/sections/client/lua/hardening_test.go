package lua

import (
	"strings"
	"testing"
)

func TestCompileRejectsBreakAcrossFunctionBoundary(t *testing.T) {
	t.Parallel()

	_, err := Compile(`while true do
  local stop = function()
    break
  end
end`)
	if err == nil || !strings.Contains(err.Error(), "break is only valid inside a loop") {
		t.Fatalf("Compile() error = %v", err)
	}
}

func TestCompileAllowsBreakInFunctionLocalLoop(t *testing.T) {
	t.Parallel()

	artifact, err := Compile(`local stop = function()
  while true do
    break
  end
end`)
	if err != nil {
		t.Fatalf("Compile() error = %v", err)
	}
	if !strings.Contains(artifact.Code, "while (globalThis.dreegoLua.truthy(true))") ||
		!strings.Contains(artifact.Code, "break;") {
		t.Fatalf("Compile() output does not preserve the function-local loop:\n%s", artifact.Code)
	}
}

func FuzzCompileIsDeterministic(f *testing.F) {
	for _, source := range []string{
		`local value = 1 + 2`,
		`local values = {"first", key = true}`,
		`if ready then print("ready") end`,
		`local function broken( return end`,
		"local value = \"unfinished",
		"\x00\xff\n-- malformed input",
	} {
		f.Add(source)
	}

	f.Fuzz(func(t *testing.T, source string) {
		first, firstErr := Compile(source)
		second, secondErr := Compile(source)
		if first != second {
			t.Fatalf("Compile() is not deterministic:\nfirst: %#v\nsecond: %#v", first, second)
		}
		if errorText(firstErr) != errorText(secondErr) {
			t.Fatalf("Compile() error is not deterministic:\nfirst: %v\nsecond: %v", firstErr, secondErr)
		}
	})
}

func errorText(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
