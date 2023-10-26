// Package main contains a simple program to generate API keys for inquiry
// which can be used to configure and protect the API.
package main

import (
	"log"

	"github.com/inquiryproj/inquiry/pkg/crypto"
)

func main() {
	apiKey, err := crypto.NewAPIKey()
	if err != nil {
		log.Fatal("unable to generate API Key", err)
	}
	log.Default().Printf("api key generated: %s\n", apiKey)
}
