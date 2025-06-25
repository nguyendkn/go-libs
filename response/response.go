package response

import (
	"time"
)

// Response represents the standard API response structure
type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ErrorResponse represents an error response with additional error details
type ErrorResponse struct {
	Code      int                    `json:"code"`
	Message   string                 `json:"message"`
	Error     string                 `json:"error,omitempty"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// ValidationErrorResponse represents validation error response
type ValidationErrorResponse struct {
	Code      int                 `json:"code"`
	Message   string              `json:"message"`
	Errors    map[string][]string `json:"errors"`
	Timestamp time.Time           `json:"timestamp"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Code       int         `json:"code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Timestamp  time.Time   `json:"timestamp"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// ListResponse represents a simple list response without pagination
type ListResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Count     int         `json:"count"`
	Timestamp time.Time   `json:"timestamp"`
}

// Success Response Builders

// Success creates a successful response with data
func Success(data interface{}, message string) Response {
	if message == "" {
		message = "Success"
	}
	return Response{
		Code:      StatusOK,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Created creates a response for successful resource creation
func Created(data interface{}, message string) Response {
	if message == "" {
		message = "Resource created successfully"
	}
	return Response{
		Code:      StatusCreated,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Updated creates a response for successful resource update
func Updated(data interface{}, message string) Response {
	if message == "" {
		message = "Resource updated successfully"
	}
	return Response{
		Code:      StatusOK,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// Deleted creates a response for successful resource deletion
func Deleted(message string) Response {
	if message == "" {
		message = "Resource deleted successfully"
	}
	return Response{
		Code:      StatusOK,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// NoContent creates a response with no content
func NoContent(message string) Response {
	if message == "" {
		message = "No content"
	}
	return Response{
		Code:      StatusNoContent,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// List creates a simple list response without pagination
func List(data interface{}, count int, message string) ListResponse {
	if message == "" {
		message = "Data retrieved successfully"
	}
	return ListResponse{
		Code:      StatusOK,
		Message:   message,
		Data:      data,
		Count:     count,
		Timestamp: time.Now(),
	}
}

// Error Response Builders

// Error creates a generic error response
func Error(code int, message string, err error) ErrorResponse {
	errorMsg := ""
	if err != nil {
		errorMsg = err.Error()
	}
	return ErrorResponse{
		Code:      code,
		Message:   message,
		Error:     errorMsg,
		Timestamp: time.Now(),
	}
}

// BadRequest creates a 400 Bad Request response
func BadRequest(message string, err error) ErrorResponse {
	if message == "" {
		message = "Bad request"
	}
	return Error(StatusBadRequest, message, err)
}

// Unauthorized creates a 401 Unauthorized response
func Unauthorized(message string) ErrorResponse {
	if message == "" {
		message = "Unauthorized access"
	}
	return Error(StatusUnauthorized, message, nil)
}

// Forbidden creates a 403 Forbidden response
func Forbidden(message string) ErrorResponse {
	if message == "" {
		message = "Access forbidden"
	}
	return Error(StatusForbidden, message, nil)
}

// NotFound creates a 404 Not Found response
func NotFound(message string) ErrorResponse {
	if message == "" {
		message = "Resource not found"
	}
	return Error(StatusNotFound, message, nil)
}

// Conflict creates a 409 Conflict response
func Conflict(message string, err error) ErrorResponse {
	if message == "" {
		message = "Resource conflict"
	}
	return Error(StatusConflict, message, err)
}

// InternalServerError creates a 500 Internal Server Error response
func InternalServerError(message string, err error) ErrorResponse {
	if message == "" {
		message = "Internal server error"
	}
	return Error(StatusInternalServerError, message, err)
}

// ValidationError creates a validation error response
func ValidationError(errors map[string][]string, message string) ValidationErrorResponse {
	if message == "" {
		message = "Validation failed"
	}
	return ValidationErrorResponse{
		Code:      StatusBadRequest,
		Message:   message,
		Errors:    errors,
		Timestamp: time.Now(),
	}
}

// Pagination Functions

// NewPagination creates a new pagination object
func NewPagination(page, limit int, total int64) Pagination {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages < 1 {
		totalPages = 1
	}

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// Paginated creates a paginated response
func Paginated(data interface{}, page, limit int, total int64, message string) PaginatedResponse {
	if message == "" {
		message = "Data retrieved successfully"
	}

	pagination := NewPagination(page, limit, total)

	return PaginatedResponse{
		Code:       StatusOK,
		Message:    message,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now(),
	}
}

// GetOffset calculates the offset for database queries
func GetOffset(page, limit int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * limit
}

// Response Utilities

// WithDetails adds additional details to an error response
func (e ErrorResponse) WithDetails(details map[string]interface{}) ErrorResponse {
	e.Details = details
	return e
}

// WithError adds an error to the error response
func (e ErrorResponse) WithError(err error) ErrorResponse {
	if err != nil {
		e.Error = err.Error()
	}
	return e
}

// IsSuccessful checks if the response indicates success (2xx status codes)
func (r Response) IsSuccessful() bool {
	return IsSuccess(r.Code)
}

// IsError checks if the response indicates an error (4xx or 5xx status codes)
func (r Response) IsError() bool {
	return IsError(r.Code)
}

// Common Response Messages
const (
	MsgSuccess             = "Success"
	MsgCreated             = "Resource created successfully"
	MsgUpdated             = "Resource updated successfully"
	MsgDeleted             = "Resource deleted successfully"
	MsgRetrieved           = "Data retrieved successfully"
	MsgBadRequest          = "Bad request"
	MsgUnauthorized        = "Unauthorized access"
	MsgForbidden           = "Access forbidden"
	MsgNotFound            = "Resource not found"
	MsgConflict            = "Resource conflict"
	MsgValidationFailed    = "Validation failed"
	MsgInternalServerError = "Internal server error"
	MsgServiceUnavailable  = "Service temporarily unavailable"
)

// Quick Response Builders with predefined messages

// OK creates a simple success response
func OK() Response {
	return Success(nil, MsgSuccess)
}

// OKWithData creates a success response with data
func OKWithData(data interface{}) Response {
	return Success(data, MsgRetrieved)
}

// CreatedWithData creates a created response with data
func CreatedWithData(data interface{}) Response {
	return Created(data, MsgCreated)
}
