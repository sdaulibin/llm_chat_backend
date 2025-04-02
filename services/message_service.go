package services

import (
	"log"
	"time"

	"star_llm_backend/models"

	"github.com/google/uuid"
)

// SaveMessageToDB 保存消息到数据库
func SaveMessageToDB(currentID, sessionID, query, answer, userID, conversationID, fileId, messageID string) error {
	// 创建消息对象
	message := &models.Message{
		CurrentID:      currentID,
		UserID:         userID,
		SessionID:      sessionID,
		MessageID:      messageID,
		ConversationID: conversationID,
		Query:          query,
		Answer:         answer,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		IsSafe:         false,
		IsLike:         false,
		FileID:         fileId,
		TaskID:         uuid.Nil.String(),
	}

	// 如果提供了messageID，则使用它，否则生成一个新的UUID
	if len(messageID) > 0 && messageID != uuid.Nil.String() {
		message.MessageID = messageID
	} else {
		message.MessageID = uuid.New().String()
	}

	log.Printf("[服务] 保存消息到数据库: message_id=%s\n >>>>content:%s", message.MessageID, message.Answer)
	return models.CreateMessage(message)
}

// UpdateMessageToDB 根据CurrentID更新消息到数据库
func UpdateMessageToDB(currentID, answer, conversationID, messageID, taskID string) error {
	// 创建消息对象
	message := &models.Message{
		CurrentID:      currentID,
		Answer:         answer,
		ConversationID: conversationID,
		MessageID:      messageID,
		UpdatedAt:      time.Now(),
		TaskID:         taskID,
	}

	log.Printf("[服务] 更新消息到数据库: CurrentID=%s\n >>>>content:%s", message.CurrentID, message.Answer)
	return models.UpdateMessageByCurrentID(message)
}
