# assembly stage
FROM golang:1.23-alpine AS builder

# Install the necessary packages for assembly
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source files
COPY . .

# Compiling the binary file
RUN CGO_ENABLED=1 GOOS=linux go build -o notifaer ./cmd/notifaer/main.go

# ======================================================
# Final image
# ======================================================
FROM alpine:latest

# Install the necessary packages
RUN apk --no-cache add \
    ca-certificates \
    sqlite \
    dcron \
    tzdata \
    && rm -rf /var/cache/apk/*

# Install timezone
ENV TZ=Europe/Kyiv

# Create a working user
RUN addgroup -g 1001 appgroup && \
    adduser -D -s /bin/sh -u 1001 -G appgroup appuser

# Create working directories
WORKDIR /app
RUN mkdir -p /app/data /var/log/cron && \
    chown -R appuser:appgroup /app /var/log/cron

# Copy binary file and env file
COPY --from=builder /build/notifaer /app/
COPY --from=builder /build/.env /app/

# (optional) copy crontab
COPY --from=builder /build/crontab /tmp/crontab
RUN crontab -u appuser /tmp/crontab && rm /tmp/crontab

# Script to start
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'echo "Starting notifaer cron service..."' >> /app/start.sh && \
    echo 'crond -f -d 8' >> /app/start.sh && \
    chmod +x /app/start.sh /app/notifaer && \
    chown appuser:appgroup /app/start.sh

# Open volume for data
VOLUME ["/app/data"]

# Start command
CMD ["/app/start.sh"]
