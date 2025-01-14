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
	httputil "github.com/nineteenseventy/minichat/core/http/util"
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
				attachment.Filename = attachmentFilename.String
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

	httputil.JSONResponse(w, httputil.NewResult(messages))
}

func postMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	channelId := chi.URLParam(r, "channelId")
	userId := serverutil.GetUserIdFromContext(ctx)

	var basemessage minichat.MessageBase
	err := httputil.JSONRequest(r, &basemessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if basemessage.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	conn := database.GetDatabase()

	var message minichat.Message
	var timestamp pgtype.Timestamptz
	err = conn.QueryRow(
		ctx,
		`
		INSERT INTO minichat.messages (channel_id, author_id, "content")
		VALUES ($1, $2, $3)
		RETURNING id, channel_id, author_id, "content", "timestamp", true
	`,
		channelId,
		userId,
		basemessage.Content,
	).Scan(&message.Id, &message.ChannelId, &message.AuthorId, &message.Content, &timestamp, &message.Read)
	message.Timestamp = coreutil.FormatTimestampz(timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update last read message timestamp
	_, err = conn.Exec(
		ctx,
		`UPDATE minichat.channels_members
		SET last_read_message_timestamp = $3
		WHERE user_id = $1 AND channel_id = $2`,
		userId,
		channelId,
		message.Timestamp,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.JSONResponse(w, message)
}

func patchMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	messageId := chi.URLParam(r, "messageId")
	userId := serverutil.GetUserIdFromContext(ctx)

	var basemessage minichat.MessageBase
	err := httputil.JSONRequest(r, &basemessage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if basemessage.Content == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	conn := database.GetDatabase()

	var message minichat.Message
	var timestamp pgtype.Timestamptz
	err = conn.QueryRow(
		ctx,
		"SELECT id, channel_id, author_id, content, timestamp, true FROM minichat.messages WHERE id = $1",
		messageId,
	).Scan(&message.Id, &message.ChannelId, &message.AuthorId, &message.Content, &timestamp, &message.Read)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message.Timestamp = coreutil.FormatTimestampz(timestamp)

	if message.AuthorId != userId {
		http.Error(w, "messages may only be edited by their author", http.StatusUnauthorized)
		return
	}

	_, err = conn.Exec(
		ctx,
		"UPDATE minichat.messages SET content = $1 WHERE id = $2",
		basemessage.Content,
		messageId,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message.Content = basemessage.Content

	httputil.JSONResponse(w, message)
}

func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	messageId := chi.URLParam(r, "messageId")
	userId := serverutil.GetUserIdFromContext(ctx)

	conn := database.GetDatabase()

	var message minichat.Message
	var timestamp pgtype.Timestamptz
	err := conn.QueryRow(
		ctx,
		"SELECT id, channel_id, author_id, content, timestamp, true FROM minichat.messages WHERE id = $1",
		messageId,
	).Scan(&message.Id, &message.ChannelId, &message.AuthorId, &message.Content, &timestamp, &message.Read)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message.Timestamp = coreutil.FormatTimestampz(timestamp)

	if message.AuthorId != userId {
		http.Error(w, "messages may only be deleted by their author", http.StatusUnauthorized)
		return
	}

	_, err = conn.Exec(
		ctx,
		"DELETE FROM minichat.messages WHERE id = $1",
		messageId,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httputil.JSONResponse(w, message)
}

func MessagesRouter(router chi.Router) {
	router.Get("/messages/{channelId}", getMessagesHandler)
	router.Post("/messages/{channelId}", postMessageHandler)
	router.Patch("/messages/{messageId}", patchMessageHandler)
	router.Delete("/messages/{messageId}", deleteMessageHandler)
}
