package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestAccessibilityRuntimeAttrs(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeFixture(t, "components")
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, `id="main"`) {
		t.Errorf("rendered HTML missing skip-link target id=\"main\": %s", body)
	}
	if !strings.Contains(body, "skip to content") {
		t.Errorf("rendered HTML missing skip-link text: %s", body)
	}
	if !strings.Contains(body, `aria-label="Primary"`) {
		t.Errorf("rendered HTML missing nav aria-label=\"Primary\": %s", body)
	}
}
