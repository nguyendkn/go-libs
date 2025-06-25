# WebRTC Library for Go

🚀 **Thư viện WebRTC mạnh mẽ và dễ sử dụng cho Go** - Được xây dựng với kiến trúc clean, dễ maintain và mở rộng.

## ✨ Tính năng

### 🔗 Core WebRTC
- **PeerConnection**: Quản lý kết nối WebRTC với đầy đủ lifecycle
- **DataChannel**: Truyền dữ liệu peer-to-peer với reliability options
- **ICE Handling**: Tự động xử lý ICE candidates và NAT traversal
- **SDP Management**: Tạo và xử lý offer/answer một cách dễ dàng

### 📡 Signaling
- **Signaling Server**: WebSocket server với REST API
- **Signaling Client**: Client với auto-reconnect và middleware support
- **Room Management**: Multi-peer rooms với permission control
- **Authentication**: Pluggable auth system

### 🎵 Media Engine
- **Codec Support**: VP8, VP9, H264, Opus, PCMU, PCMA
- **Media Tracks**: Audio/Video track management
- **Media Streams**: Stream composition và manipulation
- **Media Recording**: Record streams với multiple formats

### 🏗️ Advanced Features
- **SFU Support**: Selective Forwarding Unit cho multi-peer calls
- **Quality Control**: Adaptive bitrate và resolution
- **Statistics**: Real-time connection và media stats
- **Middleware**: Extensible message processing pipeline

## 📦 Cài đặt

```bash
go get github.com/nguyendkn/go-libs/webrtc
```

## 🚀 Quick Start

### Basic PeerConnection

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/nguyendkn/go-libs/webrtc"
)

func main() {
    // Tạo PeerConnection
    config := &webrtc.PeerConnectionConfig{
        ICEServers: []webrtc.ICEServer{
            {URLs: []string{"stun:stun.l.google.com:19302"}},
        },
    }
    
    pc, err := webrtc.NewPeerConnection(config)
    if err != nil {
        log.Fatal(err)
    }
    defer pc.Close()
    
    // Event handlers
    pc.OnConnectionStateChange(func(state webrtc.ConnectionState) {
        fmt.Printf("Connection state: %v\n", state)
    })
    
    // Tạo offer
    offer, err := pc.CreateOffer(nil)
    if err != nil {
        log.Fatal(err)
    }
    
    pc.SetLocalDescription(offer)
    fmt.Printf("Offer SDP:\n%s\n", offer.SDP)
}
```

### Signaling Server

```go
package main

import (
    "context"
    "log"
    
    "github.com/nguyendkn/go-libs/webrtc"
)

func main() {
    server := webrtc.NewSignalingServer()
    
    // Event handlers
    server.OnPeerConnected(func(peer *webrtc.PeerInfo) {
        fmt.Printf("Peer connected: %s\n", peer.ID)
    })
    
    // Start server
    log.Fatal(server.Start(":8080"))
}
```

### Media Engine

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/nguyendkn/go-libs/webrtc"
)

func main() {
    engine := webrtc.NewMediaEngine()
    
    // Get user media
    constraints := &webrtc.MediaConstraints{
        Audio: &webrtc.AudioConstraints{Enabled: true},
        Video: &webrtc.VideoConstraints{Enabled: true},
    }
    
    stream, err := engine.GetUserMedia(constraints)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Got media stream with %d tracks\n", len(stream.Tracks))
}
```

## 📚 Examples

Thư mục `examples/` chứa các ví dụ chi tiết:

- **simple_peer**: Basic PeerConnection usage
- **signaling_server**: Complete signaling server
- **media_demo**: Media engine demonstration

Chạy examples:

```bash
cd examples/simple_peer
go run main.go

cd examples/signaling_server  
go run main.go

cd examples/media_demo
go run main.go
```

## 🏗️ Kiến trúc

### Core Components

```
webrtc/
├── types.go           # Core types và constants
├── interfaces.go      # Interface definitions
├── peer_connection.go # PeerConnection implementation
├── data_channel.go    # DataChannel implementation
├── media_engine.go    # Media processing
├── media_recorder.go  # Media recording
├── signaling_client.go # Signaling client
├── signaling_server.go # Signaling server
└── room.go           # Room management
```

### Design Principles

- **Clean Architecture**: Separation of concerns với clear interfaces
- **Extensibility**: Plugin-based architecture cho easy customization
- **Performance**: Optimized cho high-throughput scenarios
- **Reliability**: Comprehensive error handling và recovery
- **Testability**: Mockable interfaces cho easy testing

## 🔧 Configuration

### PeerConnection Config

```go
config := &webrtc.PeerConnectionConfig{
    ICEServers: []webrtc.ICEServer{
        {URLs: []string{"stun:stun.l.google.com:19302"}},
        {
            URLs: []string{"turn:turn.example.com:3478"},
            Username: "user",
            Credential: "pass",
        },
    },
    ConnectionTimeout:   30 * time.Second,
    DisconnectedTimeout: 5 * time.Second,
    FailedTimeout:       30 * time.Second,
    KeepAliveInterval:   25 * time.Second,
}
```

### Server Config

```go
config := &webrtc.ServerConfig{
    MaxRooms:        1000,
    MaxPeersPerRoom: 50,
    MessageTimeout:  30 * time.Second,
    PeerTimeout:     60 * time.Second,
    EnableAuth:      true,
    EnableCORS:      true,
    AllowedOrigins:  []string{"https://example.com"},
}
```

## 📊 Monitoring

### Statistics

```go
stats, err := pc.GetStats()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Bytes sent: %d\n", stats.BytesSent)
fmt.Printf("Bytes received: %d\n", stats.BytesReceived)
fmt.Printf("Packet loss: %.2f%%\n", stats.PacketLossRate*100)
fmt.Printf("RTT: %v\n", stats.RTT)
```

### Server Stats

```go
stats := server.GetStats()
fmt.Printf("Active connections: %d\n", stats.ActiveConnections)
fmt.Printf("Total rooms: %d\n", stats.TotalRooms)
fmt.Printf("Uptime: %d seconds\n", stats.Uptime/1000)
```

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestPeerConnection ./...
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Pion WebRTC](https://github.com/pion/webrtc) - Core WebRTC implementation
- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket support
- WebRTC community for specifications and best practices

---

**Made with ❤️ by [nguyendkn](https://github.com/nguyendkn)**
