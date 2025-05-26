package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	apierrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"gorm.io/gorm"
)

func (repo *ChatSessionRepository) GetChatSession(
	ctx context.Context,
	uuid uuid.UUID,
) (*domain.ChatSession, error) {
	var err error
	var modelChatSession *ChatSession

	err = repo.db.WithContext(ctx).
		Preload("Messages").
		Model(ChatSession{}).
		Where("id = ?", uuid).
		Take(&modelChatSession).Error

	if err == gorm.ErrRecordNotFound {
		return &domain.ChatSession{}, apierrors.ResourceNotFoundErrorWrapper{
			OriginalError: errors.New("chatSessionID " + uuid.String() + " not found"),
		}
	}

	if err != nil {
		return &domain.ChatSession{}, err
	}

	messages := make([]*domain.Message, len(modelChatSession.Messages))
	for i, msg := range modelChatSession.Messages {
		messages[i] = &domain.Message{
			ID:            msg.ID,
			ChatSessionID: msg.ChatSessionID,
			Sender:        msg.Sender,
			Content:       msg.Content,
			Timestamp:     msg.Timestamp,
		}
	}

	return &domain.ChatSession{
		ID:        modelChatSession.ID,
		Title:     modelChatSession.Title,
		UserID:    modelChatSession.UserID,
		CreatedAt: modelChatSession.CreatedAt,
		UpdatedAt: modelChatSession.UpdatedAt,
		Messages:  messages,
	}, err
}

func (repo *ChatSessionRepository) GetUserChatSessions(
	ctx context.Context,
	uuid uuid.UUID,
) ([]*domain.ChatSession, error) {
	var err error
	var modelChatSessions []*ChatSession

	err = repo.db.WithContext(ctx).
		Preload("Messages").
		Model(ChatSession{}).
		Where("user_id = ?", uuid).
		Find(&modelChatSessions).Error

	if err == gorm.ErrRecordNotFound {
		return []*domain.ChatSession{}, apierrors.ResourceNotFoundErrorWrapper{
			OriginalError: errors.New("user uuid " + uuid.String() + " not found"),
		}
	}

	if err != nil {
		return []*domain.ChatSession{}, err
	}

	chatSessions := make([]*domain.ChatSession, 0, len(modelChatSessions))
	for _, modelChatSession := range modelChatSessions {

		messages := make([]*domain.Message, len(modelChatSession.Messages))
		for i, msg := range modelChatSession.Messages {
			messages[i] = &domain.Message{
				ID:            msg.ID,
				ChatSessionID: msg.ChatSessionID,
				Sender:        msg.Sender,
				Content:       msg.Content,
				Timestamp:     msg.Timestamp,
			}
		}

		chatSessions = append(
			chatSessions,
			&domain.ChatSession{
				ID:        modelChatSession.ID,
				Title:     modelChatSession.Title,
				UserID:    modelChatSession.UserID,
				CreatedAt: modelChatSession.CreatedAt,
				UpdatedAt: modelChatSession.UpdatedAt,
				Messages:  messages,
			},
		)
	}

	return chatSessions, nil
}
