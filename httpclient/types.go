package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPMethod định nghĩa các HTTP methods
type HTTPMethod string

const (
	MethodGET     HTTPMethod = "GET"
	MethodPOST    HTTPMethod = "POST"
	MethodPUT     HTTPMethod = "PUT"
	MethodPATCH   HTTPMethod = "PATCH"
	MethodDELETE  HTTPMethod = "DELETE"
	MethodHEAD    HTTPMethod = "HEAD"
	MethodOPTIONS HTTPMethod = "OPTIONS"
	MethodTRACE   HTTPMethod = "TRACE"
	MethodCONNECT HTTPMethod = "CONNECT"
)

// ContentType định nghĩa các content types phổ biến
type ContentType string

const (
	ContentTypeJSON       ContentType = "application/json"
	ContentTypeXML        ContentType = "application/xml"
	ContentTypeForm       ContentType = "application/x-www-form-urlencoded"
	ContentTypeMultipart  ContentType = "multipart/form-data"
	ContentTypeText       ContentType = "text/plain"
	ContentTypeHTML       ContentType = "text/html"
	ContentTypeBinary     ContentType = "application/octet-stream"
	ContentTypeJavaScript ContentType = "application/javascript"
	ContentTypeCSS        ContentType = "text/css"
)

// AuthType định nghĩa các loại authentication
type AuthType string

const (
	AuthTypeNone   AuthType = "none"
	AuthTypeBasic  AuthType = "basic"
	AuthTypeBearer AuthType = "bearer"
	AuthTypeAPIKey AuthType = "apikey"
	AuthTypeOAuth2 AuthType = "oauth2"
	AuthTypeCustom AuthType = "custom"
)

// RetryPolicy định nghĩa chính sách retry
type RetryPolicy struct {
	MaxAttempts     int           `json:"maxAttempts"`
	InitialDelay    time.Duration `json:"initialDelay"`
	MaxDelay        time.Duration `json:"maxDelay"`
	BackoffFactor   float64       `json:"backoffFactor"`
	RetryableStatus []int         `json:"retryableStatus"`
	RetryableErrors []string      `json:"retryableErrors"`
	Jitter          bool          `json:"jitter"`
	OnRetry         func(attempt int, err error, delay time.Duration)
}

// TimeoutConfig định nghĩa các timeout settings
type TimeoutConfig struct {
	Connect   time.Duration `json:"connect"`
	Request   time.Duration `json:"request"`
	Response  time.Duration `json:"response"`
	Idle      time.Duration `json:"idle"`
	KeepAlive time.Duration `json:"keepAlive"`
}

// ConnectionPoolConfig cấu hình connection pool
type ConnectionPoolConfig struct {
	MaxIdleConns        int           `json:"maxIdleConns"`
	MaxIdleConnsPerHost int           `json:"maxIdleConnsPerHost"`
	MaxConnsPerHost     int           `json:"maxConnsPerHost"`
	IdleConnTimeout     time.Duration `json:"idleConnTimeout"`
	DisableKeepAlives   bool          `json:"disableKeepAlives"`
	DisableCompression  bool          `json:"disableCompression"`
}

// TLSConfig cấu hình TLS
type TLSConfig struct {
	InsecureSkipVerify bool     `json:"insecureSkipVerify"`
	ServerName         string   `json:"serverName"`
	CertFile           string   `json:"certFile"`
	KeyFile            string   `json:"keyFile"`
	CAFile             string   `json:"caFile"`
	MinVersion         uint16   `json:"minVersion"`
	MaxVersion         uint16   `json:"maxVersion"`
	CipherSuites       []uint16 `json:"cipherSuites"`
}

// ProxyConfig cấu hình proxy
type ProxyConfig struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
	NoProxy  string `json:"noProxy"`
}

// CacheConfig cấu hình caching
type CacheConfig struct {
	Enabled       bool          `json:"enabled"`
	TTL           time.Duration `json:"ttl"`
	MaxSize       int64         `json:"maxSize"`
	MaxEntries    int           `json:"maxEntries"`
	Storage       string        `json:"storage"` // "memory", "redis", "file"
	StorageConfig any           `json:"storageConfig"`
	CacheKey      func(*Request) string
	ShouldCache   func(*Request, *Response) bool
	OnCacheHit    func(key string)
	OnCacheMiss   func(key string)
}

// CircuitBreakerConfig cấu hình circuit breaker
type CircuitBreakerConfig struct {
	Enabled          bool          `json:"enabled"`
	FailureThreshold int           `json:"failureThreshold"`
	SuccessThreshold int           `json:"successThreshold"`
	Timeout          time.Duration `json:"timeout"`
	MaxRequests      int           `json:"maxRequests"`
	Interval         time.Duration `json:"interval"`
	OnStateChange    func(from, to string)
	IsFailure        func(error) bool
}

