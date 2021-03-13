# 设置环境变量
export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export GOPROXY="https://goproxy.cn,direct"
export GIN_MODE=release

# 创建文件夹
mkdir -p ./_bulid/data
mkdir -p ./_bulid/cache
mkdir -p ./_bulid/conf/pro
mkdir -p ./_bulid/script/mysql


# 赋值文件至编译缓存区
cp ./conf/pro/app.conf.toml ./_bulid/conf/pro/app.conf.toml
cp ./script/mysql/godb.sql ./_bulid/script/mysql/godb.sql


# 编译文件
go build -o gobang

# 移动编译的文件至指定位置
mv gobang ./_bulid