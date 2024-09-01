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
	Format             formatArg `arg:"--format" help:"Output format (json, pretty)" default:"json"`
	Port               uint16    `arg:"--port,env:MINICHAT_PORT" help:"Port to listen on" default:"3001" `
	Host               string    `arg:"--host,env:MINICHAT_HOST" help:"Host to listen on" default:"*"`
	AllowedBucketNames []string  `arg:"--allowed-bucket-names,required,env:MINIOSERVE_ALLOWED_BUCKET_NAMES" help:"Allowed bucket names"`
	// Minio
	MinioEndpoint  string `arg:"--minio-endpoint,required,env:MINIOSERVE_MINIO_ENDPOINT" help:"Minio endpoint"`
	MinioPort      uint16 `arg:"--minio-port,env:MINIOSERVE_MINIO_PORT" help:"Minio port" default:"9000"`
	MinioAccessKey string `arg:"--minio-access-key,required,env:MINIOSERVE_MINIO_ACCESS_KEY" help:"Minio access key"`
	MinioSecretKey string `arg:"--minio-secret-key,required,env:MINIOSERVE_MINIO_SECRET_KEY" help:"Minio secret key"`
	MinioUseSSL    bool   `arg:"--minio-use-ssl,env:MINIOSERVE_MINIO_USE_SSL" help:"Use SSL for Minio" default:"false"`
}

var parsedArgs bool
var globalArgs Args

func parseArgs() Args {
	parsedArgs = true
	arg.MustParse(&globalArgs)
	return globalArgs
}

func GetArgs() Args {
	if !parsedArgs {
		parseArgs()
	}
	return globalArgs
}
