package rtsp

import (
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/nguyendkn/go-libs/ffmpeg"
	"github.com/nguyendkn/go-libs/hls"
)

// Config represents RTSP streaming configuration
type Config struct {
	// Core dependencies
	FFmpeg ffmpeg.FFmpeg `json:"-"`
	HLS    hls.HLS       `json:"-"`

	// Output settings
	OutputDir     string        `json:"output_dir"`
	BaseURL       string        `json:"base_url,omitempty"`
	StreamingMode StreamingMode `json:"streaming_mode"`
	OutputFormat  OutputFormat  `json:"output_format"`

	// Layout settings
	Layout          Layout `json:"layout"`
	AutoLayout      bool   `json:"auto_layout"`      // Auto-detect layout based on stream count
	LayoutPadding   int    `json:"layout_padding"`   // Padding between streams
	LayoutBorder    int    `json:"layout_border"`    // Border around each stream
	BackgroundColor string `json:"background_color"` // Background color for empty spaces

	// Stream settings
	DefaultTransport     TransportProtocol `json:"default_transport"`
	ConnectionTimeout    time.Duration     `json:"connection_timeout"`
	ReadTimeout          time.Duration     `json:"read_timeout"`
	ReconnectEnabled     bool              `json:"reconnect_enabled"`
	MaxReconnectAttempts int               `json:"max_reconnect_attempts"`
	ReconnectDelay       time.Duration     `json:"reconnect_delay"`
	BufferSize           int               `json:"buffer_size"`

	// HLS conversion settings
	HLSConfig         *hls.Config   `json:"hls_config,omitempty"`
	SegmentDuration   time.Duration `json:"segment_duration"`
	PlaylistSize      int           `json:"playlist_size"`
	DeleteOldSegments bool          `json:"delete_old_segments"`

	// Quality settings
	VideoCodec   ffmpeg.VideoCodec `json:"video_codec"`
	AudioCodec   ffmpeg.AudioCodec `json:"audio_codec"`
	VideoBitrate string            `json:"video_bitrate"`
	AudioBitrate string            `json:"audio_bitrate"`
	Resolution   ffmpeg.Resolution `json:"resolution"`
	FrameRate    float64           `json:"frame_rate"`

	// Performance settings
	Parallel          bool          `json:"parallel"`
	MaxConcurrent     int           `json:"max_concurrent"`
	ProcessingTimeout time.Duration `json:"processing_timeout"`
	MemoryLimit       int64         `json:"memory_limit"` // Memory limit in bytes

	// Monitoring settings
	EnableMetrics       bool          `json:"enable_metrics"`
	MetricsInterval     time.Duration `json:"metrics_interval"`
	HealthCheckInterval time.Duration `json:"health_check_interval"`

	// Callback functions
	StreamHandler    StreamHandler            `json:"-"`
	ProgressCallback func(ConversionProgress) `json:"-"`
	ErrorCallback    func(error)              `json:"-"`

	// Advanced settings
	FFmpegArgs    []string `json:"ffmpeg_args,omitempty"`
	CustomFilters []string `json:"custom_filters,omitempty"`
	LogLevel      string   `json:"log_level"`
	TempDir       string   `json:"temp_dir"`
	CleanupTemp   bool     `json:"cleanup_temp"`
}

