#!/bin/bash

# Docker 测试脚本

echo "🐳 开始 Docker 测试..."

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker 是否运行
if ! docker info &> /dev/null; then
    echo "❌ Docker 未运行，请启动 Docker"
    exit 1
fi

echo "✅ Docker 已安装并运行"

# 构建测试镜像
echo "🔨 构建测试镜像..."
if docker build -t dove:test .; then
    echo "✅ 镜像构建成功"
else
    echo "❌ 镜像构建失败"
    exit 1
fi

# 运行测试容器
echo "🚀 运行测试容器..."
if docker run -d --name dove-test -p 8080:8080 dove:test; then
    echo "✅ 容器启动成功"
    
    # 等待应用启动
    echo "⏳ 等待应用启动..."
    sleep 5
    
    # 测试健康检查
    if curl -f http://localhost:8080/health &> /dev/null; then
        echo "✅ 健康检查通过"
    else
        echo "⚠️  健康检查失败，但容器仍在运行"
    fi
    
    # 停止并清理测试容器
    echo "🧹 清理测试容器..."
    docker stop dove-test
    docker rm dove-test
    docker rmi dove:test
    
    echo "✅ Docker 测试完成"
else
    echo "❌ 容器启动失败"
    exit 1
fi
