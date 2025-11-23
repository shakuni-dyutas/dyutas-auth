package app

import (
	"fmt"
)

type AppErrorCode string

const (
	// client errors
	BadRequestError AppErrorCode = "BAD_REQUEST"

	// system errors
	InternalServerError AppErrorCode = "INTERNAL_SERVER_ERROR"

	// auth errors
	UnauthenticatedError AppErrorCode = "UNAUTHENTICATED"
	UnauthorizedError    AppErrorCode = "UNAUTHORIZED"
)

type AppError struct {
	Code    AppErrorCode
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}
