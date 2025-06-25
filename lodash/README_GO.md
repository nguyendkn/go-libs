# Go-Lodash v2.0.0

A comprehensive utility library for Go inspired by Lodash.js. Provides **100+ utility functions** across **9 specialized packages** for arrays, collections, strings, objects, math operations, and more.

## Features

- 🚀 **Zero external dependencies** - Pure Go implementation
- 🔒 **Thread-safe operations** - All functions are safe for concurrent use
- ⚡ **High performance** - Optimized implementations with minimal overhead
- 🧪 **Comprehensive test coverage** - 85%+ average test coverage across all packages
- 📚 **Clean, maintainable code** - Well-documented with clear examples
- 🔧 **Extensible design** - Modular structure for easy extension
- 🎯 **100+ Functions** - Complete feature parity with popular Lodash.js functions
- 📦 **9 Specialized Packages** - Organized by functionality for clean imports

## Installation

```bash
go get github.com/go-libs/lodash
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/go-libs/lodash/array"
    "github.com/go-libs/lodash/collection"
    "github.com/go-libs/lodash/string"
    "github.com/go-libs/lodash/math"
)

func main() {
    // Array operations
    chunks := array.Chunk([]int{1, 2, 3, 4, 5}, 2)
    fmt.Println(chunks) // [[1, 2], [3, 4], [5]]
    
    // Collection operations
    filtered := collection.Filter([]int{1, 2, 3, 4}, func(x int) bool { 
        return x%2 == 0 
    })
    fmt.Println(filtered) // [2, 4]
    
    // String operations
    camelCase := string.CamelCase("hello world")
    fmt.Println(camelCase) // "helloWorld"
    
    // Math operations
    sum := math.Sum([]int{1, 2, 3, 4, 5})
    fmt.Println(sum) // 15
}
```

## Packages

### Array Package (`array`) - 94.1% Coverage

**20+ functions** for working with arrays and slices:

- `Chunk` - Creates chunks of specified size
- `Compact` - Removes falsey values
- `Concat` - Concatenates arrays
- `Difference` - Returns values not in other arrays
- `Drop/DropRight` - Drops elements from beginning/end
- `Fill` - Fills array with value
- `Flatten/FlattenDeep` - Flattens nested arrays
- `Head/Last/Initial` - Gets first/last/initial elements
- `IndexOf/LastIndexOf` - Finds element indices
- `Intersection/Union` - Set operations
- `Join` - Joins elements into string
- `Nth` - Gets element at index
- `Pull/Remove` - Removes elements
- `Reverse` - Reverses array in place
- `Slice` - Creates array slice
- `Take/TakeRight` - Takes elements from start/end
- `Uniq` - Creates duplicate-free array
- `Without` - Excludes specified values
- `Zip` - Groups elements from multiple arrays

### Collection Package (`collection`) - 75.0% Coverage

**20+ functions** for working with collections:

- `CountBy` - Counts elements by key function
- `Every/Some` - Tests all/some elements
- `Filter/Reject` - Filters/rejects elements by predicate
- `Find/FindLast/FindIndex` - Finds elements/indices
- `FlatMap` - Maps and flattens results
- `ForEach/ForEachRight` - Iterates over elements
- `GroupBy/KeyBy` - Groups/keys elements by function
- `Includes` - Checks if value exists
- `Map` - Transforms elements
- `OrderBy/SortBy` - Sorts by criteria/function
- `Partition` - Splits into two groups
- `Reduce/ReduceRight` - Reduces to single value
- `Sample/SampleSize` - Random sampling
- `Shuffle` - Randomly shuffles elements
- `Size` - Gets collection size

### String Package (`string`) - 92.0% Coverage

**15+ functions** for working with strings:

- `CamelCase/KebabCase/SnakeCase/PascalCase` - Case conversions
- `Capitalize/LowerFirst/UpperFirst` - Capitalization
- `EndsWith/StartsWith` - String testing
- `Pad/PadEnd/PadStart` - String padding
- `Repeat` - String repetition
- `Split/Words` - String splitting and word extraction
- `Trim/TrimEnd/TrimStart` - Whitespace removal
- `Truncate` - String truncation with custom omission

### Object Package (`object`) - 65.6% Coverage

**15+ functions** for working with objects and maps:

- `Assign/Defaults` - Combine objects with different strategies
- `FromPairs/ToPairs` - Convert between arrays and objects
- `Get/Set` - Path-based property access
- `Has` - Check key existence
- `Invert` - Swap keys and values
- `IsEmpty` - Check if empty
- `Keys/Values` - Extract keys/values
- `MapKeys/MapValues` - Transform keys/values
- `Merge` - Deep merge objects
- `Omit/Pick` - Select/exclude properties

### Math Package (`math`) - 80.6% Coverage

**15+ mathematical utility functions**:

