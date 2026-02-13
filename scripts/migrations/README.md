# 数据库迁移脚本

## 使用说明

### 1. 执行 SQL 脚本

执行 `001_create_profile_field_templates.sql` 创建系统资料字段模板表：

```bash
mysql -u username -p database_name < 001_create_profile_field_templates.sql
```

或者使用 MySQL 客户端：

```sql
source /path/to/001_create_profile_field_templates.sql;
```

### 2. 验证表结构

```sql
DESCRIBE profile_field_templates;
SHOW INDEX FROM profile_field_templates;
```

### 3. 查看示例数据

```sql
-- 查看所有字段模板
SELECT * FROM profile_field_templates;

-- 按分类查看
SELECT * FROM profile_field_templates WHERE category = '基本信息';

-- 查看启用的字段模板
SELECT * FROM profile_field_templates WHERE is_active = 1 ORDER BY category, display_order;
```

## 表结构说明

- `profile_field_templates`: 系统资料字段模板表
  - 存储系统预设的**单个字段类型定义**（如：姓名、学历、毕业学校等）
  - 每个字段类型是一条记录
  - 用户引用后会在 `profile_fields` 表中复制一条记录，`user_id` 设置为用户ID
  - 包含字段配置和默认解锁规则

## 字段说明

- `field_key`: 字段唯一标识（如：name, education, school）
- `field_name`: 字段显示名称（如：姓名、学历、毕业学校）
- `field_type`: 字段类型（TEXT_SINGLE, SELECT_SINGLE, NUMBER等）
- `category`: 字段分类（基本信息、教育背景、工作信息、联系方式等）
- `default_unlock_rules`: 默认解锁规则（JSON格式）

## 使用流程

1. **系统预设字段模板**：管理员在 `profile_field_templates` 表中创建字段模板
2. **用户引用模板**：用户通过 API 引用字段模板
3. **创建用户字段**：系统在 `profile_fields` 表中创建一条记录，`user_id` 设置为用户ID
4. **用户填写值**：用户在 `profile_values` 表中填写字段值

## 注意事项

1. 执行脚本前请备份数据库
2. 脚本中包含示例数据插入，如果已存在会更新 `update_time`
3. 生产环境执行前请先测试
4. 字段模板的 `field_key` 必须唯一
