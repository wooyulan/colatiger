# 源镜像
FROM golang:1.20.6-alpine AS builder

RUN set -ex \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk --update add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime 



ADD . /app/
WORKDIR /app

# 安装依赖包
ENV GOPROXY https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download

# 把当前目录的文件拷过去，编译代码
COPY . .
RUN pwd && ls
RUN mkdir -p colatiger/ && CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o ./colatiger ./... 



FROM alpine

COPY --from=builder /app/colatiger /
COPY --from=builder /app/conf /conf



# 设置当前时区
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

EXPOSE 8000
ENV conf="conf/local.yaml"

ENTRYPOINT ["sh","-c","./cmd -conf=conf/$conf"]
