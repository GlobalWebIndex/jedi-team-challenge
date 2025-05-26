package repositories

import (
	"github.com/google/uuid"
	"time"
)

type ChatSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	Title     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Messages  []Message `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
}

type Message struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SessionID uuid.UUID `gorm:"type:uuid;not null;index"`
	Sender    string    `gorm:"type:text;not null"`
	Content   string    `gorm:"type:text;not null"`
	Timestamp time.Time `gorm:"not null"`
}
