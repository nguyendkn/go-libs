package websocket

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GinWebSocketHandler tạo Gin handler cho WebSocket
func GinWebSocketHandler(server Server) gin.HandlerFunc {
	return gin.WrapH(server.GetHTTPHandler())
}

// GinWebSocketUpgrade tạo Gin handler để upgrade connection
func GinWebSocketUpgrade(server Server) gin.HandlerFunc {
	upgrader := server.GetUpgrader()
	options := server.GetOptions()
	
	return func(c *gin.Context) {
		// Authentication check
		var authInfo *AuthInfo
		if options.AuthRequired && options.AuthHandler != nil {
			auth, err := options.AuthHandler(c.Request)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
				return
			}
			authInfo = auth
		}
		
		// Upgrade connection
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade connection"})
			return
		}
		
		// Create server client
		client := newServerClient(conn, c.Request, authInfo, options)
		
		// Register client with server's hub
		if wsServer, ok := server.(*wsServer); ok {
			wsServer.hub.RegisterClient(client)
			
			// Update metrics
			wsServer.metricsMu.Lock()
			wsServer.metrics.TotalConnections++
			wsServer.metrics.ActiveConnections++
			wsServer.metricsMu.Unlock()
			
			// Call connect handler
			wsServer.handlersMu.RLock()
			if wsServer.onConnect != nil {
				go wsServer.onConnect(client)
			}
			wsServer.handlersMu.RUnlock()
			
			// Start client goroutines
			go wsServer.handleClient(client)
		}
	}
}

// GinWebSocketMiddleware tạo middleware cho WebSocket
func GinWebSocketMiddleware(server Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if this is a WebSocket upgrade request
		if c.GetHeader("Upgrade") == "websocket" {
			GinWebSocketUpgrade(server)(c)
			return
		}
		
		c.Next()
	}
}

// GinCORSMiddleware tạo CORS middleware cho WebSocket
func GinCORSMiddleware(allowedOrigins []string, allowedHeaders []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Check allowed origins
		if len(allowedOrigins) > 0 {
			allowed := false
			for _, allowedOrigin := range allowedOrigins {
				if origin == allowedOrigin || allowedOrigin == "*" {
					allowed = true
					break
				}
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Origin not allowed"})
				return
			}
		}
		
		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		
		if len(allowedHeaders) > 0 {
			headers := ""
			for i, header := range allowedHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Header("Access-Control-Allow-Headers", headers)
		} else {
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// GinRateLimitMiddleware tạo rate limiting middleware
func GinRateLimitMiddleware(rateLimiter RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := getClientIDFromContext(c)
		
		if !rateLimiter.Allow(clientID) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}
		
		c.Next()
	}
}

// GinAuthMiddleware tạo authentication middleware
func GinAuthMiddleware(authenticator Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		authInfo, err := authenticator.Authenticate(c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
			return
		}
		
		// Store auth info in context
		c.Set("auth", authInfo)
		c.Set("user_id", authInfo.UserID)
		c.Set("username", authInfo.Username)
		c.Set("roles", authInfo.Roles)
		
		c.Next()
	}
}

// GinJWTMiddleware tạo JWT authentication middleware
func GinJWTMiddleware(secret, issuer string, expiration time.Duration) gin.HandlerFunc {
	authenticator := NewJWTAuthenticator(secret, issuer, expiration)
	return GinAuthMiddleware(authenticator)
}

// GinBasicAuthMiddleware tạo basic authentication middleware
func GinBasicAuthMiddleware(users map[string]string) gin.HandlerFunc {
	authenticator := NewBasicAuthenticator(users)
	return GinAuthMiddleware(authenticator)
}

// GinMetricsHandler tạo handler cho metrics endpoint
func GinMetricsHandler(server Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics := server.GetMetrics()
		c.JSON(http.StatusOK, metrics)
	}
}

// GinHealthHandler tạo handler cho health check endpoint
func GinHealthHandler(server Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		health := server.GetHealth()
		
		var statusCode int
		switch health.Status {
		case "healthy":
			statusCode = http.StatusOK
		case "degraded":
			statusCode = http.StatusOK
		case "unhealthy":
			statusCode = http.StatusServiceUnavailable
		default:
			statusCode = http.StatusInternalServerError
		}
		
		c.JSON(statusCode, health)
	}
}

