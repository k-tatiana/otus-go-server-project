package user

import "otus/go-server-project/internal/models"

type UserService interface {
	ValidateToken(token string) error
	Login(login, password string) (string, error)
	RegisterUser(u models.User) (string, error)
	Get(id string) (models.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{service: s}
}
