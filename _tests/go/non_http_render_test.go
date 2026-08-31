package tests

import (
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/dreegotest"
)

func TestGeneratedComponentRendersWithoutHTTP(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		"www/components/Badge.dreego": `Component Badge (label string)
<body><span class="badge">{{ label }}</span></body>`,
		"www/routes/get.dreego": `<head><title>Shop</title></head>

<body><@Badge label={"<b>hot</b>"}/></body>`,
	}

	served := dreegotest.Serve(t, files)
	code, body := served.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	wantFragment := `<span class="badge">&lt;b&gt;hot&lt;/b&gt;</span>`
	if !strings.Contains(body, wantFragment) {
		t.Fatalf("served body missing escaped badge fragment %q, got: %s", wantFragment, body)
	}

	generated := dreegotest.Build(t, files)
	components := generated["www/components/dree.go"]
	if !strings.Contains(components, "dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {") {
		t.Fatalf("generated component must wrap the render body in a ComponentFunc, got:\n%s", components)
	}

	repro := dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {
		return `<span class="badge">` + ctx.Get("label") + `</span>`, nil
	})
	renderNeutral := dreegotest.RenderComponent(t, repro, "label", "<b>hot</b>")
	if renderNeutral != wantFragment {
		t.Errorf("render-neutral output %q != served fragment %q", renderNeutral, wantFragment)
	}
}
