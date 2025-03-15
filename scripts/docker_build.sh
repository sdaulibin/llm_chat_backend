#!/bin/bash

# 设置变量
IMAGE_NAME="star_llm_backend"
IMAGE_VERSION="1.0.0"
DOCKER_REGISTRY=""

# 环境变量默认值
DIFY_API_BASE_URL=""
DIFY_API_KEY=""
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="starllm"
DB_PASSWORD="starllm123"
DB_NAME="starllm_db"
DB_SSLMODE="disable"

# 显示帮助信息
show_help() {
    echo "使用方法: $0 [选项]"
    echo "选项:"
    echo "  -h, --help     显示帮助信息"
    echo "  -r, --registry  设置Docker镜像仓库地址 (例如: docker.io/username)"
    echo "  -v, --version  设置镜像版本 (默认: ${IMAGE_VERSION})"
    exit 0
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            ;;
        -r|--registry)
            DOCKER_REGISTRY="$2/"
            shift 2
            ;;
        -v|--version)
            IMAGE_VERSION="$2"
            shift 2
            ;;
        *)
            echo "未知选项: $1"
            show_help
            ;;
    esac
done

# 构建Docker镜像
echo "正在构建Docker镜像..."
docker build -t "${DOCKER_REGISTRY}${IMAGE_NAME}:${IMAGE_VERSION}" .

# 如果指定了镜像仓库，则推送镜像
if [ -n "$DOCKER_REGISTRY" ]; then
    echo "正在推送镜像到 ${DOCKER_REGISTRY}..."
    docker push "${DOCKER_REGISTRY}${IMAGE_NAME}:${IMAGE_VERSION}"
fi

echo "
构建完成！

使用说明：
1. 构建镜像：
   ./scripts/docker-build.sh

2. 指定版本构建：
   ./scripts/docker-build.sh -v 2.0.0

3. 构建并推送到镜像仓库：
   ./scripts/docker-build.sh -r docker.io/username

4. 运行容器：
   docker run -d -p 8090:8090 ${DOCKER_REGISTRY}${IMAGE_NAME}:${IMAGE_VERSION}

5. 查看容器日志：
   docker logs $(docker ps -qf ancestor=${DOCKER_REGISTRY}${IMAGE_NAME}:${IMAGE_VERSION})

服务器将在 http://localhost:8090 上启动
"