package rtsp

import (
	"context"
	"fmt"

	"github.com/nguyendkn/go-libs/ffmpeg"
	"github.com/nguyendkn/go-libs/hls"
)

// RTSP represents the main RTSP interface
type RTSP interface {
	// Stream management
	AddStream(stream RTSPStream) error
	AddStreams(streams []RTSPStream) error
	AddStreamURL(url string) error
	AddStreamURLs(urls []string) error
	RemoveStream(streamName string) error
	
	// Stream control
	StartStream(streamName string) error
	StartAllStreams() error
	StopStream(streamName string) error
	StopAllStreams() error
	
	// Stream information
	GetStreamInfo(streamName string) (StreamInfo, error)
	GetAllStreamInfo() map[string]StreamInfo
	GetStreamNames() []string
	GetStreamCount() int
	
	// Conversion methods
	ConvertSingle(ctx context.Context, streamURL, outputDir string) (*ConversionResult, error)
	ConvertMultiple(ctx context.Context, streamURLs []string, outputDir string) (*ConversionResult, error)
	ConvertMerged(ctx context.Context, streamURLs []string, outputDir string) (*ConversionResult, error)
	Convert(ctx context.Context, streamURLs []string, outputDir string, mode StreamingMode) (*ConversionResult, error)
	
	// Layout management
	SetLayout(layout Layout) error
	GetLayout() Layout
	PreviewLayout(streamNames []string) string
	
	// Configuration
	GetConfig() *Config
	UpdateConfig(config *Config) error
	
	// Monitoring
	GetActiveConversions() map[string]ConversionProgress
	StopConversion(conversionID string) error
	StopAllConversions()
	
	// Cleanup
	Close() error
}

// rtspImpl implements the RTSP interface
type rtspImpl struct {
	config        *Config
	streamManager *StreamManager
	converter     *Converter
	layoutManager *LayoutManager
}

// New creates a new RTSP instance with default configuration
func New() (RTSP, error) {
	// Create FFmpeg instance
	ffmpegInstance, err := ffmpeg.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create FFmpeg instance: %w", err)
	}

	// Create HLS instance
	hlsInstance, err := hls.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create HLS instance: %w", err)
	}

	// Create default config
	config := DefaultConfig()
	config.FFmpeg = ffmpegInstance
	config.HLS = hlsInstance

	return NewWithConfig(config)
}

// NewWithConfig creates a new RTSP instance with custom configuration
func NewWithConfig(config *Config) (RTSP, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	streamManager, err := NewStreamManager(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream manager: %w", err)
	}

	converter, err := NewConverter(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create converter: %w", err)
	}

	layoutManager := NewLayoutManager(config)

	return &rtspImpl{
		config:        config,
		streamManager: streamManager,
		converter:     converter,
		layoutManager: layoutManager,
	}, nil
}

// NewWithDependencies creates a new RTSP instance with existing FFmpeg and HLS instances
func NewWithDependencies(ffmpegInstance ffmpeg.FFmpeg, hlsInstance hls.HLS) (RTSP, error) {
	config := DefaultConfig()
	config.FFmpeg = ffmpegInstance
	config.HLS = hlsInstance
	return NewWithConfig(config)
}

// Stream management methods

func (r *rtspImpl) AddStream(stream RTSPStream) error {
	return r.streamManager.AddStream(stream)
}

func (r *rtspImpl) AddStreams(streams []RTSPStream) error {
	return r.streamManager.AddStreams(streams)
}

func (r *rtspImpl) AddStreamURL(url string) error {
	return r.streamManager.AddStreamURLs([]string{url})
}

func (r *rtspImpl) AddStreamURLs(urls []string) error {
	return r.streamManager.AddStreamURLs(urls)
}

func (r *rtspImpl) RemoveStream(streamName string) error {
	return r.streamManager.RemoveStream(streamName)
}

// Stream control methods

func (r *rtspImpl) StartStream(streamName string) error {
	return r.streamManager.StartStream(streamName)
}

