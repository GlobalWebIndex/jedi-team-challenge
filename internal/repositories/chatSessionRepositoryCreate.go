package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"gorm.io/gorm"
)

type ChatSessionRepository struct {
	db *gorm.DB
}

func NewChatSessionRepository(db *gorm.DB) *ChatSessionRepository {
	return &ChatSessionRepository{db: db}
}

func (repo *ChatSessionRepository) CreateChatSession(
	ctx context.Context,
	chat *domain.ChatSession,
) (uuid.UUID, error) {
	var err error

	modelChat := ChatSession{
		UserID:   chat.UserID,
		Title:    chat.Title,
		Messages: nil,
	}

	err = repo.db.WithContext(ctx).Create(&modelChat).Error
	if err != nil {
		return uuid.UUID{}, err
	}

	return modelChat.ID, nil
}
