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

	"github.com/nineteenseventy/minichat/core"
	"github.com/nineteenseventy/minichat/core/cache"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/logging"
)

func initDatabase(args Args) {
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

func initRedis(args Args) {
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

func initZerolog(args Args) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if args.Format.Format == FormatArgPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func parseHost(args Args) string {
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

	args := GetArgs()

	initZerolog(args)
	logger := logging.GetLogger("server")

	initDatabase(args)
	initRedis(args)
	initMinio(args)

	r := chi.NewRouter()
	r.Mount("/api", ApiRouter())
	r.Mount("/", HealthRouter())

	host := parseHost(args)
	logger.Info().Str("host", host).Msg("Starting server")
	err = http.ListenAndServe(host, r)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
