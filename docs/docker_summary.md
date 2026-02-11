# Docker 功能总结

## 概述

为 dove 项目添加了完整的 Docker 支持，包括开发环境和生产环境的配置，以及相关的工具和文档。

## 新增文件

### 1. Docker 配置文件

- `Dockerfile` - 生产环境多阶段构建
- `Dockerfile.dev` - 开发环境配置
- `docker-compose.yml` - 开发环境服务编排
- `docker-compose.prod.yml` - 生产环境服务编排
- `.dockerignore` - Docker 构建忽略文件

### 2. 配置文件

- `nginx/nginx.conf` - Nginx 反向代理配置
- `.air.toml` - 开发环境热重载配置
- `scripts/init.sql` - 数据库初始化脚本

### 3. 文档

- `docs/docker_deployment.md` - Docker 部署指南
- `docs/docker_summary.md` - Docker 功能总结

### 4. 脚本

- `scripts/test-docker.sh` - Docker 测试脚本

## 功能特性

### 1. 多阶段构建

**生产环境 Dockerfile**：
- 使用 Go 1.21 Alpine 作为构建环境
- 多阶段构建，最终镜像基于 Alpine
- 非 root 用户运行，提高安全性
- 包含健康检查
- 静态链接二进制文件

**开发环境 Dockerfile**：
- 支持热重载（使用 Air）
- 包含开发工具
- 代码变更自动重启

### 2. 服务编排

**开发环境**：
- 应用服务（支持热重载）
- MySQL 数据库
- Redis 缓存
- Redis Commander（管理工具）
- Adminer（数据库管理工具）

**生产环境**：
- 应用服务（多实例）
- MySQL 数据库
- Redis 缓存
- Nginx 反向代理（可选）

### 3. 环境配置

**环境变量**：
- `ENV` - 环境类型（dev/production）
- `GIN_MODE` - Gin 模式（debug/release）
- 数据库和 Redis 配置

**端口配置**：
- 应用服务：8080
- MySQL：3306
- Redis：6379
- Redis Commander：8081
- Adminer：8082
- Nginx：80/443

### 4. 网络和安全

**网络配置**：
- 自定义网络 `dove-network`
- 服务间通信隔离
- 端口映射配置

**安全特性**：
- 非 root 用户运行
- 最小权限原则
- 安全头设置
- SSL/TLS 支持

### 5. 数据持久化

**数据卷**：
- `mysql_data` - MySQL 数据持久化
- `redis_data` - Redis 数据持久化
- 配置文件挂载

## 使用方法

### 1. 开发环境

```bash
# 启动开发环境
make docker-dev

# 查看日志
make docker-logs-dev

# 进入容器
make docker-exec-dev
```

### 2. 生产环境

```bash
# 启动生产环境
make docker-prod

# 查看日志
make docker-logs

# 进入容器
make docker-exec
```

### 3. 数据库操作

```bash
# 初始化数据库
make docker-init-db

# 备份数据库
make docker-backup

# 恢复数据库
make docker-restore BACKUP_FILE=backup.sql
```

## 优势

### 1. 环境一致性

- 开发、测试、生产环境完全一致
- 避免"在我机器上能运行"的问题
- 快速环境搭建

### 2. 快速部署

- 一键部署整个应用栈
- 支持多实例部署
- 自动服务发现

### 3. 易于扩展

- 水平扩展支持
- 负载均衡配置
- 高可用部署

### 4. 资源隔离

- 容器间资源隔离
- 独立网络配置
- 安全隔离

### 5. 版本管理

- 镜像版本管理
- 快速回滚能力
- 部署历史记录

## 最佳实践

### 1. 镜像优化

- 使用多阶段构建
- 减少镜像层数
- 使用 .dockerignore 排除不必要文件
- 使用 Alpine 基础镜像

### 2. 容器优化

- 设置资源限制
- 使用健康检查
- 配置日志轮转
- 非 root 用户运行

### 3. 网络优化

- 使用自定义网络
- 配置 DNS
- 优化端口映射
- 安全组配置

### 4. 数据管理

- 数据卷持久化
- 定期备份
- 数据迁移策略
- 监控和告警

## 监控和日志

### 1. 监控

- 容器资源使用监控
- 应用性能监控
- 健康检查监控
- 告警配置

### 2. 日志管理

- 结构化日志
- 日志轮转
- 日志聚合
- 日志分析

## 故障排除

### 1. 常见问题

- 端口冲突
- 数据库连接失败
- 权限问题
- 内存不足

### 2. 调试技巧

- 查看容器日志
- 进入容器调试
- 检查网络连接
- 验证配置

## 扩展功能

### 1. 高可用部署

- Docker Swarm 集群
- Kubernetes 部署
- 服务发现
- 负载均衡

### 2. 监控和告警

- Prometheus 监控
- Grafana 仪表板
- 告警通知
- 性能分析

### 3. CI/CD 集成

- 自动化构建
- 自动化测试
- 自动化部署
- 回滚机制

## 总结

Docker 支持为 dove 项目提供了：

1. **完整的容器化解决方案** - 从开发到生产的全流程支持
2. **简化的部署流程** - 一键部署整个应用栈
3. **环境一致性** - 消除环境差异问题
4. **易于扩展** - 支持水平扩展和高可用部署
5. **安全可靠** - 遵循容器安全最佳实践

通过这些功能，开发者可以更快速、更安全、更可靠地部署和管理 dove 应用。
