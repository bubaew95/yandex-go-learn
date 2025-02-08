package repository

import "strings"

func isEmpty(t string) bool {
	return strings.TrimSpace(t) == ""
}
