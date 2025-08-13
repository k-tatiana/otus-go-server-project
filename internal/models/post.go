package models

type Post struct {
	ID           int64  `json:"id"`
	AuthorUserID int64  `json:"author_user_id"`
	Text         string `json:"text"`
}
