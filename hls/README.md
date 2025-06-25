# HLS Package

A comprehensive Go library for converting videos to HLS (HTTP Live Streaming) format with support for adaptive bitrate streaming, multiple quality levels, and advanced features.

## Features

- **HLS Conversion**: Convert videos to HLS format with customizable settings
- **Adaptive Bitrate Streaming**: Generate multiple quality levels for optimal viewing experience
- **Live Streaming**: Support for real-time live stream conversion
- **Multiple Formats**: Support for MPEG-TS, fMP4, and WebM segments
- **Encryption**: AES-128 encryption support for secure streaming
- **Progress Tracking**: Real-time conversion progress monitoring
- **Parallel Processing**: Multi-threaded conversion for faster processing
- **Flexible Configuration**: Extensive customization options
- **FFmpeg Integration**: Built on top of the robust FFmpeg package

## Installation

```bash
go get github.com/nguyendkn/go-libs/hls
```

## Quick Start

### Basic HLS Conversion

```go
package main

import (
    "context"
    "log"
    
    "github.com/nguyendkn/go-libs/hls"
)

func main() {
    // Simple conversion
    ctx := context.Background()
    result, err := hls.ConvertToHLS(ctx, "input.mp4", "output")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Conversion completed: %s", result.OutputDir)
}
```

### Adaptive Bitrate Streaming

```go
// Generate adaptive stream with multiple quality levels
ctx := context.Background()
result, err := hls.ConvertToAdaptiveHLS(ctx, "input.mp4", "output", hls.PresetHD)
if err != nil {
    log.Fatal(err)
}

log.Printf("Master playlist: %s", result.MasterPlaylist)
log.Printf("Quality levels: %d", len(result.QualityLevels))
```

## Advanced Usage

### Custom Configuration

```go
// Create HLS converter with custom settings
hlsConverter, err := hls.NewBuilder().
    WithOutputDir("output/custom").
    WithQualityLevels(hls.QualityLow, hls.QualityMedium, hls.QualityHigh).
    WithSegmentDuration("6s").
    WithPlaylistType(hls.PlaylistVOD).
    WithParallel(true, 4).
    WithProgressCallback(func(progress hls.ConversionProgress) {
        log.Printf("Progress: %.1f%% - %s", progress.Progress, progress.Stage)
    }).
    Build()

if err != nil {
    log.Fatal(err)
}
defer hlsConverter.Cleanup()

ctx := context.Background()
result, err := hlsConverter.Convert(ctx, "input.mp4")
if err != nil {
    log.Fatal(err)
}
```

### Custom Quality Levels

```go
// Define custom quality levels
customLevels := []hls.QualityLevel{
    {
        Name:         "mobile",
        Resolution:   ffmpeg.Resolution360p,
        VideoBitrate: "800k",
        AudioBitrate: "96k",
        VideoCodec:   ffmpeg.VideoCodecH264,
        AudioCodec:   ffmpeg.AudioCodecAAC,
        FrameRate:    30,
        Profile:      "baseline",
        Level:        "3.0",
    },
    {
        Name:         "desktop",
        Resolution:   ffmpeg.Resolution1080p,
        VideoBitrate: "4000k",
        AudioBitrate: "192k",
        VideoCodec:   ffmpeg.VideoCodecH264,
        AudioCodec:   ffmpeg.AudioCodecAAC,
        FrameRate:    30,
        Profile:      "high",
        Level:        "4.0",
    },
}

result, err := hlsConverter.CreateCustomAdaptiveStream(ctx, "input.mp4", customLevels)
```

### Encryption

```go
// Enable AES-128 encryption
encryption := &hls.EncryptionOptions{
    Method:  hls.EncryptionAES128,
    KeyURI:  "https://example.com/key.key",
    KeyFile: "encryption.key",
}

hlsConverter, err := hls.NewBuilder().
    WithOutputDir("output/encrypted").
    WithEncryption(encryption).
    Build()
```

### Live Streaming

```go
// Configure for live streaming
hlsConverter, err := hls.NewBuilder().
    WithOutputDir("output/live").
    WithPreset(hls.PresetLive).
    Build()

// Start live conversion
ctx := context.Background()
err = hlsConverter.ConvertLive(ctx, "rtmp://live.example.com/stream")
```

## Configuration Options

### Quality Presets

