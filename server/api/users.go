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

func usersHandler(w http.ResponseWriter, r *http.Request) {
	conn := database.GetDatabase()
	rows, err := conn.Query(
		r.Context(),
		`SELECT
					id,
					username,
					picture
				FROM minichat.users`,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []minichat.User
	for rows.Next() {
		var id, username string
		var picture sql.NullString
		err := rows.Scan(&id, &username, &picture)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user := minichat.User{
			ID:       id,
			Username: username,
			Picture:  parsePictureUrl(picture),
		}
		users = append(users, user)
	}
	util.JSONResponse(w, util.NewResult(users))
}

func getMeHandler(w http.ResponseWriter, r *http.Request) {
	userProfile := r.Context().Value(minichat.UserProfileContextKey{}).(minichat.UserProfile)
	user := minichat.User{
		ID:       userProfile.ID,
		Username: userProfile.Username,
	}
	util.JSONResponse(w, user)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	conn := database.GetDatabase()
	id := chi.URLParam(r, "id")

	var username string
	var picture sql.NullString

	err := conn.QueryRow(
		r.Context(),
		`SELECT
					id,
					username,
					picture
				FROM minichat.users
				WHERE id = $1`,
		id,
	).Scan(&id, &username, &picture)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user := minichat.User{
		ID:       id,
		Username: username,
		Picture:  parsePictureUrl(picture),
	}

	util.JSONResponse(w, user)
}

func getMeProfileHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(minichat.UserProfileContextKey{}).(minichat.UserProfile)
	util.JSONResponse(w, user)
}

func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	conn := database.GetDatabase()
	id := chi.URLParam(r, "id")

	var username string
	var bio, picture, color sql.NullString

	err := conn.QueryRow(
		r.Context(),
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
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	user := minichat.UserProfile{
		ID:       id,
		Username: username,
		Picture:  parsePictureUrl(picture),
		Bio:      util.ParseSqlString(bio),
		Color:    util.ParseSqlString(color),
	}

	util.JSONResponse(w, user)
}

func UserRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/users", usersHandler)
	r.Get("/users/me", getMeHandler)
	r.Get("/users/{id}", getUserHandler)
	r.Get("/users/me/profile", getMeProfileHandler)
	r.Get("/users/{id}/profile", getUserProfileHandler)
	return r
}
