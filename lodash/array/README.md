# Array Package

A comprehensive collection of array manipulation utilities for Go, inspired by Lodash.js. This package provides 61 high-performance, thread-safe functions for working with slices and arrays.

## Features

- **🚀 High Performance**: Optimized algorithms with minimal allocations
- **🔒 Thread Safe**: All functions are safe for concurrent use
- **🎯 Type Safe**: Full Go generics support for type safety
- **📦 Zero Dependencies**: No external dependencies
- **✅ Well Tested**: 85.9% test coverage with comprehensive edge cases

## Installation

```bash
go get github.com/nguyendkn/go-libs/lodash/array
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/nguyendkn/go-libs/lodash/array"
)

func main() {
    // Chunk array into smaller arrays
    numbers := []int{1, 2, 3, 4, 5, 6}
    chunks := array.Chunk(numbers, 2)
    fmt.Println(chunks) // [[1 2] [3 4] [5 6]]

    // Remove falsy values
    mixed := []interface{}{0, 1, false, 2, "", 3, nil}
    clean := array.Compact(mixed)
    fmt.Println(clean) // [1 2 3]

    // Get unique values
    duplicates := []int{1, 2, 2, 3, 3, 3}
    unique := array.Uniq(duplicates)
    fmt.Println(unique) // [1 2 3]
}
```

## Function Categories

### 🔧 **Basic Operations**
- **`Chunk`** - Split array into chunks of specified size
- **`Compact`** - Remove falsy values from array
- **`Concat`** - Concatenate arrays together
- **`Fill`** - Fill array elements with value
- **`Flatten`** - Flatten array one level deep
- **`FlattenDeep`** - Recursively flatten array
- **`FlattenDepth`** - Flatten array up to specified depth

### 🔍 **Search & Access**
- **`Head`** - Get first element
- **`Last`** - Get last element
- **`Initial`** - Get all but last element
- **`Tail`** - Get all but first element
- **`Nth`** - Get element at index
- **`IndexOf`** - Find index of element
- **`LastIndexOf`** - Find last index of element

### ✂️ **Slicing & Extraction**
- **`Drop`** - Drop n elements from beginning
- **`DropRight`** - Drop n elements from end
- **`Take`** - Take n elements from beginning
- **`TakeRight`** - Take n elements from end
- **`TakeWhile`** - Take elements while predicate is true
- **`TakeRightWhile`** - Take elements from end while predicate is true

### 🔄 **Set Operations**
- **`Union`** - Create array of unique values from all arrays
- **`UnionBy`** - Union with iteratee for comparison
- **`UnionWith`** - Union with comparator function
- **`Intersection`** - Create array of shared values
- **`IntersectionBy`** - Intersection with iteratee
- **`IntersectionWith`** - Intersection with comparator
- **`Difference`** - Create array excluding values from other arrays
- **`Xor`** - Create array of symmetric difference
- **`XorBy`** - Symmetric difference with iteratee
- **`XorWith`** - Symmetric difference with comparator

### 🎯 **Unique Operations**
- **`Uniq`** - Create duplicate-free array
- **`UniqBy`** - Unique values with iteratee
- **`UniqWith`** - Unique values with comparator
- **`SortedUniq`** - Unique values from sorted array
- **`SortedUniqBy`** - Sorted unique with iteratee

### 🗑️ **Removal Operations**
- **`Pull`** - Remove specified values
- **`PullAll`** - Remove all values from array
- **`PullAllBy`** - Remove values with iteratee
- **`PullAllWith`** - Remove values with comparator
- **`PullAt`** - Remove elements at indexes
- **`Remove`** - Remove elements matching predicate
- **`Without`** - Create array excluding specified values

### 📊 **Sorted Array Operations**
- **`SortedIndex`** - Get insertion index for sorted array
- **`SortedIndexBy`** - Sorted index with iteratee
- **`SortedIndexOf`** - Find index in sorted array
- **`SortedLastIndex`** - Get last insertion index
- **`SortedLastIndexBy`** - Last sorted index with iteratee
- **`SortedLastIndexOf`** - Find last index in sorted array

### 🔗 **Transformation Operations**
- **`Zip`** - Create array of grouped elements
- **`ZipObject`** - Create object from arrays
- **`ZipObjectDeep`** - Create nested object from arrays
- **`ZipWith`** - Zip arrays with iteratee
- **`Unzip`** - Unzip grouped arrays
- **`UnzipWith`** - Unzip with iteratee
- **`FromPairs`** - Create object from key-value pairs
- **`FromPairsString`** - Create string-keyed object from key-value pairs

### 🔄 **Utility Operations**
- **`Reverse`** - Reverse array in place
- **`Join`** - Join array elements into string
- **`Slice`** - Extract slice of array

## Detailed Examples

### Working with Chunks
```go
// Split large dataset into batches
data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
batches := array.Chunk(data, 3)
// Result: [[1 2 3] [4 5 6] [7 8 9] [10]]

for i, batch := range batches {
    fmt.Printf("Batch %d: %v\n", i+1, batch)
}
```

### Set Operations
```go
// Find common elements
arr1 := []int{1, 2, 3, 4}
arr2 := []int{3, 4, 5, 6}
arr3 := []int{4, 5, 6, 7}

common := array.Intersection(arr1, arr2, arr3)
fmt.Println(common) // [4]

// Get unique elements from all arrays
union := array.Union(arr1, arr2, arr3)
fmt.Println(union) // [1 2 3 4 5 6 7]

// Symmetric difference
xor := array.Xor(arr1, arr2)
fmt.Println(xor) // [1 2 5 6]
```

### Advanced Filtering
```go
// Remove elements by predicate
numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// Take while condition is true
lessThanFive := array.TakeWhile(numbers, func(n int) bool {
    return n < 5
})
fmt.Println(lessThanFive) // [1 2 3 4]

// Take from right while condition is true
greaterThanSix := array.TakeRightWhile(numbers, func(n int) bool {
    return n > 6
})
fmt.Println(greaterThanSix) // [7 8 9 10]
```

### Working with Sorted Arrays
```go
// Efficient operations on sorted arrays
sorted := []int{1, 3, 5, 7, 9}

// Find insertion point
index := array.SortedIndex(sorted, 6)
fmt.Println(index) // 3

// Get unique values from sorted array (more efficient)
sortedWithDups := []int{1, 1, 2, 2, 3, 3}
unique := array.SortedUniq(sortedWithDups)
fmt.Println(unique) // [1 2 3]
```

## Performance Notes

- **Memory Efficient**: Functions avoid unnecessary allocations where possible
- **Optimized Algorithms**: Uses efficient algorithms like binary search for sorted arrays
- **Generic Support**: Type-safe operations without boxing/unboxing overhead
- **Concurrent Safe**: All functions are safe for concurrent use

## Benchmarks

```
BenchmarkChunk-8           1000000    1043 ns/op    240 B/op    3 allocs/op
BenchmarkUniq-8             500000    2891 ns/op    512 B/op    5 allocs/op
BenchmarkIntersection-8     300000    4123 ns/op    384 B/op    4 allocs/op
BenchmarkSortedUniq-8      2000000     654 ns/op    128 B/op    1 allocs/op
```

## Error Handling

Functions handle edge cases gracefully:
- Empty arrays return appropriate zero values
- Out-of-bounds access returns zero values with boolean indicators
- Nil-safe operations where applicable

## Thread Safety

All functions are thread-safe and can be used concurrently without additional synchronization.

## Contributing

See the main [lodash README](../README.md) for contribution guidelines.

## License

This package is part of the go-libs project and follows the same license terms.