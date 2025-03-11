package main

import (
	"log"
	"net/http"
	"os"

	"github.com/binginx/star_llm_backend/config"
	"github.com/binginx/star_llm_backend/controllers"
	"github.com/binginx/star_llm_backend/models"
	"github.com/binginx/star_llm_backend/router"
	"gopkg.in/yaml.v3"
)

func main() {
	// 加载配置文件
	configFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	var cfg config.Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		log.Fatalf("无法解析配置文件: %v", err)
	}

	// 设置全局配置
	config.SetConfig(cfg)

	// 初始化数据库连接
	_, err = models.InitDB(&cfg)
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}
	defer models.Close()

	// 创建控制器
	proxyController := controllers.NewProxyController(&cfg)
	fileController := controllers.NewFileController(&cfg)

	// 创建路由
	r := router.NewRouter(proxyController, fileController)
	r.SetupRoutes()

	// 启动服务器
	port := cfg.Server.Port
	if port == "" {
		port = "8090"
	}

	log.Printf("服务器启动在 :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
