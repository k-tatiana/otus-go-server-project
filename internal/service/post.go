package service

import (
	"context"
	"otus/go-server-project/internal/models"
)

type PostsRepository interface {
	GetFeed(ctx context.Context, offset, limit int) ([]models.Post, error)
}

type postsService struct {
	repo PostsRepository
}

func NewPostsService(r PostsRepository) *postsService {
	return &postsService{repo: r}
}

func (s *postsService) Feed(ctx context.Context, limit, offset int) ([]models.Post, error) {
	models, err := s.repo.GetFeed(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return models, nil
}
