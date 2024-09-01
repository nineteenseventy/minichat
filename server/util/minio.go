package util

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	logger := GetLogger("minio")

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
