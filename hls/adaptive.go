package hls

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/nguyendkn/go-libs/ffmpeg"
)

// AdaptiveStreaming handles adaptive bitrate streaming functionality
type AdaptiveStreaming struct {
	config    *Config
	converter *Converter
}

// NewAdaptiveStreaming creates a new adaptive streaming handler
func NewAdaptiveStreaming(config *Config) (*AdaptiveStreaming, error) {
	converter, err := NewConverter(config)
	if err != nil {
		return nil, err
	}

	return &AdaptiveStreaming{
		config:    config,
		converter: converter,
	}, nil
}

// GenerateAdaptiveStream creates multiple quality levels for adaptive streaming
func (as *AdaptiveStreaming) GenerateAdaptiveStream(ctx context.Context, inputFile string) (*ConversionResult, error) {
	// Analyze input to determine optimal quality levels
	inputInfo, err := as.config.FFmpeg.GetMediaInfo(inputFile)
	if err != nil {
		return nil, &HLSError{
			Message: "failed to analyze input for adaptive streaming",
			Code:    ErrCodeInvalidInput,
			Cause:   err,
		}
	}

	// Generate quality levels based on input
	qualityLevels := as.generateQualityLevels(inputInfo)
	
	// Update config with generated quality levels
	config := as.config.Clone()
	config.QualityLevels = qualityLevels
	config.AdaptiveBitrate = true

	// Create converter with updated config
	converter, err := NewConverter(config)
	if err != nil {
		return nil, err
	}

	return converter.Convert(ctx, inputFile)
}

// generateQualityLevels generates appropriate quality levels based on input characteristics
func (as *AdaptiveStreaming) generateQualityLevels(inputInfo *ffmpeg.MediaInfo) []QualityLevel {
	var levels []QualityLevel

	// Determine maximum quality based on input resolution
	maxWidth := inputInfo.Width
	maxHeight := inputInfo.Height
	maxBitrate := inputInfo.Bitrate

	// Define quality tiers
	qualityTiers := []struct {
		name       string
		maxWidth   int
		maxHeight  int
		bitrateMul float64
		profile    string
		level      string
	}{
		{"240p", 426, 240, 0.1, "baseline", "3.0"},
		{"360p", 640, 360, 0.2, "baseline", "3.0"},
		{"480p", 854, 480, 0.3, "main", "3.1"},
		{"720p", 1280, 720, 0.5, "main", "3.1"},
		{"1080p", 1920, 1080, 0.8, "high", "4.0"},
		{"1440p", 2560, 1440, 1.0, "high", "5.0"},
		{"2160p", 3840, 2160, 1.2, "high", "5.1"},
	}

	// Generate levels that don't exceed input quality
	for _, tier := range qualityTiers {
		if tier.maxWidth <= maxWidth && tier.maxHeight <= maxHeight {
			level := as.createQualityLevel(tier.name, tier.maxWidth, tier.maxHeight, maxBitrate, tier.bitrateMul, tier.profile, tier.level)
			levels = append(levels, level)
		}
	}

	// Ensure we have at least one quality level
	if len(levels) == 0 {
		// Create a single level matching input resolution
		level := as.createQualityLevel("source", maxWidth, maxHeight, maxBitrate, 1.0, "high", "4.0")
		levels = append(levels, level)
	}

	return levels
}

// createQualityLevel creates a quality level with specified parameters
func (as *AdaptiveStreaming) createQualityLevel(name string, width, height int, baseBitrate int64, bitrateMultiplier float64, profile, level string) QualityLevel {
	// Calculate video bitrate
	videoBitrate := int64(float64(baseBitrate) * bitrateMultiplier)
	if videoBitrate < 200000 { // Minimum 200kbps
		videoBitrate = 200000
	}

	// Calculate audio bitrate based on video quality
	var audioBitrate int64
	switch {
	case videoBitrate < 500000:
		audioBitrate = 64000  // 64kbps for low quality
	case videoBitrate < 2000000:
		audioBitrate = 128000 // 128kbps for medium quality
	default:
		audioBitrate = 192000 // 192kbps for high quality
	}

	return QualityLevel{
		Name:         name,
		Resolution:   ffmpeg.Resolution(fmt.Sprintf("%dx%d", width, height)),
		VideoBitrate: fmt.Sprintf("%dk", videoBitrate/1000),
		AudioBitrate: fmt.Sprintf("%dk", audioBitrate/1000),
		VideoCodec:   ffmpeg.VideoCodecH264,
		AudioCodec:   ffmpeg.AudioCodecAAC,
		FrameRate:    30,
		Profile:      profile,
		Level:        level,
	}
}

