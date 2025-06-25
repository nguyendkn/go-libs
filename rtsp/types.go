package rtsp

import (
	"context"
	"time"

	"github.com/nguyendkn/go-libs/hls"
)

// RTSPStream represents a single RTSP stream configuration
type RTSPStream struct {
	URL        string            `json:"url"`
	Name       string            `json:"name,omitempty"`
	Username   string            `json:"username,omitempty"`
	Password   string            `json:"password,omitempty"`
	Transport  TransportProtocol `json:"transport"`
	Timeout    time.Duration     `json:"timeout"`
	Reconnect  bool              `json:"reconnect"`
	MaxRetries int               `json:"max_retries"`
	RetryDelay time.Duration     `json:"retry_delay"`
	BufferSize int               `json:"buffer_size"`
	Position   *Position         `json:"position,omitempty"` // For layout positioning
}

// Position represents the position of a stream in a layout
type Position struct {
	Row    int `json:"row"`
	Column int `json:"column"`
	Width  int `json:"width,omitempty"`  // Custom width (optional)
	Height int `json:"height,omitempty"` // Custom height (optional)
}

// TransportProtocol represents RTSP transport protocols
type TransportProtocol string

const (
	TransportTCP  TransportProtocol = "tcp"
	TransportUDP  TransportProtocol = "udp"
	TransportAuto TransportProtocol = "auto"
)

// LayoutType represents different layout configurations
type LayoutType string

const (
	LayoutSingle LayoutType = "single" // 1x1 - Single stream
	Layout1x2    LayoutType = "1x2"    // 1 row, 2 columns
	Layout2x1    LayoutType = "2x1"    // 2 rows, 1 column
	Layout2x2    LayoutType = "2x2"    // 2x2 grid
	Layout2x3    LayoutType = "2x3"    // 2 rows, 3 columns
	Layout3x2    LayoutType = "3x2"    // 3 rows, 2 columns
	Layout3x3    LayoutType = "3x3"    // 3x3 grid
	Layout4x4    LayoutType = "4x4"    // 4x4 grid
	LayoutCustom LayoutType = "custom" // Custom layout
)

// Layout represents the video layout configuration
type Layout struct {
	Type        LayoutType `json:"type"`
	Rows        int        `json:"rows"`
	Columns     int        `json:"columns"`
	Width       int        `json:"width"`        // Output video width
	Height      int        `json:"height"`       // Output video height
	Padding     int        `json:"padding"`      // Padding between streams
	Background  string     `json:"background"`   // Background color (hex)
	BorderWidth int        `json:"border_width"` // Border width around each stream
	BorderColor string     `json:"border_color"` // Border color (hex)
}

// StreamingMode represents different streaming modes
type StreamingMode string

const (
	ModeSeparate StreamingMode = "separate" // Each stream as separate HLS
	ModeMerged   StreamingMode = "merged"   // All streams merged into one HLS
	ModeBoth     StreamingMode = "both"     // Both separate and merged
)

// OutputFormat represents output format options
type OutputFormat string

const (
	FormatHLS  OutputFormat = "hls"
	FormatDASH OutputFormat = "dash"
	FormatRTMP OutputFormat = "rtmp"
	FormatFile OutputFormat = "file"
)

// StreamStatus represents the status of a stream
type StreamStatus string

const (
	StatusIdle         StreamStatus = "idle"
	StatusConnecting   StreamStatus = "connecting"
	StatusConnected    StreamStatus = "connected"
	StatusStreaming    StreamStatus = "streaming"
	StatusReconnecting StreamStatus = "reconnecting"
	StatusError        StreamStatus = "error"
	StatusStopped      StreamStatus = "stopped"
)

// StreamInfo contains information about a stream
type StreamInfo struct {
	Stream     RTSPStream    `json:"stream"`
	Status     StreamStatus  `json:"status"`
	StartTime  time.Time     `json:"start_time"`
	Duration   time.Duration `json:"duration"`
	BytesRead  int64         `json:"bytes_read"`
	FrameCount int64         `json:"frame_count"`
	FPS        float64       `json:"fps"`
	Bitrate    int64         `json:"bitrate"`
	Resolution string        `json:"resolution"`
	Error      string        `json:"error,omitempty"`
}

// ConversionProgress represents the progress of RTSP to HLS conversion
type ConversionProgress struct {
	StreamURL    string                  `json:"stream_url"`
	StreamName   string                  `json:"stream_name"`
	Status       StreamStatus            `json:"status"`
	Progress     float64                 `json:"progress"`      // 0-100
	Duration     time.Duration           `json:"duration"`      // Total duration processed
	Speed        float64                 `json:"speed"`         // Processing speed multiplier
	FPS          float64                 `json:"fps"`           // Current FPS
	Bitrate      int64                   `json:"bitrate"`       // Current bitrate
	SegmentCount int                     `json:"segment_count"` // Number of segments created
	Error        string                  `json:"error,omitempty"`
	HLSProgress  *hls.ConversionProgress `json:"hls_progress,omitempty"`
}

// ConversionResult represents the result of RTSP to HLS conversion
type ConversionResult struct {
	Success        bool                    `json:"success"`
	Mode           StreamingMode           `json:"mode"`
	Layout         *Layout                 `json:"layout,omitempty"`
	OutputDir      string                  `json:"output_dir"`
	MasterPlaylist string                  `json:"master_playlist,omitempty"`
	Streams        map[string]StreamResult `json:"streams"` // stream_name -> result
	MergedStream   *StreamResult           `json:"merged_stream,omitempty"`
	Duration       time.Duration           `json:"duration"`
	StartTime      time.Time               `json:"start_time"`
	EndTime        time.Time               `json:"end_time"`
	Error          string                  `json:"error,omitempty"`
	Stats          ConversionStats         `json:"stats"`
}

