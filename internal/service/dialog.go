package service

import (
	"context"
	"time"

	"otus/go-server-project/internal/models"
)

type DialogRepository interface {
	SendMessage(context.Context, string, string, string, *time.Time) error
	ListDialogs(context.Context, string, string) ([]models.Dialog, error)
}

type DialogService struct {
	repo DialogRepository
}

func NewDialogService(r DialogRepository) *DialogService {
	return &DialogService{repo: r}
}

func (s *DialogService) Send(ctx context.Context, fromUserID, toUserID, message string, createdAt *time.Time) error {
	err := s.repo.SendMessage(ctx, fromUserID, toUserID, message, createdAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *DialogService) List(ctx context.Context, userID1, userID2 string) ([]models.Dialog, error) {
	dialogs := []models.Dialog{}
	dialogs, err := s.repo.ListDialogs(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}

	return dialogs, nil
}
