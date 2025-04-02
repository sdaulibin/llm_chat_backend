package models

import (
	"log"
	"time"
)

// Message 表示数据库中的消息记录
type Message struct {
	ID             int       `gorm:"primaryKey"`
	UserID         string    `gorm:"size:10"`
	SessionID      string    `gorm:"size:32;not null"`
	MessageID      string    `gorm:"type:uuid"`
	ConversationID string    `gorm:"type:uuid"`
	Query          string    `gorm:"type:text;not null"`
	Answer         string    `gorm:"type:text;not null"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	IsSafe         bool      `gorm:"default:false;not null"`
	IsLike         bool      `gorm:"default:false;not null"`
	CurrentID      string    `gorm:"type:uuid"`
	IsStop         bool      `gorm:"default:false;not null"`
	FileID         string    `gorm:"type:uuid"`
	TaskID         string    `gorm:"type:uuid"`
}

// CreateMessage 创建新消息
func CreateMessage(message *Message) error {
	return DB.Create(message).Error
}

// GetMessageByID 通过ID获取消息
func GetMessageByID(id int) (*Message, error) {
	var message Message
	err := DB.First(&message, id).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetMessageByMessageIDAndSessionID 通过MessageID和SessionID获取消息
func GetMessageByMessageIDAndSessionID(messageID, sessionID string) (*Message, error) {
	var message Message
	err := DB.Where("message_id = ? AND session_id = ?", messageID, sessionID).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// UpdateMessageLikeStatus 更新消息的点赞状态
func UpdateMessageLikeStatus(messageID, sessionID string, isLike bool) error {
	return DB.Model(&Message{}).Where("message_id = ? AND session_id = ?", messageID, sessionID).Update("is_like", isLike).Error
}

// UpdateMessageStopStatus 更新消息的停止状态
func UpdateMessageStopStatus(taskId string, isStop bool) error {
	log.Printf("[数据库] 更新消息: task_id=%s", taskId)
	return DB.Debug().Model(&Message{}).Where("task_id = ?", taskId).Update("is_stop", isStop).Error
}

// UpdateMessageByCurrentID 通过CurrentID更新消息
func UpdateMessageByCurrentID(message *Message) error {
	return DB.Model(&Message{}).Where("current_id = ?", message.CurrentID).Updates(map[string]interface{}{
		"answer":          message.Answer,
		"conversation_id": message.ConversationID,
		"message_id":      message.MessageID,
		"task_id":         message.TaskID,
		"updated_at":      message.UpdatedAt,
	}).Error
}

// DeleteMessage 删除消息
func DeleteMessage(id int) error {
	return DB.Delete(&Message{}, id).Error
}
