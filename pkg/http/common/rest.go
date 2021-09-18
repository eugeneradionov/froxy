package common

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	herrors "github.com/eugeneradionov/froxy/pkg/http/errors"
	"github.com/eugeneradionov/froxy/pkg/logger"
	"github.com/eugeneradionov/froxy/pkg/validator"
	"github.com/eugeneradionov/xerrors"
	v "github.com/go-playground/validator/v10"
)

// SendResponse - common method for encoding and writing any json response.
func SendResponse(w http.ResponseWriter, statusCode int, respBody interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	binRespBody, err := json.Marshal(respBody)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(binRespBody)
}

// SendError sends HTTP error with code selected based on the `xErr` ErrorCode.
func SendError(w http.ResponseWriter, xErr xerrors.XError) {
	xErr.Sanitize()
	SendResponse(w, herrors.GetHTTPCode(xErr), xErr)
}

// SendErrors sends HTTP errors
func SendErrors(w http.ResponseWriter, code int, xErrs xerrors.XErrors) {
	xErrs.Sanitize()
	SendResponse(w, code, xErrs)
}

// ProcessRequestBody - read and parse request body with errors description list.
func ProcessRequestBody(w http.ResponseWriter, r *http.Request, body interface{}) error {
	xErr := DecodeRequestBody(r, body)
	if xErr != nil {
		SendError(w, xErr)
		return xErr
	}

	xErrs := ValidateRequestBody(body)
	if xErrs != nil {
		SendErrors(w, http.StatusUnprocessableEntity, xErrs)
		return xErrs
	}

	return nil
}

func DecodeRequestBody(r *http.Request, body interface{}) xerrors.XError {
	defer r.Body.Close()

	return DecodeBody(r.Context(), r.Body, body)
}

func DecodeBody(ctx context.Context, r io.Reader, body interface{}) xerrors.XError {
	err := json.NewDecoder(r).Decode(body)
	if err != nil {
		serverError := herrors.NewBadRequestError(err, "Decode JSON body error")

		logger.Get().LogXError(ctx, serverError, "decode JSON body")
		return serverError
	}

	return nil
}

// ValidateRequestBody validates object with errors description list.
func ValidateRequestBody(body interface{}) xerrors.XErrors {
	err := validator.Get().Struct(body)
	if err != nil {
		validationErrors := err.(v.ValidationErrors)
		serverError := validator.FormatErrors(validationErrors)

		return serverError
	}

	return nil
}
