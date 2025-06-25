package httpclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"maps"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
)

// httpClient implements the Client interface
type httpClient struct {
	config      *ClientConfig
	httpClient  *http.Client
	middlewares []Middleware

	// Components
	cache          Cache
	circuitBreaker CircuitBreaker
	rateLimiter    RateLimiter
	metrics        Metrics
	logger         Logger
	tracer         Tracer

	// Synchronization
	mu sync.RWMutex
}

// NewClient tạo một HTTP client mới
func NewClient(config *ClientConfig) Client {
	if config == nil {
		config = DefaultConfig()
	}

	client := &httpClient{
		config:      config,
		middlewares: make([]Middleware, 0),
	}

	// Setup HTTP client
	client.setupHTTPClient()

	// Setup components
	client.setupComponents()

	return client
}

// DefaultConfig trả về cấu hình mặc định
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		UserAgent:       DefaultUserAgent,
		Headers:         make(map[string]string),
		Timeout:         DefaultTimeoutConfig,
		Retry:           DefaultRetryPolicy,
		ConnectionPool:  DefaultConnectionPoolConfig,
		FollowRedirects: true,
		MaxRedirects:    10,
		Debug:           false,
	}
}

// setupHTTPClient thiết lập HTTP client
func (c *httpClient) setupHTTPClient() {
	// Set default values if nil
	if c.config.Timeout == nil {
		c.config.Timeout = DefaultTimeoutConfig
	}
	if c.config.ConnectionPool == nil {
		c.config.ConnectionPool = DefaultConnectionPoolConfig
	}

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   c.config.Timeout.Connect,
			KeepAlive: c.config.Timeout.KeepAlive,
		}).DialContext,
		MaxIdleConns:        c.config.ConnectionPool.MaxIdleConns,
		MaxIdleConnsPerHost: c.config.ConnectionPool.MaxIdleConnsPerHost,
		MaxConnsPerHost:     c.config.ConnectionPool.MaxConnsPerHost,
		IdleConnTimeout:     c.config.ConnectionPool.IdleConnTimeout,
		DisableKeepAlives:   c.config.ConnectionPool.DisableKeepAlives,
		DisableCompression:  c.config.ConnectionPool.DisableCompression,
	}

	// Setup TLS
	if c.config.TLS != nil {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: c.config.TLS.InsecureSkipVerify,
			ServerName:         c.config.TLS.ServerName,
			MinVersion:         c.config.TLS.MinVersion,
			MaxVersion:         c.config.TLS.MaxVersion,
			CipherSuites:       c.config.TLS.CipherSuites,
		}
	}

	// Setup proxy
	if c.config.Proxy != nil && c.config.Proxy.URL != "" {
		proxyURL, err := url.Parse(c.config.Proxy.URL)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	c.httpClient = &http.Client{
		Transport: transport,
		Timeout:   c.config.Timeout.Request,
	}

	// Setup redirect policy
	if !c.config.FollowRedirects {
		c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else if c.config.MaxRedirects > 0 {
		c.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if len(via) >= c.config.MaxRedirects {
				return ErrTooManyRedirects
			}
			return nil
		}
	}
}

// setupComponents thiết lập các components
func (c *httpClient) setupComponents() {
	// Setup cache
	if c.config.Cache != nil && c.config.Cache.Enabled {
		c.cache = NewMemoryCache(c.config.Cache)
	}

	// Setup circuit breaker
	if c.config.CircuitBreaker != nil && c.config.CircuitBreaker.Enabled {
		c.circuitBreaker = NewCircuitBreaker(c.config.CircuitBreaker)
	}

	// Setup rate limiter
	if c.config.RateLimit != nil && c.config.RateLimit.Enabled {
		c.rateLimiter = NewRateLimiter(c.config.RateLimit)
	}

	// Setup metrics
	if c.config.Metrics != nil && c.config.Metrics.Enabled {
		c.metrics = NewMetrics(c.config.Metrics)
	}

	// Setup logger
	if c.config.Logging != nil && c.config.Logging.Enabled {
		c.logger = NewLogger(c.config.Logging)
	}

	// Setup tracer
	if c.config.Tracing != nil && c.config.Tracing.Enabled {
		c.tracer = NewTracer(c.config.Tracing)
	}
}

