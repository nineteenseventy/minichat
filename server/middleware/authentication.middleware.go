package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/nineteenseventy/minichat/server/util"
)

type AuthenticationMiddlewareOptions struct {
	Domain   string
	Audience []string
}

func errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	http.Error(w, err.Error(), http.StatusUnauthorized)
}

func AuthenticationMiddleware(options AuthenticationMiddlewareOptions) func(http.Handler) http.Handler {
	logger := util.GetLogger("http.middleware.authentication")
	issuerString := fmt.Sprintf("https://%s", options.Domain)
	issuerUrl, err := url.Parse(issuerString)
	if err != nil {
		logger.Fatal().Err(err).Str("issuer", issuerString).Msg("Failed to parse issuer URL")
		panic(err)
	}

	provider := jwks.NewCachingProvider(issuerUrl, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerUrl.String(),
		options.Audience,
		validator.WithAllowedClockSkew(time.Minute),
	)

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create JWT validator")
		panic(err)
	}

	return func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}
}
