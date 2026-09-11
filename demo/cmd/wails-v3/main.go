package main

import (
	"log"

	wailsdemo "demo/wails-v3"
)

func main() {
	if err := wailsdemo.Run(); err != nil {
		log.Fatal(err)
	}
}
