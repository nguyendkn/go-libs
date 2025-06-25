// Package math provides mathematical utility functions.
// All functions are thread-safe and designed for high performance.
package math

import (
	"math"
	"math/rand"
	"time"
)

// Numeric represents types that can be used in mathematical operations
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Add adds two numbers.
//
// Example:
//	Add(6, 4) // 10
//	Add(1.5, 2.5) // 4.0
func Add[T Numeric](a, b T) T {
	return a + b
}

// Subtract subtracts the second number from the first.
//
// Example:
//	Subtract(6, 4) // 2
//	Subtract(10.5, 2.5) // 8.0
func Subtract[T Numeric](a, b T) T {
	return a - b
}

// Multiply multiplies two numbers.
//
// Example:
//	Multiply(6, 4) // 24
//	Multiply(2.5, 4) // 10.0
func Multiply[T Numeric](a, b T) T {
	return a * b
}

// Divide divides the first number by the second.
//
// Example:
//	Divide(6, 4) // 1.5
//	Divide(10.0, 2.0) // 5.0
func Divide[T Numeric](a, b T) T {
	return a / b
}

// Max returns the maximum value from a slice of numbers.
//
// Example:
//	Max([]int{4, 2, 8, 6}) // 8
//	Max([]float64{1.5, 3.2, 2.1}) // 3.2
func Max[T Numeric](numbers []T) (T, bool) {
	if len(numbers) == 0 {
		var zero T
		return zero, false
	}
	
	max := numbers[0]
	for _, num := range numbers[1:] {
		if num > max {
			max = num
		}
	}
	return max, true
}

// Min returns the minimum value from a slice of numbers.
//
// Example:
//	Min([]int{4, 2, 8, 6}) // 2
//	Min([]float64{1.5, 3.2, 2.1}) // 1.5
func Min[T Numeric](numbers []T) (T, bool) {
	if len(numbers) == 0 {
		var zero T
		return zero, false
	}
	
	min := numbers[0]
	for _, num := range numbers[1:] {
		if num < min {
			min = num
		}
	}
	return min, true
}

