package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/isuraem/todo-api/internal/adapters/app/websocket"
	"github.com/isuraem/todo-api/internal/middleware"
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
	userID := r.Context().Value(middleware.UserIDKey).(uint)

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo.UserID = userID

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
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	todos, err := api.service.List(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todos)
}

func (api *TodoAPI) LikeTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.service.LikeTodoByUser(uint(todoID), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.broadcastTodos()
	w.WriteHeader(http.StatusOK)
}

func (api *TodoAPI) UnlikeTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(uint)
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := api.service.UnlikeTodoByUser(uint(todoID), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	api.broadcastTodos()
	w.WriteHeader(http.StatusOK)
}

func (api *TodoAPI) broadcastTodos() {
	todos, err := api.service.List(0) // List without user context for broadcast
	if err != nil {
		return
	}
	api.hub.BroadcastTodos(todos)
}