// DefaultConfig returns a default RTSP configuration
func DefaultConfig() *Config {
	return &Config{
		OutputDir:            "output",
		StreamingMode:        ModeSeparate,
		OutputFormat:         FormatHLS,
		Layout:               DefaultLayouts[Layout2x2],
		AutoLayout:           true,
		LayoutPadding:        10,
		LayoutBorder:         2,
		BackgroundColor:      "#000000",
		DefaultTransport:     TransportTCP,
		ConnectionTimeout:    30 * time.Second,
		ReadTimeout:          10 * time.Second,
		ReconnectEnabled:     true,
		MaxReconnectAttempts: 5,
		ReconnectDelay:       5 * time.Second,
		BufferSize:           1024 * 1024, // 1MB
		SegmentDuration:      6 * time.Second,
		PlaylistSize:         10,
		DeleteOldSegments:    true,
		VideoCodec:           ffmpeg.VideoCodecH264,
		AudioCodec:           ffmpeg.AudioCodecAAC,
		VideoBitrate:         "2000k",
		AudioBitrate:         "128k",
		Resolution:           ffmpeg.Resolution720p,
		FrameRate:            30,
		Parallel:             true,
		MaxConcurrent:        4,
		ProcessingTimeout:    30 * time.Minute,
		MemoryLimit:          2 * 1024 * 1024 * 1024, // 2GB
		EnableMetrics:        true,
		MetricsInterval:      10 * time.Second,
		HealthCheckInterval:  30 * time.Second,
		StreamHandler:        &DefaultStreamHandler{},
		LogLevel:             "info",
		CleanupTemp:          true,
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.FFmpeg == nil {
		return &RTSPError{
			Message: "FFmpeg instance is required",
			Code:    ErrCodeInvalidConfig,
		}
	}

	if c.HLS == nil {
		return &RTSPError{
			Message: "HLS instance is required",
			Code:    ErrCodeInvalidConfig,
		}
	}

	if c.OutputDir == "" {
		return &RTSPError{
			Message: "output directory is required",
			Code:    ErrCodeInvalidConfig,
		}
	}

	if c.Layout.Rows <= 0 || c.Layout.Columns <= 0 {
		return &RTSPError{
			Message: "layout rows and columns must be positive",
			Code:    ErrCodeLayoutError,
		}
	}

	if c.Layout.Width <= 0 || c.Layout.Height <= 0 {
		return &RTSPError{
			Message: "layout width and height must be positive",
			Code:    ErrCodeLayoutError,
		}
	}

	if c.ConnectionTimeout <= 0 {
		c.ConnectionTimeout = 30 * time.Second
	}

	if c.ReadTimeout <= 0 {
		c.ReadTimeout = 10 * time.Second
	}

	if c.MaxConcurrent <= 0 {
		c.MaxConcurrent = 1
	}

	if c.ProcessingTimeout <= 0 {
		c.ProcessingTimeout = 30 * time.Minute
	}

	if c.StreamHandler == nil {
		c.StreamHandler = &DefaultStreamHandler{}
	}

	return nil
}

// SetupOutputDir creates the output directory structure
func (c *Config) SetupOutputDir() error {
	if err := os.MkdirAll(c.OutputDir, 0755); err != nil {
		return &RTSPError{
			Message: "failed to create output directory",
			Code:    ErrCodeInvalidConfig,
			Cause:   err,
		}
	}

	// Create subdirectories for different streaming modes
	if c.StreamingMode == ModeSeparate || c.StreamingMode == ModeBoth {
		streamsDir := filepath.Join(c.OutputDir, "streams")
		if err := os.MkdirAll(streamsDir, 0755); err != nil {
			return &RTSPError{
				Message: "failed to create streams directory",
				Code:    ErrCodeInvalidConfig,
				Cause:   err,
			}
		}
	}

	if c.StreamingMode == ModeMerged || c.StreamingMode == ModeBoth {
		mergedDir := filepath.Join(c.OutputDir, "merged")
		if err := os.MkdirAll(mergedDir, 0755); err != nil {
			return &RTSPError{
				Message: "failed to create merged directory",
				Code:    ErrCodeInvalidConfig,
				Cause:   err,
			}
		}
	}

	// Setup temp directory if specified
	if c.TempDir != "" {
		if err := os.MkdirAll(c.TempDir, 0755); err != nil {
			return &RTSPError{
				Message: "failed to create temp directory",
				Code:    ErrCodeInvalidConfig,
				Cause:   err,
			}
		}
	}

	return nil
}

// GetStreamOutputDir returns the output directory for a specific stream
func (c *Config) GetStreamOutputDir(streamName string) string {
	if c.StreamingMode == ModeSeparate || c.StreamingMode == ModeBoth {
		return filepath.Join(c.OutputDir, "streams", streamName)
	}
	return c.OutputDir
}

// GetMergedOutputDir returns the output directory for merged streams
func (c *Config) GetMergedOutputDir() string {
	if c.StreamingMode == ModeMerged || c.StreamingMode == ModeBoth {
		return filepath.Join(c.OutputDir, "merged")
	}
	return c.OutputDir
}

// GetHLSConfig returns HLS configuration for conversion
func (c *Config) GetHLSConfig() *hls.Config {
	if c.HLSConfig != nil {
		return c.HLSConfig
	}

	// Create default HLS config based on RTSP config
	hlsConfig := hls.DefaultConfig()
	hlsConfig.FFmpeg = c.FFmpeg
	hlsConfig.SegmentOptions.Duration = c.SegmentDuration
	hlsConfig.SegmentOptions.ListSize = c.PlaylistSize
	hlsConfig.SegmentOptions.DeleteOld = c.DeleteOldSegments
	hlsConfig.Parallel = c.Parallel
	hlsConfig.MaxConcurrent = c.MaxConcurrent
	hlsConfig.Timeout = c.ProcessingTimeout
	hlsConfig.TempDir = c.TempDir
	hlsConfig.CleanupTemp = c.CleanupTemp

	// Set quality levels based on RTSP config
	qualityLevel := hls.QualityLevel{
		Name:         "default",
		Resolution:   c.Resolution,
		VideoBitrate: c.VideoBitrate,
		AudioBitrate: c.AudioBitrate,
		VideoCodec:   c.VideoCodec,
		AudioCodec:   c.AudioCodec,
		FrameRate:    c.FrameRate,
	}
	hlsConfig.QualityLevels = []hls.QualityLevel{qualityLevel}

	return hlsConfig
}

// AutoDetectLayout automatically detects the best layout for the given number of streams
func (c *Config) AutoDetectLayout(streamCount int) Layout {
	if !c.AutoLayout {
		return c.Layout
	}

	switch streamCount {
	case 1:
		return DefaultLayouts[LayoutSingle]
	case 2:
		return DefaultLayouts[Layout1x2]
	case 3, 4:
		return DefaultLayouts[Layout2x2]
	case 5, 6:
		return DefaultLayouts[Layout2x3]
	case 7, 8, 9:
		return DefaultLayouts[Layout3x3]
	case 10, 11, 12, 13, 14, 15, 16:
		return DefaultLayouts[Layout4x4]
	default:
		// For more than 16 streams, use a custom layout
		rows := int(float64(streamCount)/4.0 + 0.5) // Round up
		cols := 4
		if rows*cols < streamCount {
			rows++
		}
		return Layout{
			Type:    LayoutCustom,
			Rows:    rows,
			Columns: cols,
			Width:   c.Layout.Width,
			Height:  c.Layout.Height,
		}
	}
}

// ValidateStreamURL validates an RTSP URL
func (c *Config) ValidateStreamURL(streamURL string) error {
	if streamURL == "" {
		return &RTSPError{
			Message: "stream URL cannot be empty",
			Code:    ErrCodeInvalidURL,
		}
	}

	parsedURL, err := url.Parse(streamURL)
	if err != nil {
		return &RTSPError{
			Message: "invalid stream URL format",
			Code:    ErrCodeInvalidURL,
			Cause:   err,
		}
	}

	if parsedURL.Scheme != "rtsp" {
		return &RTSPError{
			Message: "URL must use rtsp:// scheme",
			Code:    ErrCodeInvalidURL,
		}
	}

	if parsedURL.Host == "" {
		return &RTSPError{
			Message: "URL must specify a host",
			Code:    ErrCodeInvalidURL,
		}
	}

	return nil
}

// Clone creates a deep copy of the configuration
func (c *Config) Clone() *Config {
	clone := *c

	// Deep copy HLS config if present
	if c.HLSConfig != nil {
		hlsConfig := c.HLSConfig.Clone()
		clone.HLSConfig = hlsConfig
	}

	// Deep copy slices
	if c.FFmpegArgs != nil {
		clone.FFmpegArgs = make([]string, len(c.FFmpegArgs))
		copy(clone.FFmpegArgs, c.FFmpegArgs)
	}

	if c.CustomFilters != nil {
		clone.CustomFilters = make([]string, len(c.CustomFilters))
		copy(clone.CustomFilters, c.CustomFilters)
	}

	return &clone
}

// ConfigBuilder provides a fluent interface for building RTSP configurations
type ConfigBuilder struct {
	config *Config
}

// NewConfigBuilder creates a new configuration builder
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: DefaultConfig(),
	}
}

