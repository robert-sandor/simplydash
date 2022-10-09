FROM golang:1.19 AS backendBuild

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download

COPY ./cmd /app/cmd
COPY ./internal /app/internal

RUN go build -o /app/simplydash ./cmd

FROM node:18 as frontendBuild

WORKDIR /app

COPY ./web .

RUN npm install
RUN npm run build

FROM alpine:3

RUN apk add libc6-compat curl

WORKDIR /app
COPY --from=backendBuild /app/simplydash ./simplydash
COPY --from=frontendBuild /app/dist ./web

VOLUME /app/config /app/icons

EXPOSE 8080

ENV CONFIG_FILE_PATH "/app/config/config.yml"
ENV ICON_CACHE_PATH "/app/icons"

ENTRYPOINT ["/app/simplydash"]

HEALTHCHECK --timeout=1s CMD curl -f http://localhost:8080/health || exit 1

LABEL org.opencontainers.image.authors="sandorrobertk94@gmail.com"
