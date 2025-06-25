# WebSocket Library

Thư viện WebSocket linh hoạt và đầy đủ tính năng cho Go, hỗ trợ cả client và server với nhiều tùy chọn cấu hình.

## 🚀 Tính năng

### 📡 Client
- ✅ Kết nối WebSocket với auto-reconnect
- 💓 Ping/Pong heartbeat tự động
- 📦 Message queuing và buffering
- 🗜️ Compression support
- 🔐 Custom headers và authentication
- ⚡ Rate limiting
- ⏱️ Timeout configuration linh hoạt
- 🔄 Exponential backoff cho reconnection

### 🖥️ Server
- 🏗️ WebSocket server với hub pattern
- 🏠 Room/Channel management
- 📢 Broadcast messaging
- 👥 Client management
- 🔒 Authentication middleware
- ⚡ Rate limiting per client
- 🏥 Health monitoring
- 📊 Real-time metrics

### 🔐 Security
- 🎫 JWT authentication
- 🛡️ Rate limiting với multiple algorithms
- 🌐 Origin validation
- 🔧 Custom authentication handlers
- 🚫 Blacklist/Whitelist support

### 📊 Monitoring
- 📈 Connection metrics
- 📋 Message statistics
- 🏥 Health checks
- 📝 Structured logging
- 📊 Prometheus-compatible metrics

### 🌐 Gin Integration
- 🔌 Seamless Gin framework integration
- 🛠️ Ready-to-use middleware
- 🎨 CORS support
- 📡 RESTful API endpoints

## 📦 Cài đặt

```bash
go get github.com/go-libs/websocket
```

## 🚀 Sử dụng nhanh

### 📡 Client đơn giản

```go
package main

import (
    "fmt"
    "github.com/go-libs/websocket"
)

func main() {
    // Tạo client với cấu hình mặc định
    client := websocket.QuickClient("ws://localhost:8080/ws")

    // Xử lý events
    client.OnConnect(func() {
        fmt.Println("✅ Connected!")
    })

    client.OnMessage(func(msg []byte) {
        fmt.Printf("📨 Received: %s\n", msg)
    })

    client.OnDisconnect(func(err error) {
        fmt.Printf("❌ Disconnected: %v\n", err)
    })

    // Kết nối và gửi message
    client.Connect()
    client.SendText("Hello WebSocket!")

    // Giữ chương trình chạy
    select {}
}
```

### 🖥️ Server đơn giản

```go
package main

import (
    "fmt"
    "github.com/go-libs/websocket"
)

func main() {
    // Tạo server với cấu hình mặc định
    server := websocket.QuickServer(":8080")

    // Xử lý client connections
    server.OnConnect(func(client websocket.ServerClient) {
        fmt.Printf("🔗 Client connected: %s\n", client.ID())
        client.SendText("Welcome!")
    })

    server.OnMessage(func(client websocket.ServerClient, msg []byte) {
        fmt.Printf("📨 Message from %s: %s\n", client.ID(), msg)
        // Echo back to sender
        client.Send(msg)
        // Broadcast to all other clients
        server.Broadcast(msg)
    })

    server.OnDisconnect(func(client websocket.ServerClient, err error) {
        fmt.Printf("❌ Client disconnected: %s\n", client.ID())
    })

    fmt.Println("🚀 Starting server on :8080")
    server.Start()
}
```

### 🌐 Gin Integration

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/go-libs/websocket"
)

func main() {
    // Tạo Gin WebSocket server
    server := websocket.QuickGinServer(":8080")

    // Lấy Gin router
    router := server.GetRouter()

    // Thêm custom routes
    router.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "WebSocket server is running!"})
    })

    // Setup WebSocket handlers
    wsServer := server.Server
    wsServer.OnConnect(func(client websocket.ServerClient) {
        client.SendJSON(map[string]string{
            "type": "welcome",
            "message": "Connected to Gin WebSocket server!",
        })
    })

    wsServer.OnTextMessage(func(client websocket.ServerClient, message string) {
        // Broadcast to all clients
        wsServer.BroadcastJSON(map[string]interface{}{
            "type": "message",
            "client": client.ID(),
            "data": message,
        })
    })

    // Start server
    server.StartGin()
}
```

## 🔧 Cấu hình nâng cao

### 📡 Client với Builder Pattern

```go
client := websocket.NewClientBuilder("ws://localhost:8080/ws").
    WithAuth("Bearer token123", "Authorization").
    WithReconnect(true, 5*time.Second, 10).
    WithCompression().
    WithRateLimit(100, 200).
    WithTimeouts(60*time.Second, 10*time.Second).
    WithHeaders(map[string]string{
        "User-Agent": "MyApp/1.0",
    }).
    Build()
