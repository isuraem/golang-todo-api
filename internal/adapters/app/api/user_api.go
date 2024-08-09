package api

import (
	"net/http"

	"github.com/isuraem/todo-api/internal/models"
	"github.com/isuraem/todo-api/internal/ports"
)

type UserAPI struct {
	service ports.UserService
}

func NewUserAPI(service ports.UserService) *UserAPI {
	return &UserAPI{service: service}
}

func (api *UserAPI) Register(w http.ResponseWriter, r *http.Request) {
	// Get the validated user from the context
	user := r.Context().Value("validatedUser").(models.User)

	if err := api.service.Register(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (api *UserAPI) Login(w http.ResponseWriter, r *http.Request) {
	// Get the validated credentials from the context
	creds := r.Context().Value("validatedUser").(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})

	token, err := api.service.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Write([]byte(token))
}
