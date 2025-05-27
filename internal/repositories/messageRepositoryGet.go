package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	apierrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"gorm.io/gorm"
)

func (repo *MessageRepository) GetMessage(
	ctx context.Context,
	uuid uuid.UUID,
) (*domain.Message, error) {
	var err error
	var modelMessage *Message

	err = repo.db.WithContext(ctx).
		Model(Message{}).
		Where("id = ?", uuid).
		Take(&modelMessage).Error

	if err == gorm.ErrRecordNotFound {
		return &domain.Message{}, apierrors.ResourceNotFoundErrorWrapper{
			OriginalError: errors.New("messageID " + uuid.String() + " not found"),
		}
	}

	if err != nil {
		return &domain.Message{}, err
	}

	return &domain.Message{
		ID:            modelMessage.ID,
		ChatSessionID: modelMessage.ChatSessionID,
		Sender:        modelMessage.Sender,
		Content:       modelMessage.Content,
		CreatedAt:     modelMessage.CreatedAt,
	}, err
}
