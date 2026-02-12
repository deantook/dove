---
name: swagger-annotations
description: Guides writing Swagger/OpenAPI annotations with swaggo and generating docs for Go Gin APIs. Use when adding or updating Swagger comments, generating swagger documentation, or working with swaggo/swag and gin-swagger.
---

# Swagger 注释与文档生成

基于 **swaggo/swag** 与 **gin-swagger**，为 Go Gin 接口编写声明式注释并生成 Swagger 文档。

## 何时使用

- 为现有 handler 补充或修改 Swagger 注释
- 新增接口时一并编写 `@Summary`、`@Router` 等
- 需要执行 `swag init` 重新生成 docs
- 接入或排查 gin-swagger 路由、导入问题

## 1. 入口：main 包总览

在 **main 函数所在文件**（如 `cmd/server/main.go`）顶部、`func main()` 上方添加：

```go
// @title           API 名称
// @version         1.0
// @description     简短描述
// @host            localhost:8080
// @BasePath        /
func main() {
```

并确保：
- 注册 Swagger 路由：`router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))`
- 导入：`_ "your-module/docs"`、`swaggerFiles "github.com/swaggo/files"`、`ginSwagger "github.com/swaggo/gin-swagger"`

## 2. Handler 注释（每个接口）

在 **handler 函数正上方** 写注释块，顺序建议：Summary → Description → Tags → Accept/Produce → Param → Success/Failure → Router。

### 必选

| 注解        | 说明           | 示例 |
|-------------|----------------|------|
| `@Summary`  | 短标题         | `@Summary 创建用户` |
| `@Router`   | 路径 + 方法     | `@Router /api/v1/users [post]` |
| `@Tags`     | 分组（如资源名）| `@Tags users` |

### 常用

| 注解        | 说明           | 示例 |
|-------------|----------------|------|
| `@Description` | 详细说明     | `@Description 根据用户名创建新用户` |
| `@Accept`   | 请求体类型     | `@Accept json` |
| `@Produce`  | 响应类型      | `@Produce json` |
| `@Param`    | 参数（见下）   | 见下方「Param 写法」 |
| `@Success`  | 成功响应      | `@Success 200 {object} handler.UserResponse` |
| `@Failure`  | 错误响应      | `@Failure 400 {object} object` |

### Param 写法

格式：`@Param 名称 位置 类型 必填 "说明"`

- **path**：路径参数，如 `@Param id path int true "用户 ID"`
- **query**：查询参数，如 `@Param page query int false "页码" default(1)`
- **body**：请求体，如 `@Param body body handler.CreateUserRequest true "请求体"`

同一 handler 可写多行 `@Param`。

### Router 路径

- 使用 **完整路径**（含前缀），与 Gin 注册一致，例如 `/api/v1/users`、`/api/v1/users/{id}`。
- 路径参数用 **`{id}`**，不要用 `:id`。

## 3. 示例片段

**POST 带 body：**

```go
// CreateUser 创建用户
// @Summary      创建用户
// @Description  根据用户名创建新用户
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        body  body  handler.CreateUserRequest  true  "创建用户请求"
// @Success      200   {object}  object  "success, message, data"
// @Failure      400   {object}  object  "success, message"
// @Router       /api/v1/users [post]
func CreateUser(c *gin.Context) {
```

**GET path + query：**

```go
// GetUser 获取用户详情
// @Summary      获取用户详情
// @Tags         users
// @Produce      json
// @Param        id   path   int  true  "用户 ID"
// @Success      200  {object}  object  "success, data"
// @Failure      404  {object}  object  "success, message"
// @Router       /api/v1/users/{id} [get]
func GetUser(c *gin.Context) {
```

```go
// ListUsers 用户列表
// @Summary      用户列表
// @Tags         users
// @Produce      json
// @Param        page       query  int  false  "页码"      default(1)
// @Param        page_size  query  int  false  "每页条数"  default(10)
// @Success      200  {object}  handler.ListUsersResponse
// @Router       /api/v1/users [get]
func ListUsers(c *gin.Context) {
```

**响应类型**：同一包内可直接写结构体名（如 `handler.CreateUserRequest`）；泛型响应可用 `object`。

## 4. 生成文档

修改注释或路由后需重新生成：

```bash
swag init -g cmd/server/main.go --parseDependency --parseInternal
```

未安装 swag 时可用：

```bash
go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/server/main.go --parseDependency --parseInternal
```

- `-g`：指定 main 入口文件。
- `--parseDependency` / `--parseInternal`：解析依赖与 internal 包中的 handler，否则不会扫描到 handler 上的注释。

生成结果在项目根目录 **`docs/`**：`docs.go`、`swagger.json`、`swagger.yaml`。不要手改这些文件。

## 5. 检查清单

- [ ] main 包有 `@title`、`@version`、`@host`、`@BasePath`
- [ ] 每个接口有 `@Summary`、`@Tags`、`@Router`，路径与 Gin 一致且 path 用 `{id}`
- [ ] 有 body 的接口写了 `@Param ... body ...` 和 `@Accept json`
- [ ] 生成后已执行 `go build ./...` 确保能编译（docs 包被正确引用）

## 更多注解

完整注解列表与可选字段见 [reference.md](reference.md)。
