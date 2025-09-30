package models

import "time"

type Dialog struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

type DialogDTO struct {
	FromUserID   string     `json:"from_user_id"`
	ToUserID     string     `json:"to_user_id"`
	Message      string     `json:"message"`
	CreatedAt    *time.Time `json:"created_at"`
	UserPairHash *string    `json:"user_pair_hash"`
}
