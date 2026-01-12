package errorx

import (
	"errors"
	"fmt"
	"net/http"
)

//
// ===== Sentinel Errors =====
//

// Custom error types following Go best practices
var (
	// ErrDataNotFound represents when requested data is not found
	ErrDataNotFound = errors.New("data not found")

	// ErrInvalidInput represents when input validation fails
	ErrInvalidInput = errors.New("invalid input")

	// ErrUnauthorized represents when user is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden represents when user doesn't have permission
	ErrForbidden = errors.New("forbidden")

	// ErrConflict represents when there's a conflict with existing data
	ErrConflict = errors.New("conflict")

	// ErrInternalServer represents internal server errors
	ErrInternalServer = errors.New("internal server error")

	// ErrUnprocessableEntity represents when the request is unprocessable
	ErrUnprocessableEntity = errors.New("unprocessable entity")

	ErrWorkerPoolAlreadyRunning = errors.New("worker pool is already running")
	ErrWorkerPoolNotRunning     = errors.New("worker pool is not running")
	ErrWorkerPoolQueueFull      = errors.New("worker pool queue is full")
	ErrWorkerPoolTimeout        = errors.New("worker pool operation timed out")
	ErrInvalidWorkerPoolConfig  = errors.New("invalid worker pool configuration")
	ErrCursorInvalid            = errors.New("invalid cursor")
	ErrCursorExpired            = errors.New("cursor expired")
)

//
// ===== ErrorType =====
//

// ErrorType represents the type of error for proper HTTP status mapping
type ErrorType int

const (
	ErrorTypeNotFound ErrorType = iota
	ErrorTypeInvalidInput
	ErrorTypeUnauthorized
	ErrorTypeForbidden
	ErrorTypeConflict
	ErrorTypeInternalServer
	ErrorTypeUnprocessableEntity
)

//
// ===== Numeric Error Code =====
//
// Format: XYYY
// 1xxx = input
// 2xxx = auth
// 3xxx = conflict
// 4xxx = not found
// 5xxx = internal
//

type ErrorCode int

const (
	CodeInvalidInput        ErrorCode = 1001
	CodeUnprocessableEntity ErrorCode = 1002

	CodeUnauthorized ErrorCode = 2001
	CodeForbidden    ErrorCode = 2003

	CodeConflict ErrorCode = 3001

	CodeNotFound ErrorCode = 4004

	CodeInternalServer ErrorCode = 5000
)

//
// ===== CustomError =====
//

// CustomError represents a custom error with type information
type CustomError struct {
	Type          ErrorType
	Code          ErrorCode
	Message       string // internal / log message
	PublicMessage string // safe for client
	Err           error
}

// Error implements the error interface
func (e *CustomError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *CustomError) Unwrap() error {
	return e.Err
}

//
// ===== Constructors =====
//

// NewCustomError creates a new custom error
func NewCustomError(errorType ErrorType, message string, err error) *CustomError {
	return &CustomError{
		Type:          errorType,
		Code:          codeFromType(errorType),
		Message:       message,
		PublicMessage: message,
		Err:           err,
	}
}

// NewDataNotFoundError creates a data not found error
func NewDataNotFoundError(resource string, err error) *CustomError {
	if err == nil {
		err = ErrDataNotFound
	}
	message := fmt.Sprintf("%s not found", resource)
	return NewCustomError(ErrorTypeNotFound, message, err)
}

// NewInvalidInputError creates an invalid input error
func NewInvalidInputError(message string, err error) *CustomError {
	if err == nil {
		err = ErrInvalidInput
	}
	return NewCustomError(ErrorTypeInvalidInput, message, err)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string, err error) *CustomError {
	if err == nil {
		err = ErrUnauthorized
	}
	return NewCustomError(ErrorTypeUnauthorized, message, err)
}

// NewForbiddenError creates a forbidden error
func NewForbiddenError(message string, err error) *CustomError {
	if err == nil {
		err = ErrForbidden
	}
	return NewCustomError(ErrorTypeForbidden, message, err)
}

// NewConflictError creates a conflict error
func NewConflictError(message string, err error) *CustomError {
	if err == nil {
		err = ErrConflict
	}
	return NewCustomError(ErrorTypeConflict, message, err)
}

// NewDuplicatePrerequisiteError creates a duplicate prerequisite error
func NewDuplicatePrerequisiteError(message string, err error) *CustomError {
	return NewCustomError(ErrorTypeConflict, message, err)
}

// NewInternalServerError creates an internal server error
func NewInternalServerError(message string, err error) *CustomError {
	if err == nil {
		err = ErrInternalServer
	}
	e := NewCustomError(ErrorTypeInternalServer, message, err)
	e.PublicMessage = "something went wrong"
	return e
}

func NewUnprocessableEntityError(message string, err error) *CustomError {
	if err == nil {
		err = ErrUnprocessableEntity
	}
	return NewCustomError(ErrorTypeUnprocessableEntity, message, err)
}

// NewBadRequestError creates a bad request error (alias for NewInvalidInputError)
func NewBadRequestError(message string, err error) *CustomError {
	return NewCustomError(ErrorTypeInvalidInput, message, err)
}

//
// ===== HTTP Mapping =====
//

// GetHTTPStatusFromError returns the appropriate HTTP status code for an error
func GetHTTPStatusFromError(err error) int {
	var customErr *CustomError
	if errors.As(err, &customErr) {
		return httpStatusFromCode(customErr.Code)
	}
	return http.StatusInternalServerError
}

func httpStatusFromCode(code ErrorCode) int {
	switch {
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest
	case code >= 2000 && code < 3000:
		if code == CodeUnauthorized {
			return http.StatusUnauthorized
		}
		return http.StatusForbidden
	case code >= 3000 && code < 4000:
		return http.StatusConflict
	case code >= 4000 && code < 5000:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

//
// ===== Helpers =====
//

func codeFromType(t ErrorType) ErrorCode {
	switch t {
	case ErrorTypeNotFound:
		return CodeNotFound
	case ErrorTypeInvalidInput:
		return CodeInvalidInput
	case ErrorTypeUnauthorized:
		return CodeUnauthorized
	case ErrorTypeForbidden:
		return CodeForbidden
	case ErrorTypeConflict:
		return CodeConflict
	case ErrorTypeUnprocessableEntity:
		return CodeUnprocessableEntity
	default:
		return CodeInternalServer
	}
}

// Existing helpers preserved

func IsDataNotFound(err error) bool {
	var customErr *CustomError
	return errors.As(err, &customErr) && customErr.Type == ErrorTypeNotFound
}

func IsInvalidInput(err error) bool {
	var customErr *CustomError
	return errors.As(err, &customErr) && customErr.Type == ErrorTypeInvalidInput
}

func IsUnauthorized(err error) bool {
	var customErr *CustomError
	return errors.As(err, &customErr) && customErr.Type == ErrorTypeUnauthorized
}

func IsForbidden(err error) bool {
	var customErr *CustomError
	return errors.As(err, &customErr) && customErr.Type == ErrorTypeForbidden
}

func IsConflict(err error) bool {
	var customErr *CustomError
	return errors.As(err, &customErr) && customErr.Type == ErrorTypeConflict
}

func IsInternalServer(err error) bool {
	var customErr *CustomError
	return errors.As(err, &customErr) && customErr.Type == ErrorTypeInternalServer
}

// Generic helper (new, optional)
func IsType(err error, t ErrorType) bool {
	var customErr *CustomError
	return errors.As(err, &customErr) && customErr.Type == t
}
