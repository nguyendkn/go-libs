package hls

import (
	"context"
	"fmt"
	"time"

	"github.com/nguyendkn/go-libs/ffmpeg"
)

// HLS represents the main HLS interface
type HLS interface {
	// Core conversion methods
	Convert(ctx context.Context, inputFile string) (*ConversionResult, error)
	ConvertWithOptions(ctx context.Context, inputFile string, options *ConversionOptions) (*ConversionResult, error)

	// Live streaming methods
	ConvertLive(ctx context.Context, inputSource string) error

	// Adaptive streaming methods
	GenerateAdaptiveStream(ctx context.Context, inputFile string) (*ConversionResult, error)
	CreateCustomAdaptiveStream(ctx context.Context, inputFile string, qualityLevels []QualityLevel) (*ConversionResult, error)
	GeneratePresetAdaptiveStream(ctx context.Context, inputFile string, preset AdaptivePreset) (*ConversionResult, error)

	// Analysis methods
	AnalyzeInput(inputFile string) (*ffmpeg.MediaInfo, error)
	AnalyzeOptimalLevels(inputFile string) ([]QualityLevel, error)
	GetBandwidthLadder(levels []QualityLevel) []BandwidthLevel

	// Configuration methods
	GetConfig() *Config
	UpdateConfig(config *Config) error

	// Utility methods
	Cleanup() error
	ValidateInput(inputFile string) error
}

// hlsImpl implements the HLS interface
type hlsImpl struct {
	converter         *Converter
	adaptiveStreaming *AdaptiveStreaming
	config            *Config
}

// New creates a new HLS instance with default configuration
func New() (HLS, error) {
	// Create FFmpeg instance
	ffmpegInstance, err := ffmpeg.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create FFmpeg instance: %w", err)
	}

	// Create default config
	config := DefaultConfig()
	config.FFmpeg = ffmpegInstance

	return NewWithConfig(config)
}

// NewWithConfig creates a new HLS instance with custom configuration
func NewWithConfig(config *Config) (HLS, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	converter, err := NewConverter(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create converter: %w", err)
	}

	adaptiveStreaming, err := NewAdaptiveStreaming(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create adaptive streaming: %w", err)
	}

	return &hlsImpl{
		converter:         converter,
		adaptiveStreaming: adaptiveStreaming,
		config:            config,
	}, nil
}

// NewWithFFmpeg creates a new HLS instance with existing FFmpeg instance
func NewWithFFmpeg(ffmpegInstance ffmpeg.FFmpeg) (HLS, error) {
	config := DefaultConfig()
	config.FFmpeg = ffmpegInstance
	return NewWithConfig(config)
}

// Convert converts a video file to HLS format
func (h *hlsImpl) Convert(ctx context.Context, inputFile string) (*ConversionResult, error) {
	return h.converter.Convert(ctx, inputFile)
}

// ConvertWithOptions converts video with custom options
func (h *hlsImpl) ConvertWithOptions(ctx context.Context, inputFile string, options *ConversionOptions) (*ConversionResult, error) {
	return h.converter.ConvertWithOptions(ctx, inputFile, options)
}

// ConvertLive converts video for live streaming
func (h *hlsImpl) ConvertLive(ctx context.Context, inputSource string) error {
	return h.converter.ConvertLive(ctx, inputSource)
}

// GenerateAdaptiveStream creates multiple quality levels for adaptive streaming
func (h *hlsImpl) GenerateAdaptiveStream(ctx context.Context, inputFile string) (*ConversionResult, error) {
	return h.adaptiveStreaming.GenerateAdaptiveStream(ctx, inputFile)
}

// CreateCustomAdaptiveStream creates adaptive stream with custom quality levels
func (h *hlsImpl) CreateCustomAdaptiveStream(ctx context.Context, inputFile string, qualityLevels []QualityLevel) (*ConversionResult, error) {
	return h.adaptiveStreaming.CreateCustomAdaptiveStream(ctx, inputFile, qualityLevels)
}

// GeneratePresetAdaptiveStream generates adaptive stream using predefined presets
func (h *hlsImpl) GeneratePresetAdaptiveStream(ctx context.Context, inputFile string, preset AdaptivePreset) (*ConversionResult, error) {
	return h.adaptiveStreaming.GeneratePresetAdaptiveStream(ctx, inputFile, preset)
}

// AnalyzeInput analyzes input file and returns media information
func (h *hlsImpl) AnalyzeInput(inputFile string) (*ffmpeg.MediaInfo, error) {
	return h.config.FFmpeg.GetMediaInfo(inputFile)
}

// AnalyzeOptimalLevels analyzes input and suggests optimal quality levels
func (h *hlsImpl) AnalyzeOptimalLevels(inputFile string) ([]QualityLevel, error) {
	return h.adaptiveStreaming.AnalyzeOptimalLevels(inputFile)
}

// GetBandwidthLadder returns the bandwidth ladder for quality levels
func (h *hlsImpl) GetBandwidthLadder(levels []QualityLevel) []BandwidthLevel {
	return h.adaptiveStreaming.GetBandwidthLadder(levels)
}

// GetConfig returns the current configuration
func (h *hlsImpl) GetConfig() *Config {
	return h.converter.GetConfig()
}

// UpdateConfig updates the HLS configuration
func (h *hlsImpl) UpdateConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return err
	}

	// Update converter
	if err := h.converter.UpdateConfig(config); err != nil {
		return err
	}

	// Recreate adaptive streaming with new config
	adaptiveStreaming, err := NewAdaptiveStreaming(config)
	if err != nil {
		return err
	}

	h.adaptiveStreaming = adaptiveStreaming
	h.config = config

	return nil
}

