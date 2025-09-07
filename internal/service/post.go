package service

import (
	"context"
	"otus/go-server-project/internal/models"
)

type PostsRepository interface {
	GetFeed(ctx context.Context, limit, offset int) ([]models.Post, error)
	WritePost(ctx context.Context, userId string, post string) (string, error)
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

func (s *postsService) WritePost(ctx context.Context, userId, post string) (string, error) {
	return s.repo.WritePost(ctx, userId, post)
}
