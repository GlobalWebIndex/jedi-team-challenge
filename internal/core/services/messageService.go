package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"github.com/loukaspe/jedi-team-challenge/internal/core/ports"
	"github.com/loukaspe/jedi-team-challenge/internal/repositories"
	apierrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"github.com/loukaspe/jedi-team-challenge/pkg/helpers"
	"github.com/loukaspe/jedi-team-challenge/pkg/logger"
	"github.com/openai/openai-go"
	"strings"
)

type MessageServiceInterface interface {
	CreateMessage(context.Context, uuid.UUID, *domain.Message) (uuid.UUID, error)
	GetAnswerForMessage(context.Context, uuid.UUID) (*domain.Message, error)
}

type Embedder interface {
	Embed(context.Context, []string) ([]*domain.Embeddings, error)
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

	chatSession, err := s.chatSessionRepository.GetChatSession(ctx, initialMessage.ChatSessionID)
	if err != nil {
		return nil, err
	}

	if chatSession.Title == "" {
		title, err := s.generateTitleFromOpenAI(ctx, initialMessage.Content)
		if err != nil {
			return nil, err
		}

		err = s.chatSessionRepository.UpdateChatSessionTitle(ctx, initialMessage.ChatSessionID, title)
		if err != nil {
			return nil, err
		}
	}

	domainEmbeddings, err := s.embedder.Embed(ctx, []string{initialMessage.Content})
	if err != nil {
		return nil, err
	}

	// we only have on text so we only care for the first embedding row
	vectorToFloat32 := helpers.Float64ToFloat32(domainEmbeddings[0].Embeddings)

	accumulatedTextFromSearch, err := s.vectorDB.SemanticSearch(ctx, vectorToFloat32)

	answer, err := s.generateAnswerFromOpenAI(ctx, accumulatedTextFromSearch, initialMessage.Content)
	if err != nil {
		return nil, err
	}

	replyMessage := &domain.Message{
		ChatSessionID: initialMessage.ChatSessionID,
		Content:       answer,
		Sender:        repositories.SYSTEM_SENDER,
	}

	insertedMessageID, err := s.messageRepository.CreateMessage(ctx, replyMessage)
	if err != nil {
		return nil, err
	}

	replyMessage.ID = insertedMessageID

	return replyMessage, nil
}

func (s *MessageService) generateAnswerFromOpenAI(ctx context.Context, text []string, initialMessage string) (string, error) {
	prompt := fmt.Sprintf(`Use the following context to answer the question.
		Context:
		%s
		
		Question:
		%s
		
		Answer:`,
		strings.Join(text, "\n"),
		initialMessage,
	)

	chatCompletion, err := s.openAIClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4_1Nano,
	})
	if err != nil {
		return "", err
	}

	if chatCompletion.Choices[0].Message.Content == "" {
		return "", errors.New("received empty response from LLM")

	}

	return chatCompletion.Choices[0].Message.Content, nil
}

func (s *MessageService) generateTitleFromOpenAI(ctx context.Context, initialMessage string) (string, error) {
	prompt := fmt.Sprintf(`Summarize the following user message into a short, descriptive chat title (max 5 words):
		%s
		Answer:`,
		initialMessage,
	)

	chatCompletion, err := s.openAIClient.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: openai.ChatModelGPT4_1Nano,
	})
	if err != nil {
		return "", err
	}

	if chatCompletion.Choices[0].Message.Content == "" {
		return "", errors.New("received empty response from LLM")

	}

	return chatCompletion.Choices[0].Message.Content, nil
}
