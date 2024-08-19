package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/isuraem/todo-api/internal/adapters/app/api"
	"github.com/isuraem/todo-api/internal/adapters/app/websocket"
	"github.com/isuraem/todo-api/internal/adapters/core/todo"
	"github.com/isuraem/todo-api/internal/adapters/core/user"
	left "github.com/isuraem/todo-api/internal/adapters/framework/left"
	"github.com/isuraem/todo-api/internal/adapters/framework/right/auth"
	"github.com/isuraem/todo-api/internal/adapters/framework/right/cache" // Import the cache package
	"github.com/isuraem/todo-api/internal/adapters/framework/right/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Setup database connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	// Initialize database adapter
	dbAdapter, err := db.NewAdapter(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer dbAdapter.CloseDbConnection()

	// Initialize JWT service
	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET"))

	// Initialize Redis client for caching
	redisClient := cache.NewRedisClient()

	// Initialize user services and APIs
	userDB := left.NewUserDB(dbAdapter)
	userService := user.NewUserService(userDB, jwtService)
	userAPI := api.NewUserAPI(userService)

	// Initialize todo services and APIs with caching
	todoDB := left.NewTodoDB(dbAdapter)
	todoService := todo.NewTodoService(todoDB, redisClient)
	hub := websocket.NewHub()
	todoAPI := api.NewTodoAPI(todoService, hub)

	// Start the websocket hub
	go hub.Run()

	// Set up router
	r := mux.NewRouter()

	// Set up API routes
	api.SetupRoutes(r, userAPI, todoAPI, hub)

	// Enable CORS for all routes
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	)(r)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
