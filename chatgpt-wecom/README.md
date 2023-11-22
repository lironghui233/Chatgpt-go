# chatgpt-official

## xml tag

在 Go 的 XML 包中，有一些常用的 XML tag 可以用于 struct 字段上，这些 tag 可以控制字段如何编码和解码。下面是其中一些常见的 tag：
- xml:"name"：将字段映射到名为 name 的 XML 元素。
- xml:"-"：忽略此字段。
- xml:",attr"：将字段作为属性编码到父元素中。
- xml:",chardata"：将字段的值编码为父元素的文本内容。
- xml:",innerxml"：将原始 XML 编码为此字段的字符串。
- xml:",comment"：将注释编码为此字段的字符串。

另外还有一些其他选项可以与这些 tag 结合使用，例如：
- omitempty：如果该选项设置，则在编码时忽略空值或默认值（例如 0 或 false）。
- -（连字符）表示该字段不会被解析/序列化

以上是部分比较常见的XML tag，在实际开发中根据具体需求可能会用到其他tag。


## docker 镜像构建
```
docker build -t chatgpt-wecom:0.1.0 .
```

# 部署服务

## 创建配置文件
```
docker config create --label env=prod chatgpt-wecom-conf config.yaml
```

## docker service 部署服务
```
docker service create --name chatgpt-wecom -p 7082:7082 \
--config src=chatgpt-wecom-conf,target=/app/config.yaml \
--replicas 2 --limit-cpu 0.3 --reserve-cpu 0.1 \
--update-parallelism=1 \
--health-cmd "curl -f http://localhost:7082/api/health" \
--health-interval 5s --health-retries 3 \
--with-registry-auth \
localhost:5000/chatgpt-wecom:0.1.0
```
