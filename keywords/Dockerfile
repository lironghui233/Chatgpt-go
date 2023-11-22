# 编译阶段
FROM golang:1.20 as stage0
RUN go env -w GOPROXY=https://proxy.golang.com.cn,https://goproxy.cn,direct
ADD ./ /src/keywords
WORKDIR /src/keywords
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o keywords ./keywords-server

FROM alpine:latest as stage1
ADD ./grpc_health_probe-linux-amd64 /usr/bin/grpc_health_probe
RUN chmod +x /usr/bin/grpc_health_probe

MAINTAINER leoh
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
ADD ./dict.txt /app/dict.txt
COPY --from=stage0 /src/keywords/keywords ./
# 指定入口程序
ENTRYPOINT ["./keywords"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml","--dict=dict.txt"]
