package models

import (
	"time"

	"github.com/google/uuid"
)

// File 表示数据库中的文件记录
type File struct {
	ID               int       `gorm:"primaryKey"`
	UserID           string    `gorm:"size:10"`
	FileID           string    `gorm:"type:uuid;not null"`
	OriginalFilename string    `gorm:"size:256;not null"`
	LocalFilename    string    `gorm:"size:256;not null"`
	FilePath         string    `gorm:"size:512;not null"`
	FileSize         int64     `gorm:"not null"`
	FileType         string    `gorm:"size:50;not null"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// CreateFile 创建新文件记录
func CreateFile(file *File) error {
	return DB.Create(file).Error
}

// GetFileByID 通过ID获取文件记录
func GetFileByID(id int) (*File, error) {
	var file File
	err := DB.First(&file, id).Error
	return &file, err
}

// GetFileByFileID 通过FileID获取文件记录
func GetFileByFileID(fileID string) (*File, error) {
	var file File
	err := DB.Where("file_id = ?", fileID).First(&file).Error
	return &file, err
}

// GetFilesByUserID 获取用户的所有文件记录
func GetFilesByUserID(userID string) ([]File, error) {
	var files []File
	err := DB.Where("user_id = ?", userID).Find(&files).Error
	return files, err
}

// GenerateFileID 生成新的文件ID
func GenerateFileID() string {
	return uuid.New().String()
}

// DeleteFileByFileID 通过FileID删除文件记录
func DeleteFileByFileID(fileID string) error {
	return DB.Where("file_id = ?", fileID).Delete(&File{}).Error
}

// UpdateFile 更新文件记录
func UpdateFile(file *File) error {
	return DB.Save(file).Error
}

// GetFilesByType 根据文件类型获取文件记录
func GetFilesByType(userID string, fileType string) ([]File, error) {
	var files []File
	err := DB.Where("user_id = ? AND file_type = ?", userID, fileType).Find(&files).Error
	return files, err
}
