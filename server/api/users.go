package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nineteenseventy/minichat/core/cache"
	"github.com/nineteenseventy/minichat/core/database"
	httputil "github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/minichat"
	coreutil "github.com/nineteenseventy/minichat/core/util"
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

	if lastEchoTime > coreutil.GetUnixTime()-60 {
		httputil.JSONResponse(writer, minichat.UserStatus{Status: "online", ID: id})
	} else {
		httputil.JSONResponse(writer, minichat.UserStatus{Status: "offline", ID: id})
	}
}

func userSettingsHandler(writer http.ResponseWriter, request *http.Request) {

}

func echoHandler(writer http.ResponseWriter, request *http.Request) {
	redis := cache.GetRedis()
	userId := serverutil.GetUserIdFromContext(request.Context())
	redis.Set(request.Context(), fmt.Sprintf("minichat:onlineStatus:%s", userId), coreutil.GetUnixTime(), 0)
	writer.WriteHeader(http.StatusNoContent)
}

func echoAndGetStatusesHandler(writer http.ResponseWriter, request *http.Request) {
	redis := cache.GetRedis()
	userId := serverutil.GetUserIdFromContext(request.Context())
	redis.Set(request.Context(), fmt.Sprintf("minichat:onlineStatus:%s", userId), coreutil.GetUnixTime(), 0)

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

		if lastEchoTime > coreutil.GetUnixTime()-60 {
			statuses = append(statuses, minichat.UserStatus{ID: id, Status: "online"})
		} else {
			statuses = append(statuses, minichat.UserStatus{ID: id, Status: "offline"})
		}
	}

	httputil.JSONResponse(writer, statuses)
}

func getUserChannelHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	meUserId := serverutil.GetUserIdFromContext(ctx)
	userId := chi.URLParam(request, "id")
	if userId == "me" {
		userId = meUserId
	}

	if meUserId == userId {
		http.Error(writer, "Cannot get channel for self", http.StatusBadRequest)
		return
	}

	conn := database.GetDatabase()
	rows, err := conn.Query(
		ctx,
		`
		SELECT
			"channel".id,
			"channel".type,
			"channel".created_at,
			"direct_partner".username AS "title",
			COUNT("unread_messages".*) as "unread_count"
		FROM minichat.channels AS "channel"
		
		-- member me
			LEFT JOIN minichat.channels_members AS "me_member"
			ON "channel".id = "me_member".channel_id

		-- direct_partner
		LEFT JOIN minichat.channels_members AS "direct_partner_member"
		ON "channel".id = "direct_partner_member".channel_id
		LEFT JOIN minichat.users AS "direct_partner"
		ON "direct_partner_member".user_id = "direct_partner".id

		-- unread messages
		LEFT JOIN minichat.messages AS "unread_messages"
		ON "unread_messages".channel_id = "channel".id AND "unread_messages"."timestamp" > "me_member".last_read_message_timestamp

		WHERE "me_member".user_id = $1 AND "direct_partner".id = $2 AND "channel".type = 'direct'
		GROUP BY "channel".id, "direct_partner".id
		`,
		meUserId,
		userId,
	)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	var channel minichat.Channel
	var createdAt pgtype.Timestamptz
	if !rows.Next() && rows.Err() == nil {
		err := conn.QueryRow(
			ctx,
			`
			INSERT INTO minichat.channels (type)
			VALUES ('direct')
			RETURNING id, type, created_at
			`,
		).Scan(&channel.Id, &channel.Type, &createdAt)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = conn.Exec(
			ctx,
			`
			INSERT INTO minichat.channels_members (channel_id, user_id)
			VALUES ($1, $2), ($1, $3)
			`,
			channel.Id,
			meUserId,
			userId,
		)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = conn.Query(
			ctx,
			`
			INSERT INTO minichat.channels_direct (id)
			VALUES ($1)
			`,
			channel.Id,
		)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		err = conn.QueryRow(
			ctx,
			`
			SELECT username
			FROM minichat.users
			WHERE id = $1
			`,
			userId,
		).Scan(&channel.Title)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		channel.CreatedAt = coreutil.FormatTimestampz(createdAt)
		channel.UnreadCount = 0
	} else {
		err := rows.Scan(&channel.Id, &channel.Type, &createdAt, &channel.Title)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		channel.CreatedAt = coreutil.FormatTimestampz(createdAt)
	}

	httputil.JSONResponse(writer, channel)
}

func UsersRouter(router chi.Router) {
	router.Get("/users", getUsersHandler)
	router.Get("/users/{id}", getUserHandler)
	router.Get("/users/{id}/profile", getUserProfileHandler)
	router.Get("/users/{id}/status", getUserStatusHandler)
	router.Get("/users/{id}/channel", getUserChannelHandler)
	router.Post("users/{id}/settings", userSettingsHandler)
	router.Post("/users/echo", echoHandler)
	router.Post("/users/echoAndGetStatuses", echoAndGetStatusesHandler)

}
