package repository

import (
	"context"
	"fmt"
	"otus/go-server-project/internal/models"
)

func (r *Repo) SendMessage(ctx context.Context, fromUserID, toUserID, message string) error {
	conn, err := r.Master.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	_, err = conn.Exec(ctx,
		"INSERT INTO dialogs (from_user_id, to_user_id, message) VALUES ($1, $2, $3)",
		fromUserID, toUserID, message,
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (r *Repo) ListDialogs(ctx context.Context, user1 string, user2 string) ([]models.Dialog, error) {
	dialogs := []models.Dialog{}
	conn, err := r.Master.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to acquire connection: %w", err)
	}
	defer conn.Release()
	rows, err := conn.Query(context.Background(),
		`SELECT from_user_id, to_user_id, message
		FROM dialogs
		WHERE (from_user_id = $1 AND to_user_id = $2) OR (from_user_id = $2 AND to_user_id = $1)
		ORDER BY created_at ASC`,
		user1, user2,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list dialogs: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var dialog models.Dialog
		err := rows.Scan(&dialog.From, &dialog.To, &dialog.Text)
		if err != nil {
			return nil, fmt.Errorf("failed to scan dialog: %w", err)
		}
		dialogs = append(dialogs, dialog)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return dialogs, nil
}
