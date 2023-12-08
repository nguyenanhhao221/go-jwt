package main

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	if _, exist := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exist {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error while loading env %v", err)
		}
	}
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatalf("failed to init postgres store %v", err)
	}

	if err := store.Init(); err != nil {
		log.Fatalf("failed to init store %v", err)
	}
	// Run all tests
	log.Print("TestMain run, initialize database and create necessary table")
	exitCode := m.Run()
	// Exit with the appropriate exit code
	os.Exit(exitCode)
}