// RateLimitConfig cấu hình rate limiting
type RateLimitConfig struct {
	Enabled   bool          `json:"enabled"`
	Rate      float64       `json:"rate"`      // requests per second
	Burst     int           `json:"burst"`     // burst size
	Window    time.Duration `json:"window"`    // sliding window
	Algorithm string        `json:"algorithm"` // "token_bucket", "sliding_window"
	KeyFunc   func(*Request) string
	OnLimit   func(key string, delay time.Duration)
}

// MetricsConfig cấu hình metrics
type MetricsConfig struct {
	Enabled          bool      `json:"enabled"`
	Namespace        string    `json:"namespace"`
	Subsystem        string    `json:"subsystem"`
	Labels           []string  `json:"labels"`
	HistogramBuckets []float64 `json:"histogramBuckets"`
	CollectBody      bool      `json:"collectBody"`
	CollectHeaders   bool      `json:"collectHeaders"`
}

// TracingConfig cấu hình tracing
type TracingConfig struct {
	Enabled     bool              `json:"enabled"`
	ServiceName string            `json:"serviceName"`
	Endpoint    string            `json:"endpoint"`
	SampleRate  float64           `json:"sampleRate"`
	Headers     map[string]string `json:"headers"`
}

// LoggingConfig cấu hình logging
type LoggingConfig struct {
	Enabled          bool     `json:"enabled"`
	Level            string   `json:"level"`
	Format           string   `json:"format"` // "json", "text"
	Output           string   `json:"output"` // "stdout", "stderr", "file"
	File             string   `json:"file"`
	MaxSize          int      `json:"maxSize"`
	MaxBackups       int      `json:"maxBackups"`
	MaxAge           int      `json:"maxAge"`
	Compress         bool     `json:"compress"`
	Fields           []string `json:"fields"`
	SensitiveHeaders []string `json:"sensitiveHeaders"`
	CollectHeaders   bool     `json:"collectHeaders"`
	CollectBody      bool     `json:"collectBody"`
}

// Request đại diện cho HTTP request
type Request struct {
	Method      HTTPMethod        `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"queryParams"`
	Body        any               `json:"body"`
	BodyReader  io.Reader         `json:"-"`
	ContentType ContentType       `json:"contentType"`
	Auth        *AuthConfig       `json:"auth"`
	Timeout     time.Duration     `json:"timeout"`
	Context     context.Context   `json:"-"`
	Metadata    map[string]any    `json:"metadata"`

	// Request options
	FollowRedirects bool `json:"followRedirects"`
	MaxRedirects    int  `json:"maxRedirects"`

	// Retry options
	RetryPolicy *RetryPolicy `json:"retryPolicy"`

	// Cache options
	CacheKey string        `json:"cacheKey"`
	CacheTTL time.Duration `json:"cacheTTL"`
	NoCache  bool          `json:"noCache"`

	// Internal fields
	attempt   int
	startTime time.Time
}

// Response đại diện cho HTTP response
type Response struct {
	StatusCode    int                 `json:"statusCode"`
	Status        string              `json:"status"`
	Headers       map[string][]string `json:"headers"`
	Body          []byte              `json:"body"`
	BodyReader    io.ReadCloser       `json:"-"`
	ContentType   string              `json:"contentType"`
	ContentLength int64               `json:"contentLength"`
	Request       *Request            `json:"request"`
	Metadata      map[string]any      `json:"metadata"`

	// Timing information
	Duration     time.Duration `json:"duration"`
	DNSLookup    time.Duration `json:"dnsLookup"`
	TCPConnect   time.Duration `json:"tcpConnect"`
	TLSHandshake time.Duration `json:"tlsHandshake"`
	ServerTime   time.Duration `json:"serverTime"`
	ResponseTime time.Duration `json:"responseTime"`

	// Cache information
	FromCache bool      `json:"fromCache"`
	CacheKey  string    `json:"cacheKey"`
	CachedAt  time.Time `json:"cachedAt"`

	// Retry information
	Attempts      int           `json:"attempts"`
	TotalDuration time.Duration `json:"totalDuration"`

	// Raw HTTP response
	Raw *http.Response `json:"-"`
}

