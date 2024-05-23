FROM golang:1.21 as builder

LABEL org.opencontainers.image.source = "https://github.com/gagraler/alert-service"

COPY . /tmp/alert-service

WORKDIR /tmp/alert-service/

ENV GOPROXY http://goproxy.cn

RUN go mod tidy

RUN go build -o alert-service ./cmd

FROM alpine:3.14

RUN apk update && \
    apk add --no-cache ca-certificates

RUN mkdir -p /opt/alert-service/config

WORKDIR /opt/alert-service

COPY --from=builder /tmp/alert-service/alert-service /opt/alert-service/alert-service
COPY --from=builder /tmp/alert-service/conf.d/alert-service.toml /opt/alert-service/conf.d/

EXPOSE 8588

ENTRYPOINT ["./alert-service"]
