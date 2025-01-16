package util

import (
	"fmt"

	"github.com/nineteenseventy/minichat/core/logging"
)

const AttachmentBucket = "attachment"

func ParseAttachmentKey(messageId string, channelId string, attachmentId string, filename string) string {
	return fmt.Sprintf("%s/%s/%s/%s", channelId, messageId, attachmentId, filename)
}

func ParseAttachmentUrl(messageId string, channelId string, attachmentId string, filename string) (string, error) {
	logger := logging.GetLogger("server.api.users.parsePictureUrl")
	attachmentUrl, err := GetCdnUrl(AttachmentBucket, ParseAttachmentKey(messageId, channelId, attachmentId, filename))
	if err != nil {
		logger.Error().Err(err).Msg("failed to get picture url")
		return "", err
	}
	return attachmentUrl, nil
}
