package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nguyenanhhao221/go-jwt/internal/auth"
)

// withJWTAuth Middleware to validate the JWT token in the client request
func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking JWT Auth")
		_, err := auth.ValidateJWT(r.Header.Get("x-jwt-token"))
		if err != nil {
			WriteErrorJson(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		handlerFunc(w, r)
	}
}
