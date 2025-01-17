package util

import (
	"context"
	"database/sql"

	"github.com/nineteenseventy/minichat/core/database"
	"github.com/nineteenseventy/minichat/core/logging"
	"github.com/nineteenseventy/minichat/core/minichat"
)

const ProfilePictureBucket = "profile"

func ParseUserPictureUrl(picture sql.NullString) *string {
	logger := logging.GetLogger("server.api.users.parsePictureUrl")
	if picture.Valid {
		pictureUrl, err := GetCdnUrl(ProfilePictureBucket, picture.String)
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

func GetUserMember(ctx context.Context, userId string, channelId string) (*minichat.Member, error) {
	conn := database.GetDatabase()

	var member minichat.Member
	err := conn.QueryRow(
		ctx,
		`
		SELECT "member".id, "member".user_id, "user".username
		FROM minichat.channels_members AS "member"
		JOIN minichat.users AS "user"
		ON "member".user_id = "user".id
		WHERE "member".channel_id = $1 AND "member".user_id = $2
		`,
		channelId,
		userId,
	).Scan(&member.ID, &member.UserId, &member.Username)
	if err != nil {
		return nil, err
	}

	return &member, nil
}
