# 构建阶段

# FROM golang:1.20.2 AS builder
FROM harbor.devops.qdb.com/devops/golang:1.20.2 AS builder

# 设置工作目录
WORKDIR /app

COPY . .
ARG ARCH
ARG RELEASE_VERSION
RUN RELEASE_VERSION=${RELEASE_VERSION} make build-star_llm_backend.amd64

# 运行阶段
FROM harbor.devops.qdb.com/devops/alpine:3.12
# FROM alpine:3.12
# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/bin/star_llm_backend .
# 复制配置文件
COPY --from=builder /app/config /app/config

# 设置环境变量
ENV DIFY_API_BASE_URL="http://nginx.docker.orb.local/"
ENV DIFY_API_KEY="app-2gyyyTpDY8OFhXB1mFB1MO3F"
ENV SERVER_PORT="8090"
ENV DB_HOST="10.238.149.28"
ENV DB_PORT="30432"
ENV DB_USER="postgres"
ENV DB_PASSWORD="ioasit123456"
ENV DB_NAME="ioa"
ENV DB_SSLMODE="disable"

# 暴露端口
EXPOSE 8090

# 创建logs目录
RUN mkdir -p /app/logs

# 运行应用
CMD ["/bin/sh", "-c", "./star_llm_backend > /app/logs/server.log 2>&1"]