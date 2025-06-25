package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nguyendkn/go-libs/logger"
	"go.uber.org/zap"
)

// Gin Integration Helpers

// JSON sends a JSON response using Gin context
func JSON(c *gin.Context, httpStatus int, response interface{}) {
	c.JSON(httpStatus, response)
}

// SuccessJSON sends a successful JSON response
func SuccessJSON(c *gin.Context, data interface{}, message string) {
	response := Success(data, message)
	JSON(c, http.StatusOK, response)
}

// CreatedJSON sends a created JSON response
func CreatedJSON(c *gin.Context, data interface{}, message string) {
	response := Created(data, message)
	JSON(c, http.StatusCreated, response)
}

// UpdatedJSON sends an updated JSON response
func UpdatedJSON(c *gin.Context, data interface{}, message string) {
	response := Updated(data, message)
	JSON(c, http.StatusOK, response)
}

// DeletedJSON sends a deleted JSON response
func DeletedJSON(c *gin.Context, message string) {
	response := Deleted(message)
	JSON(c, http.StatusOK, response)
}

// NoContentJSON sends a no content JSON response
func NoContentJSON(c *gin.Context, message string) {
	response := NoContent(message)
	JSON(c, http.StatusNoContent, response)
}

// ErrorJSON sends an error JSON response with logging
func ErrorJSON(c *gin.Context, code int, message string, err error) {
	response := Error(code, message, err)

	// Log the error with context using logger package
	logFields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
		zap.Int("status", code),
		zap.String("message", message),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		logFields = append(logFields, zap.String("request_id", requestID))
	}

	if err != nil {
		logFields = append(logFields, zap.Error(err))
	}

	// Log based on error type
	if code >= 500 {
		// Server errors - log as error
		logger.GetLogger().Error("API error response", logFields...)
	} else if code >= 400 {
		// Client errors - log as warning
		logger.GetLogger().Warn("API error response", logFields...)
	}

	JSON(c, code, response)
}

// BadRequestJSON sends a bad request JSON response
func BadRequestJSON(c *gin.Context, message string, err error) {
	response := BadRequest(message, err)
	JSON(c, http.StatusBadRequest, response)
}

// UnauthorizedJSON sends an unauthorized JSON response
func UnauthorizedJSON(c *gin.Context, message string) {
	response := Unauthorized(message)
	JSON(c, http.StatusUnauthorized, response)
}

// ForbiddenJSON sends a forbidden JSON response
func ForbiddenJSON(c *gin.Context, message string) {
	response := Forbidden(message)
	JSON(c, http.StatusForbidden, response)
}

// NotFoundJSON sends a not found JSON response
func NotFoundJSON(c *gin.Context, message string) {
	response := NotFound(message)
	JSON(c, http.StatusNotFound, response)
}

// ConflictJSON sends a conflict JSON response
func ConflictJSON(c *gin.Context, message string, err error) {
	response := Conflict(message, err)
	JSON(c, http.StatusConflict, response)
}

// InternalServerErrorJSON sends an internal server error JSON response
func InternalServerErrorJSON(c *gin.Context, message string, err error) {
	response := InternalServerError(message, err)
	JSON(c, http.StatusInternalServerError, response)
}

// ValidationErrorJSON sends a validation error JSON response
func ValidationErrorJSON(c *gin.Context, errors map[string][]string, message string) {
	response := ValidationError(errors, message)
	JSON(c, http.StatusBadRequest, response)
}

// PaginatedJSON sends a paginated JSON response
func PaginatedJSON(c *gin.Context, data interface{}, page, limit int, total int64, message string) {
	response := Paginated(data, page, limit, total, message)
	JSON(c, http.StatusOK, response)
}

// ListJSON sends a list JSON response
func ListJSON(c *gin.Context, data interface{}, count int, message string) {
	response := List(data, count, message)
	JSON(c, http.StatusOK, response)
}

// Pagination Helper Functions for Gin

// GetPaginationParams extracts pagination parameters from Gin context
func GetPaginationParams(c *gin.Context) (page, limit int) {
	page = 1
	limit = 10

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// Set reasonable limits
	if limit > 100 {
		limit = 100
	}

	return page, limit
}

// GetPaginationParamsWithDefaults extracts pagination parameters with custom defaults
func GetPaginationParamsWithDefaults(c *gin.Context, defaultPage, defaultLimit int) (page, limit int) {
	page = defaultPage
	limit = defaultLimit

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	return page, limit
}

// Quick Response Functions for Gin

// OKJSON sends a simple success JSON response
func OKJSON(c *gin.Context) {
	SuccessJSON(c, nil, MsgSuccess)
}

// OKWithDataJSON sends a success JSON response with data
func OKWithDataJSON(c *gin.Context, data interface{}) {
	SuccessJSON(c, data, MsgRetrieved)
}

// CreatedWithDataJSON sends a created JSON response with data
func CreatedWithDataJSON(c *gin.Context, data interface{}) {
	CreatedJSON(c, data, MsgCreated)
}

// AbortWithError aborts the request with an error response
func AbortWithError(c *gin.Context, code int, message string, err error) {
	ErrorJSON(c, code, message, err)
	c.Abort()
}

// AbortWithBadRequest aborts the request with a bad request response
func AbortWithBadRequest(c *gin.Context, message string, err error) {
	BadRequestJSON(c, message, err)
	c.Abort()
}

// AbortWithUnauthorized aborts the request with an unauthorized response
func AbortWithUnauthorized(c *gin.Context, message string) {
	UnauthorizedJSON(c, message)
	c.Abort()
}

// AbortWithForbidden aborts the request with a forbidden response
func AbortWithForbidden(c *gin.Context, message string) {
	ForbiddenJSON(c, message)
	c.Abort()
}

// AbortWithNotFound aborts the request with a not found response
func AbortWithNotFound(c *gin.Context, message string) {
	NotFoundJSON(c, message)
	c.Abort()
}

// AbortWithInternalServerError aborts the request with an internal server error response
func AbortWithInternalServerError(c *gin.Context, message string, err error) {
	InternalServerErrorJSON(c, message, err)
	c.Abort()
}