- `Abs/Clamp` - Value manipulation
- `Add/Subtract/Multiply/Divide` - Basic arithmetic
- `Ceil/Floor/Round` - Rounding functions
- `InRange/Random` - Range operations
- `IsInf/IsNaN` - Special value testing
- `Max/Min/Sum/Mean` - Aggregation functions
- `MaxBy/MinBy/SumBy/MeanBy` - Aggregation with iteratee
- `Pow/Sqrt` - Power and root functions

### Function Package (`function`) - 95.5% Coverage

**12+ utility functions** for working with functions:

- `After/Before` - Conditional execution
- `Compose/Pipe` - Function composition
- `Debounce/Throttle` - Rate limiting with args support
- `Delay/Defer` - Delayed execution with args support
- `Memoize/MemoizeWithResolver` - Result caching
- `Negate` - Predicate negation
- `Once/OnceVoid` - Single execution

### Lang Package (`lang`) - 85.7% Coverage

**10+ utility functions** for type checking and language constructs:

- `Clone/CloneDeep` - Shallow and deep cloning
- `IsArray/IsBoolean/IsDate/IsFunction` - Type checking
- `IsEmpty/IsEqual/IsNil` - Value testing
- `IsNumber/IsObject/IsString` - Type validation
- `ToArray/ToString` - Type conversion

### Util Package (`util`) - 90.9% Coverage

**15+ general utility functions**:

- `Attempt` - Safe function execution
- `Clamp/InRange` - Value range operations
- `Constant/Identity/Noop` - Function utilities
- `DefaultTo/DefaultToAny` - Default value handling
- `Flow/FlowRight` - Function composition
- `Random` - Random number generation
- `Range` - Number sequence generation
- `StubArray/StubFalse/StubObject/StubString/StubTrue` - Stub functions
- `Times` - Repeated function execution
- `ToPath` - Path string parsing
- `UniqueId` - Unique identifier generation

### Date Package (`date`) - 88.9% Coverage

**20+ utility functions** for date and time manipulation:

- `Add/Sub` - Date arithmetic
- `After/Before/Equal` - Date comparison
- `DaysInMonth/IsLeapYear` - Calendar utilities
- `EndOfDay/EndOfMonth/EndOfWeek/EndOfYear` - Period endings
- `Format` - Date formatting
- `IsDate/IsValid` - Date validation
- `Now` - Current timestamp
- `StartOfDay/StartOfMonth/StartOfWeek/StartOfYear` - Period beginnings
- `ToDate` - Date conversion

## Examples

### Working with Arrays

```go
import "github.com/go-libs/lodash/array"

// Chunk array into smaller arrays
chunks := array.Chunk([]int{1, 2, 3, 4, 5, 6}, 2)
// Result: [[1, 2], [3, 4], [5, 6]]

// Remove falsey values
clean := array.Compact([]interface{}{0, 1, false, 2, "", 3})
// Result: [1, 2, 3]

// Get unique values
unique := array.Uniq([]int{1, 2, 2, 3, 3, 3})
// Result: [1, 2, 3]
```

### Working with Collections

```go
import "github.com/go-libs/lodash/collection"

// Filter even numbers
evens := collection.Filter([]int{1, 2, 3, 4, 5}, func(x int) bool {
    return x%2 == 0
})
// Result: [2, 4]

// Transform values
doubled := collection.Map([]int{1, 2, 3}, func(x int) int {
    return x * 2
})
// Result: [2, 4, 6]

// Group by length
grouped := collection.GroupBy([]string{"one", "two", "three"}, func(s string) int {
    return len(s)
})
// Result: map[3:["one", "two"] 5:["three"]]
```

### Working with Strings

```go
import "github.com/go-libs/lodash/string"

// Convert to camelCase
camel := string.CamelCase("hello world")
// Result: "helloWorld"

// Convert to kebab-case
kebab := string.KebabCase("Hello World")
// Result: "hello-world"

// Truncate with custom omission
truncated := string.Truncate("A very long string", 10, "...")
// Result: "A very..."
```

## Performance

Go-Lodash is designed for high performance with minimal allocations:

- Uses generics for type safety without boxing
- Optimized algorithms for common operations
- Minimal memory allocations
- Thread-safe concurrent access

## Testing

Run all tests with coverage:

```bash
go test ./... -cover
```

Current test coverage:
- **Array**: 94.1% (20+ functions)
- **Collection**: 75.0% (20+ functions)
- **String**: 92.0% (15+ functions)
- **Object**: 65.6% (15+ functions)
- **Math**: 80.6% (15+ functions)
- **Function**: 95.5% (12+ functions)
- **Lang**: 85.7% (10+ functions)
- **Util**: 90.9% (15+ functions)
- **Date**: 88.9% (20+ functions)

**Overall**: 85%+ average coverage across all packages

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

MIT License - see LICENSE file for details.

## Inspiration

This library is inspired by [Lodash.js](https://lodash.com/) but designed specifically for Go's type system and idioms.