// Core HTTP methods
func (c *httpClient) Get(url string) RequestBuilder {
	return c.Request(MethodGET, url)
}

func (c *httpClient) Post(url string) RequestBuilder {
	return c.Request(MethodPOST, url)
}

func (c *httpClient) Put(url string) RequestBuilder {
	return c.Request(MethodPUT, url)
}

func (c *httpClient) Patch(url string) RequestBuilder {
	return c.Request(MethodPATCH, url)
}

func (c *httpClient) Delete(url string) RequestBuilder {
	return c.Request(MethodDELETE, url)
}

func (c *httpClient) Head(url string) RequestBuilder {
	return c.Request(MethodHEAD, url)
}

func (c *httpClient) Options(url string) RequestBuilder {
	return c.Request(MethodOPTIONS, url)
}

// Request tạo request builder
func (c *httpClient) Request(method HTTPMethod, url string) RequestBuilder {
	return NewRequestBuilder(c, method, url)
}

// Do thực hiện request
func (c *httpClient) Do(req *Request) (*Response, error) {
	return c.DoWithContext(context.Background(), req)
}

// DoWithContext thực hiện request với context
func (c *httpClient) DoWithContext(ctx context.Context, req *Request) (*Response, error) {
	// Set context if not provided
	if req.Context == nil {
		req.Context = ctx
	}

	// Apply default configuration
	c.applyDefaults(req)

	// Create handler chain
	handler := c.createHandler()

	// Apply middlewares
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		middleware := c.middlewares[i]
		currentHandler := handler
		handler = func(r *Request) (*Response, error) {
			return middleware.Process(r, currentHandler)
		}
	}

	// Execute request
	return handler(req)
}

// applyDefaults áp dụng cấu hình mặc định
func (c *httpClient) applyDefaults(req *Request) {
	// Apply base URL
	if c.config.BaseURL != "" && !isAbsoluteURL(req.URL) {
		req.URL = c.config.BaseURL + req.URL
	}

	// Apply default headers
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}

	// Set User-Agent
	if _, exists := req.Headers["User-Agent"]; !exists && c.config.UserAgent != "" {
		req.Headers["User-Agent"] = c.config.UserAgent
	}

	// Apply global headers
	for key, value := range c.config.Headers {
		if _, exists := req.Headers[key]; !exists {
			req.Headers[key] = value
		}
	}

	// Apply default auth
	if req.Auth == nil && c.config.Auth != nil {
		req.Auth = c.config.Auth
	}

	// Apply default timeout
	if req.Timeout == 0 && c.config.Timeout != nil {
		req.Timeout = c.config.Timeout.Request
	}

	// Apply default retry policy
	if req.RetryPolicy == nil && c.config.Retry != nil {
		req.RetryPolicy = c.config.Retry
	}

	// Apply redirect settings
	if c.config.FollowRedirects {
		req.FollowRedirects = true
		if req.MaxRedirects == 0 {
			req.MaxRedirects = c.config.MaxRedirects
		}
	}

	// Set start time
	req.startTime = time.Now()
}

// createHandler tạo handler chính
func (c *httpClient) createHandler() Handler {
	return func(req *Request) (*Response, error) {
		// Check cache first
		if c.cache != nil && !req.NoCache && req.Method == MethodGET {
			cacheKey := c.getCacheKey(req)
			if cached, found := c.cache.Get(cacheKey); found {
				cached.FromCache = true
				return cached, nil
			}
		}

		// Check circuit breaker
		if c.circuitBreaker != nil {
			return c.circuitBreaker.Execute(req, c.executeRequest)
		}

		return c.executeRequest(req)
	}
}

