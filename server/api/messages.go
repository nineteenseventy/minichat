package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go/v7"
	"github.com/nineteenseventy/minichat/core/database"
	httputil "github.com/nineteenseventy/minichat/core/http/util"
	"github.com/nineteenseventy/minichat/core/minichat"
	coreminio "github.com/nineteenseventy/minichat/core/minio"
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
			attachmentUrl, err := serverutil.ParseAttachmentUrl(message.Id, message.ChannelId, attachmentId.String, attachmentFilename.String)
			if err == nil {
				var attachment minichat.MessageAttachment
				attachment.Id = attachmentId.String
				attachment.MessageId = message.Id
				attachment.Type = attachmentType.String
				attachment.Url = &attachmentUrl
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

func getMessagesHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	channelId := chi.URLParam(request, "channelId")
	query := request.URL.Query()

	var start, count uint64 = 0, 10
	var err error

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

	var timestamp time.Time
	var isBefore bool
	var isAfter bool
	if before_param := query.Get("before"); before_param != "" {
		timestamp, err = coreutil.ParseTimestamp(before_param)
		if httputil.HandleError(writer, err) {
			return
		}
		isBefore = true
	} else if after_param := query.Get("after"); after_param != "" {
		timestamp, err = coreutil.ParseTimestamp(after_param)
		if httputil.HandleError(writer, err) {
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
	if httputil.HandleError(writer, err) {
		return
	}

	httputil.JSONResponse(writer, httputil.NewResult(messages))
}

func postMessageHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	channelId := chi.URLParam(request, "channelId")
	userId := serverutil.GetUserIdFromContext(ctx)

	var newMessage minichat.MessageBase
	err := httputil.JSONRequest(request, &newMessage)
	if httputil.HandleError(writer, err) {
		return
	}

	if newMessage.Content == "" {
		http.Error(writer, "content is required", http.StatusBadRequest)
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
		newMessage.Content,
	).Scan(&message.Id, &message.ChannelId, &message.AuthorId, &message.Content, &timestamp, &message.Read)
	message.Timestamp = coreutil.FormatTimestampz(timestamp)
	if httputil.HandleError(writer, err) {
		return
	}

	var errors []error
	// insert attachments
	for _, attachment := range newMessage.Attachments {
		attachmentId := coreutil.NewUuid()

		_, err = conn.Exec(
			ctx,
			`
			INSERT INTO minichat.attachments (id, message_id, "type", filename, url)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, message_id, "type", filename, url
		`,
			attachmentId,
			message.Id,
			attachment.Type,
			attachment.Filename,
			serverutil.ParseAttachmentKey(message.Id, channelId, attachmentId, attachment.Filename),
		)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		attachmentUrl, err := serverutil.ParseAttachmentUrl(message.Id, channelId, attachmentId, attachment.Filename)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		messageAttachment := minichat.MessageAttachment{
			Id:        attachmentId,
			MessageId: message.Id,
			Type:      attachment.Type,
			Filename:  attachment.Filename,
			Url:       &attachmentUrl,
		}
		message.Attachments = append(message.Attachments, messageAttachment)
	}

	if len(errors) > 0 {
		var errMessages []string
		for _, err := range errors {
			errMessages = append(errMessages, fmt.Sprintf("'%s'", err.Error()))
		}
		errorsString := strings.Join(errMessages, ", ")
		http.Error(writer, fmt.Sprintf("failed to insert attachments: %s", errorsString), http.StatusInternalServerError)
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
	if httputil.HandleError(writer, err) {
		return
	}

	httputil.JSONResponse(writer, message)
}

func postAttachmentHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	attachmentId := chi.URLParam(request, "attachmentId")
	userId := serverutil.GetUserIdFromContext(ctx)

	conn := database.GetDatabase()
	fmt.Println("1 attachmentId", attachmentId)

	// Get the attachment
	var attachmentKey, authorId string
	err := conn.QueryRow(
		ctx,
		`
		SELECT "attachment".url, "message".author_id
		FROM minichat.attachments AS "attachment"
		LEFT JOIN minichat.messages AS "message"
		ON "message".id = "attachment".message_id
		WHERE "attachment".id = $1
		`,
		attachmentId,
	).Scan(&attachmentKey, &authorId)
	if httputil.HandleError(writer, err) {
		return
	}

	if authorId != userId {
		http.Error(writer, "attachments may only be uploaded by their author", http.StatusUnauthorized)
		return
	}
	fmt.Println("2 attachmentId", attachmentId)

	minioClient := coreminio.GetMinio()

	// Post the attachment to Minio
	ContentType := request.Header.Get("Content-Type")
	if ContentType == "" {
		http.Error(writer, "Content-Type header is required", http.StatusBadRequest)
		return
	}

	size := request.ContentLength
	if size > 1024*1024 {
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

	fmt.Println("1 attachmentKey", attachmentKey)

	_, err = minioClient.PutObject(ctx, serverutil.AttachmentBucket, attachmentKey, request.Body, size, minio.PutObjectOptions{
		ContentType: ContentType,
	})
	if httputil.HandleError(writer, err) {
		return
	}

	fmt.Println("2 attachmentKey", attachmentKey)

	writer.WriteHeader(http.StatusCreated)
}

func patchMessageHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	messageId := chi.URLParam(request, "messageId")
	userId := serverutil.GetUserIdFromContext(ctx)

	var newMessage minichat.MessageBase
	err := httputil.JSONRequest(request, &newMessage)
	if httputil.HandleError(writer, err) {
		return
	}

	if newMessage.Content == "" {
		http.Error(writer, "content is required", http.StatusBadRequest)
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
	if httputil.HandleError(writer, err) {
		return
	}

	message.Timestamp = coreutil.FormatTimestampz(timestamp)

	if message.AuthorId != userId {
		http.Error(writer, "messages may only be edited by their author", http.StatusUnauthorized)
		return
	}

	_, err = conn.Exec(
		ctx,
		"UPDATE minichat.messages SET content = $1 WHERE id = $2",
		newMessage.Content,
		messageId,
	)
	if httputil.HandleError(writer, err) {
		return
	}

	message.Content = newMessage.Content

	httputil.JSONResponse(writer, message)
}

func deleteMessageHandler(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	messageId := chi.URLParam(request, "messageId")
	userId := serverutil.GetUserIdFromContext(ctx)

	conn := database.GetDatabase()

	var message minichat.Message
	var timestamp pgtype.Timestamptz
	err := conn.QueryRow(
		ctx,
		"SELECT id, channel_id, author_id, content, timestamp, true FROM minichat.messages WHERE id = $1",
		messageId,
	).Scan(&message.Id, &message.ChannelId, &message.AuthorId, &message.Content, &timestamp, &message.Read)
	if httputil.HandleError(writer, err) {
		return
	}

	message.Timestamp = coreutil.FormatTimestampz(timestamp)

	if message.AuthorId != userId {
		http.Error(writer, "messages may only be deleted by their author", http.StatusUnauthorized)
		return
	}

	_, err = conn.Exec(
		ctx,
		"DELETE FROM minichat.messages WHERE id = $1",
		messageId,
	)
	if httputil.HandleError(writer, err) {
		return
	}

	httputil.JSONResponse(writer, message)
}

func MessagesRouter(router chi.Router) {
	router.Get("/messages/{channelId}", getMessagesHandler)
	router.Post("/messages/{channelId}", postMessageHandler)
	router.Post("/attachment/{attachmentId}", postAttachmentHandler)
	router.Patch("/messages/{messageId}", patchMessageHandler)
	router.Delete("/messages/{messageId}", deleteMessageHandler)
}
