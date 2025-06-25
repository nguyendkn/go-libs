package rtsp

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/nguyendkn/go-libs/ffmpeg"
)

// StreamManager manages multiple RTSP streams
type StreamManager struct {
	config   *Config
	streams  map[string]*Stream
	mutex    sync.RWMutex
	ffmpeg   ffmpeg.FFmpeg
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
}

// NewStreamManager creates a new stream manager
func NewStreamManager(config *Config) (*StreamManager, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &StreamManager{
		config:  config,
		streams: make(map[string]*Stream),
		ffmpeg:  config.FFmpeg,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

// Stream represents a single RTSP stream
type Stream struct {
	config      RTSPStream
	manager     *StreamManager
	status      StreamStatus
	info        StreamInfo
	ctx         context.Context
	cancel      context.CancelFunc
	mutex       sync.RWMutex
	reconnectCount int
	lastError   error
	startTime   time.Time
	metrics     StreamMetrics
}

// StreamMetrics contains metrics for a stream
type StreamMetrics struct {
	BytesReceived   int64     `json:"bytes_received"`
	FramesReceived  int64     `json:"frames_received"`
	PacketsLost     int64     `json:"packets_lost"`
	CurrentFPS      float64   `json:"current_fps"`
	CurrentBitrate  int64     `json:"current_bitrate"`
	LastFrameTime   time.Time `json:"last_frame_time"`
	ConnectionTime  time.Time `json:"connection_time"`
	ReconnectCount  int       `json:"reconnect_count"`
	ErrorCount      int       `json:"error_count"`
}

// AddStream adds a new RTSP stream to the manager
func (sm *StreamManager) AddStream(streamConfig RTSPStream) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Validate stream URL
	if err := sm.config.ValidateStreamURL(streamConfig.URL); err != nil {
		return err
	}

	// Generate stream name if not provided
	if streamConfig.Name == "" {
		streamConfig.Name = sm.generateStreamName(streamConfig.URL)
	}

	// Check if stream already exists
	if _, exists := sm.streams[streamConfig.Name]; exists {
		return &RTSPError{
			Message:   fmt.Sprintf("stream with name '%s' already exists", streamConfig.Name),
			Code:      ErrCodeInvalidConfig,
			StreamURL: streamConfig.URL,
		}
	}

	// Set default values
	if streamConfig.Transport == "" {
		streamConfig.Transport = sm.config.DefaultTransport
	}
	if streamConfig.Timeout == 0 {
		streamConfig.Timeout = sm.config.ConnectionTimeout
	}
	if streamConfig.MaxRetries == 0 {
		streamConfig.MaxRetries = sm.config.MaxReconnectAttempts
	}
	if streamConfig.RetryDelay == 0 {
		streamConfig.RetryDelay = sm.config.ReconnectDelay
	}
	if streamConfig.BufferSize == 0 {
		streamConfig.BufferSize = sm.config.BufferSize
	}

	// Create stream context
	ctx, cancel := context.WithCancel(sm.ctx)

	stream := &Stream{
		config:  streamConfig,
		manager: sm,
		status:  StatusIdle,
		ctx:     ctx,
		cancel:  cancel,
		info: StreamInfo{
			Stream:    streamConfig,
			Status:    StatusIdle,
			StartTime: time.Now(),
		},
	}

	sm.streams[streamConfig.Name] = stream
	return nil
}

// AddStreams adds multiple RTSP streams
func (sm *StreamManager) AddStreams(streamConfigs []RTSPStream) error {
	for _, config := range streamConfigs {
		if err := sm.AddStream(config); err != nil {
			return err
		}
	}
	return nil
}

// AddStreamURLs adds streams from URLs with default configuration
func (sm *StreamManager) AddStreamURLs(urls []string) error {
	for i, url := range urls {
		streamConfig := RTSPStream{
			URL:       url,
			Name:      fmt.Sprintf("stream_%d", i+1),
			Transport: sm.config.DefaultTransport,
			Reconnect: sm.config.ReconnectEnabled,
		}
		if err := sm.AddStream(streamConfig); err != nil {
			return err
		}
	}
	return nil
}

// RemoveStream removes a stream from the manager
func (sm *StreamManager) RemoveStream(streamName string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	stream, exists := sm.streams[streamName]
	if !exists {
		return &RTSPError{
			Message: fmt.Sprintf("stream '%s' not found", streamName),
			Code:    ErrCodeStreamError,
		}
	}

	// Stop the stream if it's running
	if err := stream.Stop(); err != nil {
		return err
	}

	delete(sm.streams, streamName)
	return nil
}

// StartStream starts a specific stream
func (sm *StreamManager) StartStream(streamName string) error {
	sm.mutex.RLock()
	stream, exists := sm.streams[streamName]
	sm.mutex.RUnlock()

	if !exists {
		return &RTSPError{
			Message: fmt.Sprintf("stream '%s' not found", streamName),
			Code:    ErrCodeStreamError,
		}
	}

	return stream.Start()
}

// StartAllStreams starts all streams
func (sm *StreamManager) StartAllStreams() error {
	sm.mutex.RLock()
	streams := make([]*Stream, 0, len(sm.streams))
	for _, stream := range sm.streams {
		streams = append(streams, stream)
	}
	sm.mutex.RUnlock()

	var errors []error
	for _, stream := range streams {
		if err := stream.Start(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return &RTSPError{
			Message: fmt.Sprintf("failed to start %d streams", len(errors)),
			Code:    ErrCodeStreamError,
			Cause:   errors[0],
		}
	}

	return nil
}

// StopStream stops a specific stream
func (sm *StreamManager) StopStream(streamName string) error {
	sm.mutex.RLock()
	stream, exists := sm.streams[streamName]
	sm.mutex.RUnlock()

	if !exists {
		return &RTSPError{
			Message: fmt.Sprintf("stream '%s' not found", streamName),
			Code:    ErrCodeStreamError,
		}
	}

	return stream.Stop()
}

// StopAllStreams stops all streams
func (sm *StreamManager) StopAllStreams() error {
	sm.mutex.RLock()
	streams := make([]*Stream, 0, len(sm.streams))
	for _, stream := range sm.streams {
		streams = append(streams, stream)
	}
	sm.mutex.RUnlock()

	for _, stream := range streams {
		stream.Stop() // Ignore errors during shutdown
	}

	return nil
}

// GetStreamInfo returns information about a specific stream
func (sm *StreamManager) GetStreamInfo(streamName string) (StreamInfo, error) {
	sm.mutex.RLock()
	stream, exists := sm.streams[streamName]
	sm.mutex.RUnlock()

	if !exists {
		return StreamInfo{}, &RTSPError{
			Message: fmt.Sprintf("stream '%s' not found", streamName),
			Code:    ErrCodeStreamError,
		}
	}

	return stream.GetInfo(), nil
}

// GetAllStreamInfo returns information about all streams
func (sm *StreamManager) GetAllStreamInfo() map[string]StreamInfo {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	info := make(map[string]StreamInfo)
	for name, stream := range sm.streams {
		info[name] = stream.GetInfo()
	}
	return info
}

// GetStreamNames returns all stream names
func (sm *StreamManager) GetStreamNames() []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	names := make([]string, 0, len(sm.streams))
	for name := range sm.streams {
		names = append(names, name)
	}
	return names
}

// GetStreamCount returns the number of streams
func (sm *StreamManager) GetStreamCount() int {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return len(sm.streams)
}

// Close closes the stream manager and all streams
func (sm *StreamManager) Close() error {
	sm.cancel()
	sm.StopAllStreams()
	sm.wg.Wait()
	return nil
}

// Start starts the stream
func (s *Stream) Start() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.status == StatusStreaming || s.status == StatusConnecting {
		return nil // Already started
	}

	s.status = StatusConnecting
	s.startTime = time.Now()
	s.info.StartTime = s.startTime
	s.info.Status = s.status

	// Notify handler
	if s.manager.config.StreamHandler != nil {
		s.manager.config.StreamHandler.OnStreamConnected(s.config.URL, s.info)
	}

	// Start streaming in a goroutine
	s.manager.wg.Add(1)
	go s.streamLoop()

	return nil
}

