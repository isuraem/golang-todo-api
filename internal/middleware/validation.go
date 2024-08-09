package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/isuraem/todo-api/internal/models"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := decodeJSONBody(w, r, &user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validate.Struct(user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "validatedUser", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateTodo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var todo models.Todo
		if err := decodeJSONBody(w, r, &todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := validate.Struct(todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), "validatedTodo", todo)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}
