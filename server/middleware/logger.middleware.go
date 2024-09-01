package middleware

import (
	"net/http"

	"github.com/nineteenseventy/minichat/server/util"
)

func LoggerMiddleware() func(http.Handler) http.Handler {
	logger := util.GetLogger("http.middleware.logger")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Debug().Str("method", r.Method).Str("path", r.URL.Path).Msg("Request received")
			next.ServeHTTP(w, r)
		})
	}
}
