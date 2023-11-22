# chatgpt-proxy

## 单元测试
### 功能测试
``` 
go test ./test-case/proxy_test.go  --config=../config.yaml --addr=http://localhost:4002/v1 -v
```
### 性能测试
```
go test -run ^$ -bench . -benchmem test-case/*.go --config=../config.yaml --addr=http://localhost:4002/v1
```
### 性能测试flag 使用
```
go test -run ^$ -bench . -benchtime 2x -count 2 -cpu 1,2 -cover -benchmem test-case/*.go --config=../config.yaml --addr=http://localhost:4002/v1
```

### 性能测试指标的导出
```
go test -run ^$ -bench . -benchtime 2x -count 2 -cpu 1,2 -cover -benchmem \
-blockprofile block.out -cpuprofile cpu.out -memprofile mem.out \
-mutexprofile mutex.out -outputdir ./testout \
test-case/*.go --config=../config.yaml --addr=http://localhost:4002/v1
```
## docker镜像构建
``` 
docker build -t chatgpt-proxy:latest -t chatgpt-proxy:0.1.0 .
```
## docker swarm 服务部署
### 创建配置文件
```
docker config create --label env=prod chatgpt-proxy-conf config.yaml
```
### 服务部署
```
docker service create --name chatgpt-proxy -p 4002:4002 \
--config src=chatgpt-proxy-conf,target=/app/config.yaml \
--replicas 5 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=2 \
--health-cmd "curl -f http://localhost:4002/health" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
localhost:5000/chatgpt-proxy:0.1.2
```
### 服务更新
```
# 如有需要，创建新的配置文件
docker config create --label env=prod  chatgpt-proxy-confv1 config.yaml
# 同时更新镜像和配置
docker service update chatgpt-proxy --image localhost:5000/chatgpt-proxy:0.1.3 \
--config-rm chatgpt-proxy-conf \
--config-add src=chatgpt-proxy-confv1,target=/app/config.yaml
```
### 查看服务的更新记录
``` 
docker service ps chatgpt-proxy
```
### 服务回滚
```
docker service update chatgpt-proxy --rollback
```
### 服务伸缩
```
docker service scale chatgpt-proxy=3
```