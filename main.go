package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/binginx/star_llm_backend/config"
	"github.com/binginx/star_llm_backend/controllers"
	"github.com/binginx/star_llm_backend/models"
	"github.com/binginx/star_llm_backend/router"
	"gopkg.in/yaml.v3"
)

// 加载配置文件
func loadConfig() (*config.Config, error) {
	// 读取配置文件
	configFile, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("无法读取配置文件: %v", err)
	}

	// 解析YAML
	var cfg config.Config
	err = yaml.Unmarshal(configFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %v", err)
	}

	// 验证必要的配置项
	if cfg.API.BaseURL == "" {
		return nil, fmt.Errorf("配置文件中缺少API基础URL")
	}
	if cfg.API.Key == "" {
		return nil, fmt.Errorf("配置文件中缺少API密钥")
	}
	if cfg.Server.Port == "" {
		// 使用默认端口
		cfg.Server.Port = "8090"
	}

	// 设置全局配置
	config.SetConfig(cfg)

	return &cfg, nil
}

func main() {
	// 加载配置文件
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库连接
	dbConfig := cfg.Database
	err = models.InitDB(dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 创建控制器
	proxyController := controllers.NewProxyController(cfg)
	
	// 创建路由管理器并设置路由
	r := router.NewRouter(proxyController)
	r.SetupRoutes()

	// 启动服务器
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	fmt.Printf("Server started on port %s\n", cfg.Server.Port)
	fmt.Printf("Proxying requests to Dify API at %s\n", cfg.API.BaseURL)
	fmt.Printf("Using API Key from config file\n")
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
