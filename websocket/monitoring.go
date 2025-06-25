package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// MetricsCollector thu thập và quản lý metrics
type MetricsCollector struct {
	// Connection metrics
	totalConnections    int64
	activeConnections   int64
	totalDisconnections int64
	
	// Message metrics
	totalMessages       int64
	totalBytes          int64
	messagesPerSecond   int64
	bytesPerSecond      int64
	
	// Error metrics
	totalErrors         int64
	connectionErrors    int64
	messageErrors       int64
	authErrors          int64
	
	// Performance metrics
	averageLatency      int64 // in microseconds
	maxLatency          int64
	minLatency          int64
	latencyCount        int64
	
	// Room metrics
	totalRooms          int64
	activeRooms         int64
	
	// Rate limiting metrics
	rateLimitHits       int64
	
	// Timing
	startTime           time.Time
	lastReset           time.Time
	
	// Mutex for complex operations
	mu sync.RWMutex
	
	// Historical data
	history []MetricsSnapshot
	maxHistorySize int
}

// MetricsSnapshot lưu trữ snapshot của metrics tại một thời điểm
type MetricsSnapshot struct {
	Timestamp           time.Time `json:"timestamp"`
	TotalConnections    int64     `json:"total_connections"`
	ActiveConnections   int64     `json:"active_connections"`
	TotalMessages       int64     `json:"total_messages"`
	TotalBytes          int64     `json:"total_bytes"`
	MessagesPerSecond   float64   `json:"messages_per_second"`
	BytesPerSecond      float64   `json:"bytes_per_second"`
	AverageLatency      float64   `json:"average_latency_ms"`
	ErrorRate           float64   `json:"error_rate"`
	ActiveRooms         int64     `json:"active_rooms"`
}

// NewMetricsCollector tạo một metrics collector mới
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		startTime:      time.Now(),
		lastReset:      time.Now(),
		maxHistorySize: 1000, // Keep last 1000 snapshots
		history:        make([]MetricsSnapshot, 0, 1000),
		minLatency:     int64(^uint64(0) >> 1), // Max int64
	}
}

// Connection metrics
func (mc *MetricsCollector) IncrementConnections() {
	atomic.AddInt64(&mc.totalConnections, 1)
	atomic.AddInt64(&mc.activeConnections, 1)
}

func (mc *MetricsCollector) DecrementConnections() {
	atomic.AddInt64(&mc.activeConnections, -1)
	atomic.AddInt64(&mc.totalDisconnections, 1)
}

func (mc *MetricsCollector) GetActiveConnections() int64 {
	return atomic.LoadInt64(&mc.activeConnections)
}

func (mc *MetricsCollector) GetTotalConnections() int64 {
	return atomic.LoadInt64(&mc.totalConnections)
}

// Message metrics
func (mc *MetricsCollector) IncrementMessages(bytes int64) {
	atomic.AddInt64(&mc.totalMessages, 1)
	atomic.AddInt64(&mc.totalBytes, bytes)
}

func (mc *MetricsCollector) GetTotalMessages() int64 {
	return atomic.LoadInt64(&mc.totalMessages)
}

func (mc *MetricsCollector) GetTotalBytes() int64 {
	return atomic.LoadInt64(&mc.totalBytes)
}

// Error metrics
func (mc *MetricsCollector) IncrementErrors() {
	atomic.AddInt64(&mc.totalErrors, 1)
}

func (mc *MetricsCollector) IncrementConnectionErrors() {
	atomic.AddInt64(&mc.connectionErrors, 1)
	mc.IncrementErrors()
}

func (mc *MetricsCollector) IncrementMessageErrors() {
	atomic.AddInt64(&mc.messageErrors, 1)
	mc.IncrementErrors()
}

func (mc *MetricsCollector) IncrementAuthErrors() {
	atomic.AddInt64(&mc.authErrors, 1)
	mc.IncrementErrors()
}

func (mc *MetricsCollector) IncrementRateLimitHits() {
	atomic.AddInt64(&mc.rateLimitHits, 1)
}

