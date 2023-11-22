# chatgpt-service
内部chatgpt服务，收纳所有相关业务

## gRPC 代码生成
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\proto\chatgpt.proto
```

## 依赖对齐
```
go mod tidy
```

## 单元测试
```
go test .\test-case\. --config=../dev.config.yaml -v 
```


# docker 镜像构建
```
docker build -t chatgpt-service:0.1.0 .
```

# docker service 部署服务

## 创建配置文件资源
```
docker config create --label env=prod chatgpt-service-conf config.yaml
```

## 创建并启动服务
```
docker service create --name chatgpt-service -p 50051:50051 \
--config src=chatgpt-service-conf,target=/app/config.yaml \
--replicas 3 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=2 \
--health-cmd "grpc_health_probe -addr=:50051" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
localhost:5000/chatgpt-service:0.1.0
```
