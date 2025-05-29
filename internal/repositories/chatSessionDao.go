package repositories

import (
	"github.com/google/uuid"
	"time"
)

const SYSTEM_SENDER = "SYSTEM"
const USER_SENDER = "USER"

type ChatSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Messages  []Message
}

type Message struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ChatSessionID uuid.UUID `gorm:"type:uuid;not null;index"`
	Sender        string    `gorm:"type:text;not null"`
	Content       string    `gorm:"type:text;not null"`
	CreatedAt     time.Time `gorm:"not null"`
	Feedback      *string   `gorm:"type:text;null"`
}
