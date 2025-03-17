# 构建阶段
FROM golang:1.21-alpine AS builder

# 设置Go环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /build

# 复制项目文件
COPY . .

# 下载依赖并构建
RUN go mod download
RUN go build -o ollama-proxy .

# 运行阶段
FROM alpine:latest

# 设置工作目录
WORKDIR /app

# 从builder阶段复制二进制文件和配置文件
COPY --from=builder /build/ollama-proxy /app/
COPY --from=builder /build/config.yaml /app/

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["/app/ollama-proxy"]