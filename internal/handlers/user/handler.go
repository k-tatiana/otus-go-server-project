package user

import (
	"context"
	"otus/go-server-project/internal/models"

	"go.uber.org/zap"
)

type UserService interface {
	Login(context.Context, string, string) (string, error)
	RegisterUser(context.Context, models.User) (string, error)
	Get(context.Context, string) (models.User, error)
	SearchUser(context.Context, string, string) ([]models.User, error)
}

type authenticator interface {
	Auth(context.Context, string) error
}

type UserHandler struct {
	service       UserService
	logger        *zap.Logger
	authenticator authenticator
}

func NewUserHandler(s UserService, logger *zap.Logger, authenticator authenticator) *UserHandler {
	return &UserHandler{service: s, logger: logger, authenticator: authenticator}
}
