package util

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
	Port   uint16    `arg:"--port,env:MINICHAT_PORT" help:"Port to listen on" default:"3001" `
	Host   string    `arg:"--host,env:MINICHAT_HOST" help:"Host to listen on" default:"*"`
	// Postgres
	PostgresHost     string `arg:"--postgres-host,required,env:MINICHAT_POSTGRES_HOST" help:"Postgres host"`
	PostgresPort     uint16 `arg:"--postgres-port,env:MINICHAT_POSTGRES_PORT" help:"Postgres port" default:"5432"`
	PostgresDatabase string `arg:"--postgres-database,required,env:MINICHAT_POSTGRES_DATABASE" help:"Postgres database"`
	PostgresUser     string `arg:"--postgres-user,required,env:MINICHAT_POSTGRES_USER" help:"Postgres user"`
	PostgresPassword string `arg:"--postgres-password,required,env:MINICHAT_POSTGRES_PASSWORD" help:"Postgres password"`
	PostgresTls      bool   `arg:"--postgres-tls,env:MINICHAT_POSTGRES_TLS" help:"Use TLS for Postgres" default:"false"`
	// Redis
	RedisHost     string `arg:"--redis-host,required,env:MINICHAT_REDIS_HOST" help:"Redis host"`
	RedisPort     uint16 `arg:"--redis-port,env:MINICHAT_REDIS_PORT" help:"Redis port" default:"6379"`
	RedisPassword string `arg:"--redis-password,env:MINICHAT_REDIS_PASSWORD" help:"Redis password"`
	RedisTls      bool   `arg:"--redis-tls,env:MINICHAT_REDIS_TLS" help:"Use TLS for Redis" default:"false"`
	// Minio
	MinioEndpoint  string `arg:"--minio-endpoint,required,env:MINICHAT_MINIO_ENDPOINT" help:"Minio endpoint"`
	MinioPort      uint16 `arg:"--minio-port,env:MINICHAT_MINIO_PORT" help:"Minio port" default:"9000"`
	MinioAccessKey string `arg:"--minio-access-key,required,env:MINICHAT_MINIO_ACCESS_KEY" help:"Minio access key"`
	MinioSecretKey string `arg:"--minio-secret-key,required,env:MINICHAT_MINIO_SECRET_KEY" help:"Minio secret key"`
	MinioUseSSL    bool   `arg:"--minio-use-ssl,env:MINICHAT_MINIO_USE_SSL" help:"Use SSL for Minio" default:"false"`
	// Auth0
	Auth0Domain   string   `arg:"--auth0-domain,required,env:MINICHAT_AUTH0_DOMAIN" help:"Auth0 domain"`
	Auth0Audience []string `arg:"--auth0-audience,required,env:MINICHAT_AUTH0_AUDIENCE" help:"Auth0 audience"`
	// CDN URL
	CdnUrl string `arg:"--cdn-url,required,env:MINICHAT_CDN_URL" help:"CDN URL"`
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
