package chatSessions

type ChatSessionResponse struct {
	ID           string    `json:"id,omitempty"`
	Title        string    `json:"title,omitempty"`
	CreatedAt    string    `json:"createdAt,omitempty"`
	UpdatedAt    string    `json:"updatedAt,omitempty"`
	Messages     []Message `json:"messages,omitempty"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
}

type Message struct {
	ID        string `json:"id"`
	Sender    string `json:"sender" enum:"USER,SYSTEM" example:"USER"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

//type CreateChatSessionRequest struct {
//	UserID string `json:"userId"`
//}

type SendMessageRequest struct {
	UserID    string `json:"userId"`
	SessionID string `json:"sessionId"`
	Content   string `json:"content"`
}

type SendMessageResponse struct {
	Body struct {
		UserMessage   Message `json:"userMessage"`
		SystemMessage Message `json:"systemMessage"`
	}
}