// GinWebSocketRoutes tạo các routes cho WebSocket
func GinWebSocketRoutes(router *gin.Engine, server Server, basePath string) {
	if basePath == "" {
		basePath = "/ws"
	}
	
	wsGroup := router.Group(basePath)
	
	// WebSocket upgrade endpoint
	wsGroup.GET("/", GinWebSocketUpgrade(server))
	
	// Metrics endpoint (if enabled)
	if server.GetOptions().EnableMetrics {
		metricsPath := server.GetOptions().MetricsPath
		if metricsPath == "" {
			metricsPath = "/metrics"
		}
		wsGroup.GET(metricsPath, GinMetricsHandler(server))
	}
	
	// Health check endpoint
	wsGroup.GET("/health", GinHealthHandler(server))
}

// GinWebSocketServer tạo một WebSocket server với Gin integration
type GinWebSocketServer struct {
	Server
	router *gin.Engine
}

// NewGinWebSocketServer tạo một Gin WebSocket server mới
func NewGinWebSocketServer(addr string, options *ServerOptions) *GinWebSocketServer {
	server := NewServer(addr, options)
	router := gin.New()
	
	// Add default middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// Add CORS middleware if origins are specified
	if options != nil && len(options.AllowedOrigins) > 0 {
		router.Use(GinCORSMiddleware(options.AllowedOrigins, options.AllowedHeaders))
	}
	
	// Setup WebSocket routes
	GinWebSocketRoutes(router, server, options.Path)
	
	return &GinWebSocketServer{
		Server: server,
		router: router,
	}
}

// GetRouter trả về Gin router
func (s *GinWebSocketServer) GetRouter() *gin.Engine {
	return s.router
}

// AddMiddleware thêm middleware vào router
func (s *GinWebSocketServer) AddMiddleware(middleware ...gin.HandlerFunc) {
	s.router.Use(middleware...)
}

// AddRoutes thêm custom routes
func (s *GinWebSocketServer) AddRoutes(setupFunc func(*gin.Engine)) {
	setupFunc(s.router)
}

// StartGin khởi động server với Gin
func (s *GinWebSocketServer) StartGin() error {
	// Start the WebSocket server's hub
	if err := s.Server.(*wsServer).hub.Start(); err != nil {
		return err
	}
	
	// Start Gin server
	return s.router.Run(s.Server.GetOptions().Addr)
}

// Helper functions

// getClientIDFromContext lấy client ID từ Gin context
func getClientIDFromContext(c *gin.Context) string {
	// Try to get from auth info first
	if authInfo, exists := c.Get("auth"); exists {
		if auth, ok := authInfo.(*AuthInfo); ok && auth.UserID != "" {
			return auth.UserID
		}
	}
	
	// Fall back to IP address
	return c.ClientIP()
}

// GetAuthFromContext lấy auth info từ Gin context
func GetAuthFromContext(c *gin.Context) (*AuthInfo, bool) {
	if authInfo, exists := c.Get("auth"); exists {
		if auth, ok := authInfo.(*AuthInfo); ok {
			return auth, true
		}
	}
	return nil, false
}

// GetUserIDFromContext lấy user ID từ Gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	if userID, exists := c.Get("user_id"); exists {
		if id, ok := userID.(string); ok {
			return id, true
		}
	}
	return "", false
}

// GetUsernameFromContext lấy username từ Gin context
func GetUsernameFromContext(c *gin.Context) (string, bool) {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name, true
		}
	}
	return "", false
}

// GetRolesFromContext lấy roles từ Gin context
func GetRolesFromContext(c *gin.Context) ([]string, bool) {
	if roles, exists := c.Get("roles"); exists {
		if roleList, ok := roles.([]string); ok {
			return roleList, true
		}
	}
	return nil, false
}

// HasRole kiểm tra xem user có role cụ thể không
func HasRole(c *gin.Context, role string) bool {
	roles, exists := GetRolesFromContext(c)
	if !exists {
		return false
	}
	
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// RequireRole tạo middleware yêu cầu role cụ thể
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !HasRole(c, role) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		c.Next()
	}
}

// RequireAnyRole tạo middleware yêu cầu một trong các roles
func RequireAnyRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hasRole := false
		for _, role := range roles {
			if HasRole(c, role) {
				hasRole = true
				break
			}
		}
		
		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			return
		}
		c.Next()
	}
}
