# 编译阶段
FROM golang:1.20 as stage0
RUN go env -w GOPROXY=https://goproxy.cn,https://proxy.golang.com.cn,direct
ADD ./ /src/chatgpt-qq
WORKDIR /src/chatgpt-qq
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chatgpt-qq .

FROM alpine:latest as stage1
# RUN apk update && apk add curl
ADD ./curl-amd64 /usr/bin/curl
RUN chmod +x /usr/bin/curl
MAINTAINER leoh
WORKDIR /app/
ADD ./dev.config.cfg /app/config.cfg
COPY --from=stage0 /src/chatgpt-qq/chatgpt-qq ./
# 指定入口程序
ENTRYPOINT ["./chatgpt-qq"]
