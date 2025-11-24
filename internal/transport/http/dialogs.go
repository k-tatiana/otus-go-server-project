package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"otus/go-server-project/internal/models"
)

type SendResponse struct {
	DialogID string `json:"dialog_id"`
}

type SendDialogRequest struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Message    string `json:"text"`
}

type ListDialogRequest struct {
	FirstUserID  string `json:"first_user_id"`
	SecondUserID string `json:"second_user_id"`
}

func NewDialogClient(dialogAddress string) *DialogClient {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &DialogClient{client: client, address: dialogAddress}
}

type DialogClient struct {
	client  *http.Client
	address string
}

func (dc *DialogClient) Send(ctx context.Context, fromUserID, toUserID, message string, time *time.Time) (string, error) {
	// Реализация отправки диалога через HTTP
	path := dc.address + "/api/send"
	body := &SendDialogRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Message:    message,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	var response SendResponse
	resp, err := dc.client.Post(path, "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to send dialog: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return response.DialogID, nil
}

func (dc *DialogClient) List(ctx context.Context, firstUserID, secondUserID string) ([]models.Dialog, error) {
	// Реализация получения списка диалогов через HTTP
	path := dc.address + "/api/list"
	body := &ListDialogRequest{
		FirstUserID:  firstUserID,
		SecondUserID: secondUserID,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := dc.client.Post(path, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to list dialogs: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var dialogs []models.Dialog
	if err := json.NewDecoder(resp.Body).Decode(&dialogs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return dialogs, nil
}

func (dc *DialogClient) Delete(ctx context.Context, dialogID string) error {
	// Реализация удаления диалога через HTTP
	data, err := json.Marshal(map[string]string{
		"dialog_id": dialogID,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal delete request body: %v", err)
	}
	path := dc.address + "/api/delete"
	resp, err := dc.client.Post(path, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to delete dialog: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
