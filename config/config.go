package config

// 配置结构体
type Config struct {
	API struct {
		BaseURL string `yaml:"base_url"`
		Key     string `yaml:"key"`
	} `yaml:"api"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
}

// 全局配置实例
var AppConfig Config

// GetConfig 返回全局配置实例
func GetConfig() *Config {
	return &AppConfig
}

// SetConfig 设置全局配置实例
func SetConfig(cfg Config) {
	AppConfig = cfg
}
