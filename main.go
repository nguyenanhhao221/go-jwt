package main

import (
	"log"
)

func main() {
	sqlConnection, err := NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to get Postgres sql connection %v", err)
	}

	apiSrv := NewAPIServer("8080")
	apiSrv.Run()
}