// WithFFmpeg sets the FFmpeg instance
func (b *ConfigBuilder) WithFFmpeg(ffmpegInstance ffmpeg.FFmpeg) *ConfigBuilder {
	b.config.FFmpeg = ffmpegInstance
	return b
}

// WithHLS sets the HLS instance
func (b *ConfigBuilder) WithHLS(hlsInstance hls.HLS) *ConfigBuilder {
	b.config.HLS = hlsInstance
	return b
}

// WithOutputDir sets the output directory
func (b *ConfigBuilder) WithOutputDir(dir string) *ConfigBuilder {
	b.config.OutputDir = dir
	return b
}

// WithStreamingMode sets the streaming mode
func (b *ConfigBuilder) WithStreamingMode(mode StreamingMode) *ConfigBuilder {
	b.config.StreamingMode = mode
	return b
}

// WithLayout sets the layout configuration
func (b *ConfigBuilder) WithLayout(layout Layout) *ConfigBuilder {
	b.config.Layout = layout
	b.config.AutoLayout = false
	return b
}

// WithAutoLayout enables automatic layout detection
func (b *ConfigBuilder) WithAutoLayout(enabled bool) *ConfigBuilder {
	b.config.AutoLayout = enabled
	return b
}

// WithTransport sets the default transport protocol
func (b *ConfigBuilder) WithTransport(transport TransportProtocol) *ConfigBuilder {
	b.config.DefaultTransport = transport
	return b
}

