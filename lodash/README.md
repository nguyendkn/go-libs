# Lodash for Go

A comprehensive utility library for Go, inspired by Lodash.js, providing high-performance, type-safe functions for working with arrays, collections, objects, strings, and more. Built with Go generics for maximum type safety and performance.

## 🚀 Features

- **🎯 Type Safe**: Full Go generics support for compile-time type safety
- **⚡ High Performance**: Optimized algorithms with minimal allocations
- **🔒 Thread Safe**: All functions are safe for concurrent use
- **📦 Zero Dependencies**: No external dependencies beyond Go standard library
- **✅ Well Tested**: Comprehensive test coverage across all packages
- **🧩 Modular**: Import only what you need with separate packages

## 📦 Packages Overview

| Package | Functions | Description |
|---------|-----------|-------------|
| **[Array](./array/README.md)** | 61 | Array and slice manipulation utilities |
| **[Collection](./collection/README.md)** | 29 | Collection processing and functional programming |
| **[Date](./date/README.md)** | 20 | Date and time manipulation utilities |
| **[Function](./function/README.md)** | 28 | Function composition, memoization, and control |
| **[Lang](./lang/README.md)** | 26 | Type checking, conversion, and object operations |
| **[Math](./math/README.md)** | 23 | Mathematical operations and statistics |
| **[Object](./object/README.md)** | 15+ | Object manipulation and property access |
| **[String](./string/README.md)** | 25+ | String processing and text manipulation |
| **[Util](./util/README.md)** | 29 | General utility functions and helpers |

**Total: 250+ functions** across 9 specialized packages.

## 🚀 Quick Start

### Installation

```bash
# Install all packages
go get github.com/nguyendkn/go-libs/lodash

# Or install specific packages
go get github.com/nguyendkn/go-libs/lodash/array
go get github.com/nguyendkn/go-libs/lodash/collection
go get github.com/nguyendkn/go-libs/lodash/string
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/nguyendkn/go-libs/lodash/array"
    "github.com/nguyendkn/go-libs/lodash/collection"
    "github.com/nguyendkn/go-libs/lodash/string"
)

func main() {
    // Array operations
    numbers := []int{1, 2, 3, 4, 5, 6}
    chunks := array.Chunk(numbers, 2)
    fmt.Println(chunks) // [[1 2] [3 4] [5 6]]
    
    // Collection processing
    evens := collection.Filter(numbers, func(n int) bool {
        return n%2 == 0
    })
    fmt.Println(evens) // [2 4 6]
    
    // String manipulation
    text := "hello world"
    camelCase := string.CamelCase(text)
    fmt.Println(camelCase) // "helloWorld"
}
```

## 📚 Package Documentation

### 🔧 [Array Package](./array/README.md)
**61 functions** for array and slice manipulation:
- **Basic Operations**: Chunk, Compact, Concat, Fill, Flatten, Reverse
- **Search & Access**: Head, Last, IndexOf, Nth, Find elements
- **Slicing**: Drop, Take, Slice with various options
- **Set Operations**: Union, Intersection, Difference, Xor
- **Transformations**: Zip, Unzip, FromPairs, Join

### 🔄 [Collection Package](./collection/README.md)
**29 functions** for functional programming and collection processing:
- **Filtering**: Filter, Reject, Find, Some, Every
- **Transformation**: Map, FlatMap, Reduce, GroupBy
- **Sampling**: Sample, Shuffle, SampleSize
- **Iteration**: ForEach, ForEachRight, ForEachWithIndex

### 📅 [Date Package](./date/README.md)
**20 functions** for date and time operations:
- **Creation**: Now, ToDate, IsDate, IsValid
- **Boundaries**: StartOfDay, EndOfDay, StartOfWeek, EndOfWeek
- **Operations**: Add, Sub, Before, After, Equal
- **Utilities**: Format, DaysInMonth, IsLeapYear

### ⚡ [Function Package](./function/README.md)
**28 functions** for function manipulation and control:
- **Timing**: Debounce, Throttle, Delay, Defer
- **Execution**: Once, After, Before, Memoize
- **Composition**: Compose, Pipe, Curry, Partial
- **Arguments**: Flip, Rearg, Ary, Unary

### 🔍 [Lang Package](./lang/README.md)
**26 functions** for type checking and conversion:
- **Type Checking**: IsArray, IsString, IsNumber, IsEmpty
- **Conversion**: ToString, ToNumber, ToArray, ToInteger
- **Object Operations**: Clone, CloneDeep, IsEqual

### 🧮 [Math Package](./math/README.md)
**23 functions** for mathematical operations:
- **Basic**: Add, Subtract, Multiply, Divide
- **Statistics**: Max, Min, Sum, Mean, MaxBy, MinBy
- **Number Ops**: Abs, Ceil, Floor, Round, Clamp, Random

### 📦 [Object Package](./object/README.md)
**15+ functions** for object manipulation:
- **Property Access**: Get, Set, Has, Omit, Pick
- **Transformation**: Keys, Values, Entries, Merge
- **Path Operations**: GetPath, SetPath, UnsetPath

### 📝 [String Package](./string/README.md)
**25+ functions** for string processing:
- **Case Conversion**: CamelCase, SnakeCase, KebabCase
- **Manipulation**: Trim, Pad, Repeat, Replace
- **Validation**: StartsWith, EndsWith, Includes

### 🛠️ [Util Package](./util/README.md)
**29 functions** for general utilities:
- **Basic**: Identity, Constant, Noop, DefaultTo
- **Range**: Range, Times, Random, Clamp
- **Function**: Flow, FlowRight, Attempt
- **ID & Path**: UniqueId, ToPath, Property

## 🎯 Key Features

