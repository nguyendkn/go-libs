# HTTP Client Library for Go

🚀 **Thư viện HTTP Client mạnh mẽ và linh hoạt cho Go** - Được thiết kế với kiến trúc clean, tối ưu hóa performance và dễ dàng maintain.

## ✨ Tính năng

### 🔗 Core Features
- **Fluent API**: Request builder với method chaining
- **Connection Pooling**: Tối ưu hóa kết nối với connection reuse
- **Retry Logic**: Intelligent retry với exponential backoff
- **Timeout Management**: Flexible timeout configuration
- **Authentication**: Multiple auth methods (Basic, Bearer, OAuth2, API Key)

### 🚀 Advanced Features
- **Middleware System**: Extensible request/response processing pipeline
- **Caching**: HTTP caching với TTL và storage backends
- **Circuit Breaker**: Fault tolerance pattern
- **Rate Limiting**: Token bucket và sliding window algorithms
- **Metrics & Monitoring**: Real-time statistics và health checks
- **Distributed Tracing**: OpenTelemetry integration

### 🎯 Performance & Reliability
- **Smart Connection Management**: Automatic connection pooling
- **Request/Response Compression**: Gzip support
- **Error Handling**: Comprehensive error types và recovery
- **Context Support**: Full context.Context integration
- **Concurrent Safe**: Thread-safe operations

## 📦 Cài đặt

```bash
go get github.com/nguyendkn/go-libs/httpclient
```

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/nguyendkn/go-libs/httpclient"
)

func main() {
    // Tạo client
    client := httpclient.NewClient(&httpclient.ClientConfig{
        BaseURL: "https://api.example.com",
        Timeout: &httpclient.TimeoutConfig{
            Request: 30 * time.Second,
        },
    })
    defer client.Close()
    
    // Simple GET request
    resp, err := client.Get("/users/1").Send()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Status: %d\n", resp.StatusCode)
    fmt.Printf("Body: %s\n", resp.String())
}
```

### Fluent API

```go
// POST request với JSON body
user := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
}

resp, err := client.Post("/users").
    JSON(user).
    Header("X-API-Key", "secret").
    Timeout(10 * time.Second).
    Send()
```

### Authentication

```go
// Basic Auth
resp, err := client.Get("/protected").
    BasicAuth("username", "password").
    Send()

// Bearer Token
resp, err := client.Get("/api/data").
    BearerToken("your-jwt-token").
    Send()

// API Key
resp, err := client.Get("/api/data").
    APIKey("X-API-Key", "your-api-key").
    Send()
```

### Query Parameters

```go
// Individual parameters
resp, err := client.Get("/search").
    Query("q", "golang").
    Query("limit", "10").
    Send()

// Multiple parameters
params := map[string]string{
    "category": "tech",
    "sort":     "date",
}
resp, err := client.Get("/articles").
    QueryParams(params).
    Send()

// Struct to query params
type SearchParams struct {
    Query string `query:"q"`
    Limit int    `query:"limit"`
}

params := SearchParams{Query: "golang", Limit: 10}
resp, err := client.Get("/search").
    QueryStruct(params).
    Send()
