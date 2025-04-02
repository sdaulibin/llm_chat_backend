package controllers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"star_llm_backend/config"
	"star_llm_backend/models"
	"star_llm_backend/services"
	"star_llm_backend/types"

	"github.com/google/uuid"
)

// ProxyController 处理代理请求的控制器
type ProxyController struct {
	Config *config.Config
}

// NewProxyController 创建一个新的代理控制器
func NewProxyController(cfg *config.Config) *ProxyController {
	return &ProxyController{
		Config: cfg,
	}
}

// ProxyToDify 转发请求到Dify API并返回响应
func (pc *ProxyController) ProxyToDify(w http.ResponseWriter, r *http.Request) {
	// 处理CORS预检请求
	if r.Method == "OPTIONS" {
		pc.handleCORS(w)
		return
	}

	// 为所有响应设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 1. 提取API路径
	apiPath := strings.TrimPrefix(r.URL.Path, "/chat/api/")

	// 构建完整URL，包括查询参数
	difyURL := pc.Config.API.BaseURL + apiPath
	if r.URL.RawQuery != "" {
		difyURL += "?" + r.URL.RawQuery
	}

	log.Printf("[请求] URL: %s, 方法: %s", difyURL, r.Method)

	// 2. 读取请求体
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[错误] 读取请求体失败: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 打印请求体内容（如果是文本格式）
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/") {
		log.Printf("[请求] 请求体: %s", string(bodyBytes))
	} else {
		log.Printf("[请求] 请求体: 二进制数据或非文本格式 (Content-Type: %s)", contentType)
	}

	// 检查是否为chat-messages接口请求，并解析请求体
	var chatRequest types.ChatMessageRequest
	var userID string
	var query string
	var sessionID string
	var conversationID string
	//var taskId string
	var currentID string
	var fileId string

	if strings.Contains(apiPath, "chat-messages") && r.Method == "POST" {
		if err := json.Unmarshal(bodyBytes, &chatRequest); err == nil {
			// 提取用户ID、查询内容、会话ID和对话ID
			currentID = uuid.NewString()
			userID = chatRequest.User
			query = chatRequest.Query
			sessionID = chatRequest.SeesionID
			if len(chatRequest.ConversationID) > 0 {
				conversationID = chatRequest.ConversationID
			} else {
				conversationID = uuid.Nil.String()
			}
			log.Printf("[提取] 用户ID: %s, 输入: %s, 会话ID: %s, 对话ID: %s", userID, query, sessionID, conversationID)

			// 提取文件信息
			if len(chatRequest.Files) > 0 {
				for _, fileInfo := range chatRequest.Files {
					if fileInfo.Type == "document" && fileInfo.TransferMethod == "local_file" {
						log.Printf("[提取] 文件上传ID: %s, 类型: %s, 传输方式: %s",
							fileInfo.UploadFileID, fileInfo.Type, fileInfo.TransferMethod)
						fileId = fileInfo.UploadFileID
					}
				}
			} else {
				fileId = uuid.Nil.String()
			}

			// 检查是否需要保存消息到数据库
			if sessionID != "" && query != "" {
				err := services.SaveMessageToDB(currentID, sessionID, query, "", userID, conversationID, fileId, uuid.Nil.String())
				if err != nil {
					log.Printf("[错误] 保存消息到数据库失败: %v", err)
				}
			}
		} else {
			log.Printf("[错误] 解析chat-messages请求体失败: %v", err)
		}
	}

	if strings.Contains(apiPath, "chat-messages") && strings.Contains(apiPath, "/stop") && r.Method == "POST" {

		// 从URL中提取message_id
		apiPath := strings.TrimPrefix(r.URL.Path, "/chat/api/")
		pathParts := strings.Split(apiPath, "/")
		log.Printf("[pathParts]: %s", pathParts)
		if len(pathParts) < 3 {
			log.Printf("[错误] 无效的URL路径: %s", r.URL.Path)
			http.Error(w, "Invalid URL path", http.StatusBadRequest)
			return
		}
		taskId := pathParts[len(pathParts)-2]
		log.Printf("[提取] 从URL中提取的task_id: %s", taskId)
		err := models.UpdateMessageStopStatus(taskId, true)
		if err != nil {
			log.Printf("[错误]保存stop 状态到数据库失败: %v", err)
		}
	}

	// 3. 创建新的请求到Dify
	difyReq, err := http.NewRequest(r.Method, difyURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		http.Error(w, "Error creating request to Dify", http.StatusInternalServerError)
		return
	}

	// 4. 设置请求头信息
	log.Println("[请求] 请求头:")
	// 首先设置Authorization头，使用配置文件中的API密钥
	difyReq.Header.Set("Authorization", "Bearer "+pc.Config.API.Key)
	log.Printf("[请求] Authorization: %s", "******")

	// 复制其他原始请求的头信息，但跳过Authorization头
	for name, values := range r.Header {
		// 跳过Authorization头，因为我们已经设置了
		if strings.ToLower(name) == "authorization" {
			continue
		}
		for _, value := range values {
			difyReq.Header.Add(name, value)
			log.Printf("[请求] %s: %s", name, value)
		}
	}

	// 5. 发送请求到Dify
	log.Println("[请求] 发送请求到Dify...")
	client := &http.Client{}
	difyResp, err := client.Do(difyReq)
	if err != nil {
		log.Printf("[错误] 转发请求到Dify失败: %v", err)
		http.Error(w, "Error forwarding request to Dify", http.StatusInternalServerError)
		return
	}
	log.Printf("[响应] 收到Dify响应，状态码: %d", difyResp.StatusCode)
	defer difyResp.Body.Close()

	// 6. 复制Dify响应的头信息到我们的响应，但跳过CORS相关头信息
	log.Println("[响应] 响应头:")
	for name, values := range difyResp.Header {
		// 跳过CORS相关的头信息，因为我们已经设置了
		if strings.ToLower(name) == "access-control-allow-origin" ||
			strings.ToLower(name) == "access-control-allow-methods" ||
			strings.ToLower(name) == "access-control-allow-headers" ||
			strings.ToLower(name) == "access-control-max-age" {
			log.Printf("[响应] 跳过CORS头 %s: %s", name, values)
			continue
		}
		for _, value := range values {
			w.Header().Add(name, value)
			log.Printf("[响应] %s: %s", name, value)
		}
	}

	// 7. 设置响应状态码
	w.WriteHeader(difyResp.StatusCode)

	// 8. 复制Dify响应体到我们的响应
	contentType = difyResp.Header.Get("Content-Type")

	// 检查是否为流式响应
	if strings.Contains(contentType, "text/event-stream") {
		pc.handleStreamResponse(w, difyResp, userID, sessionID, query, conversationID, currentID)
	} else {
		pc.handleNormalResponse(w, difyResp)
	}
	// if strings.Contains(apiPath, "chat-messages") && strings.Contains(apiPath, "/stop") && r.Method == "POST" {
	// 	pc.handleStopResponse(difyResp, taskId)
	// }
}

