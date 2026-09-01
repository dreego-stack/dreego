// Package dreegotest provides helpers for testing generated Dreego apps.
// Rendering goes through the target-neutral render contract; no HTTP types are
// involved.
package dreegotest

import (
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func RenderComponent(t *testing.T, component dreego.Component) string {
	t.Helper()
	result, err := dreego.Render(component)
	if err != nil {
		t.Fatalf("RenderComponent: component returned error: %v", err)
	}
	return string(result.HTML)
}
