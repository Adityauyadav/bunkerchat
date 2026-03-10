package models

import (
	"time"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	CreatedAt    time.Time
}

type Message struct {
	ID         int
	SentFromID int
	SentToID   int
	Content    string
	Read       bool
	CreatedAt  time.Time
}
