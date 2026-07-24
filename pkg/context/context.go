package context

import "net/http"

type Context struct {
	W    http.ResponseWriter
	R    *http.Request
	data map[string]string
}

func (c *Context) Set(key, value string) {
	if c.data == nil {
		c.data = make(map[string]string)
	}
	c.data[key] = value
}

func (c *Context) Get(key string) string {
	if c.data == nil {
		return ""
	}
	return c.data[key]
}

func (c *Context) Param(key string) string {
	return c.R.PathValue(key)
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) FormValue(key string) string {
	if err := c.R.ParseForm(); err != nil {
		return ""
	}
	return c.R.FormValue(key)
}
