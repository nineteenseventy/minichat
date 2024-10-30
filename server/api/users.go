package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/logging"
	"github.com/nineteenseventy/minichat/core/minichat"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func parsePictureUrl(picture sql.NullString) *string {
	logger := logging.GetLogger("server.api.users.parsePictureUrl")
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

func getUsersHandler(writer http.ResponseWriter, request *http.Request) {
	conn := database.GetDatabase()
	rows, err := conn.Query(
		request.Context(),
		`SELECT
			id,
			username,
			picture
		FROM minichat.users`,
	)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []minichat.User
	for rows.Next() {
		var id, username string
		var picture sql.NullString
		err := rows.Scan(&id, &username, &picture)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		user := minichat.User{
			ID:       id,
			Username: username,
			Picture:  parsePictureUrl(picture),
		}
		users = append(users, user)
	}
	util.JSONResponse(writer, util.NewResult(users))
}

func getMeHandler(writer http.ResponseWriter, request *http.Request) {
	userProfile := request.Context().Value(minichat.UserProfileContextKey{}).(minichat.UserProfile)
	user := minichat.User{
		ID:       userProfile.ID,
		Username: userProfile.Username,
	}
	util.JSONResponse(writer, user)
}

func getUserHandler(writer http.ResponseWriter, request *http.Request) {
	conn := database.GetDatabase()
	id := chi.URLParam(request, "id")

	var username string
	var picture sql.NullString

	err := conn.QueryRow(
		request.Context(),
		`SELECT
			id,
			username,
			picture
		FROM minichat.users
		WHERE id = $1`,
		id,
	).Scan(&id, &username, &picture)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	user := minichat.User{
		ID:       id,
		Username: username,
		Picture:  parsePictureUrl(picture),
	}

	util.JSONResponse(writer, user)
}

func getMeProfileHandler(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(minichat.UserProfileContextKey{}).(minichat.UserProfile)
	util.JSONResponse(writer, user)
}

func getUserProfileHandler(writer http.ResponseWriter, request *http.Request) {
	conn := database.GetDatabase()
	id := chi.URLParam(request, "id")

	var username string
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
		WHERE id = $1`,
		id,
	).Scan(&id, &username, &bio, &picture, &color)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	user := minichat.UserProfile{
		ID:       id,
		Username: username,
		Picture:  parsePictureUrl(picture),
		Bio:      util.ParseSqlString(bio),
		Color:    util.ParseSqlString(color),
	}

	util.JSONResponse(writer, user)
}

func UserRouter() chi.Router {
	router := chi.NewRouter()
	router.Get("/users", getUsersHandler)
	router.Get("/users/me", getMeHandler)
	router.Get("/users/{id}", getUserHandler)
	router.Get("/users/me/profile", getMeProfileHandler)
	router.Get("/users/{id}/profile", getUserProfileHandler)
	return router
}
