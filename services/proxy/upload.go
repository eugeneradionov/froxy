package proxy

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/eugeneradionov/froxy/models"
	herrors "github.com/eugeneradionov/froxy/pkg/http/errors"
	"github.com/eugeneradionov/froxy/services/proxy/compressor"
	"github.com/eugeneradionov/froxy/services/proxy/compressor/zstd"
	"github.com/eugeneradionov/froxy/services/proxy/uploader"
	"github.com/eugeneradionov/xerrors"
)

func (srv *Service) Upload(ctx context.Context, req *models.UploadReq) (*models.UploadFileResp, xerrors.XError) {
	var (
		up   uploader.Uploader
		comp compressor.Compressor
	)

	switch req.Storage {
	case models.AWS:
		up = srv.awsUp
	case models.FS:
		up = srv.fsUp
	default:
		return nil, herrors.NewInternalServerError(fmt.Errorf("'%s' storage is not supported", req.Storage))
	}

	switch req.Compress {
	case models.Zstd:
		comp = zstd.NewCompressor()
	case "":
		// nothing to do
	default:
		return nil, herrors.NewInternalServerError(fmt.Errorf("'%s' compress algo is not supported", req.Storage))
	}

	if comp != nil {
		rc := comp.Decompress(ctx, req.File)
		req.File = rc
		defer rc.Close()
	}

	loadPath := filepath.Join(req.FilePath, req.FileName)

	savePath, err := up.Upload(ctx, req.File, loadPath)
	if err != nil {
		return nil, herrors.NewInternalServerError(err)
	}

	return &models.UploadFileResp{Path: savePath}, nil
}
