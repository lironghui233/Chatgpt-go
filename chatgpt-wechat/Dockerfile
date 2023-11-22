# 编译阶段
FROM golang:1.20 as stage1
RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ADD ./ /src/chatgpt-wechat
WORKDIR /src/chatgpt-wechat
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatgpt-wechat .

FROM alpine:latest as stage2
MAINTAINER leoh
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage1 /src/chatgpt-wechat/chatgpt-wechat /app/chatgpt-wechat
# 指定入口程序
ENTRYPOINT ["./chatgpt-wechat"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml"]
