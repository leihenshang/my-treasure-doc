ARG work_dir=/app
ARG binary_name=treasure_doc_go

FROM golang:alpine as builder
ARG work_dir
ARG binary_name

WORKDIR ${work_dir}
COPY . ${work_dir}

RUN go env -w GO111MODULE=on \
&& go env -w GOPROXY="https://goproxy.cn, direct"

RUN CGO_ENABLED=0 GOOS=linux go build -o ${binary_name} ${work_dir}


# ---------------------------------------------------------------
# setup 2 build
FROM alpine:latest as prod

#RUN apk --no-cache add ca-certificates

# define the expose port
EXPOSE 9091

# switch to the work dirctory
WORKDIR /app

# copy the binary excutable
COPY --from=0 /app/main .

# the command executed when the container is started 
CMD [ "/app/main","-addr=:9091" ]



