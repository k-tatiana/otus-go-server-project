package user

import (
	"context"
	"otus/go-server-project/internal/models"

	"go.uber.org/zap"
)

type UserService interface {
	ValidateToken(context.Context, string) error
	Login(context.Context, string, string) (string, error)
	RegisterUser(context.Context, models.User) (string, error)
	Get(context.Context, string) (models.User, error)
	SearchUser(context.Context, string, string) ([]models.User, error)
}

type UserHandler struct {
	service UserService
	logger  *zap.Logger
}

func NewUserHandler(s UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{service: s, logger: logger}
}
