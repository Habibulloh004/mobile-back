package utils

import (
	"errors"
	"fmt"
)

// Application errors
var (
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrUserNotFound          = errors.New("user not found")
	ErrUnauthorized          = errors.New("unauthorized access")
	ErrForbidden             = errors.New("forbidden action")
	ErrInvalidToken          = errors.New("invalid token")
	ErrInternalServer        = errors.New("internal server error")
	ErrInvalidInput          = errors.New("invalid input data")
	ErrResourceNotFound      = errors.New("resource not found")
	ErrResourceAlreadyExists = errors.New("resource already exists")
	ErrImageUpload           = errors.New("image upload failed")
)

// AppError represents an application error
type AppError struct {
	Err     error
	Message string
	Code    int
}

// Error returns the error message
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// NewAppError creates a new application error
func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource string, id interface{}) *AppError {
	return &AppError{
		Err:     ErrResourceNotFound,
		Message: fmt.Sprintf("%s with ID %v not found", resource, id),
		Code:    404,
	}
}

// NewUnauthorizedError creates a new unauthorized error
func NewUnauthorizedError() *AppError {
	return &AppError{
		Err:     ErrUnauthorized,
		Message: "Unauthorized access",
		Code:    401,
	}
}

// NewForbiddenError creates a new forbidden error
func NewForbiddenError() *AppError {
	return &AppError{
		Err:     ErrForbidden,
		Message: "Forbidden action",
		Code:    403,
	}
}

// NewInvalidInputError creates a new invalid input error
func NewInvalidInputError(message string) *AppError {
	if message == "" {
		message = "Invalid input data"
	}
	return &AppError{
		Err:     ErrInvalidInput,
		Message: message,
		Code:    400,
	}
}

// NewInternalServerError creates a new internal server error
func NewInternalServerError(err error) *AppError {
	return &AppError{
		Err:     err,
		Message: "Internal server error",
		Code:    500,
	}
}
