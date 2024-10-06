FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/bin/auth cmd/main.go

FROM alpine:latest
WORKDIR /root/
CMD ["./auth"]
COPY --from=builder /app/bin/auth .