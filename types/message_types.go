package types

// ChatMessageRequest 定义聊天消息请求结构体
type ChatMessageRequest struct {
	Query          string                 `json:"query"`
	Inputs         map[string]interface{} `json:"inputs"`
	ResponseMode   string                 `json:"response_mode"`
	User           string                 `json:"user"`
	ConversationID string                 `json:"conversation_id"`
	SeesionID      string                 `json:"session_id"`
}

// ChatMessageResponse 定义聊天消息响应结构体
type ChatMessageResponse struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	CreatedAt      int64  `json:"created_at"`
}

// StreamChunk 定义流式响应数据块结构体
type StreamChunk struct {
	Event          string `json:"event"`
	TaskID         string `json:"task_id"`
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Answer         string `json:"answer"`
	CreatedAt      int64  `json:"created_at"`
}