// executeRequest thực hiện request thực tế
func (c *httpClient) executeRequest(req *Request) (*Response, error) {
	// Check rate limit
	if c.rateLimiter != nil {
		key := c.getRateLimitKey(req)
		if !c.rateLimiter.Allow(key) {
			delay := c.rateLimiter.Wait(key)
			if delay > 0 {
				select {
				case <-time.After(delay):
				case <-req.Context.Done():
					return nil, req.Context.Err()
				}
			}
		}
	}

	// Execute with retry
	return c.executeWithRetry(req)
}

// executeWithRetry thực hiện request với retry logic
func (c *httpClient) executeWithRetry(req *Request) (*Response, error) {
	var lastErr error
	var lastResp *Response

	maxAttempts := 1
	if req.RetryPolicy != nil {
		maxAttempts = req.RetryPolicy.MaxAttempts
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		req.attempt = attempt

		// Execute single request
		resp, err := c.executeSingleRequest(req)

		// Record metrics
		if c.metrics != nil {
			duration := time.Since(req.startTime)
			if err != nil {
				c.metrics.RecordError(req, err, duration)
			} else {
				c.metrics.RecordResponse(req, resp, duration)
			}
		}

		// Check if we should retry
		if attempt < maxAttempts && c.shouldRetry(req, resp, err) {
			delay := c.calculateRetryDelay(req, attempt)

			// Call retry callback
			if req.RetryPolicy != nil && req.RetryPolicy.OnRetry != nil {
				req.RetryPolicy.OnRetry(attempt, err, delay)
			}

			// Wait before retry
			select {
			case <-time.After(delay):
			case <-req.Context.Done():
				return nil, req.Context.Err()
			}

			lastErr = err
			lastResp = resp
			continue
		}

		// Success or final attempt
		if err != nil {
			return resp, err
		}

		// Cache successful response
		if c.cache != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			cacheKey := c.getCacheKey(req)
			ttl := req.CacheTTL
			if ttl == 0 && c.config.Cache != nil {
				ttl = c.config.Cache.TTL
			}
			if ttl > 0 {
				c.cache.Set(cacheKey, resp, ttl)
			}
		}

		return resp, nil
	}

	// All attempts failed
	if lastErr != nil {
		return lastResp, lastErr
	}

	return lastResp, fmt.Errorf("all retry attempts failed")
}

// Configuration methods
func (c *httpClient) SetBaseURL(url string) Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.BaseURL = url
	return c
}

func (c *httpClient) SetUserAgent(userAgent string) Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.UserAgent = userAgent
	return c
}

func (c *httpClient) SetTimeout(timeout time.Duration) Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.config.Timeout == nil {
		c.config.Timeout = &TimeoutConfig{}
	}
	c.config.Timeout.Request = timeout
	c.httpClient.Timeout = timeout
	return c
}

func (c *httpClient) SetHeaders(headers map[string]string) Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.config.Headers == nil {
		c.config.Headers = make(map[string]string)
	}
	maps.Copy(c.config.Headers, headers)
	return c
}

func (c *httpClient) SetAuth(auth *AuthConfig) Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config.Auth = auth
	return c
}

// Middleware
func (c *httpClient) Use(middleware Middleware) Client {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.middlewares = append(c.middlewares, middleware)
	return c
}

// Clone tạo bản sao của client
func (c *httpClient) Clone() Client {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Deep copy config
	newConfig := *c.config

	// Create new client
	newClient := NewClient(&newConfig)

	// Copy middlewares
	for _, middleware := range c.middlewares {
		newClient.Use(middleware)
	}

	return newClient
}

// Close đóng client và giải phóng resources
func (c *httpClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close HTTP client transport
	if transport, ok := c.httpClient.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}

	// Close components
	if c.cache != nil {
		c.cache.Clear()
	}

	return nil
}