// OptimizeQualityLevels optimizes quality levels for better streaming experience
func (as *AdaptiveStreaming) OptimizeQualityLevels(levels []QualityLevel) []QualityLevel {
	if len(levels) <= 1 {
		return levels
	}

	// Sort by bitrate
	sort.Slice(levels, func(i, j int) bool {
		return as.getBitrateValue(levels[i].VideoBitrate) < as.getBitrateValue(levels[j].VideoBitrate)
	})

	var optimized []QualityLevel
	lastBitrate := int64(0)

	for _, level := range levels {
		currentBitrate := as.getBitrateValue(level.VideoBitrate)
		
		// Ensure minimum bitrate difference (at least 50% increase)
		if lastBitrate == 0 || currentBitrate >= lastBitrate*3/2 {
			optimized = append(optimized, level)
			lastBitrate = currentBitrate
		}
	}

	return optimized
}

// CreateCustomAdaptiveStream creates adaptive stream with custom quality levels
func (as *AdaptiveStreaming) CreateCustomAdaptiveStream(ctx context.Context, inputFile string, customLevels []QualityLevel) (*ConversionResult, error) {
	// Validate custom levels
	if err := as.validateQualityLevels(customLevels); err != nil {
		return nil, err
	}

	// Optimize levels
	optimizedLevels := as.OptimizeQualityLevels(customLevels)

	// Update config
	config := as.config.Clone()
	config.QualityLevels = optimizedLevels
	config.AdaptiveBitrate = true

	// Create converter
	converter, err := NewConverter(config)
	if err != nil {
		return nil, err
	}

	return converter.Convert(ctx, inputFile)
}

// validateQualityLevels validates a set of quality levels
func (as *AdaptiveStreaming) validateQualityLevels(levels []QualityLevel) error {
	if len(levels) == 0 {
		return &HLSError{
			Message: "at least one quality level is required",
			Code:    ErrCodeInvalidConfig,
		}
	}

	for i, level := range levels {
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

		if level.Resolution == "" {
			return &HLSError{
				Message: fmt.Sprintf("quality level %s: resolution is required", level.Name),
				Code:    ErrCodeInvalidConfig,
			}
		}
	}

	return nil
}

// getBitrateValue extracts numeric bitrate value from string
func (as *AdaptiveStreaming) getBitrateValue(bitrate string) int64 {
	if bitrate == "" {
		return 0
	}

	bitrate = strings.ToLower(bitrate)
	multiplier := int64(1)

	if strings.HasSuffix(bitrate, "k") {
		multiplier = 1000
		bitrate = strings.TrimSuffix(bitrate, "k")
	} else if strings.HasSuffix(bitrate, "m") {
		multiplier = 1000000
		bitrate = strings.TrimSuffix(bitrate, "m")
	}

	if value, err := strconv.ParseInt(bitrate, 10, 64); err == nil {
		return value * multiplier
	}

	return 0
}

// GeneratePresetAdaptiveStream generates adaptive stream using predefined presets
func (as *AdaptiveStreaming) GeneratePresetAdaptiveStream(ctx context.Context, inputFile string, preset AdaptivePreset) (*ConversionResult, error) {
	var qualityLevels []QualityLevel

	switch preset {
	case PresetMobile:
		qualityLevels = []QualityLevel{
			{
				Name: "240p", Resolution: ffmpeg.Resolution240p, VideoBitrate: "400k", AudioBitrate: "64k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "baseline", Level: "3.0",
			},
			{
				Name: "360p", Resolution: ffmpeg.Resolution360p, VideoBitrate: "800k", AudioBitrate: "96k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "baseline", Level: "3.0",
			},
		}
	case PresetWeb:
		qualityLevels = []QualityLevel{
			{
				Name: "360p", Resolution: ffmpeg.Resolution360p, VideoBitrate: "800k", AudioBitrate: "96k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "main", Level: "3.1",
			},
			{
				Name: "720p", Resolution: ffmpeg.Resolution720p, VideoBitrate: "2500k", AudioBitrate: "128k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "main", Level: "3.1",
			},
		}
	case PresetHD:
		qualityLevels = []QualityLevel{
			{
				Name: "480p", Resolution: ffmpeg.Resolution480p, VideoBitrate: "1200k", AudioBitrate: "128k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "main", Level: "3.1",
			},
			{
				Name: "720p", Resolution: ffmpeg.Resolution720p, VideoBitrate: "2500k", AudioBitrate: "128k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "high", Level: "4.0",
			},
			{
				Name: "1080p", Resolution: ffmpeg.Resolution1080p, VideoBitrate: "5000k", AudioBitrate: "192k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "high", Level: "4.0",
			},
		}
	case PresetUHD:
		qualityLevels = []QualityLevel{
			{
				Name: "720p", Resolution: ffmpeg.Resolution720p, VideoBitrate: "2500k", AudioBitrate: "128k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "high", Level: "4.0",
			},
			{
				Name: "1080p", Resolution: ffmpeg.Resolution1080p, VideoBitrate: "5000k", AudioBitrate: "192k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "high", Level: "4.0",
			},
			{
				Name: "2160p", Resolution: ffmpeg.Resolution2160p, VideoBitrate: "15000k", AudioBitrate: "256k",
				VideoCodec: ffmpeg.VideoCodecH264, AudioCodec: ffmpeg.AudioCodecAAC, FrameRate: 30, Profile: "high", Level: "5.1",
			},
		}
	default:
		qualityLevels = DefaultQualityLevels
	}

	return as.CreateCustomAdaptiveStream(ctx, inputFile, qualityLevels)
}

