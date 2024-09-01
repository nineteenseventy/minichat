package core

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nineteenseventy/minichat/core/logging"
)

type MinioConfig struct {
	Endpoint  string
	Port      uint16
	AccessKey string
	SecretKey string
	UseSSL    bool
}

var globalMinio *minio.Client

func InitMinio(ctx context.Context, config MinioConfig) error {
	logger := logging.GetLogger("minio")

	var err error
	endpoint := fmt.Sprintf("%s:%d", config.Endpoint, config.Port)
	globalMinio, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to Minio")
		return err
	}

	logger.Info().Str("host", endpoint).Msg("Connected to Minio")
	return nil
}

func GetMinio() *minio.Client {
	if globalMinio == nil {
		panic("Minio not initialized")
	}
	return globalMinio
}

func GetMinioEnsureBucket(ctx context.Context, bucket string) (*minio.Client, error) {
	minioClient := GetMinio()
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
	}
	return minioClient, nil
}
