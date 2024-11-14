package middleware

import (
	"context"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/logging"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func UserMiddlewareFactory() func(http.Handler) http.Handler {
	conn := database.GetDatabase()
	logger := logging.GetLogger("http.middleware.user")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			claims, ok := request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)

			if !ok {
				http.Error(writer, "Invalid token", http.StatusUnauthorized)
				return
			}

			var id string

			err := conn.QueryRow(
				request.Context(),
				`SELECT id FROM minichat.users WHERE idp_id = $1`,
				claims.RegisteredClaims.Subject,
			).Scan(&id)

			if err != nil {
				logger.Error().Err(err).Msg("Failed to get user record")
				http.Error(writer, "User not found", http.StatusNotFound)
				return
			}

			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), serverutil.UserIdContextKey{}, id)))
		})
	}

}
