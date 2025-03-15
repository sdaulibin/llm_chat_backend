package models

import (
	"database/sql"
	"fmt"
	"log"

	"star_llm_backend/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var sqlDB *sql.DB

// InitDB 初始化数据库连接
func InitDB(dbConfig *config.Config) (*gorm.DB, error) {
	fmt.Printf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Database.Host, dbConfig.Database.Port, dbConfig.Database.User, dbConfig.Database.Password, dbConfig.Database.DBName, dbConfig.Database.SSLMode)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Database.Host, dbConfig.Database.Port, dbConfig.Database.User, dbConfig.Database.Password, dbConfig.Database.DBName, dbConfig.Database.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %v", err)
	}

	// 获取底层的sql.DB对象
	sqlDB, err = DB.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库连接失败: %v", err)
	}

	// 自动迁移数据库表
	// err = DB.AutoMigrate(&Message{})
	// if err != nil {
	// 	return nil, fmt.Errorf("自动迁移数据库表失败: %v", err)
	// }

	log.Println("成功连接到数据库并迁移表")
	return DB, nil
}

// Close 关闭数据库连接
func Close() error {
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil
}
