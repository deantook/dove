# dove - Go Web API 框架

一个基于 Go 的现代化 Web API 框架，采用 Clean Architecture 架构设计，集成了完整的开发工具链。

## 特性

- 🏗️ **Clean Architecture**: 清晰的分层架构，易于维护和扩展
- 🔧 **代码生成**: 自动生成 CRUD 代码，提高开发效率
- 📝 **完整日志**: 结构化日志记录，支持链路追踪
- 🔐 **JWT 认证**: 内置 JWT 认证机制
- 📊 **分页支持**: 统一的分页和搜索接口
- 🗄️ **数据库支持**: GORM + MySQL/PostgreSQL
- 🚀 **高性能**: 基于 Gin 框架，性能优异
- 📚 **API 文档**: 自动生成 Swagger 文档
- 🧪 **测试支持**: 完整的测试框架
- 🔄 **依赖注入**: Wire 依赖注入
- 🐳 **Docker 支持**: 完整的 Docker 部署方案

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd dove
```

### 2. 安装依赖

```bash
go mod download
```

### 3. 配置环境

复制配置文件：

```bash
cp config/dev.yaml config/local.yaml
```

编辑 `config/local.yaml` 文件，配置数据库连接等信息。

### 4. 运行项目

#### 方式一：直接运行

```bash
make run
```

或者直接运行：

```bash
go run main.go
```

#### 方式二：使用 Docker（推荐）

```bash
# 开发环境（包含热重载）
make docker-dev

# 生产环境
make docker-prod
```

## Docker 部署

### 快速开始

```bash
# 开发环境
make docker-dev

# 生产环境
make docker-prod

# 查看容器状态
make docker-status

# 查看日志
make docker-logs
```

### 详细配置

- [Docker 部署指南](docs/docker_deployment.md)
- 支持开发环境和生产环境
- 包含 MySQL、Redis、Nginx 等完整服务栈
- 支持热重载和自动重启

## CRUD 代码生成器

### 概述

CRUD 代码生成器是一个自动化工具，可以根据模型文件自动生成完整的 CRUD（创建、读取、更新、删除）代码，包括：

- Domain 层（接口定义）
- Service 层（业务逻辑）
- Repository 层（数据访问）
- Handler 层（HTTP 接口）
- Wire 依赖注入配置

### 使用方法

#### 1. 创建模型文件

在 `internal/model/` 目录下创建你的模型文件，例如 `article.go`：

```go
package model

import (
	"time"
	"gorm.io/gorm"
)

// Article 文章模型
type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null;size:200"`
	Content     string         `json:"content" gorm:"type:text"`
	Author      string         `json:"author" gorm:"size:100"`
	Status      int            `json:"status" gorm:"default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
```

#### 2. 生成 CRUD 代码

使用 Makefile 命令：

```bash
# 生成 Article 的 CRUD 代码
make generate-article

# 生成自定义模型的 CRUD 代码
make generate-crud-MODEL=internal/model/your_model.go
```

或者直接使用命令行：

```bash
go run cmd/generate/main.go -model internal/model/article.go
```

#### 3. 生成的文件

生成器会创建以下文件：

- `internal/domain/article.go` - Domain 层接口定义
- `internal/service/article_service.go` - Service 层业务逻辑
- `internal/repository/article_repository.go` - Repository 层数据访问
- `internal/handler/article_handler.go` - Handler 层 HTTP 接口
- `internal/wire/providers.go` - Wire 依赖注入配置

### 功能特性

- **自动字段解析**: 解析模型文件中的字段信息
- **智能搜索和排序**: 根据字段名称自动识别可搜索和可排序字段
- **完整的 CRUD 操作**: 包含创建、读取、更新、删除功能
- **分页和搜索**: 支持分页查询、字段搜索、多字段排序
- **错误处理**: 完整的错误处理和日志记录
- **API 文档**: 自动生成 Swagger 注释

### 示例

生成 Article 的 CRUD 代码：

```bash
make generate-article
```

输出：
```
解析模型: Article
字段数量: 11
可搜索字段: [Title]
可排序字段: [ID Status]
✅ CRUD代码生成完成!
输出目录: internal

生成的文件:
  - internal/domain/article.go
  - internal/service/article_service.go
  - internal/repository/article_repository.go
  - internal/handler/article_handler.go

注意: 请检查生成的代码并根据需要进行调整
```

生成的 API 接口：

- `POST /articles/` - 创建文章
- `GET /articles/{id}` - 获取文章详情
- `GET /articles/` - 获取文章列表（支持分页、搜索、排序）
- `PUT /articles/{id}` - 更新文章
- `DELETE /articles/{id}` - 删除文章

