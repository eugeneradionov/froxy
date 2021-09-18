package compressor

import (
	"context"
	"io"
)

type Compressor interface {
	Decompress(ctx context.Context, r io.Reader) io.ReadCloser
}
