# UUID Package

A high-performance Go UUID library supporting UUIDv4 and UUIDv7 generation with comprehensive parsing capabilities.

## Features

- **UUIDv7 Support**: Time-ordered UUIDs with monotonic ordering guarantees
- **UUIDv4 Support**: Random UUIDs for general use
- **Multiple Formats**: Parse and generate UUIDs in various formats
- **High Performance**: Optimized for speed with minimal allocations
- **Thread Safe**: Concurrent generation with proper synchronization
- **Comprehensive**: Full RFC 9562 compliance

## Installation

```bash
go get github.com/nguyendkn/go-libs/uuid
```

## Quick Start

### Generate UUIDs

```go
package main

import (
    "fmt"
    "github.com/nguyendkn/go-libs/uuid"
)

func main() {
    // Generate UUIDv7 (time-ordered)
    id7 := uuid.UUIDv7()
    fmt.Println("UUIDv7:", id7)
    
    // Generate UUIDv4 (random)
    id4 := uuid.UUIDv4()
    fmt.Println("UUIDv4:", id4)
    
    // Generate without dashes
    noDash := uuid.UUIDv7(true)
    fmt.Println("No dash:", noDash)
    
    // Generate for database primary key
    pk := uuid.UUIDv7PrimaryKey()
    fmt.Println("Primary key:", pk) // Output: 'xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx'
}
```

### Parse UUIDs

```go
// Parse various formats
formats := []string{
    "0189dcd5-5311-7d40-8db0-9496a2eef37b",           // Standard
    "0189dcd553117d408db09496a2eef37b",               // No dashes
    "{0189dcd5-5311-7d40-8db0-9496a2eef37b}",         // With braces
    "urn:uuid:0189dcd5-5311-7d40-8db0-9496a2eef37b", // URN format
}

for _, format := range formats {
    parsed, err := uuid.Parse(format)
    if err != nil {
        panic(err)
    }
    fmt.Println("Parsed:", parsed.String())
}

// Check if string is valid UUID
if uuid.IsValid("0189dcd5-5311-7d40-8db0-9496a2eef37b") {
    fmt.Println("Valid UUID")
}
```

### Working with UUID Objects

```go
// Create UUID object
id := uuid.UUIDv7Obj()

// Get different string representations
fmt.Println("Standard:", id.String())           // 0189dcd5-5311-7d40-8db0-9496a2eef37b
fmt.Println("No dash:", id.StringNoDash())      // 0189dcd553117d408db09496a2eef37b
fmt.Println("Hex:", id.Hex())                   // 0189dcd553117d408db09496a2eef37b

// Get UUID properties
fmt.Println("Version:", id.GetVersion())        // 7
fmt.Println("Variant:", id.GetVariant())        // VAR_10

// Compare UUIDs
id2 := uuid.UUIDv7Obj()
fmt.Println("Equal:", id.Equals(id2))           // false
fmt.Println("Compare:", id.CompareTo(id2))      // -1, 0, or 1

// Get raw bytes
bytes := id.Bytes()
fmt.Printf("Bytes: %x\n", bytes)
```

## Advanced Usage

### Custom Generator

```go
// Create custom generator
gen := uuid.NewV7Generator()

// Generate with custom generator
id := gen.Generate()

// Generate or abort on clock rollback
id, err := gen.GenerateOrAbort()
if err != nil {
    // Handle clock rollback
}

// Generate UUIDv4 with same generator
id4 := gen.GenerateV4()
```

### Create from Fields

```go
// Create UUIDv7 from specific field values
timestamp := uint64(1640995200000) // Unix timestamp in milliseconds
randA := uint64(0x123)             // 12-bit random value
randBHi := uint64(0x456789)        // 30-bit random value
randBLo := uint64(0xabcdef01)      // 32-bit random value

id, err := uuid.FromFieldsV7(timestamp, randA, randBHi, randBLo)
if err != nil {
    panic(err)
}
```

### JSON Support

```go
import "encoding/json"

type User struct {
    ID   uuid.UUID `json:"id"`
    Name string    `json:"name"`
}

user := User{
    ID:   uuid.UUIDv7Obj(),
    Name: "John Doe",
}

// Marshal to JSON
data, _ := json.Marshal(user)
fmt.Println(string(data))

// Unmarshal from JSON
var parsed User
json.Unmarshal(data, &parsed)
```

## Performance

Benchmarks on Intel i5-13400:

```
BenchmarkUUIDv7-16                    13186508    89.83 ns/op
BenchmarkUUIDv4-16                     8005518   151.5 ns/op
BenchmarkV7Generator_Generate-16      22613391    52.63 ns/op
BenchmarkV7Generator_GenerateV4-16    10874656   108.4 ns/op
BenchmarkParse-16                      2846800   414.0 ns/op
BenchmarkParseNoDash-16                4546249   258.9 ns/op
```

## API Reference

### Main Functions

- `UUIDv7(noDash ...bool) string` - Generate UUIDv7 string
- `UUIDv4() string` - Generate UUIDv4 string
- `UUIDv7Obj() UUID` - Generate UUIDv7 object
- `UUIDv4Obj() UUID` - Generate UUIDv4 object
- `Parse(s string) (UUID, error)` - Parse UUID from string
- `IsValid(s string) bool` - Check if string is valid UUID

### UUID Methods

- `String() string` - Get standard format (with dashes)
- `StringNoDash() string` - Get format without dashes
- `Hex() string` - Alias for StringNoDash
- `Bytes() [16]byte` - Get raw bytes
- `GetVersion() int` - Get UUID version
- `GetVariant() Variant` - Get UUID variant
- `Equals(other UUID) bool` - Compare equality
- `CompareTo(other UUID) int` - Compare for ordering

## License

MIT License - see LICENSE file for details.
