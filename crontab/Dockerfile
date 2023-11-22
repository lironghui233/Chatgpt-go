# 编译阶段
FROM golang:1.20 as stage0
RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ADD ./ /src/crontab
WORKDIR /src/crontab
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o crontab ./token-server

FROM alpine:latest as stage1
ADD ./grpc_health_probe-linux-amd64 /usr/bin/grpc_health_probe
RUN chmod +x /usr/bin/grpc_health_probe

MAINTAINER leoh
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage0 /src/crontab/crontab ./
# 指定入口程序
ENTRYPOINT ["./crontab"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml"]
