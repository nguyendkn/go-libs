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

// FlattenDepth recursively flattens array up to depth times.
//
// Example:
//
//	FlattenDepth([][][]int{{{1, 2}}, {{3, 4}}}, 1) // []interface{}{[]int{1, 2}, []int{3, 4}}
//	FlattenDepth([][][]int{{{1, 2}}, {{3, 4}}}, 2) // []interface{}{1, 2, 3, 4}
func FlattenDepth(slice interface{}, depth int) []interface{} {
	if depth <= 0 {
		return []interface{}{slice}
	}

	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return []interface{}{slice}
	}

	var result []interface{}
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		itemValue := reflect.ValueOf(item)

		// If item is a slice/array and we still have depth to flatten
		if (itemValue.Kind() == reflect.Slice || itemValue.Kind() == reflect.Array) && depth > 0 {
			if depth == 1 {
				// At depth 1, just add all elements of the slice
				for j := 0; j < itemValue.Len(); j++ {
					result = append(result, itemValue.Index(j).Interface())
				}
			} else {
				// Recursively flatten with reduced depth
				flattened := FlattenDepth(item, depth-1)
				result = append(result, flattened...)
			}
		} else {
			result = append(result, item)
		}
	}

	return result
}

// FromPairs returns an object composed from key-value pairs.
//
// Example:
//
//	FromPairs([][2]interface{}{{"a", 1}, {"b", 2}}) // map[interface{}]interface{}{"a": 1, "b": 2}
func FromPairs(pairs [][2]interface{}) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})

	for _, pair := range pairs {
		result[pair[0]] = pair[1]
	}

	return result
}

// FromPairsString returns a string-keyed map from key-value pairs.
//
// Example:
//
//	FromPairsString([][2]interface{}{{"a", 1}, {"b", 2}}) // map[string]interface{}{"a": 1, "b": 2}
func FromPairsString(pairs [][2]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for _, pair := range pairs {
		if len(pair) >= 2 {
			if key, ok := pair[0].(string); ok {
				result[key] = pair[1]
			}
		}
	}

	return result
}

// IntersectionBy creates an array of unique values that are included in all given arrays,
// using iteratee to generate the criterion by which uniqueness is computed.
//
// Example:
//
//	IntersectionBy(func(x int) int { return x }, []int{2, 1}, []int{2, 3}) // []int{2}
//	IntersectionBy(func(s string) int { return len(s) }, []string{"a", "bb"}, []string{"cc", "d"}) // []string{"bb"}
func IntersectionBy[T any, K comparable](iteratee func(T) K, slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	if len(slices) == 1 {
		// For single slice, return unique elements based on iteratee
		seen := make(map[K]bool)
		var result []T
		for _, item := range slices[0] {
			criterion := iteratee(item)
			if !seen[criterion] {
				seen[criterion] = true
				result = append(result, item)
			}
		}
		return result
	}

	// Create a map to track which criteria appear in each slice
	criteriaInSlices := make([]map[K]T, len(slices))

	// Initialize maps for each slice
	for i := range slices {
		criteriaInSlices[i] = make(map[K]T)
	}

	// Populate criteria for each slice
	for i, slice := range slices {
		for _, item := range slice {
			criterion := iteratee(item)
			// Store first occurrence of each criterion in this slice
			if _, exists := criteriaInSlices[i][criterion]; !exists {
				criteriaInSlices[i][criterion] = item
			}
		}
	}

	// Find criteria that exist in all slices, preserving order from first slice
	var result []T
	seen := make(map[K]bool)

	// Iterate through first slice to preserve order
	for _, item := range slices[0] {
		criterion := iteratee(item)
		if seen[criterion] {
			continue // Skip duplicates
		}
		seen[criterion] = true

		existsInAll := true
		for i := 1; i < len(slices); i++ {
			if _, exists := criteriaInSlices[i][criterion]; !exists {
				existsInAll = false
				break
			}
		}
		if existsInAll {
			result = append(result, item)
		}
	}

	return result
}

