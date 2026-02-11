# 项目名称
PROJECT_NAME := dove

# Go 相关变量
GO := go
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# 构建相关变量
BUILD_DIR := build
BINARY_NAME := $(PROJECT_NAME)
BINARY_UNIX := $(BINARY_NAME)_unix

# 环境变量
ENV ?= dev

# 默认目标
.DEFAULT_GOAL := help

# 帮助信息
.PHONY: help
help: ## 显示帮助信息
	@echo "可用的命令:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# 构建相关
.PHONY: build
build: ## 构建应用
	@echo "构建应用..."
	@mkdir -p $(BUILD_DIR)
	@$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) .

.PHONY: build-linux
build-linux: ## 构建 Linux 版本
	@echo "构建 Linux 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_UNIX) .

.PHONY: build-windows
build-windows: ## 构建 Windows 版本
	@echo "构建 Windows 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME).exe .

.PHONY: build-mac
build-mac: ## 构建 macOS 版本
	@echo "构建 macOS 版本..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)_darwin .

# 运行相关
.PHONY: run
run: ## 运行应用 (开发环境)
	@echo "运行应用 (环境: $(ENV))..."
	@ENV=$(ENV) $(GO) run main.go

.PHONY: run-prod
run-prod: ## 运行应用 (生产环境)
	@echo "运行应用 (生产环境)..."
	@ENV=production $(GO) run main.go

.PHONY: run-test
run-test: ## 运行应用 (测试环境)
	@echo "运行应用 (测试环境)..."
	@ENV=test $(GO) run main.go

# 测试相关
.PHONY: test
test: ## 运行测试
	@echo "运行测试..."
	@$(GO) test -v ./...

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "运行测试并生成覆盖率报告..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 代码质量
.PHONY: lint
lint: ## 运行代码检查
	@echo "运行代码检查..."
	@$(GO) vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，跳过 lint 检查"; \
	fi

.PHONY: fmt
fmt: ## 格式化代码
	@echo "格式化代码..."
	@$(GO) fmt ./...

.PHONY: tidy
tidy: ## 整理 go.mod 依赖
	@echo "整理依赖..."
	@$(GO) mod tidy

# 数据库相关
.PHONY: migrate
migrate: ## 执行数据库迁移
	@echo "执行数据库迁移..."
	@$(GO) run cmd/migrate/main.go -action=migrate -env=$(ENV)

.PHONY: migrate-reset
migrate-reset: ## 重置数据库
	@echo "重置数据库..."
	@$(GO) run cmd/migrate/main.go -action=reset -env=$(ENV)

.PHONY: migrate-drop
migrate-drop: ## 删除所有表
	@echo "删除所有表..."
	@$(GO) run cmd/migrate/main.go -action=drop -env=$(ENV)

# Swagger 相关
.PHONY: swagger
swagger: ## 生成 Swagger 文档
	@echo "生成 Swagger 文档..."
	@swag init

# Wire 相关
.PHONY: wire
wire: ## 生成 Wire 依赖注入代码
	@echo "生成 Wire 依赖注入代码..."
	@wire ./internal/wire

# 清理相关
.PHONY: clean
clean: ## 清理构建文件
	@echo "清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

.PHONY: clean-all
clean-all: clean ## 清理所有生成的文件
	@echo "清理所有生成的文件..."
	@rm -rf docs/
	@rm -f internal/wire/wire_gen.go

# 安装工具
.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "安装开发工具..."
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@$(GO) install github.com/google/wire/cmd/wire@latest
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 开发相关
.PHONY: dev
dev: tidy wire swagger run ## 开发模式：整理依赖、生成代码、运行应用

.PHONY: dev-setup
dev-setup: install-tools tidy wire swagger ## 开发环境设置

# 部署相关
.PHONY: deploy-prepare
deploy-prepare: clean-all build ## 部署准备：清理、构建

# Docker 相关
.PHONY: docker-build docker-build-dev docker-build-prod docker-build-test docker-build-all docker-run docker-dev docker-prod docker-clean

# Docker 镜像标签和版本
DOCKER_REGISTRY ?= 
DOCKER_IMAGE_NAME ?= $(PROJECT_NAME)
DOCKER_TAG ?= latest
DOCKER_VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

