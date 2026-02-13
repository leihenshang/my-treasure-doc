# 宝藏文档(treasure-doc)-API

## 概述

宝藏文档的后端api

## 使用到的库

- gin framework
- gorm
- 配置文件
- 日志处理
- 缓存处理

## 项目编译

```bash
//linux
GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o treasure_user

//windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build


//mac
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build


```


### 其他

```sh
docker run  --rm \
    -w "/app" \
    --mount type=bind,source="D:\my-project\api-doc-go\backend",target=/app  \
    -p 2021:2021  \
    golang:alpine \
    sh -c  "go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && cd /app/user && go run main.go"
```

### 生成gin框架调用模型

进入 service/mall/cli 目录

1. 首先将config.example.toml改换成config.toml完善Mysql配置
2. 然后执行 `go run . -gen`

即可通过gin官方的gen工具生成模型到目录 `data/...` 下

## 使用docker启动

```shell
docker run -d --name treasure-doc \
--restart=always \
-p 2021:2021 \
-v /home/debian/project/treasure-doc/web:/app/web \
-v /home/debian/project/treasure-doc/files:/app/files \
-v /home/debian/project/treasure-doc/config.toml:/app/config.toml \
treasure-doc

docker run --rm --name treasure-doc \
-p 2021:2021 \
-v /home/debian/project/treasure-doc/web:/app/web \
-v /home/debian/project/treasure-doc/files:/app/files \
-v /home/debian/project/treasure-doc/config.toml:/app/config.toml \
treasure-doc

# 调试
docker run --rm --name treasure-doc -it \
-p 2021:2021 \
-v /home/debian/project/treasure-doc/web:/app/web \
-v /home/debian/project/treasure-doc/files:/app/files \
-v /home/debian/project/treasure-doc/config.toml:/app/config.toml \
treasure-doc /bin/sh 
```

## 构建docker镜像

```shell

docker build -t treasure-doc .

docker rm -f treasure-doc

docker save -o treasure-doc.tar.gz treasure-doc

sudo docker tag docker.linkedbus.com/golang:1.22.9-alpine3.20 golang:1.22.9-alpine3.20

sudo docker tag docker.linkedbus.com/alpine:latest alpine:latest


sudo chmod 777 treasure-doc.tar.gz

```

