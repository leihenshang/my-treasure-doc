ARG work_dir=/app
ARG binary_name=treasure-doc
ARG expose_port=2021
ARG build_dir=service/user

# if docker image canot pull,refers https://cloud.tencent.com/developer/article/2454486
#docker.linkedbus.com/
FROM docker.linkedbus.com/golang:1.22.9-alpine3.20 AS builder
ARG work_dir
ARG binary_name
ARG build_dir

WORKDIR ${work_dir}
COPY . ${work_dir}

ENV GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    #    && apk add build-base \
    #   && apk --no-cache add ca-certificates \
    &&  CGO_ENABLED=0  go build -o ${work_dir}/${build_dir}/${binary_name} ${work_dir}/${build_dir}


# setup 2 build
FROM docker.linkedbus.com/alpine:latest AS prod
ARG work_dir
ARG binary_name
ARG build_dir

# define the expose port
EXPOSE ${expose_port}

WORKDIR ${work_dir}

# copy the binary excutable
COPY --from=builder ${work_dir}/${build_dir}/${binary_name} ${work_dir}
COPY --from=builder ${work_dir}/${build_dir}/config.example.toml ${work_dir}/config.toml
COPY --from=builder ${work_dir}/${build_dir}/files ${work_dir}/files
COPY --from=builder ${work_dir}/${build_dir}/web ${work_dir}/web

# the command executed when the container is started
CMD ["/app/treasure-doc"]


