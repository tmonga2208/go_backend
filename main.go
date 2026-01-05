package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/tarunmonga/hello-world/db"
	"github.com/tarunmonga/hello-world/handlers"
	authMiddleware "github.com/tarunmonga/hello-world/middleware"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	// 2. Connect to Database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}
	db.Connect(dbURL)
	defer db.Conn.Close(context.Background())

	// Create tables
	db.CreateTable()

	// 3. Setup Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS Configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API Secure & Ready (Modular Structure)"))
	})

	// Public Routes
	r.Post("/login", handlers.Login)
	r.Post("/users", handlers.CreateUser)

	// Protected Routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware)
		r.Get("/users", handlers.GetUsers)
		r.Get("/me", handlers.GetMe)
		r.Put("/users/{id}", handlers.UpdateUser)
	})

	fmt.Println("Server running on :3333")
	http.ListenAndServe(":3333", r)
}
