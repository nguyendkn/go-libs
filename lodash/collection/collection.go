// Package collection provides utility functions for working with collections (slices, maps).
// All functions are thread-safe and designed for high performance.
package collection

import (
	"math/rand"
	"reflect"
	"sort"
	"time"
)

// Filter creates a new slice with all elements that pass the test implemented by the provided function.
//
// Example:
//
//	Filter([]int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 }) // []int{2, 4}
func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// Map creates a new slice with the results of calling a provided function on every element.
//
// Example:
//
//	Map([]int{1, 2, 3}, func(x int) int { return x * 2 }) // []int{2, 4, 6}
func Map[T, R any](slice []T, mapper func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = mapper(item)
	}
	return result
}

// Reduce executes a reducer function on each element of the slice, resulting in a single output value.
//
// Example:
//
//	Reduce([]int{1, 2, 3, 4}, func(acc, x int) int { return acc + x }, 0) // 10
func Reduce[T, R any](slice []T, reducer func(R, T) R, initial R) R {
	result := initial
	for _, item := range slice {
		result = reducer(result, item)
	}
	return result
}

// Find returns the first element in the slice that satisfies the provided testing function.
//
// Example:
//
//	Find([]int{1, 2, 3, 4}, func(x int) bool { return x > 2 }) // 3, true
//	Find([]int{1, 2}, func(x int) bool { return x > 5 }) // 0, false
func Find[T any](slice []T, predicate func(T) bool) (T, bool) {
	for _, item := range slice {
		if predicate(item) {
			return item, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element in the slice that satisfies the provided testing function.
//
// Example:
//
//	FindIndex([]int{1, 2, 3, 4}, func(x int) bool { return x > 2 }) // 2, true
//	FindIndex([]int{1, 2}, func(x int) bool { return x > 5 }) // -1, false
func FindIndex[T any](slice []T, predicate func(T) bool) (int, bool) {
	for i, item := range slice {
		if predicate(item) {
			return i, true
		}
	}
	return -1, false
}

// Every tests whether all elements in the slice pass the test implemented by the provided function.
//
// Example:
//
//	Every([]int{2, 4, 6}, func(x int) bool { return x%2 == 0 }) // true
//	Every([]int{1, 2, 3}, func(x int) bool { return x%2 == 0 }) // false
func Every[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if !predicate(item) {
			return false
		}
	}
	return true
}

// Some tests whether at least one element in the slice passes the test implemented by the provided function.
//
// Example:
//
//	Some([]int{1, 2, 3}, func(x int) bool { return x%2 == 0 }) // true
//	Some([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 }) // false
func Some[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// ForEach executes a provided function once for each slice element.
//
// Example:
//
//	ForEach([]int{1, 2, 3}, func(x int) { fmt.Println(x) })
func ForEach[T any](slice []T, fn func(T)) {
	for _, item := range slice {
		fn(item)
	}
}

// ForEachWithIndex executes a provided function once for each slice element with index.
//
// Example:
//
//	ForEachWithIndex([]string{"a", "b"}, func(i int, x string) { fmt.Printf("%d: %s\n", i, x) })
func ForEachWithIndex[T any](slice []T, fn func(int, T)) {
	for i, item := range slice {
		fn(i, item)
	}
}

// GroupBy groups the elements of the slice by the result of the provided function.
//
// Example:
//
//	GroupBy([]string{"one", "two", "three"}, func(s string) int { return len(s) })
//	// map[int][]string{3: {"one", "two"}, 5: {"three"}}
func GroupBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K][]T {
	result := make(map[K][]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = append(result[key], item)
	}
	return result
}

// CountBy counts the elements of the slice by the result of the provided function.
//
// Example:
//
//	CountBy([]string{"one", "two", "three"}, func(s string) int { return len(s) })
//	// map[int]int{3: 2, 5: 1}
func CountBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K]int {
	result := make(map[K]int)
	for _, item := range slice {
		key := keyFunc(item)
		result[key]++
	}
	return result
}

// Partition creates two slices: one with elements that pass the predicate and one with elements that don't.
//
// Example:
//
//	Partition([]int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 })
//	// []int{2, 4}, []int{1, 3}
func Partition[T any](slice []T, predicate func(T) bool) ([]T, []T) {
	var truthy, falsy []T
	for _, item := range slice {
		if predicate(item) {
			truthy = append(truthy, item)
		} else {
			falsy = append(falsy, item)
		}
	}
	return truthy, falsy
}

// Sample returns a random element from the slice.
//
// Example:
//
//	Sample([]int{1, 2, 3, 4}) // random element, true
//	Sample([]int{}) // 0, false
func Sample[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(slice))
	return slice[index], true
}

