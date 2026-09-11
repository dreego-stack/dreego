package wails_v3

import (
	dreego "github.com/dreego-stack/dreego/core"
	wailstarget "github.com/dreego-stack/dreego/target/wails"
)

func Run() error {
	app := dreego.New()
	if err := Register(app); err != nil {
		return err
	}
	host, err := wailstarget.New(app)
	if err != nil {
		return err
	}
	return host.Run(wailstarget.Options{
		Name:      "Dreego Wails Timer",
		Title:     "Focus Timer",
		Path:      "/",
		Width:     640,
		Height:    560,
		MinWidth:  320,
		MinHeight: 480,
	})
}
