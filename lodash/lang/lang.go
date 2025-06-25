// Package lang provides utility functions for language constructs and type checking.
// All functions are thread-safe and designed for high performance.
package lang

import (
	"fmt"
	"reflect"
	"time"
)

// IsArray checks if value is classified as an Array object.
//
// Example:
//
//	IsArray([]int{1, 2, 3}) // true
//	IsArray("hello") // false
func IsArray(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Array || v.Kind() == reflect.Slice
}

// IsBoolean checks if value is classified as a boolean primitive.
//
// Example:
//
//	IsBoolean(true) // true
//	IsBoolean(false) // true
//	IsBoolean(1) // false
func IsBoolean(value interface{}) bool {
	if value == nil {
		return false
	}
	_, ok := value.(bool)
	return ok
}

// IsDate checks if value is classified as a Date object.
//
// Example:
//
//	IsDate(time.Now()) // true
//	IsDate("2023-01-01") // false
func IsDate(value interface{}) bool {
	if value == nil {
		return false
	}
	_, ok := value.(time.Time)
	return ok
}

// IsEmpty checks if value is an empty object, collection, map, or set.
//
// Example:
//
//	IsEmpty([]int{}) // true
//	IsEmpty(map[string]int{}) // true
//	IsEmpty("") // true
//	IsEmpty(0) // true
//	IsEmpty(false) // true
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

// IsEqual performs a deep comparison between two values to determine if they are equivalent.
//
// Example:
//
//	IsEqual([]int{1, 2}, []int{1, 2}) // true
//	IsEqual(map[string]int{"a": 1}, map[string]int{"a": 1}) // true
//	IsEqual("hello", "hello") // true
func IsEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// IsFunction checks if value is classified as a Function object.
//
// Example:
//
//	IsFunction(func() {}) // true
//	IsFunction("hello") // false
func IsFunction(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Func
}

// IsNumber checks if value is classified as a Number primitive.
//
// Example:
//
//	IsNumber(42) // true
//	IsNumber(3.14) // true
//	IsNumber("42") // false
func IsNumber(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// IsObject checks if value is the language type of Object.
//
// Example:
//
//	IsObject(map[string]int{"a": 1}) // true
//	IsObject(struct{}{}) // true
//	IsObject([]int{1, 2}) // false
func IsObject(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map, reflect.Struct:
		return true
	case reflect.Ptr:
		return !v.IsNil() && v.Elem().Kind() == reflect.Struct
	default:
		return false
	}
}

// IsString checks if value is classified as a string primitive.
//
// Example:
//
//	IsString("hello") // true
//	IsString(42) // false
func IsString(value interface{}) bool {
	if value == nil {
		return false
	}
	_, ok := value.(string)
	return ok
}

// IsNil checks if value is nil.
//
// Example:
//
//	IsNil(nil) // true
//	IsNil(0) // false
//	IsNil("") // false
func IsNil(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}

// Clone creates a shallow clone of value.
//
// Example:
//
//	original := []int{1, 2, 3}
//	cloned := Clone(original).([]int)
//	// cloned is a new slice with same elements
func Clone(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Slice:
		newSlice := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		reflect.Copy(newSlice, v)
		return newSlice.Interface()
	case reflect.Map:
		newMap := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			newMap.SetMapIndex(key, v.MapIndex(key))
		}
		return newMap.Interface()
	case reflect.Array:
		newArray := reflect.New(v.Type()).Elem()
		reflect.Copy(newArray, v)
		return newArray.Interface()
	default:
		return value
	}
}

// CloneDeep creates a deep clone of value.
//
// Example:
//
//	original := [][]int{{1, 2}, {3, 4}}
//	cloned := CloneDeep(original).([][]int)
//	// cloned is completely independent of original
func CloneDeep(value interface{}) interface{} {
	return cloneDeepValue(reflect.ValueOf(value)).Interface()
}

// cloneDeepValue recursively clones a reflect.Value
func cloneDeepValue(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}

	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return v
		}
		newPtr := reflect.New(v.Elem().Type())
		newPtr.Elem().Set(cloneDeepValue(v.Elem()))
		return newPtr
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		return cloneDeepValue(v.Elem())
	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		newSlice := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			newSlice.Index(i).Set(cloneDeepValue(v.Index(i)))
		}
		return newSlice
	case reflect.Array:
		newArray := reflect.New(v.Type()).Elem()
		for i := 0; i < v.Len(); i++ {
			newArray.Index(i).Set(cloneDeepValue(v.Index(i)))
		}
		return newArray
	case reflect.Map:
		if v.IsNil() {
			return v
		}
		newMap := reflect.MakeMap(v.Type())
		for _, key := range v.MapKeys() {
			newMap.SetMapIndex(key, cloneDeepValue(v.MapIndex(key)))
		}
		return newMap
	case reflect.Struct:
		newStruct := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).CanSet() {
				newStruct.Field(i).Set(cloneDeepValue(v.Field(i)))
			}
		}
		return newStruct
	default:
		return v
	}
}

// ToArray converts value to an array.
//
// Example:
//
//	ToArray("hello") // []interface{}{'h', 'e', 'l', 'l', 'o'}
//	ToArray([]int{1, 2, 3}) // []interface{}{1, 2, 3}
func ToArray(value interface{}) []interface{} {
	if value == nil {
		return []interface{}{}
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		str := v.String()
		result := make([]interface{}, len(str))
		for i, r := range str {
			result[i] = r
		}
		return result
	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result
	default:
		return []interface{}{value}
	}
}

// ToString converts value to a string.
//
// Example:
//
//	ToString(42) // "42"
//	ToString(true) // "true"
//	ToString(nil) // ""
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}

	if str, ok := value.(string); ok {
		return str
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return "false"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	default:
		return fmt.Sprintf("%v", value)
	}
}
