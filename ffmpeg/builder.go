package ffmpeg

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Builder provides a fluent API for building FFmpeg commands
type Builder struct {
	config       *Config
	inputFiles   []string
	outputFile   string
	videoCodec   VideoCodec
	audioCodec   AudioCodec
	quality      Quality
	resolution   Resolution
	videoBitrate string
	audioBitrate string
	frameRate    float64
	sampleRate   int
	channels     int
	startTime    time.Duration
	duration     time.Duration
	filters      FilterChain
	customArgs   []string
	overwrite    bool
}

// NewBuilder creates a new command builder
func NewBuilder(config *Config) *Builder {
	return &Builder{
		config:    config,
		overwrite: true, // Default to overwrite output files
	}
}

// Input adds input file(s)
func (b *Builder) Input(files ...string) *Builder {
	b.inputFiles = append(b.inputFiles, files...)
	return b
}

// Output sets the output file
func (b *Builder) Output(file string) *Builder {
	b.outputFile = file
	return b
}

// VideoCodec sets the video codec
func (b *Builder) VideoCodec(codec VideoCodec) *Builder {
	b.videoCodec = codec
	return b
}

// AudioCodec sets the audio codec
func (b *Builder) AudioCodec(codec AudioCodec) *Builder {
	b.audioCodec = codec
	return b
}

// Quality sets the encoding quality/preset
func (b *Builder) Quality(quality Quality) *Builder {
	b.quality = quality
	return b
}

// Resolution sets the output resolution
func (b *Builder) Resolution(resolution Resolution) *Builder {
	b.resolution = resolution
	return b
}

// VideoBitrate sets the video bitrate (e.g., "1000k", "2M")
func (b *Builder) VideoBitrate(bitrate string) *Builder {
	b.videoBitrate = bitrate
	return b
}

// AudioBitrate sets the audio bitrate (e.g., "128k", "320k")
func (b *Builder) AudioBitrate(bitrate string) *Builder {
	b.audioBitrate = bitrate
	return b
}

// FrameRate sets the output frame rate
func (b *Builder) FrameRate(fps float64) *Builder {
	b.frameRate = fps
	return b
}

// SampleRate sets the audio sample rate
func (b *Builder) SampleRate(rate int) *Builder {
	b.sampleRate = rate
	return b
}

// Channels sets the number of audio channels
func (b *Builder) Channels(channels int) *Builder {
	b.channels = channels
	return b
}

// StartTime sets the start time for trimming
func (b *Builder) StartTime(t time.Duration) *Builder {
	b.startTime = t
	return b
}

// Duration sets the duration for trimming
func (b *Builder) Duration(d time.Duration) *Builder {
	b.duration = d
	return b
}

// VideoFilter adds video filter(s)
func (b *Builder) VideoFilter(filters ...string) *Builder {
	b.filters.Video = append(b.filters.Video, filters...)
	return b
}

// AudioFilter adds audio filter(s)
func (b *Builder) AudioFilter(filters ...string) *Builder {
	b.filters.Audio = append(b.filters.Audio, filters...)
	return b
}

// Scale adds a scale filter (shorthand for common video filter)
func (b *Builder) Scale(width, height int) *Builder {
	filter := fmt.Sprintf("scale=%d:%d", width, height)
	return b.VideoFilter(filter)
}

// Crop adds a crop filter
func (b *Builder) Crop(width, height, x, y int) *Builder {
	filter := fmt.Sprintf("crop=%d:%d:%d:%d", width, height, x, y)
	return b.VideoFilter(filter)
}

// Rotate adds a rotation filter
func (b *Builder) Rotate(degrees float64) *Builder {
	radians := degrees * 3.14159 / 180
	filter := fmt.Sprintf("rotate=%f", radians)
	return b.VideoFilter(filter)
}

// Volume adjusts audio volume (1.0 = normal, 0.5 = half, 2.0 = double)
func (b *Builder) Volume(factor float64) *Builder {
	filter := fmt.Sprintf("volume=%f", factor)
	return b.AudioFilter(filter)
}

// CustomArgs adds custom FFmpeg arguments
func (b *Builder) CustomArgs(args ...string) *Builder {
	b.customArgs = append(b.customArgs, args...)
	return b
}

// Overwrite sets whether to overwrite output files
func (b *Builder) Overwrite(overwrite bool) *Builder {
	b.overwrite = overwrite
	return b
}

