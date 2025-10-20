# Dockerfile для fin-core (Go)
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

# Копируем миграции (если они в отдельной папке)
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]