## CMD

docker build -f module/hot_top/Dockerfile  -t hot-top-go . 

docker run --rm --name hot-top-go -it -p 2025:2025 hot-top-go:latest 

docker save -o hot-top-go.tar.gz hot-top-go

docker save  hot-top-go:latest | gzip > hot-top-go.tar.gz

## PLAN
 1. each hot data source is controlled by a single timer for the collection.
 2. hot data write to database for other module to read.