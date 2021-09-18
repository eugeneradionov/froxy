package proxy

import (
	"context"

	"github.com/eugeneradionov/froxy/config"
	"github.com/eugeneradionov/froxy/models"
	"github.com/eugeneradionov/froxy/services/proxy/uploader"
	"github.com/eugeneradionov/froxy/services/proxy/uploader/aws"
	"github.com/eugeneradionov/froxy/services/proxy/uploader/fs"
	"github.com/google/uuid"
)

type store interface {
	CreatStream(ctx context.Context) (*models.Stream, error)
	CreateChunk(ctx context.Context, streamID uuid.UUID, chunk *models.Chunk, position uint) error
	GetStream(ctx context.Context, streamID uuid.UUID) (*models.Stream, error)
}

type Service struct {
	cfg config.Proxy

	store store

	awsUp uploader.Uploader
	fsUp  uploader.Uploader
}

func New(cfg config.Proxy, store store) *Service {
	return &Service{
		cfg:   cfg,
		store: store,
		awsUp: aws.NewUploader(),
		fsUp:  fs.NewUploader(),
	}
}

func (srv *Service) GetMaxFileSize() int64 {
	return srv.cfg.FileMaxSizeMB * models.MEGABYTE
}