// Helper methods

// isAbsoluteURL kiểm tra URL có phải absolute không
func isAbsoluteURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.IsAbs()
}

// getCacheKey tạo cache key cho request
func (c *httpClient) getCacheKey(req *Request) string {
	if req.CacheKey != "" {
		return req.CacheKey
	}

	if c.config.Cache != nil && c.config.Cache.CacheKey != nil {
		return c.config.Cache.CacheKey(req)
	}

	// Default cache key
	return fmt.Sprintf("%s:%s", req.Method, req.URL)
}

// getRateLimitKey tạo rate limit key cho request
func (c *httpClient) getRateLimitKey(req *Request) string {
	if c.config.RateLimit != nil && c.config.RateLimit.KeyFunc != nil {
		return c.config.RateLimit.KeyFunc(req)
	}

	// Default rate limit key
	u, err := url.Parse(req.URL)
	if err != nil {
		return req.URL
	}
	return u.Host
}

// shouldRetry kiểm tra có nên retry không
func (c *httpClient) shouldRetry(req *Request, resp *Response, err error) bool {
	if req.RetryPolicy == nil {
		return false
	}

	// Check retryable errors
	if err != nil {
		for _, retryableErr := range req.RetryPolicy.RetryableErrors {
			if contains(err.Error(), retryableErr) {
				return true
			}
		}
		return false
	}

	// Check retryable status codes
	if resp != nil {
		return slices.Contains(req.RetryPolicy.RetryableStatus, resp.StatusCode)
	}

	return false
}

// calculateRetryDelay tính toán delay cho retry
func (c *httpClient) calculateRetryDelay(req *Request, attempt int) time.Duration {
	if req.RetryPolicy == nil {
		return DefaultRetryDelay
	}

	delay := req.RetryPolicy.InitialDelay
	for i := 1; i < attempt; i++ {
		delay = time.Duration(float64(delay) * req.RetryPolicy.BackoffFactor)
		if delay > req.RetryPolicy.MaxDelay {
			delay = req.RetryPolicy.MaxDelay
			break
		}
	}

	// Add jitter if enabled
	if req.RetryPolicy.Jitter {
		jitter := time.Duration(float64(delay) * 0.1 * (2*rand.Float64() - 1))
		delay += jitter
	}

	return delay
}

// contains kiểm tra string có chứa substring không
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				strings.Contains(s, substr))))
}

// executeSingleRequest thực hiện một request đơn lẻ
func (c *httpClient) executeSingleRequest(req *Request) (*Response, error) {
	// Start tracing
	var span Span
	if c.tracer != nil {
		span = c.tracer.StartSpan(req)
		defer span.Finish()
	}

	// Build HTTP request
	httpReq, err := c.buildHTTPRequest(req)
	if err != nil {
		if span != nil {
			span.SetError(err)
		}
		return nil, err
	}

	// Inject tracing headers
	if span != nil && c.tracer != nil {
		c.tracer.InjectHeaders(span, req.Headers)
	}

	// Record request metrics
	if c.metrics != nil {
		c.metrics.RecordRequest(req)
	}

	// Execute HTTP request
	startTime := time.Now()
	httpResp, err := c.httpClient.Do(httpReq)
	duration := time.Since(startTime)

	if err != nil {
		if span != nil {
			span.SetError(err)
		}
		return nil, c.wrapError(err)
	}
	defer httpResp.Body.Close()

	// Build response
	resp, err := c.buildResponse(req, httpResp, duration)
	if err != nil {
		if span != nil {
			span.SetError(err)
		}
		return nil, err
	}

	// Set span tags
	if span != nil {
		span.SetTag("http.status_code", resp.StatusCode)
		span.SetTag("http.method", string(req.Method))
		span.SetTag("http.url", req.URL)
	}

	return resp, nil
}

