package core

import (
	"fmt"
	"html"

	"github.com/dreego-stack/dreego/core/internal/context"
)

func Render(fn ComponentFunc, props ...any) (string, error) {
	ctx := context.NewSSR(nil, nil)
	for i := 0; i+1 < len(props); i += 2 {
		key, ok := props[i].(string)
		if !ok {
			return "", fmt.Errorf("render: prop key at index %d is not a string", i)
		}
		ctx.Set(key, html.EscapeString(fmt.Sprintf("%v", props[i+1])))
	}
	return fn(ctx)
}
