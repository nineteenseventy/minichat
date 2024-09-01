package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nineteenseventy/minichat/server/api"
	"github.com/nineteenseventy/minichat/server/util"
	"github.com/nineteenseventy/minichat/server/util/database"
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
	err := util.InitDatabase(context.Background(), databaseConfig)
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
	godotenv.Load()
	args := ParseArgs()

	initZerolog(args)
	logger := util.GetLogger("server")

	initDatabase(args)

	r := chi.NewRouter()
	r.Mount("/api", api.UserRouter())
	r.Mount("/", HealthRouter())

	logger.Info().Uint16("port", args.Port).Str("host", args.Host).Msg("Starting server")
	err := http.ListenAndServe(parseHost(args), r)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