// wrapError wrap error với HTTPError
func (c *httpClient) wrapError(err error) error {
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr
	}

	// Check for specific error types
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			return ErrTimeout
		}
	}

	// Check error message for common patterns
	errMsg := err.Error()
	switch {
	case contains(errMsg, "connection refused"):
		return ErrConnectionRefused
	case contains(errMsg, "no such host"):
		return ErrDNSLookup
	case contains(errMsg, "tls"):
		return ErrTLSHandshake
	default:
		return &HTTPError{
			Code:    1000,
			Message: err.Error(),
			Type:    "unknown",
		}
	}
}

// buildHTTPRequest xây dựng HTTP request từ Request
func (c *httpClient) buildHTTPRequest(req *Request) (*http.Request, error) {
	var body io.Reader

	// Handle different body types
	if req.BodyReader != nil {
		body = req.BodyReader
	} else if req.Body != nil {
		bodyBytes, err := c.serializeBody(req)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(req.Context, string(req.Method), req.URL, body)
	if err != nil {
		return nil, &HTTPError{
			Code:    1300,
			Message: fmt.Sprintf("failed to create HTTP request: %v", err),
			Type:    "request",
		}
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Apply authentication
	if req.Auth != nil {
		if err := c.applyAuth(httpReq, req.Auth); err != nil {
			return nil, err
		}
	}

	// Set timeout
	if req.Timeout > 0 {
		ctx, cancel := context.WithTimeout(req.Context, req.Timeout)
		_ = cancel // Will be called when request completes
		httpReq = httpReq.WithContext(ctx)
	}

	return httpReq, nil
}

// serializeBody serializes request body based on content type
func (c *httpClient) serializeBody(req *Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}

	switch req.ContentType {
	case ContentTypeJSON:
		return json.Marshal(req.Body)

	case ContentTypeXML:
		return xml.Marshal(req.Body)

	case ContentTypeForm:
		if data, ok := req.Body.(map[string]string); ok {
			values := url.Values{}
			for key, value := range data {
				values.Set(key, value)
			}
			return []byte(values.Encode()), nil
		}

	case ContentTypeText:
		if str, ok := req.Body.(string); ok {
			return []byte(str), nil
		}
	}

	// Try to convert to JSON as fallback
	if req.ContentType == "" {
		req.ContentType = ContentTypeJSON
		req.Headers["Content-Type"] = string(ContentTypeJSON)
		return json.Marshal(req.Body)
	}

	return nil, &HTTPError{
		Code:    1301,
		Message: fmt.Sprintf("unsupported content type: %s", req.ContentType),
		Type:    "content_type",
	}
}

// applyAuth applies authentication to HTTP request
func (c *httpClient) applyAuth(httpReq *http.Request, auth *AuthConfig) error {
	switch auth.Type {
	case AuthTypeBasic:
		httpReq.SetBasicAuth(auth.Username, auth.Password)

	case AuthTypeBearer:
		httpReq.Header.Set("Authorization", "Bearer "+auth.Token)

	case AuthTypeAPIKey:
		if auth.Header != "" {
			httpReq.Header.Set(auth.Header, auth.APIKey)
		} else if auth.Query != "" {
			q := httpReq.URL.Query()
			q.Set(auth.Query, auth.APIKey)
			httpReq.URL.RawQuery = q.Encode()
		} else {
			httpReq.Header.Set("X-API-Key", auth.APIKey)
		}

	case AuthTypeCustom:
		for key, value := range auth.Custom {
			httpReq.Header.Set(key, value)
		}

	case AuthTypeOAuth2:
		// Check if token is expired and refresh if needed
		if auth.ExpiresAt.Before(time.Now()) && auth.RefreshFunc != nil {
			auth.Token = auth.RefreshFunc()
		}
		httpReq.Header.Set("Authorization", "Bearer "+auth.Token)
	}

	return nil
}
