package core

type Component interface {
	Render(ctx *SSRContext) (string, error)
}

type ComponentFunc func(ctx *SSRContext) (string, error)

func (f ComponentFunc) Render(ctx *SSRContext) (string, error) {
	return f(ctx)
}
