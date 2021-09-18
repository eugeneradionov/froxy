package errors

import "github.com/eugeneradionov/xerrors"

const (
	extraField    = "field"
	extraHTTPCode = "http_code"

	intExtraErr = "error"
)

func GetHTTPCode(xerr xerrors.XError) int {
	return xerr.GetExtra()[extraHTTPCode].(int)
}
