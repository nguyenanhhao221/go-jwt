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

type CustomJWTClaims struct {
	ID uuid.UUID `json:"id"`
	jwt.StandardClaims
}

func CreateJWT(accountId uuid.UUID) (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	// TODO: Expire time
	claims := CustomJWTClaims{
		ID: accountId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	hmacSecret := os.Getenv("JWT_SECRET")
	log.Println(hmacSecret)
	// Here we have to convert the secret into []bytes slice
	// Check out the signature of the SignedString it expect an interface{} type, however this is just a bait
	// Because base on different SigningMethod we choose, the value pass in need to be some specific types, so the library just put interface{} for now
	// Read more: https://github.com/dgrijalva/jwt-go/issues/65
	if tokenString, err := token.SignedString([]byte(hmacSecret)); err != nil {
		log.Printf("Error failed to sign token %v", err)
		return "", err
	} else {
		return tokenString, nil
	}
}

// func GetUserIDFromJWT(tokenString string) (uuid.UUID, error) {
// 	if token, err := jwt.ParseWithClaims(tokenString); err != nil {
// 		return uuid.Nil, err
// 	}
// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		fmt.Println(claims["foo"], claims["nbf"])
// 	} else {
// 		fmt.Println(err)
// 	}
// 	return "11231y239812749812749817kjbdwiugfiuwegfiuwe", nil
// }
