ARG frontend_img
FROM $frontend_img as stage0

# 编译阶段
FROM golang:1.20 as stage1
RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ADD ./ /src/chatgpt-web-backend
WORKDIR /src/chatgpt-web-backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatgpt-web cmd/*.go

FROM alpine:latest as stage2
ADD ./curl-amd64 /usr/bin/curl
RUN chmod +x /usr/bin/curl

MAINTAINER leoh
WORKDIR /app/
ADD ./dev.config.yaml /app/config.yaml
COPY --from=stage1 /src/chatgpt-web-backend/chatgpt-web /app/chatgpt-web
COPY --from=stage0 /app/dist /app/frontend
# 指定入口程序
ENTRYPOINT ["./chatgpt-web"]
# 指定容器的启动命令或者入口程序的参数
CMD ["--config=config.yaml"]
