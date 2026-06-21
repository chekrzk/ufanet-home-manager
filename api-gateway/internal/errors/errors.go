package errors

import (
	stderrors "errors"
	"net/http"
)

var (
	ErrUnauthorized       = New(http.StatusUnauthorized, "unauthorized", "authorization required")
	ErrForbidden          = New(http.StatusForbidden, "forbidden", "access denied")
	ErrNotFound           = New(http.StatusNotFound, "not_found", "resource not found")
	ErrValidation         = New(http.StatusBadRequest, "validation_error", "request validation failed")
	ErrServiceUnavailable = New(http.StatusServiceUnavailable, "service_unavailable", "upstream service unavailable")
)

type AppError struct {
	Status  int               `json:"-"`
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
	Err     error             `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func Wrap(err error, fallback *AppError) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr
	}

	cp := *fallback
	cp.Err = err
	return &cp
}

func Validation(fields map[string]string) *AppError {
	err := *ErrValidation
	err.Fields = fields
	return &err
}
