package ffmpeg

import (
	"context"
	"time"
)

// VideoCodec represents supported video codecs
type VideoCodec string

const (
	VideoCodecH264     VideoCodec = "libx264"
	VideoCodecH265     VideoCodec = "libx265"
	VideoCodecVP8      VideoCodec = "libvpx"
	VideoCodecVP9      VideoCodec = "libvpx-vp9"
	VideoCodecAV1      VideoCodec = "libaom-av1"
	VideoCodecCopy     VideoCodec = "copy"
	VideoCodecDefault  VideoCodec = ""
)

// AudioCodec represents supported audio codecs
type AudioCodec string

const (
	AudioCodecAAC      AudioCodec = "aac"
	AudioCodecMP3      AudioCodec = "libmp3lame"
	AudioCodecOpus     AudioCodec = "libopus"
	AudioCodecVorbis   AudioCodec = "libvorbis"
	AudioCodecFLAC     AudioCodec = "flac"
	AudioCodecCopy     AudioCodec = "copy"
	AudioCodecDefault  AudioCodec = ""
)

// Format represents output formats
type Format string

const (
	FormatMP4  Format = "mp4"
	FormatAVI  Format = "avi"
	FormatMOV  Format = "mov"
	FormatMKV  Format = "mkv"
	FormatWEBM Format = "webm"
	FormatMP3  Format = "mp3"
	FormatWAV  Format = "wav"
	FormatFLAC Format = "flac"
	FormatM4A  Format = "m4a"
)

// Quality represents quality presets
type Quality string

const (
	QualityUltraFast Quality = "ultrafast"
	QualitySuperFast Quality = "superfast"
	QualityVeryFast  Quality = "veryfast"
	QualityFaster    Quality = "faster"
	QualityFast      Quality = "fast"
	QualityMedium    Quality = "medium"
	QualitySlow      Quality = "slow"
	QualitySlower    Quality = "slower"
	QualityVerySlow  Quality = "veryslow"
	QualityHigh      Quality = "high"
	QualityLow       Quality = "low"
)

// Resolution represents video resolutions
type Resolution string

const (
	Resolution144p  Resolution = "256x144"
	Resolution240p  Resolution = "426x240"
	Resolution360p  Resolution = "640x360"
	Resolution480p  Resolution = "854x480"
	Resolution720p  Resolution = "1280x720"
	Resolution1080p Resolution = "1920x1080"
	Resolution1440p Resolution = "2560x1440"
	Resolution2160p Resolution = "3840x2160"
)

// ProgressInfo contains information about encoding progress
type ProgressInfo struct {
	Frame     int64         `json:"frame"`
	FPS       float64       `json:"fps"`
	Bitrate   string        `json:"bitrate"`
	Time      time.Duration `json:"time"`
	Size      int64         `json:"size"`
	Speed     float64       `json:"speed"`
	Progress  float64       `json:"progress"` // 0-100
}

// MediaInfo contains metadata about media files
type MediaInfo struct {
	Duration    time.Duration `json:"duration"`
	Width       int           `json:"width"`
	Height      int           `json:"height"`
	VideoCodec  string        `json:"video_codec"`
	AudioCodec  string        `json:"audio_codec"`
	Bitrate     int64         `json:"bitrate"`
	FrameRate   float64       `json:"frame_rate"`
	SampleRate  int           `json:"sample_rate"`
	Channels    int           `json:"channels"`
	Size        int64         `json:"size"`
	Format      string        `json:"format"`
}

// Command represents an FFmpeg command
type Command struct {
	Args []string `json:"args"`
	Dir  string   `json:"dir"`
}

// ExecuteOptions contains options for command execution
type ExecuteOptions struct {
	Context         context.Context
	ProgressHandler func(ProgressInfo)
	ErrorHandler    func(error)
	Timeout         time.Duration
	WorkingDir      string
}

// FFmpegError represents an FFmpeg-specific error
type FFmpegError struct {
	Message string
	Code    int
	Command string
}

func (e *FFmpegError) Error() string {
	return e.Message
}

// FilterChain represents a video/audio filter chain
type FilterChain struct {
	Video []string `json:"video"`
	Audio []string `json:"audio"`
}

// ConversionOptions contains options for media conversion
type ConversionOptions struct {
	VideoCodec    VideoCodec
	AudioCodec    AudioCodec
	Quality       Quality
	Resolution    Resolution
	VideoBitrate  string
	AudioBitrate  string
	FrameRate     float64
	SampleRate    int
	Channels      int
	StartTime     time.Duration
	Duration      time.Duration
	Filters       FilterChain
	CustomArgs    []string
}

// ThumbnailOptions contains options for thumbnail generation
type ThumbnailOptions struct {
	Time       time.Duration
	Width      int
	Height     int
	Quality    int // 1-31, lower is better
	Format     Format
	Count      int // Number of thumbnails to generate
	Interval   time.Duration
}
