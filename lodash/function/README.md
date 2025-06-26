# Function Package

Advanced function manipulation utilities for Go, providing powerful tools for functional programming patterns. This package offers 28 high-performance, thread-safe functions for controlling function execution, composition, and transformation.

## Features

- **🚀 High Performance**: Optimized algorithms with minimal overhead
- **🔒 Thread Safe**: All functions are safe for concurrent use
- **🎯 Type Safe**: Full Go generics support for type safety
- **📦 Zero Dependencies**: No external dependencies beyond Go standard library
- **✅ Well Tested**: Comprehensive test coverage with edge cases
- **⚡ Async Support**: Built-in support for asynchronous operations

## Installation

```bash
go get github.com/nguyendkn/go-libs/lodash/function
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/nguyendkn/go-libs/lodash/function"
)

func main() {
    // Debounce function calls
    debounced := function.Debounce(func() {
        fmt.Println("Debounced call")
    }, 100*time.Millisecond)

    // Memoize expensive computations
    fibonacci := function.Memoize(func(n int) int {
        if n <= 1 { return n }
        return fibonacci(n-1) + fibonacci(n-2)
    })

    // Execute function only once
    initialize := function.Once(func() string {
        fmt.Println("Initializing...")
        return "initialized"
    })

    result1 := initialize() // Prints and returns "initialized"
    result2 := initialize() // Only returns "initialized"
}
```

## Core Functions

### ⏱️ **Timing Control**
- **`Debounce`** - Delay function execution until after wait time
- **`DebounceWithArgs`** - Debounce with arguments support
- **`Throttle`** - Limit function execution frequency
- **`ThrottleWithArgs`** - Throttle with arguments support
- **`Delay`** - Execute function after specified delay
- **`DelayWithArgs`** - Delay with arguments support
- **`Defer`** - Execute function on next tick
- **`DeferWithArgs`** - Defer with arguments support

### 🔄 **Execution Control**
- **`Once`** - Execute function only once (with return value)
- **`OnceVoid`** - Execute function only once (no return value)
- **`After`** - Execute function after n calls
- **`Before`** - Execute function before n calls
- **`Ary`** - Limit function to n arguments

### 💾 **Memoization**
- **`Memoize`** - Cache function results by arguments
- **`MemoizeWithResolver`** - Memoize with custom key resolver

### 🔧 **Function Composition**
- **`Compose`** - Compose functions right to left
- **`Pipe`** - Compose functions left to right
- **`Negate`** - Create negated predicate function

### 🍛 **Currying & Partial Application**
- **`Curry2`** - Curry function with 2 arguments
- **`Curry3`** - Curry function with 3 arguments
- **`Curry4`** - Curry function with 4 arguments
- **`Partial2`** - Partial application for 2-argument functions
- **`Partial3`** - Partial application for 3-argument functions
- **`Partial4`** - Partial application for 4-argument functions

### 🔄 **Argument Manipulation**
- **`Flip2`** - Flip arguments of 2-argument function
- **`Flip3`** - Flip arguments of 3-argument function
- **`Rearg2`** - Reorder arguments of 2-argument function
- **`Unary2`** - Convert 2-argument function to unary

## Detailed Examples

### Debouncing and Throttling
```go
// Debounce - useful for search input, resize events
searchDebounced := function.Debounce(func() {
    fmt.Println("Performing search...")
}, 300*time.Millisecond)

// Multiple rapid calls will only execute once after 300ms
searchDebounced()
searchDebounced()
searchDebounced()

// Throttle - useful for scroll events, API rate limiting
scrollThrottled := function.Throttle(func() {
    fmt.Println("Handling scroll...")
}, 100*time.Millisecond)

// Will execute immediately, then ignore calls for 100ms
scrollThrottled()
scrollThrottled() // Ignored
scrollThrottled() // Ignored
```

### Memoization for Performance
```go
// Expensive computation with memoization
expensiveCalc := function.Memoize(func(n int) int {
    fmt.Printf("Computing for %d...\n", n)
    time.Sleep(100 * time.Millisecond) // Simulate expensive operation
    return n * n * n
})

result1 := expensiveCalc(5) // Computes and caches
result2 := expensiveCalc(5) // Returns cached result instantly

// Custom key resolver for complex arguments
import "math"

type Point struct{ X, Y int }
distance := function.MemoizeWithResolver(
    func(p Point) float64 {
        return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
    },
    func(p Point) string {
        return fmt.Sprintf("%d,%d", p.X, p.Y)
    },
)
```

### Function Composition
```go
// Compose functions (right to left)
addOne := func(x int) int { return x + 1 }
double := func(x int) int { return x * 2 }
square := func(x int) int { return x * x }

// Compose: square(double(addOne(x)))
composed := function.Compose(square, double, addOne)
result := composed(3) // (3+1)*2 = 8, 8^2 = 64

// Pipe functions (left to right)
piped := function.Pipe(addOne, double, square)
result2 := piped(3) // Same result: 64

// Negate predicate
isEven := func(n int) bool { return n%2 == 0 }
isOdd := function.Negate(isEven)
fmt.Println(isOdd(3)) // true
```

### Currying and Partial Application
```go
// Currying
add := func(a, b int) int { return a + b }
curriedAdd := function.Curry2(add)
addFive := curriedAdd(5)
result := addFive(3) // 8

// Partial application
greet := func(greeting, name string) string {
    return greeting + " " + name
}
sayHello := function.Partial2(greet, "Hello")
message := sayHello("World") // "Hello World"

// Flip arguments
divide := func(a, b float64) float64 { return a / b }
flippedDivide := function.Flip2(divide)
result := flippedDivide(2, 10) // 10 / 2 = 5
```

### Execution Control
```go
// Execute only once
counter := 0
increment := function.Once(func() int {
    counter++
    return counter
})

val1 := increment() // 1
val2 := increment() // 1 (same result)

// Execute after n calls
afterThree := function.After(3, func() {
    fmt.Println("Called after 3 attempts!")
})

afterThree() // Nothing
afterThree() // Nothing
afterThree() // Prints message

// Execute before n calls
beforeThree := function.Before(3, func() string {
    return "Available"
})

result1 := beforeThree() // "Available"
result2 := beforeThree() // "Available"
result3 := beforeThree() // "Available" (last time)
result4 := beforeThree() // "Available" (returns last result)
```

## Performance Notes

- **Memory Efficient**: Minimal overhead for function wrapping
- **Thread Safe**: All functions use proper synchronization
- **Generic Optimized**: Type-safe operations without boxing
- **Concurrent Safe**: Safe for use in goroutines

## Use Cases

- **Debouncing**: Search inputs, resize handlers, API calls
- **Throttling**: Scroll events, rate limiting, performance optimization
- **Memoization**: Expensive computations, recursive algorithms
- **Composition**: Data transformation pipelines, functional programming
- **Currying**: Configuration functions, partial application patterns

## Thread Safety

All functions are thread-safe and can be used concurrently without additional synchronization.

## Contributing

See the main [lodash README](../README.md) for contribution guidelines.

## License

This package is part of the go-libs project and follows the same license terms.