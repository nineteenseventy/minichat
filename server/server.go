package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nineteenseventy/minichat/core/cache"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/http/middleware"
	"github.com/nineteenseventy/minichat/core/logging"
	"github.com/nineteenseventy/minichat/core/minio"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func initDatabase() {
	args := serverutil.GetArgs()
	databaseConfig := database.DatabaseConfig{
		Host:     args.PostgresHost,
		Port:     args.PostgresPort,
		Database: args.PostgresDatabase,
		User:     args.PostgresUser,
		Password: args.PostgresPassword,
		Tls:      args.PostgresTls,
	}
	err := database.InitDatabase(context.Background(), databaseConfig)
	if err != nil {
		panic(err)
	}
}

func initRedis() {
	args := serverutil.GetArgs()
	redisConfig := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", args.RedisHost, args.RedisPort),
		Password: args.RedisPassword,
		DB:       0,
	}
	err := cache.InitRedis(context.Background(), redisConfig)
	if err != nil {
		panic(err)
	}
}

func initMinio() {
	args := serverutil.GetArgs()
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

func initZerolog() {
	args := serverutil.GetArgs()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if args.Format.Format == serverutil.FormatArgPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func parseHost() string {
	args := serverutil.GetArgs()
	host := args.Host
	if args.Host == "*" {
		host = ""
	}
	return fmt.Sprintf("%s:%d", host, args.Port)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load .env file")
	}

	initZerolog()
	logger := logging.GetLogger("server")

	initDatabase()
	initRedis()
	initMinio()

	r := chi.NewRouter()
	r.Use(middleware.CorsMiddlewareFactory())
	r.Mount("/api", ApiRouter())
	r.Mount("/", HealthRouter())

	host := parseHost()
	logger.Info().Str("host", host).Msg("Starting server")
	err = http.ListenAndServe(host, r)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
