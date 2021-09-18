package proxy

import (
	"context"
	"errors"
	"fmt"

	"github.com/eugeneradionov/froxy/models"
	herrors "github.com/eugeneradionov/froxy/pkg/http/errors"
	"github.com/eugeneradionov/xerrors"
	"github.com/google/uuid"
)

func (srv *Service) StartStream(ctx context.Context) (*models.Stream, xerrors.XError) {
	stream, err := srv.store.CreatStream(ctx)
	if err != nil {
		return nil, herrors.NewInternalServerError(err)
	}

	return stream, nil
}

func (srv *Service) AddChunk(
	ctx context.Context,
	req models.AddChunkReq,
) (*models.Chunk, xerrors.XError) {
	upReq := &models.UploadReq{
		Storage:  req.Storage,
		Compress: req.Compress,
		File:     req.File,
		FileName: req.FileName,
		FilePath: req.FilePath,
	}

	upResp, xErr := srv.Upload(ctx, upReq)
	if xErr != nil {
		return nil, xErr
	}

	chunk := &models.Chunk{
		ID:       uuid.New(),
		FilePath: upResp.Path,
	}

	err := srv.store.CreateChunk(ctx, req.StreamID, chunk, req.Position)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, herrors.NewNotFoundError(err, fmt.Sprintf("stream with id '%s' not found", req.StreamID))
		}

		return nil, herrors.NewInternalServerError(fmt.Errorf("store create chunk: %w", err))
	}

	return chunk, nil
}

func (srv *Service) GetStream(ctx context.Context, streamID uuid.UUID) (*models.Stream, xerrors.XError) {
	stream, err := srv.store.GetStream(ctx, streamID)
	if err != nil {
		if errors.Is(err, models.ErrRecordNotFound) {
			return nil, herrors.NewNotFoundError(err, fmt.Sprintf("stream with id '%s' not found", streamID))
		}

		return nil, herrors.NewInternalServerError(fmt.Errorf("store get stream with id '%s': %w", streamID, err))
	}

	return stream, nil
}
