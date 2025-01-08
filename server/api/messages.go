package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/minichat"
	coreutil "github.com/nineteenseventy/minichat/core/util"
	serverutil "github.com/nineteenseventy/minichat/server/util"
)

func parseMessages(rows pgx.Rows, buffer *[]minichat.Message) error {
	var lastMessageId string
	var message minichat.Message
	for rows.Next() {
		var attachmentId, attachmentFilename, attachmentType sql.NullString
		var timestamp pgtype.Timestamptz
		err := rows.Scan(
			&message.Id,
			&message.ChannelId,
			&message.AuthorId,
			&message.Content,
			&timestamp,
			&message.Read,
			&attachmentId,
			&attachmentFilename,
			&attachmentType,
		)
		if err != nil {
			return err
		}
		message.Timestamp = coreutil.FormatTimestampz(timestamp)

		if lastMessageId != message.Id {
			message.Attachments = []minichat.MessageAttachment{}
		}

		if attachmentId.Valid && attachmentFilename.Valid && attachmentType.Valid {
			attachmentUrl, err := serverutil.ParseAttachmentUrl(attachmentId.String, attachmentFilename.String)
			if err == nil {
				var attachment minichat.MessageAttachment
				attachment.Id = attachmentId.String
				attachment.MessageId = message.Id
				attachment.Type = attachmentType.String
				attachment.Url = attachmentUrl
				message.Attachments = append(message.Attachments, attachment)
			}
		}

		if lastMessageId != message.Id {
			*buffer = append(*buffer, message)
		}
	}
	return nil
}

func getMessages(ctx context.Context, channelId string, start uint64, count uint64, buffer *[]minichat.Message) error {
	userId := serverutil.GetUserIdFromContext(ctx)

	conn := database.GetDatabase()

	rows, err := conn.Query(
		ctx,
		`
		SELECT
			"message".id,
			"message".channel_id,
			"message".author_id,
			"message"."content",
			"message"."timestamp",
			(CASE WHEN ("me_member".last_read_message_timestamp >= "message".timestamp) THEN true ELSE false END) AS "is_read",
			"attachment".id,
			"attachment".filename,
			"attachment"."type"
		FROM minichat.messages AS "message"
		-- channel
		LEFT JOIN minichat.channels AS "channel"
		ON "channel".id = "message".channel_id
		-- member me
		LEFT JOIN minichat.channels_members AS "me_member"
		ON "me_member".channel_id = "channel".id
		-- attachments
		LEFT JOIN minichat.attachments AS "attachment"
		ON "attachment".message_id = "message".id
		WHERE "message".channel_id = $1 AND "me_member".user_id = $2
		ORDER BY "message"."timestamp" DESC
		LIMIT $3 OFFSET $4
	`,
		channelId,
		userId,
		count,
		start,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	return parseMessages(rows, buffer)
}

func getMessagesBeforeAfter(ctx context.Context, channelId string, isBefore bool, timestamp time.Time, count uint64, buffer *[]minichat.Message) error {
	userId := serverutil.GetUserIdFromContext(ctx)

	conn := database.GetDatabase()

	var operator string
	if isBefore {
		operator = "<"
	} else {
		operator = ">"
	}

	rows, err := conn.Query(
		ctx,
		fmt.Sprintf(
			`
		SELECT
			"message".id,
			"message".channel_id,
			"message".author_id,
			"message"."content",
			"message"."timestamp",
			(CASE WHEN ("me_member".last_read_message_timestamp >= "message".timestamp) THEN true ELSE false END) AS "is_read",
			"attachment".id,
			"attachment".filename,
			"attachment"."type"
		FROM minichat.messages AS "message"
		-- channel
		LEFT JOIN minichat.channels AS "channel"
		ON "channel".id = "message".channel_id
		-- member me
		LEFT JOIN minichat.channels_members AS "me_member"
		ON "me_member".channel_id = "channel".id
		-- attachments
		LEFT JOIN minichat.attachments AS "attachment"
		ON "attachment".message_id = "message".id
		WHERE "message".channel_id = $1 AND "me_member".user_id = $2 AND "message"."timestamp" %s $3
		ORDER BY "message"."timestamp" DESC
		LIMIT $4
	`, operator),
		channelId,
		userId,
		timestamp,
		count,
	)

	if err != nil {
		return err
	}

	defer rows.Close()

	return parseMessages(rows, buffer)
}

func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	channelId := chi.URLParam(r, "channelId")
	query := r.URL.Query()

	var start, count uint64 = 0, 10
	var err error

	if start_param := query.Get("start"); start_param != "" {
		start, err = strconv.ParseUint(start_param, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if count_param := query.Get("count"); count_param != "" {
		count, err = strconv.ParseUint(count_param, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var timestamp time.Time
	var isBefore bool
	var isAfter bool
	if before_param := query.Get("before"); before_param != "" {
		timestamp, err = coreutil.ParseTimestamp(before_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		isBefore = true
	} else if after_param := query.Get("after"); after_param != "" {
		timestamp, err = coreutil.ParseTimestamp(after_param)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		isAfter = true
	}

	var messages []minichat.Message
	if isBefore || isAfter {
		// if isAfter, then is before is always false
		err = getMessagesBeforeAfter(ctx, channelId, isBefore, timestamp, count, &messages)
	} else {
		err = getMessages(ctx, channelId, start, count, &messages)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	util.JSONResponse(w, util.NewResult(messages))
}

func MessagesRouter(router chi.Router) {
	router.Get("/messages/{channelId}", getMessagesHandler)
}
