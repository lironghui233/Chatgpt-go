# 编译阶段
FROM golang:1.20 as stage1
RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ADD ./ /src/chatgpt-wecom
WORKDIR /src/chatgpt-wecom
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatgpt-wecom .

FROM alpine:latest as stage2
ADD ./curl-amd64 /usr/bin/curl
RUN chmod +x /usr/bin/curl

MAINTAINER leoh
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage1 /src/chatgpt-wecom/chatgpt-wecom /app/chatgpt-wecom
# 指定入口程序
ENTRYPOINT ["./chatgpt-wecom"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml"]
