package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nguyenanhhao221/go-jwt/internal/auth"
	"github.com/nguyenanhhao221/go-jwt/util"
)

// withJWTAuth Middleware to validate the JWT token in the client request
func withJWTAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Checking JWT Auth")
		tokenString := r.Header.Get("x-jwt-token")
		token, err := auth.ValidateJWT(tokenString)
		if err != nil {
			WriteErrorJson(w, http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			return
		}
		if !token.Valid {
			WriteErrorJson(w, http.StatusForbidden, "Invalid token")
			return
		}

		claims, ok := token.Claims.(*auth.CustomJWTClaims)

		if !ok {
			fmt.Println("Token claims are not of type CustomJWTClaims")
			WriteErrorJson(w, http.StatusForbidden, "Invalid token claims")
			return
		}

		acocuntIdFromReq, err := util.GetIdFromRequest(r)
		if err != nil {
			WriteErrorJson(w, http.StatusBadRequest, fmt.Sprintf("Invalid id in request: %v", err))
			return
		}

		if claims.ID != acocuntIdFromReq {
			WriteErrorJson(w, http.StatusUnauthorized, "Permission Denied")
			return
		} else {
			next(w, r)
		}
	}
}
