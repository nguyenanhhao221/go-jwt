package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env")
	}

	// The URL in .env to connect to SQL database
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	/*  Open a connection to the postgres database */
	// sql.Open is used to establish a connection to the PostgreSQL database.
	// However, the sql.Open function only creates a connection object, it doesn't actually establish a connection to the database.
	// In order to use this we also need to import "github.com/lib/pq"
	sqlConnection, dbConnErr := sql.Open("postgres", dbURL)
	if dbConnErr != nil {
		log.Fatal("Cannot connect to the database: ", dbConnErr)
	}
	// sql.Open successfully returns an instance of sql.DB regardless of whether the database server is running or not.
	// To check if the connection was successful, you need to call the Ping method on the sql.DB instance.
	if err := sqlConnection.Ping(); err != nil {
		log.Fatal("Failed to ping the database, did you forget to run Docker? Error: ", err)
	}

	apiSrv := NewAPIServer("8080")
	apiSrv.Run()
}
