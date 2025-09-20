ARG WORK_DIR=/app
ARG BINARY_NAME=hot-top-go
ARG EXPOSE_PORT=2025
ARG BUILD_DIR=module/hot_top

# if docker image cannot pull, refers https://cloud.tencent.com/developer/article/2454486
#backup image address: docker.linkedbus.com/
FROM golang:1.23.4-alpine3.21 AS builder
ARG WORK_DIR
ARG BINARY_NAME
ARG BUILD_DIR

WORKDIR ${WORK_DIR}
COPY . ${WORK_DIR}

ENV GOPROXY=https://goproxy.cn,direct
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    #    && apk add build-base \
    #   && apk --no-cache add ca-certificates \
    &&  CGO_ENABLED=0  go build -o ${WORK_DIR}/${BUILD_DIR}/${BINARY_NAME} ${WORK_DIR}/${BUILD_DIR}


# setup 2 build
FROM alpine:latest AS prod
ARG WORK_DIR
ARG BINARY_NAME
ARG BUILD_DIR

# define the expose port
EXPOSE ${EXPOSE_PORT}

WORKDIR ${WORK_DIR}

# copy the binary executable
COPY --from=builder ${WORK_DIR}/${BUILD_DIR}/${BINARY_NAME} ${WORK_DIR}/hot-top-go
COPY --from=builder ${WORK_DIR}/${BUILD_DIR}/config.toml.example ${WORK_DIR}/config.toml
COPY --from=builder ${WORK_DIR}/${BUILD_DIR}/dist ${WORK_DIR}/dist

# the command executed when the container is started
CMD ["/app/hot-top-go"]


