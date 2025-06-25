package hls

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nguyendkn/go-libs/ffmpeg"
)

// Config represents HLS conversion configuration
type Config struct {
	// FFmpeg configuration
	FFmpeg ffmpeg.FFmpeg `json:"-"`

	// Output settings
	OutputDir       string        `json:"output_dir"`
	BaseURL         string        `json:"base_url,omitempty"`
	PlaylistName    string        `json:"playlist_name"`
	MasterPlaylist  string        `json:"master_playlist"`
	
	// Format settings
	Format          HLSFormat     `json:"format"`
	PlaylistType    PlaylistType  `json:"playlist_type"`
	
	// Segment settings
	SegmentOptions  SegmentOptions `json:"segment_options"`
	
	// Quality settings
	QualityLevels   []QualityLevel `json:"quality_levels"`
	AdaptiveBitrate bool           `json:"adaptive_bitrate"`
	
	// Encryption settings
	Encryption      *EncryptionOptions `json:"encryption,omitempty"`
	
	// Advanced settings
	Parallel        bool          `json:"parallel"`
	MaxConcurrent   int           `json:"max_concurrent"`
	Timeout         time.Duration `json:"timeout"`
	TempDir         string        `json:"temp_dir"`
	CleanupTemp     bool          `json:"cleanup_temp"`
	
	// Optimization settings
	FastStart       bool          `json:"fast_start"`
	TwoPass         bool          `json:"two_pass"`
	LookAhead       bool          `json:"look_ahead"`
	
	// Callback functions
	ProgressCallback func(ConversionProgress) `json:"-"`
	ErrorCallback    func(error)              `json:"-"`
}

// DefaultConfig returns a default HLS configuration
func DefaultConfig() *Config {
	return &Config{
		OutputDir:      "output",
		PlaylistName:   "playlist.m3u8",
		MasterPlaylist: "master.m3u8",
		Format:         FormatHLS,
		PlaylistType:   PlaylistVOD,
		SegmentOptions: DefaultSegmentOptions(),
		QualityLevels:  DefaultQualityLevels,
		AdaptiveBitrate: true,
		Parallel:       true,
		MaxConcurrent:  4,
		Timeout:        30 * time.Minute,
		CleanupTemp:    true,
		FastStart:      true,
		TwoPass:        false,
		LookAhead:      true,
	}
}