- `QualityLow`: 480p, 800k video, 96k audio
- `QualityMedium`: 720p, 2500k video, 128k audio  
- `QualityHigh`: 1080p, 5000k video, 192k audio
- `QualityUltra`: 2160p, 15000k video, 256k audio

### Adaptive Presets

- `PresetMobile`: Optimized for mobile devices (240p, 360p)
- `PresetWeb`: Balanced for web streaming (360p, 720p)
- `PresetHD`: High definition streaming (480p, 720p, 1080p)
- `PresetUHD`: Ultra high definition (720p, 1080p, 2160p)

### Config Presets

- `PresetFast`: Fast conversion, lower quality
- `PresetQuality`: High quality, slower conversion
- `PresetBalanced`: Balanced speed and quality
- `PresetLive`: Optimized for live streaming

## API Reference

### Main Interface

```go
type HLS interface {
    // Core conversion methods
    Convert(ctx context.Context, inputFile string) (*ConversionResult, error)
    ConvertWithOptions(ctx context.Context, inputFile string, options *ConversionOptions) (*ConversionResult, error)
    
    // Live streaming
    ConvertLive(ctx context.Context, inputSource string) error
    
    // Adaptive streaming
    GenerateAdaptiveStream(ctx context.Context, inputFile string) (*ConversionResult, error)
    CreateCustomAdaptiveStream(ctx context.Context, inputFile string, qualityLevels []QualityLevel) (*ConversionResult, error)
    GeneratePresetAdaptiveStream(ctx context.Context, inputFile string, preset AdaptivePreset) (*ConversionResult, error)
    
    // Analysis
    AnalyzeInput(inputFile string) (*ffmpeg.MediaInfo, error)
    AnalyzeOptimalLevels(inputFile string) ([]QualityLevel, error)
    GetBandwidthLadder(levels []QualityLevel) []BandwidthLevel
    
    // Configuration
    GetConfig() *Config
    UpdateConfig(config *Config) error
    
    // Utilities
    Cleanup() error
    ValidateInput(inputFile string) error
}
```

### Builder Pattern

```go
builder := hls.NewBuilder().
    WithOutputDir("output").
    WithQualityLevels(hls.QualityMedium, hls.QualityHigh).
    WithSegmentDuration("6s").
    WithPlaylistType(hls.PlaylistVOD).
    WithEncryption(encryptionOptions).
    WithParallel(true, 4).
    WithProgressCallback(progressFunc).
    WithPreset(hls.PresetBalanced)

hlsConverter, err := builder.Build()
```

### Convenience Functions

```go
// Simple conversions
hls.ConvertToHLS(ctx, inputFile, outputDir)
hls.ConvertToAdaptiveHLS(ctx, inputFile, outputDir, preset)
hls.ConvertWithPreset(ctx, inputFile, outputDir, configPreset)
```

## Output Structure

### Single Quality
```
output/
├── playlist.m3u8
├── segment_000.ts
├── segment_001.ts
└── segment_002.ts
```

### Adaptive Bitrate
```
output/
├── master.m3u8
├── low/
│   ├── playlist.m3u8
│   ├── segment_000.ts
│   └── segment_001.ts
├── medium/
│   ├── playlist.m3u8
│   ├── segment_000.ts
│   └── segment_001.ts
└── high/
    ├── playlist.m3u8
    ├── segment_000.ts
    └── segment_001.ts
```

## Performance

The HLS package is optimized for performance with:

- **Parallel Processing**: Convert multiple quality levels simultaneously
- **Efficient Segmentation**: Optimized FFmpeg parameters
- **Memory Management**: Minimal memory footprint
- **Progress Tracking**: Real-time monitoring without performance impact

## Error Handling

The package provides detailed error information:

```go
if err != nil {
    if hlsErr, ok := err.(*hls.HLSError); ok {
        log.Printf("HLS Error [%s]: %s", hlsErr.Code, hlsErr.Message)
        if hlsErr.Cause != nil {
            log.Printf("Caused by: %v", hlsErr.Cause)
        }
    }
}
```

## Requirements

- Go 1.21 or later
- FFmpeg 4.0 or later installed on the system
- Sufficient disk space for output files

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions are welcome! Please read the contributing guidelines and submit pull requests for any improvements.

## Support

For issues and questions:
- Create an issue on GitHub
- Check the examples in the `example_test.go` file
- Review the comprehensive test suite for usage patterns