## 项目结构

```
dove/
├── cmd/                    # 命令行工具
│   ├── generate/          # CRUD 代码生成器
│   └── migrate/           # 数据库迁移工具
├── config/                # 配置文件
├── docs/                  # 文档
├── examples/              # 示例代码
├── internal/              # 内部代码
│   ├── app/              # 应用层
│   ├── domain/           # 领域层
│   ├── handler/          # 处理器层
│   ├── middleware/       # 中间件
│   ├── model/            # 数据模型
│   ├── repository/       # 仓储层
│   ├── service/          # 服务层
│   └── wire/             # 依赖注入
├── nginx/                # Nginx 配置
├── pkg/                  # 公共包
│   ├── config/           # 配置管理
│   ├── database/         # 数据库
│   ├── generator/        # 代码生成器
│   ├── jwt/              # JWT 工具
│   ├── logger/           # 日志工具
│   ├── pagination/       # 分页工具
│   ├── password/         # 密码工具
│   ├── redis/            # Redis 工具
│   └── response/         # 响应工具
├── scripts/              # 脚本文件
├── Dockerfile            # 生产环境 Dockerfile
├── Dockerfile.dev        # 开发环境 Dockerfile
├── docker-compose.yml    # 开发环境 Docker Compose
├── docker-compose.prod.yml # 生产环境 Docker Compose
├── .dockerignore         # Docker 忽略文件
├── .air.toml            # 开发环境热重载配置
├── main.go               # 主程序
├── Makefile              # 构建脚本
└── README.md             # 项目说明
```

## 开发指南

### 添加新的模型

1. 在 `internal/model/` 目录下创建模型文件
2. 使用 CRUD 生成器生成代码
3. 在路由中注册新的 handler
4. 运行测试确保功能正常

### 自定义生成器

可以修改以下文件来自定义生成器行为：

- `pkg/generator/generator.go` - 生成器核心逻辑
- `pkg/generator/templates/` - 模板文件

### 测试

```bash
# 运行所有测试
make test

# 运行特定包的测试
go test ./internal/...

# 运行测试并生成覆盖率报告
make test-coverage
```

## 部署

### 本地部署

```bash
# 构建
make build

# 运行
make run
```

### Docker 部署

```bash
# 开发环境
make docker-dev

# 生产环境
make docker-prod

# 查看状态
make docker-status
```

### 生产环境部署

1. **使用 Docker Compose**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

2. **使用 Docker Swarm**
   ```bash
   docker stack deploy -c docker-compose.prod.yml dove
   ```

3. **使用 Kubernetes**
   - 参考 `k8s/` 目录下的配置文件

## 常用命令

### 开发命令

```bash
# 运行开发环境
make run

# 运行测试
make test

# 格式化代码
make fmt

# 代码检查
make lint
```

### Docker 命令

```bash
# 开发环境
make docker-dev

# 生产环境
make docker-prod

# 查看日志
make docker-logs

# 进入容器
make docker-exec
```

### 代码生成命令

```bash
# 生成 CRUD 代码
make generate-article

# 生成 Swagger 文档
make swagger

# 运行 Wire
make wire
```

## 贡献

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

如有问题或建议，请通过以下方式联系：

- 提交 Issue
- 发送邮件
- 参与讨论

---

**注意**: 这是一个开发中的项目，API 可能会发生变化。请在生产环境中使用前进行充分测试。 

## 待实现功能

以下是项目计划实现的功能列表：

### 🔐 安全与认证
- [ ] **OAuth2.0 支持** - 集成第三方登录（Google、GitHub、微信等）
- [ ] **多因素认证 (MFA)** - 支持短信、邮箱、TOTP 等验证方式
- [ ] **角色权限管理 (RBAC)** - 细粒度的权限控制系统
- [ ] **API 限流** - 基于 IP、用户、接口的限流机制
- [ ] **API 密钥管理** - 支持 API Key 认证
- [ ] **会话管理** - 多设备登录、会话超时、强制下线
- [ ] **密码策略** - 密码复杂度、定期更换、历史密码检查

### 📊 数据与存储
- [ ] **文件上传** - 支持图片、文档等文件上传和管理
- [ ] **数据备份** - 自动数据库备份和恢复
- [ ] **数据导入导出** - CSV、Excel 等格式的数据导入导出
- [ ] **数据版本控制** - 记录数据变更历史
- [ ] **软删除优化** - 统一的软删除机制和恢复功能
- [ ] **数据加密** - 敏感数据加密存储
- [ ] **数据脱敏** - 敏感数据脱敏显示

