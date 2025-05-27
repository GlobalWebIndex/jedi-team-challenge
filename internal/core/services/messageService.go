package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"github.com/loukaspe/jedi-team-challenge/internal/core/ports"
	apierrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"github.com/loukaspe/jedi-team-challenge/pkg/logger"
)

type MessageServiceInterface interface {
	CreateMessage(context.Context, uuid.UUID, *domain.Message) (uuid.UUID, error)
	GetAnswerForMessage(context.Context, uuid.UUID) (*domain.Message, error)
}

type MessageService struct {
	logger                logger.LoggerInterface
	messageRepository     ports.MessageRepositoryInterface
	chatSessionRepository ports.ChatSessionRepositoryInterface
}

func NewMessageService(
	logger logger.LoggerInterface,
	messageRepositoryInterface ports.MessageRepositoryInterface,
	chatSessionRepository ports.ChatSessionRepositoryInterface,
) *MessageService {
	return &MessageService{
		logger:                logger,
		messageRepository:     messageRepositoryInterface,
		chatSessionRepository: chatSessionRepository,
	}
}

func (s MessageService) CreateMessage(ctx context.Context, userID uuid.UUID, message *domain.Message) (uuid.UUID, error) {
	chatSession, err := s.chatSessionRepository.GetChatSession(ctx, message.ChatSessionID)
	if err != nil {
		return uuid.Nil, err
	}

	if chatSession.UserID != userID {
		return uuid.Nil, apierrors.NewUserMismatchError(message.ChatSessionID.String(), userID.String())
	}

	return s.messageRepository.CreateMessage(ctx, message)
}

func (s MessageService) GetAnswerForMessage(ctx context.Context, initialMessageID uuid.UUID) (*domain.Message, error) {
	initialMessage, err := s.messageRepository.GetMessage(ctx, initialMessageID)
	if err != nil {
		return nil, err
	}

	replyMessage := &domain.Message{
		ID:            initialMessageID,
		ChatSessionID: initialMessage.ChatSessionID,
		Content:       "This is a mock response from the LLM",
		Sender:        "SYSTEM", //TODO: constant
	}

	//TODO: get from llm
	insertedMessageID, err := s.messageRepository.CreateMessage(ctx, replyMessage)
	if err != nil {
		return nil, err
	}

	replyMessage.ID = insertedMessageID

	return replyMessage, nil
}
