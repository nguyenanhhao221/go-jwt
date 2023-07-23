package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/nguyenanhhao221/go-jwt/settings"
)

func isValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

func TestCreateAccount(t *testing.T) {
	store, err := NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}
	server := &APIServer{store: store}

	createcAccReqBody := &CreateAccountRequest{FirstName: "Test User First Name", LastName: "Test User Last Name"}
	reqBodyJSON, err := json.Marshal(createcAccReqBody)
	if err != nil {
		t.Fatalf("Error failed json serialized request body %v", err)
	}

	req, err := http.NewRequest("POST", settings.AppSettings.Create_Account_Route, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handleCreateAccount)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status code %d but got %d", http.StatusCreated, rr.Code)
	}

	var response struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to unmarshal response body, the id must be uuid type: %s", err)
	}

	expectIdToBeUUID := true
	if expectIdToBeUUID != isValidUUID(response.ID) {
		t.Errorf("the id return must be an uuid. Current id: %s", response.ID)
	}
}
