package ffmpeg

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// FFmpeg represents the main FFmpeg interface
type FFmpeg interface {
	// Core methods
	GetVersion() (string, error)
	GetSupportedFormats() ([]string, error)
	GetSupportedCodecs() (map[string][]string, error)
	ValidateInstallation() error
	
	// Builder methods
	New() *Builder
	NewBuilder() *Builder
	
	// Direct execution methods
	Execute(ctx context.Context, args []string, opts *ExecuteOptions) error
	ExecuteWithProgress(ctx context.Context, args []string, progressHandler func(ProgressInfo)) error
	
	// Utility methods
	GetMediaInfo(filePath string) (*MediaInfo, error)
	ProbeFile(filePath string) (*MediaInfo, error)
	
	// Extension methods
	ConvertVideo(input, output string, opts *ConversionOptions) error
	ExtractAudio(input, output string, opts *ConversionOptions) error
	GenerateThumbnail(input, output string, opts *ThumbnailOptions) error
	CompressVideo(input, output string, targetSize int64) error
}

// ffmpegImpl implements the FFmpeg interface
type ffmpegImpl struct {
	config *Config
	mutex  sync.RWMutex
}

// New creates a new FFmpeg instance with auto-detected configuration
func New() (FFmpeg, error) {
	config, err := SetupConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to setup ffmpeg config: %w", err)
	}
	
	return &ffmpegImpl{
		config: config,
	}, nil
}

// NewWithConfig creates a new FFmpeg instance with custom configuration
func NewWithConfig(config *Config) (FFmpeg, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	
	if err := ValidateFFmpeg(config.BinaryPath); err != nil {
		return nil, fmt.Errorf("invalid ffmpeg configuration: %w", err)
	}
	
	return &ffmpegImpl{
		config: config,
	}, nil
}

// GetVersion returns FFmpeg version information
func (f *ffmpegImpl) GetVersion() (string, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	return GetFFmpegVersion(f.config.BinaryPath)
}

// GetSupportedFormats returns list of supported formats
func (f *ffmpegImpl) GetSupportedFormats() ([]string, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	return GetSupportedFormats(f.config.BinaryPath)
}

// GetSupportedCodecs returns map of supported codecs by type
func (f *ffmpegImpl) GetSupportedCodecs() (map[string][]string, error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	return GetSupportedCodecs(f.config.BinaryPath)
}

// ValidateInstallation checks if FFmpeg is properly installed and accessible
func (f *ffmpegImpl) ValidateInstallation() error {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	
	return ValidateFFmpeg(f.config.BinaryPath)
}

// New creates a new command builder
func (f *ffmpegImpl) New() *Builder {
	return NewBuilder(f.config)
}

// NewBuilder creates a new command builder (alias for New)
func (f *ffmpegImpl) NewBuilder() *Builder {
	return NewBuilder(f.config)
}

// Execute runs FFmpeg with the given arguments
func (f *ffmpegImpl) Execute(ctx context.Context, args []string, opts *ExecuteOptions) error {
	f.mutex.RLock()
	binaryPath := f.config.BinaryPath
	f.mutex.RUnlock()
	
	executor := NewExecutor(binaryPath)
	return executor.Execute(ctx, args, opts)
}

// ExecuteWithProgress runs FFmpeg with progress tracking
func (f *ffmpegImpl) ExecuteWithProgress(ctx context.Context, args []string, progressHandler func(ProgressInfo)) error {
	opts := &ExecuteOptions{
		Context:         ctx,
		ProgressHandler: progressHandler,
		Timeout:         time.Duration(f.config.Timeout) * time.Second,
	}
	
	return f.Execute(ctx, args, opts)
}

// GetMediaInfo extracts metadata from media file
func (f *ffmpegImpl) GetMediaInfo(filePath string) (*MediaInfo, error) {
	f.mutex.RLock()
	binaryPath := f.config.BinaryPath
	f.mutex.RUnlock()
	
	return ProbeMediaFile(binaryPath, filePath)
}

// ProbeFile is an alias for GetMediaInfo
func (f *ffmpegImpl) ProbeFile(filePath string) (*MediaInfo, error) {
	return f.GetMediaInfo(filePath)
}

