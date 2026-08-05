package dreegotest

import (
	"fmt"
	"html"
	"net/http/httptest"
	"testing"

	"codeberg.org/dreego/dreego/core"
)

func RenderComponent(t *testing.T, fn core.ComponentFunc, props ...any) string {
	t.Helper()
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	ctx := core.NewSSR(rec, req)
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
