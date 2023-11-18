package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nguyenanhhao221/go-jwt/internal/auth"
	"github.com/nguyenanhhao221/go-jwt/settings"
	"github.com/nguyenanhhao221/go-jwt/util"
)

type APIServer struct {
	listenAdd string
	store     Storage
}

func NewAPIServer(listenAdd string, store Storage) *APIServer {
	return &APIServer{
		listenAdd: listenAdd,
		store:     store,
	}
}

func (s *APIServer) Run() {
	// start a router
	router := chi.NewRouter()

	// Allow cors
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// A good base middleware stack
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	// Add router handler for v1
	v1Router := chi.NewRouter()

	// mount the v1Router to the /v1 route
	router.Mount(settings.AppSettings.API_V1, v1Router)

	// Handlers
	v1Router.Get(settings.AppSettings.Check_Health, s.handlerReadiness)
	v1Router.Get(settings.AppSettings.Account_Route, withJWTAuth(s.handleAccount))
	v1Router.Get(settings.AppSettings.All_Account_Route, s.handleGetAllAccount)
	v1Router.Put(settings.AppSettings.Account_Route, s.handleAccount)
	v1Router.Delete(settings.AppSettings.Account_Route, s.handleAccount)
	v1Router.Post(settings.AppSettings.Create_Account_Route, s.handleCreateAccount)
	v1Router.Post(settings.AppSettings.SignIn_Account_Route, s.handleSignIn)
	v1Router.Post(settings.AppSettings.Transfer_Route, withJWTAuth(s.handleTransfer))

	// Start the server
	server := &http.Server{
		Addr:    s.listenAdd,
		Handler: router,
	}
	log.Printf("Server is listening on port %v", s.listenAdd)
	serverErr := server.ListenAndServe()
	if serverErr != nil {
		log.Fatalf("Error: Failed to start server %v", serverErr)
	}
}

func (s *APIServer) handlerReadiness(w http.ResponseWriter, r *http.Request) {
	type Ready struct {
		Status string `json:"status"`
	}
	WriteJSON(w, http.StatusOK, Ready{Status: "alive"})
}

func (s *APIServer) handleGetAllAccount(w http.ResponseWriter, r *http.Request) {
	if allAccounts, err := s.store.GetAllAccounts(); err != nil {
		WriteErrorJson(w, http.StatusInternalServerError, err.Error())
	} else {
		WriteJSON(w, http.StatusFound, allAccounts)
	}
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	accountId, err := util.GetIdFromRequest(r)
	if err != nil {
		WriteErrorJson(w, http.StatusBadRequest, err.Error())
	}
	if r.Method == http.MethodGet {
		s.handleGetAccount(w, r, accountId)
		return
	} else if r.Method == http.MethodDelete {
		s.handleDeleteAccount(w, r, accountId)
		return
	} else if r.Method == http.MethodPut {
		s.handleUpdateAccount(w, r, accountId)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request, accountId uuid.UUID) {
	if account, err := s.store.GetAccountById(accountId); err != nil {
		WriteErrorJson(w, http.StatusNotFound, err.Error())
		return
	} else {
		WriteJSON(w, http.StatusOK, account)
		return
	}
}

// handleCreateAccount create a new account with first name and last name from client's post request
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		WriteErrorJson(w, http.StatusForbidden, err.Error())
		return
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	type IError struct {
		Field string `json:"field"`
		Tag   string `json:"tag"`
		Value string `json:"value"`
	}
	var errors []*IError
	if err := validate.Struct(createAccountReq); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var el IError
			el.Field = err.Field()
			el.Tag = err.Tag()
			el.Value = err.Param()
			errors = append(errors, &el)
		}
		// Convert errors to a JSON string
		errorsJSON, _ := json.Marshal(errors)
		WriteErrorJson(w, http.StatusBadRequest, string(errorsJSON))
		return
	}
	newAccount := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Email, createAccountReq.Password)
	if _, err := s.store.GetAccountByEmail(createAccountReq.Email); err == nil {
		log.Printf("Error while checking account email %v", err)
		WriteErrorJson(w, http.StatusForbidden, "Email already existed")
		return

	}
	if id, err := s.store.CreateAccount(newAccount); err != nil {
		log.Printf("Error while creating account %v", err)
		WriteErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		type createAccountRes struct {
			ID uuid.UUID `json:"id"`
		}
		WriteJSON(w, http.StatusCreated, &createAccountRes{
			ID: id,
		})
	}
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request, accountId uuid.UUID) {
	updateAccountReq := new(Account)

	if err := json.NewDecoder(r.Body).Decode(updateAccountReq); err != nil {
		WriteErrorJson(w, http.StatusForbidden, err.Error())
		return
	}
	if err := s.store.UpdateAccountById(updateAccountReq, accountId); err != nil {
		WriteErrorJson(w, http.StatusNotFound, err.Error())
		return
	} else {
		WriteJSON(w, http.StatusNoContent, nil)
		return
	}
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request, accountId uuid.UUID) {
	if err := s.store.DeleteAccountById(accountId); err != nil {
		WriteErrorJson(w, http.StatusNotFound, err.Error())
	} else {
		WriteJSON(w, http.StatusNoContent, 1)
	}
}

func (s *APIServer) handleSignIn(w http.ResponseWriter, r *http.Request) {
	type SignInReqBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	signInReqBody := new(SignInReqBody)
	if err := json.NewDecoder(r.Body).Decode(signInReqBody); err != nil {
		WriteErrorJson(w, http.StatusForbidden, err.Error())
		return
	}

	if account, err := s.store.GetAccountByEmail(signInReqBody.Email); err != nil {
		WriteErrorJson(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		isPasswordMatch := util.CheckPasswordHash(signInReqBody.Password, account.Password)
		if !isPasswordMatch {
			WriteErrorJson(w, http.StatusUnauthorized, "Wrong email or password")
			return
		}
		if jwtToken, err := auth.CreateJWT(account.ID); err != nil {
			WriteErrorJson(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create JWT token %v", err))
			return
		} else {
			type TokenRes struct {
				Token string `json:"token"`
			}
			WriteJSON(w, http.StatusOK, TokenRes{Token: jwtToken})
		}
		return
	}
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) {
	type TransferBalanceBody struct {
		Number  int64 `json:"toAccount"`
		Balance int64 `json:"balance"`
	}
	transferBalanceReq := new(TransferBalanceBody)
	if err := json.NewDecoder(r.Body).Decode(transferBalanceReq); err != nil {
		WriteErrorJson(w, http.StatusForbidden, err.Error())
	}
	WriteJSON(w, http.StatusCreated, transferBalanceReq)
}
