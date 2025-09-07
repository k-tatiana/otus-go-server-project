package models

type WebSocketMessage struct {
	Type MessageType `json:"type"`

	PostID       string  `json:"post_id"`
	PostText     string  `json:"post_text"`
	AuthorUserID *string `json:"author_user_id"`
}
