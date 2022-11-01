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

//减少包大小 -ldflags -s 去掉符号信息 -w 去掉调试信息 
go build -ldflags "-s -w" -o main-ldflags main.go 

```

## 目录说明

```
# 待更新
```

### 计划

- [x] 加入 `gin` http框架,创建main.go
- [x] 添加配置解析库 `viper` [github地址](https://github.com/spf13/viper)
- [x] 添加日志库 `zap` [github地址](https://github.com/uber-go/zap)
- [x] 添加orm库 `gorm` [github地址](https://github.com/go-gorm/gorm)
- [x] 添加redis库 `go-redis` [github地址](https://github.com/go-redis/redis)

### 其他

```sh
docker run  --rm \
    -w "/app" \
    --mount type=bind,source="D:\my-project\api-doc-go\backend",target=/app  \
    -p 2021:2021  \
    golang:alpine \
    sh -c  "go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct && cd /app/user && go run main.go"
```

### 生成模型

`go install gorm.io/gen/tools/gentool@latest`

```bash
gentool -h

Usage of gentool:
 -db string
       input mysql or postgres or sqlite or sqlserver. consult[https://gorm.io/docs/connecting_to_the_database.html] (default "mysql")
 -dsn string
       consult[https://gorm.io/docs/connecting_to_the_database.html]
 -fieldNullable
       generate with pointer when field is nullable
 -fieldWithIndexTag
       generate field with gorm index tag
 -fieldWithTypeTag
       generate field with gorm column type tag
 -modelPkgName string
       generated model code's package name
 -outFile string
       query code file name, default: gen.go
 -outPath string
       specify a directory for output (default "./dao/query")
 -tables string
       enter the required data table or leave it blank
 -onlyModel
       only generate models (without query file)
 -withUnitTest
       generate unit test for query code
 -fieldSignable
       detect integer field's unsigned type, adjust generated data type

```

example

```bash

gentool -dsn "user:pwd@tcp(localhost:3306)/database?charset=utf8mb4&parseTime=True&loc=Local" -tables "orders,doctor"


```