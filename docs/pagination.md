# 分页功能

## 概述

本项目实现了统一的分页功能，支持所有列表查询接口的分页查询。

## 功能特性

- **统一分页参数**: 使用 `page` 和 `page_size` 作为标准分页参数
- **自定义排序**: 支持 `sort_by` 和 `sort_order` 参数进行字段排序
- **模糊搜索**: 支持 `keyword` 和 `search_by` 参数进行字段搜索
- **参数验证**: 自动验证分页、排序和搜索参数的有效性
- **默认值处理**: 提供合理的默认值（page=1, page_size=10, sort_order=desc）
- **限制保护**: 限制每页最大记录数为100
- **完整响应**: 返回总记录数、总页数、是否有下一页等信息
- **安全排序**: 只允许对预定义字段进行排序，防止SQL注入
- **安全搜索**: 只允许对预定义字段进行搜索，防止SQL注入

## 分页参数

### 请求参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| page | int | 否 | 1 | 页码，从1开始 |
| page_size | int | 否 | 10 | 每页大小，最大100 |
| sort_by | string | 否 | - | 排序字段，详见支持的排序字段 |
| sort_order | string | 否 | desc | 排序方向：asc（升序）, desc（降序） |
| keyword | string | 否 | - | 搜索关键词 |
| search_by | string | 否 | - | 搜索字段，详见支持的搜索字段 |

### 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "data": [...],           // 数据列表
    "total": 100,            // 总记录数
    "page": 1,               // 当前页码
    "page_size": 10,         // 每页大小
    "total_pages": 10,       // 总页数
    "has_next": true,        // 是否有下一页
    "has_prev": false        // 是否有上一页
  }
}
```

## 使用示例

### 获取用户列表（分页、排序和搜索）

```bash
# 搜索包含"john"的用户
GET /api/users?keyword=john

# 在用户名中搜索"john"
GET /api/users?keyword=john&search_by=username

# 在邮箱中搜索"example"
GET /api/users?keyword=example&search_by=email

# 组合使用：搜索+排序+分页
GET /api/users?keyword=john&search_by=username&sort_by=created_at&sort_order=desc&page=1&page_size=5

# 不传参数，使用默认值（按创建时间降序）
GET /api/users
```

### 获取产品列表（分页、排序和搜索）

```bash
# 搜索包含"Product"的产品
GET /api/products?keyword=Product

# 在产品名称中搜索"Product"
GET /api/products?keyword=Product&search_by=name

# 在产品描述中搜索"product"
GET /api/products?keyword=product&search_by=description

# 组合使用：搜索+排序+分页
GET /api/products?keyword=Product&search_by=name&sort_by=price&sort_order=desc&page=1&page_size=5

# 不传参数，使用默认值（按创建时间降序）
GET /api/products
```

## 实现细节

### 分页工具包

位置: `pkg/pagination/pagination.go`

主要功能:
- `ParsePageRequest()`: 从 gin.Context 解析分页参数
- `NewPageResponse()`: 创建分页响应
- `GetOffset()`: 计算数据库查询偏移量
- `GetLimit()`: 获取查询限制数量

### 数据库查询

使用 GORM 的 `Offset()` 和 `Limit()` 方法实现分页:

```go
// 获取总记录数
db.Model(&model.User{}).Count(&total)

// 获取分页数据
db.Offset(offset).Limit(limit).Find(&users)
```

### 服务层集成

所有列表查询服务都支持分页:

```go
func (s *userService) GetAllWithPagination(ctx context.Context, page *pagination.PageRequest) (*pagination.PageResponse, error) {
    users, total, err := s.repo.GetAllWithPagination(ctx, page)
    if err != nil {
        return nil, err
    }
    
    // 处理敏感数据（如密码）
    for i := range users {
        users[i].Password = ""
    }
    
    return pagination.NewPageResponse(users, total, page.Page, page.PageSize), nil
}
```

## 支持的排序字段

### 用户管理

用户列表支持以下字段排序：
- `id` - 用户ID
- `username` - 用户名
- `email` - 邮箱
- `nickname` - 昵称
- `status` - 状态
- `created_at` - 创建时间
- `updated_at` - 更新时间

### 产品管理

产品列表支持以下字段排序：
- `id` - 产品ID
- `name` - 产品名称
- `price` - 价格
- `stock` - 库存
- `created_at` - 创建时间
- `updated_at` - 更新时间

## 支持的搜索字段

### 用户管理

用户列表支持以下字段搜索：
- `username` - 用户名（模糊匹配）
- `email` - 邮箱（模糊匹配）
- `nickname` - 昵称（模糊匹配）

**注意**: 如果不指定 `search_by` 参数，将在所有可搜索字段中进行搜索。

### 产品管理

产品列表支持以下字段搜索：
- `name` - 产品名称（模糊匹配）
- `description` - 产品描述（模糊匹配）

**注意**: 如果不指定 `search_by` 参数，将在所有可搜索字段中进行搜索。

## 支持的接口

### 用户管理

- `GET /api/users` - 获取用户列表（分页、排序和搜索）

### 产品管理

- `GET /api/products` - 获取产品列表（分页、排序和搜索）

## 注意事项

1. **性能考虑**: 大数据量时建议合理设置 `page_size`，避免单次查询过多数据
2. **缓存策略**: 对于频繁访问的列表，可以考虑缓存分页结果
3. **排序性能**: 对排序字段建立索引可以提高查询性能
4. **搜索性能**: 对搜索字段建立索引可以提高搜索性能
5. **安全排序**: 只允许对预定义字段进行排序，防止SQL注入攻击
6. **安全搜索**: 只允许对预定义字段进行搜索，防止SQL注入攻击

## 扩展功能

### 多字段排序支持

可以扩展排序功能，支持多字段排序:

```go
type PageRequest struct {
    Page      int      `json:"page" form:"page"`
    PageSize  int      `json:"page_size" form:"page_size"`
    SortBy    string   `json:"sort_by" form:"sort_by"`
    SortOrder string   `json:"sort_order" form:"sort_order"`
    SortFields []SortField `json:"sort_fields" form:"sort_fields"`
}

type SortField struct {
    Field string `json:"field"`
    Order string `json:"order"`
}
```

### 高级搜索支持

可以扩展搜索功能，支持多字段组合搜索:

```go
type PageRequest struct {
    Page      int           `json:"page" form:"page"`
    PageSize  int           `json:"page_size" form:"page_size"`
    Keyword   string        `json:"keyword" form:"keyword"`
    SearchBy  string        `json:"search_by" form:"search_by"`
    SearchFields []SearchField `json:"search_fields" form:"search_fields"`
}

type SearchField struct {
    Field string `json:"field"`
    Value string `json:"value"`
    Operator string `json:"operator"` // eq, ne, gt, lt, gte, lte, like
}
```

### 过滤支持

可以添加过滤条件:

```go
type PageRequest struct {
    Page     int    `json:"page" form:"page"`
    PageSize int    `json:"page_size" form:"page_size"`
    Status   *int   `json:"status" form:"status"`
    Category string `json:"category" form:"category"`
}
```
