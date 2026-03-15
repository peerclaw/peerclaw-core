package errors

import (
	"errors"
	"fmt"
)

// Code represents a machine-readable error category.
type Code string

const (
	CodeNotFound    Code = "not_found"
	CodeValidation  Code = "validation_error"
	CodeAuth        Code = "auth_error"
	CodeConflict    Code = "conflict"
	CodeInternal    Code = "internal_error"
	CodeTimeout     Code = "timeout"
	CodeRateLimited Code = "rate_limited"
)

// Error is a structured error with a machine-readable code and human-readable message.
type Error struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	cause   error
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	return e.cause
}

// New creates a new Error with the given code and message.
func New(code Code, msg string) *Error {
	return &Error{Code: code, Message: msg}
}

// Wrap creates a new Error wrapping an underlying cause.
func Wrap(code Code, msg string, cause error) *Error {
	return &Error{Code: code, Message: msg, cause: cause}
}

// Is checks whether err is an *Error with the given Code.
func Is(err error, code Code) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}