// SampleSize returns n random elements from the slice.
//
// Example:
//
//	SampleSize([]int{1, 2, 3, 4, 5}, 3) // []int{2, 4, 1} (random order)
func SampleSize[T any](slice []T, n int) []T {
	if n <= 0 || len(slice) == 0 {
		return []T{}
	}

	if n >= len(slice) {
		// Return a shuffled copy of the entire slice
		result := make([]T, len(slice))
		copy(result, slice)
		Shuffle(result)
		return result
	}

	rand.Seed(time.Now().UnixNano())
	indices := rand.Perm(len(slice))[:n]
	result := make([]T, n)
	for i, idx := range indices {
		result[i] = slice[idx]
	}
	return result
}

// Shuffle randomly shuffles the elements of the slice in place.
//
// Example:
//
//	arr := []int{1, 2, 3, 4}
//	Shuffle(arr) // arr becomes randomly shuffled
func Shuffle[T any](slice []T) {
	rand.Seed(time.Now().UnixNano())
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Size returns the length of the slice.
//
// Example:
//
//	Size([]int{1, 2, 3}) // 3
//	Size([]string{}) // 0
func Size[T any](slice []T) int {
	return len(slice)
}

// Includes checks if a value is in the slice.
//
// Example:
//
//	Includes([]int{1, 2, 3}, 2) // true
//	Includes([]int{1, 2, 3}, 4) // false
func Includes[T comparable](slice []T, value T) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// Reject creates a new slice with all elements that do not pass the test implemented by the provided function.
// It's the opposite of Filter.
//
// Example:
//
//	Reject([]int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 }) // []int{1, 3}
func Reject[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range slice {
		if !predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

// FindLast iterates over elements of collection from right to left and returns the first element predicate returns truthy for.
//
// Example:
//
//	FindLast([]int{1, 2, 3, 4}, func(x int) bool { return x%2 == 0 }) // 4, true
//	FindLast([]int{1, 3, 5}, func(x int) bool { return x%2 == 0 }) // 0, false
func FindLast[T any](slice []T, predicate func(T) bool) (T, bool) {
	for i := len(slice) - 1; i >= 0; i-- {
		if predicate(slice[i]) {
			return slice[i], true
		}
	}
	var zero T
	return zero, false
}

// FlatMap creates a flattened array of values by running each element in collection through iteratee.
//
// Example:
//
//	FlatMap([][]int{{1, 2}, {3, 4}}, func(x []int) []int { return x }) // []int{1, 2, 3, 4}
func FlatMap[T, R any](slice []T, mapper func(T) []R) []R {
	var result []R
	for _, item := range slice {
		mapped := mapper(item)
		result = append(result, mapped...)
	}
	return result
}

// ForEachRight iterates over elements of collection from right to left and invokes iteratee for each element.
//
// Example:
//
//	ForEachRight([]int{1, 2, 3}, func(x int) { fmt.Println(x) }) // prints 3, 2, 1
func ForEachRight[T any](slice []T, fn func(T)) {
	for i := len(slice) - 1; i >= 0; i-- {
		fn(slice[i])
	}
}

// KeyBy creates an object composed of keys generated from the results of running each element through iteratee.
//
// Example:
//
//	KeyBy([]string{"a", "bb", "ccc"}, func(s string) int { return len(s) }) // map[int]string{1: "a", 2: "bb", 3: "ccc"}
func KeyBy[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	for _, item := range slice {
		key := keyFunc(item)
		result[key] = item
	}
	return result
}

// OrderBy sorts slice by multiple criteria. Each criterion is defined by an iteratee function and sort order.
//
// Example:
//
//	type Person struct { Name string; Age int }
//	people := []Person{{"John", 30}, {"Jane", 25}, {"Bob", 30}}
//	OrderBy(people, []func(Person) interface{}{
//		func(p Person) interface{} { return p.Age },
//		func(p Person) interface{} { return p.Name },
//	}, []bool{true, false}) // Sort by age asc, then name desc
func OrderBy[T any](slice []T, iteratees []func(T) interface{}, orders []bool) []T {
	if len(iteratees) == 0 {
		return append([]T{}, slice...)
	}

	result := append([]T{}, slice...)

	// Use stable sort for multiple criteria
	for i := len(iteratees) - 1; i >= 0; i-- {
		iteratee := iteratees[i]
		ascending := true
		if i < len(orders) {
			ascending = orders[i]
		}

		sort.SliceStable(result, func(a, b int) bool {
			valA := iteratee(result[a])
			valB := iteratee(result[b])

			// Compare based on type
			switch vA := valA.(type) {
			case int:
				vB := valB.(int)
				if ascending {
					return vA < vB
				}
				return vA > vB
			case string:
				vB := valB.(string)
				if ascending {
					return vA < vB
				}
				return vA > vB
			case float64:
				vB := valB.(float64)
				if ascending {
					return vA < vB
				}
				return vA > vB
			default:
				return false
			}
		})
	}

	return result
}

// ReduceRight reduces collection to a value which is the accumulated result of running each element from right to left.
//
// Example:
//
//	ReduceRight([]string{"a", "b", "c"}, func(acc, x string) string { return acc + x }, "") // "cba"
func ReduceRight[T, R any](slice []T, reducer func(R, T) R, initial R) R {
	result := initial
	for i := len(slice) - 1; i >= 0; i-- {
		result = reducer(result, slice[i])
	}
	return result
}

// SortBy creates an array of elements, sorted in ascending order by the results of running each element through iteratee.
//
// Example:
//
//	SortBy([]string{"banana", "apple", "cherry"}, func(s string) int { return len(s) }) // ["apple", "banana", "cherry"]
func SortBy[T any](slice []T, iteratee func(T) interface{}) []T {
	result := append([]T{}, slice...)

	sort.SliceStable(result, func(i, j int) bool {
		valI := iteratee(result[i])
		valJ := iteratee(result[j])

		switch vI := valI.(type) {
		case int:
			if vJ, ok := valJ.(int); ok {
				return vI < vJ
			}
		case string:
			if vJ, ok := valJ.(string); ok {
				return vI < vJ
			}
		case float64:
			if vJ, ok := valJ.(float64); ok {
				return vI < vJ
			}
		}
		return false
	})

	return result
}

// FlatMapDeep creates a flattened array of values by running each element in collection through iteratee and flattening the mapped results recursively.
//
// Example:
//
//	FlatMapDeep([][]int{{1, 2}, {3, 4}}, func(x []int) [][]int { return [][]int{x, x} }) // []int{1, 2, 1, 2, 3, 4, 3, 4}
//	FlatMapDeep([]string{"hello", "world"}, func(s string) []interface{} { return []interface{}{s, []string{s}} }) // []interface{}{"hello", "hello", "world", "world"}
func FlatMapDeep[T any](slice []T, mapper func(T) interface{}) []interface{} {
	var result []interface{}
	for _, item := range slice {
		mapped := mapper(item)
		flattened := flattenDeepRecursive(mapped)
		result = append(result, flattened...)
	}
	return result
}

// flattenDeepRecursive recursively flattens any nested structure using reflection
func flattenDeepRecursive(value interface{}) []interface{} {
	var result []interface{}

	// Use reflection to handle different slice types dynamically
	v := reflect.ValueOf(value)

	// Check if it's a slice or array
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			result = append(result, flattenDeepRecursive(item)...)
		}
	} else {
		// For non-slice types, add directly
		result = append(result, value)
	}

	return result
}

// FlatMapDepth creates a flattened array of values by running each element in collection through iteratee and flattening the mapped results up to depth times.
//
// Example:
//
//	FlatMapDepth([][]int{{1, 2}, {3, 4}}, func(x []int) [][]int { return [][]int{x, x} }, 1) // [][]int{{1, 2}, {1, 2}, {3, 4}, {3, 4}}
//	FlatMapDepth([][]int{{1, 2}, {3, 4}}, func(x []int) [][]int { return [][]int{x, x} }, 2) // []int{1, 2, 1, 2, 3, 4, 3, 4}
func FlatMapDepth[T any](slice []T, mapper func(T) interface{}, depth int) []interface{} {
	var result []interface{}
	for _, item := range slice {
		mapped := mapper(item)
		flattened := flattenDepthRecursive(mapped, depth)
		result = append(result, flattened...)
	}
	return result
}

// flattenDepthRecursive recursively flattens any nested structure up to specified depth
func flattenDepthRecursive(value interface{}, depth int) []interface{} {
	if depth <= 0 {
		return []interface{}{value}
	}

	var result []interface{}
	v := reflect.ValueOf(value)

	// Check if it's a slice or array
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i).Interface()
			result = append(result, flattenDepthRecursive(item, depth-1)...)
		}
	} else {
		// For non-slice types, add directly
		result = append(result, value)
	}

	return result
}

// InvokeMap invokes the method at path of each element in collection, returning an array of the results of each invoked method.
// This is a simplified version that works with function calls rather than method paths.
//
// Example:
//
//	InvokeMap([]string{"hello", "world"}, func(s string) string { return strings.ToUpper(s) }) // []string{"HELLO", "WORLD"}
//	InvokeMap([]int{1, 2, 3}, func(x int) int { return x * x }) // []int{1, 4, 9}
func InvokeMap[T any, R any](slice []T, method func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = method(item)
	}
	return result
}

// InvokeMapWithArgs invokes the method at path of each element in collection with additional arguments.
//
// Example:
//
//	InvokeMapWithArgs([]string{"hello", "world"}, func(s string, suffix string) string { return s + suffix }, "!") // []string{"hello!", "world!"}
func InvokeMapWithArgs[T any, R any](slice []T, method func(T, ...interface{}) R, args ...interface{}) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = method(item, args...)
	}
	return result
}
