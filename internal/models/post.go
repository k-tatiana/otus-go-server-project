package models

type Post struct {
	ID           int64   `json:"id"`
	AuthorUserID *string `json:"author_user_id"`
	Text         *string `json:"text"`
}
