package http

import (
	"fmt"
	"net/http"

	"github.com/eugeneradionov/froxy/pkg/http/common"
	herrors "github.com/eugeneradionov/froxy/pkg/http/errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Transport) getStream(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()

		streamID = chi.URLParam(r, "streamID")
	)

	streamUUID, err := uuid.Parse(streamID)
	if err != nil {
		common.SendError(w, herrors.NewUnprocessableEntityError(err, fmt.Sprintf("stream ID '%s' is invalid", streamID)))
		return
	}

	stream, xErr := h.svc.GetStream(ctx, streamUUID)
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "start stream")
		common.SendError(w, xErr)

		return
	}

	common.SendResponse(w, http.StatusCreated, stream)
}
