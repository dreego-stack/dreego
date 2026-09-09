package process

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/transpiler/codegen"
)

func TestLuaOutputCarriesOriginalSourceReference(t *testing.T) {
	artifact, err := compile(codegen.NewState(), "lua", `local value = 1 + "invalid"`, nil, nil, "www/routes/account.dreego", 12)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(artifact.Code, strings.Repeat("\n", 11)) {
		t.Fatalf("Lua output is not aligned to source line 12:\n%s", artifact.Code)
	}
	for _, want := range []string{"sourceURL=dreego:///www/routes/account.dreego", "source line 12"} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("Lua output missing %q:\n%s", want, artifact.Code)
		}
	}
}
