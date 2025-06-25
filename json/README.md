# Go JSON Library

A comprehensive, high-performance JSON library for Go with zero external dependencies.

## Features

- **Fast JSON Parsing**: Parse JSON from string, []byte, and io.Reader
- **Flexible Serialization**: Marshal/unmarshal structs, maps, slices with custom marshaling support
- **JSON Validation**: Validate JSON format with detailed error reporting and schema validation
- **Pretty Printing**: Format JSON with customizable indentation
- **JSON Query**: Extract data using JSON path/query syntax with filtering
- **JSON Manipulation**: Merge and manipulate JSON objects with path-based operations
- **Type Safety**: Safe conversion between JSON and Go types with conversion options
- **Nested Structures**: Full support for nested JSON structures
- **Thread Safe**: All operations are thread-safe
- **Zero Dependencies**: No external dependencies for core functionality

## Installation

```bash
go get github.com/go-libs/json
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/go-libs/json"
)

func main() {
    // Parse JSON
    data := `{"name": "John", "age": 30, "hobbies": ["reading", "coding"]}`
    parsed, err := json.Parse(data)
    if err != nil {
        panic(err)
    }

    // Extract values using paths
    name, _ := parsed.GetPath("name")
    age, _ := parsed.GetPath("age")
    firstHobby, _ := parsed.GetPath("hobbies[0]")

    fmt.Printf("Name: %s, Age: %s, First Hobby: %s\n",
        name.String(), age.String(), firstHobby.String())

    // Modify JSON
    parsed.SetPath("age", 31)
    parsed.SetPath("city", "New York")

    // Pretty print
    fmt.Println(parsed.PrettyString())
}
```

## Core API

### Parsing

```go
// Parse from string
value, err := json.Parse(`{"key": "value"}`)

// Parse from bytes
value, err := json.ParseBytes([]byte(`{"key": "value"}`))

// Parse from io.Reader
value, err := json.ParseReader(reader)

// Validate JSON
result := json.ValidateString(`{"key": "value"}`)
if !result.Valid {
    for _, err := range result.Errors {
        fmt.Println(err.Reason)
    }
}
```

### Type Checking

```go
value, _ := json.Parse(`{"name": "John", "age": 30, "active": true}`)

fmt.Println(value.IsObject())  // true
fmt.Println(value.IsArray())   // false
fmt.Println(value.IsString())  // false

nameValue, _ := value.GetPath("name")
fmt.Println(nameValue.IsString()) // true
```

### Value Extraction

```go
// Get typed values
name, err := value.GetPath("name")
nameStr, err := name.GetString()

age, err := value.GetPath("age")
ageInt, err := age.GetInt()

// Direct conversion
nameStr, err := value.GetPath("name").GetString()
ageInt, err := value.GetPath("age").GetInt()
```

### Path Operations

```go
// Set values using paths
value.SetPath("user.name", "John")
value.SetPath("user.hobbies[0]", "reading")
value.SetPath("config.debug", true)

// Check if path exists
if value.PathExists("user.email") {
    email, _ := value.GetPath("user.email")
    fmt.Println(email.String())
}

// Delete paths
value.DeletePath("user.age")
value.DeletePath("hobbies[1]")
```

## Advanced Features

### JSON Queries

```go
// Query with filters
query := json.NewQuery("products").
    Where("category", "=", "Electronics").
    Where("price", ">", 100).
    Select("name", "price")

results, err := query.Execute(jsonValue)
for _, result := range results {
    fmt.Println(result.PrettyString())
}

// Find patterns
names, _ := jsonValue.Find("products[*].name")
```

### Schema Validation

```go
schema := &json.Schema{
    Type: "object",
    Properties: map[string]*json.Schema{
        "name": {Type: "string"},
        "age":  {Type: "number", Minimum: &[]float64{0}[0]},
    },
    Required: []string{"name"},
}

result := value.ValidateSchema(schema)
if !result.Valid {
    for _, err := range result.Errors {
        fmt.Printf("Validation error: %s\n", err.Reason)
    }
}
```

### Type Conversion

```go
// Safe type conversion with options
opts := &json.ConversionOptions{
    StrictMode:  false,
    TimeFormat:  time.RFC3339,
    NullAsZero:  true,
}

// Convert to specific types
var userID int
err := value.GetPath("user_id").ConvertTo(&userID, opts)

// Convert to struct
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

var user User
err := value.UnmarshalTo(&user)
```

### JSON Manipulation

```go
// Clone
copy := value.Clone()

// Merge objects
additional := json.New(map[string]interface{}{
    "department": "Engineering",
    "salary":     75000,
})
value.Merge(additional)

// Transform values
value.Transform("*.price", func(v *json.Value) *json.Value {
    price, _ := v.GetFloat64()
    return json.New(price * 1.1) // Add 10% tax
})
```

## Examples

See the [examples](./examples/) directory for comprehensive usage examples:

- [Basic Usage](./examples/basic_usage.go) - Parsing, type checking, path operations
- [Advanced Usage](./examples/advanced_usage.go) - Queries, conversions, performance

## Performance

This library is optimized for performance while maintaining ease of use:

- Fast JSON parsing using Go's standard library
- Efficient path-based operations
- Minimal memory allocations
- Thread-safe operations

## Documentation

See [pkg.go.dev](https://pkg.go.dev/github.com/go-libs/json) for full API documentation.

## License

MIT License
