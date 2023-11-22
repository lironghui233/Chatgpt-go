# chatgpt-qq QQ协议解析端

## docker 镜像构建
```
docker build -t chatgpt-qq:0.1.0 .
```

## docker service 服务部署

### 配置文件创建
```
docker config create --label env=prod chatgpt-qq-conf config.cfg
```

### docker service 部署服务
```
docker service create --name chatgpt-qq -p 8080:8080 \
--config src=chatgpt-qq-conf,target=/app/config.cfg \
--replicas 3 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=2 \
--health-cmd "curl -f http://localhost:8080/health" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
chatgpt-qq:0.1.0
```
