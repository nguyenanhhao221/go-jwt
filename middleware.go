package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
)

// withJWTAuth Middleware to validate the JWT token in the client request
func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking JWT Auth")
		_, err := validateJWT(r.Header.Get("x-jwt-token"))
		if err != nil {
			WriteErrorJson(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}

		handlerFunc(w, r)
	}
}

// validateJWT Validate the token string
// Depends on the signing method we choose for the JWT signing
// We will parse the token received from the client and using the secret hash to check if it valid
func validateJWT(tokenString string) (*jwt.Token, error) {
	// The secret hash to be used
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}
