package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	customerrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"gorm.io/gorm"
)

func (repo *ChatSessionRepository) UpdateChatSessionTitle(
	ctx context.Context,
	uuid uuid.UUID,
	title string,
) error {
	err := repo.db.WithContext(ctx).Model(&ChatSession{}).
		Where("id = ?", uuid).
		Update("title", title).Error

	if err == gorm.ErrRecordNotFound {
		return customerrors.ResourceNotFoundErrorWrapper{
			OriginalError: errors.New("chatSessionID " + uuid.String() + " not found"),
		}
	}

	return err
}