// AuthConfig cấu hình authentication
type AuthConfig struct {
	Type     AuthType          `json:"type"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Token    string            `json:"token"`
	APIKey   string            `json:"apiKey"`
	Header   string            `json:"header"`
	Query    string            `json:"query"`
	Scheme   string            `json:"scheme"`
	Custom   map[string]string `json:"custom"`

	// OAuth2 specific
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	TokenURL     string   `json:"tokenUrl"`
	Scopes       []string `json:"scopes"`

	// Token refresh
	RefreshToken string        `json:"refreshToken"`
	ExpiresAt    time.Time     `json:"expiresAt"`
	RefreshFunc  func() string `json:"-"`
}

// ClientConfig cấu hình tổng thể cho HTTP client
type ClientConfig struct {
	BaseURL        string                `json:"baseUrl"`
	UserAgent      string                `json:"userAgent"`
	Headers        map[string]string     `json:"headers"`
	Auth           *AuthConfig           `json:"auth"`
	Timeout        *TimeoutConfig        `json:"timeout"`
	Retry          *RetryPolicy          `json:"retry"`
	ConnectionPool *ConnectionPoolConfig `json:"connectionPool"`
	TLS            *TLSConfig            `json:"tls"`
	Proxy          *ProxyConfig          `json:"proxy"`
	Cache          *CacheConfig          `json:"cache"`
	CircuitBreaker *CircuitBreakerConfig `json:"circuitBreaker"`
	RateLimit      *RateLimitConfig      `json:"rateLimit"`
	Metrics        *MetricsConfig        `json:"metrics"`
	Tracing        *TracingConfig        `json:"tracing"`
	Logging        *LoggingConfig        `json:"logging"`

	// Behavior options
	FollowRedirects bool `json:"followRedirects"`
	MaxRedirects    int  `json:"maxRedirects"`

	// Debug options
	Debug        bool `json:"debug"`
	DumpRequest  bool `json:"dumpRequest"`
	DumpResponse bool `json:"dumpResponse"`
}

// Default values
const (
	DefaultTimeout         = 30 * time.Second
	DefaultConnectTimeout  = 10 * time.Second
	DefaultIdleTimeout     = 90 * time.Second
	DefaultKeepAlive       = 30 * time.Second
	DefaultMaxIdleConns    = 100
	DefaultMaxConnsPerHost = 10
	DefaultRetryAttempts   = 3
	DefaultRetryDelay      = 1 * time.Second
	DefaultMaxRetryDelay   = 30 * time.Second
	DefaultBackoffFactor   = 2.0
	DefaultUserAgent       = "go-httpclient/1.0"
)

// Default configurations
var (
	DefaultTimeoutConfig = &TimeoutConfig{
		Connect:   DefaultConnectTimeout,
		Request:   DefaultTimeout,
		Response:  DefaultTimeout,
		Idle:      DefaultIdleTimeout,
		KeepAlive: DefaultKeepAlive,
	}

	DefaultConnectionPoolConfig = &ConnectionPoolConfig{
		MaxIdleConns:        DefaultMaxIdleConns,
		MaxIdleConnsPerHost: DefaultMaxConnsPerHost,
		MaxConnsPerHost:     DefaultMaxConnsPerHost,
		IdleConnTimeout:     DefaultIdleTimeout,
		DisableKeepAlives:   false,
		DisableCompression:  false,
	}

	DefaultRetryPolicy = &RetryPolicy{
		MaxAttempts:     DefaultRetryAttempts,
		InitialDelay:    DefaultRetryDelay,
		MaxDelay:        DefaultMaxRetryDelay,
		BackoffFactor:   DefaultBackoffFactor,
		RetryableStatus: []int{408, 429, 500, 502, 503, 504},
		Jitter:          true,
	}
)

// Error types
type HTTPError struct {
	Code       int       `json:"code"`
	Message    string    `json:"message"`
	Type       string    `json:"type"`
	StatusCode int       `json:"statusCode"`
	Response   *Response `json:"response,omitempty"`
}

func (e *HTTPError) Error() string {
	if e.StatusCode > 0 {
		return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("HTTP Client Error [%d]: %s", e.Code, e.Message)
}

// Common errors
var (
	ErrTimeout           = &HTTPError{Code: 1001, Message: "request timeout", Type: "timeout"}
	ErrConnectionRefused = &HTTPError{Code: 1002, Message: "connection refused", Type: "connection"}
	ErrDNSLookup         = &HTTPError{Code: 1003, Message: "DNS lookup failed", Type: "dns"}
	ErrTLSHandshake      = &HTTPError{Code: 1004, Message: "TLS handshake failed", Type: "tls"}
	ErrTooManyRedirects  = &HTTPError{Code: 1005, Message: "too many redirects", Type: "redirect"}
	ErrInvalidURL        = &HTTPError{Code: 1006, Message: "invalid URL", Type: "url"}
	ErrInvalidResponse   = &HTTPError{Code: 1007, Message: "invalid response", Type: "response"}
	ErrRateLimited       = &HTTPError{Code: 1008, Message: "rate limited", Type: "ratelimit"}
	ErrCircuitOpen       = &HTTPError{Code: 1009, Message: "circuit breaker open", Type: "circuit"}
	ErrCacheMiss         = &HTTPError{Code: 1010, Message: "cache miss", Type: "cache"}
)