```

### 🖥️ Server với Builder Pattern

```go
server := websocket.NewServerBuilder(":8080").
    WithPath("/websocket").
    WithCompression(6).
    WithJWTAuth("secret", "myapp", 24*time.Hour).
    WithRateLimit(100, 200).
    WithCORS([]string{"*"}, []string{"Authorization"}).
    WithTLS("cert.pem", "key.pem").
    WithMetrics(true, "/metrics").
    WithTimeouts(60*time.Second, 10*time.Second, 120*time.Second, 30*time.Second).
    Build()
```

## 📚 Examples

### 🏠 Room Management

```go
// Tạo room với options
roomOptions := &websocket.RoomOptions{
    MaxClients:     100,
    RequireAuth:    true,
    AllowedRoles:   []string{"user", "admin"},
    MessageHistory: 50,
    TTL:           24 * time.Hour,
}

server.CreateRoom("chat-room", roomOptions)

// Client join room
server.OnConnect(func(client websocket.ServerClient) {
    // Auto-join general room
    if room, exists := server.GetRoom("general"); exists {
        room.AddClient(client)
        client.JoinRoom("general")
    }
})

// Broadcast to specific room
server.BroadcastToRoomText("chat-room", "Hello room!")
```

### 🔐 Authentication

```go
// JWT Authentication
jwtAuth := websocket.NewJWTAuthenticator("secret", "myapp", 24*time.Hour)

server := websocket.NewServerBuilder(":8080").
    WithAuth(jwtAuth.Authenticate).
    Build()

// Custom Authentication
customAuth := func(req *http.Request) (*websocket.AuthInfo, error) {
    token := req.Header.Get("X-API-Key")
    if token == "valid-key" {
        return &websocket.AuthInfo{
            UserID:   "user123",
            Username: "john",
            Roles:    []string{"user"},
        }, nil
    }
    return nil, errors.New("invalid API key")
}

server.SetOptions(&websocket.ServerOptions{
    AuthRequired: true,
    AuthHandler:  customAuth,
})
```

### ⚡ Rate Limiting

```go
// Token bucket rate limiter
rateLimiter := websocket.NewTokenBucketRateLimiter(100, 200)

// Sliding window rate limiter
rateLimiter := websocket.NewSlidingWindowRateLimiter(time.Minute, 100)

// Apply to server
server := websocket.NewServerBuilder(":8080").
    WithRateLimit(100, 200).
    Build()

// Apply to Gin middleware
router.Use(websocket.GinRateLimitMiddleware(rateLimiter))
```

### 📊 Monitoring

```go
// Get server metrics
metrics := server.GetMetrics()
fmt.Printf("Active connections: %d\n", metrics.ActiveConnections)
fmt.Printf("Total messages: %d\n", metrics.TotalMessages)

// Get health status
health := server.GetHealth()
fmt.Printf("Status: %s\n", health.Status)

// Custom health checks
healthChecker := websocket.NewHealthChecker()
healthChecker.AddCheck("database", "Database connectivity", func() (bool, string) {
    // Check database connection
    return true, "Database is healthy"
})

// Start monitoring server
monitoringServer := websocket.NewMonitoringServer(":9090", metricsCollector, healthChecker)
go monitoringServer.Start()
```

### 🌐 Gin Middleware

```go
router := gin.Default()

// CORS middleware
router.Use(websocket.GinCORSMiddleware(
    []string{"http://localhost:3000"},
    []string{"Authorization", "Content-Type"},
))

