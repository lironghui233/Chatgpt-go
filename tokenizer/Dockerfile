FROM python:3.10-alpine

# 在python:3.10-alpine环境下，安装curl
RUN apk update && apk add curl

# 指定工作目录
WORKDIR app

# 环境变量，自定义变量
ENV PORT 3002

ADD tokenizer.py /app/tokenizer.py
ADD requirements.txt /app/requirements.txt
RUN pip install -i https://mirrors.aliyun.com/pypi/simple --upgrade pip
RUN pip install --root-user-action=ignore -i https://mirrors.aliyun.com/pypi/simple -r requirements.txt

#健康检查
HEALTHCHECK --interval=5s --timeout=3s --retries=5 \
CMD curl -f -X POST -H "Content-Type: application/json" -d "{\"role\":\"user\",\"content\":\"hello\"}" http://localhost:${PORT}/tokenizer/gpt-3.5-turbo-0301 || exit 1

CMD ["sh", "-c", "nuxt --port ${PORT} --module tokenizer.py --workers 2"]