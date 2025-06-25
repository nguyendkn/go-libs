package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// wsServer implements the Server interface
type wsServer struct {
	options *ServerOptions
	hub     Hub

	// HTTP server
	server   *http.Server
	upgrader *websocket.Upgrader

	// State
	running int32 // atomic

	// Event handlers
	onConnect       func(ServerClient)
	onDisconnect    func(ServerClient, error)
	onMessage       func(ServerClient, []byte)
	onTextMessage   func(ServerClient, string)
	onBinaryMessage func(ServerClient, []byte)
	onError         func(ServerClient, error)
	onRoomJoin      func(ServerClient, string)
	onRoomLeave     func(ServerClient, string)
	handlersMu      sync.RWMutex

	// Metrics
	metrics   *ServerMetrics
	metricsMu sync.RWMutex
	startTime time.Time

	// Context
	ctx    context.Context
	cancel context.CancelFunc
}

// NewServer tạo một WebSocket server mới
func NewServer(addr string, options ...*ServerOptions) Server {
	opts := &ServerOptions{
		Addr:              addr,
		Path:              "/ws",
		ReadTimeout:       DefaultReadTimeout,
		WriteTimeout:      DefaultWriteTimeout,
		IdleTimeout:       DefaultIdleTimeout,
		MaxMessageSize:    DefaultMaxMessageSize,
		MessageQueueSize:  DefaultMessageQueueSize,
		PingInterval:      DefaultPingInterval,
		PongTimeout:       DefaultPongTimeout,
		RateLimit:         DefaultRateLimit,
		RateBurst:         DefaultRateBurst,
		RateWindow:        DefaultRateWindow,
		EnableCompression: false,
		CompressionLevel:  1,
		ShutdownTimeout:   DefaultShutdownTimeout,
		EnableMetrics:     true,
		MetricsPath:       "/metrics",
	}

	if len(options) > 0 && options[0] != nil {
		mergeServerOptions(opts, options[0])
	}

	ctx, cancel := context.WithCancel(context.Background())

	server := &wsServer{
		options: opts,
		hub:     NewHub(),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: opts.EnableCompression,
			CheckOrigin:       opts.CheckOrigin,
			Subprotocols:      opts.Subprotocols,
		},
		metrics: &ServerMetrics{
			LastUpdated: time.Now(),
		},
		startTime: time.Now(),
		ctx:       ctx,
		cancel:    cancel,
	}

	// Set default CheckOrigin if not provided
	if server.upgrader.CheckOrigin == nil {
		server.upgrader.CheckOrigin = func(r *http.Request) bool {
			if len(opts.AllowedOrigins) == 0 {
				return true // allow all origins by default
			}
			origin := r.Header.Get("Origin")
			return slices.Contains(opts.AllowedOrigins, origin)
		}
	}

	return server
}

// Start khởi động server
func (s *wsServer) Start() error {
	return s.StartWithContext(context.Background())
}

// StartWithContext khởi động server với context
func (s *wsServer) StartWithContext(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&s.running, 0, 1) {
		return fmt.Errorf("server is already running")
	}

	// Start hub
	if err := s.hub.Start(); err != nil {
		atomic.StoreInt32(&s.running, 0)
		return fmt.Errorf("failed to start hub: %w", err)
	}

	// Setup HTTP routes
	mux := http.NewServeMux()
	mux.HandleFunc(s.options.Path, s.handleWebSocket)

	if s.options.EnableMetrics && s.options.MetricsPath != "" {
		mux.HandleFunc(s.options.MetricsPath, s.handleMetrics)
	}

	// Create HTTP server
	s.server = &http.Server{
		Addr:         s.options.Addr,
		Handler:      mux,
		ReadTimeout:  s.options.ReadTimeout,
		WriteTimeout: s.options.WriteTimeout,
		IdleTimeout:  s.options.IdleTimeout,
	}

	// Start metrics updater
	go s.updateMetrics()

	// Start server
	if s.options.TLSCertFile != "" && s.options.TLSKeyFile != "" {
		return s.server.ListenAndServeTLS(s.options.TLSCertFile, s.options.TLSKeyFile)
	}

	return s.server.ListenAndServe()
}

// Stop dừng server
func (s *wsServer) Stop() error {
	return s.Shutdown(context.Background())
}

