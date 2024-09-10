FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

RUN export GO111MODULE=on && \
    export GOPROXY=https://goproxy.cn,direct && \
    go build -o alist-gallery

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/alist-gallery .

EXPOSE 5243
CMD ["/app/alist-gallery"]