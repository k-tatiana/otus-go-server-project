package service

import (
	"context"
	"otus/go-server-project/internal/models"
)

type DialogRepository interface {
	SendMessage(context.Context, string, string, string) error
	ListDialogs(context.Context, string, string) ([]models.Dialog, error)
}

type dialogService struct {
	repo DialogRepository
}

func NewDialogService(r DialogRepository) *dialogService {
	return &dialogService{repo: r}
}

func (s *dialogService) Send(ctx context.Context, fromUserID, toUserID, message string) error {
	err := s.repo.SendMessage(ctx, fromUserID, toUserID, message)
	if err != nil {
		return err
	}
	return nil
}

func (s *dialogService) List(ctx context.Context, userID1, userID2 string) ([]models.Dialog, error) {
	dialogs := []models.Dialog{}
	dialogs, err := s.repo.ListDialogs(ctx, userID1, userID2)
	if err != nil {
		return nil, err
	}

	return dialogs, nil
}
