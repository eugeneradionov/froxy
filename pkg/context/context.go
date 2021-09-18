package context

import (
	"context"

	"github.com/go-chi/chi/v5/middleware"
)

type ctxKey string

// GetRequestID - returns Request ID from context.
func GetRequestID(ctx context.Context) string {
	return middleware.GetReqID(ctx)
}
