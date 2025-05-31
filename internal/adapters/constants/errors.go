package constants

import "errors"

// Ошибки.
var (
	ErrUniqueIndex      = errors.New("url already exists") // Такой url уже существует
	ErrIsDeleted        = errors.New("url is deleted")     // Url удален
	ErrNotTrustedSubnet = errors.New("not trusted subnet")
)
