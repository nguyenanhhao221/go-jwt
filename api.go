package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/nguyenanhhao221/go-jwt/settings"
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
	v1Router.Get(settings.AppSettings.Account_Route, s.handleAccount)
	v1Router.Post(settings.AppSettings.Create_Account_Route, s.handleCreateAccount)

	// Start the server
	server := &http.Server{
		Addr:    ":" + s.listenAdd,
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

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.handleGetAccount(w, r)
		return
	} else if r.Method == "DELETE" {
		s.handleDeleteAccount(w, r)
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "accountId")
	log.Println(accountId)
	account := NewAccount("Hao", "Nguyen")

	WriteJSON(w, http.StatusFound, account)
}

// handleCreateAccount create a new account with first name and last name from client's post request
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		WriteErrorJson(w, http.StatusBadRequest, err.Error())
	}
	newAccount := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if id, err := s.store.CreateAccount(newAccount); err != nil {
		WriteErrorJson(w, http.StatusBadRequest, err.Error())
	} else {
		type createAccountRes struct {
			ID uuid.UUID `json:"id"`
		}
		WriteJSON(w, http.StatusCreated, &createAccountRes{
			ID: id,
		})
	}
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
