# FFmpeg Go Package

A comprehensive Go package for integrating FFmpeg into Wails applications with clean architecture, fluent API, and extensive functionality.

## Features

- **Auto-detection** of FFmpeg binary across platforms (Windows, macOS, Linux)
- **Fluent API** for building FFmpeg commands with ease
- **Progress tracking** with real-time updates
- **Context support** for timeouts and cancellation
- **Thread-safe** implementation
- **Extension functions** for common use cases
- **Comprehensive error handling**
- **Cross-platform compatibility**

## Installation

```bash
go get ./apps/backend/third_party/ffmpeg
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "your-project/apps/backend/third_party/ffmpeg"
)

func main() {
    // Create FFmpeg instance with auto-detection
    ff, err := ffmpeg.New()
    if err != nil {
        log.Fatal("Failed to initialize FFmpeg:", err)
    }
    
    // Validate installation
    if err := ff.ValidateInstallation(); err != nil {
        log.Fatal("FFmpeg validation failed:", err)
    }
    
    // Simple video conversion
    ctx := context.Background()
    err = ff.New().
        Input("input.mp4").
        Output("output.mp4").
        VideoCodec(ffmpeg.VideoCodecH264).
        AudioCodec(ffmpeg.AudioCodecAAC).
        Quality(ffmpeg.QualityMedium).
        Execute(ctx)
    
    if err != nil {
        log.Fatal("Conversion failed:", err)
    }
    
    fmt.Println("Conversion completed successfully!")
}
```

## API Documentation

### Core Interface

#### Creating FFmpeg Instance

```go
// Auto-detect FFmpeg binary
ff, err := ffmpeg.New()

// Or with custom configuration
config := &ffmpeg.Config{
    BinaryPath: "/path/to/ffmpeg",
    Timeout:    600, // seconds
    LogLevel:   "info",
}
ff, err := ffmpeg.NewWithConfig(config)
```

#### Basic Information

```go
// Get FFmpeg version
version, err := ff.GetVersion()

// Get supported formats
formats, err := ff.GetSupportedFormats()

// Get supported codecs
codecs, err := ff.GetSupportedCodecs()

// Get media file information
info, err := ff.GetMediaInfo("video.mp4")
```

### Command Builder API

The fluent API allows you to build complex FFmpeg commands easily:

```go
builder := ff.New().
    Input("input.mp4").
    Output("output.mp4").
    VideoCodec(ffmpeg.VideoCodecH264).
    AudioCodec(ffmpeg.AudioCodecAAC).
    Quality(ffmpeg.QualityHigh).
    Resolution(ffmpeg.Resolution1080p).
    VideoBitrate("2000k").
    AudioBitrate("192k").
    FrameRate(30).
    StartTime(time.Second * 10).
    Duration(time.Minute * 5)

// Execute the command
ctx := context.Background()
err := builder.Execute(ctx)
```

#### Available Builder Methods

- `Input(files ...string)` - Add input file(s)
- `Output(file string)` - Set output file
- `VideoCodec(codec VideoCodec)` - Set video codec
- `AudioCodec(codec AudioCodec)` - Set audio codec
- `Quality(quality Quality)` - Set encoding quality/preset
- `Resolution(resolution Resolution)` - Set output resolution
- `VideoBitrate(bitrate string)` - Set video bitrate
- `AudioBitrate(bitrate string)` - Set audio bitrate
- `FrameRate(fps float64)` - Set frame rate
- `SampleRate(rate int)` - Set audio sample rate
- `Channels(channels int)` - Set audio channels
- `StartTime(t time.Duration)` - Set start time for trimming
- `Duration(d time.Duration)` - Set duration for trimming
- `VideoFilter(filters ...string)` - Add video filters
- `AudioFilter(filters ...string)` - Add audio filters
- `Scale(width, height int)` - Add scale filter
- `Crop(width, height, x, y int)` - Add crop filter
- `Rotate(degrees float64)` - Add rotation filter
- `Volume(factor float64)` - Adjust audio volume
- `CustomArgs(args ...string)` - Add custom arguments
- `Overwrite(overwrite bool)` - Set overwrite behavior

### Progress Tracking

```go
err := ff.New().
    Input("input.mp4").
    Output("output.mp4").
    VideoCodec(ffmpeg.VideoCodecH264).
    ExecuteWithProgress(ctx, func(progress ffmpeg.ProgressInfo) {
        fmt.Printf("Progress: %.2f%% - Frame: %d - FPS: %.2f - Speed: %.2fx\n",
            progress.Progress, progress.Frame, progress.FPS, progress.Speed)
    })
```

### Extension Functions

#### Video Converter

```go
converter := ffmpeg.NewVideoConverter(ff)

// Convert to MP4
err := converter.ConvertToMP4("input.avi", "output.mp4")

// Convert to WebM
err := converter.ConvertToWebM("input.mp4", "output.webm")

// Convert with specific quality
err := converter.ConvertWithQuality("input.avi", "output.mp4", ffmpeg.QualityHigh)

// Resize video
err := converter.ResizeVideo("input.mp4", "output.mp4", ffmpeg.Resolution720p)
```

#### Audio Extractor

```go
extractor := ffmpeg.NewAudioExtractor(ff)

// Extract to MP3
err := extractor.ExtractToMP3("video.mp4", "audio.mp3", "192k")

// Extract to AAC
err := extractor.ExtractToAAC("video.mp4", "audio.m4a", "128k")

// Extract to FLAC (lossless)
err := extractor.ExtractToFLAC("video.mp4", "audio.flac")

// Extract segment
startTime := time.Second * 30
duration := time.Minute * 2
err := extractor.ExtractSegment("video.mp4", "segment.mp3", 
    startTime, duration, ffmpeg.AudioCodecMP3)
```