// IntersectionWith creates an array of unique values that are included in all given arrays,
// using comparator to determine equality.
//
// Example:
//
//	IntersectionWith(func(a, b int) bool { return a == b }, []int{2, 1}, []int{2, 3}) // []int{2}
//	IntersectionWith(func(a, b string) bool { return len(a) == len(b) }, []string{"a", "bb"}, []string{"cc", "d"}) // []string{"bb"}
func IntersectionWith[T any](comparator func(T, T) bool, slices ...[]T) []T {
	if len(slices) == 0 {
		return []T{}
	}

	if len(slices) == 1 {
		// For single slice, return unique elements based on comparator
		var result []T
		for _, item := range slices[0] {
			found := false
			for _, existing := range result {
				if comparator(item, existing) {
					found = true
					break
				}
			}
			if !found {
				result = append(result, item)
			}
		}
		return result
	}

	// For multiple slices, find elements from first slice that have matches in all other slices
	var result []T
	seen := make(map[int]bool) // Track indices to avoid duplicates

	for i, item := range slices[0] {
		if seen[i] {
			continue
		}

		// Check if this item (or equivalent) exists in all other slices
		existsInAll := true
		for j := 1; j < len(slices); j++ {
			found := false
			for _, otherItem := range slices[j] {
				if comparator(item, otherItem) {
					found = true
					break
				}
			}
			if !found {
				existsInAll = false
				break
			}
		}

		if existsInAll {
			// Check if we already have an equivalent item in result
			alreadyInResult := false
			for _, existing := range result {
				if comparator(item, existing) {
					alreadyInResult = true
					break
				}
			}
			if !alreadyInResult {
				result = append(result, item)
			}
		}

		seen[i] = true
	}

	return result
}

// PullAll removes all given values from array using SameValueZero for equality comparisons.
// Note: This method mutates array.
//
// Example:
//
//	arr := []int{1, 2, 3, 1, 2, 3}
//	PullAll(&arr, []int{2, 3}) // arr becomes []int{1, 1}
func PullAll[T comparable](slice *[]T, values []T) {
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

// PullAllBy removes all given values from array using iteratee to generate the criterion by which uniqueness is computed.
// Note: This method mutates array.
//
// Example:
//
//	arr := []int{1, 2, 3, 4, 5}
//	PullAllBy(&arr, []int{2, 4}, func(x int) int { return x % 2 }) // arr becomes []int{1, 3, 5} (removes even numbers)
func PullAllBy[T any, K comparable](slice *[]T, values []T, iteratee func(T) K) {
	excludeCriteria := make(map[K]bool)
	for _, value := range values {
		excludeCriteria[iteratee(value)] = true
	}

	result := (*slice)[:0] // Reuse the underlying array
	for _, item := range *slice {
		if !excludeCriteria[iteratee(item)] {
			result = append(result, item)
		}
	}
	*slice = result
}

// PullAllWith removes all given values from array using comparator to determine equality.
// Note: This method mutates array.
//
// Example:
//
//	arr := []int{1, 2, 3, 4, 5}
//	PullAllWith(&arr, []int{2, 4}, func(a, b int) bool { return a%2 == b%2 }) // removes elements with same parity
func PullAllWith[T any](slice *[]T, values []T, comparator func(T, T) bool) {
	result := (*slice)[:0] // Reuse the underlying array

	for _, item := range *slice {
		shouldExclude := false
		for _, value := range values {
			if comparator(item, value) {
				shouldExclude = true
				break
			}
		}
		if !shouldExclude {
			result = append(result, item)
		}
	}
	*slice = result
}

// PullAt removes elements from array corresponding to indexes and returns an array of removed elements.
// Note: This method mutates array.
//
// Example:
//
//	arr := []string{"a", "b", "c", "d"}
//	removed := PullAt(&arr, 1, 3) // arr becomes []string{"a", "c"}, removed is []string{"b", "d"}
func PullAt[T any](slice *[]T, indexes ...int) []T {
	if len(*slice) == 0 || len(indexes) == 0 {
		return []T{}
	}

	length := len(*slice)

	// Create a set of valid indexes to remove
	indexSet := make(map[int]bool)
	for _, idx := range indexes {
		// Handle negative indexes
		if idx < 0 {
			idx = length + idx
		}
		// Only include valid indexes
		if idx >= 0 && idx < length {
			indexSet[idx] = true
		}
	}

	// Collect removed elements and create new slice
	var removed []T
	var result []T

	for i, item := range *slice {
		if indexSet[i] {
			removed = append(removed, item)
		} else {
			result = append(result, item)
		}
	}

	*slice = result
	return removed
}

// Tail gets all but the first element of array.
//
// Example:
//
//	Tail([]int{1, 2, 3}) // []int{2, 3}
//	Tail([]int{1}) // []int{}
//	Tail([]int{}) // []int{}
func Tail[T any](slice []T) []T {
	if len(slice) <= 1 {
		return []T{}
	}
	return append([]T{}, slice[1:]...)
}

// UniqBy creates a duplicate-free version of an array, using iteratee to generate the criterion by which uniqueness is computed.
//
// Example:
//
//	UniqBy([]int{2, 1, 2}, func(x int) int { return x }) // []int{2, 1}
//	UniqBy([]string{"a", "bb", "c", "dd"}, func(s string) int { return len(s) }) // []string{"a", "bb"}
func UniqBy[T any, K comparable](slice []T, iteratee func(T) K) []T {
	seen := make(map[K]bool)
	var result []T

	for _, item := range slice {
		criterion := iteratee(item)
		if !seen[criterion] {
			seen[criterion] = true
			result = append(result, item)
		}
	}

	return result
}

// UniqWith creates a duplicate-free version of an array, using comparator to determine equality.
//
// Example:
//
//	UniqWith([]int{1, 2, 2, 3}, func(a, b int) bool { return a == b }) // []int{1, 2, 3}
//	UniqWith([]string{"a", "A", "b", "B"}, func(a, b string) bool { return strings.ToLower(a) == strings.ToLower(b) }) // []string{"a", "b"}
func UniqWith[T any](slice []T, comparator func(T, T) bool) []T {
	var result []T

	for _, item := range slice {
		found := false
		for _, existing := range result {
			if comparator(item, existing) {
				found = true
				break
			}
		}
		if !found {
			result = append(result, item)
		}
	}

	return result
}

// Helper functions for comparison
func isLess[T comparable](a, b T) bool {
	// For comparable types, we use reflection to determine ordering
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	switch va.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return va.Int() < vb.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return va.Uint() < vb.Uint()
	case reflect.Float32, reflect.Float64:
		return va.Float() < vb.Float()
	case reflect.String:
		return va.String() < vb.String()
	default:
		// For other types, convert to string and compare
		return fmt.Sprintf("%v", a) < fmt.Sprintf("%v", b)
	}
}