// Shutdown dừng server gracefully
func (s *wsServer) Shutdown(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&s.running, 1, 0) {
		return nil // already stopped
	}

	s.cancel()

	// Create shutdown context with timeout
	shutdownCtx, cancel := context.WithTimeout(ctx, s.options.ShutdownTimeout)
	defer cancel()

	// Stop hub
	if err := s.hub.Stop(); err != nil {
		return fmt.Errorf("failed to stop hub: %w", err)
	}

	// Shutdown HTTP server
	if s.server != nil {
		return s.server.Shutdown(shutdownCtx)
	}

	return nil
}

// handleWebSocket xử lý WebSocket upgrade
func (s *wsServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Authentication check
	var authInfo *AuthInfo
	if s.options.AuthRequired && s.options.AuthHandler != nil {
		auth, err := s.options.AuthHandler(r)
		if err != nil {
			http.Error(w, "Authentication failed", http.StatusUnauthorized)
			return
		}
		authInfo = auth
	}

	// Upgrade connection
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// Create server client
	client := newServerClient(conn, r, authInfo, s.options, s)

	// Register client with hub
	s.hub.RegisterClient(client)

	// Update metrics
	atomic.AddInt64(&s.metrics.TotalConnections, 1)
	atomic.AddInt64(&s.metrics.ActiveConnections, 1)

	// Call connect handler
	s.handlersMu.RLock()
	if s.onConnect != nil {
		go s.onConnect(client)
	}
	s.handlersMu.RUnlock()

	// Start client goroutines
	go s.handleClient(client)
}

// handleClient xử lý client connection
func (s *wsServer) handleClient(client ServerClient) {
	defer func() {
		// Unregister client
		s.hub.UnregisterClient(client.ID())

		// Update metrics
		atomic.AddInt64(&s.metrics.ActiveConnections, -1)

		// Call disconnect handler
		s.handlersMu.RLock()
		if s.onDisconnect != nil {
			go s.onDisconnect(client, nil)
		}
		s.handlersMu.RUnlock()
	}()

	// Cast to internal client type
	wsClient, ok := client.(*serverClient)
	if !ok {
		return
	}

	// Start client pumps
	go wsClient.writePump(s)
	wsClient.readPump(s)
}

// handleMetrics xử lý metrics endpoint
func (s *wsServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	s.metricsMu.RLock()
	metrics := *s.metrics
	s.metricsMu.RUnlock()

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
	}
}

// updateMetrics cập nhật metrics định kỳ
func (s *wsServer) updateMetrics() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.metricsMu.Lock()
			s.metrics.Uptime = int64(time.Since(s.startTime).Seconds())
			s.metrics.RoomsCount = len(s.hub.GetRooms())
			s.metrics.LastUpdated = time.Now()
			s.metricsMu.Unlock()
		}
	}
}

// Client management methods
func (s *wsServer) GetClient(id string) (ServerClient, bool) {
	return s.hub.GetClient(id)
}

func (s *wsServer) GetClients() []ServerClient {
	return s.hub.GetClients()
}

func (s *wsServer) GetClientCount() int {
	return s.hub.GetClientCount()
}

func (s *wsServer) DisconnectClient(id string) error {
	client, ok := s.hub.GetClient(id)
	if !ok {
		return fmt.Errorf("client not found")
	}
	return client.Close()
}

// Room management methods
func (s *wsServer) CreateRoom(name string, options *RoomOptions) error {
	return s.hub.CreateRoom(name, options)
}

func (s *wsServer) DeleteRoom(name string) error {
	return s.hub.DeleteRoom(name)
}

func (s *wsServer) GetRoom(name string) (Room, bool) {
	return s.hub.GetRoom(name)
}

func (s *wsServer) GetRooms() []Room {
	return s.hub.GetRooms()
}

// Broadcasting methods
func (s *wsServer) Broadcast(data []byte) error {
	msg := &Message{
		Type:      BinaryMessage,
		Data:      data,
		Timestamp: time.Now(),
	}
	return s.hub.Broadcast(msg)
}

func (s *wsServer) BroadcastText(text string) error {
	msg := &Message{
		Type:      TextMessage,
		Data:      []byte(text),
		Timestamp: time.Now(),
	}
	return s.hub.Broadcast(msg)
}

