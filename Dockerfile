FROM golang:alpine AS builder
WORKDIR /app
COPY . .
EXPOSE 8090
ENTRYPOINT ["go", "run", "main.go"]