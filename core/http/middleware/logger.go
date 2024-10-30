package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/nineteenseventy/minichat/core/logging"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LoggerMiddlewareFactory() func(http.Handler) http.Handler {
	logger := logging.GetLogger("http.middleware.logger")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			t1 := time.Now()

			lrw := loggingResponseWriter{
				ResponseWriter: writer,
				statusCode:     http.StatusOK,
			}

			defer func() {
				scheme := "http"
				if request.TLS != nil {
					scheme = "https"
				}
				logger.Debug().
					Str("from", request.RemoteAddr).
					Int("status", lrw.statusCode).
					Str("size", writer.Header().Get("Content-Length")).
					Str("method", request.Method).
					Str("proto", request.Proto).
					TimeDiff("time", time.Now(), t1).
					Msg(fmt.Sprintf(
						"%s://%s%s",
						scheme,
						request.Host,
						request.RequestURI,
					))
			}()

			next.ServeHTTP(&lrw, request)
		})
	}
}
