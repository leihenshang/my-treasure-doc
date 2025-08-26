#!/bin/bash

# 用法: ./build.sh [dev]
# 默认构建 prod，传入 dev 构建开发环境

if [ "$1" = "dev" ]; then
    echo "构建开发环境镜像..."
    docker build -t treasure-doc-dev . --build-arg BINARY_NAME="treasure-doc-dev" --build-arg EXPOSE_PORT=2025
    docker save -o treasure-doc-dev.tar.gz treasure-doc-dev
    sudo chmod 777 treasure-doc-dev.tar.gz
    echo "开发环境镜像已保存: treasure-doc-dev.tar.gz"
else
    echo "构建生产环境镜像..."
    docker build -t treasure-doc .
    docker save -o treasure-doc.tar.gz treasure-doc
    sudo chmod 777 treasure-doc.tar.gz
    echo "生产环境镜像已保存: treasure-doc.tar.gz"
fi