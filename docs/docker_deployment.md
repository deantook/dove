# Docker 部署指南

## 概述

本文档介绍如何使用 Docker 部署 dove 应用，包括开发环境和生产环境的配置。

## 目录结构

```
dove/
├── Dockerfile                 # 生产环境 Dockerfile
├── Dockerfile.dev            # 开发环境 Dockerfile
├── docker-compose.yml        # 开发环境 Docker Compose
├── docker-compose.prod.yml   # 生产环境 Docker Compose
├── .dockerignore             # Docker 忽略文件
├── .air.toml                 # 开发环境热重载配置
├── nginx/
│   └── nginx.conf           # Nginx 反向代理配置
└── scripts/
    └── init.sql             # 数据库初始化脚本
```

## 快速开始

### 1. 开发环境

#### 使用 Docker Compose（推荐）

```bash
# 启动开发环境（包含热重载）
make docker-dev

# 或者直接使用 docker-compose
docker-compose --profile dev up -d
```

#### 使用 Docker 命令

```bash
# 构建开发环境镜像
make docker-build-dev

# 运行开发容器
docker run -d --name dove-dev -p 8080:8080 \
  -v $(pwd):/app \
  -v /app/vendor \
  dove:dev
```

### 2. 生产环境

#### 使用 Docker Compose

```bash
# 启动生产环境
make docker-prod

# 或者直接使用 docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

#### 使用 Docker 命令

```bash
# 构建生产镜像
make docker-build

# 运行生产容器
docker run -d --name dove-prod -p 8080:8080 dove:latest
```

## 环境配置

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `ENV` | `dev` | 环境类型（dev/production） |
| `GIN_MODE` | `debug` | Gin 模式（debug/release） |
| `MYSQL_ROOT_PASSWORD` | `root123456` | MySQL root 密码 |
| `MYSQL_DATABASE` | `dove` | MySQL 数据库名 |
| `MYSQL_USER` | `dove` | MySQL 用户名 |
| `MYSQL_PASSWORD` | `dove123456` | MySQL 用户密码 |
| `REDIS_PASSWORD` | `redis123456` | Redis 密码 |

### 端口配置

| 服务 | 端口 | 说明 |
|------|------|------|
| 应用服务 | 8080 | 主应用端口 |
| MySQL | 3306 | 数据库端口 |
| Redis | 6379 | 缓存端口 |
| Redis Commander | 8081 | Redis 管理工具 |
| Adminer | 8082 | 数据库管理工具 |
| Nginx | 80/443 | 反向代理端口 |

## 详细配置

### 1. 开发环境配置

开发环境使用 `Dockerfile.dev` 和 `docker-compose.yml`：

```yaml
# docker-compose.yml
services:
  app-dev:
    build:
      context: .
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    environment:
      - ENV=dev
      - GIN_MODE=debug
    volumes:
      - .:/app
      - /app/vendor
    depends_on:
      - mysql
      - redis
```

特点：
- 支持热重载（使用 Air）
- 代码变更自动重启
- 包含开发工具（Redis Commander、Adminer）

### 2. 生产环境配置

生产环境使用 `Dockerfile` 和 `docker-compose.prod.yml`：

```yaml
# docker-compose.prod.yml
services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - ENV=production
      - GIN_MODE=release
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

特点：
- 多阶段构建，镜像更小
- 非 root 用户运行
- 健康检查
- 资源限制

### 3. 数据库配置

#### MySQL 配置

```yaml
mysql:
  image: mysql:8.0
  environment:
    MYSQL_ROOT_PASSWORD: root123456
    MYSQL_DATABASE: dove
    MYSQL_USER: dove
    MYSQL_PASSWORD: dove123456
  volumes:
    - mysql_data:/var/lib/mysql
    - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
```

#### Redis 配置

```yaml
redis:
  image: redis:7-alpine
  command: redis-server --appendonly yes --requirepass redis123456
  volumes:
    - redis_data:/data
```

### 4. Nginx 反向代理

生产环境可以使用 Nginx 作为反向代理：

```bash
# 启动带 Nginx 的生产环境
docker-compose -f docker-compose.prod.yml --profile nginx up -d
```

Nginx 配置特点：
- SSL/TLS 支持
- Gzip 压缩
- 负载均衡
- 安全头设置
- 静态文件缓存

