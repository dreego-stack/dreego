package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestNestedComponentSlotCodegenUsesValidBuilders(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Inner.dreego": "Component Inner ()\n<div><aside>{#slot header}{/slot}</aside></div>",
		"www/components/Outer.dreego": "Component Outer ()\n<div><header>{#slot header}{/slot}</header><main>{#slot}</main></div>",
		"www/routes/get.dreego": `<div><@Outer>{#slot header}<strong>outer</strong>{/slot}` +
			`<@Inner>{#slot header}<em>inner</em>{/slot}</@Inner></@Outer></div>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	for _, want := range []string{"<strong>outer</strong>", "<em>inner</em>"} {
		if !strings.Contains(body, want) {
			t.Fatalf("expected body to contain %q, got: %s", want, body)
		}
	}
}
