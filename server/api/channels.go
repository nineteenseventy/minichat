package api

import (
	"context"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nineteenseventy/minichat/core/database"
	httputil "github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/minichat"
	coreutil "github.com/nineteenseventy/minichat/core/util"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func getChannelsHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userId := serverutil.GetUserIdFromContext(ctx)
	var channels = []minichat.Channel{}

	conn := database.GetDatabase()
	rows, err := conn.Query(
		ctx,
		`
		SELECT
				"channel".id,
				"channel".type,
				"channel".created_at,
				COALESCE("public".title, "group".title, "direct_partner".username) AS "title",
				"public".description,
				COUNT("unread_messages".*) AS "unread_count"
			FROM minichat.channels AS "channel"
			
			-- member me
			LEFT JOIN minichat.channels_members AS "me_member"
			ON "channel".id = "me_member".channel_id
			
			-- group
			LEFT JOIN minichat.channels_group AS "group"
			ON "group".id = "channel".id
			
			-- direct
			LEFT JOIN minichat.channels_direct AS "direct"
			ON "channel".id = "direct".id
			
			-- direct_partner
			LEFT JOIN minichat.channels_members AS "direct_partner_member"
			ON "direct".id = "direct_partner_member".channel_id AND "direct_partner_member".id != "me_member".id
			LEFT JOIN minichat.users AS "direct_partner"
			ON "direct_partner_member".user_id = "direct_partner".id
			
			-- unread messages
			LEFT JOIN minichat.messages AS "unread_messages"
			ON "unread_messages".channel_id = "channel".id AND "unread_messages"."timestamp" > "me_member".last_read_message_timestamp
			
			-- public channel
			LEFT JOIN minichat.channels_public AS "public"
			ON "public".id = "channel".id
			
			WHERE "me_member".user_id = $1
			GROUP BY "channel".id, "public".id, "group".id, "direct_partner".id, "public".id
		`,
		userId,
	)
	if httputil.HandleError(writer, err) {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var channel minichat.Channel
		var description sql.NullString
		var createdAt pgtype.Timestamptz
		err := rows.Scan(&channel.Id, &channel.Type, &createdAt, &channel.Title, &description, &channel.UnreadCount)
		if httputil.HandleError(writer, err) {
			return
		}
		channel.Description = httputil.ParseSqlString(description)
		channel.CreatedAt = coreutil.FormatTimestampz(createdAt)
		channels = append(channels, channel)
	}

	httputil.JSONResponse(writer, httputil.NewResult(channels))
}

func getChannel(ctx context.Context, channelId string, userId string, buffer *minichat.Channel) error {
	conn := database.GetDatabase()

	var timestamp pgtype.Timestamptz
	var description sql.NullString
	err := conn.QueryRow(
		ctx,
		`
			SELECT
				"channel".id,
				"channel".type,
				"channel".created_at,
				COALESCE("public".title, "group".title, "direct_partner".username) AS "title",
				"public".description
			FROM minichat.channels AS "channel"
			
			-- member me
			LEFT JOIN minichat.channels_members AS "me_member"
			ON "channel".id = "me_member".channel_id
			
			-- group
			LEFT JOIN minichat.channels_group AS "group"
			ON "group".id = "channel".id
			
			-- direct
			LEFT JOIN minichat.channels_direct AS "direct"
			ON "channel".id = "direct".id
			
			-- direct_partner
			LEFT JOIN minichat.channels_members AS "direct_partner_member"
			ON "direct".id = "direct_partner_member".channel_id AND "direct_partner_member".id != "me_member".id
			LEFT JOIN minichat.users AS "direct_partner"
			ON "direct_partner_member".user_id = "direct_partner".id
			
			-- public channel
			LEFT JOIN minichat.channels_public AS "public"
			ON "public".id = "channel".id
			
			WHERE "me_member".user_id = $1 AND "channel".id = $2
		`,
		userId,
		channelId,
	).Scan(&buffer.Id, &buffer.Type, &timestamp, &buffer.Title, &description)
	if err != nil {
		return err
	}
	buffer.CreatedAt = coreutil.FormatTimestampz(timestamp)
	buffer.Description = httputil.ParseSqlString(description)
	return nil
}

func getChannelHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userId := serverutil.GetUserIdFromContext(ctx)
	channelId := chi.URLParam(request, "channelId")
	var channel minichat.Channel

	_, err := serverutil.GetUserMember(ctx, userId, channelId)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	err = getChannel(ctx, channelId, userId, &channel)
	if httputil.HandleError(writer, err) {
		return
	}

	httputil.JSONResponse(writer, channel)
}

func setReadHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userId := serverutil.GetUserIdFromContext(ctx)
	channelId := chi.URLParam(request, "channelId")

	conn := database.GetDatabase()
	_, err := conn.Exec(
		ctx,
		`
		UPDATE minichat.channels_members
		SET last_read_message_timestamp = NOW()
		WHERE user_id = $1 AND channel_id = $2
		`,
		userId,
		channelId,
	)
	if httputil.HandleError(writer, err) {
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

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
		SELECT "member".id, "member".user_id, "user".username, "user".picture
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
		var picture sql.NullString
		err = rows.Scan(&member.ID, &member.UserId, &member.Username, &picture)
		if httputil.HandleError(writer, err) {
			return
		}
		member.Picture = httputil.ParseSqlString(picture)
		members = append(members, member)
	}

	err = rows.Err()
	if httputil.HandleError(writer, err) {
		return
	}

	httputil.JSONResponse(writer, httputil.NewResult(members))
}

func ChannelsRouter(router chi.Router) {
	router.Get("/channels", getChannelsHandler)
	router.Get("/channels/{channelId}", getChannelHandler)
	router.Get("/channels/{channelId}/members", getMembersHandler)
	router.Patch("/channels/{channelId}/setRead", setReadHandler)
}
