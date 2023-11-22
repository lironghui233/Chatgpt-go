# keywords
关键词查找

## 生成代码
```
protoc.exe --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative .\proto\keywords.proto 
```


## 初始化词库
```
go run .\keywords-server\main.go --config=dev.config.yaml --dict=dict.txt --init-dict=true
```


## 部署

### docker镜像构建
```
docker build -t keywords:0.1.0 .
```

## 创建配置
```
docker config create --label env=prod keywords-conf config.yaml

docker config create --label env=prod keywords-dict dict.txt
```


## 部署服务
```
docker service create --name keywords -p 50054:50054 \
--config src=keywords-conf,target=/app/config.yaml \
--config src=keywords-dict,target=/app/dict.txt \
--replicas 3 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=2 \
--health-cmd "grpc_health_probe -addr=:50054" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
localhost:5000/keywords:0.1.0
```