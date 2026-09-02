package transpiler

import (
	"strings"
	"testing"
)

func TestGenerateLayoutStyleBeforeCloseHTML(t *testing.T) {
	layoutSrc := "<body><html><body>{#slot}</body></html></body>\n\n<style>\n.x { color: red; }\n</style>\n"
	file := parseFile(t, layoutSrc)
	out, err := GenerateLayout(NewGenerator(), file, "Default")
	if err != nil {
		t.Fatalf("GenerateLayout: %v", err)
	}
	htmlIdx := strings.Index(out, "</html>")
	styleIdx := strings.Index(out, "<style>")
	if htmlIdx < 0 {
		t.Fatalf("no </html> in output:\n%s", out)
	}
	if styleIdx < 0 {
		t.Fatalf("no <style> in output:\n%s", out)
	}
	if styleIdx > htmlIdx {
		t.Fatalf("style emitted after </html> (style=%d html=%d):\n%s", styleIdx, htmlIdx, out)
	}
}
