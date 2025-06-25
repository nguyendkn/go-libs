// Package array provides utility functions for working with arrays and slices.
// All functions are thread-safe and designed for high performance.
package array

import (
	"fmt"
	"reflect"
	"strings"
)

// Chunk creates an array of elements split into groups the length of size.
// If array can't be split evenly, the final chunk will be the remaining elements.
//
// Example:
//
//	Chunk([]int{1, 2, 3, 4, 5}, 2) // [][]int{{1, 2}, {3, 4}, {5}}
//	Chunk([]string{"a", "b", "c", "d"}, 3) // [][]string{{"a", "b", "c"}, {"d"}}
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return [][]T{}
	}

	if len(slice) == 0 {
		return [][]T{}
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// Compact creates an array with all falsey values removed.
// The values false, nil, 0, "", and empty slices/maps are falsey.
//
// Example:
//
//	Compact([]interface{}{0, 1, false, 2, "", 3}) // []interface{}{1, 2, 3}
func Compact[T any](slice []T) []T {
	var result []T

	for _, item := range slice {
		if !isFalsey(item) {
			result = append(result, item)
		}
	}

	return result
}

// isFalsey checks if a value is considered falsey
func isFalsey(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// Concat creates a new array concatenating array with any additional arrays and/or values.
//
// Example:
//
//	Concat([]int{1}, []int{2, 3}, []int{4}) // []int{1, 2, 3, 4}
func Concat[T any](slice []T, others ...[]T) []T {
	result := make([]T, len(slice))
	copy(result, slice)

	for _, other := range others {
		result = append(result, other...)
	}

	return result
}

// Difference creates an array of slice values not included in the other given arrays.
// The order and references of result values are determined by the first array.
//
// Example:
//
//	Difference([]int{2, 1}, []int{2, 3}) // []int{1}
func Difference[T comparable](slice []T, others ...[]T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	// Create a set of values to exclude
	exclude := make(map[T]bool)
	for _, other := range others {
		for _, item := range other {
			exclude[item] = true
		}
	}

	var result []T
	for _, item := range slice {
		if !exclude[item] {
			result = append(result, item)
		}
	}

	return result
}

// Drop creates a slice of array with n elements dropped from the beginning.
//
// Example:
//
//	Drop([]int{1, 2, 3}, 1) // []int{2, 3}
//	Drop([]int{1, 2, 3}, 2) // []int{3}
//	Drop([]int{1, 2, 3}, 5) // []int{}
func Drop[T any](slice []T, n int) []T {
	if n <= 0 {
		return append([]T{}, slice...)
	}

	if n >= len(slice) {
		return []T{}
	}

	return append([]T{}, slice[n:]...)
}

// DropRight creates a slice of array with n elements dropped from the end.
//
// Example:
//
//	DropRight([]int{1, 2, 3}, 1) // []int{1, 2}
//	DropRight([]int{1, 2, 3}, 2) // []int{1}
//	DropRight([]int{1, 2, 3}, 5) // []int{}
func DropRight[T any](slice []T, n int) []T {
	if n <= 0 {
		return append([]T{}, slice...)
	}

	if n >= len(slice) {
		return []T{}
	}

	return append([]T{}, slice[:len(slice)-n]...)
}

// Fill fills elements of array with value from start up to, but not including, end.
// Note: This method mutates array.
//
// Example:
//
//	arr := []int{1, 2, 3}
//	Fill(arr, 0, 1, 3) // arr becomes []int{1, 0, 0}
func Fill[T any](slice []T, value T, start, end int) {
	length := len(slice)
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start >= end {
		return
	}

	for i := start; i < end; i++ {
		slice[i] = value
	}
}

// Head gets the first element of array.
//
// Example:
//
//	Head([]int{1, 2, 3}) // 1, true
//	Head([]int{}) // 0, false
func Head[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[0], true
}

// Last gets the last element of array.
//
// Example:
//
//	Last([]int{1, 2, 3}) // 3, true
//	Last([]int{}) // 0, false
func Last[T any](slice []T) (T, bool) {
	if len(slice) == 0 {
		var zero T
		return zero, false
	}
	return slice[len(slice)-1], true
}

// Reverse reverses array so that the first element becomes the last,
// the second element becomes the second to last, and so on.
// Note: This method mutates array.
//
// Example:
//
//	arr := []int{1, 2, 3}
//	Reverse(arr) // arr becomes []int{3, 2, 1}
func Reverse[T any](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Uniq creates a duplicate-free version of an array.
//
// Example:
//
//	Uniq([]int{2, 1, 2}) // []int{2, 1}
func Uniq[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	var result []T

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}

// Flatten flattens array a single level deep.
//
// Example:
//
//	Flatten([][]int{{1, 2}, {3, 4}}) // []int{1, 2, 3, 4}
func Flatten[T any](slice [][]T) []T {
	var result []T
	for _, subSlice := range slice {
		result = append(result, subSlice...)
	}
	return result
}

// FlattenDeep recursively flattens array.
//
// Example:
//
//	FlattenDeep([][][]int{{{1, 2}}, {{3, 4}}}) // []int{1, 2, 3, 4}
func FlattenDeep(slice interface{}) []interface{} {
	var result []interface{}
	flattenRecursive(slice, &result)
	return result
}

// flattenRecursive is a helper function for FlattenDeep
func flattenRecursive(slice interface{}, result *[]interface{}) {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		*result = append(*result, slice)
		return
	}

	for i := 0; i < v.Len(); i++ {
		flattenRecursive(v.Index(i).Interface(), result)
	}
}

// IndexOf gets the index at which the first occurrence of value is found in array.
//
// Example:
//
//	IndexOf([]int{1, 2, 1, 2}, 2) // 1
//	IndexOf([]int{1, 2, 1, 2}, 3) // -1
func IndexOf[T comparable](slice []T, value T) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

// LastIndexOf gets the index at which the last occurrence of value is found in array.
//
// Example:
//
//	LastIndexOf([]int{1, 2, 1, 2}, 2) // 3
//	LastIndexOf([]int{1, 2, 1, 2}, 3) // -1
func LastIndexOf[T comparable](slice []T, value T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

// Join converts all elements in array into a string separated by separator.
//
// Example:
//
//	Join([]string{"a", "b", "c"}, "~") // "a~b~c"
//	Join([]int{1, 2, 3}, "-") // "1-2-3"
func Join[T any](slice []T, separator string) string {
	if len(slice) == 0 {
		return ""
	}

	var parts []string
	for _, item := range slice {
		parts = append(parts, fmt.Sprintf("%v", item))
	}
	return strings.Join(parts, separator)
}

// Slice creates a slice of array from start up to, but not including, end.
//
// Example:
//
//	Slice([]int{1, 2, 3, 4}, 1, 3) // []int{2, 3}
//	Slice([]int{1, 2, 3, 4}, 2, -1) // []int{3, 4}
func Slice[T any](slice []T, start, end int) []T {
	length := len(slice)

	// Handle negative indices
	if start < 0 {
		start = max(0, length+start)
	}
	if end < 0 {
		end = length + end
	}

	// Clamp to bounds
	start = max(0, min(start, length))
	end = max(start, min(end, length))

	return append([]T{}, slice[start:end]...)
}

// Take creates a slice of array with n elements taken from the beginning.
//
// Example:
//
//	Take([]int{1, 2, 3}, 2) // []int{1, 2}
//	Take([]int{1, 2, 3}, 5) // []int{1, 2, 3}
//	Take([]int{1, 2, 3}, 0) // []int{}
func Take[T any](slice []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		return append([]T{}, slice...)
	}
	return append([]T{}, slice[:n]...)
}

// TakeRight creates a slice of array with n elements taken from the end.
//
// Example:
//
//	TakeRight([]int{1, 2, 3}, 2) // []int{2, 3}
//	TakeRight([]int{1, 2, 3}, 5) // []int{1, 2, 3}
//	TakeRight([]int{1, 2, 3}, 0) // []int{}
func TakeRight[T any](slice []T, n int) []T {
	if n <= 0 {
		return []T{}
	}
	if n >= len(slice) {
		return append([]T{}, slice...)
	}
	return append([]T{}, slice[len(slice)-n:]...)
}

// Without creates an array excluding all given values.
//
// Example:
//
//	Without([]int{2, 1, 2, 3}, 1, 2) // []int{3}
func Without[T comparable](slice []T, values ...T) []T {
	exclude := make(map[T]bool)
	for _, value := range values {
		exclude[value] = true
	}

	var result []T
	for _, item := range slice {
		if !exclude[item] {
			result = append(result, item)
		}
	}
	return result
}

// Zip creates an array of grouped elements, the first of which contains
// the first elements of the given arrays.
//
// Example:
//
//	Zip([]string{"a", "b"}, []int{1, 2}) // [][2]interface{}{{"a", 1}, {"b", 2}}
func Zip[T, U any](slice1 []T, slice2 []U) [][2]interface{} {
	minLen := min(len(slice1), len(slice2))
	result := make([][2]interface{}, minLen)

	for i := 0; i < minLen; i++ {
		result[i] = [2]interface{}{slice1[i], slice2[i]}
	}

	return result
}

// Initial gets all but the last element of array.
//
// Example:
//
//	Initial([]int{1, 2, 3}) // []int{1, 2}
//	Initial([]int{}) // []int{}
func Initial[T any](slice []T) []T {
	if len(slice) <= 1 {
		return []T{}
	}
	return append([]T{}, slice[:len(slice)-1]...)
}

// Intersection creates an array of unique values that are included in all given arrays.
//
// Example:
//
//	Intersection([]int{2, 1}, []int{2, 3}) // []int{2}
func Intersection[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	// Count occurrences in each slice
	counts := make(map[T]int)
	seen := make(map[T]map[int]bool)

	for i, slice := range slices {
		sliceSeen := make(map[T]bool)
		for _, item := range slice {
			if !sliceSeen[item] {
				sliceSeen[item] = true
				if seen[item] == nil {
					seen[item] = make(map[int]bool)
				}
				if !seen[item][i] {
					seen[item][i] = true
					counts[item]++
				}
			}
		}
	}

	// Find items that appear in all slices
	var result []T
	for item, count := range counts {
		if count == len(slices) {
			result = append(result, item)
		}
	}

	return result
}

// Union creates an array of unique values, in order, from all given arrays.
//
// Example:
//
//	Union([]int{2}, []int{1, 2}) // []int{2, 1}
func Union[T comparable](slices ...[]T) []T {
	seen := make(map[T]bool)
	var result []T

	for _, slice := range slices {
		for _, item := range slice {
			if !seen[item] {
				seen[item] = true
				result = append(result, item)
			}
		}
	}

	return result
}

// Nth gets the element at index n of array. If n is negative, the nth element from the end is returned.
//
// Example:
//
//	Nth([]string{"a", "b", "c", "d"}, 1) // "b", true
//	Nth([]string{"a", "b", "c", "d"}, -2) // "c", true
//	Nth([]string{"a", "b", "c", "d"}, 10) // "", false
func Nth[T any](slice []T, n int) (T, bool) {
	length := len(slice)
	var zero T

	if length == 0 {
		return zero, false
	}

	// Handle negative index
	if n < 0 {
		n = length + n
	}

	if n < 0 || n >= length {
		return zero, false
	}

	return slice[n], true
}

// Pull removes all given values from array using SameValueZero for equality comparisons.
// Note: This method mutates array.
//
// Example:
//
//	arr := []int{1, 2, 3, 1, 2, 3}
//	Pull(arr, 2, 3) // arr becomes []int{1, 1}
func Pull[T comparable](slice *[]T, values ...T) {
	exclude := make(map[T]bool)
	for _, value := range values {
		exclude[value] = true
	}

	result := (*slice)[:0] // Reuse the underlying array
	for _, item := range *slice {
		if !exclude[item] {
			result = append(result, item)
		}
	}
	*slice = result
}

// Remove removes all elements from array that predicate returns truthy for.
// Note: This method mutates array.
//
// Example:
//
//	arr := []int{1, 2, 3, 4}
//	removed := Remove(&arr, func(x int) bool { return x%2 == 0 })
//	// arr becomes []int{1, 3}, removed is []int{2, 4}
func Remove[T any](slice *[]T, predicate func(T) bool) []T {
	var removed []T
	result := (*slice)[:0] // Reuse the underlying array

	for _, item := range *slice {
		if predicate(item) {
			removed = append(removed, item)
		} else {
			result = append(result, item)
		}
	}

	*slice = result
	return removed
}
