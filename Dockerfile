# ===== STAGE 1: build =====
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Опционально: go env для более повторяемой сборки
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Сначала модули — чтобы кешировать зависимости
COPY go.mod go.sum ./
RUN go mod download

# Теперь весь код
COPY . .

# Собираем бинарник
RUN go build -o go-highload ./cmd/server

# Если у тебя main.go лежит в корне (как сейчас), а не в cmd/server,
# то вместо строки выше используй:
# RUN go build -o go-highload .

# ===== STAGE 2: minimal runtime =====
FROM alpine:3.20

WORKDIR /app

# Для логов и таймзоны по желанию
RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/go-highload /app/go-highload

EXPOSE 8080

ENV GIN_MODE=release

CMD ["/app/go-highload"]
