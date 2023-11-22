# tokenizer
## docker build
```
docker build -t tokenizer:latest .
```
## docker 运行
```
# 可通过环境变量，动态指定应用程序监听的端口
docker run -d -p 3003:3003 -e PORT=3003 --restart=always tokenizer:latest
```
## docker service 部署
```
docker service create --name tokenizer -p 3002:3002 \
--replicas 3 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=2 \
--with-registry-auth \
tokenizer:0.1.0
```