func isEqual[T comparable](a, b T) bool {
	return a == b
}

// SortedIndex uses a binary search to determine the lowest index at which value should be inserted into array in order to maintain its sort order.
//
// Example:
//
//	SortedIndex([]int{30, 50}, 40) // 1
//	SortedIndex([]int{4, 5, 5, 5, 6}, 5) // 1
func SortedIndex[T comparable](slice []T, value T) int {
	left, right := 0, len(slice)

	for left < right {
		mid := (left + right) / 2
		if isLess(slice[mid], value) {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return left
}

// SortedIndexBy uses a binary search to determine the lowest index at which value should be inserted into array in order to maintain its sort order, using iteratee to compute the sort ranking.
//
// Example:
//
//	SortedIndexBy([]string{"apple", "banana", "cherry"}, "blueberry", func(s string) int { return len(s) }) // 2
func SortedIndexBy[T any, K comparable](slice []T, value T, iteratee func(T) K) int {
	left, right := 0, len(slice)
	valueKey := iteratee(value)

	for left < right {
		mid := (left + right) / 2
		if isLess(iteratee(slice[mid]), valueKey) {
			left = mid + 1
		} else {
			right = mid
		}
	}

	return left
}

// SortedIndexOf performs a binary search of a sorted array to find the index of the first occurrence of value.
//
// Example:
//
//	SortedIndexOf([]int{4, 5, 5, 5, 6}, 5) // 1
//	SortedIndexOf([]int{4, 5, 5, 5, 6}, 3) // -1
func SortedIndexOf[T comparable](slice []T, value T) int {
	index := SortedIndex(slice, value)
	if index < len(slice) && isEqual(slice[index], value) {
		return index
	}
	return -1
}

// SortedLastIndex uses a binary search to determine the highest index at which value should be inserted into array in order to maintain its sort order.
//
// Example:
//
//	SortedLastIndex([]int{4, 5, 5, 5, 6}, 5) // 4
func SortedLastIndex[T comparable](slice []T, value T) int {
	left, right := 0, len(slice)

	for left < right {
		mid := (left + right) / 2
		if isLess(value, slice[mid]) {
			right = mid
		} else {
			left = mid + 1
		}
	}

	return left
}

// SortedLastIndexBy uses a binary search to determine the highest index at which value should be inserted into array in order to maintain its sort order, using iteratee to compute the sort ranking.
//
// Example:
//
//	SortedLastIndexBy([]string{"a", "bb", "ccc"}, "dd", func(s string) int { return len(s) }) // 2
func SortedLastIndexBy[T any, K comparable](slice []T, value T, iteratee func(T) K) int {
	left, right := 0, len(slice)
	valueKey := iteratee(value)

	for left < right {
		mid := (left + right) / 2
		if isLess(valueKey, iteratee(slice[mid])) {
			right = mid
		} else {
			left = mid + 1
		}
	}

	return left
}

// SortedLastIndexOf performs a binary search of a sorted array to find the index of the last occurrence of value.
//
// Example:
//
//	SortedLastIndexOf([]int{4, 5, 5, 5, 6}, 5) // 3
//	SortedLastIndexOf([]int{4, 5, 5, 5, 6}, 3) // -1
func SortedLastIndexOf[T comparable](slice []T, value T) int {
	index := SortedLastIndex(slice, value) - 1
	if index >= 0 && index < len(slice) && isEqual(slice[index], value) {
		return index
	}
	return -1
}

// SortedUniq creates a duplicate-free version of an array, using SameValueZero for equality comparisons, in which only the first occurrence of each element is kept. The order of result values is determined by the order they occur in the array.
//
// Example:
//
//	SortedUniq([]int{1, 1, 2, 2, 3}) // []int{1, 2, 3}
func SortedUniq[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return []T{}
	}

	result := []T{slice[0]}

	for i := 1; i < len(slice); i++ {
		if !isEqual(slice[i], slice[i-1]) {
			result = append(result, slice[i])
		}
	}

	return result
}

// SortedUniqBy creates a duplicate-free version of an array, using iteratee to generate the criterion by which uniqueness is computed. The order of result values is determined by the order they occur in the array.
//
// Example:
//
//	SortedUniqBy([]int{1, 1, 2, 2, 3}, func(x int) int { return x }) // []int{1, 2, 3}
//	SortedUniqBy([]string{"a", "A", "b", "B"}, func(s string) string { return strings.ToLower(s) }) // []string{"a", "b"}
func SortedUniqBy[T any, K comparable](slice []T, iteratee func(T) K) []T {
	if len(slice) == 0 {
		return []T{}
	}

	result := []T{slice[0]}
	lastKey := iteratee(slice[0])

	for i := 1; i < len(slice); i++ {
		currentKey := iteratee(slice[i])
		if currentKey != lastKey {
			result = append(result, slice[i])
			lastKey = currentKey
		}
	}

	return result
}

// ZipWith creates an array of grouped elements, each of which is the result of running each corresponding element through iteratee.
//
// Example:
//
//	ZipWith(func(a, b int) int { return a + b }, []int{1, 2}, []int{3, 4}) // []int{4, 6}
//	ZipWith(func(a, b, c int) int { return a + b + c }, []int{1, 2}, []int{3, 4}, []int{5, 6}) // []int{9, 12}
func ZipWith[T any, R any](iteratee func(...T) R, slices ...[]T) []R {
	if len(slices) == 0 {
		return []R{}
	}

	// Find the minimum length
	minLen := len(slices[0])
	for _, slice := range slices[1:] {
		if len(slice) < minLen {
			minLen = len(slice)
		}
	}

	result := make([]R, minLen)
	for i := 0; i < minLen; i++ {
		args := make([]T, len(slices))
		for j, slice := range slices {
			args[j] = slice[i]
		}
		result[i] = iteratee(args...)
	}

	return result
}

// ZipObject creates an object composed from arrays of property names and values.
//
// Example:
//
//	ZipObject([]string{"a", "b"}, []int{1, 2}) // map[string]int{"a": 1, "b": 2}
func ZipObject[K comparable, V any](keys []K, values []V) map[K]V {
	result := make(map[K]V)

	minLen := len(keys)
	if len(values) < minLen {
		minLen = len(values)
	}

	for i := 0; i < minLen; i++ {
		result[keys[i]] = values[i]
	}

	return result
}

// ZipObjectDeep creates an object composed from arrays of property paths and values.
// This is a simplified version that works with string paths.
//
// Example:
//
//	ZipObjectDeep([]string{"a.b", "c"}, []interface{}{1, 2}) // map[string]interface{}{"a": map[string]interface{}{"b": 1}, "c": 2}
func ZipObjectDeep(paths []string, values []interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	minLen := len(paths)
	if len(values) < minLen {
		minLen = len(values)
	}

	for i := 0; i < minLen; i++ {
		setDeepValue(result, paths[i], values[i])
	}

	return result
}

// Helper function to set deep values in nested maps
func setDeepValue(obj map[string]interface{}, path string, value interface{}) {
	parts := strings.Split(path, ".")
	current := obj

	for _, part := range parts[:len(parts)-1] {
		if _, exists := current[part]; !exists {
			current[part] = make(map[string]interface{})
		}
		if nested, ok := current[part].(map[string]interface{}); ok {
			current = nested
		} else {
			// If the path conflicts with existing non-object value, create new object
			current[part] = make(map[string]interface{})
			current = current[part].(map[string]interface{})
		}
	}

	current[parts[len(parts)-1]] = value
}

// TakeWhile creates a slice of array with elements taken from the beginning. Elements are taken until predicate returns falsey.
//
// Example:
//
//	TakeWhile([]int{1, 2, 3, 4, 5}, func(x int) bool { return x < 4 }) // []int{1, 2, 3}
//	TakeWhile([]string{"a", "b", "c", "1"}, func(s string) bool { return s >= "a" && s <= "z" }) // []string{"a", "b", "c"}
func TakeWhile[T any](slice []T, predicate func(T) bool) []T {
	for i, item := range slice {
		if !predicate(item) {
			return append([]T{}, slice[:i]...)
		}
	}
	return append([]T{}, slice...)
}

// TakeRightWhile creates a slice of array with elements taken from the end. Elements are taken until predicate returns falsey.
//
// Example:
//
//	TakeRightWhile([]int{1, 2, 3, 4, 5}, func(x int) bool { return x > 2 }) // []int{3, 4, 5}
//	TakeRightWhile([]string{"a", "1", "2", "3"}, func(s string) bool { return s >= "0" && s <= "9" }) // []string{"1", "2", "3"}
func TakeRightWhile[T any](slice []T, predicate func(T) bool) []T {
	for i := len(slice) - 1; i >= 0; i-- {
		if !predicate(slice[i]) {
			return append([]T{}, slice[i+1:]...)
		}
	}
	return append([]T{}, slice...)
}

// UnionBy creates an array of unique values, in order, from all given arrays, using iteratee to generate the criterion by which uniqueness is computed.
//
// Example:
//
//	UnionBy(func(x int) int { return x }, []int{2, 1}, []int{2, 3}) // []int{2, 1, 3}
//	UnionBy(func(s string) int { return len(s) }, []string{"a", "bb"}, []string{"cc", "d"}) // []string{"a", "bb", "d"}
func UnionBy[T any, K comparable](iteratee func(T) K, slices ...[]T) []T {
	seen := make(map[K]bool)
	var result []T

	for _, slice := range slices {
		for _, item := range slice {
			criterion := iteratee(item)
			if !seen[criterion] {
				seen[criterion] = true
				result = append(result, item)
			}
		}
	}

	return result
}

// UnionWith creates an array of unique values, in order, from all given arrays, using comparator to determine equality.
//
// Example:
//
//	UnionWith(func(a, b int) bool { return a == b }, []int{2, 1}, []int{2, 3}) // []int{2, 1, 3}
//	UnionWith(func(a, b string) bool { return len(a) == len(b) }, []string{"a", "bb"}, []string{"cc", "d"}) // []string{"a", "bb", "d"}
func UnionWith[T any](comparator func(T, T) bool, slices ...[]T) []T {
	var result []T

	for _, slice := range slices {
		for _, item := range slice {
			found := false
			for _, existing := range result {
				if comparator(item, existing) {
					found = true
					break
				}
			}
			if !found {
				result = append(result, item)
			}
		}
	}

	return result
}

// Unzip accepts an array of grouped elements and creates an array regrouping the elements to their pre-zip configuration.
//
// Example:
//
//	Unzip([][]int{{1, 4}, {2, 5}, {3, 6}}) // [][]int{{1, 2, 3}, {4, 5, 6}}
//	Unzip([][]string{{"a", "d"}, {"b", "e"}, {"c", "f"}}) // [][]string{{"a", "b", "c"}, {"d", "e", "f"}}
func Unzip[T any](slice [][]T) [][]T {
	if len(slice) == 0 {
		return [][]T{}
	}

	// Find the maximum length of inner slices
	maxLen := 0
	for _, inner := range slice {
		if len(inner) > maxLen {
			maxLen = len(inner)
		}
	}

	if maxLen == 0 {
		return [][]T{}
	}

	result := make([][]T, maxLen)
	for i := 0; i < maxLen; i++ {
		result[i] = make([]T, 0, len(slice))
		for _, inner := range slice {
			if i < len(inner) {
				result[i] = append(result[i], inner[i])
			}
		}
	}

	return result
}

// UnzipWith accepts an array of grouped elements and creates an array regrouping the elements to their pre-zip configuration, using iteratee to specify how regrouped values should be combined.
//
// Example:
//
//	UnzipWith([][]int{{1, 4}, {2, 5}, {3, 6}}, func(args ...int) int { return args[0] + args[1] }) // []int{5, 7, 9}
func UnzipWith[T any, R any](slice [][]T, iteratee func(...T) R) []R {
	unzipped := Unzip(slice)
	result := make([]R, len(unzipped))

	for i, group := range unzipped {
		result[i] = iteratee(group...)
	}

	return result
}

// Xor creates an array of unique values that is the symmetric difference of the given arrays.
// The order of result values is determined by the order they occur in the arrays.
//
// Example:
//
//	Xor([]int{2, 1}, []int{2, 3}) // []int{1, 3}
//	Xor([]int{1, 2}, []int{2, 3}, []int{3, 4}) // []int{1, 4}
func Xor[T comparable](slices ...[]T) []T {
	// Count occurrences of each element across all slices
	counts := make(map[T]int)
	firstOccurrence := make(map[T]T)

	for _, slice := range slices {
		seen := make(map[T]bool)
		for _, item := range slice {
			if !seen[item] {
				seen[item] = true
				counts[item]++
				if _, exists := firstOccurrence[item]; !exists {
					firstOccurrence[item] = item
				}
			}
		}
	}

	// Collect elements that appear in exactly one slice
	var result []T
	for item, count := range counts {
		if count == 1 {
			result = append(result, firstOccurrence[item])
		}
	}

	return result
}

// XorBy creates an array of unique values that is the symmetric difference of the given arrays, using iteratee to generate the criterion by which uniqueness is computed.
//
// Example:
//
//	XorBy(func(x int) int { return x }, []int{2, 1}, []int{2, 3}) // []int{1, 3}
//	XorBy(func(s string) int { return len(s) }, []string{"a", "bb"}, []string{"cc", "d"}) // []string{"bb", "d"}
func XorBy[T any, K comparable](iteratee func(T) K, slices ...[]T) []T {
	// Count occurrences of each criterion across all slices
	counts := make(map[K]int)
	firstOccurrence := make(map[K]T)

	for _, slice := range slices {
		seen := make(map[K]bool)
		for _, item := range slice {
			criterion := iteratee(item)
			if !seen[criterion] {
				seen[criterion] = true
				counts[criterion]++
				if _, exists := firstOccurrence[criterion]; !exists {
					firstOccurrence[criterion] = item
				}
			}
		}
	}

	// Collect elements whose criteria appear in exactly one slice
	var result []T
	for criterion, count := range counts {
		if count == 1 {
			result = append(result, firstOccurrence[criterion])
		}
	}

	return result
}

// XorWith creates an array of unique values that is the symmetric difference of the given arrays, using comparator to determine equality.
//
// Example:
//
//	XorWith(func(a, b int) bool { return a == b }, []int{2, 1}, []int{2, 3}) // []int{1, 3}
//	XorWith(func(a, b string) bool { return len(a) == len(b) }, []string{"a", "bb"}, []string{"cc", "d"}) // []string{"bb", "d"}
func XorWith[T any](comparator func(T, T) bool, slices ...[]T) []T {
	var allItems []T
	sliceIndices := make(map[int][]T) // Track which slice each item came from

	// Collect all items with their slice indices
	for i, slice := range slices {
		sliceIndices[i] = make([]T, 0)
		for _, item := range slice {
			// Check if this item already exists in this slice (avoid duplicates within slice)
			found := false
			for _, existing := range sliceIndices[i] {
				if comparator(item, existing) {
					found = true
					break
				}
			}
			if !found {
				sliceIndices[i] = append(sliceIndices[i], item)
				allItems = append(allItems, item)
			}
		}
	}

	var result []T
	for _, item := range allItems {
		// Count in how many slices this item appears
		sliceCount := 0
		for _, sliceItems := range sliceIndices {
			for _, sliceItem := range sliceItems {
				if comparator(item, sliceItem) {
					sliceCount++
					break
				}
			}
		}

		// If item appears in exactly one slice, include it in result
		if sliceCount == 1 {
			// Check if already in result to avoid duplicates
			found := false
			for _, existing := range result {
				if comparator(item, existing) {
					found = true
					break
				}
			}
			if !found {
				result = append(result, item)
			}
		}
	}

	return result
}
