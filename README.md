# S3 Bucket Manager (MinIO/AWS)

Go-программа для создания S3-бакета и загрузки файлов трех типов: TXT, PNG и JSON.

## Возможности

- ✅ Создание S3-бакета (MinIO или AWS S3)
- ✅ Автоматическая генерация PNG файла (если отсутствует)
- ✅ Загрузка файлов трех типов:
  - TXT (текстовые файлы)
  - PNG (изображения)
  - JSON (структурированные данные)

## Требования

- Go 1.16 или выше
- MinIO сервер (для локального использования) или AWS S3 аккаунт

## Установка

1. Клонируйте репозиторий или перейдите в директорию проекта:
```bash
cd minn
```

2. Установите зависимости:
```bash
go mod download
```

## Настройка

### Для MinIO (локальный S3-сервер)

1. Установите и запустите MinIO:
```bash
# Установка MinIO (macOS)
brew install minio/stable/minio

# Запуск MinIO
minio server ~/minio-data
```

MinIO будет доступен по адресу `http://localhost:9000` (веб-интерфейс) и `localhost:9000` (API).

По умолчанию используются учетные данные:
- Access Key: `minioadmin`
- Secret Key: `minioadmin`

2. Запустите программу:
```bash
go run main.go
```

### Для AWS S3

Установите переменные окружения:

```bash
export S3_ENDPOINT="s3.amazonaws.com"
export S3_ACCESS_KEY="your-access-key"
export S3_SECRET_KEY="your-secret-key"
export S3_USE_SSL="true"
```

Затем запустите:
```bash
go run main.go
```

### Кастомная конфигурация

Вы можете настроить подключение через переменные окружения:

```bash
export S3_ENDPOINT="your-endpoint"          # По умолчанию: localhost:9000
export S3_ACCESS_KEY="your-access-key"      # По умолчанию: minioadmin
export S3_SECRET_KEY="your-secret-key"      # По умолчанию: minioadmin
export S3_USE_SSL="true"                    # По умолчанию: false
```

## Использование

1. Убедитесь, что у вас есть файлы для загрузки:
   - `sample.txt` - текстовый файл (уже включен)
   - `sample.json` - JSON файл (уже включен)
   - `sample.png` - PNG файл (будет автоматически создан, если отсутствует)

2. Запустите программу:
```bash
go run main.go
```

Программа выполнит следующие действия:
1. Создаст подключение к S3 (MinIO или AWS)
2. Создаст бакет `my-bucket` (если его еще нет)
3. Загрузит все три типа файлов в бакет

## Сборка

Создать исполняемый файл:

```bash
go build -o s3uploader main.go
```

Запустить:
```bash
./s3uploader
```

## Структура проекта

```
minn/
├── main.go          # Основная программа
├── sample.txt       # Пример текстового файла
├── sample.json      # Пример JSON файла
├── sample.png       # PNG файл (создается автоматически)
├── go.mod           # Зависимости Go
└── README.md        # Документация
```

## Примеры использования

### Минимальный запуск (MinIO по умолчанию)
```bash
go run main.go
```

### С кастомными параметрами MinIO
```bash
export S3_ENDPOINT="192.168.1.100:9000"
export S3_ACCESS_KEY="myaccesskey"
export S3_SECRET_KEY="mysecretkey"
go run main.go
```

### С AWS S3
```bash
export S3_ENDPOINT="s3.amazonaws.com"
export S3_ACCESS_KEY="AKIAIOSFODNN7EXAMPLE"
export S3_SECRET_KEY="wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
export S3_USE_SSL="true"
go run main.go
```

## Зависимости

- `github.com/minio/minio-go/v7` - клиент для работы с MinIO и AWS S3

## Лицензия

MIT

