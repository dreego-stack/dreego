package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestLuaLocalPrintShadowsRuntimeBuiltin(t *testing.T) {
	out := dreegotest.Generate(t, `<body></body>
<client lang="lua">
local print = function(value)
    document.title = value
end
print("local")
</client>`)
	if strings.Contains(out, "globalThis.dreegoLua.print") {
		t.Fatalf("generated client calls the runtime print helper:\n%s", out)
	}
	dreegotest.MustContain(t, out, `print("local");`)
}
