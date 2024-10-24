FROM alpine:3.15

# RUN apk add --no-cache <必要的包>

ARG MODULE_NAME
# 使用构建参数进行复制
COPY build/${MODULE_NAME} /app/webook

WORKDIR /app

ENTRYPOINT ["/app/webook"]