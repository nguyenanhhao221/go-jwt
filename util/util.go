package util

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(inputPassword, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(inputPassword))
	return err == nil
}

func GetIdFromRequest(r *http.Request) (uuid.UUID, error) {
	if accountId, err := uuid.Parse(chi.URLParam(r, "accountId")); err != nil {
		log.Printf("Failed to get account id from request %v", err)
		return uuid.Nil, err
	} else {
		return accountId, nil
	}
}

// Determine the host name to expose our app on local and Railway deploy
func GetHostString(portAsString string) string {
	if _, exist := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exist {
		return ":" + portAsString
	} else {
		return "0.0.0.0:" + portAsString
	}
}
