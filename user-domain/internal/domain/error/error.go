package domainerror

import "errors"

var (
	ErrCodeNotFound = errors.New("the requested resource was not found in the system")

	ErrCodeInvalidInput = errors.New("the provided input data is invalid or does not meet business requirements")

	ErrCodeConflict = errors.New("the operation could not be completed due to a conflict with the current system state")

	ErrCodeForbidden = errors.New("you do not have permission to perform this action")

	ErrCodeInternal = errors.New("an unexpected internal server error occurred")
)
