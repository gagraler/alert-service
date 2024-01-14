FROM golang:1.21 as builder

COPY . /tmp/alertService

WORKDIR /tmp/alertService/

ENV GOPROXY http://goproxy.cn

RUN go mod vendor

RUN make build

FROM alpine:3.14

RUN apk update && \
    apk add --no-cache ca-certificates

RUN mkdir -p /opt/alertService/config

WORKDIR /opt/alertService

COPY --from=builder /tmp/alertService/alertService_* /opt/alertService/alertService
COPY --from=builder /tmp/alertService/conf.d/alertService.toml /opt/alertService/conf.d/

EXPOSE 8588

ENTRYPOINT ["./alertService"]
