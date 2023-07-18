package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	router.Mount("/v1", v1Router)

	// Handlers
	v1Router.Get("/health", s.handlerReadiness)
	v1Router.Get("/account/{accountId}", s.handleAccount)

	// Start the server
	server := &http.Server{
		Addr:    ":" + s.listenAdd,
		Handler: router,
	}
	log.Printf("Server is listening on %v", s.listenAdd)
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

	} else if r.Method == "POST" {
		s.handleCreateAccount(w, r)
		return
	} else if r.Method == "DELETE" {
		s.handleDeleteAccount(w, r)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "accountId")
	log.Println(accountId)
	account := NewAccount("Hao", "Nguyen")

	WriteJSON(w, http.StatusFound, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) {
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
