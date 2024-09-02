package util

import "net/url"

func GetCdnUrl(
	bucket string,
	object string,
) (string, error) {
	args := GetArgs()
	return url.JoinPath(args.CdnUrl, bucket, object)
}