func (r *rtspImpl) StartAllStreams() error {
	return r.streamManager.StartAllStreams()
}

func (r *rtspImpl) StopStream(streamName string) error {
	return r.streamManager.StopStream(streamName)
}

func (r *rtspImpl) StopAllStreams() error {
	return r.streamManager.StopAllStreams()
}

// Stream information methods

func (r *rtspImpl) GetStreamInfo(streamName string) (StreamInfo, error) {
	return r.streamManager.GetStreamInfo(streamName)
}

func (r *rtspImpl) GetAllStreamInfo() map[string]StreamInfo {
	return r.streamManager.GetAllStreamInfo()
}

func (r *rtspImpl) GetStreamNames() []string {
	return r.streamManager.GetStreamNames()
}

func (r *rtspImpl) GetStreamCount() int {
	return r.streamManager.GetStreamCount()
}

// Conversion methods

func (r *rtspImpl) ConvertSingle(ctx context.Context, streamURL, outputDir string) (*ConversionResult, error) {
	return r.converter.ConvertSingle(ctx, streamURL, outputDir)
}

func (r *rtspImpl) ConvertMultiple(ctx context.Context, streamURLs []string, outputDir string) (*ConversionResult, error) {
	return r.converter.ConvertMultiple(ctx, streamURLs, outputDir)
}

func (r *rtspImpl) ConvertMerged(ctx context.Context, streamURLs []string, outputDir string) (*ConversionResult, error) {
	return r.converter.ConvertMerged(ctx, streamURLs, outputDir)
}

func (r *rtspImpl) Convert(ctx context.Context, streamURLs []string, outputDir string, mode StreamingMode) (*ConversionResult, error) {
	return r.converter.Convert(ctx, streamURLs, outputDir, mode)
}

// Layout management methods

func (r *rtspImpl) SetLayout(layout Layout) error {
	if err := r.layoutManager.ValidateLayout(layout, r.GetStreamCount()); err != nil {
		return err
	}
	r.config.Layout = layout
	r.config.AutoLayout = false
	return nil
}

func (r *rtspImpl) GetLayout() Layout {
	return r.config.Layout
}

func (r *rtspImpl) PreviewLayout(streamNames []string) string {
	layout := r.layoutManager.CalculateLayout(len(streamNames))
	return r.layoutManager.PreviewLayout(layout, streamNames)
}

// Configuration methods

func (r *rtspImpl) GetConfig() *Config {
	return r.config.Clone()
}

func (r *rtspImpl) UpdateConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return err
	}

	// Update internal components
	streamManager, err := NewStreamManager(config)
	if err != nil {
		return err
	}

	converter, err := NewConverter(config)
	if err != nil {
		return err
	}

	layoutManager := NewLayoutManager(config)

	// Close old components
	r.streamManager.Close()
	r.converter.Close()

	// Update with new components
	r.config = config
	r.streamManager = streamManager
	r.converter = converter
	r.layoutManager = layoutManager

	return nil
}

// Monitoring methods

func (r *rtspImpl) GetActiveConversions() map[string]ConversionProgress {
	return r.converter.GetActiveConversions()
}

func (r *rtspImpl) StopConversion(conversionID string) error {
	return r.converter.StopConversion(conversionID)
}

func (r *rtspImpl) StopAllConversions() {
	r.converter.StopAllConversions()
}

// Cleanup

func (r *rtspImpl) Close() error {
	r.StopAllStreams()
	r.StopAllConversions()
	if err := r.streamManager.Close(); err != nil {
		return err
	}
	return r.converter.Close()
}

// Builder provides a fluent interface for RTSP operations
type Builder struct {
	config *Config
}

// NewBuilder creates a new RTSP builder
func NewBuilder() *Builder {
	return &Builder{
		config: DefaultConfig(),
	}
}