// HandleFeedbacks 处理消息反馈（点赞）请求
func (pc *ProxyController) HandleFeedbacks(w http.ResponseWriter, r *http.Request) {
	// 处理CORS预检请求
	if r.Method == "OPTIONS" {
		pc.handleCORS(w)
		return
	}

	// 为所有响应设置CORS头
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 从URL中提取message_id
	apiPath := strings.TrimPrefix(r.URL.Path, "/chat/api/")
	pathParts := strings.Split(apiPath, "/")
	log.Printf("[pathParts]: %s", pathParts)
	if len(pathParts) < 3 {
		log.Printf("[错误] 无效的URL路径: %s", r.URL.Path)
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	messageID := pathParts[len(pathParts)-2]
	log.Printf("[提取] 从URL中提取的message_id: %s", messageID)

	// 读取请求体
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[错误] 读取请求体失败: %v", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 解析请求体
	var feedbackRequest struct {
		Rating    string `json:"rating"`
		User      string `json:"user"`
		Content   string `json:"content"`
		SessionID string `json:"session_id"`
	}

	if err := json.Unmarshal(bodyBytes, &feedbackRequest); err != nil {
		log.Printf("[错误] 解析反馈请求体失败: %v", err)
		http.Error(w, "Error parsing feedback request", http.StatusBadRequest)
		return
	}

	log.Printf("[提取] 反馈信息: rating=%s, user=%s, session_id=%s",
		feedbackRequest.Rating, feedbackRequest.User, feedbackRequest.SessionID)

	// 根据rating更新is_like字段
	// 当rating为"like"时，is_like为true；否则为false
	isLike := feedbackRequest.Rating == "like"
	err = models.UpdateMessageLikeStatus(messageID, feedbackRequest.SessionID, isLike)
	if err != nil {
		log.Printf("[错误] 更新消息点赞状态失败: %v", err)
		http.Error(w, "Error updating message like status", http.StatusInternalServerError)
		return
	}

	log.Printf("[成功] 已更新消息点赞状态: message_id=%s, session_id=%s, is_like=%v",
		messageID, feedbackRequest.SessionID, isLike)

	// 重新设置请求体，以便转发到Dify
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 转发请求到Dify
	pc.ProxyToDify(w, r)
}

// 处理CORS预检请求
func (pc *ProxyController) handleCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Max-Age", "86400")
	w.WriteHeader(http.StatusOK)
}

