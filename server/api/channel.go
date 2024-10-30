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
		`SELECT
			c.id,
			cp.title,
			cp.description,
			cp.created_at
		FROM minichat.channels AS c
		LEFT JOIN minichat.channels_public AS cp
		ON c.id = cp.id AND c.type = 'public'
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
		err := rows.Scan(&channel.Id, &channel.Title, &description, &createdAt)
		if err != nil {
			return err
		}
		channel.Description = util.ParseSqlString(description)
		channel.CreatedAt = coreutil.FormatTimestampz(createdAt)
		channel.Type = "public"
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

func getChannelsDirect(ctx context.Context, userId string, buffer *[]minichat.ChannelDirect) error {
	conn := database.GetDatabase()
	rows, err := conn.Query(
		ctx,
		`SELECT
			c.id,
			cd.user1_id,
			cd.user2_id
		FROM minichat.channels AS c
		LEFT JOIN minichat.channels_direct AS cd
		ON c.id = cd.id AND c.type = 'direct'
		WHERE cd.user1_id = $1 OR cd.user2_id = $1
		`,
		userId,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var channel minichat.ChannelDirect
		err := rows.Scan(&channel.Id, &channel.User1Id, &channel.User2Id)
		if err != nil {
			return err
		}
		*buffer = append(*buffer, channel)
	}

	return nil
}

func getChannelsDirectHandler(writer http.ResponseWriter, request *http.Request) {
	user := request.Context().Value(minichat.UserProfileContextKey{}).(minichat.UserProfile)
	var directChannels []minichat.ChannelDirect
	ctx := request.Context()

	err := getChannelsDirect(ctx, user.ID, &directChannels)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(writer, util.NewResult(directChannels))
}

func getChannelsPrivate(ctx context.Context, userId string, buffer *[]minichat.ChannelPrivate) error {
	conn := database.GetDatabase()
	rows, err := conn.Query(
		ctx,
		`SELECT
			c.id,
			cp.title,
			cp.created_at
		FROM minichat.channels AS c
		LEFT JOIN minichat.channels_private AS cp
		ON c.id = cp.id AND c.type = 'private'
		LEFT JOIN minichat.channels_private_members AS cpm
		ON c.id = cpm.channel_id
		WHERE cpm.user_id = $1
		`,
		userId,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var channel minichat.ChannelPrivate
		err := rows.Scan(&channel.Id, &channel.Title, &channel.CreatedAt)
		if err != nil {
			return err
		}
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

	err = getChannelsDirect(ctx, user.ID, &channelsResponse.Direct)
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
	router.Get("/direct", getChannelsDirectHandler)
	router.Get("/", getChannelsHandler)
	return "/channels", router
}
