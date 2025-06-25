package ffmpeg

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// VideoConverter provides high-level video conversion functions
type VideoConverter struct {
	ffmpeg FFmpeg
}

// NewVideoConverter creates a new video converter
func NewVideoConverter(ffmpeg FFmpeg) *VideoConverter {
	return &VideoConverter{ffmpeg: ffmpeg}
}

// ConvertToMP4 converts any video to MP4 format with H.264/AAC
func (vc *VideoConverter) ConvertToMP4(input, output string) error {
	return vc.ffmpeg.ConvertVideo(input, output, &ConversionOptions{
		VideoCodec: VideoCodecH264,
		AudioCodec: AudioCodecAAC,
		Quality:    QualityMedium,
	})
}

// ConvertToWebM converts video to WebM format with VP9/Opus
func (vc *VideoConverter) ConvertToWebM(input, output string) error {
	return vc.ffmpeg.ConvertVideo(input, output, &ConversionOptions{
		VideoCodec: VideoCodecVP9,
		AudioCodec: AudioCodecOpus,
		Quality:    QualityMedium,
	})
}

// ConvertWithQuality converts video with specific quality settings
func (vc *VideoConverter) ConvertWithQuality(input, output string, quality Quality) error {
	// Determine output format from extension
	ext := strings.ToLower(filepath.Ext(output))
	var videoCodec VideoCodec
	var audioCodec AudioCodec
	
	switch ext {
	case ".mp4":
		videoCodec = VideoCodecH264
		audioCodec = AudioCodecAAC
	case ".webm":
		videoCodec = VideoCodecVP9
		audioCodec = AudioCodecOpus
	case ".mkv":
		videoCodec = VideoCodecH264
		audioCodec = AudioCodecAAC
	case ".avi":
		videoCodec = VideoCodecH264
		audioCodec = AudioCodecMP3
	default:
		videoCodec = VideoCodecH264
		audioCodec = AudioCodecAAC
	}
	
	return vc.ffmpeg.ConvertVideo(input, output, &ConversionOptions{
		VideoCodec: videoCodec,
		AudioCodec: audioCodec,
		Quality:    quality,
	})
}

// ResizeVideo resizes video to specific resolution
func (vc *VideoConverter) ResizeVideo(input, output string, resolution Resolution) error {
	return vc.ffmpeg.ConvertVideo(input, output, &ConversionOptions{
		VideoCodec: VideoCodecH264,
		AudioCodec: AudioCodecAAC,
		Resolution: resolution,
		Quality:    QualityMedium,
	})
}

// AudioExtractor provides audio extraction and conversion functions
type AudioExtractor struct {
	ffmpeg FFmpeg
}

// NewAudioExtractor creates a new audio extractor
func NewAudioExtractor(ffmpeg FFmpeg) *AudioExtractor {
	return &AudioExtractor{ffmpeg: ffmpeg}
}

// ExtractToMP3 extracts audio to MP3 format
func (ae *AudioExtractor) ExtractToMP3(input, output string, bitrate string) error {
	if bitrate == "" {
		bitrate = "192k"
	}
	
	return ae.ffmpeg.ExtractAudio(input, output, &ConversionOptions{
		AudioCodec:   AudioCodecMP3,
		AudioBitrate: bitrate,
	})
}

// ExtractToAAC extracts audio to AAC format
func (ae *AudioExtractor) ExtractToAAC(input, output string, bitrate string) error {
	if bitrate == "" {
		bitrate = "128k"
	}
	
	return ae.ffmpeg.ExtractAudio(input, output, &ConversionOptions{
		AudioCodec:   AudioCodecAAC,
		AudioBitrate: bitrate,
	})
}

// ExtractToFLAC extracts audio to FLAC format (lossless)
func (ae *AudioExtractor) ExtractToFLAC(input, output string) error {
	return ae.ffmpeg.ExtractAudio(input, output, &ConversionOptions{
		AudioCodec: AudioCodecFLAC,
	})
}

// ExtractSegment extracts audio segment with start time and duration
func (ae *AudioExtractor) ExtractSegment(input, output string, startTime, duration time.Duration, codec AudioCodec) error {
	return ae.ffmpeg.ExtractAudio(input, output, &ConversionOptions{
		AudioCodec: codec,
		StartTime:  startTime,
		Duration:   duration,
	})
}

// VideoCompressor provides video compression functions
type VideoCompressor struct {
	ffmpeg FFmpeg
}

// NewVideoCompressor creates a new video compressor
func NewVideoCompressor(ffmpeg FFmpeg) *VideoCompressor {
	return &VideoCompressor{ffmpeg: ffmpeg}
}

// CompressToSize compresses video to target file size
func (vc *VideoCompressor) CompressToSize(input, output string, targetSizeMB int64) error {
	targetSizeBytes := targetSizeMB * 1024 * 1024
	return vc.ffmpeg.CompressVideo(input, output, targetSizeBytes)
}

// CompressWithBitrate compresses video with specific bitrate
func (vc *VideoCompressor) CompressWithBitrate(input, output string, videoBitrate, audioBitrate string) error {
	return vc.ffmpeg.ConvertVideo(input, output, &ConversionOptions{
		VideoCodec:   VideoCodecH264,
		AudioCodec:   AudioCodecAAC,
		VideoBitrate: videoBitrate,
		AudioBitrate: audioBitrate,
		Quality:      QualityFast,
	})
}