#### Video Compressor

```go
compressor := ffmpeg.NewVideoCompressor(ff)

// Compress to target size (in MB)
err := compressor.CompressToSize("input.mp4", "output.mp4", 100)

// Compress with specific bitrates
err := compressor.CompressWithBitrate("input.mp4", "output.mp4", "1000k", "128k")

// Compress for web
err := compressor.CompressForWeb("input.mp4", "output.mp4")
```

#### Thumbnail Generator

```go
generator := ffmpeg.NewThumbnailGenerator(ff)

// Generate at specific time
err := generator.GenerateAtTime("video.mp4", "thumb.jpg", time.Second*30)

// Generate at percentage of video
err := generator.GenerateAtPercentage("video.mp4", "thumb.jpg", 50.0)

// Generate with specific size
err := generator.GenerateWithSize("video.mp4", "thumb.jpg", 320, 240, time.Second*30)

// Generate multiple thumbnails
err := generator.GenerateMultiple("video.mp4", "thumb_%d.jpg", 5)
```

#### Video Editor

```go
editor := ffmpeg.NewVideoEditor(ff)

// Trim video
err := editor.TrimVideo("input.mp4", "output.mp4", 
    time.Second*10, time.Minute*5)

// Concatenate videos
inputs := []string{"video1.mp4", "video2.mp4", "video3.mp4"}
err := editor.ConcatenateVideos(inputs, "combined.mp4")

// Add watermark
err := editor.AddWatermark("video.mp4", "watermark.png", 
    "output.mp4", "bottom-right")

// Change speed
err := editor.ChangeSpeed("input.mp4", "output.mp4", 2.0) // 2x speed
```

## Constants

### Video Codecs

```go
ffmpeg.VideoCodecH264     // libx264
ffmpeg.VideoCodecH265     // libx265
ffmpeg.VideoCodecVP8      // libvpx
ffmpeg.VideoCodecVP9      // libvpx-vp9
ffmpeg.VideoCodecAV1      // libaom-av1
ffmpeg.VideoCodecCopy     // copy (no re-encoding)
```

### Audio Codecs

```go
ffmpeg.AudioCodecAAC      // aac
ffmpeg.AudioCodecMP3      // libmp3lame
ffmpeg.AudioCodecOpus     // libopus
ffmpeg.AudioCodecVorbis   // libvorbis
ffmpeg.AudioCodecFLAC     // flac
ffmpeg.AudioCodecCopy     // copy (no re-encoding)
```

### Quality Presets

```go
ffmpeg.QualityUltraFast   // Fastest encoding, largest file
ffmpeg.QualitySuperFast
ffmpeg.QualityVeryFast
ffmpeg.QualityFaster
ffmpeg.QualityFast
ffmpeg.QualityMedium      // Balanced
ffmpeg.QualitySlow
ffmpeg.QualitySlower
ffmpeg.QualityVerySlow    // Slowest encoding, smallest file
ffmpeg.QualityHigh        // High quality (CRF 18)
ffmpeg.QualityLow         // Low quality (CRF 28)
```

### Resolutions

```go
ffmpeg.Resolution144p     // 256x144
ffmpeg.Resolution240p     // 426x240
ffmpeg.Resolution360p     // 640x360
ffmpeg.Resolution480p     // 854x480
ffmpeg.Resolution720p     // 1280x720
ffmpeg.Resolution1080p    // 1920x1080
ffmpeg.Resolution1440p    // 2560x1440
ffmpeg.Resolution2160p    // 3840x2160 (4K)
```

## Error Handling

The package provides detailed error information:

```go
err := ff.New().Input("input.mp4").Output("output.mp4").Execute(ctx)
if err != nil {
    if ffmpegErr, ok := err.(*ffmpeg.FFmpegError); ok {
        fmt.Printf("FFmpeg Error: %s\n", ffmpegErr.Message)
        fmt.Printf("Exit Code: %d\n", ffmpegErr.Code)
        fmt.Printf("Command: %s\n", ffmpegErr.Command)
    } else {
        fmt.Printf("General Error: %s\n", err.Error())
    }
}
```

## Testing

Run tests with:

```bash
# Unit tests only
go test ./apps/backend/third_party/ffmpeg

# Include integration tests (requires FFmpeg installation)
go test -v ./apps/backend/third_party/ffmpeg

# Skip integration tests
go test -short ./apps/backend/third_party/ffmpeg

# Run benchmarks
go test -bench=. ./apps/backend/third_party/ffmpeg
```

## Requirements

- Go 1.19 or later
- FFmpeg binary installed on the system
- For development: FFmpeg with development headers (for advanced features)

## Platform Support

- **Windows**: Auto-detects FFmpeg in common installation paths
- **macOS**: Supports Homebrew and MacPorts installations
- **Linux**: Supports package manager installations and manual installs

## Examples

See the `examples/` directory for more detailed examples:

- `basic_conversion.go` - Basic video conversion
- `progress_tracking.go` - Progress tracking example
- `batch_processing.go` - Batch processing multiple files
- `advanced_filters.go` - Advanced video/audio filtering
- `web_optimization.go` - Web-optimized video processing

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This package is part of the VoiceAgents project and follows the same license terms.
