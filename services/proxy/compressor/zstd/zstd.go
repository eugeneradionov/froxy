package zstd

import (
	"context"
	"io"

	"github.com/DataDog/zstd"
)

type Compressor struct{}

func NewCompressor() *Compressor {
	return &Compressor{}
}

func (f *Compressor) Decompress(_ context.Context, r io.Reader) io.ReadCloser {
	return zstd.NewReader(r)
}
