// Package util provides general utility functions.
// All functions are thread-safe and designed for high performance.
package util

import (
	"fmt"
	"math/rand"
	"reflect"
	"sync/atomic"
	"time"
)

var uniqueIdCounter int64

// Identity returns the first argument it receives.
//
// Example:
//
//	Identity(42) // 42
//	Identity("hello") // "hello"
func Identity[T any](value T) T {
	return value
}

// Constant returns a function that always returns the same value.
//
// Example:
//
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
//
//	Noop() // does nothing
func Noop() {
	// Intentionally empty
}

// Range creates an array of numbers progressing from start up to, but not including, end.
//
// Example:
//
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
//
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
//
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
//
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
//
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
//
//	result, err := Attempt(func() (int, error) { return 42, nil })
//	// result: 42, err: nil
func Attempt[T any](fn func() (T, error)) (T, error) {
	return fn()
}

// Flow creates a function that is the composition of the provided functions,
// where each successive invocation is supplied the return value of the previous.
//
// Example:
//
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
//
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
//
//	StubArray() // []interface{}{}
func StubArray() []interface{} {
	return []interface{}{}
}

// StubFalse returns false.
//
// Example:
//
//	StubFalse() // false
func StubFalse() bool {
	return false
}

// StubObject returns a new empty object.
//
// Example:
//
//	StubObject() // map[string]interface{}{}
func StubObject() map[string]interface{} {
	return map[string]interface{}{}
}

// StubString returns an empty string.
//
// Example:
//
//	StubString() // ""
func StubString() string {
	return ""
}

// StubTrue returns true.
//
// Example:
//
//	StubTrue() // true
func StubTrue() bool {
	return true
}

// Random generates a random number between min and max (inclusive).
//
// Example:
//
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
//
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
//
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
//
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

// Sample gets a random element from collection.
//
// Example:
//
//	Sample([]int{1, 2, 3, 4, 5}) // random element like 3
//	Sample([]string{"a", "b", "c"}) // random element like "b"
func Sample[T any](collection []T) T {
	var zero T
	if len(collection) == 0 {
		return zero
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(collection))
	return collection[index]
}

// SampleSize gets n random elements at unique keys from collection up to the size of collection.
//
// Example:
//
//	SampleSize([]int{1, 2, 3, 4, 5}, 3) // random 3 elements like [2, 4, 1]
//	SampleSize([]string{"a", "b", "c"}, 2) // random 2 elements like ["c", "a"]
func SampleSize[T any](collection []T, n int) []T {
	if len(collection) == 0 || n <= 0 {
		return []T{}
	}

	if n >= len(collection) {
		return Shuffle(collection)
	}

	rand.Seed(time.Now().UnixNano())

	// Create a copy to avoid modifying original
	temp := make([]T, len(collection))
	copy(temp, collection)

	// Fisher-Yates shuffle for first n elements
	for i := 0; i < n; i++ {
		j := rand.Intn(len(temp)-i) + i
		temp[i], temp[j] = temp[j], temp[i]
	}

	return temp[:n]
}

