# CRUD代码生成器

## 概述

CRUD代码生成器是一个自动化工具，可以根据模型文件自动生成完整的CRUD（创建、读取、更新、删除）代码，包括：

- Domain层（接口定义）
- Service层（业务逻辑）
- Repository层（数据访问）
- Handler层（HTTP接口）
- Wire依赖注入配置

## 使用方法

### 1. 创建模型文件

首先在 `internal/model/` 目录下创建你的模型文件，例如 `article.go`：

```go
package model

import (
	"time"
	"gorm.io/gorm"
)

// Article 文章模型
// @Description 文章信息
type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey" example:"1"`
	Title       string         `json:"title" gorm:"not null;size:200" example:"这是一篇测试文章"`
	Content     string         `json:"content" gorm:"type:text" example:"文章内容..."`
	Author      string         `json:"author" gorm:"size:100" example:"张三"`
	Category    string         `json:"category" gorm:"size:50" example:"技术"`
	Tags        string         `json:"tags" gorm:"size:200" example:"golang,web"`
	Status      int            `json:"status" gorm:"default:1" example:"1"` // 1: 发布, 0: 草稿
	ViewCount   int            `json:"view_count" gorm:"default:0" example:"100"`
	CreatedAt   time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt   time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`
}
```

### 2. 运行生成器

使用以下命令生成CRUD代码：

```bash
go run cmd/generate/main.go -model internal/model/article.go
```

### 3. 参数说明

- `-model`: 模型文件路径（必需）
- `-output`: 输出目录（可选，默认为 internal）
- `-help`: 显示帮助信息

### 4. 生成的文件

生成器会创建以下文件：

- `internal/domain/article.go` - Domain层接口定义
- `internal/service/article_service.go` - Service层业务逻辑
- `internal/repository/article_repository.go` - Repository层数据访问
- `internal/handler/article_handler.go` - Handler层HTTP接口
- `internal/wire/providers.go` - Wire依赖注入配置

## 功能特性

### 自动字段解析

生成器会自动解析模型文件中的字段信息：

- 字段名称和类型
- JSON标签
- GORM标签
- 字段注释

### 智能搜索和排序

生成器会根据字段名称自动识别：

- **可搜索字段**: 包含 `name`, `title`, `description`, `username`, `email`, `nickname` 的字段
- **可排序字段**: 包含 `id`, `name`, `created_at`, `updated_at`, `price`, `stock`, `status` 的字段

### 完整的CRUD操作

生成的代码包含：

1. **创建 (Create)**
   - 验证请求数据
   - 创建记录
   - 返回创建结果

2. **读取 (Read)**
   - 根据ID获取单条记录
   - 获取所有记录（支持分页）
   - 支持搜索和排序

3. **更新 (Update)**
   - 验证记录存在
   - 部分更新字段
   - 返回更新结果

4. **删除 (Delete)**
   - 验证记录存在
   - 软删除（如果支持）
   - 返回删除结果

### 分页和搜索

生成的代码支持：

- 分页查询
- 字段搜索
- 多字段排序
- 参数验证

### 错误处理

包含完整的错误处理：

- 数据库错误
- 验证错误
- 记录不存在错误
- 业务逻辑错误

### 日志记录

集成了结构化日志：

- 操作日志
- 错误日志
- 性能监控

## 自定义配置

### 修改搜索字段

在 `pkg/generator/generator.go` 中修改 `isSearchableField` 函数：

```go
func isSearchableField(fieldName string) bool {
	searchableFields := []string{"name", "title", "description", "username", "email", "nickname"}
	for _, field := range searchableFields {
		if strings.Contains(strings.ToLower(fieldName), field) {
			return true
		}
	}
	return false
}
```

### 修改排序字段

在 `pkg/generator/generator.go` 中修改 `isSortableField` 函数：

```go
func isSortableField(fieldName string) bool {
	sortableFields := []string{"id", "name", "created_at", "updated_at", "price", "stock", "status"}
	for _, field := range sortableFields {
		if strings.Contains(strings.ToLower(fieldName), field) {
			return true
		}
	}
	return false
}
```

## 模板文件

生成器使用以下模板文件：

- `pkg/generator/templates/domain.go.tmpl` - Domain层模板
- `pkg/generator/templates/service.go.tmpl` - Service层模板
- `pkg/generator/templates/repository.go.tmpl` - Repository层模板
- `pkg/generator/templates/handler.go.tmpl` - Handler层模板

## 注意事项

1. **检查生成的代码**: 生成后请检查代码并根据需要进行调整
2. **添加业务逻辑**: 生成的代码是基础框架，需要根据具体业务需求添加逻辑
3. **测试**: 生成代码后请进行充分测试
4. **路由配置**: 需要在路由中注册生成的handler
5. **数据库迁移**: 确保数据库表结构已创建

## 示例

### 生成Article的CRUD代码

```bash
go run cmd/generate/main.go -model internal/model/article.go
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

### 使用生成的API

生成后可以直接使用以下API：

- `POST /articles/` - 创建文章
- `GET /articles/{id}` - 获取文章详情
- `GET /articles/` - 获取文章列表（支持分页、搜索、排序）
- `PUT /articles/{id}` - 更新文章
- `DELETE /articles/{id}` - 删除文章

## 扩展功能

生成器支持以下扩展：

1. **自定义验证**: 在Service层添加业务验证逻辑
2. **权限控制**: 在Handler层添加权限检查
3. **缓存支持**: 在Repository层添加缓存逻辑
4. **事件系统**: 在Service层添加事件发布
5. **审计日志**: 添加操作审计功能

## 故障排除

### 常见问题

1. **模板解析错误**: 检查模板文件语法
2. **字段解析失败**: 确保模型文件格式正确
3. **Wire配置错误**: 检查依赖注入配置
4. **编译错误**: 检查生成的代码语法

### 调试技巧

1. 使用 `-help` 参数查看帮助
2. 检查生成的代码语法
3. 运行 `go mod tidy` 更新依赖
4. 使用 `go build` 检查编译错误