```

## 🏗️ Advanced Usage

### Client Configuration

```go
config := &httpclient.ClientConfig{
    BaseURL:   "https://api.example.com",
    UserAgent: "MyApp/1.0",
    Headers: map[string]string{
        "Accept": "application/json",
    },
    Timeout: &httpclient.TimeoutConfig{
        Request:   30 * time.Second,
        Connect:   10 * time.Second,
        KeepAlive: 30 * time.Second,
    },
    ConnectionPool: &httpclient.ConnectionPoolConfig{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
    Retry: &httpclient.RetryPolicy{
        MaxAttempts:     3,
        InitialDelay:    1 * time.Second,
        MaxDelay:        30 * time.Second,
        BackoffFactor:   2.0,
        RetryableStatus: []int{429, 500, 502, 503, 504},
        Jitter:          true,
    },
}

client := httpclient.NewClient(config)
```

### Middleware

```go
// Logging middleware
logger := httpclient.NewLogger(&httpclient.LoggingConfig{
    Enabled: true,
    Level:   "info",
})
client.Use(httpclient.NewLoggingMiddleware(logger, nil))

// Metrics middleware
metrics := httpclient.NewMetrics(&httpclient.MetricsConfig{
    Enabled: true,
})
client.Use(httpclient.NewMetricsMiddleware(metrics, nil))

// Custom middleware
client.Use(httpclient.MiddlewareFunc(func(req *httpclient.Request, next httpclient.Handler) (*httpclient.Response, error) {
    // Pre-request processing
    req.Headers["X-Request-ID"] = generateRequestID()
    
    // Execute request
    resp, err := next(req)
    
    // Post-response processing
    if resp != nil {
        log.Printf("Request completed: %d", resp.StatusCode)
    }
    
    return resp, err
}))
```

### Caching

```go
config := &httpclient.ClientConfig{
    Cache: &httpclient.CacheConfig{
        Enabled:    true,
        TTL:        5 * time.Minute,
        MaxEntries: 1000,
    },
}

client := httpclient.NewClient(config)

// Request với cache
resp, err := client.Get("/data").
    Cache(10 * time.Minute).
    Send()
```

### Error Handling

```go
resp, err := client.Get("/api/data").Send()
if err != nil {
    if httpErr, ok := err.(*httpclient.HTTPError); ok {
        switch httpErr.StatusCode {
        case 401:
            log.Println("Unauthorized")
        case 404:
            log.Println("Not found")
        case 500:
            log.Println("Server error")
        default:
            log.Printf("HTTP error: %d - %s", httpErr.StatusCode, httpErr.Message)
        }
    } else {
        log.Printf("Request error: %v", err)
    }
    return
}

// Response validation
var data MyStruct
resp, err := client.Get("/api/data").ExpectJSON(&data)
if err != nil {
    log.Printf("Failed to get valid JSON: %v", err)
    return
}
```

### Context & Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.Get("/slow-endpoint").
    Context(ctx).
    Send()

if err != nil {
    if err == context.DeadlineExceeded {
        log.Println("Request timed out")
    }
}
```

## 📊 Monitoring & Metrics

### Built-in Metrics

```go
// Get client metrics
stats := client.GetStats()
fmt.Printf("Total requests: %d\n", stats.TotalRequests)
fmt.Printf("Error rate: %.2f%%\n", stats.ErrorRate*100)
fmt.Printf("Average latency: %v\n", stats.AverageLatency)
```

### Health Checks

```go
// Add health check
healthChecker := &MyHealthChecker{}
client.AddHealthCheck(healthChecker)

// Check health
status := client.CheckHealth(context.Background())
fmt.Printf("Health: %s\n", status.Status)
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestClient ./...

# Benchmark tests
go test -bench=. ./...
```

## 📁 Examples

Thư mục `examples/` chứa các ví dụ chi tiết:

- **basic_usage**: Các tính năng cơ bản
- **advanced_features**: Tính năng nâng cao với middleware

```bash
cd examples/basic_usage
go run main.go

cd examples/advanced_features
go run main.go
```

## 🏗️ Architecture

### Core Components

```
httpclient/
├── types.go              # Core types và constants
├── interfaces.go         # Interface definitions
├── client.go             # HTTP client implementation
├── request_builder.go    # Fluent API builder
├── response.go           # Response processing
├── middleware.go         # Middleware system
├── components.go         # Cache, metrics, etc.
└── examples/            # Usage examples
```

### Design Principles

- **Clean Architecture**: Clear separation of concerns
- **Interface-based Design**: Easy testing và mocking
- **Performance First**: Optimized cho high-throughput
- **Extensibility**: Plugin-based architecture
- **Developer Experience**: Intuitive API design

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by popular HTTP clients in other languages
- Built with Go's excellent standard library
- Community feedback và contributions

---

**Made with ❤️ by [nguyendkn](https://github.com/nguyendkn)**
