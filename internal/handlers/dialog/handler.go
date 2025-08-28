package dialog

import (
	"context"
	"otus/go-server-project/internal/models"

	"go.uber.org/zap"
)

type DialogService interface {
	Send(context.Context, string, string, string) error
	List(context.Context, string, string) ([]models.Dialog, error)
}

type authenticator interface {
	Auth(context.Context, string) error
}

type DialogHandler struct {
	service       DialogService
	logger        *zap.Logger
	authenticator authenticator
}

func NewDialogHandler(
	s DialogService,
	logger *zap.Logger,
	authenticator authenticator,
) *DialogHandler {
	return &DialogHandler{service: s, logger: logger, authenticator: authenticator}
}
