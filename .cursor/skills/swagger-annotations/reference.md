# Swaggo 注解速查

## 通用 API 信息（main 包）

| 注解 | 说明 | 示例 |
|------|------|------|
| `@title` | API 标题 | `@title Dove API` |
| `@version` | 版本 | `@version 1.0` |
| `@description` | 描述 | `@description Dove 服务 API` |
| `@host` | 主机:端口 | `@host localhost:8080` |
| `@BasePath` | 基础路径 | `@BasePath /` |
| `@schemes` | 协议 | `@schemes http https` |
| `@contact.name` | 联系人 | 可选 |

## 接口级注解（handler 上）

| 注解 | 说明 | 示例 |
|------|------|------|
| `@Summary` | 简短摘要 | `@Summary 创建用户` |
| `@Description` | 详细描述 | 多行可写多条 `@Description` |
| `@Tags` | 分组标签 | `@Tags users` |
| `@Accept` | 接受的请求类型 | `@Accept json` |
| `@Produce` | 响应的 Content-Type | `@Produce json` |
| `@Param` | 参数定义 | 见下 |
| `@Success` | 成功响应 | `@Success 200 {object} Type` |
| `@Failure` | 失败响应 | `@Failure 400 {object} object` |
| `@Router` | 路径与方法 | `@Router /api/v1/users [post]` |
| `@Security` | 安全定义 | `@Security ApiKeyAuth` |

## @Param 格式

```
@Param 名称 位置 类型 必填 "说明" 其他
```

- **位置**：`path` | `query` | `header` | `body` | `formData`
- **类型**：`int`、`string`、`bool`、或结构体（如 `handler.Req`）
- **必填**：`true` | `false`
- **其他**：如 `default(1)` 用于 query

多参数写多行 `@Param`。

## @Success / @Failure

- `@Success 状态码 返回类型 "说明"`
- 返回类型常用：`{object} Type`、`{array} Type`、`object`（泛型）
- 同一接口可写多条 `@Success`、`@Failure`（不同状态码）。

## 类型引用

- 同包：直接写结构体名，如 `UserResponse`。
- 其他包：包名.结构体名，如 `handler.CreateUserRequest`。
- 不关心具体结构时用 `object`。
