package server

import (
	"errors"
	"fmt"
	"path"

	rendercontext "github.com/dreego-stack/dreego/core/internal/context"
	rendercore "github.com/dreego-stack/dreego/core/internal/render"
)

var ErrRenderRouteNotFound = errors.New("dreego: render route not found")

func (a *App) RegisterRender(routePath string, component rendercore.Renderable) error {
	if component == nil {
		return errors.New("dreego: render component is nil")
	}
	if routePath == "" || routePath[0] != '/' || path.Clean(routePath) != routePath {
		return fmt.Errorf("dreego: invalid render route %q", routePath)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.mutable(); err != nil {
		return err
	}
	if _, exists := a.renderPages[routePath]; exists {
		return fmt.Errorf("%w: GET %s", ErrRouteConflict, routePath)
	}
	a.renderPages[routePath] = component
	return nil
}

func (a *App) RenderPage(routePath string) (rendercore.Result, error) {
	a.mu.RLock()
	component, exists := a.renderPages[routePath]
	a.mu.RUnlock()
	if !exists {
		return rendercore.Result{}, fmt.Errorf("%w: %s", ErrRenderRouteNotFound, routePath)
	}
	return component.Render(rendercontext.NewRender(nil))
}
