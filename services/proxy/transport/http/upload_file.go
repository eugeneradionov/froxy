package http

import (
	"net/http"
	"strings"

	"github.com/eugeneradionov/froxy/models"
	"github.com/eugeneradionov/froxy/pkg/http/common"
	"github.com/eugeneradionov/froxy/pkg/http/errors"
	"github.com/eugeneradionov/xerrors"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	fileFieldName = "file"
	pathFieldName = "path"
)

func (h *Transport) uploadFile(w http.ResponseWriter, r *http.Request) { // nolint:funlen
	var (
		ctx = r.Context()

		storage  = chi.URLParam(r, "storage")
		compress = chi.URLParam(r, "compress")
	)

	xErr := h.validateStorage(models.StorageProvider(strings.ToLower(storage)))
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "storage provider is invalid")
		common.SendError(w, xErr)

		return
	}

	xErr = h.validateCompress(models.CompressAlgo(strings.ToLower(compress)))
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "compress algo is invalid")
		common.SendError(w, xErr)

		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.svc.GetMaxFileSize())
	err := r.ParseMultipartForm(h.svc.GetMaxFileSize())
	if err != nil {
		common.SendError(w, errors.NewUnprocessableEntityError(models.ErrFileTooLarge, models.ErrFileTooLarge.Error()))
		return
	}

	file, header, err := r.FormFile(fileFieldName)
	if err != nil {
		xErr := errors.NewUnprocessableEntityError(err, "read uploaded file")
		h.log.LogXError(ctx, xErr, "read uploaded file")
		common.SendError(w, xErr)

		return
	}
	defer file.Close()

	path := r.FormValue(pathFieldName)

	req := &models.UploadReq{
		Storage:  models.StorageProvider(storage),
		Compress: models.CompressAlgo(compress),
		File:     file,
		FileName: header.Filename,
		FilePath: path,
	}

	resp, xErr := h.svc.Upload(ctx, req)
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "upload file",
			zap.String("file_name", header.Filename),
			zap.String("path", path))
		common.SendError(w, errors.NewInternalServerError(xErr))

		return
	}

	common.SendResponse(w, http.StatusOK, resp)
}

func (h *Transport) validateStorage(storage models.StorageProvider) xerrors.XError {
	switch storage {
	case models.AWS, models.FS:
		return nil
	default:
		return errors.NewUnprocessableEntityError(
			models.ErrStorageProviderNotSupported,
			models.ErrStorageProviderNotSupported.Error(),
		)
	}
}

func (h *Transport) validateCompress(compress models.CompressAlgo) xerrors.XError {
	switch compress {
	case models.Zstd, models.Zip, "":
		return nil
	default:
		return errors.NewUnprocessableEntityError(
			models.ErrCompressAlgoNotSupported,
			models.ErrStorageProviderNotSupported.Error(),
		)
	}
}
