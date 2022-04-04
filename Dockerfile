# Build stage
FROM golang:1.18-alpine3.15 AS base
WORKDIR /app
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go build -o main main.go

# Run stage
FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=base /app/app.env /app/main ./
EXPOSE 9000
CMD ["/app/main"]