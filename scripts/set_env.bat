@echo off
:: 该脚本用于将config.yaml中的配置项设置为环境变量
:: 适用于Windows系统

:: Dify API配置
set DIFY_API_BASE_URL=http://localhost/
set DIFY_API_KEY=app-2gyyyTpDY8OFhXB1mFB1MO3F

:: 服务器配置
set SERVER_PORT=8090

:: 数据库配置
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=difyai123456
set DB_NAME=star_llm
set DB_SSLMODE=disable

echo 环境变量已设置完成！
echo 使用方法: 直接运行此批处理文件