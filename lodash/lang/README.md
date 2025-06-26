# Lang Package

Comprehensive language utilities for Go, providing essential functions for type checking, conversion, and object manipulation. This package offers 26 high-performance, thread-safe functions for working with various data types and performing deep operations.

## Features

- **🚀 High Performance**: Optimized algorithms with minimal allocations
- **🔒 Thread Safe**: All functions are safe for concurrent use
- **🎯 Type Safe**: Robust type checking and conversion
- **📦 Zero Dependencies**: No external dependencies beyond Go standard library
- **✅ Well Tested**: Comprehensive test coverage with edge cases
- **🔍 Deep Operations**: Support for deep cloning and comparison

## Installation

```bash
go get github.com/nguyendkn/go-libs/lodash/lang
```

## Core Functions

### 🔍 **Type Checking**
- **`IsArray`** - Check if value is array or slice
- **`IsBoolean`** - Check if value is boolean
- **`IsDate`** - Check if value is time.Time
- **`IsEmpty`** - Check if value is empty
- **`IsEqual`** - Deep equality comparison
- **`IsError`** - Check if value is error
- **`IsFloat`** - Check if value is floating point
- **`IsFunction`** - Check if value is function
- **`IsInteger`** - Check if value is integer
- **`IsMap`** - Check if value is map
- **`IsNil`** - Check if value is nil
- **`IsNumber`** - Check if value is numeric
- **`IsObject`** - Check if value is object/struct
- **`IsPlainObject`** - Check if value is plain object
- **`IsRegExp`** - Check if value is regexp
- **`IsString`** - Check if value is string
- **`IsSymbol`** - Check if value is symbol
- **`IsArrayBuffer`** - Check if value is byte array

### 🔄 **Type Conversion**
- **`ToArray`** - Convert value to array
- **`ToInteger`** - Convert value to integer
- **`ToNumber`** - Convert value to number
- **`ToString`** - Convert value to string

### 📋 **Object Operations**
- **`Clone`** - Shallow clone of value
- **`CloneDeep`** - Deep clone of value

## Quick Examples

```go
// Type checking
lang.IsArray([]int{1, 2, 3})        // true
lang.IsString("hello")               // true
lang.IsEmpty("")                     // true
lang.IsEqual([]int{1, 2}, []int{1, 2}) // true

// Type conversion
lang.ToString(42)                    // "42"
lang.ToNumber("3.14")               // 3.14
lang.ToArray("hello")               // []interface{}{'h', 'e', 'l', 'l', 'o'}

// Object operations
original := []int{1, 2, 3}
cloned := lang.Clone(original)      // Shallow copy
deepCloned := lang.CloneDeep(original) // Deep copy
```

## Contributing

See the main [lodash README](../README.md) for contribution guidelines.

## License

This package is part of the go-libs project and follows the same license terms.