# chatgpt-wechat

## golang 个人微信SDK
地址（https://github.com/eatmoreapple/openwechat）
文档（https://openwechat.readthedocs.io/zh/latest/）

## 镜像构建
```
docker build -t chatgpt-wechat:0.1.0 .
```

## 部署
```
docker run -d --name chatgpt-wechat01 chatgpt-wechat:0.1.0 
docker run -d -v /home/leoh/work/chatgpt-wechat/dev.config.yaml:/app/config.yaml --name chatgpt-wechat01 chatgpt-wechat:0.1.0
# 查看容器日志输出，获取首次扫码登录链接
docker logs chatgpt-wechat01  
```
