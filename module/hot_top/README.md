docker build -f module/hot_top/Dockerfile  -t hot-top-go . 

docker run --rm --name hot-top-go -it -p 2025:2025 hot-top-go:latest 