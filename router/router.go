package router

import (
	"net/http"
	"strings"

	"star_llm_backend/controllers"
)

// Router 路由管理器
type Router struct {
	ProxyController *controllers.ProxyController
	FileController  *controllers.FileController
}

// NewRouter 创建一个新的路由管理器
func NewRouter(proxyController *controllers.ProxyController, fileController *controllers.FileController) *Router {
	return &Router{
		ProxyController: proxyController,
		FileController:  fileController,
	}
}

// SetupRoutes 设置路由处理器
func (r *Router) SetupRoutes() {
	// 处理文件上传请求
	http.HandleFunc("/chat/api/v1/files/upload", r.FileController.HandleFileUpload)

	// 处理/v1/messages/路径下的请求
	http.HandleFunc("/chat/api/v1/messages/", func(w http.ResponseWriter, req *http.Request) {
		// 检查是否为feedbacks接口
		if strings.Contains(req.URL.Path, "/feedbacks") {
			r.ProxyController.HandleFeedbacks(w, req)
		} else {
			r.ProxyController.ProxyToDify(w, req)
		}
	})

	// 处理/v1/路径下的请求
	http.HandleFunc("/chat/api/v1/", func(w http.ResponseWriter, req *http.Request) {
		// 检查是否为feedbacks接口
		if strings.Contains(req.URL.Path, "/chat-messages") {
			r.ProxyController.ProxyToDify(w, req)
		} else {
			// 拒绝处理其他请求
			http.Error(w, "不支持的请求路径", http.StatusNotFound)
		}
	})
	// 处理所有其他请求
	//http.HandleFunc("/", r.ProxyController.ProxyToDify)
}
