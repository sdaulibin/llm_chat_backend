package router

import (
	"net/http"
	"strings"

	"github.com/binginx/star_llm_backend/controllers"
)

// Router 路由管理器
type Router struct {
	ProxyController *controllers.ProxyController
}

// NewRouter 创建一个新的路由管理器
func NewRouter(proxyController *controllers.ProxyController) *Router {
	return &Router{
		ProxyController: proxyController,
	}
}

// SetupRoutes 设置路由处理器
func (r *Router) SetupRoutes() {
	// 处理/v1/messages/路径下的请求
	http.HandleFunc("/v1/messages/", func(w http.ResponseWriter, req *http.Request) {
		// 检查是否为feedbacks接口
		if strings.Contains(req.URL.Path, "/feedbacks") {
			r.ProxyController.HandleFeedbacks(w, req)
		} else {
			r.ProxyController.ProxyToDify(w, req)
		}
	})

	// 处理所有其他请求
	http.HandleFunc("/", r.ProxyController.ProxyToDify)
}