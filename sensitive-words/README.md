# sensitive
敏感词过滤

## gRPC 代码生成
```
protoc.exe --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\proto\sensitive.proto
```

## 初始化敏感词库
```
go run .\sensitive-server\main.go --config=dev.config.yaml --dict=dict.txt --init-dict=true
```


## 单元测试
```
go test  .\test-case\. --config=../dev.config.yaml -v
```

## 镜像构建
```
docker build -t sensitive-words:0.1.0 .
```

## docker service 部署

### 创建配置
```
docker config create --label env=prod sensitive-conf config.yaml

docker config create --label env=prod sensitive-dict dict.txt
```


### 部署服务
```
docker service create --name sensitive-words -p 50053:50053 \
--config src=sensitive-conf,target=/app/config.yaml \
--config src=sensitive-dict,target=/app/dict.txt \
--replicas 3 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=2 \
--health-cmd "grpc_health_probe -addr=:50053" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
localhost:5000/sensitive-words:0.1.0
```