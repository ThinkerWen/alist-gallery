FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev
RUN export GO111MODULE=on && \
    export GOPROXY=https://goproxy.cn,direct && \
    go build -ldflags='-s -w -extldflags "-static"' -o alist-gallery

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/alist-gallery .
COPY --from=builder /app/sync.sh .
RUN chmod +x /app/sync.sh

EXPOSE 5243
CMD ["/bin/sh", "-c", "/app/sync.sh && /app/alist-gallery"]