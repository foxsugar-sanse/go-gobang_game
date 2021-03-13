FROM apline

# 为镜像设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://goproxy.cn,direct"

# 复制文件
COPY ./conf/pro/ /app/conf/pro
COPY ./script/mysql/ /app/script/mysql
COPY ./_build/ /app

# 容器的工作目录
WORKDIR ./app

# 创建数据目录
RUN mkdir data

# 对外服务的端口
EXPOSE 8080

# 启动容器时打开的命令
CMD ["./gobang"]
