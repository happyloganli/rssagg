package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/happyloganli/rssagg/internal/database"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

type apiConfig struct {
	DB *database.Queries
}

func runMigrations(db *sql.DB) {
	// Goose expects the source to be a path to migration files
	if err := goose.Up(db, "sql/schema"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
}

func main() {

	// Load the environment
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT is not found in the environment variables")
	}

	// Get postgres database url
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("$DB_URL is not found in the environment variables")
	}
	// Open a database instance
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can not connect to database: ", err)
	}
	defer db.Close()
	// Run database migration
	runMigrations(db)

	queries := database.New(db)
	apiCfg := apiConfig{
		DB: queries,
	}

	go startScraping(queries, 10, time.Minute)

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
	v1Router.Get("/error", ErrorHandler)
	v1Router.Post("/users", apiCfg.CreateUserHandler)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.GetUserHandler))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.CreateFeedHandler))
	v1Router.Get("/feeds", apiCfg.GetFeedsHandler)
	v1Router.Delete("/feeds/{feedID}", apiCfg.middlewareAuth(apiCfg.DeleteFeedHandler))
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.CreateFeedFollowHandler))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.GetFeedFollowsHandler))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.DeleteFeedFollowsHandler))
	v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.GetUserPostsHandler))

	router.Mount("/v1", v1Router)

	// Set up the http server on the specific port
	httpServer := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Start http listen on the specific port
	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting server on port %s...", portString)
}
