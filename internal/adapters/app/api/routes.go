package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/isuraem/todo-api/internal/adapters/app/websocket"
	"github.com/isuraem/todo-api/internal/middleware"
)

func NewRouter(userAPI *UserAPI, todoAPI *TodoAPI, hub *websocket.Hub) *mux.Router {
	r := mux.NewRouter()
	SetupRoutes(r, userAPI, todoAPI, hub)
	return r
}

// SetupRoutes sets up the API routes and handlers.
func SetupRoutes(r *mux.Router, userAPI *UserAPI, todoAPI *TodoAPI, hub *websocket.Hub) {
	r.HandleFunc("/login", userAPI.Login).Methods("POST")
	r.HandleFunc("/register", userAPI.Register).Methods("POST")

	todos := r.PathPrefix("/todos").Subrouter()
	todos.Use(middleware.JWTAuthMiddleware)
	todos.HandleFunc("", todoAPI.CreateTodo).Methods("POST")
	todos.HandleFunc("/{id}", todoAPI.UpdateTodo).Methods("PUT")
	todos.HandleFunc("/{id}", todoAPI.DeleteTodo).Methods("DELETE")
	todos.HandleFunc("", todoAPI.ListTodos).Methods("GET")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})
}
