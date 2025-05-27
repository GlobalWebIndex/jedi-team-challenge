package ports

import (
	"context"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
)

type MessageRepositoryInterface interface {
	CreateMessage(context.Context, *domain.Message) (uuid.UUID, error)
	GetMessage(context.Context, uuid.UUID) (*domain.Message, error)
}
