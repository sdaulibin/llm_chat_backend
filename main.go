package main

import (
	"log"
	"net/http"
	"os"

	"star_llm_backend/config"
	"star_llm_backend/controllers"
	"star_llm_backend/models"
	"star_llm_backend/router"

	"gopkg.in/yaml.v3"
)

// getEnvOrDefault 从环境变量获取值，如果环境变量不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

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

	// 优先使用环境变量覆盖配置
	cfg.API.BaseURL = getEnvOrDefault("DIFY_API_BASE_URL", cfg.API.BaseURL)
	cfg.API.Key = getEnvOrDefault("DIFY_API_KEY", cfg.API.Key)
	cfg.Server.Port = getEnvOrDefault("SERVER_PORT", cfg.Server.Port)
	cfg.Database.Host = getEnvOrDefault("DB_HOST", cfg.Database.Host)
	cfg.Database.Port = getEnvOrDefault("DB_PORT", cfg.Database.Port)
	cfg.Database.User = getEnvOrDefault("DB_USER", cfg.Database.User)
	cfg.Database.Password = getEnvOrDefault("DB_PASSWORD", cfg.Database.Password)
	cfg.Database.DBName = getEnvOrDefault("DB_NAME", cfg.Database.DBName)
	cfg.Database.SSLMode = getEnvOrDefault("DB_SSLMODE", cfg.Database.SSLMode)

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