// Latency metrics
func (mc *MetricsCollector) RecordLatency(latency time.Duration) {
	latencyMicros := latency.Microseconds()
	
	// Update average
	count := atomic.AddInt64(&mc.latencyCount, 1)
	currentAvg := atomic.LoadInt64(&mc.averageLatency)
	newAvg := (currentAvg*(count-1) + latencyMicros) / count
	atomic.StoreInt64(&mc.averageLatency, newAvg)
	
	// Update max
	for {
		currentMax := atomic.LoadInt64(&mc.maxLatency)
		if latencyMicros <= currentMax {
			break
		}
		if atomic.CompareAndSwapInt64(&mc.maxLatency, currentMax, latencyMicros) {
			break
		}
	}
	
	// Update min
	for {
		currentMin := atomic.LoadInt64(&mc.minLatency)
		if latencyMicros >= currentMin {
			break
		}
		if atomic.CompareAndSwapInt64(&mc.minLatency, currentMin, latencyMicros) {
			break
		}
	}
}

// Room metrics
func (mc *MetricsCollector) SetActiveRooms(count int64) {
	atomic.StoreInt64(&mc.activeRooms, count)
	
	// Update total rooms if current is higher
	for {
		currentTotal := atomic.LoadInt64(&mc.totalRooms)
		if count <= currentTotal {
			break
		}
		if atomic.CompareAndSwapInt64(&mc.totalRooms, currentTotal, count) {
			break
		}
	}
}

// GetSnapshot tạo snapshot của metrics hiện tại
func (mc *MetricsCollector) GetSnapshot() MetricsSnapshot {
	now := time.Now()
	duration := now.Sub(mc.lastReset).Seconds()
	
	snapshot := MetricsSnapshot{
		Timestamp:         now,
		TotalConnections:  atomic.LoadInt64(&mc.totalConnections),
		ActiveConnections: atomic.LoadInt64(&mc.activeConnections),
		TotalMessages:     atomic.LoadInt64(&mc.totalMessages),
		TotalBytes:        atomic.LoadInt64(&mc.totalBytes),
		ActiveRooms:       atomic.LoadInt64(&mc.activeRooms),
	}
	
	// Calculate rates
	if duration > 0 {
		snapshot.MessagesPerSecond = float64(snapshot.TotalMessages) / duration
		snapshot.BytesPerSecond = float64(snapshot.TotalBytes) / duration
	}
	
	// Calculate latency
	avgLatencyMicros := atomic.LoadInt64(&mc.averageLatency)
	snapshot.AverageLatency = float64(avgLatencyMicros) / 1000.0 // Convert to milliseconds
	
	// Calculate error rate
	totalErrors := atomic.LoadInt64(&mc.totalErrors)
	if snapshot.TotalMessages > 0 {
		snapshot.ErrorRate = float64(totalErrors) / float64(snapshot.TotalMessages) * 100
	}
	
	return snapshot
}

// AddSnapshot thêm snapshot vào history
func (mc *MetricsCollector) AddSnapshot(snapshot MetricsSnapshot) {
	mc.mu.Lock()
	mc.history = append(mc.history, snapshot)
	
	// Keep only the last N snapshots
	if len(mc.history) > mc.maxHistorySize {
		mc.history = mc.history[len(mc.history)-mc.maxHistorySize:]
	}
	mc.mu.Unlock()
}

// GetHistory trả về history của metrics
func (mc *MetricsCollector) GetHistory(limit int) []MetricsSnapshot {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	
	if limit <= 0 || limit > len(mc.history) {
		limit = len(mc.history)
	}
	
	start := len(mc.history) - limit
	history := make([]MetricsSnapshot, limit)
	copy(history, mc.history[start:])
	
	return history
}

// Reset reset metrics
func (mc *MetricsCollector) Reset() {
	atomic.StoreInt64(&mc.totalConnections, 0)
	atomic.StoreInt64(&mc.totalDisconnections, 0)
	atomic.StoreInt64(&mc.totalMessages, 0)
	atomic.StoreInt64(&mc.totalBytes, 0)
	atomic.StoreInt64(&mc.totalErrors, 0)
	atomic.StoreInt64(&mc.connectionErrors, 0)
	atomic.StoreInt64(&mc.messageErrors, 0)
	atomic.StoreInt64(&mc.authErrors, 0)
	atomic.StoreInt64(&mc.rateLimitHits, 0)
	atomic.StoreInt64(&mc.averageLatency, 0)
	atomic.StoreInt64(&mc.maxLatency, 0)
	atomic.StoreInt64(&mc.minLatency, int64(^uint64(0)>>1))
	atomic.StoreInt64(&mc.latencyCount, 0)
	atomic.StoreInt64(&mc.totalRooms, 0)
	
	mc.lastReset = time.Now()
	
	mc.mu.Lock()
	mc.history = mc.history[:0]
	mc.mu.Unlock()
}

