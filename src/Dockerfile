# Use the official Go image as the base image
FROM golang:1.23.0-alpine as builder

# Set the working directory to /app
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -v -o auth-api ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auth-api /app/auth-api
COPY --from=builder /app/pkg/config/config-docker.yml /app/pkg/config/config-docker.yml
COPY --from=builder /app/docs /app/docs

ENV APP_ENV=docker

CMD ["./auth-api"]