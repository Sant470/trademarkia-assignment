package store

import "fmt"

const (
	_prefix = "service:trademark"
)

func userKey(username string) string {
	return fmt.Sprintf("%s:username:%s", _prefix, username)
}
