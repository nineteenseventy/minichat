package middleware

import (
	"context"
	"database/sql"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/logging"
	"github.com/nineteenseventy/minichat/core/minichat"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func parsePictureUrl(picture sql.NullString) *string {
	logger := logging.GetLogger("http.middleware.user.parsePictureUrl")
	if picture.Valid {
		pictureUrl, err := serverutil.GetCdnUrl("profile", picture.String)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get picture url")
			return nil
		}
		return &pictureUrl
	}
	return nil
}

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

			var id, username string
			var bio, picture, color sql.NullString

			err := conn.QueryRow(
				request.Context(),
				`SELECT
					id,
					username,
					bio,
					picture,
					color
				FROM minichat.users
				WHERE idp_id = $1`,
				claims.RegisteredClaims.Subject,
			).Scan(&id, &username, &bio, &picture, &color)

			if err != nil {
				logger.Error().Err(err).Msg("Failed to get user record")
				http.Error(writer, "User not found", http.StatusNotFound)
				return
			}

			user := minichat.UserProfile{
				ID:       id,
				Username: username,
				Picture:  parsePictureUrl(picture),
				Bio:      util.ParseSqlString(bio),
				Color:    util.ParseSqlString(color),
			}
			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), minichat.UserProfileContextKey{}, user)))
		})
	}

}
