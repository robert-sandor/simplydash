FROM golang:1.19 AS backendBuild

WORKDIR /app

RUN pwd

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./cmd /app/cmd
COPY ./internal /app/internal

RUN go build -o /app/simplydash ./cmd

FROM alpine:3

RUN apk add libc6-compat curl

WORKDIR /app
COPY --from=backendBuild /app/simplydash ./simplydash

VOLUME /app/config /app/icons

EXPOSE 8080

ENTRYPOINT ["/app/simplydash", "--config", "/app/config/config.yml", "--icon-cache", "/app/icons"]

HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://localhost:8080/health || exit 1
