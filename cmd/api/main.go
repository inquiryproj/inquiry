// Package main is the entrypoint for the API server.
package main

import (
	"log"

	factory "github.com/inquiryproj/inquiry/internal"
)

func main() {
	app, err := factory.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
