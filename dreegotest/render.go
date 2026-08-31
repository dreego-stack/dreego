// Package dreegotest provides helpers for testing generated Dreego apps.
// Rendering goes through the target-neutral render contract; no HTTP types are
// involved.
package dreegotest

import (
	"fmt"
	"html"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func RenderComponent(t *testing.T, fn dreego.ComponentFunc, props ...any) string {
	t.Helper()
	ctx := dreego.NewContext()
	for i := 0; i+1 < len(props); i += 2 {
		key, ok := props[i].(string)
		if !ok {
			t.Fatalf("RenderComponent: prop key at index %d is not a string", i)
		}
		ctx.Set(key, html.EscapeString(fmt.Sprintf("%v", props[i+1])))
	}
	out, err := fn(ctx)
	if err != nil {
		t.Fatalf("RenderComponent: component returned error: %v", err)
	}
	return out
}
