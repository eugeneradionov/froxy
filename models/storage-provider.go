package models

type StorageProvider string

const (
	AWS StorageProvider = "aws"
	FS  StorageProvider = "fs"
)
