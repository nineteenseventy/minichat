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
	Format formatArg `arg:"--format" help:"Output format (json, pretty)" default:"json"`
	Port   uint16    `arg:"--port" help:"Port to listen on" default:"3001" env:"MINICHAT_PORT"`
	Host   string    `arg:"--host" help:"Host to listen on" default:"*" env:"MINICHAT_HOST"`
}

func ParseArgs() Args {
	var args Args
	arg.MustParse(&args)
	return args
}