// CompressForWeb compresses video optimized for web streaming
func (vc *VideoCompressor) CompressForWeb(input, output string) error {
	return vc.ffmpeg.ConvertVideo(input, output, &ConversionOptions{
		VideoCodec:   VideoCodecH264,
		AudioCodec:   AudioCodecAAC,
		VideoBitrate: "1000k",
		AudioBitrate: "128k",
		Quality:      QualityFast,
		Resolution:   Resolution720p,
	})
}

// ThumbnailGenerator provides thumbnail generation functions
type ThumbnailGenerator struct {
	ffmpeg FFmpeg
}

// NewThumbnailGenerator creates a new thumbnail generator
func NewThumbnailGenerator(ffmpeg FFmpeg) *ThumbnailGenerator {
	return &ThumbnailGenerator{ffmpeg: ffmpeg}
}

// GenerateAtTime generates thumbnail at specific time
func (tg *ThumbnailGenerator) GenerateAtTime(input, output string, timeOffset time.Duration) error {
	return tg.ffmpeg.GenerateThumbnail(input, output, &ThumbnailOptions{
		Time:    timeOffset,
		Quality: 2,
		Format:  FormatMP4, // Will be determined by output extension
	})
}

// GenerateAtPercentage generates thumbnail at percentage of video duration
func (tg *ThumbnailGenerator) GenerateAtPercentage(input, output string, percentage float64) error {
	// Get video duration first
	info, err := tg.ffmpeg.GetMediaInfo(input)
	if err != nil {
		return fmt.Errorf("failed to get media info: %w", err)
	}
	
	timeOffset := time.Duration(float64(info.Duration) * percentage / 100)
	return tg.GenerateAtTime(input, output, timeOffset)
}

// GenerateMultiple generates multiple thumbnails at intervals
func (tg *ThumbnailGenerator) GenerateMultiple(input, outputPattern string, count int) error {
	// Get video duration first
	info, err := tg.ffmpeg.GetMediaInfo(input)
	if err != nil {
		return fmt.Errorf("failed to get media info: %w", err)
	}
	
	interval := info.Duration / time.Duration(count+1)
	
	for i := 1; i <= count; i++ {
		timeOffset := time.Duration(i) * interval
		output := fmt.Sprintf(outputPattern, i)
		
		if err := tg.GenerateAtTime(input, output, timeOffset); err != nil {
			return fmt.Errorf("failed to generate thumbnail %d: %w", i, err)
		}
	}
	
	return nil
}

// GenerateWithSize generates thumbnail with specific dimensions
func (tg *ThumbnailGenerator) GenerateWithSize(input, output string, width, height int, timeOffset time.Duration) error {
	return tg.ffmpeg.GenerateThumbnail(input, output, &ThumbnailOptions{
		Time:    timeOffset,
		Width:   width,
		Height:  height,
		Quality: 2,
	})
}

// VideoEditor provides basic video editing functions
type VideoEditor struct {
	ffmpeg FFmpeg
}

// NewVideoEditor creates a new video editor
func NewVideoEditor(ffmpeg FFmpeg) *VideoEditor {
	return &VideoEditor{ffmpeg: ffmpeg}
}

// TrimVideo trims video to specific start time and duration
func (ve *VideoEditor) TrimVideo(input, output string, startTime, duration time.Duration) error {
	return ve.ffmpeg.ConvertVideo(input, output, &ConversionOptions{
		VideoCodec: VideoCodecCopy,
		AudioCodec: AudioCodecCopy,
		StartTime:  startTime,
		Duration:   duration,
	})
}

// ConcatenateVideos concatenates multiple videos into one
func (ve *VideoEditor) ConcatenateVideos(inputs []string, output string) error {
	builder := ve.ffmpeg.New()
	
	// Add all input files
	for _, input := range inputs {
		builder = builder.Input(input)
	}
	
	// Use concat filter
	filterComplex := fmt.Sprintf("concat=n=%d:v=1:a=1", len(inputs))
	builder = builder.CustomArgs("-filter_complex", filterComplex)
	builder = builder.Output(output)
	
	ctx := context.Background()
	return builder.Execute(ctx)
}

// AddWatermark adds watermark to video
func (ve *VideoEditor) AddWatermark(input, watermark, output string, position string) error {
	var overlay string
	
	switch position {
	case "top-left":
		overlay = "overlay=10:10"
	case "top-right":
		overlay = "overlay=main_w-overlay_w-10:10"
	case "bottom-left":
		overlay = "overlay=10:main_h-overlay_h-10"
	case "bottom-right":
		overlay = "overlay=main_w-overlay_w-10:main_h-overlay_h-10"
	case "center":
		overlay = "overlay=(main_w-overlay_w)/2:(main_h-overlay_h)/2"
	default:
		overlay = "overlay=10:10" // default to top-left
	}
	
	builder := ve.ffmpeg.New().
		Input(input).
		Input(watermark).
		CustomArgs("-filter_complex", overlay).
		VideoCodec(VideoCodecH264).
		AudioCodec(AudioCodecCopy).
		Output(output)
	
	ctx := context.Background()
	return builder.Execute(ctx)
}

// ChangeSpeed changes video playback speed
func (ve *VideoEditor) ChangeSpeed(input, output string, speed float64) error {
	videoFilter := fmt.Sprintf("setpts=%.2f*PTS", 1.0/speed)
	audioFilter := fmt.Sprintf("atempo=%.2f", speed)
	
	builder := ve.ffmpeg.New().
		Input(input).
		VideoFilter(videoFilter).
		AudioFilter(audioFilter).
		VideoCodec(VideoCodecH264).
		AudioCodec(AudioCodecAAC).
		Output(output)
	
	ctx := context.Background()
	return builder.Execute(ctx)
}
