# RTSP Package

A comprehensive Go package for RTSP (Real Time Streaming Protocol) stream handling and conversion to HLS (HTTP Live Streaming) format.

## Features

### 🎥 RTSP Stream Management
- **Single & Multiple Streams**: Handle one or multiple RTSP streams simultaneously
- **Connection Management**: Automatic reconnection, timeout handling, and error recovery
- **Transport Protocols**: Support for TCP, UDP, and auto-detection
- **Authentication**: Username/password authentication support
- **Stream Monitoring**: Real-time metrics and health monitoring

### 🔄 HLS Conversion
- **RTSP to HLS**: Convert RTSP streams to HLS format using integrated FFmpeg and HLS packages
- **Multiple Modes**: Separate streams, merged streams, or both
- **Live Streaming**: Real-time conversion for live RTSP feeds
- **Quality Control**: Configurable video/audio quality settings

### 🖼️ Video Layout & Merging
- **Flexible Layouts**: Support for various grid layouts (1x1, 2x2, 3x3, 4x4, custom)
- **Auto Layout Detection**: Automatically choose optimal layout based on stream count
- **Video Merging**: Combine multiple RTSP streams into a single video output
- **Customizable Appearance**: Configurable padding, borders, background colors

### ⚡ Performance & Monitoring
- **Parallel Processing**: Concurrent stream handling and conversion
- **Progress Tracking**: Real-time conversion progress callbacks
- **Resource Management**: Memory and CPU usage optimization
- **Error Handling**: Comprehensive error reporting and recovery

## Installation

```bash
go get github.com/nguyendkn/go-libs/rtsp
```

## Quick Start

### Simple RTSP to HLS Conversion

```go
package main

import (
    "context"
    "log"
    
    "github.com/nguyendkn/go-libs/rtsp"
)

func main() {
    ctx := context.Background()
    
    // Convert single RTSP stream to HLS
    result, err := rtsp.ConvertSingleStream(ctx, 
        "rtsp://example.com/stream", 
        "output")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Conversion completed: %s", result.OutputDir)
}
```

### Multiple Streams with Layout

```go
package main

import (
    "context"
    "log"
    
    "github.com/nguyendkn/go-libs/rtsp"
)

func main() {
    ctx := context.Background()
    
    // Multiple RTSP streams
    streamURLs := []string{
        "rtsp://camera1.example.com/stream",
        "rtsp://camera2.example.com/stream",
        "rtsp://camera3.example.com/stream",
        "rtsp://camera4.example.com/stream",
    }
    
    // Convert with 2x2 layout
    layout := rtsp.DefaultLayouts[rtsp.Layout2x2]
    result, err := rtsp.ConvertWithLayout(ctx, streamURLs, "output", layout)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Merged stream created: %s", result.MergedStream.PlaylistPath)
}
```

### Advanced Configuration

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/nguyendkn/go-libs/rtsp"
    "github.com/nguyendkn/go-libs/ffmpeg"
)

func main() {
    ctx := context.Background()
    
    // Create RTSP instance with custom configuration
    rtspInstance, err := rtsp.NewBuilder().
        WithOutputDir("output").
        WithStreamingMode(rtsp.ModeBoth).
        WithLayout(rtsp.DefaultLayouts[rtsp.Layout3x3]).
        WithTransport(rtsp.TransportTCP).
        WithQuality(ffmpeg.Resolution1080p, "4000k", "192k", 30).
        WithParallel(true, 4).
        WithProgressCallback(func(progress rtsp.ConversionProgress) {
            log.Printf("Progress: %.1f%% - %s", 
                progress.Progress, progress.StreamName)
        }).
        Build()
    if err != nil {
        log.Fatal(err)
    }
    defer rtspInstance.Close()
    
    // Add streams
    streams := []rtsp.RTSPStream{
        {
            URL:       "rtsp://camera1.example.com/stream",
            Name:      "camera1",
            Username:  "admin",
            Password:  "password",
            Transport: rtsp.TransportTCP,
            Reconnect: true,
        },
        {
            URL:       "rtsp://camera2.example.com/stream",
            Name:      "camera2",
            Transport: rtsp.TransportUDP,
            Reconnect: true,
        },
    }
    
    err = rtspInstance.AddStreams(streams)
    if err != nil {
        log.Fatal(err)
    }
    
    // Convert streams
    streamURLs := []string{
        "rtsp://camera1.example.com/stream",
        "rtsp://camera2.example.com/stream",
    }
    
    result, err := rtspInstance.Convert(ctx, streamURLs, "output", rtsp.ModeBoth)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Conversion completed successfully")
    log.Printf("Separate streams: %d", len(result.Streams))
    if result.MergedStream != nil {
        log.Printf("Merged stream: %s", result.MergedStream.PlaylistPath)
    }
}
```

## Configuration Options

### Layout Types

```go
// Predefined layouts
rtsp.LayoutSingle  // 1x1 - Single stream
rtsp.Layout1x2     // 1 row, 2 columns  
rtsp.Layout2x1     // 2 rows, 1 column
rtsp.Layout2x2     // 2x2 grid
rtsp.Layout2x3     // 2 rows, 3 columns
rtsp.Layout3x2     // 3 rows, 2 columns
rtsp.Layout3x3     // 3x3 grid
rtsp.Layout4x4     // 4x4 grid
rtsp.LayoutCustom  // Custom layout