## 常用命令

### 开发环境

```bash
# 启动开发环境
make docker-dev

# 查看开发环境日志
make docker-logs-dev

# 进入开发容器
make docker-exec-dev

# 重新构建开发环境
make docker-rebuild-dev

# 停止开发环境
docker-compose down
```

### 生产环境

```bash
# 启动生产环境
make docker-prod

# 查看生产环境日志
make docker-logs

# 进入生产容器
make docker-exec

# 重新构建生产环境
make docker-rebuild

# 停止生产环境
docker-compose -f docker-compose.prod.yml down
```

### 数据库操作

```bash
# 初始化数据库
make docker-init-db

# 备份数据库
make docker-backup

# 恢复数据库
make docker-restore BACKUP_FILE=backup_20231201_120000.sql
```

### 其他操作

```bash
# 查看容器状态
make docker-status

# 清理 Docker 资源
make docker-clean

# 运行测试（在容器中）
make docker-test

# 生成 Swagger 文档（在容器中）
make docker-swagger

# 运行 Wire（在容器中）
make docker-wire
```

## 故障排除

### 常见问题

1. **端口冲突**
   ```bash
   # 检查端口占用
   netstat -tulpn | grep :8080
   
   # 修改端口映射
   docker-compose up -d -p 8081:8080
   ```

2. **数据库连接失败**
   ```bash
   # 检查数据库容器状态
   docker-compose ps mysql
   
   # 查看数据库日志
   docker-compose logs mysql
   
   # 重新初始化数据库
   make docker-init-db
   ```

3. **权限问题**
   ```bash
   # 修复文件权限
   sudo chown -R $USER:$USER .
   
   # 重新构建镜像
   make docker-rebuild
   ```

4. **内存不足**
   ```bash
   # 清理 Docker 资源
   make docker-clean
   
   # 增加 Docker 内存限制
   # 在 Docker Desktop 设置中调整内存限制
   ```

### 日志查看

```bash
# 查看应用日志
docker-compose logs -f app

# 查看数据库日志
docker-compose logs -f mysql

# 查看 Redis 日志
docker-compose logs -f redis

# 查看所有服务日志
docker-compose logs -f
```

### 性能优化

1. **镜像优化**
   - 使用多阶段构建
   - 减少镜像层数
   - 使用 .dockerignore 排除不必要文件

2. **容器优化**
   - 设置资源限制
   - 使用健康检查
   - 配置日志轮转

3. **网络优化**
   - 使用自定义网络
   - 配置 DNS
   - 优化端口映射

## 安全考虑

1. **镜像安全**
   - 使用官方基础镜像
   - 定期更新依赖
   - 扫描安全漏洞

2. **运行时安全**
   - 非 root 用户运行
   - 最小权限原则
   - 安全头设置

3. **网络安全**
   - 使用 HTTPS
   - 配置防火墙
   - 限制端口访问

## 监控和日志

### 监控

```bash
# 查看容器资源使用
docker stats

# 查看容器进程
docker top dove-app

# 查看容器信息
docker inspect dove-app
```

### 日志管理

```bash
# 配置日志驱动
docker run --log-driver=json-file --log-opt max-size=10m --log-opt max-file=3

# 查看日志
docker logs --tail 100 -f dove-app
```

## 扩展部署

### 1. 多实例部署

```yaml
# docker-compose.prod.yml
services:
  app:
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

### 2. 负载均衡

```nginx
# nginx/nginx.conf
upstream dove_backend {
    server app1:8080;
    server app2:8080;
    server app3:8080;
}
```

### 3. 高可用部署

- 使用 Docker Swarm 或 Kubernetes
- 配置健康检查和自动重启
- 实现服务发现和负载均衡

## 总结

Docker 部署提供了以下优势：

1. **环境一致性** - 开发、测试、生产环境一致
2. **快速部署** - 一键部署整个应用栈
3. **易于扩展** - 支持水平扩展和负载均衡
4. **资源隔离** - 容器间资源隔离
5. **版本管理** - 镜像版本管理
6. **回滚能力** - 快速回滚到之前的版本

通过合理配置 Docker 和 Docker Compose，可以大大简化应用的部署和运维工作。
