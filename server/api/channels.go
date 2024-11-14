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
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

type ChannelsResponse struct {
	Public  []minichat.ChannelPublic  `json:"public"`
	Private []minichat.ChannelPrivate `json:"private"`
}

func getChannelsPublic(ctx context.Context, buffer *[]minichat.ChannelPublic) error {
	userId := serverutil.GetUserIdFromContext(ctx)
	conn := database.GetDatabase()
	rows, err := conn.Query(
		ctx,
		`
		SELECT
			"channel".id,
			"channel".type,
			"channel".created_at,
			"public".title,
			"public".description,
			COUNT("unread_messages".*) AS "unread_count"
		FROM minichat.channels AS "channel"
		
		-- member me
		LEFT JOIN minichat.channels_members AS "me_member"
		ON "me_member".channel_id = "channel".id
		
		-- unread messages
		LEFT JOIN minichat.messages AS "unread_messages"
		ON "unread_messages".channel_id = "channel".id AND "unread_messages"."timestamp" > "me_member".last_read_message_timestamp
		
		-- public channel
		LEFT JOIN minichat.channels_public AS "public"
		ON "public".id = "channel".id
		WHERE "channel"."type" = 'public' AND "me_member".user_id = $1
		GROUP BY "channel".id, "public".id
		`,
		userId,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var channel minichat.ChannelPublic
		var description sql.NullString
		var createdAt pgtype.Timestamptz
		err := rows.Scan(&channel.Id, &channel.Type, &createdAt, &channel.Title, &description, &channel.UnreadCount)
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
	ctx := request.Context()
	var publicChannels []minichat.ChannelPublic

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
			COALESCE("group".title, "direct_partner".username) AS "title",
			COUNT("unread_messages".*) AS "unread_count"
		FROM minichat.channels AS "channel"
		
		-- member me
		LEFT JOIN minichat.channels_members AS "me_member"
		ON "channel".id = "me_member".channel_id
		
		-- unread messages
		LEFT JOIN minichat.messages AS "unread_messages"
		ON "unread_messages".channel_id = "channel".id AND "unread_messages"."timestamp" > "me_member".last_read_message_timestamp
		
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
		WHERE ("channel"."type" = 'group' OR "channel"."type" = 'direct') AND "me_member".user_id = $1
		GROUP BY "channel".id, "group".id, "direct_partner".id
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
		err := rows.Scan(&channel.Id, &channel.Type, &timestamp, &channel.Title, &channel.UnreadCount)
		if err != nil {
			return err
		}
		channel.CreatedAt = coreutil.FormatTimestampz(timestamp)
		*buffer = append(*buffer, channel)
	}
	return nil
}

func getChannelsPrivateHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userId := serverutil.GetUserIdFromContext(ctx)
	var privateChannels []minichat.ChannelPrivate

	err := getChannelsPrivate(ctx, userId, &privateChannels)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(writer, util.NewResult(privateChannels))
}

func getChannelsHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	userId := serverutil.GetUserIdFromContext(ctx)
	var channelsResponse ChannelsResponse

	err := getChannelsPublic(ctx, &channelsResponse.Public)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = getChannelsPrivate(ctx, userId, &channelsResponse.Private)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(writer, channelsResponse)
}

func ChannelsRouter(router chi.Router) {
	router.Get("/channels/public", getChannelsPublicHandler)
	router.Get("/channels/private", getChannelsPrivateHandler)
	router.Get("/channels", getChannelsHandler)
}
