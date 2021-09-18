package models

import "io"

type UploadFileResp struct {
	Path string `json:"path"`
}

type UploadReq struct {
	Storage  StorageProvider
	Compress CompressAlgo

	File io.Reader

	FileName string
	FilePath string
}
