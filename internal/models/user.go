package models

import (
	"time"
)

type User struct {
	ID        *string   `json:"id,omitempty"`
	Name      string    `json:"name"`
	Surname   string    `json:"surname"`
	Birthday  time.Time `json:"birthday"`
	Gender    string    `json:"gender"`
	Interests []string  `json:"interests"`
	City      string    `json:"city"`
	Login     string    `json:"login"`
	Password  string    `json:"password,omitempty"`
}
