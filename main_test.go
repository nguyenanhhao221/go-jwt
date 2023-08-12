package main

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
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
