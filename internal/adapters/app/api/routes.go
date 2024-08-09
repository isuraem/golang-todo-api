package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/isuraem/todo-api/internal/adapters/app/websocket"
	"github.com/isuraem/todo-api/internal/middleware"
)

// SetupRoutes sets up the API routes and handlers.
func SetupRoutes(r *mux.Router, userAPI *UserAPI, todoAPI *TodoAPI, hub *websocket.Hub) {
	r.Handle("/login", middleware.ValidateUser(http.HandlerFunc(userAPI.Login))).Methods("POST")
	r.Handle("/register", middleware.ValidateUser(http.HandlerFunc(userAPI.Register))).Methods("POST")
	r.Handle("/todos", middleware.ValidateTodo(http.HandlerFunc(todoAPI.CreateTodo))).Methods("POST")
	r.Handle("/todos/{id}", middleware.ValidateTodo(http.HandlerFunc(todoAPI.UpdateTodo))).Methods("PUT")
	r.Handle("/todos/{id}", http.HandlerFunc(todoAPI.DeleteTodo)).Methods("DELETE")
	r.Handle("/todos", http.HandlerFunc(todoAPI.ListTodos)).Methods("GET")
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})
}
