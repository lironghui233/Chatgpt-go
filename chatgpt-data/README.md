# chatgpt-data
## 通过虚拟终端与容器交互
```
docker exec -it mysql8.0.32 bash
```
## 命令行连接到数据
```
mysql -h [host] -P [port] -u root --default-character-set=utf8mb4 -p
```
## protobuf 编译器下载
```
https://github.com/protocolbuffers/protobuf/releases/tag/v22.0
```
## protobuf golang插件安装
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
## protoc 生成代码
```
protoc.exe --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\proto\chatgpt_data.proto
```
## 依赖对齐
```
go mod tidy
```
# 证书生成

## 自签根证书(替代权威机构的证书)
```
openssl req -x509                                                       \
  -newkey rsa:4096                                                      \
  -nodes                                                                \
  -days 3650                                                            \
  -keyout ca_key.pem                                                    \
  -out ca_cert.pem                                                      \
  -subj /C=CN/ST=hunan/L=changsha/O=0voice/OU=teacher/CN=www.0voice.com \
  -sha256
```

## 生成密钥
```
openssl genrsa -out server_key.pem 4096
```

## 基于私钥生成证书（公钥）
```
openssl req -new  \
-key server_key.pem \
-out server_csr.pem \
-subj /C=CN/ST=hunan/L=changsha/O=0voice/OU=teacher/CN=*.grpc.0voice.com
```

## 为证书 颁发机构签署的证书
```
openssl x509 -req -in server_csr.pem -CA ca_cert.pem -CAkey ca_key.pem -CAcreateserial -out server_cert.pem -days 3650
```

## 校验证书的的有效性
```
openssl verify -verbose -CAfile ca_cert.pem  server_cert.pem
```

## 课堂上gRPC通讯安全使用的证书，均由create.sh生成

## grpc 健康检查客户端下载
```
https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.18/grpc_health_probe-linux-amd64
```

## docker 镜像构建
```
docker build -t chatgpt-data:0.1.0 .
```

## docker service 部署

### 将配置文件创建为资源
```
docker config create --label env=prod chatgpt-data-conf dev.config.yaml
```

### docker service 创建服务
```
docker service create --name chatgpt-data -p 50052:50052 \
--config src=chatgpt-data-conf,target=/app/config.yaml \
--replicas 3 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=2 \
--health-cmd "grpc_health_probe -addr=:50052" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
chatgpt-data:0.1.0
```

## 单元测试
```
go test  .\test-case\. --config=../dev.config.yaml -v
```