# 使用官方Golang镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go模块文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用轻量级alpine镜像作为运行环境
FROM alpine:latest

# 安装ca-certificates，用于HTTPS请求
RUN apk --no-cache add ca-certificates

# 创建非root用户
RUN adduser -D -s /bin/sh appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 更改文件所有者
RUN chown -R appuser:appuser /app
USER appuser

# 暴露端口
EXPOSE 8888

# 启动应用
CMD ["./main"]