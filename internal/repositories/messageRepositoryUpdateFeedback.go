package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	customerrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"gorm.io/gorm"
)

func (repo *MessageRepository) UpdateMessageFeedback(
	ctx context.Context,
	uuid uuid.UUID,
	feedback string,
) error {
	var err error

	err = repo.db.WithContext(ctx).
		Model(Message{}).
		Where("id = ?", uuid).
		Update("feedback", feedback).Error

	if err == gorm.ErrRecordNotFound {
		// TODO: customerrors
		return customerrors.ResourceNotFoundErrorWrapper{
			OriginalError: errors.New("messageID " + uuid.String() + " not found"),
		}
	}

	return err
}
