package models

import (
	"time"
)

type UserDTO struct {
	ID           *string
	Name         string
	Surname      string
	Birthday     time.Time
	Gender       string
	Interests    []string
	City         string
	Login        string
	PasswordHash string
	Token        string
}