// Sum calculates the sum of all numbers in a slice.
//
// Example:
//	Sum([]int{4, 2, 8, 6}) // 20
//	Sum([]float64{1.5, 2.5, 3.0}) // 7.0
func Sum[T Numeric](numbers []T) T {
	var sum T
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// Mean calculates the arithmetic mean of all numbers in a slice.
//
// Example:
//	Mean([]int{4, 2, 8, 6}) // 5.0
//	Mean([]float64{1.0, 2.0, 3.0}) // 2.0
func Mean[T Numeric](numbers []T) (float64, bool) {
	if len(numbers) == 0 {
		return 0, false
	}
	
	sum := Sum(numbers)
	return float64(sum) / float64(len(numbers)), true
}

// Ceil computes the ceiling of a number (rounds up to the nearest integer).
//
// Example:
//	Ceil(4.006) // 5.0
//	Ceil(6.004) // 7.0
//	Ceil(6) // 6.0
func Ceil(n float64) float64 {
	return math.Ceil(n)
}

// Floor computes the floor of a number (rounds down to the nearest integer).
//
// Example:
//	Floor(4.006) // 4.0
//	Floor(0.046) // 0.0
//	Floor(4) // 4.0
func Floor(n float64) float64 {
	return math.Floor(n)
}

// Round rounds a number to the nearest integer.
//
// Example:
//	Round(4.006) // 4.0
//	Round(4.6) // 5.0
//	Round(-4.6) // -5.0
func Round(n float64) float64 {
	return math.Round(n)
}

// Abs returns the absolute value of a number.
//
// Example:
//	Abs(-5) // 5
//	Abs(5) // 5
//	Abs(-3.14) // 3.14
func Abs[T Numeric](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

// Clamp clamps a number within the inclusive lower and upper bounds.
//
// Example:
//	Clamp(-10, -5, 5) // -5
//	Clamp(10, -5, 5) // 5
//	Clamp(3, -5, 5) // 3
func Clamp[T Numeric](number, lower, upper T) T {
	if number < lower {
		return lower
	}
	if number > upper {
		return upper
	}
	return number
}

// InRange checks if a number is between start and end (not including end).
//
// Example:
//	InRange(3, 2, 4) // true
//	InRange(4, 8) // true (start defaults to 0)
//	InRange(4, 2) // false
//	InRange(2, 2) // false
func InRange[T Numeric](number T, args ...T) bool {
	var start, end T
	
	switch len(args) {
	case 1:
		start = 0
		end = args[0]
	case 2:
		start = args[0]
		end = args[1]
	default:
		return false
	}
	
	if start > end {
		start, end = end, start
	}
	
	return number >= start && number < end
}

// Random produces a random number between the inclusive lower and upper bounds.
//
// Example:
//	Random(0, 5) // random number between 0 and 5
//	Random(1.2, 5.2) // random float between 1.2 and 5.2
//	Random(5) // random number between 0 and 5
func Random[T Numeric](args ...T) T {
	rand.Seed(time.Now().UnixNano())
	
	var lower, upper T
	
	switch len(args) {
	case 1:
		lower = 0
		upper = args[0]
	case 2:
		lower = args[0]
		upper = args[1]
	default:
		return 0
	}
	
	if lower > upper {
		lower, upper = upper, lower
	}
	
	// Handle different numeric types
	switch any(lower).(type) {
	case float32, float64:
		diff := float64(upper - lower)
		return T(float64(lower) + rand.Float64()*diff)
	default:
		diff := int64(upper - lower + 1)
		return T(int64(lower) + rand.Int63n(diff))
	}
}

// MaxBy returns the maximum value from a slice using an iteratee function.
//
// Example:
//	MaxBy([]string{"a", "bb", "ccc"}, func(s string) int { return len(s) }) // "ccc"
func MaxBy[T any, R Numeric](slice []T, iteratee func(T) R) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	
	maxItem := slice[0]
	maxValue := iteratee(slice[0])
	
	for _, item := range slice[1:] {
		value := iteratee(item)
		if value > maxValue {
			maxValue = value
			maxItem = item
		}
	}
	
	return maxItem, true
}

// MinBy returns the minimum value from a slice using an iteratee function.
//
// Example:
//	MinBy([]string{"a", "bb", "ccc"}, func(s string) int { return len(s) }) // "a"
func MinBy[T any, R Numeric](slice []T, iteratee func(T) R) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	
	minItem := slice[0]
	minValue := iteratee(slice[0])
	
	for _, item := range slice[1:] {
		value := iteratee(item)
		if value < minValue {
			minValue = value
			minItem = item
		}
	}
	
	return minItem, true
}

// SumBy calculates the sum of all values in a slice using an iteratee function.
//
// Example:
//	SumBy([]string{"a", "bb", "ccc"}, func(s string) int { return len(s) }) // 6
func SumBy[T any, R Numeric](slice []T, iteratee func(T) R) R {
	var sum R
	for _, item := range slice {
		sum += iteratee(item)
	}
	return sum
}

// MeanBy calculates the arithmetic mean of all values in a slice using an iteratee function.
//
// Example:
//	MeanBy([]string{"a", "bb", "ccc"}, func(s string) int { return len(s) }) // 2.0
func MeanBy[T any, R Numeric](slice []T, iteratee func(T) R) (float64, bool) {
	if len(slice) == 0 {
		return 0, false
	}
	
	sum := SumBy(slice, iteratee)
	return float64(sum) / float64(len(slice)), true
}

// Pow returns base raised to the power of exponent.
//
// Example:
//	Pow(2, 3) // 8.0
//	Pow(4, 0.5) // 2.0
func Pow(base, exponent float64) float64 {
	return math.Pow(base, exponent)
}

// Sqrt returns the square root of a number.
//
// Example:
//	Sqrt(9) // 3.0
//	Sqrt(2) // 1.4142135623730951
func Sqrt(n float64) float64 {
	return math.Sqrt(n)
}

// IsNaN reports whether f is an IEEE 754 "not-a-number" value.
//
// Example:
//	IsNaN(math.NaN()) // true
//	IsNaN(1.0) // false
func IsNaN(f float64) bool {
	return math.IsNaN(f)
}

// IsInf reports whether f is an infinity.
//
// Example:
//	IsInf(math.Inf(1), 0) // true
//	IsInf(1.0, 0) // false
func IsInf(f float64, sign int) bool {
	return math.IsInf(f, sign)
}
