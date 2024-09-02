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
	"github.com/minio/minio-go/v7"
	"github.com/nineteenseventy/minichat/core"
	"github.com/nineteenseventy/minichat/core/http/middleware"
	"github.com/nineteenseventy/minichat/core/logging"
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
	minioConfig := core.MinioConfig{
		Endpoint:  args.MinioEndpoint,
		Port:      args.MinioPort,
		AccessKey: args.MinioAccessKey,
		SecretKey: args.MinioSecretKey,
		UseSSL:    args.MinioUseSSL,
	}
	err := core.InitMinio(context.Background(), minioConfig)
	if err != nil {
		panic(err)
	}
}

func parseHost(args Args) string {
	host := args.Host
	if args.Host == "*" {
		host = ""
	}
	return fmt.Sprintf("%s:%d", host, args.Port)
}

func serve(w http.ResponseWriter, request *http.Request) {
	bucket := chi.URLParam(request, "bucket")
	object := chi.URLParam(request, "*")

	args := GetArgs()

	if !ContainsCaseInsensitive(args.AllowedBucketNames, bucket) {
		http.Error(w, "Bucket not allowed", http.StatusForbidden)
		return
	}

	if object == "" {
		http.Error(w, "Object not specified", http.StatusBadRequest)
		return
	}

	minioClient := core.GetMinio()

	bucketExists, err := minioClient.BucketExists(request.Context(), bucket)

	if err != nil {
		http.Error(w, "Object not found", http.StatusInternalServerError)
		return
	}

	if !bucketExists {
		http.Error(w, "Object not found", http.StatusNotFound)
		return
	}

	objectInfo, err := minioClient.StatObject(request.Context(), bucket, object, minio.StatObjectOptions{})
	if err != nil {
		http.Error(w, "Object not found", http.StatusNotFound)
		return
	}

	objectReader, err := minioClient.GetObject(request.Context(), bucket, object, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, "Error getting object", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", objectInfo.ContentType)
	w.Header().Set("Last-Modified", objectInfo.LastModified.Format(http.TimeFormat))

	written, err := io.Copy(w, objectReader)
	if err != nil {
		http.Error(w, "Error reading object", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", written))
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

	r := chi.NewRouter()
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CorsMiddleware())
	r.Get("/{bucket}/*", serve)

	host := parseHost(args)
	logger.Info().Str("host", host).Msg("Starting server")
	err = http.ListenAndServe(host, r)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