# 构建所有环境的 Docker 镜像
docker-build-all: docker-build-dev docker-build-test docker-build-prod ## 构建所有环境的 Docker 镜像
	@echo "✅ 所有环境的 Docker 镜像构建完成"

# 构建开发环境 Docker 镜像
docker-build-dev: ## 构建开发环境 Docker 镜像
	@echo "🔨 构建开发环境 Docker 镜像..."
	@docker build \
		-f Dockerfile \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):dev \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):dev-$(DOCKER_VERSION) \
		--build-arg ENV=dev \
		--build-arg GIN_MODE=debug \
		.
	@echo "✅ 开发环境镜像构建完成: $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):dev"

# 构建测试环境 Docker 镜像
docker-build-test: ## 构建测试环境 Docker 镜像
	@echo "🔨 构建测试环境 Docker 镜像..."
	@docker build \
		-f Dockerfile \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):test \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):test-$(DOCKER_VERSION) \
		--build-arg ENV=test \
		--build-arg GIN_MODE=release \
		.
	@echo "✅ 测试环境镜像构建完成: $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):test"

# 构建生产环境 Docker 镜像
docker-build-prod: ## 构建生产环境 Docker 镜像
	@echo "🔨 构建生产环境 Docker 镜像..."
	@docker build \
		-f Dockerfile \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):latest \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):$(DOCKER_VERSION) \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):prod \
		--build-arg ENV=production \
		--build-arg GIN_MODE=release \
		.
	@echo "✅ 生产环境镜像构建完成: $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):latest"

# 构建指定环境的 Docker 镜像
docker-build: ## 构建当前环境的 Docker 镜像 (默认: dev)
	@echo "🔨 构建 $(ENV) 环境 Docker 镜像..."
	@if [ "$(ENV)" = "production" ]; then \
		$(MAKE) docker-build-prod; \
	elif [ "$(ENV)" = "test" ]; then \
		$(MAKE) docker-build-test; \
	else \
		$(MAKE) docker-build-dev; \
	fi

# 构建并推送镜像到仓库
docker-build-push: docker-build-all ## 构建并推送所有镜像到仓库
	@echo "🚀 推送镜像到仓库..."
	@if [ -n "$(DOCKER_REGISTRY)" ]; then \
		docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):dev; \
		docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):dev-$(DOCKER_VERSION); \
		docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):test; \
		docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):test-$(DOCKER_VERSION); \
		docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):latest; \
		docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):$(DOCKER_VERSION); \
		docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):prod; \
		echo "✅ 所有镜像推送完成"; \
	else \
		echo "⚠️  DOCKER_REGISTRY 未设置，跳过推送"; \
	fi

# 构建并推送指定环境镜像
docker-build-push-$(ENV): docker-build ## 构建并推送指定环境镜像
	@echo "🚀 推送 $(ENV) 环境镜像到仓库..."
	@if [ -n "$(DOCKER_REGISTRY)" ]; then \
		if [ "$(ENV)" = "production" ]; then \
			docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):latest; \
			docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):$(DOCKER_VERSION); \
			docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):prod; \
		elif [ "$(ENV)" = "test" ]; then \
			docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):test; \
			docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):test-$(DOCKER_VERSION); \
		else \
			docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):dev; \
			docker push $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):dev-$(DOCKER_VERSION); \
		fi; \
		echo "✅ $(ENV) 环境镜像推送完成"; \
	else \
		echo "⚠️  DOCKER_REGISTRY 未设置，跳过推送"; \
	fi

# 构建多架构镜像 (需要 Docker Buildx)
docker-build-multiarch: ## 构建多架构 Docker 镜像 (linux/amd64, linux/arm64)
	@echo "🔨 构建多架构 Docker 镜像..."
	@docker buildx create --use --name multiarch-builder || true
	@docker buildx build \
		--platform linux/amd64,linux/arm64 \
		-f Dockerfile \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):latest \
		-t $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):$(DOCKER_VERSION) \
		--build-arg ENV=production \
		--build-arg GIN_MODE=release \
		--push \
		.
	@echo "✅ 多架构镜像构建完成"

