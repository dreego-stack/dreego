module demo-wailsv3

go 1.27

require (
	github.com/dreego-stack/dreego/adapter/wails v0.8.0
	github.com/dreego-stack/dreego/core v0.8.0
	github.com/wailsapp/wails/v3 v3.0.0-beta.20
)

require (
	github.com/adrg/xdg v0.5.3 // indirect
	github.com/coder/websocket v1.8.14 // indirect
	github.com/dreego-stack/dreego v0.8.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.46.0 // indirect
	golang.org/x/text v0.39.0 // indirect
)

replace github.com/dreego-stack/dreego => ../..

replace github.com/dreego-stack/dreego/core => ../../core

replace github.com/dreego-stack/dreego/adapter/wails => ../../adapter/wails
