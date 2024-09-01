package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/nineteenseventy/minichat/server/api"
	"github.com/nineteenseventy/minichat/server/util"
)

func main() {
	args := ParseArgs()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if args.Format.Format == FormatArgPretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger := util.GetLogger("server")

	r := chi.NewRouter()
	r.Mount("/api", api.ApiRouter())
	r.Mount("/", HealthRouter())

	logger.Info().Uint16("port", args.Port).Str("host", args.Host).Msg("Starting server")

	var host string
	if args.Host == "*" {
		host = ""
	} else {
		host = args.Host
	}

	err := http.ListenAndServe(host+":"+strconv.FormatInt(int64(args.Port), 10), r)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to start server")
	}
}
