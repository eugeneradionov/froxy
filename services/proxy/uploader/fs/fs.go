package fs

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var baseDir = os.TempDir()

type Uploader struct{}

func NewUploader() *Uploader {
	return &Uploader{}
}

func (f *Uploader) Upload(ctx context.Context, r io.Reader, filePath string) (_ string, err error) {
	filePath = filepath.Join(baseDir, filePath)

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0200)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return "", fmt.Errorf("write file: %w", err)
	}

	err = file.Chmod(0400)
	if err != nil {
		return "", fmt.Errorf("change file permissions: %w", err)
	}

	return file.Name(), nil
}
