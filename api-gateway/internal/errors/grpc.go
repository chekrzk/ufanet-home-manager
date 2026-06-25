package errors

import (
	"context"
	"errors"
	"net/http"
)

func FromGRPC(err error) *AppError {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return New(http.StatusGatewayTimeout, "upstream_timeout", "upstream service timeout")
	}
	if errors.Is(err, context.Canceled) {
		return New(499, "request_canceled", "request canceled")
	}
	return Wrap(err, ErrServiceUnavailable)
}
