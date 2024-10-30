package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nineteenseventy/minichat/core/cache"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/http/middleware"
	httputil "github.com/nineteenseventy/minichat/core/http/util"
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

func debugRoute(prefix string, route chi.Routes) {
	for _, method := range route.Routes() {
		patternIsClean := !strings.HasSuffix(method.Pattern, "/*")
		var cleanPattern string
		if patternIsClean {
			cleanPattern = method.Pattern
		} else {
			cleanPattern = strings.TrimSuffix(method.Pattern, "/*")
		}
		path := fmt.Sprintf("%s%s", prefix, cleanPattern)

		if patternIsClean {
			log.Debug().Str("path", path).Msg("Registered route")
		}
		if method.SubRoutes != nil {
			debugRoute(path, method.SubRoutes)
		}
	}
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

	router := chi.NewRouter()
	router.Use(middleware.CorsMiddlewareFactory())
	router.Mount("/api", ApiRouter())
	router.Mount("/", HealthRouter())

	// debug route
	debugRoute("", router)

	args := serverutil.GetArgs()

	host := httputil.ParseHost(args.Host, args.Port)
	logger.Info().Str("host", host).Msg("Starting server")
	err = http.ListenAndServe(host, router)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
