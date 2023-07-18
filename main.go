package main

import (
	"log"
)

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to get Postgres sql connection %v", err)
	}

	apiSrv := NewAPIServer("8080", store)
	apiSrv.Run()
}
