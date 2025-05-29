package domain

import (
	"github.com/google/uuid"
	"time"
)

type ChatSession struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []*Message
}

type Message struct {
	ID            uuid.UUID
	ChatSessionID uuid.UUID
	Sender        string
	Content       string
	CreatedAt     time.Time
	Feedback      *string
}
