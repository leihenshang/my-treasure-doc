# 宝藏文档 (Treasure Doc) - API 后端服务

## 项目概述 Overview

宝藏文档的后端API服务，基于Gin框架构建的高性能文档管理系统。

Backend API service for Treasure Doc, a high-performance document management system built on Gin framework.

## 技术栈 Tech Stack

- **Web框架**: Gin Framework
- **数据库ORM**: GORM  
- **配置管理**: TOML配置文件
- **日志系统**: Zap日志处理
- **缓存支持**: Redis缓存处理

## 快速开始 Quick Start

### 环境要求 Requirements

- Go 1.22+
- MySQL 5.7+
- Redis (可选)

### 配置文件 Configuration

复制示例配置文件并修改相应参数：

```bash
cp config.example.toml config.toml
```

主要配置项：
- 应用端口：2021
- 数据库连接信息
- Redis配置（可选）

### 本地运行 Local Development

```bash
# 进入用户模块目录
cd module/user

# 安装依赖
go mod tidy

# 运行服务
go run main.go
```

## 编译部署 Build & Deployment

### 跨平台编译 Cross-platform Compilation

```bash
# Linux
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o treasure_user

# Windows  
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

# macOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build
```

### Docker部署 Docker Deployment

#### 构建镜像 Build Image

```bash
docker build -t treasure-doc .
```

#### 运行容器 Run Container

```bash
# 后台运行
docker run -d --name treasure-doc \
  --restart=always \
  -p 2021:2021 \
  -v /path/to/web:/app/web \
  -v /path/to/files:/app/files \
  -v /path/to/config.toml:/app/config.toml \
  treasure-doc

# 前台调试模式
docker run --rm --name treasure-doc -it \
  -p 2021:2021 \
  -v /path/to/web:/app/web \
  -v /path/to/files:/app/files \
  -v /path/to/config.toml:/app/config.toml \
  treasure-doc /bin/sh
```

## 开发工具 Development Tools

### 数据库模型生成 Generate Database Models

进入CLI目录生成Gin模型：

```bash
cd module/user/cli
go run . -gen
```

生成的模型将保存在 `data/model/` 目录下。

## 项目结构 Project Structure

```
module/user/
├── api/          # API接口层
├── config/       # 配置文件
├── data/         # 数据模型和传输对象
├── global/       # 全局变量和初始化
├── internal/     # 内部服务逻辑
├── router/       # 路由配置
├── utils/        # 工具函数
└── web/          # 前端静态文件
```

## 数据库维护 Database Maintenance

数据修复SQL：

```sql
-- 修复文档分组关联
UPDATE td_doc SET group_id = 'root' WHERE group_id = '' OR group_id = '0';

-- 修复文档组父子关系  
UPDATE td_doc_group SET p_id = 'root' WHERE p_id = '' OR p_id = '0';
```

## 许可证 License

MIT License
