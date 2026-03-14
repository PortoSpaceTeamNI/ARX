# Base
FROM golang:1.25.5-alpine3.22 AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Development
FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY . .
CMD ["air", "-c", ".air.toml"]