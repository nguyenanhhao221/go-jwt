package auth

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// validateJWT Validate the token string
// Depends on the signing method we choose for the JWT signing
// We will parse the token received from the client and using the secret hash to check if it valid
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// The secret hash to be used
	hmacSecret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(hmacSecret), nil
	})
}

func CreateJWT(accountId uuid.UUID) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	// TODO: Expire time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"accountId": accountId,
	})

	// Sign and get the complete encoded token as a string using the secret
	hmacSecret := os.Getenv("JWT_SECRET")
	log.Println(hmacSecret)
	if tokenString, err := token.SignedString([]byte(hmacSecret)); err != nil {
		log.Printf("Error failed to sign token %v", err)
		return "", err
	} else {
		return tokenString, nil
	}
}
