package util

import (
	"context"
	"database/sql"

	"github.com/nineteenseventy/minichat/core/logging"
)

const ProfulePictureBucket = "profile"

func ParseUserPictureUrl(picture sql.NullString) *string {
	logger := logging.GetLogger("server.api.users.parsePictureUrl")
	if picture.Valid {
		pictureUrl, err := GetCdnUrl(ProfulePictureBucket, picture.String)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get picture url")
			return nil
		}
		return &pictureUrl
	}
	return nil
}

type UserIdContextKey struct{}

func GetUserIdFromContext(ctx context.Context) string {
	return ctx.Value(UserIdContextKey{}).(string)
}
