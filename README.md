# ChatGPT 微服务生态系统

这是一个基于微服务架构的 ChatGPT 应用生态系统，提供多平台接入能力和丰富的功能支持。

## 服务模块说明

### 1. chatgpt-data

主要功能：

- 数据持久化服务
- 对话历史存储
- 用户信息管理
- 配置信息存储

技术栈：

- Go
- MySQL/PostgreSQL
- Redis
- gRPC

主要依赖：

- GORM
- Redis-go
- gRPC 相关包
- Viper

### 2. chatgpt-proxy

主要功能：

- API 代理服务
- 负载均衡
- 请求限流
- 访问控制
- 服务监控

技术栈：

- Go
- Nginx
- Redis
- Prometheus

主要依赖：

- gorilla/mux
- redis-go
- prometheus-client
- ratelimit

### 3. chatgpt-service

主要功能：

- ChatGPT 核心业务逻辑
- 会话管理
- 模型调用封装
- 多平台适配接口

技术栈：

- Go
- gRPC
- Redis
- Event Bus

主要依赖：

- gRPC
- Redis
- NATS/RabbitMQ
- protobuf

### 4. chatgpt-stack

主要功能：

- 技术栈管理服务
- 依赖包版本控制
- 服务部署配置
- 监控告警集成

技术栈：

- Docker
- Kubernetes
- Prometheus
- Grafana

主要依赖：

- Docker Compose
- Kubernetes
- Helm
- Prometheus

### 5. chatgpt-web-backend

主要功能：

- Web 后端服务
- API 网关
- 用户认证授权
- 业务逻辑处理

技术栈：

- Go
- RESTful API
- JWT
- MySQL

主要依赖：

- Gin
- GORM
- JWT-go
- Swagger

### 6. chatgpt-web-frontend-leoh

主要功能：

- Web 前端界面
- 多会话管理
- Markdown 渲染
- 代码高亮
- 移动端适配
- 多语言支持
- 主题切换

技术栈：

- Vue.js
- TypeScript
- Vite
- TailwindCSS

主要依赖：

- Vue 3
- Vite
- markdown-it
- highlight.js
- TailwindCSS

### 7. chatgpt-wechat

主要功能：

- 微信消息处理
- 微信群聊处理
- 用户管理
- 自动回复

技术栈：

- Go
- WeChat SDK
- Redis
- MySQL

主要依赖：

- wechat-sdk-go
- GORM
- Redis
- gRPC
- Gin

### 8. chatgpt-wecom

主要功能：

- 企业微信集成
- 企业应用管理
- 员工消息处理
- 群聊集成

技术栈：

- Go
- WeCom SDK
- Redis
- MySQL

主要依赖：

- wecom-sdk-go
- GORM
- Redis
- gRPC
- Gin

### 9. chatgpt-official

主要功能：

- 微信公众号集成
- 用户管理
- 自动回复

技术栈：

- Go
- gRPC
- YAML 配置
- JS-SDK

主要依赖：

- JS-SDK
- gRPC
- Viper
- logrus
- Gin

### 10. chatgpt-qq

主要功能：
- QQ 机器人集成服务
- 连接 ChatGPT 服务
- 提供 QQ 消息转发能力

技术栈：
- Go
- WebSocket
- gRPC

主要依赖：
- gorilla/websocket
- gRPC 相关包
- Viper (配置管理)
- Gin

### 11. crontab

主要功能：

- 定时任务管理
- 任务调度
- 失败重试
- 任务监控

技术栈：

- Go
- Redis
- MySQL
- gRPC

主要依赖：

- cron
- GORM
- Redis
- prometheus

### 12. keywords

主要功能：

- 关键词管理
- 敏感词过滤
- 关键词匹配
- 词频统计

技术栈：

- Go
- Trie树
- Redis
- MySQL

主要依赖：

- GORM
- Redis
- Bloom Filter
- gRPC

### 13. sensitive-words

主要功能：

- 敏感词库管理
- 内容审核
- 实时过滤
- 规则配置

技术栈：

- Go
- DFA算法
- Redis
- MySQL

主要依赖：

- GORM
- Redis
- gRPC
- Bloom Filter

### 14. tokenizer

主要功能：

- Token计数
- 多语言支持
- 编码转换

技术栈：

- Go
- NLP
- Unicode
- gRPC

主要依赖：

- tiktoken
- unicode
- protobuf
- gRPC

## 架构图

![架构图（修订）](D:\project\go\chatgpt\image\README\架构图（修订）.png)

## 系统架构

系统采用微服务架构，主要分为以下几层：

1. 代理层：chatgpt-proxy
2. 应用层：chatgpt-web-backend, chatgpt-wechat, chatgpt-wecom, chatgpt-qq, chatgpt-official
3. 核心服务层： chatgpt-service
4. 数据层：chatgpt-data
5. 基础设施层：chatgpt-stack
6. 工具服务层：keywords, sensitive-words, tokenizer, crontab

## 注意事项

**架构设计：**

需要对外提供 HTTP 接口的服务使用 Gin

内部服务间通信使用 gRPC

## 部署要求

- Go 1.16+
- Node.js 16+
- Docker 20.10+
- Kubernetes 1.20+
- MySQL 8.0+
- Redis 6.0+
- gRPC
- Nginx

## 监控告警

系统使用 Prometheus + Grafana 进行监控，主要监控指标：

- API 请求量和延迟
- 服务健康状态
- 资源使用情况
- 错误率统计
- 业务指标

## 配置管理

所有服务均支持：

- YAML 配置文件
- 环境变量覆盖
- 动态配置更新
- 多环境配置

配置示例见各服务的 config.yaml 文件。



