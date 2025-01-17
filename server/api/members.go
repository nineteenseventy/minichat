package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/nineteenseventy/minichat/core/database"
	httputil "github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/minichat"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func getMembersHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	channelId := chi.URLParam(request, "channelId")
	userId := serverutil.GetUserIdFromContext(ctx)

	_, err := serverutil.GetUserMember(ctx, userId, channelId)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	query := request.URL.Query()
	var username string = query.Get("username")

	var start, count uint64 = 0, 10
	if start_param := query.Get("start"); start_param != "" {
		start, err = strconv.ParseUint(start_param, 10, 64)
		if httputil.HandleError(writer, err) {
			return
		}
	}

	if count_param := query.Get("count"); count_param != "" {
		count, err = strconv.ParseUint(count_param, 10, 64)
		if httputil.HandleError(writer, err) {
			return
		}
	}

	conn := database.GetDatabase()
	rows, err := conn.Query(
		ctx,
		`
		SELECT "member".id, "member".user_id, "user".username
		FROM minichat.channels_members AS "member"
		JOIN minichat.users AS "user"
		ON "member".user_id = "user".id
		WHERE "member".channel_id = $1
		AND LOWER("user".username) LIKE LOWER('%' || $2 || '%')
		LIMIT $3 OFFSET $4
		`,
		channelId,
		username,
		count,
		start,
	)
	if httputil.HandleError(writer, err) {
		return
	}
	defer rows.Close()

	var members []minichat.Member = []minichat.Member{}
	for rows.Next() {
		var member minichat.Member
		err = rows.Scan(&member.ID, &member.UserId, &member.Username)
		if httputil.HandleError(writer, err) {
			return
		}
		members = append(members, member)
	}

	err = rows.Err()
	if httputil.HandleError(writer, err) {
		return
	}

	httputil.JSONResponse(writer, httputil.NewResult(members))
}

func MembersRouter(router chi.Router) {
	router.Get("/members/{channelId}", getMembersHandler)
}
