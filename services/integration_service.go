package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/igralkin/go-highload/models"
)

type IntegrationService struct {
	client     *minio.Client
	bucketName string
}

func NewIntegrationService(endpoint, accessKey, secretKey, bucketName string, useSSL bool) (*IntegrationService, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	svc := &IntegrationService{
		client:     minioClient,
		bucketName: bucketName,
	}

	// Ждём готовности MinIO с ретраями, чтобы не падать при старте
	const (
		maxAttempts = 10
		delay       = 3 * time.Second
	)

	ctx := context.Background()
	var lastErr error

	for i := 1; i <= maxAttempts; i++ {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil {
			if !exists {
				if err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
					return nil, fmt.Errorf("failed to create bucket: %w", err)
				}
				log.Printf("Created bucket %q in MinIO\n", bucketName)
			}
			log.Printf("MinIO bucket %q is ready\n", bucketName)
			return svc, nil
		}

		lastErr = errBucketExists
		log.Printf("MinIO not ready yet (attempt %d/%d): %v", i, maxAttempts, errBucketExists)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("failed to init MinIO after %d attempts: %w", maxAttempts, lastErr)
}

// SaveUsers сохраняет текущий список пользователей в MinIO в формате JSON.
func (s *IntegrationService) SaveUsers(ctx context.Context, users []models.User) (string, error) {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal users: %w", err)
	}

	objectName := fmt.Sprintf("users-%d.json", time.Now().Unix())

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{
			ContentType: "application/json",
		})
	if err != nil {
		return "", fmt.Errorf("failed to put object: %w", err)
	}

	return objectName, nil
}

// ListObjects возвращает список объектов в бакете.
func (s *IntegrationService) ListObjects(ctx context.Context) ([]string, error) {
	var result []string
	ch := s.client.ListObjects(ctx, s.bucketName, minio.ListObjectsOptions{
		Recursive: true,
	})

	for obj := range ch {
		if obj.Err != nil {
			return nil, obj.Err
		}
		result = append(result, obj.Key)
	}

	return result, nil
}
