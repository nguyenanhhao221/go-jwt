package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if statusCode == http.StatusNoContent {
		return
	}
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		log.Printf("Error: failed to encode the payload %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func WriteErrorJson(w http.ResponseWriter, statusCode int, msg string) {
	type ApiError struct {
		Error string `json:"error"`
	}
	WriteJSON(w, statusCode, ApiError{Error: msg})
}
