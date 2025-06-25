# OBS Studio WebSocket Client for Go

A comprehensive Go implementation of the OBS WebSocket protocol, inspired by and compatible with the official obs-websocket-js library. This library provides a clean, type-safe interface for controlling OBS Studio programmatically.

## Features

- **Full WebSocket Protocol Support**: Complete implementation of OBS WebSocket protocol v5
- **Authentication**: Automatic handling of challenge/salt authentication
- **Event System**: Comprehensive event handling with type-safe event data
- **Request/Response**: Synchronous request handling with timeout support
- **Batch Requests**: Support for batch operations
- **Thread-Safe**: Concurrent-safe design with proper synchronization
- **Typed API**: Strong typing for all requests, responses, and events
- **Error Handling**: Robust error handling and connection management

## Installation

```bash
go get -u obs_studio
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    obs "obs_studio"
)

func main() {
    // Create new client
    client := obs.NewClient()
    
    // Register event handlers
    client.On(obs.EventConnectionOpened, func(data interface{}) {
        fmt.Println("Connected to OBS")
    })
    
    client.On(obs.EventIdentified, func(data interface{}) {
        fmt.Println("Authenticated successfully")
        
        // Get scene list
        scenes, err := client.GetSceneList()
        if err != nil {
            log.Printf("Error getting scenes: %v", err)
            return
        }
        
        fmt.Printf("Current scene: %s\n", scenes.CurrentProgramSceneName)
        for _, scene := range scenes.Scenes {
            fmt.Printf("- %s\n", scene.SceneName)
        }
    })
    
    // Connect to OBS
    err := client.Connect("ws://localhost:4455", "your_password")
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    
    // Keep alive
    time.Sleep(10 * time.Second)
    
    // Disconnect
    client.Disconnect()
}
```

## Architecture

The library is designed with a modular architecture for maintainability and extensibility:

### Core Components

1. **`types.go`** - Core data structures and constants
2. **`auth.go`** - Authentication handling
3. **`events.go`** - Event system and event processing
4. **`requests.go`** - Request/response management
5. **`client.go`** - Main client implementation

### Design Patterns

- **Event Emitter Pattern**: For handling OBS events
- **Request/Response Pattern**: For API calls with automatic correlation
- **Builder Pattern**: For configuration options
- **Context Pattern**: For cancellation and timeouts

## API Reference

### Client Creation and Connection

```go
// Create client
client := obs.NewClient()

// Connect with default options
err := client.Connect("ws://localhost:4455", "password")

// Connect with custom options
err := client.Connect("ws://localhost:4455", "password",
    obs.WithRpcVersion(1),
    obs.WithEventSubscriptions(obs.EventSubscriptionAll),
    obs.WithConnectTimeout(10*time.Second),
    obs.WithRequestTimeout(30*time.Second),
)
```

### Event Handling

```go
// Connection events
client.On(obs.EventConnectionOpened, func(data interface{}) {})
client.On(obs.EventConnectionClosed, func(data interface{}) {})
client.On(obs.EventConnectionError, func(data interface{}) {})
client.On(obs.EventIdentified, func(data interface{}) {})

// Scene events
client.On(obs.EventCurrentProgramSceneChanged, func(data interface{}) {
    sceneData := data.(*obs.SceneChangedEventData)
    fmt.Printf("Scene changed to: %s\n", sceneData.SceneName)
})

// Stream events
client.On(obs.EventStreamStateChanged, func(data interface{}) {
    streamData := data.(*obs.StreamStateChangedEventData)
    if streamData.OutputActive {
        fmt.Println("Stream started")
    } else {
        fmt.Println("Stream stopped")
    }
})

// Remove event handlers
client.Off(obs.EventStreamStateChanged, handler)
```

### Common Operations

```go
// Scene management
scenes, err := client.GetSceneList()
err = client.SetCurrentScene("Scene Name")

// Streaming
err = client.StartStream()
err = client.StopStream()

// Recording
err = client.StartRecord()
err = client.StopRecord()

// Version info
version, err := client.GetVersion()
```

### Advanced Features

#### Batch Requests

```go
// Create batch requests
requests := []obs.RequestData{
    *client.CreateRequest(obs.RequestTypeGetVersion, nil),
    *client.CreateRequest(obs.RequestTypeGetSceneList, nil),
}

// Execute batch
response, err := client.CallBatch(requests, false)
```

#### Custom Requests

```go
// Send custom request
response, err := client.Call("CustomRequestType", map[string]interface{}{
    "customParam": "value",
})
```

#### Event Subscriptions

