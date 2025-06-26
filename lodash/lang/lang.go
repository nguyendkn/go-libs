// Package lang provides utility functions for language constructs and type checking.
// All functions are thread-safe and designed for high performance.
package lang

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
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

// IsError checks if value is an Error object.
//
// Example:
//
//	IsError(errors.New("test")) // true
//	IsError("error string") // false
func IsError(value interface{}) bool {
	if value == nil {
		return false
	}
	_, ok := value.(error)
	return ok
}

// IsFunction checks if value is classified as a Function object.
//
// Example:
//
//	IsFunction(func() {}) // true
//	IsFunction("not a function") // false
func IsFunction(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Func
}

// IsInteger checks if value is an integer number.
//
// Example:
//
//	IsInteger(42) // true
//	IsInteger(3.14) // false
//	IsInteger("42") // false
func IsInteger(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		return f == float64(int64(f))
	default:
		return false
	}
}

// IsFloat checks if value is a floating point number.
//
// Example:
//
//	IsFloat(3.14) // true
//	IsFloat(42) // false (integer)
//	IsFloat("3.14") // false
func IsFloat(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

// IsPlainObject checks if value is a plain object (map or struct).
//
// Example:
//
//	IsPlainObject(map[string]int{"a": 1}) // true
//	IsPlainObject(struct{}{}) // true
//	IsPlainObject([]int{1, 2}) // false
func IsPlainObject(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Map:
		return true
	case reflect.Struct:
		// Exclude special types like time.Time
		if _, ok := value.(time.Time); ok {
			return false
		}
		return true
	case reflect.Ptr:
		if v.IsNil() {
			return false
		}
		elem := v.Elem()
		if elem.Kind() == reflect.Struct {
			// Exclude special types like time.Time
			if _, ok := elem.Interface().(time.Time); ok {
				return false
			}
			return true
		}
		return false
	default:
		return false
	}
}

// IsMap checks if value is a map.
//
// Example:
//
//	IsMap(map[string]int{"a": 1}) // true
//	IsMap(struct{}{}) // false
func IsMap(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Map
}

// ToNumber converts value to a number.
//
// Example:
//
//	ToNumber("42") // 42.0
//	ToNumber("3.14") // 3.14
//	ToNumber(true) // 1.0
//	ToNumber(false) // 0.0
func ToNumber(value interface{}) float64 {
	if value == nil {
		return 0
	}

	// Handle direct numeric types
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Bool:
		if v.Bool() {
			return 1.0
		}
		return 0.0
	case reflect.String:
		str := v.String()
		if str == "" {
			return 0
		}
		// Try to parse as float
		if f, err := parseFloat(str); err == nil {
			return f
		}
		return 0 // NaN equivalent
	default:
		return 0
	}
}

// ToInteger converts value to an integer.
//
// Example:
//
//	ToInteger("42") // 42
//	ToInteger(3.14) // 3
//	ToInteger(true) // 1
//	ToInteger(false) // 0
func ToInteger(value interface{}) int64 {
	return int64(ToNumber(value))
}

// parseFloat is a simple float parser without external dependencies
func parseFloat(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty string")
	}

	var result float64
	var sign float64 = 1
	var decimal bool
	var decimalPlace float64 = 0.1

	i := 0
	// Handle sign
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		if s[i] == '-' {
			sign = -1
		}
		i++
	}

	// Parse digits
	for i < len(s) {
		c := s[i]
		if c >= '0' && c <= '9' {
			digit := float64(c - '0')
			if decimal {
				result += digit * decimalPlace
				decimalPlace *= 0.1
			} else {
				result = result*10 + digit
			}
		} else if c == '.' && !decimal {
			decimal = true
		} else {
			return 0, fmt.Errorf("invalid character: %c", c)
		}
		i++
	}

	return result * sign, nil
}

// IsRegExp checks if value is a regular expression.
//
// Example:
//
//	IsRegExp(regexp.MustCompile("abc")) // true
//	IsRegExp("abc") // false
func IsRegExp(value interface{}) bool {
	if value == nil {
		return false
	}

	// Check if it's a *regexp.Regexp
	_, ok := value.(*regexp.Regexp)
	return ok
}

// IsSymbol checks if value is classified as a Symbol primitive or object.
// Note: Go doesn't have native Symbol type like JavaScript, but we can check for custom Symbol types.
//
// Example:
//
//	IsSymbol(Symbol("test")) // true (if Symbol type is defined)
//	IsSymbol("test") // false
func IsSymbol(value interface{}) bool {
	if value == nil {
		return false
	}

	// Check if the type name contains "Symbol"
	v := reflect.ValueOf(value)
	typeName := v.Type().String()

	// Check for Symbol type patterns - must be exact match or end with ".Symbol"
	// Only match types that are specifically Symbol types, not just containing "Symbol"
	// Check for Symbol type patterns
	// Match types that are exactly "Symbol" or end with "Symbol" but not containing "Symbol" in the middle
	if typeName == "Symbol" || strings.HasSuffix(typeName, ".Symbol") {
		return true
	}

	// For structs, check if type name ends with "Symbol"
	if v.Kind() == reflect.Struct {
		return strings.HasSuffix(typeName, "Symbol")
	}

	// For non-structs, be more permissive
	return strings.Contains(typeName, "Symbol")
}

// IsArrayBuffer checks if value is classified as an ArrayBuffer object.
// Note: Go doesn't have native ArrayBuffer, but we can check for byte slices.
//
// Example:
//
//	IsArrayBuffer([]byte{1, 2, 3}) // true
//	IsArrayBuffer([]int{1, 2, 3}) // false
func IsArrayBuffer(value interface{}) bool {
	if value == nil {
		return false
	}

	// Check if it's a byte slice
	_, ok := value.([]byte)
	return ok
}
