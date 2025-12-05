# Многоступенчатая сборка для оптимизации размера образа

# Этап сборки
FROM golang:1.23-alpine AS builder

# Установка рабочих инструментов
RUN apk add --no-cache git

# Установка рабочей директории
WORKDIR /app

# Копирование файлов зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Обновление зависимостей и сборка приложения
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o s3uploader main.go

# Финальный этап
FROM alpine:latest

# Установка CA сертификатов для HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копирование исполняемого файла из этапа сборки
COPY --from=builder /app/s3uploader .

# Переменные окружения по умолчанию
ENV S3_ENDPOINT=minio:9000
ENV S3_ACCESS_KEY=minioadmin
ENV S3_SECRET_KEY=minioadmin
ENV S3_USE_SSL=false

# Команда запуска
CMD ["./s3uploader"]

