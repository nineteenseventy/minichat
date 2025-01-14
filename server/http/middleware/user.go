package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	auth0auth "github.com/auth0/go-auth0/authentication"
	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/logging"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func getUserInfo(token string) (*auth0auth.UserInfoResponse, error) {
	args := serverutil.GetArgs()

	url, err := url.JoinPath(args.Auth0Domain, "userinfo")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo auth0auth.UserInfoResponse

	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}

/*
createNewUser creates a new user record in the database and returns the user's ID.
*/
func createNewUser(user auth0auth.UserInfoResponse) (string, error) {
	conn := database.GetDatabase()

	var id string

	username := user.Nickname

	err := conn.QueryRow(
		context.Background(),
		`INSERT INTO minichat.users (idp_id, username) VALUES ($1, $2) RETURNING id`,
		user.Sub,
		username,
	).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
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

			var id string
			idpId := claims.RegisteredClaims.Subject

			err := conn.QueryRow(
				request.Context(),
				`SELECT id FROM minichat.users WHERE idp_id = $1`,
				idpId,
			).Scan(&id)

			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr == pgx.ErrNoRows {
					token, err := jwtmiddleware.AuthHeaderTokenExtractor(request)
					if err != nil {
						http.Error(writer, err.Error(), http.StatusInternalServerError)
						return
					}
					user, err := getUserInfo(token)
					if err != nil {
						logger.Err(err).Str("idpId", idpId).Msg("Failed to get user info")
						http.Error(writer, err.Error(), http.StatusInternalServerError)
						return
					}
					id, err = createNewUser(*user)
					if err != nil {
						logger.Err(err).Str("idpId", idpId).Msg("Failed to create new user")
						http.Error(writer, err.Error(), http.StatusInternalServerError)
						return
					}
					logger.Info().Str("id", id).Str("idpId", idpId).Msg("Created new user")
				} else {
					http.Error(writer, err.Error(), http.StatusInternalServerError)
					return
				}
			}

			next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), serverutil.UserIdContextKey{}, id)))
		})
	}

}
