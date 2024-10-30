package api

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/minichat"
	coreutil "github.com/nineteenseventy/minichat/core/util"
)

type ChannelsResponse struct {
	Public  []minichat.ChannelPublic  `json:"public"`
	Private []minichat.ChannelPrivate `json:"private"`
}

func getChannelsPublic(ctx context.Context, buffer *[]minichat.ChannelPublic) error {
	conn := database.GetDatabase()
	rows, err := conn.Query(
		ctx,
		`
		SELECT
			"channel".id,
			"channel".type,
			"channel".created_at,
			"public".title,
			"public".description
		FROM minichat.channels AS "channel"
		LEFT JOIN minichat.channels_public AS "public"
		ON "channel".id = "public".id AND "channel".type = 'public'
		WHERE "channel".type = 'public'
		`,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var channel minichat.ChannelPublic
		var description sql.NullString
		var createdAt pgtype.Timestamptz
		err := rows.Scan(&channel.Id, &channel.Type, &createdAt, &channel.Title, &description)
		if err != nil {
			return err
		}
		channel.Description = util.ParseSqlString(description)
		channel.CreatedAt = coreutil.FormatTimestampz(createdAt)
		*buffer = append(*buffer, channel)
	}
	return nil
}

func getChannelsPublicHandler(writer http.ResponseWriter, request *http.Request) {
	var publicChannels []minichat.ChannelPublic
	ctx := request.Context()

	err := getChannelsPublic(ctx, &publicChannels)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(writer, util.NewResult(publicChannels))
}

func getChannelsPrivate(ctx context.Context, userId string, buffer *[]minichat.ChannelPrivate) error {
	conn := database.GetDatabase()
	rows, err := conn.Query(
		ctx,
		`
		SELECT
			"channel".id,
			"channel"."type",
			"channel".created_at,
			COALESCE("direct_user".username, "group".title) AS "title"
		FROM minichat.channels AS "channel"

		-- group channels
		LEFT JOIN minichat.channels_group AS "group"
		ON "channel".id = "group".id AND "channel"."type" = 'group'
		LEFT JOIN minichat.channels_group_members AS "group_member"
		ON "channel".id = "group_member".channel_id AND "group_member".user_id = $1
		
		-- direct channels
		LEFT JOIN minichat.channels_direct AS "direct"
		ON "channel".id = "direct".id AND "channel"."type" = 'direct' AND (
			"direct".user1_id = $1 OR "direct".user2_id = $1
		)
		LEFT JOIN minichat.users AS "direct_user"
		ON ("direct_user".id = "direct".user1_id AND "direct".user1_id != $1)
		OR ("direct_user".id = "direct".user2_id AND "direct".user2_id != $1)
		WHERE "channel"."type" = 'group' OR "channel"."type" = 'direct'
		`,
		userId,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var channel minichat.ChannelPrivate
		var timestamp pgtype.Timestamptz
		err := rows.Scan(&channel.Id, &channel.Type, &timestamp, &channel.Title)
		if err != nil {
			return err
		}
		channel.CreatedAt = coreutil.FormatTimestampz(timestamp)
		*buffer = append(*buffer, channel)
	}
	return nil
}

func getChannelsPrivateHandler(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(minichat.UserProfileContextKey{}).(minichat.UserProfile)
	var privateChannels []minichat.ChannelPrivate
	ctx := request.Context()

	err := getChannelsPrivate(ctx, user.ID, &privateChannels)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(writer, util.NewResult(privateChannels))
}

func getChannelsHandler(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(minichat.UserProfileContextKey{}).(minichat.UserProfile)
	var channelsResponse ChannelsResponse
	ctx := request.Context()

	err := getChannelsPublic(ctx, &channelsResponse.Public)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = getChannelsPrivate(ctx, user.ID, &channelsResponse.Private)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(writer, channelsResponse)
}

func ChannelRouter() (string, chi.Router) {
	router := chi.NewRouter()
	router.Get("/public", getChannelsPublicHandler)
	router.Get("/private", getChannelsPrivateHandler)
	router.Get("/", getChannelsHandler)
	return "/channels", router
}
