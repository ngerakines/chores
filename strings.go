package chores

import "strings"

func normalize(value string) string {
	value = strings.ToLower(value)
	value = strings.TrimSpace(value)
	return value
}
