package apperror

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidInput             = errors.New("invalid input")
	ErrUserDoesNotExist         = errors.New("user does not exist")
	ErrUserEmpty                = errors.New("user is required")
	ErrOrganizationDoesNotExist = errors.New("organization does not exist")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrForbidden                = errors.New("don't have enough permissions")
	ErrInternal                 = errors.New("internal error")
	ErrNotFound                 = errors.New("not found")
)

type AppError struct {
	Code    int
	Err     error
	Message string
}

func (h AppError) Unwrap() error {
	return h.Err
}

func (h AppError) Error() string {
	return h.Err.Error()
}

func BadRequest(err error) error {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
		Err:     err,
	}
}

func InternalServerError(err error) error {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: "internal_server_error",
		Err:     err,
	}
}

func Unauthorized(err error) error {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
		Err:     err,
	}
}

func Forbidden(err error) error {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: "forbidden",
		Err:     err,
	}
}

func NotFound(err error) error {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: "not_found",
		Err:     err,
	}
}