// GetUptime trả về uptime
func (mc *MetricsCollector) GetUptime() time.Duration {
	return time.Since(mc.startTime)
}

// HealthChecker kiểm tra health của hệ thống
type HealthChecker struct {
	checks map[string]HealthCheck
	mu     sync.RWMutex
}

// HealthCheck định nghĩa một health check
type HealthCheck struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	CheckFunc   func() (bool, string)  `json:"-"`
	LastCheck   time.Time              `json:"last_check"`
	LastResult  bool                   `json:"last_result"`
	LastMessage string                 `json:"last_message"`
}

// NewHealthChecker tạo một health checker mới
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]HealthCheck),
	}
}

// AddCheck thêm một health check
func (hc *HealthChecker) AddCheck(name, description string, checkFunc func() (bool, string)) {
	hc.mu.Lock()
	hc.checks[name] = HealthCheck{
		Name:        name,
		Description: description,
		CheckFunc:   checkFunc,
	}
	hc.mu.Unlock()
}

// RunChecks chạy tất cả health checks
func (hc *HealthChecker) RunChecks() map[string]HealthCheck {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	
	results := make(map[string]HealthCheck)
	
	for name, check := range hc.checks {
		result, message := check.CheckFunc()
		
		check.LastCheck = time.Now()
		check.LastResult = result
		check.LastMessage = message
		
		hc.checks[name] = check
		results[name] = check
	}
	
	return results
}

// GetOverallHealth trả về overall health status
func (hc *HealthChecker) GetOverallHealth() (bool, string) {
	checks := hc.RunChecks()
	
	healthy := true
	var issues []string
	
	for _, check := range checks {
		if !check.LastResult {
			healthy = false
			issues = append(issues, fmt.Sprintf("%s: %s", check.Name, check.LastMessage))
		}
	}
	
	if healthy {
		return true, "All checks passed"
	}
	
	return false, fmt.Sprintf("Failed checks: %v", issues)
}

// MonitoringServer cung cấp HTTP endpoints cho monitoring
type MonitoringServer struct {
	metricsCollector *MetricsCollector
	healthChecker    *HealthChecker
	server           *http.Server
}

// NewMonitoringServer tạo một monitoring server mới
func NewMonitoringServer(addr string, metricsCollector *MetricsCollector, healthChecker *HealthChecker) *MonitoringServer {
	ms := &MonitoringServer{
		metricsCollector: metricsCollector,
		healthChecker:    healthChecker,
	}
	
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", ms.handleMetrics)
	mux.HandleFunc("/health", ms.handleHealth)
	mux.HandleFunc("/metrics/history", ms.handleMetricsHistory)
	mux.HandleFunc("/metrics/reset", ms.handleMetricsReset)
	
	ms.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	
	return ms
}

// Start khởi động monitoring server
func (ms *MonitoringServer) Start() error {
	return ms.server.ListenAndServe()
}

// Stop dừng monitoring server
func (ms *MonitoringServer) Stop() error {
	return ms.server.Close()
}

// handleMetrics xử lý metrics endpoint
func (ms *MonitoringServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	snapshot := ms.metricsCollector.GetSnapshot()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snapshot)
}

// handleHealth xử lý health endpoint
func (ms *MonitoringServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	healthy, message := ms.healthChecker.GetOverallHealth()
	
	status := map[string]interface{}{
		"healthy":   healthy,
		"message":   message,
		"timestamp": time.Now(),
		"uptime":    ms.metricsCollector.GetUptime().Seconds(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	
	if healthy {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	
	json.NewEncoder(w).Encode(status)
}

// handleMetricsHistory xử lý metrics history endpoint
func (ms *MonitoringServer) handleMetricsHistory(w http.ResponseWriter, r *http.Request) {
	limit := 100 // default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := fmt.Sscanf(limitStr, "%d", &limit); err != nil || parsedLimit != 1 {
			limit = 100
		}
	}
	
	history := ms.metricsCollector.GetHistory(limit)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// handleMetricsReset xử lý metrics reset endpoint
func (ms *MonitoringServer) handleMetricsReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	
	ms.metricsCollector.Reset()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Metrics reset successfully",
	})
}