# 构建并运行 Docker 容器
docker-build-run: docker-build ## 构建并运行 Docker 容器
	@echo "🚀 运行 Docker 容器..."
	@docker rm -f $(PROJECT_NAME)-$(ENV) || true
	@docker run -d --name $(PROJECT_NAME)-$(ENV) -p 8080:8080 $(DOCKER_REGISTRY)$(DOCKER_IMAGE_NAME):$(ENV)

# 运行开发环境容器
docker-dev: ## 运行开发环境 Docker 容器
	@echo "🚀 运行开发环境 Docker 容器..."
	@docker-compose --profile dev up -d

# 运行测试环境容器
docker-test: ## 运行测试环境 Docker 容器
	@echo "🚀 运行测试环境 Docker 容器..."
	@docker-compose -f docker-compose.test.yml up -d

# 运行生产环境容器
docker-prod: ## 运行生产环境 Docker 容器
	@echo "🚀 运行生产环境 Docker 容器..."
	@docker-compose -f docker-compose.prod.yml up -d

# 停止所有容器
docker-stop: ## 停止所有容器
	@echo "🛑 停止所有容器..."
	@docker-compose down
	@docker-compose -f docker-compose.prod.yml down
	@docker-compose -f docker-compose.test.yml down 2>/dev/null || true

# 清理 Docker 资源
docker-clean: ## 清理 Docker 资源
	@echo "🧹 清理 Docker 资源..."
	@docker-compose down -v
	@docker-compose -f docker-compose.prod.yml down -v
	@docker-compose -f docker-compose.test.yml down -v 2>/dev/null || true
	@docker system prune -f

# 查看容器日志
docker-logs: ## 查看应用容器日志
	@echo "📋 查看应用容器日志..."
	@docker-compose logs -f app

# 查看开发环境日志
docker-logs-dev: ## 查看开发环境日志
	@echo "📋 查看开发环境日志..."
	@docker-compose logs -f app-dev

# 查看测试环境日志
docker-logs-test: ## 查看测试环境日志
	@echo "📋 查看测试环境日志..."
	@docker-compose -f docker-compose.test.yml logs -f app

# 查看生产环境日志
docker-logs-prod: ## 查看生产环境日志
	@echo "📋 查看生产环境日志..."
	@docker-compose -f docker-compose.prod.yml logs -f app

# 进入容器
docker-exec: ## 进入应用容器
	@echo "🔍 进入应用容器..."
	@docker-compose exec app sh

# 进入开发容器
docker-exec-dev: ## 进入开发容器
	@echo "🔍 进入开发容器..."
	@docker-compose exec app-dev sh

# 进入测试容器
docker-exec-test: ## 进入测试容器
	@echo "🔍 进入测试容器..."
	@docker-compose -f docker-compose.test.yml exec app sh

# 进入生产容器
docker-exec-prod: ## 进入生产容器
	@echo "🔍 进入生产容器..."
	@docker-compose -f docker-compose.prod.yml exec app sh

# 重新构建并运行
docker-rebuild: ## 重新构建并运行 Docker 容器
	@echo "🔄 重新构建并运行 Docker 容器..."
	@docker-compose down
	@docker-compose build --no-cache
	@docker-compose up -d

# 重新构建开发环境
docker-rebuild-dev: ## 重新构建并运行开发环境
	@echo "🔄 重新构建并运行开发环境..."
	@docker-compose down
	@docker-compose build --no-cache app-dev
	@docker-compose --profile dev up -d

# 重新构建测试环境
docker-rebuild-test: ## 重新构建并运行测试环境
	@echo "🔄 重新构建并运行测试环境..."
	@docker-compose -f docker-compose.test.yml down
	@docker-compose -f docker-compose.test.yml build --no-cache
	@docker-compose -f docker-compose.test.yml up -d

# 重新构建生产环境
docker-rebuild-prod: ## 重新构建并运行生产环境
	@echo "🔄 重新构建并运行生产环境..."
	@docker-compose -f docker-compose.prod.yml down
	@docker-compose -f docker-compose.prod.yml build --no-cache
	@docker-compose -f docker-compose.prod.yml up -d

