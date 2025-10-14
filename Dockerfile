# Этап сборки
FROM golang:1.23-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходники
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=1 GOOS=linux go build -o notifaer ./cmd/notifaer/main.go

# ======================================================
# Финальный образ
# ======================================================
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add \
    ca-certificates \
    sqlite \
    dcron \
    tzdata \
    && rm -rf /var/cache/apk/*

# Устанавливаем временную зону
ENV TZ=Europe/Kyiv

# Создаём рабочего пользователя
RUN addgroup -g 1001 appgroup && \
    adduser -D -s /bin/sh -u 1001 -G appgroup appuser

# Создаём рабочие папки
WORKDIR /app
RUN mkdir -p /app/data /var/log/cron && \
    chown -R appuser:appgroup /app /var/log/cron

# Копируем бинарь и env
COPY --from=builder /build/notifaer /app/
COPY --from=builder /build/.env /app/

# (Опционально) копируем crontab, если используется
COPY --from=builder /build/crontab /tmp/crontab
RUN crontab -u appuser /tmp/crontab && rm /tmp/crontab

# Скрипт запуска
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "Starting notifaer cron service..."' >> /app/start.sh && \
    echo 'crond -f -d 8' >> /app/start.sh && \
    chmod +x /app/start.sh /app/notifaer && \
    chown appuser:appgroup /app/start.sh

# Открываем volume для данных
VOLUME ["/app/data"]

# Запуск
CMD ["/app/start.sh"]