// AdaptivePreset represents predefined adaptive streaming presets
type AdaptivePreset string

const (
	PresetMobile AdaptivePreset = "mobile" // Low bandwidth, mobile-optimized
	PresetWeb    AdaptivePreset = "web"    // Web streaming, balanced quality
	PresetHD     AdaptivePreset = "hd"     // High definition streaming
	PresetUHD    AdaptivePreset = "uhd"    // Ultra high definition streaming
)

// GetBandwidthLadder returns the bandwidth ladder for quality levels
func (as *AdaptiveStreaming) GetBandwidthLadder(levels []QualityLevel) []BandwidthLevel {
	var ladder []BandwidthLevel

	for _, level := range levels {
		videoBitrate := as.getBitrateValue(level.VideoBitrate)
		audioBitrate := as.getBitrateValue(level.AudioBitrate)
		totalBandwidth := videoBitrate + audioBitrate

		ladder = append(ladder, BandwidthLevel{
			Name:      level.Name,
			Bandwidth: totalBandwidth,
			Resolution: string(level.Resolution),
			VideoCodec: string(level.VideoCodec),
			AudioCodec: string(level.AudioCodec),
		})
	}

	// Sort by bandwidth
	sort.Slice(ladder, func(i, j int) bool {
		return ladder[i].Bandwidth < ladder[j].Bandwidth
	})

	return ladder
}

// BandwidthLevel represents a bandwidth level in the adaptive ladder
type BandwidthLevel struct {
	Name       string `json:"name"`
	Bandwidth  int64  `json:"bandwidth"`
	Resolution string `json:"resolution"`
	VideoCodec string `json:"video_codec"`
	AudioCodec string `json:"audio_codec"`
}

// AnalyzeOptimalLevels analyzes input and suggests optimal quality levels
func (as *AdaptiveStreaming) AnalyzeOptimalLevels(inputFile string) ([]QualityLevel, error) {
	inputInfo, err := as.config.FFmpeg.GetMediaInfo(inputFile)
	if err != nil {
		return nil, &HLSError{
			Message: "failed to analyze input file",
			Code:    ErrCodeInvalidInput,
			Cause:   err,
		}
	}

	// Generate levels based on input characteristics
	levels := as.generateQualityLevels(inputInfo)
	
	// Optimize the levels
	optimized := as.OptimizeQualityLevels(levels)

	return optimized, nil
}

// CreateBandwidthOptimizedStream creates a stream optimized for specific bandwidth constraints
func (as *AdaptiveStreaming) CreateBandwidthOptimizedStream(ctx context.Context, inputFile string, maxBandwidth int64) (*ConversionResult, error) {
	// Analyze input
	levels, err := as.AnalyzeOptimalLevels(inputFile)
	if err != nil {
		return nil, err
	}

	// Filter levels by bandwidth constraint
	var filteredLevels []QualityLevel
	for _, level := range levels {
		videoBitrate := as.getBitrateValue(level.VideoBitrate)
		audioBitrate := as.getBitrateValue(level.AudioBitrate)
		totalBandwidth := videoBitrate + audioBitrate

		if totalBandwidth <= maxBandwidth {
			filteredLevels = append(filteredLevels, level)
		}
	}

	if len(filteredLevels) == 0 {
		return nil, &HLSError{
			Message: fmt.Sprintf("no quality levels fit within bandwidth constraint of %d bps", maxBandwidth),
			Code:    ErrCodeInvalidConfig,
		}
	}

	return as.CreateCustomAdaptiveStream(ctx, inputFile, filteredLevels)
}
