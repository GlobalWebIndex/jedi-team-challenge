package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"github.com/loukaspe/jedi-team-challenge/internal/core/ports"
	apierrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"github.com/loukaspe/jedi-team-challenge/pkg/logger"
	"github.com/openai/openai-go"
	"strings"
)

type MessageServiceInterface interface {
	CreateMessage(context.Context, uuid.UUID, *domain.Message) (uuid.UUID, error)
	GetAnswerForMessage(context.Context, uuid.UUID) (*domain.Message, error)
}

type Embedder interface {
	Embed(ctx context.Context, inputs []string) ([][]float64, error)
}

type VectorDB interface {
	SemanticSearch(ctx context.Context, embeddings []float32) ([]string, error)
}

type MessageService struct {
	logger                logger.LoggerInterface
	messageRepository     ports.MessageRepositoryInterface
	chatSessionRepository ports.ChatSessionRepositoryInterface
	embedder              Embedder
	vectorDB              VectorDB
	openAIClient          *openai.Client
}

func NewMessageService(
	logger logger.LoggerInterface,
	messageRepositoryInterface ports.MessageRepositoryInterface,
	chatSessionRepository ports.ChatSessionRepositoryInterface,
	embedder Embedder,
	vectorDB VectorDB,
	openAIClient *openai.Client,
) *MessageService {
	return &MessageService{
		logger:                logger,
		messageRepository:     messageRepositoryInterface,
		chatSessionRepository: chatSessionRepository,
		embedder:              embedder,
		vectorDB:              vectorDB,
		openAIClient:          openAIClient,
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

	embeddings, err := s.embedder.Embed(context.Background(), []string{
		"What do you know about Gen Z in Nashville",
	})
	if err != nil {
		return nil, err
	}

	// we only have on text so we only care for the first embedding row
	vectorToFloat32 := make([]float32, len(embeddings[0]))
	for i, v := range embeddings[0] {
		vectorToFloat32[i] = float32(v)
	}

	accumulatedTextFromSearch, err := s.vectorDB.SemanticSearch(context.Background(), vectorToFloat32)

	prompt := fmt.Sprintf(`Use the following context to answer the question.
		Context:
		%s
		
		Question:
		%s
		
		Answer:`,
		strings.Join(accumulatedTextFromSearch, "\n"),
		initialMessage.Content,
	)

	chatCompletion, err := s.openAIClient.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4o,
	})
	if err != nil {
		return nil, err
	}

	answer := chatCompletion.Choices[0].Message.Content

	if answer == "" {
		// TODO: mh ta les apierrors
		return nil, errors.New("received empty response from LLM")
	}

	replyMessage := &domain.Message{
		ID:            initialMessageID,
		ChatSessionID: initialMessage.ChatSessionID,
		Content:       answer,
		Sender:        "SYSTEM", //TODO: constant
	}

	insertedMessageID, err := s.messageRepository.CreateMessage(ctx, replyMessage)
	if err != nil {
		return nil, err
	}

	replyMessage.ID = insertedMessageID

	return replyMessage, nil
}