// WithFFmpeg sets the FFmpeg instance
func (b *Builder) WithFFmpeg(ffmpegInstance ffmpeg.FFmpeg) *Builder {
	b.config.FFmpeg = ffmpegInstance
	return b
}

// WithHLS sets the HLS instance
func (b *Builder) WithHLS(hlsInstance hls.HLS) *Builder {
	b.config.HLS = hlsInstance
	return b
}

// WithOutputDir sets the output directory
func (b *Builder) WithOutputDir(dir string) *Builder {
	b.config.OutputDir = dir
	return b
}

// WithStreamingMode sets the streaming mode
func (b *Builder) WithStreamingMode(mode StreamingMode) *Builder {
	b.config.StreamingMode = mode
	return b
}

// WithLayout sets the layout configuration
func (b *Builder) WithLayout(layout Layout) *Builder {
	b.config.Layout = layout
	b.config.AutoLayout = false
	return b
}

// WithAutoLayout enables automatic layout detection
func (b *Builder) WithAutoLayout(enabled bool) *Builder {
	b.config.AutoLayout = enabled
	return b
}

// WithTransport sets the default transport protocol
func (b *Builder) WithTransport(transport TransportProtocol) *Builder {
	b.config.DefaultTransport = transport
	return b
}

// WithQuality sets video quality parameters
func (b *Builder) WithQuality(resolution ffmpeg.Resolution, videoBitrate, audioBitrate string, fps float64) *Builder {
	b.config.Resolution = resolution
	b.config.VideoBitrate = videoBitrate
	b.config.AudioBitrate = audioBitrate
	b.config.FrameRate = fps
	return b
}

// WithParallel enables parallel processing
func (b *Builder) WithParallel(parallel bool, maxConcurrent int) *Builder {
	b.config.Parallel = parallel
	b.config.MaxConcurrent = maxConcurrent
	return b
}

// WithStreamHandler sets the stream event handler
func (b *Builder) WithStreamHandler(handler StreamHandler) *Builder {
	b.config.StreamHandler = handler
	return b
}

// WithProgressCallback sets the progress callback function
func (b *Builder) WithProgressCallback(callback func(ConversionProgress)) *Builder {
	b.config.ProgressCallback = callback
	return b
}

// Build creates the RTSP instance
func (b *Builder) Build() (RTSP, error) {
	return NewWithConfig(b.config)
}

// Convenience functions

// ConvertSingleStream converts a single RTSP stream to HLS
func ConvertSingleStream(ctx context.Context, streamURL, outputDir string) (*ConversionResult, error) {
	rtsp, err := New()
	if err != nil {
		return nil, err
	}
	defer rtsp.Close()

	return rtsp.ConvertSingle(ctx, streamURL, outputDir)
}

// ConvertMultipleStreams converts multiple RTSP streams to separate HLS streams
func ConvertMultipleStreams(ctx context.Context, streamURLs []string, outputDir string) (*ConversionResult, error) {
	rtsp, err := New()
	if err != nil {
		return nil, err
	}
	defer rtsp.Close()

	return rtsp.ConvertMultiple(ctx, streamURLs, outputDir)
}

// ConvertMergedStreams converts multiple RTSP streams to a single merged HLS stream
func ConvertMergedStreams(ctx context.Context, streamURLs []string, outputDir string) (*ConversionResult, error) {
	rtsp, err := New()
	if err != nil {
		return nil, err
	}
	defer rtsp.Close()

	return rtsp.ConvertMerged(ctx, streamURLs, outputDir)
}

// ConvertWithLayout converts multiple RTSP streams with a specific layout
func ConvertWithLayout(ctx context.Context, streamURLs []string, outputDir string, layout Layout) (*ConversionResult, error) {
	rtsp, err := NewBuilder().
		WithOutputDir(outputDir).
		WithLayout(layout).
		WithStreamingMode(ModeMerged).
		Build()
	if err != nil {
		return nil, err
	}
	defer rtsp.Close()

	return rtsp.ConvertMerged(ctx, streamURLs, outputDir)
}
