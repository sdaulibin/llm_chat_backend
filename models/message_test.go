package models_test

import (
	"testing"

	"github.com/binginx/star_llm_backend/models"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Import the models package to access the Message struct

func setupTestDB() (*gorm.DB, error) {
	dsn := "host=localhost port=5432 user=postgres password=difyai123456 dbname=star_llm sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func TestMessageOperations(t *testing.T) {
	// 初始化测试数据库连接
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("无法连接到测试数据库: %v", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(&models.Message{}) // Use models.Message
	if err != nil {
		t.Fatalf("自动迁移数据库表失败: %v", err)
	}

	// 创建新消息
	message := &models.Message{ // Use models.Message
		UserID:        "testuser",
		SessionID:     "testsession",
		MessageID:     uuid.New().String(),
		ConversationID: uuid.New().String(),
		Query:         "Test query",
		Answer:        "Test answer",
		IsSafe:        true,
		IsLike:        false,
	}

	err = db.Create(message).Error
	if err != nil {
		t.Fatalf("创建消息失败: %v", err)
	}

	// 获取消息
	var fetchedMessage models.Message // Use models.Message
	err = db.First(&fetchedMessage, message.ID).Error
	if err != nil {
		t.Fatalf("获取消息失败: %v", err)
	}

	if fetchedMessage.Query != message.Query {
		t.Errorf("获取的消息内容不匹配: got %v, want %v", fetchedMessage.Query, message.Query)
	}

	if fetchedMessage.ConversationID != message.ConversationID {
		t.Errorf("获取的对话ID不匹配: got %v, want %v", fetchedMessage.ConversationID, message.ConversationID)
	}

	// 更新消息
	fetchedMessage.Answer = "Updated answer"
	err = db.Save(&fetchedMessage).Error
	if err != nil {
		t.Fatalf("更新消息失败: %v", err)
	}

	// 验证更新
	var updatedMessage models.Message // Use models.Message
	err = db.First(&updatedMessage, message.ID).Error
	if err != nil {
		t.Fatalf("获取更新后的消息失败: %v", err)
	}

	if updatedMessage.Answer != "Updated answer" {
		t.Errorf("更新后的消息内容不匹配: got %v, want %v", updatedMessage.Answer, "Updated answer")
	}

	// 删除消息
	err = db.Delete(&models.Message{}, message.ID).Error // Use models.Message
	if err != nil {
		t.Fatalf("删除消息失败: %v", err)
	}

	// 验证删除
	var deletedMessage models.Message // Use models.Message
	err = db.First(&deletedMessage, message.ID).Error
	if err == nil {
		t.Errorf("消息未被删除")
	}
}