// StreamResult represents the result for a single stream
type StreamResult struct {
	StreamName     string                `json:"stream_name"`
	StreamURL      string                `json:"stream_url"`
	Success        bool                  `json:"success"`
	OutputDir      string                `json:"output_dir"`
	PlaylistPath   string                `json:"playlist_path"`
	Segments       []string              `json:"segments"`
	Duration       time.Duration         `json:"duration"`
	SegmentCount   int                   `json:"segment_count"`
	TotalSize      int64                 `json:"total_size"`
	AverageBitrate int64                 `json:"average_bitrate"`
	Resolution     string                `json:"resolution"`
	FPS            float64               `json:"fps"`
	Error          string                `json:"error,omitempty"`
	HLSResult      *hls.ConversionResult `json:"hls_result,omitempty"`
}

// ConversionStats contains statistics about the conversion
type ConversionStats struct {
	TotalStreams      int           `json:"total_streams"`
	SuccessfulStreams int           `json:"successful_streams"`
	FailedStreams     int           `json:"failed_streams"`
	TotalDuration     time.Duration `json:"total_duration"`
	TotalSize         int64         `json:"total_size"`
	AverageFPS        float64       `json:"average_fps"`
	AverageBitrate    int64         `json:"average_bitrate"`
	SegmentCount      int           `json:"segment_count"`
	ReconnectCount    int           `json:"reconnect_count"`
}

// RTSPError represents RTSP-specific errors
type RTSPError struct {
	Message   string `json:"message"`
	Code      string `json:"code"`
	StreamURL string `json:"stream_url,omitempty"`
	Cause     error  `json:"cause,omitempty"`
}

func (e *RTSPError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *RTSPError) Unwrap() error {
	return e.Cause
}

// Error codes
const (
	ErrCodeInvalidURL       = "INVALID_URL"
	ErrCodeConnectionFailed = "CONNECTION_FAILED"
	ErrCodeAuthFailed       = "AUTH_FAILED"
	ErrCodeTimeout          = "TIMEOUT"
	ErrCodeStreamError      = "STREAM_ERROR"
	ErrCodeConversionFailed = "CONVERSION_FAILED"
	ErrCodeLayoutError      = "LAYOUT_ERROR"
	ErrCodeFFmpegError      = "FFMPEG_ERROR"
	ErrCodeHLSError         = "HLS_ERROR"
	ErrCodeInvalidConfig    = "INVALID_CONFIG"
)

// StreamHandler defines the interface for handling stream events
type StreamHandler interface {
	OnStreamConnected(streamURL string, info StreamInfo)
	OnStreamDisconnected(streamURL string, reason string)
	OnStreamError(streamURL string, err error)
	OnStreamData(streamURL string, data []byte)
	OnConversionProgress(progress ConversionProgress)
}

// DefaultStreamHandler provides a default implementation
type DefaultStreamHandler struct{}

func (h *DefaultStreamHandler) OnStreamConnected(streamURL string, info StreamInfo)  {}
func (h *DefaultStreamHandler) OnStreamDisconnected(streamURL string, reason string) {}
func (h *DefaultStreamHandler) OnStreamError(streamURL string, err error)            {}
func (h *DefaultStreamHandler) OnStreamData(streamURL string, data []byte)           {}
func (h *DefaultStreamHandler) OnConversionProgress(progress ConversionProgress)     {}

// ContextKey represents context keys for RTSP operations
type ContextKey string

const (
	ContextKeyStreamURL  ContextKey = "stream_url"
	ContextKeyStreamName ContextKey = "stream_name"
	ContextKeyOperation  ContextKey = "operation"
)

// Operation represents different RTSP operations
type Operation string

const (
	OperationConnect   Operation = "connect"
	OperationStream    Operation = "stream"
	OperationConvert   Operation = "convert"
	OperationReconnect Operation = "reconnect"
	OperationStop      Operation = "stop"
)

// Predefined layouts
var (
	DefaultLayouts = map[LayoutType]Layout{
		LayoutSingle: {Type: LayoutSingle, Rows: 1, Columns: 1, Width: 1920, Height: 1080},
		Layout1x2:    {Type: Layout1x2, Rows: 1, Columns: 2, Width: 1920, Height: 1080},
		Layout2x1:    {Type: Layout2x1, Rows: 2, Columns: 1, Width: 1920, Height: 1080},
		Layout2x2:    {Type: Layout2x2, Rows: 2, Columns: 2, Width: 1920, Height: 1080},
		Layout2x3:    {Type: Layout2x3, Rows: 2, Columns: 3, Width: 1920, Height: 1080},
		Layout3x2:    {Type: Layout3x2, Rows: 3, Columns: 2, Width: 1920, Height: 1080},
		Layout3x3:    {Type: Layout3x3, Rows: 3, Columns: 3, Width: 1920, Height: 1080},
		Layout4x4:    {Type: Layout4x4, Rows: 4, Columns: 4, Width: 1920, Height: 1080},
	}
)

// StreamContext contains context information for stream operations
type StreamContext struct {
	Context    context.Context
	StreamURL  string
	StreamName string
	Operation  Operation
	StartTime  time.Time
	Metadata   map[string]interface{}
}

// NewStreamContext creates a new stream context
func NewStreamContext(ctx context.Context, streamURL, streamName string, op Operation) *StreamContext {
	return &StreamContext{
		Context:    ctx,
		StreamURL:  streamURL,
		StreamName: streamName,
		Operation:  op,
		StartTime:  time.Now(),
		Metadata:   make(map[string]interface{}),
	}
}
