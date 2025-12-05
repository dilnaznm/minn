package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucketName = "my-bucket"
	region     = "us-east-1"
	serverPort = "8080"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
}

type UploadResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	FileName  string `json:"fileName,omitempty"`
	FileSize  int64  `json:"fileSize,omitempty"`
	FileType  string `json:"fileType,omitempty"`
	ObjectKey string `json:"objectKey,omitempty"`
}

var s3Client *minio.Client

func main() {

	config := loadConfig()

	client, err := createClient(config)
	if err != nil {
		log.Fatalf("Ошибка создания клиента: %v", err)
	}
	s3Client = client

	ctx := context.Background()

	err = createBucket(ctx, client, bucketName, region)
	if err != nil {
		log.Fatalf("Ошибка создания бакета: %v", err)
	}
	fmt.Printf("Бакет '%s' успешно создан или уже существует\n", bucketName)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/health", healthHandler)

	serverAddr := ":" + serverPort
	fmt.Printf("HTTP сервер запущен на порту %s\n", serverPort)
	fmt.Printf("Endpoint для загрузки файлов: http://localhost%s/upload\n", serverAddr)
	fmt.Println("\nПоддерживаемые типы файлов: .txt, .png, .json")

	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "ok",
		"service": "s3-uploader",
	})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Разрешены только POST запросы")
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Ошибка парсинга формы: %v", err))
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Ошибка получения файла: %v", err))
		return
	}
	defer file.Close()

	fileExt := strings.ToLower(filepath.Ext(handler.Filename))
	contentType, err := validateFileType(fileExt)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := context.Background()
	objectKey := handler.Filename

	fileSize := handler.Size
	_, err = s3Client.PutObject(ctx, bucketName, objectKey, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Ошибка загрузки в S3: %v", err))
		return
	}

	response := UploadResponse{
		Success:   true,
		Message:   "Файл успешно загружен в MinIO",
		FileName:  handler.Filename,
		FileSize:  fileSize,
		FileType:  contentType,
		ObjectKey: objectKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	fmt.Printf("Загружен файл: %s (%s, %d bytes)\n", handler.Filename, contentType, fileSize)
}

func validateFileType(ext string) (string, error) {
	allowedTypes := map[string]string{
		".txt":  "text/plain",
		".png":  "image/png",
		".json": "application/json",
	}

	contentType, ok := allowedTypes[ext]
	if !ok {
		return "", fmt.Errorf("неподдерживаемый тип файла. Разрешены только: .txt, .png, .json")
	}

	return contentType, nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(UploadResponse{
		Success: false,
		Message: message,
	})
}

func loadConfig() Config {

	endpoint := os.Getenv("S3_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKey := os.Getenv("S3_ACCESS_KEY")
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("S3_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

	useSSL := os.Getenv("S3_USE_SSL") == "true"

	return Config{
		Endpoint:        endpoint,
		AccessKeyID:     accessKey,
		SecretAccessKey: secretKey,
		UseSSL:          useSSL,
	}
}

func createClient(config Config) (*minio.Client, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
		Region: region,
	})
	if err != nil {
		return nil, fmt.Errorf("не удалось создать клиент: %w", err)
	}

	return client, nil
}

func createBucket(ctx context.Context, client *minio.Client, bucketName, region string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования бакета: %w", err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
		if err != nil {
			return fmt.Errorf("ошибка создания бакета: %w", err)
		}
	}

	return nil
}
