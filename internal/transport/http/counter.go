package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"otus/go-server-project/internal/models"
	"time"
)

type IncrementCounterRequest struct {
	FromUserID string `json:"from_user_id"`
	ToUserID   string `json:"to_user_id"`
	Count      int    `json:"count,omitempty"`
}

type CounterResponse struct {
	ReadCount int `json:"read_count"`
}

func NewCounterClient(counterAddress string) *CounterClient {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &CounterClient{client: client, address: counterAddress}
}

type CounterClient struct {
	client  *http.Client
	address string
}

func (cc *CounterClient) Increment(ctx context.Context, from_user_id, to_user_id string) error {
	// Реализация инкремента счетчика через HTTP
	path := cc.address + "/increment"
	body := &IncrementCounterRequest{
		FromUserID: from_user_id,
		ToUserID:   to_user_id,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := cc.client.Post(path, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to increment counter: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (cc *CounterClient) Decrement(ctx context.Context, from_user_id, to_user_id string) error {
	// Реализация инкремента счетчика через HTTP
	path := cc.address + "/decrement"
	body := &IncrementCounterRequest{
		FromUserID: from_user_id,
		ToUserID:   to_user_id,
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := cc.client.Post(path, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to increment counter: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 202 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func (cc *CounterClient) GetCount(ctx context.Context, from_user_id, to_user_id string) (models.GetCountResponse, error) {
	var getCountResp models.GetCountResponse

	path := fmt.Sprintf("%s/get?from_user_id=%s&to_user_id=%s", cc.address, from_user_id, to_user_id)

	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return getCountResp, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := cc.client.Do(req)
	if err != nil {
		return getCountResp, fmt.Errorf("failed to get counter: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return getCountResp, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&getCountResp)
	if err != nil {
		return getCountResp, fmt.Errorf("failed to decode response: %v", err)
	}

	return getCountResp, nil
}
