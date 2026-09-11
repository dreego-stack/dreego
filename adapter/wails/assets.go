package wails

import (
	"errors"
	"net/http"
	"path"
	"strings"

	dreego "github.com/dreego-stack/dreego/core"
)

type handler struct {
	app *dreego.App
}

func New(app *dreego.App) (http.Handler, error) {
	if app == nil {
		return nil, errors.New("dreego wails: app is nil")
	}
	return &handler{app: app}, nil
}

func (h *handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet && request.Method != http.MethodHead {
		response.Header().Set("Allow", "GET, HEAD")
		http.Error(response, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if !validAssetRequest(request) {
		http.NotFound(response, request)
		return
	}
	result, err := h.app.RenderPage(request.URL.Path)
	if err == nil {
		writeAssetResponse(response, request, "text/html; charset=utf-8", result.HTML)
		return
	}
	if !errors.Is(err, dreego.ErrRenderRouteNotFound) {
		http.Error(response, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	asset, err := h.app.StaticAsset(request.URL.Path)
	if err != nil {
		http.NotFound(response, request)
		return
	}
	writeAssetResponse(response, request, asset.MIME, asset.Content)
}

func validAssetRequest(request *http.Request) bool {
	requestPath := request.URL.Path
	return request.URL.RawQuery == "" &&
		strings.HasPrefix(requestPath, "/") &&
		!strings.ContainsAny(requestPath, "\\\x00") &&
		!strings.Contains(request.URL.EscapedPath(), "%") &&
		path.Clean(requestPath) == requestPath
}

func writeAssetResponse(response http.ResponseWriter, request *http.Request, mime string, content []byte) {
	response.Header().Set("Content-Type", mime)
	response.Header().Set("X-Content-Type-Options", "nosniff")
	if request.Method == http.MethodGet {
		_, _ = response.Write(content)
	}
}
