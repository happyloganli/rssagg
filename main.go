package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {

	// Load the environment
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set")
	}

	// Set up the routers
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Set up version 1 routers
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", ReadinessHandler)
	v1Router.Get("/errorz", ErrorHandler)

	router.Mount("/v1", v1Router)

	// Set up the http server on the specific port
	httpServer := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Starting server on port %s...", portString)

	// Start http listen on the specific port
	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
