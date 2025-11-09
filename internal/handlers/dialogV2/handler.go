package dialogv2

import (
	"context"
	"time"

	"go.uber.org/zap"

	"otus/go-server-project/internal/models"
)

type DialogService interface {
	Send(context.Context, string, string, string, *time.Time) error
	List(context.Context, string, string) ([]models.Dialog, error)
}

type authenticator interface {
	Auth(context.Context, string) error
}

type DialogAPIHandler struct {
	service       DialogService
	logger        *zap.Logger
	authenticator authenticator
}

func NewDialogAPIHandler(
	s DialogService,
	logger *zap.Logger,
	authenticator authenticator,
) *DialogAPIHandler {
	return &DialogAPIHandler{service: s, logger: logger, authenticator: authenticator}
}
