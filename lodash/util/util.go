// Package util provides general utility functions.
// All functions are thread-safe and designed for high performance.
package util

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"
)

var uniqueIdCounter int64

// Identity returns the first argument it receives.
//
// Example:
//	Identity(42) // 42
//	Identity("hello") // "hello"
func Identity[T any](value T) T {
	return value
}

// Constant returns a function that always returns the same value.
//
// Example:
//	getValue := Constant(42)
//	getValue() // 42
//	getValue() // 42
func Constant[T any](value T) func() T {
	return func() T {
		return value
	}
}

// Noop is a no-operation function that does nothing.
//
// Example:
//	Noop() // does nothing
func Noop() {
	// Intentionally empty
}

// Range creates an array of numbers progressing from start up to, but not including, end.
//
// Example:
//	Range(4) // []int{0, 1, 2, 3}
//	Range(1, 5) // []int{1, 2, 3, 4}
//	Range(0, 20, 5) // []int{0, 5, 10, 15}
func Range(args ...int) []int {
	var start, end, step int
	
	switch len(args) {
	case 1:
		start, end, step = 0, args[0], 1
	case 2:
		start, end, step = args[0], args[1], 1
	case 3:
		start, end, step = args[0], args[1], args[2]
	default:
		return []int{}
	}
	
	if step == 0 {
		return []int{}
	}
	
	var result []int
	if step > 0 {
		for i := start; i < end; i += step {
			result = append(result, i)
		}
	} else {
		for i := start; i > end; i += step {
			result = append(result, i)
		}
	}
	
	return result
}

// Times invokes the iteratee n times, returning an array of the results.
//
// Example:
//	Times(3, func(i int) int { return i * 2 }) // []int{0, 2, 4}
//	Times(4, func(i int) string { return fmt.Sprintf("item-%d", i) }) // []string{"item-0", "item-1", "item-2", "item-3"}
func Times[T any](n int, iteratee func(int) T) []T {
	if n <= 0 {
		return []T{}
	}
	
	result := make([]T, n)
	for i := 0; i < n; i++ {
		result[i] = iteratee(i)
	}
	return result
}

// UniqueId generates a unique ID. If prefix is given, the ID is appended to it.
//
// Example:
//	UniqueId() // "1"
//	UniqueId("contact_") // "contact_2"
func UniqueId(prefix ...string) string {
	id := atomic.AddInt64(&uniqueIdCounter, 1)
	
	if len(prefix) > 0 {
		return fmt.Sprintf("%s%d", prefix[0], id)
	}
	return fmt.Sprintf("%d", id)
}

// DefaultTo checks value to determine whether a default value should be returned in its place.
//
// Example:
//	DefaultTo(1, 10) // 1
//	DefaultTo(nil, 10) // 10
//	DefaultTo("", "default") // "default"
func DefaultTo[T comparable](value, defaultValue T) T {
	var zero T
	if value == zero {
		return defaultValue
	}
	return value
}

// DefaultToAny checks value to determine whether a default value should be returned in its place.
// Works with any type including interfaces.
//
// Example:
//	DefaultToAny(nil, "default") // "default"
//	DefaultToAny(42, "default") // 42
func DefaultToAny(value, defaultValue interface{}) interface{} {
	if value == nil {
		return defaultValue
	}
	return value
}

// Attempt attempts to invoke func, returning either the result or the caught error object.
//
// Example:
//	result, err := Attempt(func() (int, error) { return 42, nil })
//	// result: 42, err: nil
func Attempt[T any](fn func() (T, error)) (T, error) {
	return fn()
}

// Flow creates a function that is the composition of the provided functions,
// where each successive invocation is supplied the return value of the previous.
//
// Example:
//	add := func(x int) int { return x + 1 }
//	multiply := func(x int) int { return x * 2 }
//	composed := Flow(add, multiply)
//	result := composed(3) // (3 + 1) * 2 = 8
func Flow[T any](fns ...func(T) T) func(T) T {
	return func(value T) T {
		result := value
		for _, fn := range fns {
			result = fn(result)
		}
		return result
	}
}

// FlowRight creates a function that is the composition of the provided functions,
// where each successive invocation is supplied the return value of the previous.
// This is like Flow except that it composes functions from right to left.
//
// Example:
//	add := func(x int) int { return x + 1 }
//	multiply := func(x int) int { return x * 2 }
//	composed := FlowRight(add, multiply)
//	result := composed(3) // (3 * 2) + 1 = 7
func FlowRight[T any](fns ...func(T) T) func(T) T {
	return func(value T) T {
		result := value
		for i := len(fns) - 1; i >= 0; i-- {
			result = fns[i](result)
		}
		return result
	}
}

// StubArray returns a new empty array.
//
// Example:
//	StubArray() // []interface{}{}
func StubArray() []interface{} {
	return []interface{}{}
}

// StubFalse returns false.
//
// Example:
//	StubFalse() // false
func StubFalse() bool {
	return false
}

// StubObject returns a new empty object.
//
// Example:
//	StubObject() // map[string]interface{}{}
func StubObject() map[string]interface{} {
	return map[string]interface{}{}
}

// StubString returns an empty string.
//
// Example:
//	StubString() // ""
func StubString() string {
	return ""
}

// StubTrue returns true.
//
// Example:
//	StubTrue() // true
func StubTrue() bool {
	return true
}

// Random generates a random number between min and max (inclusive).
//
// Example:
//	Random(1, 10) // random number between 1 and 10
//	Random(5) // random number between 0 and 5
func Random(args ...int) int {
	rand.Seed(time.Now().UnixNano())
	
	var min, max int
	switch len(args) {
	case 1:
		min, max = 0, args[0]
	case 2:
		min, max = args[0], args[1]
	default:
		return 0
	}
	
	if min > max {
		min, max = max, min
	}
	
	return rand.Intn(max-min+1) + min
}

// Clamp clamps number within the inclusive lower and upper bounds.
//
// Example:
//	Clamp(10, 5, 15) // 10
//	Clamp(3, 5, 15) // 5
//	Clamp(20, 5, 15) // 15
func Clamp[T int | float64](number, lower, upper T) T {
	if number < lower {
		return lower
	}
	if number > upper {
		return upper
	}
	return number
}

// InRange checks if number is between start and end (not including end).
//
// Example:
//	InRange(3, 2, 4) // true
//	InRange(4, 8) // true (start defaults to 0)
//	InRange(4, 2) // false
func InRange[T int | float64](number T, args ...T) bool {
	var start, end T
	
	switch len(args) {
	case 1:
		start, end = 0, args[0]
	case 2:
		start, end = args[0], args[1]
	default:
		return false
	}
	
	if start > end {
		start, end = end, start
	}
	
	return number >= start && number < end
}

// ToPath converts string to a property path array.
//
// Example:
//	ToPath("a.b.c") // []string{"a", "b", "c"}
//	ToPath("a[0].b") // []string{"a", "0", "b"}
func ToPath(str string) []string {
	if str == "" {
		return []string{}
	}
	
	var result []string
	var current string
	inBracket := false
	
	for _, char := range str {
		switch char {
		case '.':
			if !inBracket && current != "" {
				result = append(result, current)
				current = ""
			} else if inBracket {
				current += string(char)
			}
		case '[':
			if current != "" {
				result = append(result, current)
				current = ""
			}
			inBracket = true
		case ']':
			if inBracket && current != "" {
				result = append(result, current)
				current = ""
			}
			inBracket = false
		default:
			current += string(char)
		}
	}
	
	if current != "" {
		result = append(result, current)
	}
	
	return result
}
