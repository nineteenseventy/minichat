package util

import (
	"strings"
)

func ParseStringArray(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}