// Create custom layout
customLayout := rtsp.Layout{
    Type:    rtsp.LayoutCustom,
    Rows:    3,
    Columns: 4,
    Width:   1920,
    Height:  1080,
    Padding: 10,
    Background: "#000000",
    BorderWidth: 2,
    BorderColor: "#FFFFFF",
}
```

### Streaming Modes

```go
rtsp.ModeSeparate  // Each stream as separate HLS
rtsp.ModeMerged    // All streams merged into one HLS  
rtsp.ModeBoth      // Both separate and merged
```

### Transport Protocols

```go
rtsp.TransportTCP   // TCP transport (reliable)
rtsp.TransportUDP   // UDP transport (faster)
rtsp.TransportAuto  // Auto-detection
```

## Stream Management

### Adding Streams

```go
// Add single stream URL
err := rtspInstance.AddStreamURL("rtsp://example.com/stream")

// Add multiple stream URLs
urls := []string{
    "rtsp://camera1.example.com/stream",
    "rtsp://camera2.example.com/stream",
}
err := rtspInstance.AddStreamURLs(urls)

// Add stream with configuration
stream := rtsp.RTSPStream{
    URL:       "rtsp://camera.example.com/stream",
    Name:      "main_camera",
    Username:  "admin",
    Password:  "password",
    Transport: rtsp.TransportTCP,
    Timeout:   30 * time.Second,
    Reconnect: true,
    MaxRetries: 5,
    RetryDelay: 5 * time.Second,
}
err := rtspInstance.AddStream(stream)
```

### Stream Control

```go
// Start specific stream
err := rtspInstance.StartStream("camera1")

// Start all streams
err := rtspInstance.StartAllStreams()

// Stop specific stream  
err := rtspInstance.StopStream("camera1")

// Stop all streams
err := rtspInstance.StopAllStreams()
```

### Stream Information

```go
// Get stream info
info, err := rtspInstance.GetStreamInfo("camera1")
log.Printf("Status: %s, FPS: %.2f, Bitrate: %d", 
    info.Status, info.FPS, info.Bitrate)

// Get all stream info
allInfo := rtspInstance.GetAllStreamInfo()
for name, info := range allInfo {
    log.Printf("Stream %s: %s", name, info.Status)
}

// Get stream names
names := rtspInstance.GetStreamNames()
log.Printf("Streams: %v", names)
```

## Error Handling

The package provides detailed error information:

```go
result, err := rtsp.ConvertSingleStream(ctx, "rtsp://invalid", "output")
if err != nil {
    if rtspErr, ok := err.(*rtsp.RTSPError); ok {
        log.Printf("RTSP Error: %s (Code: %s)", rtspErr.Message, rtspErr.Code)
        if rtspErr.StreamURL != "" {
            log.Printf("Stream URL: %s", rtspErr.StreamURL)
        }
    } else {
        log.Printf("General error: %v", err)
    }
}
```

## Dependencies

This package integrates with:
- **FFmpeg Package**: `github.com/nguyendkn/go-libs/ffmpeg`
- **HLS Package**: `github.com/nguyendkn/go-libs/hls`

Make sure FFmpeg is installed on your system for video processing capabilities.

## License

This package is part of the go-libs collection and follows the same licensing terms.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
