package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"star_llm_backend/config"
	"star_llm_backend/services"
)

// FileController 处理文件上传的控制器
type FileController struct {
	Config    *config.Config
	UploadDir string
}

// NewFileController 创建一个新的文件控制器
func NewFileController(cfg *config.Config) *FileController {
	// 创建上传目录
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, 0755)
	return &FileController{
		Config:    cfg,
		UploadDir: uploadDir,
	}
}

// HandleFileUpload 处理文件上传请求
func (fc *FileController) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	// 处理CORS预检请求
	if r.Method == "OPTIONS" {
		fc.handleCORS(w)
		return
	}

	// 为所有响应设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 检查请求方法
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 检查Content-Type
	if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
		http.Error(w, "Content-Type must be multipart/form-data", http.StatusBadRequest)
		return
	}

	log.Printf("[文件上传] 接收到文件上传请求")

	// 解析multipart表单，32MB最大内存
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("[错误] 解析multipart表单失败: %v", err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// 获取用户标识
	userID := r.FormValue("user")
	if userID == "" {
		log.Printf("[错误] 缺少用户标识")
		http.Error(w, "User identifier is required", http.StatusBadRequest)
		return
	}

	// 获取上传的文件
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Printf("[错误] 获取上传文件失败: %v", err)
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 检查文件类型
	fileExt := filepath.Ext(handler.Filename)
	validExts := map[string]bool{
		// 文本和文档格式
		".txt":      true,
		".markdown": true,
		".mdx":      true,
		".pdf":      true,
		".html":     true,
		".xlsx":     true,
		".xls":      true,
		".docx":     true,
		".csv":      true,
		".md":       true,
		".htm":      true,
	}

	if !validExts[strings.ToLower(fileExt)] {
		log.Printf("[错误] 不支持的文件类型: %s", fileExt)
		http.Error(w, "Unsupported file type", http.StatusUnsupportedMediaType)
		return
	}

	// 创建用户目录
	userDir := filepath.Join(fc.UploadDir, userID)
	os.MkdirAll(userDir, 0755)

	// 保存文件到本地
	localFilePath := filepath.Join(userDir, handler.Filename)
	localFile, err := os.Create(localFilePath)
	if err != nil {
		log.Printf("[错误] 创建本地文件失败: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer localFile.Close()

	// 将上传的文件内容复制到本地文件
	file.Seek(0, 0) // 重置文件指针到开始位置
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Printf("[错误] 读取上传文件内容失败: %v", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// 写入文件内容到本地文件
	_, err = localFile.Write(fileBytes)
	if err != nil {
		log.Printf("[错误] 写入文件内容到本地失败: %v", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	log.Printf("[文件上传] 文件已保存到本地: %s", localFilePath)

	// 转发文件到Dify
	apiPath := strings.TrimPrefix(r.URL.Path, "/chat/api/")
	difyURL := fc.Config.API.BaseURL + apiPath
	log.Printf("[文件上传] 地址: %s", difyURL)
	// 创建一个新的multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加用户标识
	err = writer.WriteField("user", userID)
	if err != nil {
		log.Printf("[错误] 写入用户标识失败: %v", err)
		http.Error(w, "Error forwarding file", http.StatusInternalServerError)
		return
	}

	// 添加文件
	part, err := writer.CreateFormFile("file", handler.Filename)
	if err != nil {
		log.Printf("[错误] 创建文件表单失败: %v", err)
		http.Error(w, "Error forwarding file", http.StatusInternalServerError)
		return
	}

	// 将文件内容写入part
	_, err = part.Write(fileBytes)
	if err != nil {
		log.Printf("[错误] 写入文件内容失败: %v", err)
		http.Error(w, "Error forwarding file", http.StatusInternalServerError)
		return
	}

	// 关闭writer
	err = writer.Close()
	if err != nil {
		log.Printf("[错误] 关闭writer失败: %v", err)
		http.Error(w, "Error forwarding file", http.StatusInternalServerError)
		return
	}

	// 创建请求
	difyReq, err := http.NewRequest("POST", difyURL, body)
	if err != nil {
		log.Printf("[错误] 创建Dify请求失败: %v", err)
		http.Error(w, "Error forwarding file", http.StatusInternalServerError)
		return
	}

	// 设置请求头
	difyReq.Header.Set("Content-Type", writer.FormDataContentType())
	difyReq.Header.Set("Authorization", "Bearer "+fc.Config.API.Key)

	// 发送请求
	client := &http.Client{}
	difyResp, err := client.Do(difyReq)
	if err != nil {
		log.Printf("[错误] 发送Dify请求失败: %v", err)
		http.Error(w, "Error forwarding file", http.StatusInternalServerError)
		return
	}
	defer difyResp.Body.Close()

	log.Printf("[文件上传] Dify响应状态码: %d", difyResp.StatusCode)

	// 读取Dify响应
	respBody, err := io.ReadAll(difyResp.Body)
	if err != nil {
		log.Printf("[错误] 读取Dify响应失败: %v", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	for name, values := range difyResp.Header {
		// 跳过CORS相关的头信息，因为我们已经设置了
		if strings.ToLower(name) == "access-control-allow-origin" ||
			strings.ToLower(name) == "access-control-allow-methods" ||
			strings.ToLower(name) == "access-control-allow-headers" ||
			strings.ToLower(name) == "access-control-max-age" {
			continue
		}
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	// 设置响应状态码
	w.WriteHeader(difyResp.StatusCode)

	// 返回Dify响应
	w.Write(respBody)

	// 记录上传成功信息
	if difyResp.StatusCode == http.StatusCreated {
		var fileInfo map[string]interface{}
		if err := json.Unmarshal(respBody, &fileInfo); err == nil {
			fileID, _ := fileInfo["id"].(string)
			log.Printf("[文件上传] 文件上传成功: 本地路径=%s, Dify文件ID=%s", localFilePath, fileID)

			// 获取文件大小
			fileSize := int64(len(fileBytes))

			// 保存文件信息到数据库
			dbFileID, err := services.SaveFileToDB(userID, handler.Filename, handler.Filename, localFilePath, fileID, fileSize)
			if err != nil {
				log.Printf("[错误] 保存文件信息到数据库失败: %v", err)
			} else {
				log.Printf("[文件上传] 文件信息已保存到数据库: file_id=%s, 原始文件名=%s, 文件大小=%d字节", dbFileID, handler.Filename, fileSize)
			}
		}
	}
}

// 处理CORS预检请求
func (fc *FileController) handleCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusOK)
}
