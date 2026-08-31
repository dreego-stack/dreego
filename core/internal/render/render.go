// Package render defines the internal, target-neutral render contract.
// Targets (SSG, Wails) consume Result directly. It must not import net/http.
package render

import (
	"github.com/dreego-stack/dreego/core/internal/context"
)

type Asset struct {
	Name    string
	Content []byte
}

type Result struct {
	HTML   []byte
	Head   []byte
	Assets []Asset
}

func Component(fn func(c *context.SSRContext) (string, error)) (Result, error) {
	ctx := context.NewSSR(nil, nil)
	out, err := fn(ctx)
	if err != nil {
		return Result{}, err
	}
	return Result{HTML: []byte(out)}, nil
}
