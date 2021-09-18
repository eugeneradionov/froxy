package errors

import (
	"net/http"

	"github.com/eugeneradionov/xerrors"
)

func NewBadRequestError(err error, msg string, opts ...ErrorOption) *xerrors.XErr {
	return newError(err, msg, http.StatusNotFound, opts)
}

func NewNotFoundError(err error, msg string, opts ...ErrorOption) *xerrors.XErr {
	return newError(err, msg, http.StatusNotFound, opts)
}

func NewUnprocessableEntityError(err error, msg string, opts ...ErrorOption) *xerrors.XErr {
	return newError(err, msg, http.StatusUnprocessableEntity, opts)
}

func NewInternalServerError(err error, opts ...ErrorOption) *xerrors.XErr {
	return newError(err, "Internal Server Error", http.StatusInternalServerError, opts)
}

func newError(
	err error,
	msg string,
	httpCode int,
	opts []ErrorOption,
) *xerrors.XErr {
	if err == nil {
		return nil
	}
	xerr := xerrors.NewXErr(msg, err.Error(), nil, nil)

	opts = append(opts, WithHTTPCode(httpCode), WithInternalErr(err))
	applyOptions(xerr, opts)

	return xerr
}
