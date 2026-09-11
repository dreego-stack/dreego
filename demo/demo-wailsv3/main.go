package main

import (
	"log"

	"demo-wailsv3/app"
	wailsadapter "github.com/dreego-stack/dreego/adapter/wails"
	dreego "github.com/dreego-stack/dreego/core"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	dreegoApp := dreego.New()
	if err := app.Register(dreegoApp); err != nil {
		log.Fatal(err)
	}
	handler, err := wailsadapter.New(dreegoApp)
	if err != nil {
		log.Fatal(err)
	}
	wailsApp := application.New(application.Options{
		Name:   "Dreego Wails Timer",
		Assets: application.AssetOptions{Handler: handler, DisableLogging: true},
		Services: []application.Service{
			application.NewService(app.NewTimerService(5 * 60)),
		},
	})
	wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Focus Timer",
		URL:       "/",
		Width:     640,
		Height:    560,
		MinWidth:  320,
		MinHeight: 480,
	})
	if err := wailsApp.Run(); err != nil {
		log.Fatal(err)
	}
}
