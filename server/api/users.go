package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/core/cache"
	"github.com/nineteenseventy/minichat/core/database"
	httputil "github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/minichat"
	"github.com/nineteenseventy/minichat/core/util"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

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
			Picture:  httputil.ParseSqlString(picture),
		}
		users = append(users, user)
	}
	httputil.JSONResponse(writer, httputil.NewResult(users))
}

func getUserHandler(writer http.ResponseWriter, request *http.Request) {
	conn := database.GetDatabase()
	id := chi.URLParam(request, "id")

	if id == "me" {
		id = serverutil.GetUserIdFromContext(request.Context())
	}

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
		Picture:  serverutil.ParseUserPictureUrl(picture),
	}

	httputil.JSONResponse(writer, user)
}

func getUserProfileHandler(writer http.ResponseWriter, request *http.Request) {
	conn := database.GetDatabase()
	id := chi.URLParam(request, "id")
	if id == "me" {
		id = serverutil.GetUserIdFromContext(request.Context())
	}

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
		Picture:  serverutil.ParseUserPictureUrl(picture),
		Bio:      httputil.ParseSqlString(bio),
		Color:    httputil.ParseSqlString(color),
	}

	httputil.JSONResponse(writer, user)
}

func getUserStatusHandler(writer http.ResponseWriter, request *http.Request) {
	redis := cache.GetRedis()
	id := chi.URLParam(request, "id")
	if id == "me" {
		id = serverutil.GetUserIdFromContext(request.Context())
	}

	result := redis.Get(request.Context(), fmt.Sprintf("minichat:onlineStatus:%s", id))
	if result.Err() != nil {
		httputil.JSONResponse(writer, minichat.UserStatus{Status: "offline", ID: id})
		return
	}

	lastEchoTime, err := result.Int64()
	if err != nil {
		httputil.JSONResponse(writer, minichat.UserStatus{Status: "offline", ID: id})
		return
	}

	if lastEchoTime > util.GetUnixTime()-60 {
		httputil.JSONResponse(writer, minichat.UserStatus{Status: "online", ID: id})
	} else {
		httputil.JSONResponse(writer, minichat.UserStatus{Status: "offline", ID: id})
	}
}

func echoHandler(writer http.ResponseWriter, request *http.Request) {
	redis := cache.GetRedis()
	userId := serverutil.GetUserIdFromContext(request.Context())
	redis.Set(request.Context(), fmt.Sprintf("minichat:onlineStatus:%s", userId), util.GetUnixTime(), 0)
	writer.WriteHeader(http.StatusNoContent)
}

func echoAndGetStatusesHandler(writer http.ResponseWriter, request *http.Request) {
	redis := cache.GetRedis()
	userId := serverutil.GetUserIdFromContext(request.Context())
	redis.Set(request.Context(), fmt.Sprintf("minichat:onlineStatus:%s", userId), util.GetUnixTime(), 0)

	queryIds := request.URL.Query().Get("ids")
	if queryIds == "" {
		http.Error(writer, "ids query parameter is required", http.StatusBadRequest)
		return
	}

	statuses := []minichat.UserStatus{}
	for _, id := range httputil.ParseStringArray(queryIds) {
		result := redis.Get(request.Context(), fmt.Sprintf("minichat:onlineStatus:%s", id))
		if result.Err() != nil {
			statuses = append(statuses, minichat.UserStatus{ID: id, Status: "offline"})
			continue
		}

		lastEchoTime, err := result.Int64()
		if err != nil {
			statuses = append(statuses, minichat.UserStatus{ID: id, Status: "offline"})
			continue
		}

		if lastEchoTime > util.GetUnixTime()-60 {
			statuses = append(statuses, minichat.UserStatus{ID: id, Status: "online"})
		} else {
			statuses = append(statuses, minichat.UserStatus{ID: id, Status: "offline"})
		}
	}

	httputil.JSONResponse(writer, statuses)
}

func UsersRouter(router chi.Router) {
	router.Get("/users", getUsersHandler)
	router.Get("/users/{id}", getUserHandler)
	router.Get("/users/{id}/profile", getUserProfileHandler)
	router.Get("/users/{id}/status", getUserStatusHandler)
	router.Post("/users/echo", echoHandler)
	router.Post("/users/echoAndGetStatuses", echoAndGetStatusesHandler)

}
