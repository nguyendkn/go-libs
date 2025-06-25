package hls

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nguyendkn/go-libs/ffmpeg"
)

// SegmentProcessor handles video segmentation for HLS
type SegmentProcessor struct {
	config *Config
	ffmpeg ffmpeg.FFmpeg
}

// NewSegmentProcessor creates a new segment processor
func NewSegmentProcessor(config *Config) *SegmentProcessor {
	return &SegmentProcessor{
		config: config,
		ffmpeg: config.FFmpeg,
	}
}

// ProcessSegments converts video to HLS segments for a specific quality level
func (sp *SegmentProcessor) ProcessSegments(ctx context.Context, inputFile string, qualityLevel QualityLevel) ([]Segment, error) {
	outputDir := sp.config.GetQualityOutputDir(qualityLevel.Name)
	segmentPattern := sp.config.GetSegmentPattern(qualityLevel.Name)

	// Build FFmpeg command for segmentation
	builder := sp.ffmpeg.New().
		Input(inputFile).
		VideoCodec(qualityLevel.VideoCodec).
		AudioCodec(qualityLevel.AudioCodec).
		Resolution(qualityLevel.Resolution).
		VideoBitrate(qualityLevel.VideoBitrate).
		AudioBitrate(qualityLevel.AudioBitrate)

	// Add frame rate if specified
	if qualityLevel.FrameRate > 0 {
		builder = builder.FrameRate(qualityLevel.FrameRate)
	}

	// Add HLS-specific options
	builder = sp.addHLSOptions(builder, outputDir, segmentPattern, qualityLevel)

	// Execute segmentation
	if err := builder.Execute(ctx); err != nil {
		return nil, &HLSError{
			Message: fmt.Sprintf("failed to process segments for quality %s", qualityLevel.Name),
			Code:    ErrCodeSegmentError,
			Cause:   err,
		}
	}

	// Collect generated segments
	segments, err := sp.collectSegments(outputDir, segmentPattern)
	if err != nil {
		return nil, err
	}

	return segments, nil
}

// addHLSOptions adds HLS-specific FFmpeg options
func (sp *SegmentProcessor) addHLSOptions(builder *ffmpeg.Builder, outputDir, segmentPattern string, qualityLevel QualityLevel) *ffmpeg.Builder {
	segmentPath := filepath.Join(outputDir, segmentPattern)
	playlistPath := sp.config.GetPlaylistPath(qualityLevel.Name)

	// Basic HLS options
	args := []string{
		"-f", "hls",
		"-hls_time", fmt.Sprintf("%.0f", sp.config.SegmentOptions.Duration.Seconds()),
		"-hls_playlist_type", string(sp.config.PlaylistType),
		"-hls_segment_filename", segmentPath,
	}

	// List size (0 = keep all segments)
	if sp.config.SegmentOptions.ListSize > 0 {
		args = append(args, "-hls_list_size", strconv.Itoa(sp.config.SegmentOptions.ListSize))
	} else {
		args = append(args, "-hls_list_size", "0")
	}

	// Start number
	if sp.config.SegmentOptions.StartNumber > 0 {
		args = append(args, "-hls_start_number_source", "generic",
			"-hls_start_number", strconv.Itoa(sp.config.SegmentOptions.StartNumber))
	}

	// Segment format specific options
	switch sp.config.SegmentOptions.Format {
	case SegmentMP4:
		args = append(args, "-hls_segment_type", "fmp4")
		if sp.config.SegmentOptions.InitSegment != "" {
			initPath := filepath.Join(outputDir, sp.config.SegmentOptions.InitSegment)
			args = append(args, "-hls_fmp4_init_filename", initPath)
		}
	case SegmentWebM:
		args = append(args, "-f", "webm_dash_manifest")
	default: // SegmentTS
		args = append(args, "-hls_segment_type", "mpegts")
	}

	// Encryption options
	if sp.config.Encryption != nil && sp.config.Encryption.Method != EncryptionNone {
		args = sp.addEncryptionOptions(args, outputDir)
	}

	// Advanced options
	if sp.config.FastStart {
		args = append(args, "-movflags", "+faststart")
	}

	// Force keyframes at segment boundaries
	segmentDuration := int(sp.config.SegmentOptions.Duration.Seconds())
	args = append(args, "-force_key_frames", fmt.Sprintf("expr:gte(t,n_forced*%d)", segmentDuration))

	// GOP size (should be segment duration * frame rate)
	if qualityLevel.FrameRate > 0 {
		gopSize := int(qualityLevel.FrameRate * sp.config.SegmentOptions.Duration.Seconds())
		args = append(args, "-g", strconv.Itoa(gopSize))
	}

	// Quality settings
	if qualityLevel.Profile != "" {
		args = append(args, "-profile:v", qualityLevel.Profile)
	}
	if qualityLevel.Level != "" {
		args = append(args, "-level:v", qualityLevel.Level)
	}

	// Add output path
	args = append(args, playlistPath)

	return builder.CustomArgs(args...)
}