// JWT middleware
router.Use(websocket.GinJWTMiddleware("secret", "myapp", 24*time.Hour))

// Rate limiting middleware
router.Use(websocket.GinRateLimitMiddleware(rateLimiter))

// Role-based access control
adminRoutes := router.Group("/admin")
adminRoutes.Use(websocket.RequireRole("admin"))

// WebSocket routes
websocket.GinWebSocketRoutes(router, server, "/ws")
```

## 🔧 API Reference

### Client Interface

```go
type Client interface {
    // Connection management
    Connect() error
    ConnectWithContext(ctx context.Context) error
    Disconnect() error
    IsConnected() bool
    GetState() ConnectionState

    // Message handling
    Send(data []byte) error
    SendText(text string) error
    SendJSON(v interface{}) error
    SendWithType(messageType MessageType, data []byte) error

    // Event handlers
    OnConnect(handler func())
    OnDisconnect(handler func(error))
    OnMessage(handler func([]byte))
    OnTextMessage(handler func(string))
    OnBinaryMessage(handler func([]byte))
    OnError(handler func(error))
    OnReconnect(handler func(int))

    // Client info
    ID() string
    Info() *ClientInfo
    Metrics() *ConnectionMetrics

    // Configuration
    SetOptions(options *ClientOptions)
    GetOptions() *ClientOptions

    // Lifecycle
    Close() error
}
```

### Server Interface

```go
type Server interface {
    // Server lifecycle
    Start() error
    StartWithContext(ctx context.Context) error
    Stop() error
    Shutdown(ctx context.Context) error

    // Client management
    GetClient(id string) (ServerClient, bool)
    GetClients() []ServerClient
    GetClientCount() int
    DisconnectClient(id string) error

    // Room management
    CreateRoom(name string, options *RoomOptions) error
    DeleteRoom(name string) error
    GetRoom(name string) (Room, bool)
    GetRooms() []Room

    // Broadcasting
    Broadcast(data []byte) error
    BroadcastText(text string) error
    BroadcastJSON(v interface{}) error
    BroadcastToRoom(room string, data []byte) error
    BroadcastToRoomText(room string, text string) error
    BroadcastToRoomJSON(room string, v interface{}) error

    // Event handlers
    OnConnect(handler func(ServerClient))
    OnDisconnect(handler func(ServerClient, error))
    OnMessage(handler func(ServerClient, []byte))
    OnTextMessage(handler func(ServerClient, string))
    OnBinaryMessage(handler func(ServerClient, []byte))
    OnError(handler func(ServerClient, error))
    OnRoomJoin(handler func(ServerClient, string))
    OnRoomLeave(handler func(ServerClient, string))

    // Configuration
    SetOptions(options *ServerOptions)
    GetOptions() *ServerOptions

    // Monitoring
    GetMetrics() *ServerMetrics
    GetHealth() *HealthStatus

    // HTTP handler
    GetHTTPHandler() http.Handler
    GetUpgrader() *websocket.Upgrader
}
```

## 🧪 Testing

Chạy tests:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...

# Run specific test
go test -run TestClientConnection ./...
```

## 📁 Project Structure

```
websocket/
├── README.md              # Documentation
├── go.mod                 # Go module
├── websocket.go           # Main package file
├── types.go               # Type definitions
├── interfaces.go          # Interface definitions
├── client.go              # WebSocket client
├── server.go              # WebSocket server
├── server_client.go       # Server-side client
├── hub.go                 # Message hub
├── room.go                # Room management
├── auth.go                # Authentication
├── rate_limiter.go        # Rate limiting
├── monitoring.go          # Monitoring & metrics
├── gin.go                 # Gin integration
├── *_test.go              # Unit tests
└── examples/              # Example applications
    ├── basic_client/      # Basic client example
    ├── basic_server/      # Basic server example
    └── gin_server/        # Gin integration example
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket implementation
- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Rate](https://golang.org/x/time/rate) - Rate limiting

## 📞 Support

- 📧 Email: support@go-libs.com
- 🐛 Issues: [GitHub Issues](https://github.com/go-libs/websocket/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/go-libs/websocket/discussions)

---

Made with ❤️ by the Go-Libs team
