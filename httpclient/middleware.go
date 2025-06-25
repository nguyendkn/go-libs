package httpclient

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Built-in middlewares

// LoggingMiddleware logs requests and responses
type LoggingMiddleware struct {
	logger Logger
	config *LoggingConfig
}

// NewLoggingMiddleware tạo logging middleware
func NewLoggingMiddleware(logger Logger, config *LoggingConfig) *LoggingMiddleware {
	if config == nil {
		config = &LoggingConfig{
			Enabled: true,
			Level:   "info",
			Format:  "json",
		}
	}

	return &LoggingMiddleware{
		logger: logger,
		config: config,
	}
}

// Process implements Middleware interface
func (m *LoggingMiddleware) Process(req *Request, next Handler) (*Response, error) {
	if !m.config.Enabled {
		return next(req)
	}

	start := time.Now()

	// Log request
	m.logRequest(req)

	// Execute request
	resp, err := next(req)

	duration := time.Since(start)

	// Log response
	m.logResponse(req, resp, err, duration)

	return resp, err
}

func (m *LoggingMiddleware) logRequest(req *Request) {
	fields := []Field{
		{Key: "method", Value: req.Method},
		{Key: "url", Value: req.URL},
		{Key: "attempt", Value: req.attempt},
	}

	// Add headers if enabled
	if m.config.CollectHeaders {
		headers := make(map[string]string)
		for key, value := range req.Headers {
			// Mask sensitive headers
			if m.isSensitiveHeader(key) {
				headers[key] = "***"
			} else {
				headers[key] = value
			}
		}
		fields = append(fields, Field{Key: "headers", Value: headers})
	}

	// Add body if enabled and not too large
	if m.config.CollectBody && req.Body != nil {
		fields = append(fields, Field{Key: "body", Value: req.Body})
	}

	m.logger.Info("HTTP Request", fields...)
}

func (m *LoggingMiddleware) logResponse(req *Request, resp *Response, err error, duration time.Duration) {
	fields := []Field{
		{Key: "method", Value: req.Method},
		{Key: "url", Value: req.URL},
		{Key: "duration", Value: duration},
		{Key: "attempt", Value: req.attempt},
	}

	if err != nil {
		fields = append(fields, Field{Key: "error", Value: err.Error()})
		m.logger.Error("HTTP Request Failed", fields...)
		return
	}

	fields = append(fields,
		Field{Key: "status_code", Value: resp.StatusCode},
		Field{Key: "status", Value: resp.Status},
		Field{Key: "content_length", Value: resp.ContentLength},
	)

	// Add response headers if enabled
	if m.config.CollectHeaders {
		fields = append(fields, Field{Key: "response_headers", Value: resp.Headers})
	}

	// Add response body if enabled and not too large
	if m.config.CollectBody && len(resp.Body) > 0 && len(resp.Body) < 1024 {
		fields = append(fields, Field{Key: "response_body", Value: string(resp.Body)})
	}

	if resp.IsError() {
		m.logger.Warn("HTTP Request Error", fields...)
	} else {
		m.logger.Info("HTTP Request Success", fields...)
	}
}

func (m *LoggingMiddleware) isSensitiveHeader(key string) bool {
	key = strings.ToLower(key)
	for _, sensitive := range m.config.SensitiveHeaders {
		if strings.ToLower(sensitive) == key {
			return true
		}
	}

	// Default sensitive headers
	sensitiveHeaders := []string{"authorization", "cookie", "x-api-key", "x-auth-token"}
	for _, sensitive := range sensitiveHeaders {
		if sensitive == key {
			return true
		}
	}

	return false
}

// MetricsMiddleware collects metrics
type MetricsMiddleware struct {
	metrics Metrics
	config  *MetricsConfig
}

// NewMetricsMiddleware tạo metrics middleware
func NewMetricsMiddleware(metrics Metrics, config *MetricsConfig) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
		config:  config,
	}
}

// Process implements Middleware interface
func (m *MetricsMiddleware) Process(req *Request, next Handler) (*Response, error) {
	if !m.config.Enabled {
		return next(req)
	}

	start := time.Now()

	// Record request
	m.metrics.RecordRequest(req)

	// Execute request
	resp, err := next(req)

	duration := time.Since(start)

	// Record response or error
	if err != nil {
		m.metrics.RecordError(req, err, duration)
	} else {
		m.metrics.RecordResponse(req, resp, duration)
	}

	return resp, err
}

// TracingMiddleware adds distributed tracing
type TracingMiddleware struct {
	tracer Tracer
	config *TracingConfig
}

// NewTracingMiddleware tạo tracing middleware
func NewTracingMiddleware(tracer Tracer, config *TracingConfig) *TracingMiddleware {
	return &TracingMiddleware{
		tracer: tracer,
		config: config,
	}
}

// Process implements Middleware interface
func (m *TracingMiddleware) Process(req *Request, next Handler) (*Response, error) {
	if !m.config.Enabled {
		return next(req)
	}

	// Start span
	span := m.tracer.StartSpan(req)
	defer span.Finish()

	// Set span tags
	span.SetTag("http.method", string(req.Method))
	span.SetTag("http.url", req.URL)
	span.SetTag("component", "httpclient")

	// Inject tracing headers
	m.tracer.InjectHeaders(span, req.Headers)

	// Execute request
	resp, err := next(req)

	// Set response tags
	if err != nil {
		span.SetError(err)
		span.SetTag("error", true)
	} else {
		span.SetTag("http.status_code", resp.StatusCode)
		if resp.IsError() {
			span.SetTag("error", true)
		}
	}

	return resp, err
}

