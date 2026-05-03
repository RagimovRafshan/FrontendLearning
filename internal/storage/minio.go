package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"messenger/config"
)

// MinIOClient обертка над minio.Client для работы с MinIO S3
type MinIOClient struct {
	client     *minio.Client
	bucketName string
	region     string
}

// NewMinIOClient создает новый клиент MinIO
func NewMinIOClient(cfg config.MinIOConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	minioClient := &MinIOClient{
		client:     client,
		bucketName: cfg.BucketName,
		region:     cfg.Region,
	}

	return minioClient, nil
}

// Client возвращает базовый клиент minio
func (c *MinIOClient) Client() *minio.Client {
	return c.client
}

// BucketName возвращает имя бакета по умолчанию
func (c *MinIOClient) BucketName() string {
	return c.bucketName
}

// HealthCheck проверяет подключение к MinIO
func (c *MinIOClient) HealthCheck(ctx context.Context) error {
	// Проверяем существование бакета или возможность подключения
	found, err := c.client.BucketExists(ctx, c.bucketName)
	if err != nil {
		return fmt.Errorf("MinIO health check failed: %w", err)
	}
	
	if !found {
		// Бакет не существует, пробуем создать
		if err := c.EnsureBucket(ctx, c.bucketName); err != nil {
			return fmt.Errorf("MinIO bucket creation failed: %w", err)
		}
	}
	
	return nil
}

// EnsureBucket создает бакет если он не существует
func (c *MinIOClient) EnsureBucket(ctx context.Context, bucketName string) error {
	exists, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if exists {
		return nil
	}

	err = c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
		Region: c.region,
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	// Устанавливаем политику публичного доступа (опционально)
	// policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]} ]}`, bucketName)
	// err = c.client.SetBucketPolicy(ctx, bucketName, policy)
	// if err != nil {
	// 	return fmt.Errorf("failed to set bucket policy: %w", err)
	// }

	return nil
}

// PutFile загружает файл в MinIO
func (c *MinIOClient) PutFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	info, err := c.client.PutObject(ctx, c.bucketName, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return info.Key, nil
}

// GetFile получает файл из MinIO
func (c *MinIOClient) GetFile(ctx context.Context, objectName string) (*minio.Object, error) {
	obj, err := c.client.GetObject(ctx, c.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return obj, nil
}

// RemoveFile удаляет файл из MinIO
func (c *MinIOClient) RemoveFile(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	return nil
}

// FileExists проверяет существование файла
func (c *MinIOClient) FileExists(ctx context.Context, objectName string) (bool, error) {
	_, err := c.client.StatObject(ctx, c.bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, fmt.Errorf("failed to stat object: %w", err)
	}

	return true, nil
}

// GetPresignedURL генерирует предподписанный URL для доступа к файлу
func (c *MinIOClient) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := c.client.PresignedGetObject(ctx, c.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// GetPresignedPutURL генерирует предподписанный URL для загрузки файла
func (c *MinIOClient) GetPresignedPutURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := c.client.PresignedPutObject(ctx, c.bucketName, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned PUT URL: %w", err)
	}

	return url.String(), nil
}

// ListFiles получает список файлов в бакете
func (c *MinIOClient) ListFiles(ctx context.Context, prefix string, recursive bool) <-chan minio.ObjectInfo {
	return c.client.ListObjects(ctx, c.bucketName, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: recursive,
	})
}