```go
// Subscribe to specific events only
client.Connect("ws://localhost:4455", "password",
    obs.WithEventSubscriptions(obs.EventSubscriptionScenes | obs.EventSubscriptionOutputs),
)

// Available subscriptions:
// - EventSubscriptionGeneral
// - EventSubscriptionConfig  
// - EventSubscriptionScenes
// - EventSubscriptionInputs
// - EventSubscriptionTransitions
// - EventSubscriptionFilters
// - EventSubscriptionOutputs
// - EventSubscriptionSceneItems
// - EventSubscriptionMediaInputs
// - EventSubscriptionVendors
// - EventSubscriptionUi
// - EventSubscriptionAll
```

## Error Handling

The library provides comprehensive error handling:

```go
client.On(obs.EventConnectionError, func(data interface{}) {
    errorData := data.(*obs.ConnectionEventData)
    log.Printf("Connection error: %s", errorData.Message)
})

// Check connection status
if !client.IsConnected() {
    log.Println("Not connected to OBS")
}

if !client.IsIdentified() {
    log.Println("Not authenticated with OBS")
}

// Handle request errors
response, err := client.GetSceneList()
if err != nil {
    log.Printf("Request failed: %v", err)
}
```

## Configuration

### Connection Options

```go
type ConnectionConfig struct {
    Address            string        // WebSocket address
    Password           string        // Authentication password
    RpcVersion         int           // Protocol version (default: 1)
    EventSubscriptions int           // Event subscription flags
    ConnectTimeout     time.Duration // Connection timeout
    RequestTimeout     time.Duration // Request timeout
}
```

### Default Values

- **RPC Version**: 1
- **Event Subscriptions**: All events
- **Connect Timeout**: 10 seconds
- **Request Timeout**: 30 seconds

## Thread Safety

The client is designed to be thread-safe:

- All public methods can be called from multiple goroutines
- Event handlers are executed in separate goroutines
- Internal state is protected with mutexes
- Channel-based message handling prevents race conditions

## Best Practices

### 1. Event Handler Registration

Register event handlers before connecting:

```go
client := obs.NewClient()

// Register handlers first
client.On(obs.EventIdentified, identifiedHandler)
client.On(obs.EventStreamStateChanged, streamHandler)

// Then connect
client.Connect("ws://localhost:4455", "password")
```

### 2. Error Handling

Always handle connection errors:

```go
client.On(obs.EventConnectionError, func(data interface{}) {
    errorData := data.(*obs.ConnectionEventData)
    log.Printf("Connection error: %s", errorData.Message)
    
    // Implement reconnection logic if needed
})
```

### 3. Graceful Shutdown

Properly disconnect when done:

```go
defer func() {
    if err := client.Disconnect(); err != nil {
        log.Printf("Disconnect error: %v", err)
    }
}()
```

### 4. Request Validation

Check connection status before making requests:

```go
if !client.IsIdentified() {
    return fmt.Errorf("client not ready for requests")
}

response, err := client.GetSceneList()
```

## Examples

See `example_test.go` for comprehensive usage examples including:

- Basic connection and scene management
- Advanced event handling
- Batch request processing
- Error handling patterns
- Stream control automation

## Compatibility

- **OBS Studio**: 28.0+ (with obs-websocket plugin)
- **OBS WebSocket Protocol**: Version 5.x
- **Go Version**: 1.21+

## Dependencies

- `github.com/gorilla/websocket` - WebSocket client implementation
- `github.com/google/uuid` - UUID generation for request IDs

## License

This library is provided under the same terms as the original obs-websocket-js project.

## Contributing

1. Follow Go coding standards
2. Add tests for new features
3. Update documentation
4. Ensure thread safety
5. Handle errors appropriately

## Troubleshooting

### Common Issues

1. **Authentication Failed**
   - Verify OBS WebSocket password
   - Ensure OBS WebSocket plugin is enabled
   - Check if authentication is required

2. **Connection Timeout**
   - Verify OBS is running
   - Check WebSocket address (default: ws://localhost:4455)
   - Ensure firewall allows connections

3. **Request Timeout**
   - Increase request timeout in configuration
   - Check OBS responsiveness
   - Verify request parameters

### Debug Mode

Enable debug logging for troubleshooting:

```go
// Add debug event handler
client.On("*", func(data interface{}) {
    log.Printf("Event: %+v", data)
})
```

## Roadmap

- [ ] Unit tests
- [ ] Integration tests
- [ ] Benchmarks
- [ ] More typed request/response structures
- [ ] Reconnection logic
- [ ] Metrics and monitoring
- [ ] gRPC wrapper (optional) 
