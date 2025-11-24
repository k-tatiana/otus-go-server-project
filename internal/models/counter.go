package models

type GetCountResponse struct {
	CountTotal  int `json:"count_total"`
	CountUnread int `json:"count_unread"`
}
