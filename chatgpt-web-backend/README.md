# 构建流程

## 1.构建前端项目docker镜像
```
# 进入chatgpt-web-frontend 项目根目录构建前端镜像
docker build -t chatgpt-web-frontend:0.1.0 .
```

## 2. 将前端镜像作为参数，传递到后端镜像的构建过程
后端依赖前端镜像：
```
docker build -t chatgpt-web:0.1.0  --build-arg "frontend_img=chatgpt-web-frontend:0.1.0" .
```
前后端分离后，后端不依赖前端镜像，前端直接部署在nginx
```
docker build -t chatgpt-web:0.1.0  .
```

# docker service 部署

## 创建配置文件资源
```
docker config create --label env=prod chatgpt-web-conf config.yaml
```

## 部署服务
```
docker service create --name chatgpt-web -p 7080:7080 \
--config src=chatgpt-web-conf,target=/app/config.yaml \
--replicas 3 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=2 \
--health-cmd "curl -f http://localhost:7080/api/health" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
chatgpt-web:0.1.0
```
