package tests

import (
	"regexp"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestTypeScriptRootClientEmitsJavaScript(t *testing.T) {
	out := dreegotest.Generate(t, `<body><main>Ready</main></body>
<client lang="ts">
const count: number = 2
document.body.dataset.count = String(count)
</client>`)
	dreegotest.MustContain(t, out, "const count = 2")
	dreegotest.MustNotContain(t, out, ": number")
}

func TestTypeScriptNestedBodyScriptEmitsInPlace(t *testing.T) {
	out := dreegotest.Generate(t, `<body><main>Ready</main>
<script lang="ts">const ready: boolean = true</script>
</body>`)
	dreegotest.MustContain(t, out, "<script>")
	dreegotest.MustContain(t, out, "const ready = true")
	dreegotest.MustNotContain(t, out, ": boolean")
}

func TestTypeScriptChecksGeneratedGoModels(t *testing.T) {
	source := `<server>
type User struct {
    Name string ` + "`json:\"name\"`" + `
    Age int ` + "`json:\"age\"`" + `
}
</server>
<body></body>
<client lang="ts">const user: User = { name: "Ada", age: "wrong" }</client>`
	_, err := dreegotest.MustGenerate(t, source)
	if err == nil || !strings.Contains(err.Error(), "TS2322") {
		t.Fatalf("error = %v", err)
	}
}

func TestInlineTypeScriptChecksGeneratedGoModels(t *testing.T) {
	source := `<server>type User struct { Name string ` + "`json:\"name\"`" + ` }</server>
<body><script lang="ts">const user: User = { name: 42 }</script></body>`
	_, err := dreegotest.MustGenerate(t, source)
	if err == nil || !strings.Contains(err.Error(), "TS2322") {
		t.Fatalf("error = %v", err)
	}
}

func TestTypeScriptRejectsInvalidCode(t *testing.T) {
	dreegotest.MustFailWith(t, `<body></body>
<client lang="ts">const count: number = "wrong"</client>`, "TS2322")
}

func TestTypeScriptDiagnosticUsesDreegoSourceLine(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body><main>Ready</main></body>
<client lang="ts">
const count: number = "wrong"
</client>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate succeeded:\n%s", out)
	}
	if !regexp.MustCompile(`www/routes/get\.dreego\(3,\d+\).*TS2322`).MatchString(out) {
		t.Fatalf("diagnostic does not identify the TypeScript source line:\n%s", out)
	}
}

func TestPlainJavaScriptNeedsNoTypeScriptSyntax(t *testing.T) {
	out := dreegotest.Generate(t, `<body><script type="application/ld+json">{"ready":true}</script></body>
<client>window.ready = true</client>`)
	dreegotest.MustContain(t, out, `type="application/ld+json"`)
	dreegotest.MustContain(t, out, "window.ready = true")
}
