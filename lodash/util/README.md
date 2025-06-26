# Util Package

Comprehensive utility functions for Go, providing essential helper functions and utilities for common programming tasks. This package offers 29 high-performance, thread-safe functions for various utility operations, data manipulation, and functional programming patterns.

## Features

- **🚀 High Performance**: Optimized algorithms with minimal allocations
- **🔒 Thread Safe**: All functions are safe for concurrent use
- **🎯 Type Safe**: Full Go generics support for type safety
- **📦 Zero Dependencies**: No external dependencies beyond Go standard library
- **✅ Well Tested**: Comprehensive test coverage with edge cases
- **🔧 Versatile**: Wide range of utility functions for common tasks

## Installation

```bash
go get github.com/nguyendkn/go-libs/lodash/util
```

## Core Functions

### 🔧 **Basic Utilities**
- **`Identity`** - Return the same value that is used as argument
- **`Constant`** - Create function that returns constant value
- **`Noop`** - No operation function
- **`DefaultTo`** - Return default value if input is zero value
- **`DefaultToAny`** - Return default value for any type

### 🔢 **Number & Range Utilities**
- **`Range`** - Create array of numbers in range
- **`Times`** - Execute function n times
- **`Random`** - Generate random number
- **`Clamp`** - Clamp number between bounds
- **`InRange`** - Check if number is in range

### 🎲 **Collection Utilities**
- **`Sample`** - Get random element from collection
- **`SampleSize`** - Get n random elements from collection
- **`Shuffle`** - Shuffle collection randomly
- **`Size`** - Get size of collection

### 🔄 **Function Utilities**
- **`Flow`** - Create function pipeline (left to right)
- **`FlowRight`** - Create function pipeline (right to left)
- **`Attempt`** - Execute function and handle errors

### 🏷️ **ID & Path Utilities**
- **`UniqueId`** - Generate unique ID with optional prefix
- **`ToPath`** - Convert string to property path array
- **`Property`** - Create property accessor function
- **`PropertyOf`** - Create property accessor for object

### 🎭 **Stub Functions**
- **`StubArray`** - Return empty array
- **`StubFalse`** - Return false
- **`StubObject`** - Return empty object
- **`StubString`** - Return empty string
- **`StubTrue`** - Return true

### 🔍 **Matching Utilities**
- **`Matches`** - Create predicate function for partial matching

## Quick Examples

```go
// Basic utilities
util.Identity(42)                   // 42
util.DefaultTo("", "default")       // "default"

// Range and iteration
util.Range(1, 5)                    // [1, 2, 3, 4]
util.Times(3, func(i int) int { return i * 2 }) // [0, 2, 4]

// Collection utilities
numbers := []int{1, 2, 3, 4, 5}
util.Sample(numbers)                // random element
util.Shuffle(numbers)               // shuffled slice

// Function composition
addOne := func(x int) int { return x + 1 }
double := func(x int) int { return x * 2 }
pipeline := util.Flow(addOne, double)
result := pipeline(3)               // (3+1)*2 = 8

// Utilities
util.UniqueId("user_")              // "user_1"
util.ToPath("a.b.c")               // ["a", "b", "c"]
```

## Contributing

See the main [lodash README](../README.md) for contribution guidelines.

## License

This package is part of the go-libs project and follows the same license terms.