// DefaultSegmentOptions returns default segment options
func DefaultSegmentOptions() SegmentOptions {
	return SegmentOptions{
		Duration:    6 * time.Second,
		ListSize:    0, // Keep all segments
		Format:      SegmentTS,
		Pattern:     "segment_%03d.ts",
		StartNumber: 0,
		DeleteOld:   false,
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.FFmpeg == nil {
		return &HLSError{
			Message: "FFmpeg instance is required",
			Code:    ErrCodeInvalidConfig,
		}
	}

	if c.OutputDir == "" {
		return &HLSError{
			Message: "output directory is required",
			Code:    ErrCodeInvalidConfig,
		}
	}

	if c.PlaylistName == "" {
		return &HLSError{
			Message: "playlist name is required",
			Code:    ErrCodeInvalidConfig,
		}
	}

	if len(c.QualityLevels) == 0 {
		return &HLSError{
			Message: "at least one quality level is required",
			Code:    ErrCodeInvalidConfig,
		}
	}

	if c.SegmentOptions.Duration <= 0 {
		return &HLSError{
			Message: "segment duration must be positive",
			Code:    ErrCodeInvalidConfig,
		}
	}

	if c.MaxConcurrent <= 0 {
		c.MaxConcurrent = 1
	}

	if c.Timeout <= 0 {
		c.Timeout = 30 * time.Minute
	}

	// Validate quality levels
	for i, level := range c.QualityLevels {
		if level.Name == "" {
			return &HLSError{
				Message: fmt.Sprintf("quality level %d: name is required", i),
				Code:    ErrCodeInvalidConfig,
			}
		}
		if level.VideoBitrate == "" {
			return &HLSError{
				Message: fmt.Sprintf("quality level %s: video bitrate is required", level.Name),
				Code:    ErrCodeInvalidConfig,
			}
		}
	}

	return nil
}

// SetupOutputDir creates the output directory structure
func (c *Config) SetupOutputDir() error {
	if err := os.MkdirAll(c.OutputDir, 0755); err != nil {
		return &HLSError{
			Message: "failed to create output directory",
			Code:    ErrCodeFileSystem,
			Cause:   err,
		}
	}

	// Create subdirectories for each quality level if adaptive bitrate is enabled
	if c.AdaptiveBitrate {
		for _, level := range c.QualityLevels {
			qualityDir := filepath.Join(c.OutputDir, level.Name)
			if err := os.MkdirAll(qualityDir, 0755); err != nil {
				return &HLSError{
					Message: fmt.Sprintf("failed to create quality directory: %s", level.Name),
					Code:    ErrCodeFileSystem,
					Cause:   err,
				}
			}
		}
	}

	// Setup temp directory if specified
	if c.TempDir != "" {
		if err := os.MkdirAll(c.TempDir, 0755); err != nil {
			return &HLSError{
				Message: "failed to create temp directory",
				Code:    ErrCodeFileSystem,
				Cause:   err,
			}
		}
	}

	return nil
}

// GetQualityOutputDir returns the output directory for a specific quality level
func (c *Config) GetQualityOutputDir(qualityName string) string {
	if c.AdaptiveBitrate {
		return filepath.Join(c.OutputDir, qualityName)
	}
	return c.OutputDir
}

// GetPlaylistPath returns the full path to the playlist file for a quality level
func (c *Config) GetPlaylistPath(qualityName string) string {
	outputDir := c.GetQualityOutputDir(qualityName)
	if c.AdaptiveBitrate {
		return filepath.Join(outputDir, c.PlaylistName)
	}
	return filepath.Join(outputDir, qualityName+"_"+c.PlaylistName)
}

// GetMasterPlaylistPath returns the full path to the master playlist
func (c *Config) GetMasterPlaylistPath() string {
	return filepath.Join(c.OutputDir, c.MasterPlaylist)
}

// GetSegmentPattern returns the segment filename pattern for a quality level
func (c *Config) GetSegmentPattern(qualityName string) string {
	if c.SegmentOptions.Pattern == "" {
		return "segment_%03d.ts"
	}
	
	// Add quality prefix if adaptive bitrate
	if c.AdaptiveBitrate {
		return c.SegmentOptions.Pattern
	}
	
	// Add quality prefix to pattern
	ext := filepath.Ext(c.SegmentOptions.Pattern)
	base := c.SegmentOptions.Pattern[:len(c.SegmentOptions.Pattern)-len(ext)]
	return qualityName + "_" + base + ext
}

// Clone creates a deep copy of the configuration
func (c *Config) Clone() *Config {
	clone := *c
	
	// Deep copy quality levels
	clone.QualityLevels = make([]QualityLevel, len(c.QualityLevels))
	copy(clone.QualityLevels, c.QualityLevels)
	
	// Deep copy encryption options if present
	if c.Encryption != nil {
		encryption := *c.Encryption
		clone.Encryption = &encryption
	}
	
	return &clone
}

// ConfigBuilder provides a fluent interface for building HLS configurations
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

// WithOutputDir sets the output directory
func (b *ConfigBuilder) WithOutputDir(dir string) *ConfigBuilder {
	b.config.OutputDir = dir
	return b
}

// WithQualityLevels sets the quality levels
func (b *ConfigBuilder) WithQualityLevels(levels ...QualityLevel) *ConfigBuilder {
	b.config.QualityLevels = levels
	return b
}

// WithSegmentDuration sets the segment duration
func (b *ConfigBuilder) WithSegmentDuration(duration time.Duration) *ConfigBuilder {
	b.config.SegmentOptions.Duration = duration
	return b
}

// WithPlaylistType sets the playlist type
func (b *ConfigBuilder) WithPlaylistType(playlistType PlaylistType) *ConfigBuilder {
	b.config.PlaylistType = playlistType
	return b
}

// WithEncryption enables encryption with the specified options
func (b *ConfigBuilder) WithEncryption(encryption *EncryptionOptions) *ConfigBuilder {
	b.config.Encryption = encryption
	return b
}

// WithParallel enables or disables parallel processing
func (b *ConfigBuilder) WithParallel(parallel bool, maxConcurrent int) *ConfigBuilder {
	b.config.Parallel = parallel
	b.config.MaxConcurrent = maxConcurrent
	return b
}

// WithProgressCallback sets the progress callback function
func (b *ConfigBuilder) WithProgressCallback(callback func(ConversionProgress)) *ConfigBuilder {
	b.config.ProgressCallback = callback
	return b
}

// WithTimeout sets the conversion timeout
func (b *ConfigBuilder) WithTimeout(timeout time.Duration) *ConfigBuilder {
	b.config.Timeout = timeout
	return b
}

// Build returns the built configuration
func (b *ConfigBuilder) Build() *Config {
	return b.config.Clone()
}
