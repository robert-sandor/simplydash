# backend build
FROM golang:1.22-alpine as backendBuild

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o ./simplydash ./cmd/simplydash

# frontend build
FROM node:20-alpine as frontendBuild

WORKDIR /app
COPY ./web .

RUN npm install && npm run build

# final image
FROM alpine:3

WORKDIR /app

COPY --from=backendBuild /app/simplydash simplydash
COPY --from=frontendBuild /app/build web/build

VOLUME ["/app/config", "/app/images"]

EXPOSE 8080

ENTRYPOINT ["./simplydash"]
