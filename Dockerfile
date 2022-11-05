FROM golang:1.19-alpine

RUN set -ex; \
    apk update; \
    apk upgrade; \
    apk add --no-cache supervisor

WORKDIR /app
COPY . /app

RUN mkdir -p storage/logs

ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"
RUN go mod tidy
RUN go build -o run ./src

EXPOSE 80

ENTRYPOINT ["/app/run"]