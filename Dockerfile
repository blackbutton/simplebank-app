# Build stage
FROM golang:1.18-alpine3.15 AS base
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go build -o main main.go && \
    set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk add curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine
WORKDIR /app
RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates && \
    apk add -U tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    apk del tzdata

COPY --from=base /app/app.env /app/main ./
COPY --from=base /app/migrate ./migrate
COPY db/migrations ./migration
COPY start.sh ./
COPY wait-for.sh ./

EXPOSE 9000
# CMD ["/app/main"]
# ENTRYPOINT ["/app/start.sh"]