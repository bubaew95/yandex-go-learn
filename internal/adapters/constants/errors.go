package constants

import "errors"

var (
	ErrUniqueIndex = errors.New("url already exists")
	ErrIsDeleted   = errors.New("url is deleted")
)
