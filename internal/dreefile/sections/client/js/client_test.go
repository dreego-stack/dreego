package js

import "testing"

func TestProcessPreservesJavaScriptInput(t *testing.T) {
	t.Parallel()

	code := `const message = "ready";`
	if got := Process(code); got.Code != code {
		t.Fatalf("Process().Code = %q, want %q", got.Code, code)
	}
}
