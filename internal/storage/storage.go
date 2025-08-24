package storage

import (
	"context"
	"otus/go-server-project/internal/models"
)

type storage interface {
	GetFeed(ctx context.Context, limit, offset int) ([]models.Post, error)
	SetFeed(ctx context.Context, limit, offset int, models []models.Post)
}

type Storage struct {
	repo     storage
	cache    storage
	useCache bool
}

// GetFeed implements the PostsRepository interface.
// Adjust the signature and implementation as required by your application.
func (s *Storage) GetFeed(ctx context.Context, limit, offset int) ([]models.Post, error) {
	if !s.useCache {
		return s.repo.GetFeed(ctx, limit, offset)
	}

	models, err := s.cache.GetFeed(ctx, limit, offset)
	if err == nil {
		return models, nil
	}
	models, err = s.repo.GetFeed(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	s.cache.SetFeed(ctx, limit, offset, models)
	// Optionally, you might want to cache the result here.
	return models, nil
}

func NewStorage(r storage, c storage, useCache bool) *Storage {
	return &Storage{
		repo: r, cache: c, useCache: useCache,
	}
}
