package context

import (
	gcontext "context"
	"net/http"
)

type Context interface {
	gcontext.Context
	Param(name string) string
	Data(key string) any
	Render(name string, data any) error
}

type SSRContext struct {
	gcontext.Context
	W    http.ResponseWriter
	R    *http.Request
	data map[string]any
}

func NewSSR(w http.ResponseWriter, r *http.Request) *SSRContext {
	return &SSRContext{
		Context: r.Context(),
		W:       w,
		R:       r,
		data:    make(map[string]any),
	}
}

func (c *SSRContext) Param(name string) string {
	return c.R.PathValue(name)
}

func (c *SSRContext) Data(key string) any {
	if c.data == nil {
		return nil
	}
	return c.data[key]
}

func (c *SSRContext) Set(key string, value any) {
	if c.data == nil {
		c.data = make(map[string]any)
	}
	c.data[key] = value
}

func (c *SSRContext) Get(key string) string {
	if c.data == nil {
		return ""
	}
	s, _ := c.data[key].(string)
	return s
}

func (c *SSRContext) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *SSRContext) FormValue(key string) string {
	if err := c.R.ParseForm(); err != nil {
		return ""
	}
	return c.R.FormValue(key)
}

func (c *SSRContext) Render(name string, data any) error {
	return nil
}
