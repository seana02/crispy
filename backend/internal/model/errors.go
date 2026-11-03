package model

import "fmt"

// AppError represents a structured application error.
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

func NewNotFoundError(entity string, id any) *AppError {
    return &AppError{Code: 404, Message: fmt.Sprintf("%s with id %v not found", entity, id)}
}

func NewValidationError(msg string) *AppError {
    return &AppError{Code: 400, Message: msg}
}

func NewConflictError(msg string) *AppError {
    return &AppError{Code: 409, Message: msg}
}

func NewInternalError(msg string) *AppError {
    return &AppError{Code: 500, Message: msg}
}
