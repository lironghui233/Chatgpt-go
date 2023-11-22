# 编译阶段
FROM golang:1.20 as stage1
RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ADD ./ /src/chatgpt-official
WORKDIR /src/chatgpt-official
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatgpt-official .

FROM alpine:latest as stage2
ADD ./curl-amd64 /usr/bin/curl
RUN chmod +x /usr/bin/curl

MAINTAINER leoh
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage1 /src/chatgpt-official/chatgpt-official /app/chatgpt-official
# 指定入口程序
ENTRYPOINT ["./chatgpt-official"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml"]