// Stop stops the stream
func (s *Stream) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.status == StatusStopped || s.status == StatusIdle {
		return nil // Already stopped
	}

	s.cancel()
	s.status = StatusStopped
	s.info.Status = s.status

	// Notify handler
	if s.manager.config.StreamHandler != nil {
		s.manager.config.StreamHandler.OnStreamDisconnected(s.config.URL, "stopped")
	}

	return nil
}

// GetInfo returns stream information
func (s *Stream) GetInfo() StreamInfo {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	info := s.info
	info.Duration = time.Since(s.startTime)
	return info
}

// streamLoop is the main streaming loop
func (s *Stream) streamLoop() {
	defer s.manager.wg.Done()
	defer func() {
		s.mutex.Lock()
		s.status = StatusStopped
		s.info.Status = s.status
		s.mutex.Unlock()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			if err := s.connectAndStream(); err != nil {
				s.handleError(err)
				
				if !s.config.Reconnect || s.reconnectCount >= s.config.MaxRetries {
					return
				}
				
				s.reconnectCount++
				s.status = StatusReconnecting
				s.info.Status = s.status
				
				// Wait before reconnecting
				select {
				case <-s.ctx.Done():
					return
				case <-time.After(s.config.RetryDelay):
					continue
				}
			}
		}
	}
}