# 查看容器状态
docker-status: ## 查看容器状态
	@echo "📊 查看容器状态..."
	@docker-compose ps
	@echo ""
	@echo "生产环境容器状态:"
	@docker-compose -f docker-compose.prod.yml ps 2>/dev/null || echo "生产环境容器未运行"
	@echo ""
	@echo "测试环境容器状态:"
	@docker-compose -f docker-compose.test.yml ps 2>/dev/null || echo "测试环境容器未运行"

# 查看镜像列表
docker-images: ## 查看项目相关镜像
	@echo "📋 查看项目相关镜像..."
	@docker images | grep $(DOCKER_IMAGE_NAME) || echo "未找到相关镜像"

# 清理镜像
docker-clean-images: ## 清理项目相关镜像
	@echo "🧹 清理项目相关镜像..."
	@docker images | grep $(DOCKER_IMAGE_NAME) | awk '{print $$3}' | xargs -r docker rmi -f

# 备份数据库
docker-backup: ## 备份数据库
	@echo "💾 备份数据库..."
	@docker-compose exec mysql mysqldump -u root -proot123456 dove > backup_$(shell date +%Y%m%d_%H%M%S).sql

# 恢复数据库
docker-restore: ## 恢复数据库
	@echo "📥 恢复数据库..."
	@docker-compose exec -T mysql mysql -u root -proot123456 dove < $(BACKUP_FILE)

# 初始化数据库
docker-init-db: ## 初始化数据库
	@echo "🗄️ 初始化数据库..."
	@docker-compose exec mysql mysql -u root -proot123456 -e "CREATE DATABASE IF NOT EXISTS dove CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

# 运行测试（在容器中）
docker-test: ## 在容器中运行测试
	@echo "🧪 在容器中运行测试..."
	@docker-compose exec app go test ./...

# 生成 Swagger 文档（在容器中）
docker-swagger: ## 在容器中生成 Swagger 文档
	@echo "📚 在容器中生成 Swagger 文档..."
	@docker-compose exec app swag init

# 运行 Wire（在容器中）
docker-wire: ## 在容器中运行 Wire
	@echo "🔧 在容器中运行 Wire..."
	@docker-compose exec app wire ./internal/wire

# 版本信息
.PHONY: version
version: ## 显示版本信息
	@echo "项目: $(PROJECT_NAME)"
	@echo "Go 版本: $(shell go version)"
	@echo "操作系统: $(GOOS)"
	@echo "架构: $(GOARCH)"

# 依赖检查
.PHONY: deps
deps: ## 检查依赖
	@echo "检查依赖..."
	@$(GO) mod download
	@$(GO) mod verify 

# CRUD代码生成器
.PHONY: generate-crud
generate-crud:
	@echo "CRUD代码生成器"
	@echo "用法: make generate-crud MODEL=internal/model/your_model.go"
	@echo "示例: make generate-crud MODEL=internal/model/article.go"

generate-crud-%:
	@echo "生成CRUD代码: $*"
	@go run cmd/generate/main.go -model $*

# 生成Article的CRUD代码
generate-article:
	@echo "生成Article的CRUD代码..."
	@go run cmd/generate/main.go -model internal/model/article.go

# 生成User的CRUD代码
generate-user:
	@echo "生成User的CRUD代码..."
	@go run cmd/generate/main.go -model internal/model/user.go

# 生成Product的CRUD代码
generate-product:
	@echo "生成Product的CRUD代码..."
	@go run cmd/generate/main.go -model internal/model/product.go

# 清理生成的文件
clean-generated:
	@echo "清理生成的文件..."
	@rm -f internal/domain/article.go internal/service/article_service.go internal/repository/article_repository.go internal/handler/article_handler.go
	@rm -f internal/domain/user.go internal/service/user_service.go internal/repository/user_repository.go internal/handler/user_handler.go
	@rm -f internal/domain/product.go internal/service/product_service.go internal/repository/product_repository.go internal/handler/product_handler.go
	@rm -f internal/wire/providers.go

# 测试生成器
test-generator:
	@echo "测试CRUD代码生成器..."
	@go run cmd/generate/main.go -help
	@echo ""
	@echo "生成Article示例..."
	@make generate-article
	@echo ""
	@echo "编译测试..."
	@go build ./...
	@echo "✅ 生成器测试通过!" 