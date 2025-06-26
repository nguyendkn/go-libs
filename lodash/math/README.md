# Math Package

Comprehensive mathematical utilities for Go, providing essential functions for numeric operations, statistics, and calculations. This package offers 23 high-performance, thread-safe functions for working with numbers and performing mathematical computations.

## Features

- **🚀 High Performance**: Optimized algorithms with minimal allocations
- **🔒 Thread Safe**: All functions are safe for concurrent use
- **🎯 Type Safe**: Full Go generics support for numeric types
- **📦 Zero Dependencies**: No external dependencies beyond Go standard library
- **✅ Well Tested**: Comprehensive test coverage with edge cases
- **🔢 Generic Numeric**: Works with all numeric types (int, float, etc.)

## Installation

```bash
go get github.com/nguyendkn/go-libs/lodash/math
```

## Core Functions

### ➕ **Basic Operations**
- **`Add`** - Add two numbers
- **`Subtract`** - Subtract two numbers
- **`Multiply`** - Multiply two numbers
- **`Divide`** - Divide two numbers

### 📊 **Statistics**
- **`Max`** - Find maximum value in slice
- **`Min`** - Find minimum value in slice
- **`Sum`** - Calculate sum of slice
- **`Mean`** - Calculate average of slice
- **`MaxBy`** - Find max using iteratee function
- **`MinBy`** - Find min using iteratee function
- **`SumBy`** - Calculate sum using iteratee function
- **`MeanBy`** - Calculate average using iteratee function

### 🔢 **Number Operations**
- **`Abs`** - Absolute value
- **`Ceil`** - Round up to nearest integer
- **`Floor`** - Round down to nearest integer
- **`Round`** - Round to nearest integer
- **`Clamp`** - Clamp number between bounds
- **`InRange`** - Check if number is in range
- **`Random`** - Generate random number
- **`Pow`** - Power operation
- **`Sqrt`** - Square root
- **`IsNaN`** - Check if value is NaN
- **`IsInf`** - Check if value is infinite

## Quick Examples

```go
// Basic operations
math.Add(5, 3)                      // 8
math.Multiply(4, 2.5)               // 10.0

// Statistics
numbers := []int{1, 2, 3, 4, 5}
math.Max(numbers)                   // 5, true
math.Sum(numbers)                   // 15
math.Mean(numbers)                  // 3.0, true

// Number operations
math.Abs(-5)                        // 5
math.Clamp(10, 0, 5)               // 5
math.InRange(3, 1, 5)              // true
math.Random(1, 10)                 // random number between 1-10
```

## Contributing

See the main [lodash README](../README.md) for contribution guidelines.

## License

This package is part of the go-libs project and follows the same license terms.