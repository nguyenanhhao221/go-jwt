package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nguyenanhhao221/go-jwt/settings"
)

func TestCheckHealth(t *testing.T) {
	server := &APIServer{}
	req, err := http.NewRequest("GET", settings.AppSettings.API_V1+"/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handlerReadiness)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d but got %d", http.StatusOK, rr.Code)
	}

	var response struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response body: %s", err)
	}
	expectedResponse := "alive"
	if response.Status != expectedResponse {
		t.Errorf("expected body %s but got %s", expectedResponse, rr.Body.String())
	}
}
