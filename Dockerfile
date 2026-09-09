ARG work_dir=/app
ARG binary_name=treasure-doc
ARG expose_port=2026
ARG build_dir=module/user

# ---- 后端构建 ----
FROM golang:1.26.7-alpine3.23 AS builder
ARG work_dir
ARG binary_name
ARG build_dir

ENV GOPROXY=https://goproxy.cn,direct
ENV CGO_ENABLED=0
WORKDIR ${work_dir}

# 先只拷贝依赖清单并下载，利用镜像层缓存：
# 仅当 go.mod / go.sum 变化时才重新拉取依赖，源码改动不会触发重下。
COPY go.mod go.sum ./
RUN go mod download

# 再拷贝全部源码并构建
COPY . ${work_dir}
RUN go build -o ${work_dir}/${build_dir}/${binary_name} ${work_dir}/${build_dir}


# ---- 运行镜像 ----
FROM alpine:3.23 AS prod
ARG work_dir
ARG binary_name
ARG build_dir

# 定义暴露端口
EXPOSE ${expose_port}

WORKDIR ${work_dir}

# 拷贝可执行文件、配置与静态资源
COPY --from=builder ${work_dir}/${build_dir}/${binary_name} ${work_dir}
COPY --from=builder ${work_dir}/${build_dir}/config.example.toml ${work_dir}/config.toml
COPY --from=builder ${work_dir}/${build_dir}/files ${work_dir}/files
COPY --from=builder ${work_dir}/${build_dir}/web ${work_dir}/web

# 容器启动时执行的命令
CMD ["/app/treasure-doc"]
