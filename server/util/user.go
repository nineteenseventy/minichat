package util

import (
	"context"
	"database/sql"

	"github.com/nineteenseventy/minichat/core/logging"
	"github.com/nineteenseventy/minichat/core/minichat"
)

func ParseUserPictureUrl(picture sql.NullString) *string {
	logger := logging.GetLogger("server.api.users.parsePictureUrl")
	if picture.Valid {
		pictureUrl, err := GetCdnUrl("profile", picture.String)
		if err != nil {
			logger.Error().Err(err).Msg("failed to get picture url")
			return nil
		}
		return &pictureUrl
	}
	return nil
}

func GetUserFromContext(ctx context.Context) minichat.UserProfile {
	return ctx.Value(minichat.UserProfileContextKey{}).(minichat.UserProfile)
}
