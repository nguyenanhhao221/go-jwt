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

func TestAccountCI(t *testing.T) {
	store, err := NewPostgresStore()
	if err != nil {
		t.Fatal(err)
	}
	server := &APIServer{store: store}
	var createAccountResponse struct {
		ID uuid.UUID `json:"id"`
	}
	mockUser := &CreateAccountRequest{FirstName: "Test User First Name", LastName: "Test User Last Name"}

	t.Run("CreateAccount", func(t *testing.T) {
		createcAccReqBody := mockUser
		reqBodyJSON, err := json.Marshal(createcAccReqBody)
		if err != nil {
			t.Fatalf("Error failed json serialized request body %v", err)
		}

		req, err := http.NewRequest("POST", settings.AppSettings.Create_Account_Route, bytes.NewBuffer(reqBodyJSON))
		if err != nil {
			t.Fatal(err)
		}

		handler := http.HandlerFunc(server.handleCreateAccount)
		testRecorder := httptest.NewRecorder()

		handler.ServeHTTP(testRecorder, req)

		if testRecorder.Code != http.StatusCreated {
			t.Errorf("expected status code %d but got %d", http.StatusCreated, testRecorder.Code)
		}

		if err := json.Unmarshal(testRecorder.Body.Bytes(), &createAccountResponse); err != nil {
			t.Errorf("failed to unmarshal response body, the id must be uuid type: %s", err)
		}

		expectIdToBeUUID := true
		if expectIdToBeUUID != isValidUUID(createAccountResponse.ID.String()) {
			t.Errorf("the id return must be an uuid. Current id: %s", createAccountResponse.ID)
		}
	})
}

// func TestAccount(t *testing.T) {
// 	store, err := NewPostgresStore()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	server := &APIServer{store: store}
// 	var createAccountResponse struct {
// 		ID uuid.UUID `json:"id"`
// 	}
// 	mockUser := &CreateAccountRequest{FirstName: "Test User First Name", LastName: "Test User Last Name"}
//
// 	t.Run("CreateAccount", func(t *testing.T) {
// 		createcAccReqBody := mockUser
// 		reqBodyJSON, err := json.Marshal(createcAccReqBody)
// 		if err != nil {
// 			t.Fatalf("Error failed json serialized request body %v", err)
// 		}
//
// 		req, err := http.NewRequest("POST", settings.AppSettings.Create_Account_Route, bytes.NewBuffer(reqBodyJSON))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
//
// 		handler := http.HandlerFunc(server.handleCreateAccount)
// 		rr := httptest.NewRecorder()
//
// 		handler.ServeHTTP(rr, req)
//
// 		if rr.Code != http.StatusCreated {
// 			t.Errorf("expected status code %d but got %d", http.StatusCreated, rr.Code)
// 		}
//
// 		if err := json.Unmarshal(rr.Body.Bytes(), &createAccountResponse); err != nil {
// 			t.Errorf("failed to unmarshal response body, the id must be uuid type: %s", err)
// 		}
//
// 		expectIdToBeUUID := true
// 		if expectIdToBeUUID != isValidUUID(createAccountResponse.ID.String()) {
// 			t.Errorf("the id return must be an uuid. Current id: %s", createAccountResponse.ID)
// 		}
// 	})
//
// 	t.Run("GetAccount", func(t *testing.T) {
// 		accountId := createAccountResponse.ID
// 		req, err := http.NewRequest("GET", "v1/account/"+accountId.String(), nil)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		// NOTE: Because of Chi, if we need to test url path that have param variable we must add the context manually like this
// 		// So the handler can pick up the param (accountId).
// 		// Make sure we import the chi package with same version in the test context and the handler main code
// 		rctx := chi.NewRouteContext()
// 		rctx.URLParams.Add("accountId", accountId.String())
//
// 		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
//
// 		handler := http.HandlerFunc(server.handleAccount)
//
// 		rr := httptest.NewRecorder()
// 		handler.ServeHTTP(rr, req)
//
// 		expectHttpStatus := http.StatusFound
// 		if rr.Code != expectHttpStatus {
// 			t.Errorf("expected status code %d but got %d", expectHttpStatus, rr.Code)
// 		}
//
// 		expectCreatedUser := Account{
// 			ID:        accountId,
// 			FirstName: mockUser.FirstName,
// 			LastName:  mockUser.LastName,
// 			Number:    0,
// 			Balance:   0,
// 		}
//
// 		var responseUser Account
// 		if err := json.NewDecoder(rr.Body).Decode(&responseUser); err != nil {
// 			t.Errorf("Failed to decode response user body %v", err)
// 		}
//
// 		if cmp.Equal(expectCreatedUser, responseUser, cmpopts.IgnoreFields(Account{}, "CreatedAt")) == false {
// 			t.Errorf("expected create user %v but got %v", expectCreatedUser, &responseUser)
// 		}
// 	})
// }
