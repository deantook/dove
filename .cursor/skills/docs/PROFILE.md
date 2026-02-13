# 用户资料系统技术设计文档

## 1. 系统概述

### 1.1 核心概念
用户资料系统是恋爱类交友APP的核心功能模块，采用**渐进式解锁机制**，通过聊天互动、时间积累、付费等方式逐步解锁对方的详细资料，增强用户粘性和互动深度。

### 1.2 设计原则
- **隐私保护**：用户可自主控制资料可见性
- **渐进式解锁**：通过互动逐步解锁，增强用户粘性
- **灵活配置**：支持系统预设和用户自定义资料字段
- **多维度解锁**：支持多种解锁方式和条件组合

---

## 2. 资料类型设计

### 2.1 资料分类

#### 2.1.1 按来源分类
- **默认资料（System Profile）**：系统预设的基础资料字段，所有用户必须填写
- **自定义资料（Custom Profile）**：用户自行创建和配置的个性化资料字段

#### 2.1.2 按数据类型分类

| 类型 | 标识 | 说明 | 示例 |
|------|------|------|------|
| 单行文本 | `TEXT_SINGLE` | 单行文本输入 | 昵称、职业 |
| 多行文本 | `TEXT_MULTI` | 多行文本输入 | 个人简介、兴趣爱好描述 |
| 单选 | `SELECT_SINGLE` | 单选下拉/单选按钮 | 性别、年龄范围、城市 |
| 多选 | `SELECT_MULTI` | 多选下拉/复选框 | 兴趣爱好、性格标签 |
| 标签 | `TAG` | 标签云，可多选 | 技能标签、性格标签 |
| 位置 | `LOCATION` | 地理位置信息 | 当前城市、家乡、常去地点 |
| 图片 | `IMAGE` | 单张或多张图片 | 照片、作品集 |
| 视频 | `VIDEO` | 视频文件 | 自我介绍视频、生活片段 |
| 媒体 | `MEDIA` | 混合媒体（图片+视频） | 动态相册 |
| 数字 | `NUMBER` | 数值类型 | 身高、体重、收入范围 |
| 日期 | `DATE` | 日期选择 | 生日、纪念日 |
| 时间 | `TIME` | 时间选择 | 作息时间、空闲时间 |
| 日期时间 | `DATETIME` | 日期时间选择 | 活动时间 |
| 布尔值 | `BOOLEAN` | 是/否选择 | 是否接受异地恋 |
| 范围 | `RANGE` | 数值范围选择 | 年龄范围、距离范围 |
| 颜色 | `COLOR` | 颜色选择器 | 喜欢的颜色 |
| 链接 | `LINK` | URL链接 | 社交媒体、个人网站 |
| 文件 | `FILE` | 文件上传 | 简历、证书 |

### 2.2 资料字段属性

每个资料字段包含以下属性：

```go
type ProfileField struct {
    ID              int       // 字段ID
    UserID          int       // 用户ID（自定义字段）
    FieldKey        string    // 字段唯一标识（如：nickname, age, hobbies）
    FieldName       string    // 字段显示名称
    FieldType       string    // 字段类型（见上表）
    IsSystem        bool      // 是否系统默认字段
    IsRequired      bool      // 是否必填
    IsSearchable    bool      // 是否可搜索
    IsPublic        bool      // 是否公开（无需解锁即可查看）
    DefaultValue    string    // 默认值（JSON格式）
    Options         string    // 选项配置（JSON格式，用于SELECT类型）
    Validation      string    // 验证规则（JSON格式）
    DisplayOrder    int       // 显示顺序
    Icon            string    // 图标URL
    Description     string    // 字段描述
    CreateTime      time.Time
    UpdateTime      time.Time
}
```

### 2.3 资料值存储

```go
type ProfileValue struct {
    ID          int       // 值ID
    UserID      int       // 用户ID
    FieldID     int       // 字段ID
    FieldKey    string    // 字段标识（冗余，便于查询）
    Value       string    // 值内容（JSON格式，支持复杂类型）
    ValueType   string    // 值类型（与FieldType对应）
    IsVerified  bool      // 是否已验证（如身份认证）
    CreateTime  time.Time
    UpdateTime  time.Time
}
```

**值存储格式说明**：
- 单行文本：`"value": "张三"`
- 多行文本：`"value": "热爱生活，喜欢旅行..."`
- 单选：`"value": "option_key"` 或 `"value": "选项值"`
- 多选：`"value": ["option1", "option2"]`
- 标签：`"value": ["tag1", "tag2", "tag3"]`
- 位置：`"value": {"lat": 39.9042, "lng": 116.4074, "address": "北京市朝阳区", "city": "北京"}`
- 图片：`"value": [{"url": "https://...", "thumbnail": "https://...", "order": 1}]`
- 视频：`"value": {"url": "https://...", "thumbnail": "https://...", "duration": 120}`
- 数字：`"value": 180`
- 日期：`"value": "1990-01-01"`
- 布尔值：`"value": true`
- 范围：`"value": {"min": 25, "max": 30}`

