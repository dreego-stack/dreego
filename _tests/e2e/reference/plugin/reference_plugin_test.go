package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestReferencePlugin(t *testing.T) {
	t.Parallel()
	pluginSrc := `package plugin

import (
	"fmt"
	"net/http"

	dreego "github.com/dreego-stack/dreego/core"
)

type Options struct {
	Prefix string
}

func Register(app *dreego.App, options Options) error {
	if options.Prefix == "" {
		options.Prefix = "/plugin"
	}
	if err := app.Register(http.MethodGet, options.Prefix+"/hello", helloHandler); err != nil {
		return err
	}
	if err := app.Register(http.MethodGet, options.Prefix+"/hello/{id}", helloIDHandler); err != nil {
		return err
	}
	return app.Register(http.MethodGet, options.Prefix+"/health", healthHandler)
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hello from the plugin"))
}

func helloIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hello %s", r.PathValue("id"))
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("plugin ok"))
}
`
	pluginMain := `package main

import (
	"os"

	"github.com/dreego-stack/dreego/adapter/ssr"
	dreego "github.com/dreego-stack/dreego/core"
	"t/plugin"
	"t/www"
)

func main() {
	app := dreego.New()
	if err := plugin.Register(app, plugin.Options{Prefix: "/plugin"}); err != nil {
		panic(err)
	}
	if err := www.Register(app); err != nil {
		panic(err)
	}
	if err := ssr.Listen(app, os.Getenv("DREEGO_TEST_ADDR")); err != nil {
		panic(err)
	}
}
`
	c := dreegotest.Serve(t, map[string]string{
		"main.go":                pluginMain,
		"www/dreego.config.json": `{"logging":{"enabled":false},"redirects":[],"rewrites":[]}`,
		"plugin/plugin.go":       pluginSrc,
		"www/routes/+page.dreego": `<head>
    <title>Plugin demo</title>
</head>

<body>
    <h1>Plugin demo</h1>
    <p>This app registers a local plugin package before the generated routes.</p>
</body>`,
	})

	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Plugin demo") {
		t.Fatalf("plugin home missing heading: %s", body)
	}
	code, body = c.Get(t, "/plugin/hello")
	if code != 200 {
		t.Fatalf("GET /plugin/hello = %d, want 200", code)
	}
	if !strings.Contains(body, "Hello from the plugin") {
		t.Fatalf("plugin route missing body: %s", body)
	}
	code, body = c.Get(t, "/plugin/hello/42")
	if code != 200 {
		t.Fatalf("GET /plugin/hello/42 = %d, want 200", code)
	}
	if !strings.Contains(body, "Hello 42") {
		t.Fatalf("plugin dynamic route missing param: %s", body)
	}
	code, body = c.Get(t, "/plugin/health")
	if code != 200 {
		t.Fatalf("GET /plugin/health = %d, want 200", code)
	}
	if !strings.Contains(body, "plugin ok") {
		t.Fatalf("plugin health route missing body: %s", body)
	}
}
