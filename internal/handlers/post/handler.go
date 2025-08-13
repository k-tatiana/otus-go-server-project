package post

import (
	"context"
	"otus/go-server-project/internal/models"

	"go.uber.org/zap"
)

type PostsService interface {
	Feed(context.Context, int, int) ([]models.Post, error)
}

type authenticator interface {
	Auth(context.Context, string) error
}

type PostsHandler struct {
	service       PostsService
	logger        *zap.Logger
	authenticator authenticator
}

func NewPostsHandler(s PostsService, logger *zap.Logger, authenticator authenticator) *PostsHandler {
	return &PostsHandler{service: s, logger: logger, authenticator: authenticator}
}
