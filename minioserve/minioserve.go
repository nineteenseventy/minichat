package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	miniolib "github.com/minio/minio-go/v7"
	"github.com/nineteenseventy/minichat/core/http/middleware"
	httputil "github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/logging"
	"github.com/nineteenseventy/minichat/core/minio"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ContainsCaseInsensitive(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

func initZerolog(args Args) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if args.Format.Format == FormatArgPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func initMinio(args Args) {
	minioConfig := minio.MinioConfig{
		Endpoint:  args.MinioEndpoint,
		Port:      args.MinioPort,
		AccessKey: args.MinioAccessKey,
		SecretKey: args.MinioSecretKey,
		UseSSL:    args.MinioUseSSL,
	}
	err := minio.InitMinio(context.Background(), minioConfig)
	if err != nil {
		panic(err)
	}
}

func serve(writer http.ResponseWriter, request *http.Request) {
	bucket := chi.URLParam(request, "bucket")
	object := chi.URLParam(request, "*")

	args := GetArgs()

	if !ContainsCaseInsensitive(args.AllowedBucketNames, bucket) {
		http.Error(writer, "Bucket is not allowed or does not exist", http.StatusNotFound)
		return
	}

	if object == "" {
		http.Error(writer, "Object not specified", http.StatusBadRequest)
		return
	}

	minioClient := minio.GetMinio()

	bucketExists, err := minioClient.BucketExists(request.Context(), bucket)

	if err != nil {
		http.Error(writer, "Error checking bucket", http.StatusInternalServerError)
		return
	}

	if !bucketExists {
		http.Error(writer, "Bucket does not exist", http.StatusNotFound)
		return
	}

	objectInfo, err := minioClient.StatObject(request.Context(), bucket, object, miniolib.StatObjectOptions{})
	if err != nil {
		http.Error(writer, "Object not found", http.StatusNotFound)
		return
	}

	objectReader, err := minioClient.GetObject(request.Context(), bucket, object, miniolib.GetObjectOptions{})
	if err != nil {
		http.Error(writer, "Error getting object", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", objectInfo.ContentType)
	writer.Header().Set("Last-Modified", objectInfo.LastModified.Format(http.TimeFormat))

	// Handle HEAD requests
	if request.Method == http.MethodHead {
		writer.Header().Set("Content-Length", fmt.Sprintf("%d", objectInfo.Size))
		return
	}

	written, err := io.Copy(writer, objectReader)
	if err != nil {
		http.Error(writer, "Error reading object", http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Length", fmt.Sprintf("%d", written))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load .env file")
	}

	args := GetArgs()

	initZerolog(args)
	logger := logging.GetLogger("server")

	initMinio(args)

	router := chi.NewRouter()
	router.Use(middleware.LoggerMiddlewareFactory())
	router.Use(middleware.CorsMiddlewareFactory())
	router.Mount("/", HealthRouter())
	router.Get("/{bucket}/*", serve)
	router.Head("/{bucket}/*", serve)

	host := httputil.ParseHost(args.Host, args.Port)
	logger.Info().Str("host", host).Msg("Starting server")
	err = http.ListenAndServe(host, router)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