// Build constructs the FFmpeg command arguments
func (b *Builder) Build() ([]string, error) {
	if len(b.inputFiles) == 0 {
		return nil, fmt.Errorf("no input files specified")
	}

	if b.outputFile == "" {
		return nil, fmt.Errorf("no output file specified")
	}

	var args []string

	// Global options
	args = append(args, "-loglevel", b.config.LogLevel)

	// Overwrite option
	if b.overwrite {
		args = append(args, "-y")
	} else {
		args = append(args, "-n")
	}

	// Input files
	for _, input := range b.inputFiles {
		args = append(args, "-i", input)
	}

	// Start time (seek)
	if b.startTime > 0 {
		args = append(args, "-ss", formatDuration(b.startTime))
	}

	// Duration
	if b.duration > 0 {
		args = append(args, "-t", formatDuration(b.duration))
	}

	// Video codec
	if b.videoCodec != "" {
		if b.videoCodec == VideoCodecCopy {
			args = append(args, "-c:v", "copy")
		} else {
			args = append(args, "-c:v", string(b.videoCodec))
		}
	}

	// Audio codec
	if b.audioCodec != "" {
		if b.audioCodec == AudioCodecCopy {
			args = append(args, "-c:a", "copy")
		} else {
			args = append(args, "-c:a", string(b.audioCodec))
		}
	}

	// Quality/preset
	if b.quality != "" {
		if b.videoCodec == VideoCodecH264 || b.videoCodec == VideoCodecH265 {
			args = append(args, "-preset", string(b.quality))
		} else {
			// For other codecs, use CRF values based on quality
			var crf string
			switch b.quality {
			case QualityHigh:
				crf = "18"
			case QualityLow:
				crf = "28"
			case QualityUltraFast, QualitySuperFast, QualityVeryFast, QualityFaster, QualityFast, QualityMedium, QualitySlow, QualitySlower, QualityVerySlow:
				// For preset-style qualities on non-H264/H265 codecs, use default CRF
				crf = "23"
			default:
				crf = "23"
			}
			args = append(args, "-crf", crf)
		}
	}

	// Resolution
	if b.resolution != "" {
		args = append(args, "-s", string(b.resolution))
	}

	// Video bitrate
	if b.videoBitrate != "" {
		args = append(args, "-b:v", b.videoBitrate)
	}

	// Audio bitrate
	if b.audioBitrate != "" {
		args = append(args, "-b:a", b.audioBitrate)
	}

	// Frame rate
	if b.frameRate > 0 {
		args = append(args, "-r", strconv.FormatFloat(b.frameRate, 'f', -1, 64))
	}

	// Sample rate
	if b.sampleRate > 0 {
		args = append(args, "-ar", strconv.Itoa(b.sampleRate))
	}

	// Channels
	if b.channels > 0 {
		args = append(args, "-ac", strconv.Itoa(b.channels))
	}

	// Filters
	if len(b.filters.Video) > 0 {
		videoFilter := strings.Join(b.filters.Video, ",")
		args = append(args, "-vf", videoFilter)
	}

	if len(b.filters.Audio) > 0 {
		audioFilter := strings.Join(b.filters.Audio, ",")
		args = append(args, "-af", audioFilter)
	}

	// Custom arguments
	args = append(args, b.customArgs...)

	// Output file
	args = append(args, b.outputFile)

	return args, nil
}

// Execute builds and executes the FFmpeg command
func (b *Builder) Execute(ctx context.Context) error {
	args, err := b.Build()
	if err != nil {
		return fmt.Errorf("failed to build command: %w", err)
	}

	executor := NewExecutor(b.config.BinaryPath)
	opts := &ExecuteOptions{
		Context: ctx,
		Timeout: time.Duration(b.config.Timeout) * time.Second,
	}

	return executor.Execute(ctx, args, opts)
}

// ExecuteWithProgress builds and executes the FFmpeg command with progress tracking
func (b *Builder) ExecuteWithProgress(ctx context.Context, progressHandler func(ProgressInfo)) error {
	args, err := b.Build()
	if err != nil {
		return fmt.Errorf("failed to build command: %w", err)
	}

	executor := NewExecutor(b.config.BinaryPath)
	opts := &ExecuteOptions{
		Context:         ctx,
		ProgressHandler: progressHandler,
		Timeout:         time.Duration(b.config.Timeout) * time.Second,
	}

	return executor.Execute(ctx, args, opts)
}

// formatDuration formats time.Duration to FFmpeg time format (HH:MM:SS.mmm)
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := d.Seconds() - float64(hours*3600) - float64(minutes*60)

	return fmt.Sprintf("%02d:%02d:%06.3f", hours, minutes, seconds)
}