### Type Safety with Generics
```go
// Type-safe operations with Go generics
numbers := []int{1, 2, 3, 4, 5}
doubled := collection.Map(numbers, func(n int) int {
    return n * 2
}) // Returns []int, not []interface{}

// Works with any type
strings := []string{"hello", "world"}
lengths := collection.Map(strings, func(s string) int {
    return len(s)
}) // Returns []int
```

### Performance Optimized
```go
// Efficient operations with minimal allocations
large := make([]int, 1000000)
for i := range large {
    large[i] = i
}

// Fast filtering with early termination where possible
evens := collection.Filter(large, func(n int) bool {
    return n%2 == 0
})

// Memoization for expensive computations
fibonacci := function.Memoize(func(n int) int {
    if n <= 1 { return n }
    return fibonacci(n-1) + fibonacci(n-2)
})
```

### Thread Safety
```go
// All functions are safe for concurrent use
var wg sync.WaitGroup
results := make([][]int, 10)

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(index int) {
        defer wg.Done()
        data := []int{1, 2, 3, 4, 5}
        results[index] = array.Chunk(data, 2)
    }(i)
}
wg.Wait()
```

## 🚀 Performance

- **Memory Efficient**: Minimal allocations and memory reuse where possible
- **Algorithm Optimized**: Uses efficient algorithms (e.g., binary search for sorted arrays)
- **Generic Benefits**: No boxing/unboxing overhead with Go generics
- **Concurrent Safe**: Thread-safe operations without performance penalties

## 📖 Examples

### Data Processing Pipeline
```go
type User struct {
    ID   int
    Name string
    Age  int
    Role string
}

users := []User{
    {1, "Alice", 25, "admin"},
    {2, "Bob", 30, "user"},
    {3, "Charlie", 35, "admin"},
    {4, "Diana", 28, "user"},
}

// Complex data processing pipeline
adminNames := collection.Map(
    collection.Filter(users, func(u User) bool {
        return u.Role == "admin"
    }),
    func(u User) string {
        return strings.ToUpper(u.Name)
    },
)
// Result: ["ALICE", "CHARLIE"]

// Group by role and get average age
byRole := collection.GroupBy(users, func(u User) string {
    return u.Role
})

for role, roleUsers := range byRole {
    avgAge := collection.Reduce(roleUsers, func(acc float64, u User) float64 {
        return acc + float64(u.Age)
    }, 0.0) / float64(len(roleUsers))
    
    fmt.Printf("%s average age: %.1f\n", role, avgAge)
}
```

### String Processing
```go
// Advanced string manipulation
text := "hello-world_example"

// Case conversions
camelCase := string.CamelCase(text)     // "helloWorldExample"
snakeCase := string.SnakeCase(text)     // "hello_world_example"
kebabCase := string.KebabCase(text)     // "hello-world-example"

// String utilities
padded := string.Pad("Go", 10, "*")     // "***Go****"
repeated := string.Repeat("Hi", 3)      // "HiHiHi"
truncated := string.Truncate("Long text here", 8) // "Long tex..."
```

### Mathematical Operations
```go
// Statistical operations
numbers := []float64{1.5, 2.3, 3.7, 4.1, 5.9}

max, _ := math.Max(numbers)             // 5.9
min, _ := math.Min(numbers)             // 1.5
sum := math.Sum(numbers)                // 17.5
mean, _ := math.Mean(numbers)           // 3.5

// Custom comparisons
type Product struct {
    Name  string
    Price float64
}

products := []Product{
    {"Laptop", 999.99},
    {"Mouse", 29.99},
    {"Keyboard", 79.99},
}

cheapest, _ := math.MinBy(products, func(p Product) float64 {
    return p.Price
}) // Returns Mouse product
```

## 🛠️ Installation & Usage

### Prerequisites
- Go 1.18+ (for generics support)

### Installation Options

```bash
# Install all packages
go get github.com/nguyendkn/go-libs/lodash

# Install specific packages only
go get github.com/nguyendkn/go-libs/lodash/array
go get github.com/nguyendkn/go-libs/lodash/collection
go get github.com/nguyendkn/go-libs/lodash/string

# Install multiple packages
go get github.com/nguyendkn/go-libs/lodash/{array,collection,string}
```

### Import Patterns

```go
// Import specific packages
import (
    "github.com/nguyendkn/go-libs/lodash/array"
    "github.com/nguyendkn/go-libs/lodash/collection"
)

// Use with aliases for clarity
import (
    lodashArray "github.com/nguyendkn/go-libs/lodash/array"
    lodashStr "github.com/nguyendkn/go-libs/lodash/string"
)

// Import all (not recommended for production)
import _ "github.com/nguyendkn/go-libs/lodash"
```

## 🧪 Testing

All packages include comprehensive test suites:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./array
go test ./collection

# Benchmark tests
go test -bench=. ./array
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/nguyendkn/go-libs.git
cd go-libs/lodash

# Run tests
go test ./...

# Run linting
golangci-lint run

# Format code
go fmt ./...
```

### Guidelines

- Follow Go conventions and best practices
- Add comprehensive tests for new functions
- Update documentation and examples
- Ensure thread safety for all functions
- Maintain backward compatibility

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Lodash.js](https://lodash.com/) - The original JavaScript utility library
- Built with ❤️ for the Go community
- Thanks to all contributors and users

## 📞 Support

- 📖 [Documentation](./docs/)
- 🐛 [Issue Tracker](https://github.com/nguyendkn/go-libs/issues)
- 💬 [Discussions](https://github.com/nguyendkn/go-libs/discussions)
- 📧 [Email Support](mailto:support@example.com)

---

**Made with ❤️ by the Go community**
