package main

import (
	"log"
	"strconv"

	"github.com/nguyenanhhao221/go-jwt/settings"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to get Postgres sql connection %v", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	portAsString := strconv.Itoa(settings.AppSettings.PORT)
	apiSrv := NewAPIServer(portAsString, store)
	apiSrv.Run()
}
