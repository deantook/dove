# 开发环境 Dockerfile
FROM golang:1.21-alpine

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖和开发工具
RUN apk add --no-cache git ca-certificates tzdata \
    && go install github.com/cosmtrek/air@latest

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 暴露端口
EXPOSE 8080

# 设置环境变量
ENV GIN_MODE=debug
ENV ENV=dev

# 使用 air 进行热重载开发
CMD ["air"]
