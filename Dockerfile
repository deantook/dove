# 使用官方Golang镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 使用阿里云镜像源并安装必要的构建工具
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache git ca-certificates tzdata

# 设置环境变量
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,https://goproxy.io,direct \
    GOSUMDB=sum.golang.google.cn \
    GO111MODULE=on

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖（利用Docker缓存层）
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（使用模块路径，添加构建参数）
RUN go build \
    -a \
    -installsuffix cgo \
    -ldflags="-w -s -extldflags '-static'" \
    -o dove cmd/main.go

# 使用轻量级的Alpine镜像作为运行时环境
FROM alpine:latest

# 使用阿里云镜像源并安装必要的运行时依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件和配置文件
COPY --from=builder /app/dove .
COPY --from=builder /app/config ./config

# 更改文件所有者
RUN chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用
CMD ["./dove"]