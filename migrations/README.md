# 数据库迁移

本项目使用 `golang-migrate` 进行数据库迁移管理。

## 迁移文件命名规则

迁移文件必须遵循以下命名格式：

```
{version}_{name}.up.sql   - 升级迁移
{version}_{name}.down.sql - 回滚迁移
```

- `version`: 版本号（数字，递增）
- `name`: 迁移名称（描述性名称，使用下划线分隔）

示例：
- `000001_init_users_table.up.sql`
- `000001_init_users_table.down.sql`

## 使用方法

### 使用 Makefile（推荐）

```bash
# 执行所有待执行的迁移
make migrate-up

# 回滚一次迁移
make migrate-down

# 查看当前迁移版本
make migrate-version

# 创建新的迁移文件
make migrate-create NAME=add_user_table
```

### 直接使用 migrate 命令

```bash
# 执行所有待执行的迁移
go run cmd/migrate/main.go -command=up

# 回滚一次迁移
go run cmd/migrate/main.go -command=down

# 回滚到指定版本
go run cmd/migrate/main.go -command=down -version=1

# 强制设置版本（用于修复迁移状态）
go run cmd/migrate/main.go -command=force -version=1

# 查看当前版本
go run cmd/migrate/main.go -command=version

# 创建新的迁移文件
go run cmd/migrate/main.go -command=create -name=add_user_table
```

### 选项说明

- `-config`: 配置文件路径（默认: `configs/config.yaml`）
- `-path`: 迁移文件目录（默认: `migrations`）
- `-version`: 版本号（用于 down 和 force 命令）
- `-name`: 迁移文件名（用于 create 命令）

## 迁移文件编写规范

### Up Migration（升级迁移）

Up migration 应该包含：
- 创建表、索引、约束等
- 添加列、修改列等
- 插入初始数据

示例：

```sql
-- Up Migration: 创建用户表
CREATE TABLE IF NOT EXISTS `u_user` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `username` VARCHAR(255) NOT NULL,
    `phone` VARCHAR(20) NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### Down Migration（回滚迁移）

Down migration 应该包含：
- 删除表、索引、约束等
- 删除列等
- 回滚 Up migration 的所有操作

示例：

```sql
-- Down Migration: 删除用户表
DROP TABLE IF EXISTS `u_user`;
```

## 最佳实践

1. **总是编写 Down Migration**: 每个 Up migration 都应该有对应的 Down migration
2. **使用事务**: 确保迁移操作的原子性
3. **测试迁移**: 在开发环境测试 Up 和 Down migration
4. **版本控制**: 将迁移文件纳入版本控制
5. **命名规范**: 使用描述性的迁移文件名
6. **避免数据丢失**: Down migration 应该能够安全地回滚 Up migration 的操作

## 迁移历史

迁移工具会在数据库中创建 `schema_migrations` 表来记录迁移历史：

```sql
CREATE TABLE `schema_migrations` (
    `version` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
    `dirty` BOOLEAN NOT NULL
);
```

- `version`: 当前迁移版本号
- `dirty`: 标记迁移是否处于脏状态（迁移失败时）

## 故障处理

### 迁移失败

如果迁移失败，数据库可能处于"脏"状态。可以使用 `force` 命令修复：

```bash
# 查看当前版本
make migrate-version

# 强制设置版本（如果迁移失败，可能需要手动修复数据库后强制设置版本）
go run cmd/migrate/main.go -command=force -version=<version>
```

### 回滚迁移

```bash
# 回滚一次迁移
make migrate-down

# 回滚到指定版本
go run cmd/migrate/main.go -command=down -version=<version>
```

## 注意事项

1. **生产环境**: 在生产环境执行迁移前，务必备份数据库
2. **测试环境**: 先在测试环境验证迁移脚本
3. **团队协作**: 团队成员应该按顺序执行迁移
4. **版本冲突**: 避免创建相同版本的迁移文件
