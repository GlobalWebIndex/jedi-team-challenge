package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (repo *MessageRepository) CreateMessage(
	ctx context.Context,
	chat *domain.Message,
) (uuid.UUID, error) {
	var err error

	modelChat := Message{
		ChatSessionID: chat.ChatSessionID,
		Sender:        chat.Sender,
		Content:       chat.Content,
		CreatedAt:     chat.CreatedAt,
	}

	err = repo.db.WithContext(ctx).Create(&modelChat).Error
	if err != nil {
		return uuid.Nil, err
	}

	return modelChat.ID, nil
}
