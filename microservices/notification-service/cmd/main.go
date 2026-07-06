package main

import (
	"log"

	"github.com/chekrzk/ufanet-home-manager/notification-service/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