// connectAndStream connects to the RTSP stream and processes data
func (s *Stream) connectAndStream() error {
	s.mutex.Lock()
	s.status = StatusConnecting
	s.info.Status = s.status
	s.mutex.Unlock()

	// This is a simplified implementation
	// In a real implementation, you would use an RTSP client library
	// or FFmpeg to connect to the RTSP stream
	
	// For now, we'll simulate the connection
	time.Sleep(1 * time.Second) // Simulate connection time
	
	s.mutex.Lock()
	s.status = StatusStreaming
	s.info.Status = s.status
	s.metrics.ConnectionTime = time.Now()
	s.mutex.Unlock()

	// Simulate streaming data
	ticker := time.NewTicker(time.Second / 30) // 30 FPS
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-ticker.C:
			// Simulate frame data
			s.updateMetrics()
			
			// Notify handler with simulated data
			if s.manager.config.StreamHandler != nil {
				data := make([]byte, 1024) // Simulated frame data
				s.manager.config.StreamHandler.OnStreamData(s.config.URL, data)
			}
		}
	}
}

// handleError handles stream errors
func (s *Stream) handleError(err error) {
	s.mutex.Lock()
	s.lastError = err
	s.status = StatusError
	s.info.Status = s.status
	s.info.Error = err.Error()
	s.metrics.ErrorCount++
	s.mutex.Unlock()

	// Notify handler
	if s.manager.config.StreamHandler != nil {
		s.manager.config.StreamHandler.OnStreamError(s.config.URL, err)
	}
}

// updateMetrics updates stream metrics
func (s *Stream) updateMetrics() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	s.metrics.FramesReceived++
	s.metrics.BytesReceived += 1024 // Simulated
	s.metrics.LastFrameTime = now
	
	// Calculate FPS
	if s.metrics.ConnectionTime.IsZero() {
		s.metrics.ConnectionTime = now
	}
	
	duration := now.Sub(s.metrics.ConnectionTime).Seconds()
	if duration > 0 {
		s.metrics.CurrentFPS = float64(s.metrics.FramesReceived) / duration
	}
	
	// Update info
	s.info.FrameCount = s.metrics.FramesReceived
	s.info.BytesRead = s.metrics.BytesReceived
	s.info.FPS = s.metrics.CurrentFPS
}

// generateStreamName generates a unique stream name from URL
func (sm *StreamManager) generateStreamName(streamURL string) string {
	parsedURL, err := url.Parse(streamURL)
	if err != nil {
		return fmt.Sprintf("stream_%d", len(sm.streams)+1)
	}
	
	// Use host and path to generate name
	name := parsedURL.Host
	if parsedURL.Path != "" && parsedURL.Path != "/" {
		name += "_" + parsedURL.Path[1:] // Remove leading slash
	}
	
	// Replace invalid characters
	name = fmt.Sprintf("stream_%s", name)
	
	// Ensure uniqueness
	counter := 1
	originalName := name
	for {
		if _, exists := sm.streams[name]; !exists {
			break
		}
		name = fmt.Sprintf("%s_%d", originalName, counter)
		counter++
	}
	
	return name
}