func (s *wsServer) BroadcastJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	return s.BroadcastText(string(data))
}

func (s *wsServer) BroadcastToRoom(room string, data []byte) error {
	msg := &Message{
		Type:      BinaryMessage,
		Data:      data,
		Timestamp: time.Now(),
		Room:      room,
	}
	return s.hub.BroadcastToRoom(room, msg)
}

func (s *wsServer) BroadcastToRoomText(room string, text string) error {
	msg := &Message{
		Type:      TextMessage,
		Data:      []byte(text),
		Timestamp: time.Now(),
		Room:      room,
	}
	return s.hub.BroadcastToRoom(room, msg)
}

func (s *wsServer) BroadcastToRoomJSON(room string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}
	return s.BroadcastToRoomText(room, string(data))
}

// Event handler setters
func (s *wsServer) OnConnect(handler func(ServerClient)) {
	s.handlersMu.Lock()
	s.onConnect = handler
	s.handlersMu.Unlock()
}

func (s *wsServer) OnDisconnect(handler func(ServerClient, error)) {
	s.handlersMu.Lock()
	s.onDisconnect = handler
	s.handlersMu.Unlock()
}

func (s *wsServer) OnMessage(handler func(ServerClient, []byte)) {
	s.handlersMu.Lock()
	s.onMessage = handler
	s.handlersMu.Unlock()
}

func (s *wsServer) OnTextMessage(handler func(ServerClient, string)) {
	s.handlersMu.Lock()
	s.onTextMessage = handler
	s.handlersMu.Unlock()
}

func (s *wsServer) OnBinaryMessage(handler func(ServerClient, []byte)) {
	s.handlersMu.Lock()
	s.onBinaryMessage = handler
	s.handlersMu.Unlock()
}

func (s *wsServer) OnError(handler func(ServerClient, error)) {
	s.handlersMu.Lock()
	s.onError = handler
	s.handlersMu.Unlock()
}

func (s *wsServer) OnRoomJoin(handler func(ServerClient, string)) {
	s.handlersMu.Lock()
	s.onRoomJoin = handler
	s.handlersMu.Unlock()
}

func (s *wsServer) OnRoomLeave(handler func(ServerClient, string)) {
	s.handlersMu.Lock()
	s.onRoomLeave = handler
	s.handlersMu.Unlock()
}

// Configuration methods
func (s *wsServer) SetOptions(options *ServerOptions) {
	if options != nil {
		mergeServerOptions(s.options, options)
	}
}

func (s *wsServer) GetOptions() *ServerOptions {
	return s.options
}

// Monitoring methods
func (s *wsServer) GetMetrics() *ServerMetrics {
	s.metricsMu.RLock()
	metrics := *s.metrics
	s.metricsMu.RUnlock()
	return &metrics
}

func (s *wsServer) GetHealth() *HealthStatus {
	status := &HealthStatus{
		Status:      "healthy",
		Checks:      make(map[string]string),
		Timestamp:   time.Now(),
		Uptime:      int64(time.Since(s.startTime).Seconds()),
		Version:     "1.0.0",
		Environment: "production",
	}

	// Check if server is running
	if atomic.LoadInt32(&s.running) == 0 {
		status.Status = "unhealthy"
		status.Checks["server"] = "not running"
	} else {
		status.Checks["server"] = "running"
	}

	// Check hub status
	if s.hub == nil {
		status.Status = "unhealthy"
		status.Checks["hub"] = "not initialized"
	} else {
		status.Checks["hub"] = "running"
	}

	// Check connection count
	activeConnections := atomic.LoadInt64(&s.metrics.ActiveConnections)
	if activeConnections > 10000 {
		status.Status = "degraded"
		status.Checks["connections"] = fmt.Sprintf("high load: %d connections", activeConnections)
	} else {
		status.Checks["connections"] = fmt.Sprintf("normal: %d connections", activeConnections)
	}

	return status
}

// HTTP handler methods
func (s *wsServer) GetHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(s.options.Path, s.handleWebSocket)

	if s.options.EnableMetrics && s.options.MetricsPath != "" {
		mux.HandleFunc(s.options.MetricsPath, s.handleMetrics)
	}

	return mux
}

func (s *wsServer) GetUpgrader() *websocket.Upgrader {
	return s.upgrader
}

