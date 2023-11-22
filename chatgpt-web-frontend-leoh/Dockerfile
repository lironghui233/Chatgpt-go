FROM node:lts-alpine AS frontend
RUN npm install pnpm -g
COPY . ./app
WORKDIR /app
RUN pnpm bootstrap
RUN pnpm run build-only

FROM nginx:1.23.4 as web
COPY --from=frontend /app/dist/ /usr/share/nginx/html/
# 指定入口程序
ENTRYPOINT ["/docker-entrypoint.sh"]
# 指定容器的启动命令或者入口程序的参数
CMD ["nginx", "-g", "daemon off;"]
