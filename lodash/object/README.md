# Object Package

Advanced object manipulation utilities for Go, providing powerful functions for working with maps, structs, and complex data structures. This package offers 8+ essential functions for object transformation, analysis, and manipulation.

## Features

- **🚀 High Performance**: Optimized object operations with minimal allocations
- **🔒 Thread Safe**: All functions are safe for concurrent use
- **🎯 Type Safe**: Full Go generics support for type safety
- **📦 Zero Dependencies**: No external dependencies
- **✅ Well Tested**: 68.7% test coverage with comprehensive scenarios
- **🔄 Flexible**: Works with maps, structs, and interfaces

## Installation

```bash
go get github.com/nguyendkn/go-libs/lodash/object
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/nguyendkn/go-libs/lodash/object"
)

func main() {
    data := map[string]interface{}{
        "name": "John",
        "age":  30,
        "city": "New York",
    }

    // Get object keys
    keys := object.Keys(data)
    fmt.Println(keys) // ["name", "age", "city"]

    // Get object values
    values := object.Values(data)
    fmt.Println(values) // ["John", 30, "New York"]

    // Pick specific fields
    subset := object.Pick(data, []string{"name", "age"})
    fmt.Println(subset) // map[string]interface{}{"name": "John", "age": 30}

    // Omit specific fields
    filtered := object.Omit(data, []string{"age"})
    fmt.Println(filtered) // map[string]interface{}{"name": "John", "city": "New York"}
}
```

## Function Categories

### 🔍 **Object Analysis**
- **`Keys`** - Get all keys from object
- **`Values`** - Get all values from object
- **`Entries`** - Get key-value pairs as slice
- **`Size`** - Get number of properties
- **`IsEmpty`** - Check if object is empty
- **`Has`** - Check if object has property
- **`HasIn`** - Check if object has property (including inherited)

### ✂️ **Object Filtering**
- **`Pick`** - Create object with only specified keys
- **`PickBy`** - Pick properties by predicate
- **`Omit`** - Create object without specified keys
- **`OmitBy`** - Omit properties by predicate

### 🔄 **Object Transformation**
- **`MapKeys`** - Transform object keys
- **`MapValues`** - Transform object values
- **`Transform`** - Transform object to accumulator
- **`TransformSlice`** - Transform slice to accumulator
- **`Invert`** - Swap keys and values
- **`InvertBy`** - Invert with custom key generation

### 🎯 **Object Utilities**
- **`Assign`** - Shallow merge objects
- **`AssignIn`** - Merge with inherited properties
- **`Defaults`** - Fill undefined properties
- **`DefaultsDeep`** - Deep fill undefined properties
- **`Merge`** - Deep merge objects
- **`MergeWith`** - Merge with custom merger

### 🔍 **Deep Operations**
- **`Get`** - Get value at path
- **`Set`** - Set value at path
- **`Has`** - Check if path exists
- **`Unset`** - Remove property at path
- **`Clone`** - Shallow clone object
- **`CloneDeep`** - Deep clone object

## Detailed Examples

### Object Manipulation
```go
user := map[string]interface{}{
    "id":      1,
    "name":    "Alice",
    "email":   "alice@example.com",
    "age":     25,
    "active":  true,
    "profile": map[string]interface{}{
        "bio":     "Software Engineer",
        "website": "https://alice.dev",
    },
}

// Extract specific fields
publicInfo := object.Pick(user, []string{"name", "profile"})

// Remove sensitive data
safeUser := object.Omit(user, []string{"email", "id"})

// Transform keys
kebabKeys := object.MapKeys(user, func(value interface{}, key string) string {
    return strings.ReplaceAll(key, "_", "-")
})
```