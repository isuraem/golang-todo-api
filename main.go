// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gorilla/handlers"
// 	"github.com/gorilla/mux"
// 	"github.com/isuraem/todo-api/internal/adapters/app/api"
// 	"github.com/isuraem/todo-api/internal/adapters/app/websocket"
// 	"github.com/isuraem/todo-api/internal/adapters/core/todo"
// 	"github.com/isuraem/todo-api/internal/adapters/core/user"
// 	left "github.com/isuraem/todo-api/internal/adapters/framework/left"
// 	"github.com/isuraem/todo-api/internal/adapters/framework/right/auth"
// 	"github.com/isuraem/todo-api/internal/adapters/framework/right/db"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
// 		os.Getenv("DB_USER"),
// 		os.Getenv("DB_PASSWORD"),
// 		os.Getenv("DB_HOST"),
// 		os.Getenv("DB_PORT"),
// 		os.Getenv("DB_NAME"))

// 	dbAdapter, err := db.NewAdapter(connStr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer dbAdapter.CloseDbConnection()

// 	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET"))

// 	userDB := left.NewUserDB(dbAdapter)
// 	userService := user.NewUserService(userDB, jwtService)
// 	userAPI := api.NewUserAPI(userService)

// 	todoDB := left.NewTodoDB(dbAdapter)
// 	todoService := todo.NewTodoService(todoDB)
// 	hub := websocket.NewHub()
// 	todoAPI := api.NewTodoAPI(todoService, hub) // Ensure the hub is passed here

// 	go hub.Run()

// 	r := mux.NewRouter()
// 	r.HandleFunc("/login", userAPI.Login).Methods("POST")
// 	r.HandleFunc("/register", userAPI.Register).Methods("POST")
// 	r.HandleFunc("/todos", todoAPI.CreateTodo).Methods("POST")
// 	r.HandleFunc("/todos/{id}", todoAPI.UpdateTodo).Methods("PUT")
// 	r.HandleFunc("/todos/{id}", todoAPI.DeleteTodo).Methods("DELETE")
// 	r.HandleFunc("/todos", todoAPI.ListTodos).Methods("GET")
// 	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		websocket.ServeWs(hub, w, r)
// 	})

// 	// Enable CORS for all routes
// 	corsHandler := handlers.CORS(
// 		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
// 		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
// 		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
// 	)(r)

//		log.Fatal(http.ListenAndServe(":8080", corsHandler))
//	}
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
	"github.com/isuraem/todo-api/internal/adapters/framework/right/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	dbAdapter, err := db.NewAdapter(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer dbAdapter.CloseDbConnection()

	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET"))

	userDB := left.NewUserDB(dbAdapter)
	userService := user.NewUserService(userDB, jwtService)
	userAPI := api.NewUserAPI(userService)

	todoDB := left.NewTodoDB(dbAdapter)
	todoService := todo.NewTodoService(todoDB)
	hub := websocket.NewHub()
	todoAPI := api.NewTodoAPI(todoService, hub)

	go hub.Run()

	r := mux.NewRouter()

	// Set up API routes
	api.SetupRoutes(r, userAPI, todoAPI, hub)

	// Enable CORS for all routes
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
	)(r)

	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
