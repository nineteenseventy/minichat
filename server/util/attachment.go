package util

import (
	"fmt"

	"github.com/nineteenseventy/minichat/core/logging"
)

const AttachmentBucket = "attachment"

func ParseAttachmentUrl(attachmentId string, filename string) (*string, error) {
	logger := logging.GetLogger("server.api.users.parsePictureUrl")
	attachmentUrl, err := GetCdnUrl(AttachmentBucket, fmt.Sprintf("%s/%s", attachmentId, filename))
	if err != nil {
		logger.Error().Err(err).Msg("failed to get picture url")
		return nil, err
	}
	return &attachmentUrl, nil
}
