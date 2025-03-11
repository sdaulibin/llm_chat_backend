# Star LLM Backend

这是一个基于Golang开发的后端服务，用于实现与Dify API的交互。该服务主要功能是透传转发前端请求到Dify API，并将响应返回给前端。

## 功能特点

- 提供HTTP服务器接收前端请求
- 转发请求到Dify API
- 保留原始请求的头信息和认证信息
- 将Dify API的响应完整返回给前端

## 使用方法

### 安装

```bash
# 克隆仓库
git clone https://github.com/binginx/star_llm_backend.git
cd star_llm_backend

# 构建项目
go build
```

### 运行

```bash
# 直接运行
go run main.go

# 或者构建后运行
./star_llm_backend
```

服务器将在8080端口启动，并将请求转发到`http://localhost/v1`。

### API使用

前端应用应该将原本发送到Dify API的请求改为发送到此服务的对应端点：

```
http://localhost:8080/api/[dify-endpoint]
```

例如，发送聊天消息的请求应该从：
```
POST http://localhost/v1/chat-messages
```

改为：
```
POST http://localhost:8080/api/chat-messages
```

所有的请求头（包括认证信息）和请求体保持不变。

## 配置

服务的配置现在使用YAML格式的配置文件`config.yaml`：

```yaml
# Dify API配置
api:
  # Dify API的基础URL
  base_url: "http://10.238.149.28:30000/"
  # Dify API的认证密钥
  key: "your-api-key-here"

# 服务器配置
server:
  # 服务器监听端口
  port: "8090"
```

### 配置说明

- `api.base_url`: Dify API的基础URL
- `api.key`: Dify API的认证密钥（API-Key）
- `server.port`: 本服务器监听的端口

使用前请确保将`api.key`设置为您的实际Dify API密钥。这样可以避免在前端暴露API密钥，提高安全性。

## 许可证

MIT