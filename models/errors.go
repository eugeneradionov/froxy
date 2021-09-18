package models

import "errors"

var (
	ErrFileTooLarge                = errors.New("file exceeds maximum allowed size")
	ErrStorageProviderNotSupported = errors.New("storage provider is not supported")
	ErrCompressAlgoNotSupported    = errors.New("compress algorithm is not supported")

	ErrRecordNotFound = errors.New("record not found")
)
