package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/isuraem/todo-api/internal/adapters/app/websocket"
	"github.com/isuraem/todo-api/internal/models"
	"github.com/isuraem/todo-api/internal/ports"
)

type TodoAPI struct {
	service ports.TodoService
	hub     *websocket.Hub
}

func NewTodoAPI(service ports.TodoService, hub *websocket.Hub) *TodoAPI {
	return &TodoAPI{service: service, hub: hub}
}

func (api *TodoAPI) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// Extract JWT from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Extract sub from token claims and assign to UserID
	sub, ok := claims["sub"].(string)
	if !ok {
		http.Error(w, "Invalid sub claim", http.StatusUnauthorized)
		return
	}

	userId, err := strconv.ParseUint(sub, 10, 32)
	if err != nil {
		http.Error(w, "Invalid sub claim", http.StatusUnauthorized)
		return
	}

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assign the extracted userId to the Todo's UserID
	todo.UserID = uint(userId)

	if err := api.service.Create(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.broadcastTodos()
	w.WriteHeader(http.StatusCreated)
}
func (api *TodoAPI) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := api.service.Update(uint(id), todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.broadcastTodos()
	w.WriteHeader(http.StatusOK)
}

func (api *TodoAPI) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := api.service.Delete(uint(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.broadcastTodos()
	w.WriteHeader(http.StatusOK)
}

func (api *TodoAPI) ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := api.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

func (api *TodoAPI) broadcastTodos() {
	todos, err := api.service.List()
	if err != nil {
		return
	}
	api.hub.BroadcastTodos(todos)
}