// RetryMiddleware handles retry logic
type RetryMiddleware struct {
	config *RetryPolicy
}

// NewRetryMiddleware tạo retry middleware
func NewRetryMiddleware(config *RetryPolicy) *RetryMiddleware {
	return &RetryMiddleware{
		config: config,
	}
}

// Process implements Middleware interface
func (m *RetryMiddleware) Process(req *Request, next Handler) (*Response, error) {
	// Use request-specific retry policy if available
	retryPolicy := req.RetryPolicy
	if retryPolicy == nil {
		retryPolicy = m.config
	}

	if retryPolicy == nil || retryPolicy.MaxAttempts <= 1 {
		return next(req)
	}

	var lastErr error
	var lastResp *Response

	for attempt := 1; attempt <= retryPolicy.MaxAttempts; attempt++ {
		req.attempt = attempt

		resp, err := next(req)

		// Success case
		if err == nil && !m.shouldRetryResponse(resp, retryPolicy) {
			return resp, nil
		}

		// Last attempt
		if attempt == retryPolicy.MaxAttempts {
			return resp, err
		}

		// Check if we should retry
		if !m.shouldRetryError(err, retryPolicy) && !m.shouldRetryResponse(resp, retryPolicy) {
			return resp, err
		}

		// Calculate delay
		delay := m.calculateDelay(attempt, retryPolicy)

		// Call retry callback
		if retryPolicy.OnRetry != nil {
			retryPolicy.OnRetry(attempt, err, delay)
		}

		// Wait before retry
		select {
		case <-time.After(delay):
		case <-req.Context.Done():
			return resp, req.Context.Err()
		}

		lastErr = err
		lastResp = resp
	}

	return lastResp, lastErr
}

func (m *RetryMiddleware) shouldRetryError(err error, policy *RetryPolicy) bool {
	if err == nil {
		return false
	}

	errMsg := err.Error()
	for _, retryableErr := range policy.RetryableErrors {
		if strings.Contains(errMsg, retryableErr) {
			return true
		}
	}

	return false
}

func (m *RetryMiddleware) shouldRetryResponse(resp *Response, policy *RetryPolicy) bool {
	if resp == nil {
		return false
	}

	for _, status := range policy.RetryableStatus {
		if resp.StatusCode == status {
			return true
		}
	}

	return false
}

func (m *RetryMiddleware) calculateDelay(attempt int, policy *RetryPolicy) time.Duration {
	delay := policy.InitialDelay
	for i := 1; i < attempt; i++ {
		delay = time.Duration(float64(delay) * policy.BackoffFactor)
		if delay > policy.MaxDelay {
			delay = policy.MaxDelay
			break
		}
	}

	// Add jitter if enabled
	if policy.Jitter {
		jitter := time.Duration(float64(delay) * 0.1 * (2*rand.Float64() - 1))
		delay += jitter
	}

	return delay
}

// TimeoutMiddleware handles request timeouts
type TimeoutMiddleware struct {
	timeout time.Duration
}

// NewTimeoutMiddleware tạo timeout middleware
func NewTimeoutMiddleware(timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{
		timeout: timeout,
	}
}

// Process implements Middleware interface
func (m *TimeoutMiddleware) Process(req *Request, next Handler) (*Response, error) {
	timeout := req.Timeout
	if timeout == 0 {
		timeout = m.timeout
	}

	if timeout == 0 {
		return next(req)
	}

	// Create timeout context
	ctx, cancel := context.WithTimeout(req.Context, timeout)
	defer cancel()

	// Update request context
	req.Context = ctx

	// Execute with timeout
	type result struct {
		resp *Response
		err  error
	}

	resultCh := make(chan result, 1)
	go func() {
		resp, err := next(req)
		resultCh <- result{resp: resp, err: err}
	}()

	select {
	case res := <-resultCh:
		return res.resp, res.err
	case <-ctx.Done():
		return nil, &HTTPError{
			Code:    1400,
			Message: fmt.Sprintf("request timeout after %v", timeout),
			Type:    "timeout",
		}
	}
}

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	auth *AuthConfig
}

// NewAuthMiddleware tạo auth middleware
func NewAuthMiddleware(auth *AuthConfig) *AuthMiddleware {
	return &AuthMiddleware{
		auth: auth,
	}
}

// Process implements Middleware interface
func (m *AuthMiddleware) Process(req *Request, next Handler) (*Response, error) {
	// Use request-specific auth if available
	auth := req.Auth
	if auth == nil {
		auth = m.auth
	}

	if auth == nil {
		return next(req)
	}

	// Apply authentication
	if err := m.applyAuth(req, auth); err != nil {
		return nil, err
	}

	return next(req)
}

func (m *AuthMiddleware) applyAuth(req *Request, auth *AuthConfig) error {
	switch auth.Type {
	case AuthTypeBasic:
		// Will be handled in buildHTTPRequest

	case AuthTypeBearer:
		req.Headers["Authorization"] = "Bearer " + auth.Token

	case AuthTypeAPIKey:
		if auth.Header != "" {
			req.Headers[auth.Header] = auth.APIKey
		} else {
			req.Headers["X-API-Key"] = auth.APIKey
		}

	case AuthTypeCustom:
		for key, value := range auth.Custom {
			req.Headers[key] = value
		}

	case AuthTypeOAuth2:
		// Check if token is expired and refresh if needed
		if auth.ExpiresAt.Before(time.Now()) && auth.RefreshFunc != nil {
			auth.Token = auth.RefreshFunc()
		}
		req.Headers["Authorization"] = "Bearer " + auth.Token
	}

	return nil
}