// ConvertVideo converts video with specified options
func (f *ffmpegImpl) ConvertVideo(input, output string, opts *ConversionOptions) error {
	builder := f.New().Input(input).Output(output)
	
	if opts != nil {
		if opts.VideoCodec != "" {
			builder = builder.VideoCodec(opts.VideoCodec)
		}
		if opts.AudioCodec != "" {
			builder = builder.AudioCodec(opts.AudioCodec)
		}
		if opts.Quality != "" {
			builder = builder.Quality(opts.Quality)
		}
		if opts.Resolution != "" {
			builder = builder.Resolution(opts.Resolution)
		}
		if opts.VideoBitrate != "" {
			builder = builder.VideoBitrate(opts.VideoBitrate)
		}
		if opts.AudioBitrate != "" {
			builder = builder.AudioBitrate(opts.AudioBitrate)
		}
		if opts.FrameRate > 0 {
			builder = builder.FrameRate(opts.FrameRate)
		}
		if opts.StartTime > 0 {
			builder = builder.StartTime(opts.StartTime)
		}
		if opts.Duration > 0 {
			builder = builder.Duration(opts.Duration)
		}
		if len(opts.CustomArgs) > 0 {
			builder = builder.CustomArgs(opts.CustomArgs...)
		}
	}
	
	ctx := context.Background()
	return builder.Execute(ctx)
}

// ExtractAudio extracts audio from video file
func (f *ffmpegImpl) ExtractAudio(input, output string, opts *ConversionOptions) error {
	builder := f.New().
		Input(input).
		Output(output).
		VideoCodec("").  // No video
		AudioCodec(AudioCodecAAC)
	
	if opts != nil {
		if opts.AudioCodec != "" {
			builder = builder.AudioCodec(opts.AudioCodec)
		}
		if opts.AudioBitrate != "" {
			builder = builder.AudioBitrate(opts.AudioBitrate)
		}
		if opts.SampleRate > 0 {
			builder = builder.SampleRate(opts.SampleRate)
		}
		if opts.Channels > 0 {
			builder = builder.Channels(opts.Channels)
		}
		if opts.StartTime > 0 {
			builder = builder.StartTime(opts.StartTime)
		}
		if opts.Duration > 0 {
			builder = builder.Duration(opts.Duration)
		}
	}
	
	ctx := context.Background()
	return builder.Execute(ctx)
}

// GenerateThumbnail generates thumbnail from video
func (f *ffmpegImpl) GenerateThumbnail(input, output string, opts *ThumbnailOptions) error {
	builder := f.New().Input(input).Output(output)
	
	if opts != nil {
		if opts.Time > 0 {
			builder = builder.StartTime(opts.Time)
		}
		if opts.Width > 0 && opts.Height > 0 {
			resolution := fmt.Sprintf("%dx%d", opts.Width, opts.Height)
			builder = builder.Resolution(Resolution(resolution))
		}
		if opts.Quality > 0 {
			builder = builder.CustomArgs("-q:v", fmt.Sprintf("%d", opts.Quality))
		}
		
		// Generate single frame
		builder = builder.CustomArgs("-vframes", "1")
	}
	
	ctx := context.Background()
	return builder.Execute(ctx)
}

// CompressVideo compresses video to target file size
func (f *ffmpegImpl) CompressVideo(input, output string, targetSize int64) error {
	// Get input file info first
	info, err := f.GetMediaInfo(input)
	if err != nil {
		return fmt.Errorf("failed to get media info: %w", err)
	}
	
	// Calculate target bitrate (rough estimation)
	targetBitrate := (targetSize * 8) / int64(info.Duration.Seconds()) // bits per second
	targetBitrate = targetBitrate * 90 / 100 // Leave 10% margin
	
	builder := f.New().
		Input(input).
		Output(output).
		VideoCodec(VideoCodecH264).
		AudioCodec(AudioCodecAAC).
		VideoBitrate(fmt.Sprintf("%dk", targetBitrate/1000)).
		CustomArgs("-pass", "1", "-f", "null")
	
	// Two-pass encoding for better quality
	ctx := context.Background()
	
	// First pass
	if err := builder.Execute(ctx); err != nil {
		return fmt.Errorf("first pass failed: %w", err)
	}
	
	// Second pass
	builder = f.New().
		Input(input).
		Output(output).
		VideoCodec(VideoCodecH264).
		AudioCodec(AudioCodecAAC).
		VideoBitrate(fmt.Sprintf("%dk", targetBitrate/1000)).
		CustomArgs("-pass", "2")
	
	return builder.Execute(ctx)
}
