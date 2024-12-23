FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev
RUN export CGO_ENABLED=1 && \
    export GO111MODULE=on && \
    export GOPROXY=https://goproxy.cn,direct && \
    go build -ldflags='-s -w -extldflags "-static"' -o alist-gallery

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/alist-gallery .

EXPOSE 5243
CMD ["/app/alist-gallery"]