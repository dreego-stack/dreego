package transpiler

import (
	"strings"
	"testing"
)

// A pure GET page rendered inside a layout that carries a <server> section must
// still generate coherent code: the page renderer takes dreego.RenderContext and
// passes it to the layout renderer, which also takes dreego.RenderContext.
func TestGenerateMethodHandlerPureGETWithServerLayout(t *testing.T) {
	routeSrc := "<head><title>Home</title></head>\n<body><h1>Home</h1></body>\n"
	layoutSrc := "<server>\nmsg := \"layout\"\n_ = msg\n</server>\n<body><html><body>{#slot}</body></html></body>\n"

	file := parseFile(t, routeSrc)
	layout := &layoutEntry{file: parseFile(t, layoutSrc), name: "Default"}

	got, _, err := GenerateMethodHandler(NewGenerator(), file, layout, "home", "index", "/", scopeHashFor(routeSrc))
	if err != nil {
		t.Fatalf("GenerateMethodHandler: %v", err)
	}

	for _, want := range []string{
		"func renderIndex(c dreego.RenderContext) (string, error) {",
		"layouts.Default(c, pageContent, head)",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("generated page missing %q, got:\n%s", want, got)
		}
	}

	layoutOut, err := GenerateLayout(NewGenerator(), layout.file, "Default")
	if err != nil {
		t.Fatalf("GenerateLayout: %v", err)
	}
	if !strings.Contains(layoutOut, "func Default(c dreego.RenderContext, content, head string) (string, error) {") {
		t.Errorf("layout renderer must take dreego.RenderContext, got:\n%s", layoutOut)
	}
}
