package user

import (
	"errors"

	"github.com/isuraem/todo-api/internal/models"
	"github.com/isuraem/todo-api/internal/ports"
)

type Service struct {
	userDB     ports.UserDB
	jwtService ports.JWTService
}

func NewUserService(userDB ports.UserDB, jwtService ports.JWTService) *Service {
	return &Service{
		userDB:     userDB,
		jwtService: jwtService,
	}
}

func (s *Service) Register(user models.User) error {
	if err := s.userDB.CreateUser(user); err != nil {
		return err
	}
	return nil
}

func (s *Service) Login(email, password string) (string, error) {
	user, err := s.userDB.GetUserByEmail(email)
	if err != nil || user.Password != password {
		return "", errors.New("invalid credentials")
	}
	token, err := s.jwtService.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}
