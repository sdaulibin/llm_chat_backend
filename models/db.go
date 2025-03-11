package models

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(host, port, user, password, dbname, sslmode string) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}

	// 自动迁移数据库表
	err = DB.AutoMigrate(&Message{})
	if err != nil {
		return fmt.Errorf("自动迁移数据库表失败: %v", err)
	}

	log.Println("成功连接到数据库并迁移表")
	return nil
}