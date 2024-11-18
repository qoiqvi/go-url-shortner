package storage

import "errors"

// Объявление ошибок из пакета Error
var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists = errors.New("url exists")
)