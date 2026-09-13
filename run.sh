#!/usr/bin/env bash
# 运行 treasure-doc 容器，bind mount 配置与状态目录（data/backup/files）到宿主机，
# 保证重建/更新容器时配置、SQLite 数据与备份不丢失。前置：先 docker build -t treasure-doc .
set -euo pipefail

# bind mount 单文件要求宿主机上文件已存在，缺失时先用示例配置初始化
if [ ! -f "$(pwd)/config.toml" ]; then
  cp module/user/config.example.toml config.toml
fi

# 首次把上传目录切到宿主机前，先把容器里已有的文件复制出来，否则会被空目录遮住：
#   mkdir -p files && docker cp treasure-doc:/app/files/. ./files/
# 上传目录的路径是代码里的常量（module/user/config/app.go: FilesPath = "files"），
# 相对容器工作目录 /app，因此只能挂到 /app/files；实际写入的是 /app/files/blog。
docker run -d \
  --name treasure-doc \
  --restart unless-stopped \
  -p 2026:2026 \
  -v "$(pwd)/config.toml:/app/config.toml" \
  -v "$(pwd)/data:/app/data" \
  -v "$(pwd)/backup:/app/backup" \
  -v "$(pwd)/web:/app/web" \
  -v "$(pwd)/files:/app/files" \
  treasure-doc:v0.0.1

  # docker exec -it treasure-doc sh