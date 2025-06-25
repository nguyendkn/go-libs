package response

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyendkn/go-libs/logger"
	"go.uber.org/zap"
)

// Enhanced response functions with integrated logging

// SuccessJSONWithLog sends a successful JSON response with logging
func SuccessJSONWithLog(c *gin.Context, data interface{}, message string) {
	response := Success(data, message)

	// Log success response
	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.Int("status", StatusOK),
		logger.String("message", message),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	if userID := c.GetString("user_id"); userID != "" {
		fields = append(fields, logger.String("user_id", userID))
	}

	logger.Info("Success response", fields...)
	JSON(c, StatusOK, response)
}

// ErrorJSONWithLog sends an error JSON response with comprehensive logging
func ErrorJSONWithLog(c *gin.Context, code int, message string, err error) {
	response := Error(code, message, err)

	// Prepare log fields
	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.String("ip", c.ClientIP()),
		logger.Int("status", code),
		logger.String("message", message),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	if userID := c.GetString("user_id"); userID != "" {
		fields = append(fields, logger.String("user_id", userID))
	}

	if err != nil {
		fields = append(fields, logger.Err(err))
	}

	// Log based on error severity
	switch {
	case code >= 500:
		logger.Error("Server error response", fields...)
	case code >= 400:
		logger.Warn("Client error response", fields...)
	default:
		logger.Info("Error response", fields...)
	}

	JSON(c, code, response)
}

// BadRequestJSONWithLog sends a bad request JSON response with logging
func BadRequestJSONWithLog(c *gin.Context, message string, err error) {
	ErrorJSONWithLog(c, StatusBadRequest, message, err)
}

// UnauthorizedJSONWithLog sends an unauthorized JSON response with logging
func UnauthorizedJSONWithLog(c *gin.Context, message string) {
	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.String("ip", c.ClientIP()),
		logger.String("user_agent", c.Request.UserAgent()),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	logger.Warn("Unauthorized access attempt", fields...)
	ErrorJSONWithLog(c, StatusUnauthorized, message, nil)
}

// ForbiddenJSONWithLog sends a forbidden JSON response with logging
func ForbiddenJSONWithLog(c *gin.Context, message string) {
	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.String("ip", c.ClientIP()),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	if userID := c.GetString("user_id"); userID != "" {
		fields = append(fields, logger.String("user_id", userID))
	}

	logger.Warn("Forbidden access attempt", fields...)
	ErrorJSONWithLog(c, StatusForbidden, message, nil)
}

// NotFoundJSONWithLog sends a not found JSON response with logging
func NotFoundJSONWithLog(c *gin.Context, message string) {
	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.String("ip", c.ClientIP()),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	logger.Info("Resource not found", fields...)
	ErrorJSONWithLog(c, StatusNotFound, message, nil)
}

// InternalServerErrorJSONWithLog sends an internal server error JSON response with logging
func InternalServerErrorJSONWithLog(c *gin.Context, message string, err error) {
	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.String("ip", c.ClientIP()),
		logger.String("user_agent", c.Request.UserAgent()),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	if userID := c.GetString("user_id"); userID != "" {
		fields = append(fields, logger.String("user_id", userID))
	}

	if err != nil {
		fields = append(fields, logger.Err(err))
	}

	logger.Error("Internal server error", fields...)
	ErrorJSONWithLog(c, StatusInternalServerError, message, err)
}

// ValidationErrorJSONWithLog sends a validation error JSON response with logging
func ValidationErrorJSONWithLog(c *gin.Context, errors map[string][]string, message string) {
	response := ValidationError(errors, message)

	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.String("ip", c.ClientIP()),
		logger.Any("validation_errors", errors),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	logger.Warn("Validation error", fields...)
	JSON(c, StatusBadRequest, response)
}

// PaginatedJSONWithLog sends a paginated JSON response with logging
func PaginatedJSONWithLog(c *gin.Context, data interface{}, page, limit int, total int64, message string) {
	response := Paginated(data, page, limit, total, message)

	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.Int("page", page),
		logger.Int("limit", limit),
		logger.Int64("total", total),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	logger.Info("Paginated response", fields...)
	JSON(c, StatusOK, response)
}

// AbortWithErrorAndLog aborts the request with an error response and comprehensive logging
func AbortWithErrorAndLog(c *gin.Context, code int, message string, err error) {
	ErrorJSONWithLog(c, code, message, err)
	c.Abort()
}

// AbortWithInternalServerErrorAndLog aborts the request with an internal server error and logging
func AbortWithInternalServerErrorAndLog(c *gin.Context, message string, err error) {
	InternalServerErrorJSONWithLog(c, message, err)
	c.Abort()
}

// LogRequest logs incoming request details
func LogRequest(c *gin.Context) {
	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.String("ip", c.ClientIP()),
		logger.String("user_agent", c.Request.UserAgent()),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	if userID := c.GetString("user_id"); userID != "" {
		fields = append(fields, logger.String("user_id", userID))
	}

	if c.Request.URL.RawQuery != "" {
		fields = append(fields, logger.String("query", c.Request.URL.RawQuery))
	}

	logger.Info("Incoming request", fields...)
}

// LogResponse logs outgoing response details
func LogResponse(c *gin.Context, responseData interface{}) {
	fields := []zap.Field{
		logger.String("method", c.Request.Method),
		logger.String("path", c.Request.URL.Path),
		logger.Int("status", c.Writer.Status()),
		logger.Int("size", c.Writer.Size()),
	}

	if requestID := c.GetString("request_id"); requestID != "" {
		fields = append(fields, logger.String("request_id", requestID))
	}

	logger.Info("Outgoing response", fields...)
}