// addEncryptionOptions adds encryption-related FFmpeg options
func (sp *SegmentProcessor) addEncryptionOptions(args []string, outputDir string) []string {
	enc := sp.config.Encryption

	// Key info file path
	keyInfoFile := filepath.Join(outputDir, "key_info.txt")

	// Create key info file
	if err := sp.createKeyInfoFile(keyInfoFile, enc); err == nil {
		args = append(args, "-hls_key_info_file", keyInfoFile)
	}

	return args
}

// createKeyInfoFile creates the key info file for HLS encryption
func (sp *SegmentProcessor) createKeyInfoFile(keyInfoFile string, enc *EncryptionOptions) error {
	content := fmt.Sprintf("%s\n%s\n", enc.KeyURI, enc.KeyFile)
	if enc.IV != "" {
		content += enc.IV + "\n"
	}

	return os.WriteFile(keyInfoFile, []byte(content), 0644)
}

// collectSegments collects information about generated segments
func (sp *SegmentProcessor) collectSegments(outputDir, segmentPattern string) ([]Segment, error) {
	var segments []Segment

	// Find all segment files matching the pattern
	segmentFiles, err := sp.findSegmentFiles(outputDir, segmentPattern)
	if err != nil {
		return nil, err
	}

	// Get segment information
	for i, segmentFile := range segmentFiles {
		segment := Segment{
			URI:      filepath.Base(segmentFile),
			Duration: sp.config.SegmentOptions.Duration.Seconds(),
		}

		// Get actual duration from file if possible
		if actualDuration, err := sp.getSegmentDuration(segmentFile); err == nil {
			segment.Duration = actualDuration
		}

		// Add discontinuity for first segment if needed
		if i == 0 && sp.config.SegmentOptions.StartNumber > 0 {
			segment.Discontinuity = false
		}

		segments = append(segments, segment)
	}

	return segments, nil
}

// findSegmentFiles finds all segment files in the output directory
func (sp *SegmentProcessor) findSegmentFiles(outputDir, segmentPattern string) ([]string, error) {
	// Extract base pattern and extension
	ext := filepath.Ext(segmentPattern)
	baseName := strings.TrimSuffix(segmentPattern, ext)

	// Remove %03d or similar format specifiers
	baseName = strings.ReplaceAll(baseName, "%03d", "*")
	baseName = strings.ReplaceAll(baseName, "%d", "*")

	pattern := filepath.Join(outputDir, baseName+ext)

	// Find matching files
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, &HLSError{
			Message: "failed to find segment files",
			Code:    ErrCodeSegmentError,
			Cause:   err,
		}
	}

	return matches, nil
}

// getSegmentDuration gets the actual duration of a segment file
func (sp *SegmentProcessor) getSegmentDuration(segmentFile string) (float64, error) {
	mediaInfo, err := sp.ffmpeg.GetMediaInfo(segmentFile)
	if err != nil {
		return 0, err
	}

	return mediaInfo.Duration.Seconds(), nil
}

