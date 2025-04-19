FROM golang:1.24.2-alpine AS deb_builder

WORKDIR /app
COPY . .
RUN export GOPROXY=https://goproxy.cn \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/' /etc/apk/repositories \
    && go build -o eta \
    && go run main.go stats > /app/.env

FROM golang:1.24.2-alpine

# 镜像commit信息
ARG COMMIT_TAG
ARG COMMIT_HASH
ENV COMMIT_TAG=$COMMIT_TAG
ENV COMMIT_HASH=$COMMIT_HASH

#ARG TARGETPLATFORM=linux/amd64
ARG TARGET_ARCH=x86_64
#ENV TARGETPLATFORM=$TARGETPLATFORM
ENV TARGET_ARCH=$TARGET_ARCH

WORKDIR /app

COPY --from=deb_builder /app/eta /app/eta
COPY --from=deb_builder /app/.env /app/.env
COPY deploy/docker/mime.types /etc/mime.types

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/' /etc/apk/repositories \
    && apk update \
    && apk --no-cache add curl tini \
    && rm -rf /var/cache/apk/*

#ENTRYPOINT ["/app/eta"]
ENTRYPOINT ["tini", "--", "/app/eta"]
CMD ["server"]
