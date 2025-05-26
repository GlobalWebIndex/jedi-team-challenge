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
	ID        uuid.UUID
	SessionID uuid.UUID
	Sender    string
	Content   string
	Timestamp time.Time
}
