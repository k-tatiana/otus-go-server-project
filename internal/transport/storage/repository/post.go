package repository

import (
	"context"
	"fmt"

	"otus/go-server-project/internal/models"
)

func (r *Repo) GetFeed(ctx context.Context, limit, offset int) ([]models.Post, error) {
	conn, err := r.Replicas.Acquire(ctx) // balancing through haproxy
	if err != nil {
		return []models.Post{}, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	rows, err := conn.Query(ctx,
		"SELECT id, author_user_id, text FROM posts ORDER BY id DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get feed: %w", err)
	}
	defer rows.Close()

	posts := make([]models.Post, 0, limit)
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.AuthorUserID, &post.Text); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return posts, nil
}

func (r *Repo) SetFeed(ctx context.Context, limit, offset int, models []models.Post) {
	// This method is intentionally left blank as caching logic is handled elsewhere.
}

func (r *Repo) WritePost(ctx context.Context, userId, post string) (string, error) {
	conn, err := r.Master.Acquire(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	var postId string
	err = conn.QueryRow(ctx,
		`INSERT INTO posts (author_user_id, text) VALUES ($1::varchar, $2) RETURNING id`,
		userId, post,
	).Scan(&postId)
	if err != nil {
		return "", fmt.Errorf("failed to write post: %w", err)
	}

	return postId, nil
}
