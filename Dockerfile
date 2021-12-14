FROM golang:1.16-alpine AS builder

# 设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的所有文件
COPY . .

# 执行redis连接的可用性
RUN go test -run "^TestRedisPing$" -v ./utils

# 将代码编译成二进制可执行文件
RUN go build -tags=jsoniter,nomsgpack -o web-server .

###################
# 创建一个小镜像
# 镜像中是不需要go编译器的，通过多阶段构建仅保留二进制文件来减小镜像大小
###################
FROM debian:stretch-slim

# 复制必要的静态文件和配置文件
#COPY ./conf /conf
#COPY ./static /static

# 创建放置文件的目录
RUN mkdir /data /data/logs

WORKDIR /data

# 从builder镜像中把二进制文件拷贝到当前目录
COPY --from=builder /build/web-server .

RUN set -eux; \
    apt-get update; \
    chmod +x ./web-server

# 声明服务端口
EXPOSE 8000

# 需要运行的命令
ENTRYPOINT ["./web-server"]

# 默认的启动命令只输出帮助信息
# 通过在`docker run`命令尾部指定命令行参数来覆盖该参数
CMD ["-h"]
