package user

import (
	"otus/go-server-project/internal/models"

	"go.uber.org/zap"
)

type UserService interface {
	ValidateToken(token string) error
	Login(login, password string) (string, error)
	RegisterUser(u models.User) (string, error)
	Get(id string) (models.User, error)
	SearchUser(name, surname string) ([]models.User, error)
}

type UserHandler struct {
	service UserService
	logger  *zap.Logger
}

func NewUserHandler(s UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: s, logger: logger}
}
