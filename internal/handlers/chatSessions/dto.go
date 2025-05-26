package chatSessions

import (
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
)

type UserChatSessionsResponse struct {
	Sessions     []ChatSessionResponse `json:"sessions,omitempty"`
	ErrorMessage string                `json:"errorMessage,omitempty"`
}

type ChatSessionResponse struct {
	ID           string    `json:"id,omitempty"`
	Title        string    `json:"title,omitempty"`
	CreatedAt    string    `json:"createdAt,omitempty"`
	UpdatedAt    string    `json:"updatedAt,omitempty"`
	Messages     []Message `json:"messages,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

func ChatSessionResponseFromModel(domainChatSession *domain.ChatSession) *ChatSessionResponse {
	messages := make([]Message, len(domainChatSession.Messages))
	for i, msg := range domainChatSession.Messages {
		messages[i] = Message{
			ID:        msg.ID.String(),
			Sender:    msg.Sender,
			Content:   msg.Content,
			Timestamp: msg.Timestamp.String(),
		}
	}

	return &ChatSessionResponse{
		ID:        domainChatSession.ID.String(),
		Title:     domainChatSession.Title,
		CreatedAt: domainChatSession.CreatedAt.String(),
		UpdatedAt: domainChatSession.UpdatedAt.String(),
		Messages:  messages,
	}
}

func UserChatSessionsResponseFromModel(domainChatSession []*domain.ChatSession) *UserChatSessionsResponse {
	sessions := make([]ChatSessionResponse, len(domainChatSession))
	for i, session := range domainChatSession {
		sessions[i] = *ChatSessionResponseFromModel(session)
	}

	return &UserChatSessionsResponse{
		Sessions: sessions,
	}
}

type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender" enum:"USER,SYSTEM" example:"USER"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

//type CreateChatSessionRequest struct {
//	UserID string `json:"chatSessionID"`
//}

type SendMessageRequest struct {
	UserID    string `json:"chatSessionID"`
	SessionID string `json:"sessionId"`
	Content   string `json:"content"`
}

type SendMessageResponse struct {
	Body struct {
		UserMessage   Message `json:"userMessage"`
		SystemMessage Message `json:"systemMessage"`
	}
}
