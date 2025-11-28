# ===== STAGE 1: build =====
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy && go build -o go-highload .


# ===== STAGE 2: minimal runtime =====
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/go-highload /app/go-highload

EXPOSE 8080

ENV GIN_MODE=release

CMD ["/app/go-highload"]