// 处理流式响应
func (pc *ProxyController) handleStreamResponse(w http.ResponseWriter, difyResp *http.Response, userID, sessionID, query, conversationID, currentID string) {
	// 流式响应处理
	log.Println("[响应] 检测到流式响应，开始转发数据流...")

	// 创建缓冲读取器
	reader := bufio.NewReader(difyResp.Body)

	// 用于累积流式响应的完整答案
	var fullAnswer string
	var messageID string
	var task_id string
	// 标记是否已经处理了第一条消息
	var isFirstMessage bool = false

	// 逐行读取并转发
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("[错误] 读取流式响应行失败: %v", err)
			}
			break
		}

		// 记录日志（仅记录data开头的行）
		lineStr := string(line)
		if strings.HasPrefix(lineStr, "data: ") {
			log.Printf("[响应] 流式数据: %s", strings.TrimSpace(lineStr))

			// 解析流式数据
			dataContent := strings.TrimPrefix(lineStr, "data: ")
			var chunk types.StreamChunk
			if err := json.Unmarshal([]byte(dataContent), &chunk); err == nil {
				// 只处理message事件
				if chunk.Event == "message" {
					// 累积答案
					fullAnswer += chunk.Answer
					// 保存消息ID
					if messageID == "" && chunk.MessageID != "" {
						messageID = chunk.MessageID
					}
					if conversationID == uuid.Nil.String() && chunk.ConversationID != "" {
						conversationID = chunk.ConversationID
					}
					if task_id == "" && chunk.TaskID != "" {
						task_id = chunk.TaskID
					}

					// 在第一条消息时就更新数据库
					if !isFirstMessage && sessionID != "" && query != "" && messageID != "" {
						isFirstMessage = true
						log.Printf("[数据库] 第一条流式消息: 会话ID=%s, 查询=%s, 对话ID=%s, 消息ID=%s, task_id=%s", sessionID, query, conversationID, messageID, task_id)
						err := services.UpdateMessageToDB(currentID, "", conversationID, messageID, task_id)
						if err != nil {
							log.Printf("[错误] 保存第一条消息到数据库失败: %v", err)
						}
					}
				}
				// 如果是消息结束事件，保存到数据库
				if chunk.Event == "message_end" && sessionID != "" && query != "" && fullAnswer != "" {
					log.Printf("[数据库] 更新消息: 会话ID=%s, 查询=%s, 答案长度=%d, 对话ID=%s, 消息ID=%s, task_id=%s", sessionID, query, len(fullAnswer), chunk.ConversationID, messageID, task_id)
					err := services.UpdateMessageToDB(currentID, fullAnswer, conversationID, messageID, task_id)
					if err != nil {
						log.Printf("[错误] 保存消息到数据库失败: %v", err)
					}
				}
			} else {
				log.Printf("[错误] 解析流式数据失败: %v", err)
			}
		}

		// 直接写入响应
		_, err = w.Write(line)
		if err != nil {
			log.Printf("[错误] 写入流式响应失败: %v", err)
			break
		}

		// 立即刷新，确保数据发送到客户端
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
}

// 处理普通响应
func (pc *ProxyController) handleNormalResponse(w http.ResponseWriter, difyResp *http.Response) {
	// 非流式响应处理
	respBodyBytes, err := io.ReadAll(difyResp.Body)
	if err != nil {
		log.Printf("[错误] 读取响应体失败: %v", err)
		return
	}

	// 打印响应体内容（如果是文本格式）
	// contentType := difyResp.Header.Get("Content-Type")
	// if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/") {
	// 	log.Printf("[响应] 响应体: %s", string(respBodyBytes))

	// 	// 如果是chat-messages接口的响应，解析并保存到数据库
	// 	if strings.Contains(apiPath, "chat-messages") && sessionID != "" && query != "" {
	// 		var chatResponse types.ChatMessageResponse
	// 		if err := json.Unmarshal(respBodyBytes, &chatResponse); err == nil {
	// 			log.Printf("[数据库] 保存消息: 会话ID=%s, 查询=%s, 答案长度=%d, 对话ID=%s, 消息ID=%s", sessionID, query, len(chatResponse.Answer), chatResponse.ConversationID, chatResponse.MessageID)
	// 			err := services.SaveMessageToDB(sessionID, query, chatResponse.Answer, userID, chatResponse.ConversationID, chatResponse.MessageID)
	// 			if err != nil {
	// 				log.Printf("[错误] 保存消息到数据库失败: %v", err)
	// 			}
	// 		} else {
	// 			log.Printf("[错误] 解析chat-messages响应失败: %v", err)
	// 		}
	// 	}
	// } else {
	// 	log.Printf("[响应] 响应体: 二进制数据或非文本格式 (Content-Type: %s)", contentType)
	// }

	// 重新创建一个新的响应体，因为原来的已经被读取
	difyResp.Body = io.NopCloser(bytes.NewBuffer(respBodyBytes))

	// 复制响应体到客户端
	_, err = io.Copy(w, difyResp.Body)
	if err != nil {
		log.Printf("[错误] 复制响应体失败: %v", err)
	}
}

func (pc *ProxyController) handleStopResponse(difyResp *http.Response, taskId string) {
	respBodyBytes, err := io.ReadAll(difyResp.Body)
	if err != nil {
		log.Printf("[错误]读取响应体失败: %v", err)
		return
	}
	var stopResponse types.StopResponse
	if err := json.Unmarshal(respBodyBytes, &stopResponse); err == nil {
		if stopResponse.Result == "success" {
			log.Printf("[数据库] 更新消息: task_id=%s", taskId)
			err := models.UpdateMessageStopStatus(taskId, true)
			if err != nil {
				log.Printf("[错误]保存stop 状态到数据库失败: %v", err)
			}
		}
	} else {
		log.Printf("[错误]解析stop-messages响应失败: %v", err)
	}
}
