# 构建阶段
FROM golang:alpine AS builder

# 构建可执行文件
ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

WORKDIR /build

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制所有源代码
COPY . .

# 编译应用（指定入口文件）
RUN go build -o main cmd/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装运行时依赖
RUN apk --no-cache add ca-certificates tzdata

# 从构建阶段复制文件
COPY --from=builder /build/main .
COPY --from=builder /build/config ./config

EXPOSE 3000

CMD ["./main"]