// ProcessSegmentsWithProgress processes segments with progress reporting
func (sp *SegmentProcessor) ProcessSegmentsWithProgress(ctx context.Context, inputFile string, qualityLevel QualityLevel, progressCallback func(ConversionProgress)) ([]Segment, error) {
	// Get input file info for progress calculation
	inputInfo, err := sp.ffmpeg.GetMediaInfo(inputFile)
	if err != nil {
		return nil, &HLSError{
			Message: "failed to get input file info",
			Code:    ErrCodeInvalidInput,
			Cause:   err,
		}
	}

	totalDuration := inputInfo.Duration
	segmentDuration := sp.config.SegmentOptions.Duration
	estimatedSegments := int(totalDuration / segmentDuration)

	// Create progress handler
	progressHandler := func(ffmpegProgress ffmpeg.ProgressInfo) {
		if progressCallback != nil {
			// Calculate overall progress
			progress := (ffmpegProgress.Time.Seconds() / totalDuration.Seconds()) * 100
			if progress > 100 {
				progress = 100
			}

			// Calculate current segment
			currentSegment := int(ffmpegProgress.Time / segmentDuration)

			// Calculate ETA
			var eta time.Duration
			if ffmpegProgress.Speed > 0 {
				remainingTime := totalDuration - ffmpegProgress.Time
				eta = time.Duration(remainingTime.Seconds()/ffmpegProgress.Speed) * time.Second
			}

			conversionProgress := ConversionProgress{
				Stage:          "segmenting",
				QualityLevel:   qualityLevel.Name,
				Segment:        currentSegment,
				TotalSegments:  estimatedSegments,
				Progress:       progress,
				Speed:          ffmpegProgress.Speed,
				ETA:            eta,
				FFmpegProgress: ffmpegProgress,
			}

			progressCallback(conversionProgress)
		}
	}

	// Process with progress
	outputDir := sp.config.GetQualityOutputDir(qualityLevel.Name)
	segmentPattern := sp.config.GetSegmentPattern(qualityLevel.Name)

	builder := sp.ffmpeg.New().
		Input(inputFile).
		VideoCodec(qualityLevel.VideoCodec).
		AudioCodec(qualityLevel.AudioCodec).
		Resolution(qualityLevel.Resolution).
		VideoBitrate(qualityLevel.VideoBitrate).
		AudioBitrate(qualityLevel.AudioBitrate)

	if qualityLevel.FrameRate > 0 {
		builder = builder.FrameRate(qualityLevel.FrameRate)
	}

	builder = sp.addHLSOptions(builder, outputDir, segmentPattern, qualityLevel)

	// Execute with progress
	args, err := builder.Build()
	if err != nil {
		return nil, &HLSError{
			Message: fmt.Sprintf("failed to build FFmpeg command for quality %s", qualityLevel.Name),
			Code:    ErrCodeSegmentError,
			Cause:   err,
		}
	}

	if err := sp.ffmpeg.ExecuteWithProgress(ctx, args, progressHandler); err != nil {
		return nil, &HLSError{
			Message: fmt.Sprintf("failed to process segments for quality %s", qualityLevel.Name),
			Code:    ErrCodeSegmentError,
			Cause:   err,
		}
	}

	// Collect segments
	return sp.collectSegments(outputDir, segmentPattern)
}

// CleanupSegments removes old segment files
func (sp *SegmentProcessor) CleanupSegments(outputDir string, keepCount int) error {
	segmentFiles, err := sp.findSegmentFiles(outputDir, "*.ts")
	if err != nil {
		return err
	}

	if len(segmentFiles) <= keepCount {
		return nil // Nothing to cleanup
	}

	// Sort by modification time (oldest first)
	// Implementation would sort files by timestamp

	// Remove oldest files
	filesToRemove := segmentFiles[:len(segmentFiles)-keepCount]
	for _, file := range filesToRemove {
		if err := os.Remove(file); err != nil {
			// Log error but continue
			continue
		}
	}

	return nil
}

// ValidateSegments validates that all segments were created successfully
func (sp *SegmentProcessor) ValidateSegments(segments []Segment, outputDir string) error {
	for _, segment := range segments {
		segmentPath := filepath.Join(outputDir, segment.URI)

		// Check if file exists
		if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
			return &HLSError{
				Message: fmt.Sprintf("segment file not found: %s", segment.URI),
				Code:    ErrCodeSegmentError,
			}
		}

		// Check if file is not empty
		if info, err := os.Stat(segmentPath); err == nil && info.Size() == 0 {
			return &HLSError{
				Message: fmt.Sprintf("segment file is empty: %s", segment.URI),
				Code:    ErrCodeSegmentError,
			}
		}
	}

	return nil
}