// WithTimeouts sets connection and read timeouts
func (b *ConfigBuilder) WithTimeouts(connection, read time.Duration) *ConfigBuilder {
	b.config.ConnectionTimeout = connection
	b.config.ReadTimeout = read
	return b
}

// WithReconnect configures reconnection settings
func (b *ConfigBuilder) WithReconnect(enabled bool, maxAttempts int, delay time.Duration) *ConfigBuilder {
	b.config.ReconnectEnabled = enabled
	b.config.MaxReconnectAttempts = maxAttempts
	b.config.ReconnectDelay = delay
	return b
}

// WithQuality sets video quality parameters
func (b *ConfigBuilder) WithQuality(resolution ffmpeg.Resolution, videoBitrate, audioBitrate string, fps float64) *ConfigBuilder {
	b.config.Resolution = resolution
	b.config.VideoBitrate = videoBitrate
	b.config.AudioBitrate = audioBitrate
	b.config.FrameRate = fps
	return b
}

// WithParallel enables parallel processing
func (b *ConfigBuilder) WithParallel(parallel bool, maxConcurrent int) *ConfigBuilder {
	b.config.Parallel = parallel
	b.config.MaxConcurrent = maxConcurrent
	return b
}

// WithStreamHandler sets the stream event handler
func (b *ConfigBuilder) WithStreamHandler(handler StreamHandler) *ConfigBuilder {
	b.config.StreamHandler = handler
	return b
}

// WithProgressCallback sets the progress callback function
func (b *ConfigBuilder) WithProgressCallback(callback func(ConversionProgress)) *ConfigBuilder {
	b.config.ProgressCallback = callback
	return b
}

// Build returns the built configuration
func (b *ConfigBuilder) Build() *Config {
	return b.config.Clone()
}