### 2.4 资料模板（ProfileTemplate）

资料模板用于存储系统预设的默认资料配置，可以包含多个字段及其解锁规则。模板可以被应用到新用户注册时，也可以作为用户创建自定义资料的参考。

```go
type ProfileTemplate struct {
    ID              int       // 模板ID
    TemplateKey     string    // 模板唯一标识（如：default, premium, basic）
    TemplateName    string    // 模板名称
    TemplateType    string    // 模板类型：DEFAULT（默认模板）、CUSTOM（自定义模板）
    Description     string    // 模板描述
    Fields          string    // 字段配置（JSON格式，包含字段定义）
    UnlockRules     string    // 解锁规则配置（JSON格式，包含各字段的解锁规则）
    IsActive        bool      // 是否启用
    IsDefault       bool      // 是否默认模板（新用户注册时自动应用）
    Version         int       // 模板版本号（用于模板升级）
    ApplyCount      int       // 应用次数统计
    CreateTime      time.Time
    UpdateTime      time.Time
}
```

**字段配置格式（Fields JSON）**：
```json
{
    "fields": [
        {
            "field_key": "nickname",
            "field_name": "昵称",
            "field_type": "TEXT_SINGLE",
            "is_required": true,
            "is_public": true,
            "display_order": 1,
            "validation": {
                "min_length": 2,
                "max_length": 20
            }
        },
        {
            "field_key": "age",
            "field_name": "年龄",
            "field_type": "NUMBER",
            "is_required": true,
            "is_public": false,
            "display_order": 2,
            "validation": {
                "min": 18,
                "max": 100
            }
        }
    ]
}
```

**解锁规则配置格式（UnlockRules JSON）**：
```json
{
    "rules": [
        {
            "field_key": "age",
            "unlock_type": "CHAT",
            "conditions": {
                "message_count": 30
            },
            "priority": 1
        },
        {
            "field_key": "bio",
            "unlock_type": "TIME",
            "conditions": {
                "friend_days": 7
            },
            "priority": 2
        }
    ]
}
```

---

## 3. 解锁机制设计

### 3.1 解锁方式（UnlockType）

| 解锁方式 | 标识 | 说明 |
|---------|------|------|
| 公开可见 | `PUBLIC` | 无需解锁，所有人可见 |
| 聊天解锁 | `CHAT` | 通过聊天互动解锁 |
| 时间解锁 | `TIME` | 成为好友后时间累积解锁 |
| 付费解锁 | `PAID` | 付费解锁 |
| 申请解锁 | `REQUEST` | 双方申请同意后解锁 |
| 组合解锁 | `COMBINED` | 多种方式组合（如：聊天+时间） |

### 3.2 解锁条件（UnlockCondition）

```go
type UnlockRule struct {
    ID              int       // 规则ID
    FieldID         int       // 字段ID（为空表示全局规则）
    UserID          int       // 用户ID（为空表示系统规则）
    UnlockType      string    // 解锁方式
    Conditions      string    // 解锁条件（JSON格式）
    Priority        int       // 优先级（数字越大优先级越高）
    IsActive        bool      // 是否启用
    CreateTime      time.Time
    UpdateTime      time.Time
}
```

#### 3.2.1 聊天解锁条件
```json
{
    "type": "CHAT",
    "conditions": {
        "message_count": 50,        // 聊天消息数量
        "chat_days": 7,             // 聊天天数
        "last_message_hours": 24    // 最后一条消息在24小时内
    }
}
```

#### 3.2.2 时间解锁条件
```json
{
    "type": "TIME",
    "conditions": {
        "friend_days": 30,          // 成为好友天数
        "continuous_days": 7        // 连续互动天数
    }
}
```

#### 3.2.3 付费解锁条件
```json
{
    "type": "PAID",
    "conditions": {
        "price": 9.9,               // 解锁价格（元）
        "currency": "CNY",          // 货币类型
        "discount": 0.8             // 折扣（可选）
    }
}
```

#### 3.2.4 申请解锁条件
```json
{
    "type": "REQUEST",
    "conditions": {
        "auto_approve": false,      // 是否自动同意
        "require_reason": true,     // 是否需要申请理由
        "max_requests_per_day": 3   // 每天最多申请次数
    }
}
```

#### 3.2.5 组合解锁条件
```json
{
    "type": "COMBINED",
    "conditions": {
        "logic": "AND",             // 逻辑关系：AND/OR
        "rules": [
            {
                "type": "CHAT",
                "message_count": 30
            },
            {
                "type": "TIME",
                "friend_days": 7
            }
        ]
    }
}
```

### 3.3 解锁记录（UnlockRecord）

```go
type UnlockRecord struct {
    ID              int       // 记录ID
    ViewerID        int       // 查看者用户ID
    OwnerID         int       // 资料所有者用户ID
    FieldID         int       // 字段ID（0表示解锁全部）
    UnlockType      string    // 解锁方式
    UnlockMethod    string    // 具体解锁方法（如：CHAT_MESSAGE, PAID, REQUEST_APPROVED）
    UnlockTime      time.Time // 解锁时间
    ExpireTime      time.Time // 过期时间（可选）
    Status          string    // 状态：ACTIVE, EXPIRED, REVOKED
    Metadata        string    // 元数据（JSON格式，记录解锁时的上下文）
    CreateTime      time.Time
    UpdateTime      time.Time
}
```

