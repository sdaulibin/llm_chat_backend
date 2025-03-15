#!/bin/bash
# 该脚本用于将config.yaml中的配置项设置为环境变量
# 适用于Linux/Unix系统

# Dify API配置
export DIFY_API_BASE_URL="http://localhost/"
export DIFY_API_KEY="app-2gyyyTpDY8OFhXB1mFB1MO3F"

# 服务器配置
export SERVER_PORT="8090"

# 数据库配置
export DB_HOST="10.238.149.28"
export DB_PORT="30432"
export DB_USER="postgres"
export DB_PASSWORD="ioasit123456"
export DB_NAME="ioa"
export DB_SSLMODE="disable"

echo "环境变量已设置完成！"
echo "使用方法: source ./scripts/set_env.sh"