package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nguyenanhhao221/go-jwt/util"
)

func main() {
	// NOTE: As this app is being deploy with railway, we need to handle this.
	// Since railway doens't create a .env file when deploy, but the gotdotenv library expects this file.
	// Without this check gotdotenv will cause the app to Fatal when run deploy
	if _, exist := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exist {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error while loading env %v", err)
		}
	}

	store, err := NewPostgresStore()
	if err != nil {
		log.Fatalf("Failed to get Postgres sql connection %v", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	portAsString := util.GetHostString(port)
	apiSrv := NewAPIServer(portAsString, store)
	apiSrv.Run()
}