### 3.4 解锁申请（UnlockRequest）

```go
type UnlockRequest struct {
    ID              int       // 申请ID
    RequesterID     int       // 申请者用户ID
    OwnerID         int       // 资料所有者用户ID
    FieldID         int       // 字段ID（0表示申请解锁全部）
    Reason          string    // 申请理由
    Status          string    // 状态：PENDING, APPROVED, REJECTED, EXPIRED
    ResponseMessage string    // 回复消息
    RequestTime     time.Time // 申请时间
    ResponseTime     time.Time // 回复时间
    CreateTime      time.Time
    UpdateTime      time.Time
}
```

---

## 4. 数据库设计

### 4.1 表结构

#### 4.1.1 资料字段表（profile_fields）
```sql
CREATE TABLE `profile_fields` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT DEFAULT 0 COMMENT '用户ID，0表示系统字段',
    `field_key` VARCHAR(100) NOT NULL COMMENT '字段唯一标识',
    `field_name` VARCHAR(100) NOT NULL COMMENT '字段显示名称',
    `field_type` VARCHAR(50) NOT NULL COMMENT '字段类型',
    `is_system` TINYINT(1) DEFAULT 0 COMMENT '是否系统字段',
    `is_required` TINYINT(1) DEFAULT 0 COMMENT '是否必填',
    `is_searchable` TINYINT(1) DEFAULT 0 COMMENT '是否可搜索',
    `is_public` TINYINT(1) DEFAULT 0 COMMENT '是否公开',
    `default_value` TEXT COMMENT '默认值（JSON）',
    `options` TEXT COMMENT '选项配置（JSON）',
    `validation` TEXT COMMENT '验证规则（JSON）',
    `display_order` INT DEFAULT 0 COMMENT '显示顺序',
    `icon` VARCHAR(500) COMMENT '图标URL',
    `description` VARCHAR(500) COMMENT '字段描述',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_field_key` (`field_key`),
    INDEX `idx_is_system` (`is_system`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资料字段表';
```

#### 4.1.2 资料值表（profile_values）
```sql
CREATE TABLE `profile_values` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT NOT NULL COMMENT '用户ID',
    `field_id` INT NOT NULL COMMENT '字段ID',
    `field_key` VARCHAR(100) NOT NULL COMMENT '字段标识（冗余）',
    `value` TEXT NOT NULL COMMENT '值内容（JSON）',
    `value_type` VARCHAR(50) NOT NULL COMMENT '值类型',
    `is_verified` TINYINT(1) DEFAULT 0 COMMENT '是否已验证',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_user_field` (`user_id`, `field_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_field_id` (`field_id`),
    INDEX `idx_field_key` (`field_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资料值表';
```

#### 4.1.3 资料模板表（profile_templates）
```sql
CREATE TABLE `profile_templates` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `template_key` VARCHAR(100) NOT NULL COMMENT '模板唯一标识',
    `template_name` VARCHAR(100) NOT NULL COMMENT '模板名称',
    `template_type` VARCHAR(50) DEFAULT 'DEFAULT' COMMENT '模板类型：DEFAULT（默认模板）、CUSTOM（自定义模板）',
    `description` VARCHAR(500) COMMENT '模板描述',
    `fields` TEXT NOT NULL COMMENT '字段配置（JSON格式，包含字段定义）',
    `unlock_rules` TEXT COMMENT '解锁规则配置（JSON格式，包含各字段的解锁规则）',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `is_default` TINYINT(1) DEFAULT 0 COMMENT '是否默认模板（新用户注册时自动应用）',
    `version` INT DEFAULT 1 COMMENT '模板版本号（用于模板升级）',
    `apply_count` INT DEFAULT 0 COMMENT '应用次数统计',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_template_key` (`template_key`),
    INDEX `idx_template_type` (`template_type`),
    INDEX `idx_is_default` (`is_default`),
    INDEX `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资料模板表';
```

#### 4.1.4 解锁规则表（unlock_rules）
```sql
CREATE TABLE `unlock_rules` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `field_id` INT DEFAULT 0 COMMENT '字段ID，0表示全局规则',
    `user_id` INT DEFAULT 0 COMMENT '用户ID，0表示系统规则',
    `unlock_type` VARCHAR(50) NOT NULL COMMENT '解锁方式',
    `conditions` TEXT NOT NULL COMMENT '解锁条件（JSON）',
    `priority` INT DEFAULT 0 COMMENT '优先级',
    `is_active` TINYINT(1) DEFAULT 1 COMMENT '是否启用',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_field_id` (`field_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_unlock_type` (`unlock_type`),
    INDEX `idx_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='解锁规则表';
```

#### 4.1.5 解锁记录表（unlock_records）
```sql
CREATE TABLE `unlock_records` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `viewer_id` INT NOT NULL COMMENT '查看者用户ID',
    `owner_id` INT NOT NULL COMMENT '资料所有者用户ID',
    `field_id` INT DEFAULT 0 COMMENT '字段ID，0表示解锁全部',
    `unlock_type` VARCHAR(50) NOT NULL COMMENT '解锁方式',
    `unlock_method` VARCHAR(50) NOT NULL COMMENT '具体解锁方法',
    `unlock_time` DATETIME NOT NULL COMMENT '解锁时间',
    `expire_time` DATETIME COMMENT '过期时间',
    `status` VARCHAR(20) DEFAULT 'ACTIVE' COMMENT '状态',
    `metadata` TEXT COMMENT '元数据（JSON）',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_viewer_owner` (`viewer_id`, `owner_id`),
    INDEX `idx_field_id` (`field_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_expire_time` (`expire_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='解锁记录表';
```

#### 4.1.6 解锁申请表（unlock_requests）
```sql
CREATE TABLE `unlock_requests` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `requester_id` INT NOT NULL COMMENT '申请者用户ID',
    `owner_id` INT NOT NULL COMMENT '资料所有者用户ID',
    `field_id` INT DEFAULT 0 COMMENT '字段ID，0表示申请解锁全部',
    `reason` VARCHAR(500) COMMENT '申请理由',
    `status` VARCHAR(20) DEFAULT 'PENDING' COMMENT '状态',
    `response_message` VARCHAR(500) COMMENT '回复消息',
    `request_time` DATETIME NOT NULL COMMENT '申请时间',
    `response_time` DATETIME COMMENT '回复时间',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX `idx_requester_owner` (`requester_id`, `owner_id`),
    INDEX `idx_status` (`status`),
    INDEX `idx_field_id` (`field_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='解锁申请表';
```

#### 4.1.7 好友关系表（friendships）
```sql
CREATE TABLE `friendships` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT NOT NULL COMMENT '用户ID',
    `friend_id` INT NOT NULL COMMENT '好友ID',
    `status` VARCHAR(20) DEFAULT 'PENDING' COMMENT '状态：PENDING, ACCEPTED, BLOCKED',
    `friend_time` DATETIME COMMENT '成为好友时间',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_user_friend` (`user_id`, `friend_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_friend_id` (`friend_id`),
    INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友关系表';
```

#### 4.1.8 聊天统计表（chat_statistics）
```sql
CREATE TABLE `chat_statistics` (
    `id` INT PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT NOT NULL COMMENT '用户ID',
    `friend_id` INT NOT NULL COMMENT '好友ID',
    `message_count` INT DEFAULT 0 COMMENT '消息总数',
    `first_message_time` DATETIME COMMENT '首次消息时间',
    `last_message_time` DATETIME COMMENT '最后消息时间',
    `chat_days` INT DEFAULT 0 COMMENT '聊天天数',
    `continuous_days` INT DEFAULT 0 COMMENT '连续聊天天数',
    `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_user_friend` (`user_id`, `friend_id`),
    INDEX `idx_user_id` (`user_id`),
    INDEX `idx_friend_id` (`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天统计表';
```

---

## 5. API设计

### 5.1 资料模板管理

#### 5.1.1 获取所有模板列表
```
GET /api/v1/profile/templates
Query: ?type=DEFAULT&is_active=1&page=1&page_size=20
Response: {
    "code": 200,
    "data": {
        "list": [
            {
                "id": 1,
                "template_key": "default",
                "template_name": "默认模板",
                "template_type": "DEFAULT",
                "description": "系统默认资料模板",
                "is_active": true,
                "is_default": true,
                "version": 1,
                "apply_count": 1000,
                "fields_count": 10,
                "create_time": "2026-01-01 00:00:00"
            }
        ],
        "total": 5,
        "page": 1,
        "page_size": 20
    }
}
```

#### 5.1.2 获取模板详情
```
GET /api/v1/profile/templates/{template_id}
Response: {
    "code": 200,
    "data": {
        "id": 1,
        "template_key": "default",
        "template_name": "默认模板",
        "template_type": "DEFAULT",
        "description": "系统默认资料模板",
        "fields": {
            "fields": [...]
        },
        "unlock_rules": {
            "rules": [...]
        },
        "is_active": true,
        "is_default": true,
        "version": 1,
        "apply_count": 1000,
        "create_time": "2026-01-01 00:00:00",
        "update_time": "2026-01-01 00:00:00"
    }
}
```

#### 5.1.3 获取默认模板
```
GET /api/v1/profile/templates/default
Response: {
    "code": 200,
    "data": {
        // 同模板详情
    }
}
```

#### 5.1.4 创建模板（管理员）
```
POST /api/v1/profile/templates
Request: {
    "template_key": "premium",
    "template_name": "高级模板",
    "template_type": "DEFAULT",
    "description": "包含更多字段的高级资料模板",
    "fields": {
        "fields": [
            {
                "field_key": "nickname",
                "field_name": "昵称",
                "field_type": "TEXT_SINGLE",
                "is_required": true,
                "is_public": true,
                "display_order": 1
            },
            ...
        ]
    },
    "unlock_rules": {
        "rules": [
            {
                "field_key": "age",
                "unlock_type": "CHAT",
                "conditions": {
                    "message_count": 30
                },
                "priority": 1
            },
            ...
        ]
    },
    "is_default": false
}
```

#### 5.1.5 更新模板（管理员）
```
PUT /api/v1/profile/templates/{template_id}
Request: {
    "template_name": "更新后的模板名称",
    "description": "更新后的描述",
    "fields": {...},
    "unlock_rules": {...},
    "version": 2  // 版本号递增
}
```

#### 5.1.6 删除模板（管理员）
```
DELETE /api/v1/profile/templates/{template_id}
Response: {
    "code": 200,
    "message": "删除成功"
}
```

#### 5.1.7 应用模板到用户
```
POST /api/v1/profile/templates/{template_id}/apply
Request: {
    "user_id": 123  // 可选，不传则应用到当前用户
}
Response: {
    "code": 200,
    "data": {
        "applied_fields": 10,
        "applied_rules": 5,
        "message": "模板应用成功"
    }
}
```

#### 5.1.8 设置默认模板（管理员）
```
PUT /api/v1/profile/templates/{template_id}/set-default
Request: {
    "is_default": true
}
Response: {
    "code": 200,
    "message": "设置成功"
}
```

### 5.2 资料字段管理

#### 5.1.1 获取系统默认字段列表
```
GET /api/v1/profile/fields/system
Response: {
    "code": 200,
    "data": [
        {
            "id": 1,
            "field_key": "nickname",
            "field_name": "昵称",
            "field_type": "TEXT_SINGLE",
            "is_required": true,
            "is_public": true,
            ...
        }
    ]
}
```

#### 5.1.2 获取用户自定义字段列表
```
GET /api/v1/profile/fields/custom
Response: {
    "code": 200,
    "data": [...]
}
```

#### 5.1.3 创建自定义字段
```
POST /api/v1/profile/fields
Request: {
    "field_key": "my_hobby",
    "field_name": "我的爱好",
    "field_type": "TAG",
    "is_required": false,
    "options": {
        "tags": ["阅读", "旅行", "摄影", "音乐"]
    },
    "unlock_rules": [
        {
            "unlock_type": "CHAT",
            "conditions": {
                "message_count": 50
            }
        }
    ]
}
```

#### 5.1.4 更新字段配置
```
PUT /api/v1/profile/fields/{field_id}
```

#### 5.1.5 删除自定义字段
```
DELETE /api/v1/profile/fields/{field_id}
```

### 5.3 资料值管理

#### 5.2.1 获取自己的完整资料
```
GET /api/v1/profile/me
Response: {
    "code": 200,
    "data": {
        "fields": [
            {
                "field_id": 1,
                "field_key": "nickname",
                "field_name": "昵称",
                "field_type": "TEXT_SINGLE",
                "value": "张三",
                "is_public": true
            },
            ...
        ]
    }
}
```

#### 5.2.2 获取他人资料（根据解锁状态）
```
GET /api/v1/profile/{user_id}
Response: {
    "code": 200,
    "data": {
        "user_id": 123,
        "fields": [
            {
                "field_id": 1,
                "field_key": "nickname",
                "field_name": "昵称",
                "field_type": "TEXT_SINGLE",
                "value": "张三",
                "is_unlocked": true,
                "unlock_time": "2026-02-10 10:00:00"
            },
            {
                "field_id": 2,
                "field_key": "age",
                "field_name": "年龄",
                "field_type": "NUMBER",
                "value": null,
                "is_unlocked": false,
                "unlock_rules": [
                    {
                        "unlock_type": "CHAT",
                        "conditions": {
                            "message_count": 30
                        },
                        "progress": {
                            "current": 15,
                            "target": 30,
                            "percent": 50
                        }
                    }
                ]
            }
        ]
    }
}
```

#### 5.2.3 更新资料值
```
PUT /api/v1/profile/values
Request: {
    "field_id": 1,
    "value": "新昵称"
}
```

#### 5.2.4 批量更新资料值
```
PUT /api/v1/profile/values/batch
Request: {
    "values": [
        {"field_id": 1, "value": "新昵称"},
        {"field_id": 2, "value": 25}
    ]
}
```

### 5.4 解锁管理

#### 5.3.1 检查解锁状态
```
GET /api/v1/profile/{user_id}/unlock-status
Response: {
    "code": 200,
    "data": {
        "unlocked_fields": [1, 2, 3],
        "locked_fields": [
            {
                "field_id": 4,
                "field_name": "联系方式",
                "unlock_rules": [...],
                "can_unlock_now": false,
                "unlock_progress": {...}
            }
        ]
    }
}
```

#### 5.3.2 付费解锁
```
POST /api/v1/profile/{user_id}/unlock/pay
Request: {
    "field_id": 4,
    "payment_method": "wechat"
}
Response: {
    "code": 200,
    "data": {
        "order_id": "xxx",
        "payment_url": "https://...",
        "unlock_record": {...}
    }
}
```

#### 5.3.3 申请解锁
```
POST /api/v1/profile/{user_id}/unlock/request
Request: {
    "field_id": 4,
    "reason": "想了解更多关于你的信息"
}
Response: {
    "code": 200,
    "data": {
        "request_id": 123,
        "status": "PENDING"
    }
}
```

#### 5.3.4 处理解锁申请
```
PUT /api/v1/profile/unlock-requests/{request_id}
Request: {
    "action": "APPROVE", // APPROVE or REJECT
    "response_message": "好的，同意解锁"
}
```

#### 5.3.5 获取解锁记录
```
GET /api/v1/profile/unlock-records
Query: ?type=VIEWED&type=UNLOCKED&page=1&page_size=20
Response: {
    "code": 200,
    "data": {
        "list": [...],
        "total": 100,
        "page": 1,
        "page_size": 20
    }
}
```

### 5.5 解锁规则管理

#### 5.4.1 设置字段解锁规则
```
POST /api/v1/profile/fields/{field_id}/unlock-rules
Request: {
    "rules": [
        {
            "unlock_type": "CHAT",
            "conditions": {
                "message_count": 50
            },
            "priority": 1
        },
        {
            "unlock_type": "TIME",
            "conditions": {
                "friend_days": 7
            },
            "priority": 2
        }
    ]
}
```

#### 5.4.2 更新解锁规则
```
PUT /api/v1/profile/unlock-rules/{rule_id}
```

#### 5.4.3 删除解锁规则
```
DELETE /api/v1/profile/unlock-rules/{rule_id}
```

---

## 6. 业务逻辑设计

### 6.1 解锁判断流程

```
1. 用户请求查看他人资料
2. 获取资料所有者的解锁规则（优先级排序）
3. 遍历规则，检查是否满足解锁条件：
   a. PUBLIC：直接解锁
   b. CHAT：检查聊天统计
   c. TIME：检查好友关系时间
   d. PAID：检查付费记录
   e. REQUEST：检查申请状态
   f. COMBINED：按逻辑关系组合判断
4. 返回解锁状态和进度信息
```

### 6.2 自动解锁触发

#### 6.2.1 聊天消息触发
```
当用户发送消息时：
1. 更新聊天统计表
2. 检查是否有新的字段满足解锁条件
3. 自动创建解锁记录
4. 推送解锁通知
```

#### 6.2.2 时间触发
```
定时任务（每小时执行）：
1. 扫描好友关系表
2. 计算成为好友天数
3. 检查时间解锁条件
4. 自动创建解锁记录
```

### 6.3 解锁进度计算

```go
type UnlockProgress struct {
    Current    int     // 当前值
    Target     int     // 目标值
    Percent    float64 // 完成百分比
    Remaining  int     // 剩余值
    EstimatedDays int  // 预计完成天数（基于历史数据）
}
```

**示例**：
- 聊天解锁：当前15条消息，目标30条，进度50%，剩余15条
- 时间解锁：成为好友3天，目标7天，进度42.9%，剩余4天

### 6.4 资料模板应用流程

#### 6.4.1 新用户注册时应用默认模板
```
1. 用户注册成功后
2. 获取系统默认模板（is_default = true）
3. 解析模板中的字段配置
4. 为每个字段创建 profile_fields 记录（user_id = 0, is_system = true）
5. 解析模板中的解锁规则配置
6. 为每个字段创建 unlock_rules 记录
7. 更新模板的 apply_count 计数
8. 返回应用结果
```

#### 6.4.2 手动应用模板
```
1. 用户选择要应用的模板
2. 检查模板是否启用（is_active = true）
3. 解析模板配置
4. 检查字段是否已存在：
   a. 如果字段已存在，跳过或更新（根据策略）
   b. 如果字段不存在，创建新字段
5. 检查解锁规则是否已存在：
   a. 如果规则已存在，跳过或更新（根据策略）
   b. 如果规则不存在，创建新规则
6. 更新模板的 apply_count 计数
7. 返回应用结果（应用了多少字段和规则）
```

#### 6.4.3 模板版本管理
```
1. 模板更新时，version 字段递增
2. 记录模板变更历史（可选，通过版本表）
3. 已应用模板的用户可以选择升级到新版本
4. 升级时：
   a. 对比新旧版本差异
   b. 新增字段自动添加
   c. 删除字段标记为废弃（不直接删除）
   d. 修改字段更新配置
   e. 解锁规则同步更新
```

### 6.5 资料可见性控制

```go
func GetProfileVisibility(userID, viewerID int) map[int]bool {
    // 1. 如果是自己查看，返回全部可见
    if userID == viewerID {
        return allVisible()
    }
    
    // 2. 获取所有字段
    fields := getProfileFields(userID)
    
    // 3. 获取解锁记录
    unlocks := getUnlockRecords(viewerID, userID)
    
    // 4. 检查每个字段的可见性
    visibility := make(map[int]bool)
    for _, field := range fields {
        if field.IsPublic {
            visibility[field.ID] = true
            continue
        }
        
        // 检查是否已解锁
        if isUnlocked(field.ID, unlocks) {
            visibility[field.ID] = true
        } else {
            visibility[field.ID] = false
        }
    }
    
    return visibility
}
```

---

## 7. 扩展场景设计

### 7.1 资料验证机制

- **身份认证**：身份证、学历证书等验证
- **照片认证**：真人照片认证
- **视频认证**：真人视频认证
- **社交认证**：绑定社交媒体账号

验证后的资料值标记 `is_verified = true`，在展示时显示认证标识。

### 7.2 资料隐私等级

- **完全公开**：所有人可见
- **好友可见**：仅好友可见
- **解锁可见**：满足解锁条件后可见
- **仅自己可见**：完全私密

### 7.3 资料过期机制

某些解锁记录可以设置过期时间：
- **临时解锁**：付费解锁可设置7天有效期
- **活跃解锁**：聊天解锁需要保持活跃度，30天无互动自动失效
- **永久解锁**：满足特定条件后永久解锁（如成为好友90天）

### 7.4 资料分享机制

用户可以主动分享资料给他人：
- **分享链接**：生成临时分享链接，24小时有效
- **分享码**：生成分享码，对方输入后解锁
- **互相关注**：双方互相关注后自动解锁部分资料

### 7.5 资料推荐系统

基于解锁行为推荐相似用户：
- 分析用户解锁的字段类型偏好
- 推荐具有相似资料结构的用户
- 推荐解锁进度相近的用户

### 7.6 资料统计与分析

- **解锁统计**：各字段的解锁率、解锁方式分布
- **互动分析**：解锁后的互动频率变化
- **用户画像**：基于资料数据构建用户画像

---

## 8. 性能优化

### 8.1 缓存策略

- **字段配置缓存**：系统字段配置缓存到Redis，TTL 1小时
- **解锁状态缓存**：用户间的解锁状态缓存，TTL 5分钟
- **资料值缓存**：热点用户的资料值缓存，TTL 10分钟

### 8.2 数据库优化

- **索引优化**：在 `user_id`, `field_id`, `viewer_id`, `owner_id` 上建立联合索引
- **分表策略**：解锁记录表按时间分表（按月）
- **读写分离**：查询走从库，写入走主库

### 8.3 异步处理

- **解锁检查**：异步检查解锁条件，避免阻塞主流程
- **通知推送**：解锁通知异步推送
- **统计更新**：聊天统计异步更新

---

## 9. 安全设计

### 9.1 权限控制

- **字段权限**：用户只能修改自己的资料值
- **规则权限**：用户只能设置自己自定义字段的解锁规则
- **申请权限**：防止恶意申请，限制申请频率

### 9.2 数据校验

- **字段类型校验**：严格校验值类型与字段类型匹配
- **选项校验**：单选/多选值必须在选项范围内
- **格式校验**：URL、邮箱、手机号等格式校验

### 9.3 防刷机制

- **解锁频率限制**：同一用户对同一字段的解锁申请频率限制
- **付费风控**：付费解锁需要风控审核
- **异常检测**：检测异常解锁行为并告警

---

## 10. 实施计划

### 10.1 第一阶段：基础功能
1. 实现系统默认字段
2. 实现资料值的CRUD
3. 实现基础的公开/私密控制

### 10.2 第二阶段：解锁机制
1. 实现聊天解锁
2. 实现时间解锁
3. 实现解锁记录和状态查询

### 10.3 第三阶段：高级功能
1. 实现付费解锁
2. 实现申请解锁
3. 实现组合解锁

### 10.4 第四阶段：优化与扩展
1. 性能优化
2. 缓存机制
3. 统计分析功能

---

## 11. 注意事项

1. **数据一致性**：解锁记录与聊天统计、好友关系需要保持一致性
2. **用户体验**：解锁进度要清晰展示，给用户明确的目标感
3. **隐私保护**：严格遵守隐私保护规范，用户可随时撤回解锁授权
4. **扩展性**：字段类型和解锁方式要易于扩展
5. **性能**：大量用户场景下要考虑查询性能，合理使用缓存和索引

---

## 12. 附录

### 12.1 系统默认字段建议

| 字段Key | 字段名称 | 类型 | 是否必填 | 是否公开 |
|---------|---------|------|---------|---------|
| nickname | 昵称 | TEXT_SINGLE | 是 | 是 |
| avatar | 头像 | IMAGE | 是 | 是 |
| gender | 性别 | SELECT_SINGLE | 是 | 是 |
| age | 年龄 | NUMBER | 是 | 否 |
| city | 城市 | LOCATION | 是 | 否 |
| bio | 个人简介 | TEXT_MULTI | 否 | 否 |
| hobbies | 兴趣爱好 | TAG | 否 | 否 |
| photos | 照片 | IMAGE | 否 | 否 |
| video | 自我介绍视频 | VIDEO | 否 | 否 |

### 12.2 解锁规则示例

**示例1：基础资料聊天解锁**
- 昵称、头像：公开
- 年龄：聊天30条消息解锁
- 城市：聊天50条消息解锁
- 个人简介：聊天100条消息解锁

**示例2：深度资料时间解锁**
- 联系方式：成为好友30天解锁
- 详细地址：成为好友60天解锁

**示例3：付费解锁**
- 联系方式：付费9.9元解锁
- 详细资料包：付费19.9元解锁全部

**示例4：组合解锁**
- 联系方式：聊天50条消息 + 成为好友7天
- 详细地址：聊天100条消息 + 成为好友30天 + 付费9.9元（三选一）

### 12.3 资料模板示例

#### 示例1：默认模板（default）
```json
{
    "template_key": "default",
    "template_name": "默认模板",
    "template_type": "DEFAULT",
    "description": "系统默认资料模板，包含基础字段",
    "fields": {
        "fields": [
            {
                "field_key": "nickname",
                "field_name": "昵称",
                "field_type": "TEXT_SINGLE",
                "is_required": true,
                "is_public": true,
                "display_order": 1,
                "validation": {
                    "min_length": 2,
                    "max_length": 20
                }
            },
            {
                "field_key": "avatar",
                "field_name": "头像",
                "field_type": "IMAGE",
                "is_required": true,
                "is_public": true,
                "display_order": 2
            },
            {
                "field_key": "gender",
                "field_name": "性别",
                "field_type": "SELECT_SINGLE",
                "is_required": true,
                "is_public": true,
                "display_order": 3,
                "options": {
                    "options": [
                        {"key": "male", "label": "男"},
                        {"key": "female", "label": "女"},
                        {"key": "other", "label": "其他"}
                    ]
                }
            },
            {
                "field_key": "age",
                "field_name": "年龄",
                "field_type": "NUMBER",
                "is_required": true,
                "is_public": false,
                "display_order": 4,
                "validation": {
                    "min": 18,
                    "max": 100
                }
            },
            {
                "field_key": "city",
                "field_name": "城市",
                "field_type": "LOCATION",
                "is_required": true,
                "is_public": false,
                "display_order": 5
            },
            {
                "field_key": "bio",
                "field_name": "个人简介",
                "field_type": "TEXT_MULTI",
                "is_required": false,
                "is_public": false,
                "display_order": 6,
                "validation": {
                    "max_length": 500
                }
            }
        ]
    },
    "unlock_rules": {
        "rules": [
            {
                "field_key": "age",
                "unlock_type": "CHAT",
                "conditions": {
                    "message_count": 30
                },
                "priority": 1
            },
            {
                "field_key": "city",
                "unlock_type": "CHAT",
                "conditions": {
                    "message_count": 50
                },
                "priority": 1
            },
            {
                "field_key": "bio",
                "unlock_type": "CHAT",
                "conditions": {
                    "message_count": 100
                },
                "priority": 1
            }
        ]
    },
    "is_default": true,
    "is_active": true
}
```

#### 示例2：高级模板（premium）
```json
{
    "template_key": "premium",
    "template_name": "高级模板",
    "template_type": "DEFAULT",
    "description": "包含更多字段的高级资料模板",
    "fields": {
        "fields": [
            // 包含默认模板的所有字段
            // 额外添加：
            {
                "field_key": "hobbies",
                "field_name": "兴趣爱好",
                "field_type": "TAG",
                "is_required": false,
                "is_public": false,
                "display_order": 7,
                "options": {
                    "tags": ["阅读", "旅行", "摄影", "音乐", "运动", "美食", "电影", "游戏"]
                }
            },
            {
                "field_key": "photos",
                "field_name": "照片",
                "field_type": "IMAGE",
                "is_required": false,
                "is_public": false,
                "display_order": 8,
                "validation": {
                    "max_count": 9
                }
            },
            {
                "field_key": "video",
                "field_name": "自我介绍视频",
                "field_type": "VIDEO",
                "is_required": false,
                "is_public": false,
                "display_order": 9,
                "validation": {
                    "max_duration": 60
                }
            }
        ]
    },
    "unlock_rules": {
        "rules": [
            // 包含默认模板的所有规则
            // 额外添加：
            {
                "field_key": "hobbies",
                "unlock_type": "TIME",
                "conditions": {
                    "friend_days": 7
                },
                "priority": 2
            },
            {
                "field_key": "photos",
                "unlock_type": "COMBINED",
                "conditions": {
                    "logic": "OR",
                    "rules": [
                        {
                            "type": "CHAT",
                            "message_count": 200
                        },
                        {
                            "type": "TIME",
                            "friend_days": 30
                        }
                    ]
                },
                "priority": 1
            },
            {
                "field_key": "video",
                "unlock_type": "PAID",
                "conditions": {
                    "price": 9.9,
                    "currency": "CNY"
                },
                "priority": 1
            }
        ]
    },
    "is_default": false,
    "is_active": true
}
```

---

**文档版本**：v1.1  
**最后更新**：2026-02-13  
**维护者**：开发团队
