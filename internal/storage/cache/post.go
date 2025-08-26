package cache

import (
	"context"
	"errors"
	"otus/go-server-project/internal/models"
	"strconv"
)

func createPostFeedKey(limit, offset int) string {
	return "limit" + strconv.Itoa(limit) + "offset" + strconv.Itoa(offset)
}

func (c *LRUCache) GetFeed(_ context.Context, limit, offset int) ([]models.Post, error) {
	var resp []models.Post
	key := createPostFeedKey(limit, offset)
	data, ok := c.get(key)
	if !ok {
		return resp, errors.New("unable to get data from cache")
	}

	byteData, ok := data.([]models.Post)
	if !ok {
		return resp, errors.New("cached data is not of type []byte")
	}

	return byteData, nil
}

func (c *LRUCache) SetFeed(ctx context.Context, limit, offset int, models []models.Post) {
	key := createPostFeedKey(limit, offset)
	c.set(key, models)
}
