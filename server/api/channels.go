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
		select
			"channel".id,
			"channel"."type",
			"channel".created_at,
			coalesce("group".title, "direct_partner".username) as "title"
		from minichat.channels AS "channel"
		
		-- member
		left join minichat.channels_members as "me_member"
		on "channel".id = "me_member".channel_id 
		
		-- group
		left join minichat.channels_group as "group"
		on "group".id = "channel".id
		
		-- direct
		left join minichat.channels_direct as "direct"
		on "channel".id = "direct".id 

		-- direct_partner
		left join minichat.channels_members as "direct_partner_member"
		on "direct".id = "direct_partner_member".channel_id and "direct_partner_member".id != "me_member".id
		left join minichat.users as "direct_partner"
		on "direct_partner_member".user_id = "direct_partner".id
		WHERE ("channel"."type" = 'group' OR "channel"."type" = 'direct') and "me_member".user_id = $1
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
