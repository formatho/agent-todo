package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error
type AppError struct {
	Code       int                    `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Internal   error                  `json:"-"`
	HTTPStatus int                    `json:"-"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Internal)
	}
	return e.Message
}

// NewAppError creates a new AppError
func NewAppError(httpStatus, code int, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// WithDetails adds details to the error
func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	e.Details = details
	return e
}

// WithInternal adds an internal error
func (e *AppError) WithInternal(err error) *AppError {
	e.Internal = err
	return e
}

// Common error constructors
func BadRequest(message string) *AppError {
	return &AppError{
		Code:       http.StatusBadRequest,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:       http.StatusUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func Forbidden(message string) *AppError {
	return &AppError{
		Code:       http.StatusForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

func NotFound(message string) *AppError {
	return &AppError{
		Code:       http.StatusNotFound,
		Message:    message,
		HTTPStatus: http.StatusNotFound,
	}
}

func Conflict(message string) *AppError {
	return &AppError{
		Code:       http.StatusConflict,
		Message:    message,
		HTTPStatus: http.StatusConflict,
	}
}

func InternalServerError(message string) *AppError {
	return &AppError{
		Code:       http.StatusInternalServerError,
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
	}
}

// Predefined errors
var (
	ErrInvalidCredentials = Unauthorized("Invalid email or password")
	ErrUserNotFound       = NotFound("User not found")
	ErrAgentNotFound      = NotFound("Agent not found")
	ErrTaskNotFound       = NotFound("Task not found")
	ErrProjectNotFound    = NotFound("Project not found")
	ErrUnauthorized       = Unauthorized("Unauthorized")
	ErrInvalidToken       = Unauthorized("Invalid token")
	ErrTokenExpired       = Unauthorized("Token expired")
)
