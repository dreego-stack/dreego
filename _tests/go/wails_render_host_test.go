package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dreego-stack/dreego/adapter/wails"
	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/dreegotest"
)

func TestWailsRenderMatchesSSRForRegisteredPage(t *testing.T) {
	page := dreego.ComponentFunc(func(dreego.RenderContext) (dreego.Result, error) {
		return dreego.Result{HTML: []byte("<main><h1>Timer</h1></main>")}, nil
	})
	app := dreego.New()
	if err := app.RegisterRender("/timer", page); err != nil {
		t.Fatalf("RegisterRender: %v", err)
	}
	if err := app.Register("GET", "/timer", func(w http.ResponseWriter, _ *http.Request) {
		result, err := dreego.Render(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(result.HTML)
	}); err != nil {
		t.Fatalf("Register: %v", err)
	}

	response := httptest.NewRecorder()
	app.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/timer", nil))
	ssrResult, err := dreego.Render(page)
	if err != nil {
		t.Fatalf("SSR render: %v", err)
	}
	host, err := wails.New(app)
	if err != nil {
		t.Fatalf("wails.New: %v", err)
	}
	wailsResponse := httptest.NewRecorder()
	host.ServeHTTP(wailsResponse, httptest.NewRequest(http.MethodGet, "/timer", nil))
	dreegotest.MustEqual(t, string(response.Body.Bytes()), string(ssrResult.HTML))
	dreegotest.MustEqual(t, wailsResponse.Body.String(), string(ssrResult.HTML))
}

func TestWailsRejectsNilApp(t *testing.T) {
	if _, err := wails.New(nil); err == nil {
		t.Fatal("wails.New(nil) must return an error")
	}
}

func TestWailsRejectsHTTPOnlyRoute(t *testing.T) {
	app := dreego.New()
	if err := app.Register("GET", "/http-only", func(http.ResponseWriter, *http.Request) {}); err != nil {
		t.Fatalf("Register: %v", err)
	}
	host, err := wails.New(app)
	if err != nil {
		t.Fatalf("wails.New: %v", err)
	}
	response := httptest.NewRecorder()
	host.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/http-only", nil))
	if response.Code != http.StatusNotFound {
		t.Fatal("Wails must reject an HTTP-only route without a render registration")
	}
}

func TestGeneratedPureGetPageRegistersRenderComponent(t *testing.T) {
	generated := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": `<body><main><h1>Timer</h1></main></body>`,
	})
	routes := generated["www/routes/dree.go"]
	dreegotest.MustContain(t, routes, `app.RegisterRender("/", PageIndex())`)
}

func TestGeneratedHTTPPageDoesNotRegisterRenderComponent(t *testing.T) {
	generated := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": `<server>message := "HTTP only"</server><body>{{ message }}</body>`,
	})
	routes := generated["www/routes/dree.go"]
	dreegotest.MustNotContain(t, routes, "app.RegisterRender(")
}
