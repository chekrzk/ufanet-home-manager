package main

import (
	"log"

	"github.com/chekrzk/ufanet-home-manager/requests-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
