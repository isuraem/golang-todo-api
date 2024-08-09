package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/isuraem/todo-api/internal/adapters/app/websocket"
)

// SetupRoutes sets up the API routes and handlers.
func SetupRoutes(r *mux.Router, userAPI *UserAPI, todoAPI *TodoAPI, hub *websocket.Hub) {
	r.HandleFunc("/login", userAPI.Login).Methods("POST")
	r.HandleFunc("/register", userAPI.Register).Methods("POST")
	r.HandleFunc("/todos", todoAPI.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", todoAPI.UpdateTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", todoAPI.DeleteTodo).Methods("DELETE")
	r.HandleFunc("/todos", todoAPI.ListTodos).Methods("GET")
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})
}
