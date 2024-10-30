package util

import "fmt"

func ParseHost(host string, port uint16) string {
	if host == "*" {
		host = ""
	}
	return fmt.Sprintf("%s:%d", host, port)
}