// handleMessage xử lý message từ client
func (s *wsServer) handleMessage(client ServerClient, messageType int, data []byte) {
	// Update metrics
	atomic.AddInt64(&s.metrics.TotalMessages, 1)
	atomic.AddInt64(&s.metrics.TotalBytes, int64(len(data)))

	// Call appropriate handlers
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()

	switch messageType {
	case websocket.TextMessage:
		if s.onTextMessage != nil {
			go s.onTextMessage(client, string(data))
		}
		if s.onMessage != nil {
			go s.onMessage(client, data)
		}

	case websocket.BinaryMessage:
		if s.onBinaryMessage != nil {
			go s.onBinaryMessage(client, data)
		}
		if s.onMessage != nil {
			go s.onMessage(client, data)
		}
	}
}

// handleError xử lý error từ client
func (s *wsServer) handleError(client ServerClient, err error) {
	atomic.AddInt64(&s.metrics.ErrorCount, 1)

	s.handlersMu.RLock()
	if s.onError != nil {
		go s.onError(client, err)
	}
	s.handlersMu.RUnlock()
}

// handleRoomJoin xử lý client join room
func (s *wsServer) handleRoomJoin(client ServerClient, room string) {
	s.handlersMu.RLock()
	if s.onRoomJoin != nil {
		go s.onRoomJoin(client, room)
	}
	s.handlersMu.RUnlock()
}

// handleRoomLeave xử lý client leave room
func (s *wsServer) handleRoomLeave(client ServerClient, room string) {
	s.handlersMu.RLock()
	if s.onRoomLeave != nil {
		go s.onRoomLeave(client, room)
	}
	s.handlersMu.RUnlock()
}

// mergeServerOptions merge server options
func mergeServerOptions(dst, src *ServerOptions) {
	if src.Addr != "" {
		dst.Addr = src.Addr
	}
	if src.Path != "" {
		dst.Path = src.Path
	}
	if src.ReadTimeout > 0 {
		dst.ReadTimeout = src.ReadTimeout
	}
	if src.WriteTimeout > 0 {
		dst.WriteTimeout = src.WriteTimeout
	}
	if src.IdleTimeout > 0 {
		dst.IdleTimeout = src.IdleTimeout
	}
	if src.MaxMessageSize > 0 {
		dst.MaxMessageSize = src.MaxMessageSize
	}
	if src.MessageQueueSize > 0 {
		dst.MessageQueueSize = src.MessageQueueSize
	}
	if src.PingInterval > 0 {
		dst.PingInterval = src.PingInterval
	}
	if src.PongTimeout > 0 {
		dst.PongTimeout = src.PongTimeout
	}
	if src.RateLimit > 0 {
		dst.RateLimit = src.RateLimit
	}
	if src.RateBurst > 0 {
		dst.RateBurst = src.RateBurst
	}
	if src.RateWindow > 0 {
		dst.RateWindow = src.RateWindow
	}
	if src.CompressionLevel > 0 {
		dst.CompressionLevel = src.CompressionLevel
	}
	if src.JWTSecret != "" {
		dst.JWTSecret = src.JWTSecret
	}
	if src.TLSCertFile != "" {
		dst.TLSCertFile = src.TLSCertFile
	}
	if src.TLSKeyFile != "" {
		dst.TLSKeyFile = src.TLSKeyFile
	}
	if src.MetricsPath != "" {
		dst.MetricsPath = src.MetricsPath
	}
	if src.ShutdownTimeout > 0 {
		dst.ShutdownTimeout = src.ShutdownTimeout
	}
	if src.CheckOrigin != nil {
		dst.CheckOrigin = src.CheckOrigin
	}
	if src.AuthHandler != nil {
		dst.AuthHandler = src.AuthHandler
	}
	if src.Subprotocols != nil {
		dst.Subprotocols = src.Subprotocols
	}
	if src.AllowedOrigins != nil {
		dst.AllowedOrigins = src.AllowedOrigins
	}
	if src.AllowedHeaders != nil {
		dst.AllowedHeaders = src.AllowedHeaders
	}
	dst.AuthRequired = src.AuthRequired
	dst.EnableCompression = src.EnableCompression
	dst.EnableMetrics = src.EnableMetrics
}
