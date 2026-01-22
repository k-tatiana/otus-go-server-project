package dialogv2

import (
	"context"
	"time"

	"go.uber.org/zap"

	"otus/go-server-project/internal/models"
)

type DialogService interface {
	Send(context.Context, string, string, string, *time.Time) (string, error)
	List(context.Context, string, string) ([]models.Dialog, error)
	Delete(context.Context, string) error
}

type authenticator interface {
	Auth(context.Context, string) error
}

type CounterService interface {
	Increment(ctx context.Context, from_user_id, to_user_id string) error
	Decrement(ctx context.Context, from_user_id, to_user_id string) error
	GetCount(ctx context.Context, from_user_id, to_user_id string) (models.GetCountResponse, error)
}

type DialogAPIHandler struct {
	service       DialogService
	logger        *zap.Logger
	authenticator authenticator
	counter       CounterService
}

func NewDialogAPIHandler(
	s DialogService,
	logger *zap.Logger,
	authenticator authenticator,
	counter CounterService,
) *DialogAPIHandler {
	return &DialogAPIHandler{
		service: s, logger: logger, authenticator: authenticator, counter: counter,
	}
}
