package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/nineteenseventy/minichat/core/cache"
	"github.com/nineteenseventy/minichat/core/database"
	httputil "github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/minichat"
	coreminio "github.com/nineteenseventy/minichat/core/minio"
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
	if httputil.HandleError(writer, err) {
		return
	}
	defer rows.Close()

	var users []minichat.User
	for rows.Next() {
		var id, username string
		var picture sql.NullString
		err := rows.Scan(&id, &username, &picture)
		if httputil.HandleError(writer, err) {
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

	if httputil.HandleError(writer, err) {
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
	var bio, picture sql.NullString

	err := conn.QueryRow(
		request.Context(),
		`SELECT
			id,
			username,
			bio,
			picture
		FROM minichat.users
		WHERE id = $1`,
		id,
	).Scan(&id, &username, &bio, &picture)

	if httputil.HandleError(writer, err) {
		return
	}

	user := minichat.UserProfile{
		ID:       id,
		Username: username,
		Picture:  serverutil.ParseUserPictureUrl(picture),
		Bio:      httputil.ParseSqlString(bio),
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

func putProfileHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userId := serverutil.GetUserIdFromContext(ctx)
	id := chi.URLParam(request, "id")
	if id == "me" {
		id = serverutil.GetUserIdFromContext(request.Context())
	}

	if id != userId {
		http.Error(writer, "You cannot change another users settings.", http.StatusForbidden)
		return
	}

	var currentProfile minichat.PatchUserProfile
	err := httputil.JSONRequest(request, &currentProfile)
	if httputil.HandleError(writer, err) {
		return
	}

	var newProfile minichat.UserProfile
	var bio, picture sql.NullString

	conn := database.GetDatabase()
	err = conn.QueryRow(
		ctx,
		`
		UPDATE minichat.users
		SET
			username = $2,
			bio = $3
		WHERE id = $1
		RETURNING id, username, bio, picture
		`,
		userId,
		currentProfile.Username,
		currentProfile.Bio,
	).Scan(&newProfile.ID, &newProfile.Username, &bio, &picture)

	if httputil.HandleError(writer, err) {
		return
	}

	newProfile.Bio = httputil.ParseSqlString(bio)
	newProfile.Picture = serverutil.ParseUserPictureUrl(picture)

	httputil.JSONResponse(writer, newProfile)
}

func postUserPictureHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userId := serverutil.GetUserIdFromContext(ctx)
	id := chi.URLParam(request, "id")
	if id == "me" {
		id = serverutil.GetUserIdFromContext(request.Context())
	}

	if id != userId {
		http.Error(writer, "You cannot change another users settings.", http.StatusForbidden)
		return
	}

	minioClient := coreminio.GetMinio()

	newPictureKey := fmt.Sprintf("%s/%s/%d", userId, coreutil.NewUuid(), coreutil.GetUnixTime())

	ContentType := request.Header.Get("Content-Type")
	if ContentType == "" {
		http.Error(writer, "Content-Type header is required", http.StatusBadRequest)
		return
	}

	size := request.ContentLength
	var _10MB int64 = 1024 * 1024 * 10
	if size > _10MB {
		http.Error(writer, "File size is too large", http.StatusBadRequest)
		return
	}

	if size <= -1 {
		http.Error(writer, "File size is unknown, please provide Content-Length header", http.StatusBadRequest)
		return
	}

	if size == 0 {
		http.Error(writer, "File size is too small", http.StatusBadRequest)
		return
	}

	minioInfo, err := minioClient.PutObject(ctx, serverutil.ProfilePictureBucket, newPictureKey, request.Body, size, minio.PutObjectOptions{
		ContentType: ContentType,
	})
	if httputil.HandleError(writer, err) {
		return
	}

	conn := database.GetDatabase()
	var user minichat.User
	var picture sql.NullString
	err = conn.QueryRow(
		ctx,
		`
		UPDATE minichat.users
		SET picture = $2
		WHERE id = $1
		RETURNING id, username, picture
		`,
		userId,
		minioInfo.Key,
	).Scan(&user.ID, &user.Username, &picture)
	if httputil.HandleError(writer, err) {
		return
	}

	user.Picture = serverutil.ParseUserPictureUrl(picture)

	httputil.JSONResponse(writer, user)
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
	if httputil.HandleError(writer, err) {
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
		if httputil.HandleError(writer, err) {
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
		if httputil.HandleError(writer, err) {
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
		if httputil.HandleError(writer, err) {
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
		if httputil.HandleError(writer, err) {
			return
		}

		channel.CreatedAt = coreutil.FormatTimestampz(createdAt)
		channel.UnreadCount = 0
	} else {
		err := rows.Scan(&channel.Id, &channel.Type, &createdAt, &channel.Title)
		if httputil.HandleError(writer, err) {
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
	router.Put("/users/{id}/profile", putProfileHandler)
	router.Post("/users/{id}/picture", postUserPictureHandler)
	router.Post("/users/echo", echoHandler)
	router.Post("/users/echoAndGetStatuses", echoAndGetStatusesHandler)

}
