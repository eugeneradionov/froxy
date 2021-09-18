package http

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/eugeneradionov/froxy/models"
	"github.com/eugeneradionov/froxy/pkg/http/common"
	"github.com/eugeneradionov/froxy/pkg/http/errors"
	herrors "github.com/eugeneradionov/froxy/pkg/http/errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestBodyFieldName = "request"

func (h *Transport) addChunk(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()

		streamID = chi.URLParam(r, "streamID")

		req models.AddChunkReq
	)

	streamUUID, err := uuid.Parse(streamID)
	if err != nil {
		common.SendError(w, herrors.NewUnprocessableEntityError(err, fmt.Sprintf("stream ID '%s' is invalid", streamID)))
		return
	}

	body := r.FormValue(requestBodyFieldName)

	xErr := common.DecodeBody(r.Context(), bytes.NewReader([]byte(body)), &req)
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "decode body")
		common.SendError(w, xErr)
	}

	xErr = h.validateStorage(req.Storage)
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "storage provider is invalid")
		common.SendError(w, xErr)

		return
	}

	xErr = h.validateCompress(req.Compress)
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "compress algo is invalid")
		common.SendError(w, xErr)

		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, h.svc.GetMaxFileSize())
	err = r.ParseMultipartForm(h.svc.GetMaxFileSize())
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

	req.StreamID = streamUUID
	req.FileName = header.Filename
	req.File = file

	chunk, xErr := h.svc.AddChunk(ctx, req)
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "add chunk", zap.String("stream_id", streamUUID.String()))
		common.SendError(w, xErr)

		return
	}

	common.SendResponse(w, http.StatusCreated, chunk)
}
