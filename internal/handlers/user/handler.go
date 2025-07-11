package user

import "otus/go-server-project/internal/models"

type UserService interface {
	Login(login, password string) (string, error)
	RegisterUser(u models.User) (string, error)
	Get(id int) (models.User, error)
}

type UserHandler struct {
	service UserService
}

func NewUserHandler(s UserService) *UserHandler {
	return &UserHandler{service: s}
}
