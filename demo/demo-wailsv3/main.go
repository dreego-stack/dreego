package main

import (
	"log"

	"demo-wailsv3/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
