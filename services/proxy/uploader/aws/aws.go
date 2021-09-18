package aws

import (
	"context"
	"io"
)

type Uploader struct{}

func NewUploader() *Uploader {
	return &Uploader{}
}

func (f *Uploader) Upload(ctx context.Context, r io.Reader, path string) (_ string, err error) {
	// TODO: IMPLEMENT!
	return "", nil
}