### 🔍 搜索与查询
- [ ] **全文搜索** - 集成 Elasticsearch 或 Meilisearch
- [ ] **高级搜索** - 多条件组合搜索
- [ ] **搜索建议** - 智能搜索建议和自动补全
- [ ] **搜索历史** - 用户搜索历史记录
- [ ] **搜索结果高亮** - 搜索结果关键词高亮显示

### 📈 监控与日志
- [ ] **性能监控** - 接口响应时间、吞吐量监控
- [ ] **错误追踪** - 集成 Sentry 等错误追踪服务
- [ ] **健康检查增强** - 详细的系统健康状态检查
- [ ] **日志分析** - 日志聚合和分析工具
- [ ] **告警系统** - 异常情况自动告警
- [ ] **指标收集** - Prometheus 指标收集

### 🚀 性能与缓存
- [ ] **缓存策略** - Redis 缓存策略优化
- [ ] **CDN 集成** - 静态资源 CDN 加速
- [ ] **数据库连接池** - 连接池优化和监控
- [ ] **异步处理** - 消息队列集成（RabbitMQ、Kafka）
- [ ] **批量操作** - 批量创建、更新、删除接口
- [ ] **数据预加载** - 关联数据预加载优化

### 🔄 工作流与业务
- [ ] **工作流引擎** - 业务流程自动化
- [ ] **定时任务** - Cron 任务调度系统
- [ ] **消息通知** - 邮件、短信、推送通知
- [ ] **审批流程** - 多级审批流程
- [ ] **数据统计** - 数据统计和报表功能
- [ ] **数据可视化** - 图表和仪表板

### 🌐 集成与扩展
- [ ] **Webhook 支持** - 事件驱动的 Webhook 机制
- [ ] **第三方服务集成** - 支付、短信、邮件等服务
- [ ] **API 网关** - 统一的 API 网关
- [ ] **微服务架构** - 服务拆分和治理
- [ ] **服务发现** - 服务注册和发现
- [ ] **配置中心** - 动态配置管理

### 🧪 测试与质量
- [ ] **集成测试** - 完整的集成测试套件
- [ ] **性能测试** - 压力测试和性能基准
- [ ] **安全测试** - 安全漏洞扫描和测试
- [ ] **代码质量** - 代码质量检查和静态分析
- [ ] **自动化测试** - CI/CD 自动化测试
- [ ] **测试覆盖率** - 提高测试覆盖率

### 📱 用户体验
- [ ] **API 版本控制** - 多版本 API 支持
- [ ] **API 文档增强** - 交互式 API 文档
- [ ] **SDK 生成** - 自动生成客户端 SDK
- [ ] **API 沙箱** - 在线 API 测试环境
- [ ] **用户反馈** - 用户反馈和建议系统
- [ ] **帮助文档** - 详细的用户帮助文档

### 🛠️ 开发工具
- [ ] **代码生成器增强** - 更多模板和自定义选项
- [ ] **开发环境优化** - 本地开发环境一键搭建
- [ ] **调试工具** - 更好的调试和诊断工具
- [ ] **代码规范** - 代码规范和格式化工具
- [ ] **依赖管理** - 依赖版本管理和更新
- [ ] **项目模板** - 快速创建新项目的模板

### 🚢 部署与运维
- [ ] **Kubernetes 支持** - 完整的 K8s 部署方案
- [ ] **服务网格** - Istio 服务网格集成
- [ ] **蓝绿部署** - 零停机部署策略
- [ ] **回滚机制** - 快速回滚和版本管理
- [ ] **环境管理** - 多环境配置管理
- [ ] **容器编排** - Docker Swarm 支持

### 🔒 合规与审计
- [ ] **审计日志** - 完整的操作审计日志
- [ ] **数据隐私** - GDPR 等隐私法规支持
- [ ] **合规报告** - 自动生成合规报告
- [ ] **数据保留** - 数据保留策略
- [ ] **访问控制** - 细粒度的访问控制
- [ ] **安全扫描** - 定期安全扫描

### 📊 数据分析
- [ ] **用户行为分析** - 用户行为数据收集和分析
- [ ] **业务指标** - 关键业务指标监控
- [ ] **数据挖掘** - 数据挖掘和机器学习
- [ ] **预测分析** - 基于历史数据的预测
- [ ] **实时分析** - 实时数据处理和分析
- [ ] **数据仓库** - 数据仓库和 ETL 流程

---

**注意**: 这些功能将根据项目需求和优先级逐步实现。欢迎贡献代码和提出建议！ 