# S3 Upload API (MinIO/AWS)

HTTP API сервер на Go для загрузки файлов в S3-бакет (MinIO или AWS S3) через REST API. Поддерживает загрузку файлов трех типов: TXT, PNG и JSON.

## Возможности

- ✅ HTTP REST API для загрузки файлов через Postman
- ✅ Автоматическое создание S3-бакета при запуске
- ✅ Валидация типов файлов (только .txt, .png, .json)
- ✅ Прямая загрузка файлов в MinIO/AWS S3
- ✅ Поддержка MinIO и AWS S3
- ✅ Docker Compose для быстрого запуска

## Требования

- Docker и Docker Compose (для использования Docker)
- Postman или любой HTTP клиент для тестирования API
- Go 1.23+ (для локальной разработки)

## Быстрый старт

### Запуск с Docker Compose (рекомендуется)

```bash
# Запуск MinIO и HTTP API сервера
docker-compose up --build
```

После запуска:
- **HTTP API**: http://localhost:8080
- **MinIO Web UI**: http://localhost:9001 (логин: `minioadmin`, пароль: `minioadmin`)
- **MinIO API**: http://localhost:9000

### Локальный запуск

1. Убедитесь, что MinIO запущен локально или в Docker
2. Установите зависимости:
```bash
go mod download
```

3. Запустите сервер:
```bash
go run main.go
```

## API Документация

### Health Check

**GET** `/health`

Проверка состояния сервера.

**Ответ:**
```json
{
  "status": "ok",
  "service": "s3-uploader"
}
```

### Загрузка файла

**POST** `/upload`

Загрузка файла в S3-бакет.

**Параметры:**
- `file` (multipart/form-data) - файл для загрузки

**Поддерживаемые типы файлов:**
- `.txt` - текстовые файлы
- `.png` - изображения PNG
- `.json` - JSON файлы

**Успешный ответ (200):**
```json
{
  "success": true,
  "message": "Файл успешно загружен в MinIO",
  "fileName": "example.txt",
  "fileSize": 1024,
  "fileType": "text/plain",
  "objectKey": "example.txt"
}
```

**Ошибка (400/500):**
```json
{
  "success": false,
  "message": "Описание ошибки"
}
```

## Использование с Postman

### Шаг 1: Запустите сервисы

```bash
docker-compose up -d
```

### Шаг 2: Настройте запрос в Postman

1. Метод: **POST**
2. URL: `http://localhost:8080/upload`
3. Body: выберите **form-data**
4. Добавьте ключ `file` типа **File**
5. Выберите файл для загрузки (`.txt`, `.png` или `.json`)

### Пример запроса в Postman

```
POST http://localhost:8080/upload
Content-Type: multipart/form-data

file: [выберите файл]
```

### Шаг 3: Проверка результата

После успешной загрузки файл будет доступен в MinIO:
1. Откройте http://localhost:9001
2. Войдите с учетными данными: `minioadmin` / `minioadmin`
3. Откройте бакет `my-bucket`
4. Найдите загруженный файл

## Примеры загрузки файлов

### Пример 1: Загрузка текстового файла

Создайте файл `test.txt`:
```
Привет, мир!
Это тестовый файл.
```

Загрузите через Postman:
- URL: `POST http://localhost:8080/upload`
- Body → form-data → `file`: выберите `test.txt`

### Пример 2: Загрузка JSON файла

Создайте файл `data.json`:
```json
{
  "name": "Тест",
  "value": 123
}
```

Загрузите через Postman аналогично примеру 1.

### Пример 3: Загрузка PNG изображения

Загрузите любое PNG изображение через Postman.

## Конфигурация

### Переменные окружения

```bash
S3_ENDPOINT=minio:9000          # MinIO endpoint (по умолчанию: minio:9000)
S3_ACCESS_KEY=minioadmin        # Access key (по умолчанию: minioadmin)
S3_SECRET_KEY=minioadmin        # Secret key (по умолчанию: minioadmin)
S3_USE_SSL=false                # Использовать SSL (по умолчанию: false)
```

### Для AWS S3

```bash
export S3_ENDPOINT="s3.amazonaws.com"
export S3_ACCESS_KEY="your-access-key"
export S3_SECRET_KEY="your-secret-key"
export S3_USE_SSL="true"
```

## Структура проекта

```
minn/
├── main.go              # HTTP API сервер
├── Dockerfile           # Docker образ приложения
├── docker-compose.yml   # Docker Compose конфигурация
├── .dockerignore        # Исключения для Docker
├── go.mod               # Зависимости Go
└── README.md            # Документация
```

## Полезные команды

### Запуск сервисов
```bash
docker-compose up -d
```

### Просмотр логов
```bash
docker-compose logs -f s3uploader
```

### Остановка сервисов
```bash
docker-compose down
```

### Остановка и удаление данных
```bash
docker-compose down -v
```

### Пересборка после изменений
```bash
docker-compose build
docker-compose up -d
```

## Ограничения

- Максимальный размер файла: 32 MB
- Поддерживаются только файлы с расширениями: `.txt`, `.png`, `.json`
- Имя файла в S3 совпадает с оригинальным именем файла

## Зависимости

- `github.com/minio/minio-go/v7` - клиент для работы с MinIO и AWS S3

## Лицензия

MIT
# minn
