#!/usr/bin/env bash
# 运行 treasure-doc 容器，bind mount 状态目录（data/backup/files）到宿主机，
# 保证重建/更新容器时 SQLite 数据与备份不丢失。前置：先 docker build -t treasure-doc .
set -euo pipefail

docker run -d \
  --name treasure-doc \
  --restart unless-stopped \
  -p 2026:2026 \
  -v "$(pwd)/data:/app/data" \
  -v "$(pwd)/backup:/app/backup" \
  -v "$(pwd)/files:/app/files" \
  treasure-doc:latest
