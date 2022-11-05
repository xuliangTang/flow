FROM golang:1.19-alpine

RUN mkdir /app
ADD . /app
WORKDIR /app
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.cn"
RUN go mod tidy
RUN go build -o run ./src


FROM alpine:3.12
RUN set -ex; \
    apk update; \
    apk upgrade; \
    apk add --no-cache supervisor

WORKDIR /app
COPY --from=0 /app/run /app

RUN mkdir -p storage/logs

EXPOSE 8080

ENTRYPOINT ["/app/run"]