// Shuffle creates a shuffled array of values, using a version of the Fisher-Yates shuffle.
//
// Example:
//
//	Shuffle([]int{1, 2, 3, 4}) // random order like [3, 1, 4, 2]
//	Shuffle([]string{"a", "b", "c"}) // random order like ["c", "a", "b"]
func Shuffle[T any](collection []T) []T {
	if len(collection) <= 1 {
		result := make([]T, len(collection))
		copy(result, collection)
		return result
	}

	rand.Seed(time.Now().UnixNano())

	// Create a copy to avoid modifying original
	result := make([]T, len(collection))
	copy(result, collection)

	// Fisher-Yates shuffle
	for i := len(result) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// Size gets the size of collection by returning its length for array-like values or the number of own enumerable string keyed properties for objects.
//
// Example:
//
//	Size([]int{1, 2, 3}) // 3
//	Size(map[string]int{"a": 1, "b": 2}) // 2
//	Size("hello") // 5
func Size(collection interface{}) int {
	if collection == nil {
		return 0
	}

	switch v := collection.(type) {
	case string:
		return len(v)
	case []interface{}:
		return len(v)
	case []int:
		return len(v)
	case []string:
		return len(v)
	case []float64:
		return len(v)
	case []bool:
		return len(v)
	case map[string]interface{}:
		return len(v)
	case map[string]int:
		return len(v)
	case map[string]string:
		return len(v)
	case map[int]interface{}:
		return len(v)
	default:
		// Use reflection for other types
		rv := reflect.ValueOf(collection)
		switch rv.Kind() {
		case reflect.Array, reflect.Slice, reflect.Map, reflect.String, reflect.Chan:
			return rv.Len()
		default:
			return 0
		}
	}
}

// Property creates a function that returns the value at path of a given object.
//
// Example:
//
//	objects := []map[string]interface{}{{"a": map[string]interface{}{"b": 2}}, {"a": map[string]interface{}{"b": 1}}}
//	getB := Property("a.b")
//	values := Map(objects, getB) // [2, 1]
func Property(path string) func(interface{}) interface{} {
	pathArray := ToPath(path)
	return func(obj interface{}) interface{} {
		return getValueAtPath(obj, pathArray)
	}
}

// PropertyOf creates a function that returns the value at path of object.
//
// Example:
//
//	object := map[string]interface{}{"a": map[string]interface{}{"b": 2}}
//	getValue := PropertyOf(object)
//	result := getValue("a.b") // 2
func PropertyOf(object interface{}) func(string) interface{} {
	return func(path string) interface{} {
		pathArray := ToPath(path)
		return getValueAtPath(object, pathArray)
	}
}

// getValueAtPath gets the value at the specified path in an object
func getValueAtPath(obj interface{}, path []string) interface{} {
	if obj == nil || len(path) == 0 {
		return nil
	}

	current := obj
	for _, key := range path {
		switch v := current.(type) {
		case map[string]interface{}:
			if val, exists := v[key]; exists {
				current = val
			} else {
				return nil
			}
		case map[string]string:
			if val, exists := v[key]; exists {
				current = val
			} else {
				return nil
			}
		case map[string]int:
			if val, exists := v[key]; exists {
				current = val
			} else {
				return nil
			}
		default:
			// Use reflection for other map types or struct fields
			rv := reflect.ValueOf(current)
			if rv.Kind() == reflect.Ptr {
				if rv.IsNil() {
					return nil
				}
				rv = rv.Elem()
			}

			switch rv.Kind() {
			case reflect.Map:
				keyValue := reflect.ValueOf(key)
				mapValue := rv.MapIndex(keyValue)
				if !mapValue.IsValid() {
					return nil
				}
				current = mapValue.Interface()
			case reflect.Struct:
				fieldValue := rv.FieldByName(key)
				if !fieldValue.IsValid() || !fieldValue.CanInterface() {
					return nil
				}
				current = fieldValue.Interface()
			default:
				return nil
			}
		}
	}

	return current
}

// Matches creates a predicate function that tells you if a given object has equivalent property values.
//
// Example:
//
//	objects := []map[string]interface{}{{"a": 1, "b": 2}, {"a": 1, "b": 3}, {"a": 2, "b": 2}}
//	predicate := Matches(map[string]interface{}{"a": 1})
//	filtered := Filter(objects, predicate) // [{"a": 1, "b": 2}, {"a": 1, "b": 3}]
func Matches(source interface{}) func(interface{}) bool {
	return func(object interface{}) bool {
		return isMatch(object, source)
	}
}

// isMatch checks if object matches source properties
func isMatch(object, source interface{}) bool {
	if source == nil {
		return true
	}
	if object == nil {
		return false
	}

	sourceValue := reflect.ValueOf(source)
	objectValue := reflect.ValueOf(object)

	// Handle maps
	if sourceValue.Kind() == reflect.Map && objectValue.Kind() == reflect.Map {
		for _, key := range sourceValue.MapKeys() {
			sourceVal := sourceValue.MapIndex(key)
			objectVal := objectValue.MapIndex(key)

			if !objectVal.IsValid() {
				return false
			}

			if !reflect.DeepEqual(sourceVal.Interface(), objectVal.Interface()) {
				return false
			}
		}
		return true
	}

	// Handle structs
	if sourceValue.Kind() == reflect.Struct && objectValue.Kind() == reflect.Struct {
		sourceType := sourceValue.Type()
		for i := 0; i < sourceValue.NumField(); i++ {
			field := sourceType.Field(i)
			if !field.IsExported() {
				continue
			}

			sourceFieldVal := sourceValue.Field(i)
			objectFieldVal := objectValue.FieldByName(field.Name)

			if !objectFieldVal.IsValid() {
				return false
			}

			if !reflect.DeepEqual(sourceFieldVal.Interface(), objectFieldVal.Interface()) {
				return false
			}
		}
		return true
	}

	// For other types, use direct comparison
	return reflect.DeepEqual(object, source)
}
