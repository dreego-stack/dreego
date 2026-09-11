package wails

import (
	"errors"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type Host struct {
	app *dreego.App
}

type Options struct {
	Name      string
	Title     string
	Path      string
	Width     int
	Height    int
	MinWidth  int
	MinHeight int
}

func New(app *dreego.App) (*Host, error) {
	if app == nil {
		return nil, errors.New("dreego wails: app is nil")
	}
	return &Host{app: app}, nil
}

func (h *Host) Render(path string) (dreego.Result, error) {
	return h.app.RenderPage(path)
}

func (h *Host) Run(options Options) error {
	result, err := h.Render(options.Path)
	if err != nil {
		return err
	}
	app := application.New(application.Options{Name: options.Name})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     options.Title,
		Width:     options.Width,
		Height:    options.Height,
		MinWidth:  options.MinWidth,
		MinHeight: options.MinHeight,
		HTML:      string(result.HTML),
	})
	return app.Run()
}
