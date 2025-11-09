package app

import "fmt"

type AppErrorCode string

const (
	AppErrorCodeUnknown AppErrorCode = "UNKNOWN"
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
