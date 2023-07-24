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

func TestAccount(t *testing.T) {
	store, err := NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}
	server := &APIServer{store: store}
	var createAccountResponse struct {
		ID string `json:"id"`
	}
	t.Run("CreateAccount", func(t *testing.T) {
		createcAccReqBody := &CreateAccountRequest{FirstName: "Test User First Name", LastName: "Test User Last Name"}
		reqBodyJSON, err := json.Marshal(createcAccReqBody)
		if err != nil {
			t.Fatalf("Error failed json serialized request body %v", err)
		}

		req, err := http.NewRequest("POST", settings.AppSettings.Create_Account_Route, bytes.NewBuffer(reqBodyJSON))
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(server.handleCreateAccount)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, rr.Code)
		}

		if err := json.Unmarshal(rr.Body.Bytes(), &createAccountResponse); err != nil {
			t.Errorf("failed to unmarshal response body, the id must be uuid type: %s", err)
		}

		expectIdToBeUUID := true
		if expectIdToBeUUID != isValidUUID(createAccountResponse.ID) {
			t.Errorf("the id return must be an uuid. Current id: %s", createAccountResponse.ID)
		}
	})

	t.Run("GetAccount", func(t *testing.T) {
		accountId := createAccountResponse.ID
		req, err := http.NewRequest("GET", "v1/account/"+accountId, nil)
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(server.handleAccount)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		expectHttpStatus := http.StatusFound
		if rr.Code != expectHttpStatus {
			t.Errorf("expected status code %d but got %d", expectHttpStatus, rr.Code)
		}
	})
}
