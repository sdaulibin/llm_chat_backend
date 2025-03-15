package services

import (
	"log"
	"time"

	"star_llm_backend/models"

	"github.com/google/uuid"
)

// SaveMessageToDB 保存消息到数据库
func SaveMessageToDB(sessionID, query, answer, userID, conversationID string, messageID ...string) error {
	// 创建消息对象
	message := &models.Message{
		UserID:         userID,
		SessionID:      sessionID,
		ConversationID: conversationID,
		Query:          query,
		Answer:         answer,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		IsSafe:         false,
		IsLike:         false,
	}

	// 如果提供了messageID，则使用它，否则生成一个新的UUID
	if len(messageID) > 0 && messageID[0] != "" {
		message.MessageID = messageID[0]
	} else {
		message.MessageID = uuid.New().String()
	}

	log.Printf("[服务] 保存消息到数据库: message_id=%s\n >>>>content:%s", message.MessageID, message.Answer)
	return models.CreateMessage(message)
}
