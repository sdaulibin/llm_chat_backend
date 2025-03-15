package services

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"star_llm_backend/models"
)

// SaveFileToDB 保存文件信息到数据库
func SaveFileToDB(userID, originalFilename, localFilename, filePath, fileID string, fileSize int64) (string, error) {
	// 获取文件类型
	fileType := filepath.Ext(originalFilename)
	if fileType != "" && fileType[0] == '.' {
		fileType = fileType[1:] // 去掉开头的点
	}
	// 创建文件对象
	file := &models.File{
		UserID:           userID,
		FileID:           fileID,
		OriginalFilename: originalFilename,
		LocalFilename:    localFilename,
		FilePath:         filePath,
		FileSize:         fileSize,
		FileType:         fileType,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	log.Printf("[服务] 保存文件信息到数据库: file_id=%s, 原始文件名=%s", fileID, originalFilename)
	err := models.CreateFile(file)
	return fileID, err
}

// GetFileByID 通过文件ID获取文件信息
func GetFileByID(fileID string) (*models.File, error) {
	return models.GetFileByFileID(fileID)
}

// GetFilesByUserID 获取用户的所有文件
func GetFilesByUserID(userID string) ([]models.File, error) {
	return models.GetFilesByUserID(userID)
}

// DeleteFile 删除文件及其数据库记录
func DeleteFile(fileID string) error {
	// 获取文件信息
	file, err := GetFileByID(fileID)
	if err != nil {
		log.Printf("[服务] 删除文件失败: 无法获取文件信息 file_id=%s, 错误: %v", fileID, err)
		return err
	}

	// 删除物理文件
	err = os.Remove(file.FilePath)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("[服务] 删除物理文件失败: file_id=%s, 路径=%s, 错误: %v", fileID, file.FilePath, err)
		return err
	}

	// 删除数据库记录
	return models.DeleteFileByFileID(fileID)
}

// UpdateFileInfo 更新文件信息
func UpdateFileInfo(fileID string, newFilename string) error {
	// 获取文件信息
	file, err := GetFileByID(fileID)
	if err != nil {
		log.Printf("[服务] 更新文件信息失败: 无法获取文件信息 file_id=%s, 错误: %v", fileID, err)
		return err
	}

	// 更新文件名
	file.OriginalFilename = newFilename
	// 更新文件类型
	fileType := filepath.Ext(newFilename)
	if fileType != "" && fileType[0] == '.' {
		fileType = fileType[1:] // 去掉开头的点
	}
	file.FileType = fileType
	file.UpdatedAt = time.Now()

	return models.UpdateFile(file)
}

// GetFilesByType 根据文件类型获取文件
func GetFilesByType(userID string, fileType string) ([]models.File, error) {
	return models.GetFilesByType(userID, fileType)
}

// GetFileStats 获取用户文件统计信息
func GetFileStats(userID string) (map[string]interface{}, error) {
	// 获取用户所有文件
	files, err := GetFilesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// 统计信息
	totalFiles := len(files)
	totalSize := int64(0)
	fileTypeCount := make(map[string]int)

	for _, file := range files {
		totalSize += file.FileSize
		fileTypeCount[file.FileType]++
	}

	// 返回统计结果
	return map[string]interface{}{
		"total_files":     totalFiles,
		"total_size":      totalSize,
		"file_type_count": fileTypeCount,
	}, nil
}

// CheckFileExists 检查文件是否存在
func CheckFileExists(fileID string) (bool, error) {
	file, err := GetFileByID(fileID)
	if err != nil {
		return false, err
	}

	// 检查物理文件是否存在
	_, err = os.Stat(file.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// BatchDeleteFiles 批量删除文件
func BatchDeleteFiles(fileIDs []string) (map[string]error, error) {
	results := make(map[string]error)
	hasError := false

	for _, fileID := range fileIDs {
		err := DeleteFile(fileID)
		results[fileID] = err
		if err != nil {
			hasError = true
		}
	}

	if hasError {
		return results, errors.New("部分文件删除失败")
	}

	return results, nil
}
