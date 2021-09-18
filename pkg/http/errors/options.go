package errors

import "github.com/eugeneradionov/xerrors"

type ErrorOption func(err *xerrors.XErr)

func WithField(field string) ErrorOption {
	return func(xerr *xerrors.XErr) {
		xerr.Extra = setExtra(xerr.Extra, extraField, field)
	}
}

func WithHTTPCode(code int) ErrorOption {
	return func(xerr *xerrors.XErr) {
		xerr.Extra = setExtra(xerr.Extra, extraHTTPCode, code)
	}
}

func WithInternalErr(err error) ErrorOption {
	return func(xerr *xerrors.XErr) {
		xerr.InternalExtra = setExtra(xerr.InternalExtra, intExtraErr, err)
	}
}

func applyOptions(xerr *xerrors.XErr, opts []ErrorOption) {
	for _, opt := range opts {
		opt(xerr)
	}
}

// nolint
func setExtra(extra map[string]interface{}, key string, val interface{}) map[string]interface{} {
	if extra == nil {
		extra = map[string]interface{}{
			key: val,
		}
	} else {
		extra[key] = val
	}

	return extra
}
