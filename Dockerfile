# -------------------------------------------------------------
# ЭТАП 1: Сборка статического бинарного файла
# -------------------------------------------------------------
FROM golang:1.23-alpine AS builder

WORKDIR /build

# Кэшируем зависимости
COPY go.mod ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем статический бинарник под Linux:
# CGO_ENABLED=0 исключает зависимость от системных библиотек C (glibc)
# -ldflags="-w -s" удаляет отладочные символы для минимизации размера
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /build/server main.go

# -------------------------------------------------------------
# ЭТАП 2: Финальный безопасный образ для продакшна
# -------------------------------------------------------------
FROM alpine:3.20

# Устанавливаем сертификаты (для TLS/HTTPS) и curl (для проверок healthcheck)
RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Создаем непривилегированного пользователя appuser (ID 10001)
RUN addgroup -S -g 10001 appgroup && \
    adduser -S -u 10001 -G appgroup -s /sbin/nologin appuser

# Копируем только готовый бинарник из этапа builder
COPY --from=builder --chown=appuser:appgroup /build/server /app/server

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт 8080
EXPOSE 8080

# Встроенная проба здоровья контейнера
HEALTHCHECK --interval=20s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/healthz || exit 1

# Запуск микросервиса
ENTRYPOINT ["/app/server"]

