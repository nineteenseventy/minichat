package main

import (
	"fmt"

	"github.com/alexflint/go-arg"
)

type formatArg struct {
	Format string
}

const (
	FormatArgJson   = "json"
	FormatArgPretty = "pretty"
)

func (f *formatArg) UnmarshalText(text []byte) error {
	f.Format = string(text)

	if f.Format != FormatArgJson && f.Format != FormatArgPretty {
		return fmt.Errorf("invalid format: %s", f.Format)
	}
	return nil
}

type Args struct {
	Format           formatArg `arg:"--format" help:"Output format (json, pretty)" default:"json"`
	Port             uint16    `arg:"--port,env:MINICHAT_PORT" help:"Port to listen on" default:"3001" `
	Host             string    `arg:"--host,env:MINICHAT_HOST" help:"Host to listen on" default:"*"`
	PostgresHost     string    `arg:"--postgres-host,required,env:MINICHAT_POSTGRES_HOST" help:"Postgres host"`
	PostgresPort     uint16    `arg:"--postgres-port,env:MINICHAT_POSTGRES_PORT" help:"Postgres port" default:"5432"`
	PostgresDatabase string    `arg:"--postgres-database,required,env:MINICHAT_POSTGRES_DATABASE" help:"Postgres database"`
	PostgresUser     string    `arg:"--postgres-user,required,env:MINICHAT_POSTGRES_USER" help:"Postgres user"`
	PostgresPassword string    `arg:"--postgres-password,required,env:MINICHAT_POSTGRES_PASSWORD" help:"Postgres password"`
	PostgresTls      bool      `arg:"--postgres-tls,env:MINICHAT_POSTGRES_TLS" help:"Use TLS for Postgres" default:"false"`
}

func ParseArgs() Args {
	var args Args
	arg.MustParse(&args)
	return args
}
