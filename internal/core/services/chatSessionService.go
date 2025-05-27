package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"github.com/loukaspe/jedi-team-challenge/internal/core/ports"
	"github.com/loukaspe/jedi-team-challenge/pkg/logger"
)

type ChatSessionServiceInterface interface {
	GetChatSession(context.Context, uuid.UUID) (*domain.ChatSession, error)
	GetUserChatSessions(context.Context, uuid.UUID) ([]*domain.ChatSession, error)
	CreateChatSession(context.Context, *domain.ChatSession) (uuid.UUID, error)
	UpdateChatSessionTitle(context.Context, uuid.UUID, string) error
}

type ChatSessionService struct {
	logger     logger.LoggerInterface
	repository ports.ChatSessionRepositoryInterface
}

func NewChatSessionService(
	logger logger.LoggerInterface,
	repository ports.ChatSessionRepositoryInterface,
) *ChatSessionService {
	return &ChatSessionService{
		logger:     logger,
		repository: repository,
	}
}

func (s ChatSessionService) CreateChatSession(ctx context.Context, session *domain.ChatSession) (uuid.UUID, error) {
	return s.repository.CreateChatSession(ctx, session)
}

func (s ChatSessionService) UpdateChatSessionTitle(ctx context.Context, uuid uuid.UUID, title string) error {
	return s.repository.UpdateChatSessionTitle(ctx, uuid, title)
}

func (s ChatSessionService) GetChatSession(ctx context.Context, uuid uuid.UUID) (*domain.ChatSession, error) {
	return s.repository.GetChatSession(ctx, uuid)
}

func (s ChatSessionService) GetUserChatSessions(ctx context.Context, uuid uuid.UUID) ([]*domain.ChatSession, error) {
	return s.repository.GetUserChatSessions(ctx, uuid)
}
