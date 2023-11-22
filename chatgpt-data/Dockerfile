# 编译阶段
FROM golang:1.20 as stage0
RUN go env -w GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct
ADD ./ /src/chatgpt-data
WORKDIR /src/chatgpt-data
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatgpt-data ./chatgpt-data-server

FROM alpine:latest as stage1
# grpc 健康检查探针
ADD ./grpc_health_probe-linux-amd64 /usr/bin/grpc_health_probe
RUN chmod +x /usr/bin/grpc_health_probe

MAINTAINER leoh
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage0 /src/chatgpt-data/chatgpt-data ./
# 指定入口程序
ENTRYPOINT ["./chatgpt-data"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml"]
