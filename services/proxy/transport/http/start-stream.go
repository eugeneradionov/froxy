package http

import (
	"net/http"

	"github.com/eugeneradionov/froxy/pkg/http/common"
)

func (h *Transport) startStream(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	stream, xErr := h.svc.StartStream(ctx)
	if xErr != nil {
		h.log.LogXError(ctx, xErr, "start stream")
		common.SendError(w, xErr)

		return
	}

	common.SendResponse(w, http.StatusCreated, stream)
}
