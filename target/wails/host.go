package wails

import (
	"errors"
	"fmt"
	"os"

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
	if os.Getenv("FRONTEND_DEVSERVER_URL") != "" {
		return errors.New("dreego wails: FRONTEND_DEVSERVER_URL is not supported")
	}
	if _, err := h.Render(options.Path); err != nil {
		return err
	}
	app := application.New(application.Options{
		Name: options.Name,
		Assets: application.AssetOptions{
			Handler:        h,
			DisableLogging: true,
		},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     options.Title,
		Width:     options.Width,
		Height:    options.Height,
		MinWidth:  options.MinWidth,
		MinHeight: options.MinHeight,
		URL:       options.Path,
	})
	if err := app.Run(); err != nil {
		return fmt.Errorf("dreego wails: run: %w", err)
	}
	return nil
}
