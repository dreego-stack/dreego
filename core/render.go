package core

import (
	"fmt"
	"html"

	"github.com/dreego-stack/dreego/core/internal/context"
)

// NewContext builds a render-neutral SSRContext with no HTTP request or
// response. Data/Set/Get work on the in-memory data map; request-bound methods
// panic on the nil request and must not be called outside an SSR host.
func NewContext() *SSRContext {
	return context.NewSSR(nil, nil)
}

func Render(fn ComponentFunc, props ...any) (string, error) {
	ctx := NewContext()
	for i := 0; i+1 < len(props); i += 2 {
		key, ok := props[i].(string)
		if !ok {
			return "", fmt.Errorf("render: prop key at index %d is not a string", i)
		}
		ctx.Set(key, html.EscapeString(fmt.Sprintf("%v", props[i+1])))
	}
	return fn(ctx)
}
