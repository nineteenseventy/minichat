package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/minichat"
	coreutil "github.com/nineteenseventy/minichat/core/util"
)

func getGroupChannelHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	// user := serverutil.GetUserFromContext(ctx)
	// channelId := chi.URLParam(request, "channelId")

	conn := database.GetDatabase()

	var channel minichat.GroupChannel
	var timestamp pgtype.Timestamptz
	err := conn.QueryRow(
		ctx,
		`
		SELECT
			"channel".id,
			"channel".type,
			"channel".created_at,
			"group".title
		FROM minichat.channels AS "channel"
		LEFT JOIN minichat.channels_group AS "group"
		ON "channel".id = "group".id AND "channel".type = 'group'
		WHERE "channel".id = $1
		`,
	).Scan(channel.Id, &channel.Type, &timestamp, &channel.Title)

	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(writer, "Channel not found", http.StatusNotFound)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	channel.CreatedAt = coreutil.FormatTimestampz(timestamp)

}

func GroupChannelRouter() (string, chi.Router) {
	router := chi.NewRouter()
	router.Get("/{channelId}", getGroupChannelHandler)
	return "/channels/group", router
}
