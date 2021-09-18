package uploader

import (
	"context"
	"io"
)

type Uploader interface {
	Upload(ctx context.Context, r io.Reader, filePath string) (savePath string, err error)
}
