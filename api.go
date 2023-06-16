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
}

func NewAPIServer(listenAdd string) *APIServer {
	return &APIServer{
		listenAdd: listenAdd,
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
	// so if we access /v1/healthz the handlerReadiness will be called
	router.Mount("/v1", v1Router)

	//Handlers
	v1Router.Get("/health", s.handlerReadiness)

	//Start the server
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
	w.WriteHeader(200)
	w.Write([]byte("Ready"))
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