// Cleanup removes temporary files and directories
func (h *hlsImpl) Cleanup() error {
	return h.converter.Cleanup()
}

// ValidateInput validates the input file
func (h *hlsImpl) ValidateInput(inputFile string) error {
	if inputFile == "" {
		return &HLSError{
			Message: "input file path is required",
			Code:    ErrCodeInvalidInput,
		}
	}

	// Use FFmpeg to validate the file
	_, err := h.config.FFmpeg.GetMediaInfo(inputFile)
	if err != nil {
		return &HLSError{
			Message: fmt.Sprintf("invalid input file: %s", inputFile),
			Code:    ErrCodeInvalidInput,
			Cause:   err,
		}
	}

	return nil
}

// Builder provides a fluent interface for HLS operations
type Builder struct {
	config *Config
}

// NewBuilder creates a new HLS builder
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

// WithOutputDir sets the output directory
func (b *Builder) WithOutputDir(dir string) *Builder {
	b.config.OutputDir = dir
	return b
}

// WithQualityLevels sets the quality levels
func (b *Builder) WithQualityLevels(levels ...QualityLevel) *Builder {
	b.config.QualityLevels = levels
	return b
}

// WithSegmentDuration sets the segment duration
func (b *Builder) WithSegmentDuration(duration string) *Builder {
	if d, err := parseSegmentDuration(duration); err == nil {
		b.config.SegmentOptions.Duration = d
	}
	return b
}

// WithPlaylistType sets the playlist type
func (b *Builder) WithPlaylistType(playlistType PlaylistType) *Builder {
	b.config.PlaylistType = playlistType
	return b
}

// WithEncryption enables encryption
func (b *Builder) WithEncryption(encryption *EncryptionOptions) *Builder {
	b.config.Encryption = encryption
	return b
}

// WithParallel enables parallel processing
func (b *Builder) WithParallel(parallel bool, maxConcurrent int) *Builder {
	b.config.Parallel = parallel
	if maxConcurrent > 0 {
		b.config.MaxConcurrent = maxConcurrent
	}
	return b
}

// WithProgressCallback sets the progress callback
func (b *Builder) WithProgressCallback(callback func(ConversionProgress)) *Builder {
	b.config.ProgressCallback = callback
	return b
}

// WithPreset applies a predefined configuration preset
func (b *Builder) WithPreset(preset ConfigPreset) *Builder {
	switch preset {
	case PresetFast:
		b.config.TwoPass = false
		b.config.Parallel = true
		b.config.MaxConcurrent = 4
		b.config.QualityLevels = []QualityLevel{QualityMedium}
	case PresetQuality:
		b.config.TwoPass = true
		b.config.Parallel = false
		b.config.QualityLevels = []QualityLevel{QualityHigh}
	case PresetBalanced:
		b.config.TwoPass = false
		b.config.Parallel = true
		b.config.MaxConcurrent = 2
		b.config.QualityLevels = DefaultQualityLevels
	case PresetLive:
		b.config.PlaylistType = PlaylistLive
		b.config.TwoPass = false
		b.config.Parallel = true
		b.config.SegmentOptions.Duration = parseSegmentDurationDefault("2s")
		b.config.SegmentOptions.ListSize = 10
	}
	return b
}

// Build creates the HLS instance
func (b *Builder) Build() (HLS, error) {
	return NewWithConfig(b.config)
}

// ConfigPreset represents predefined configuration presets
type ConfigPreset string

const (
	PresetFast     ConfigPreset = "fast"     // Fast conversion, lower quality
	PresetQuality  ConfigPreset = "quality"  // High quality, slower conversion
	PresetBalanced ConfigPreset = "balanced" // Balanced speed and quality
	PresetLive     ConfigPreset = "live"     // Live streaming optimized
)

// Convenience functions

// ConvertToHLS converts a video file to HLS with default settings
func ConvertToHLS(ctx context.Context, inputFile, outputDir string) (*ConversionResult, error) {
	hls, err := NewBuilder().
		WithOutputDir(outputDir).
		Build()
	if err != nil {
		return nil, err
	}
	defer hls.Cleanup()

	return hls.Convert(ctx, inputFile)
}

// ConvertToAdaptiveHLS converts a video file to adaptive HLS
func ConvertToAdaptiveHLS(ctx context.Context, inputFile, outputDir string, preset AdaptivePreset) (*ConversionResult, error) {
	hls, err := NewBuilder().
		WithOutputDir(outputDir).
		Build()
	if err != nil {
		return nil, err
	}
	defer hls.Cleanup()

	return hls.GeneratePresetAdaptiveStream(ctx, inputFile, preset)
}

// ConvertWithPreset converts a video file using a predefined preset
func ConvertWithPreset(ctx context.Context, inputFile, outputDir string, preset ConfigPreset) (*ConversionResult, error) {
	hls, err := NewBuilder().
		WithOutputDir(outputDir).
		WithPreset(preset).
		Build()
	if err != nil {
		return nil, err
	}
	defer hls.Cleanup()

	return hls.Convert(ctx, inputFile)
}

// Helper functions

// parseSegmentDuration parses segment duration string
func parseSegmentDuration(duration string) (time.Duration, error) {
	// Implementation would parse duration strings like "6s", "10s", etc.
	// For now, return a default
	return parseSegmentDurationDefault(duration), nil
}

// parseSegmentDurationDefault returns default duration
func parseSegmentDurationDefault(duration string) time.Duration {
	// Simple implementation - in real code, you'd parse the string properly
	switch duration {
	case "2s":
		return 2 * time.Second
	case "4s":
		return 4 * time.Second
	case "6s":
		return 6 * time.Second
	case "10s":
		return 10 * time.Second
	default:
		return 6 * time.Second
	}
}
