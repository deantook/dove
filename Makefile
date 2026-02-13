.PHONY: help swagger wire build run test clean install-tools

# 变量定义
APP_NAME := dove
MAIN_PATH := cmd/server/main.go
SWAGGER_MAIN := cmd/server/main.go
SWAGGER_OUTPUT := api/swagger
WIRE_DIR := wire
CONFIG_PATH := configs/config.yaml

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make swagger      - 生成 Swagger 文档"
	@echo "  make wire         - 生成 Wire 依赖注入代码"
	@echo "  make build        - 构建应用"
	@echo "  make run          - 运行应用"
	@echo "  make test         - 运行测试"
	@echo "  make clean        - 清理构建文件"
	@echo "  make install-tools - 安装开发工具 (swag, wire)"
	@echo "  make fmt          - 格式化代码"
	@echo "  make lint         - 代码检查"

# 安装开发工具
install-tools:
	@echo "安装 Swagger 工具..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "安装 Wire 工具..."
	@go install github.com/google/wire/cmd/wire@latest
	@echo "工具安装完成"

# 生成 Swagger 文档
swagger:
	@echo "生成 Swagger 文档..."
	@swag init -g $(SWAGGER_MAIN) -o $(SWAGGER_OUTPUT)
	@echo "Swagger 文档生成完成: $(SWAGGER_OUTPUT)"

# 生成 Wire 代码
wire:
	@echo "生成 Wire 依赖注入代码..."
	@cd $(WIRE_DIR) && wire
	@echo "Wire 代码生成完成: $(WIRE_DIR)/wire_gen.go"

# 格式化代码
fmt:
	@echo "格式化代码..."
	@go fmt ./...
	@echo "代码格式化完成"

# 代码检查
lint:
	@echo "代码检查..."
	@go vet ./...
	@echo "代码检查完成"

# 构建应用
build:
	@echo "构建应用..."
	@go build -o $(APP_NAME) $(MAIN_PATH)
	@echo "构建完成: $(APP_NAME)"

# 运行应用
run:
	@echo "运行应用..."
	@go run $(MAIN_PATH) $(CONFIG_PATH)

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -f $(APP_NAME)
	@echo "清理完成"

# 完整构建流程（生成文档和代码后构建）
all: swagger wire build
	@echo "完整构建流程完成"
