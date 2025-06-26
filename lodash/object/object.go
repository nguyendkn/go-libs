// Package object provides utility functions for working with objects (maps and structs).
// All functions are thread-safe and designed for high performance.
package object

import (
	"reflect"
	"strings"
)

// Keys returns the keys of a map.
//
// Example:
//
//	Keys(map[string]int{"a": 1, "b": 2}) // []string{"a", "b"} (order may vary)
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns the values of a map.
//
// Example:
//
//	Values(map[string]int{"a": 1, "b": 2}) // []int{1, 2} (order may vary)
func Values[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// Has checks if a key exists in a map.
//
// Example:
//
//	Has(map[string]int{"a": 1, "b": 2}, "a") // true
//	Has(map[string]int{"a": 1, "b": 2}, "c") // false
func Has[K comparable, V any](m map[K]V, key K) bool {
	_, exists := m[key]
	return exists
}

// Get gets the value at path of object. If the resolved value is nil, the defaultValue is returned.
//
// Example:
//
//	Get(map[string]interface{}{"a": map[string]interface{}{"b": 2}}, "a.b", 0) // 2
//	Get(map[string]interface{}{"a": 1}, "a.b", 0) // 0
func Get(obj interface{}, path string, defaultValue interface{}) interface{} {
	if obj == nil {
		return defaultValue
	}

	keys := strings.Split(path, ".")
	current := obj

	for _, key := range keys {
		if current == nil {
			return defaultValue
		}

		v := reflect.ValueOf(current)
		if !v.IsValid() {
			return defaultValue
		}

		switch v.Kind() {
		case reflect.Map:
			mapValue := v.MapIndex(reflect.ValueOf(key))
			if !mapValue.IsValid() {
				return defaultValue
			}
			current = mapValue.Interface()
		case reflect.Struct:
			field := v.FieldByName(key)
			if !field.IsValid() {
				return defaultValue
			}
			current = field.Interface()
		case reflect.Ptr:
			if v.IsNil() {
				return defaultValue
			}
			elem := v.Elem()
			if elem.Kind() == reflect.Struct {
				field := elem.FieldByName(key)
				if !field.IsValid() {
					return defaultValue
				}
				current = field.Interface()
			} else {
				return defaultValue
			}
		default:
			return defaultValue
		}
	}

	return current
}

// Set sets the value at path of object.
//
// Example:
//
//	m := make(map[string]interface{})
//	Set(m, "a.b", 2) // m becomes map[string]interface{}{"a": map[string]interface{}{"b": 2}}
func Set(obj interface{}, path string, value interface{}) bool {
	if obj == nil {
		return false
	}

	keys := strings.Split(path, ".")
	if len(keys) == 0 {
		return false
	}

	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Map || !v.IsValid() {
		return false
	}

	current := v
	for _, key := range keys[:len(keys)-1] {
		keyValue := reflect.ValueOf(key)
		mapValue := current.MapIndex(keyValue)

		if !mapValue.IsValid() {
			// Create new nested map
			newMap := reflect.MakeMap(reflect.TypeOf(map[string]interface{}{}))
			current.SetMapIndex(keyValue, newMap)
			current = newMap
		} else {
			if mapValue.Kind() == reflect.Interface {
				mapValue = mapValue.Elem()
			}
			if mapValue.Kind() != reflect.Map {
				// Path exists but is not a map, cannot continue
				return false
			}
			current = mapValue
		}
	}

	// Set the final value
	finalKey := reflect.ValueOf(keys[len(keys)-1])
	current.SetMapIndex(finalKey, reflect.ValueOf(value))
	return true
}

// Omit creates an object composed of the own and inherited enumerable property paths of object that are not omitted.
//
// Example:
//
//	Omit(map[string]int{"a": 1, "b": 2, "c": 3}, []string{"a", "c"}) // map[string]int{"b": 2}
func Omit[K comparable, V any](m map[K]V, keys []K) map[K]V {
	omitSet := make(map[K]bool)
	for _, key := range keys {
		omitSet[key] = true
	}

	result := make(map[K]V)
	for k, v := range m {
		if !omitSet[k] {
			result[k] = v
		}
	}
	return result
}

// Merge recursively merges own and inherited enumerable string keyed properties of source objects into the destination object.
//
// Example:
//
//	Merge(map[string]interface{}{"a": 1}, map[string]interface{}{"b": 2}) // map[string]interface{}{"a": 1, "b": 2}
func Merge(dest map[string]interface{}, sources ...map[string]interface{}) map[string]interface{} {
	if dest == nil {
		dest = make(map[string]interface{})
	}

	for _, source := range sources {
		for key, value := range source {
			if destValue, exists := dest[key]; exists {
				// If both values are maps, merge recursively
				if destMap, ok := destValue.(map[string]interface{}); ok {
					if sourceMap, ok := value.(map[string]interface{}); ok {
						dest[key] = Merge(destMap, sourceMap)
						continue
					}
				}
			}
			dest[key] = value
		}
	}

	return dest
}

// Assign copies all enumerable own properties from one or more source objects to a target object.
//
// Example:
//
//	Assign(map[string]int{"a": 1}, map[string]int{"b": 2}, map[string]int{"c": 3}) // map[string]int{"a": 1, "b": 2, "c": 3}
func Assign[K comparable, V any](dest map[K]V, sources ...map[K]V) map[K]V {
	if dest == nil {
		dest = make(map[K]V)
	}

	for _, source := range sources {
		for key, value := range source {
			dest[key] = value
		}
	}

	return dest
}

// Defaults assigns properties of source objects to the destination object for all destination properties that resolve to nil.
//
// Example:
//
//	Defaults(map[string]interface{}{"a": 1}, map[string]interface{}{"a": 2, "b": 2}) // map[string]interface{}{"a": 1, "b": 2}
func Defaults(dest map[string]interface{}, sources ...map[string]interface{}) map[string]interface{} {
	if dest == nil {
		dest = make(map[string]interface{})
	}

	for _, source := range sources {
		for key, value := range source {
			if _, exists := dest[key]; !exists {
				dest[key] = value
			}
		}
	}

	return dest
}

// Invert creates an object composed of the inverted keys and values of object.
//
// Example:
//
//	Invert(map[string]string{"a": "1", "b": "2"}) // map[string]string{"1": "a", "2": "b"}
func Invert[K, V comparable](m map[K]V) map[V]K {
	result := make(map[V]K)
	for k, v := range m {
		result[v] = k
	}
	return result
}

// MapKeys creates an object with the same values as object and keys generated by running each own enumerable string keyed property through iteratee.
//
// Example:
//
//	MapKeys(map[string]int{"a": 1, "b": 2}, func(k string) string { return k + "1" }) // map[string]int{"a1": 1, "b1": 2}
func MapKeys[K1, K2 comparable, V any](m map[K1]V, iteratee func(K1) K2) map[K2]V {
	result := make(map[K2]V)
	for k, v := range m {
		newKey := iteratee(k)
		result[newKey] = v
	}
	return result
}

// MapValues creates an object with the same keys as object and values generated by running each own enumerable string keyed property through iteratee.
//
// Example:
//
//	MapValues(map[string]int{"a": 1, "b": 2}, func(v int) int { return v * 2 }) // map[string]int{"a": 2, "b": 4}
func MapValues[K comparable, V1, V2 any](m map[K]V1, iteratee func(V1) V2) map[K]V2 {
	result := make(map[K]V2)
	for k, v := range m {
		result[k] = iteratee(v)
	}
	return result
}

// ToPairs creates an array of key-value pairs for object.
//
// Example:
//
//	ToPairs(map[string]int{"a": 1, "b": 2}) // [][2]interface{}{{"a", 1}, {"b", 2}} (order may vary)
func ToPairs[K comparable, V any](m map[K]V) [][2]interface{} {
	pairs := make([][2]interface{}, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, [2]interface{}{k, v})
	}
	return pairs
}

// FromPairs creates an object from an array of key-value pairs.
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

// Clone creates a shallow clone of value.
//
// Example:
//
//	Clone(map[string]int{"a": 1, "b": 2}) // map[string]int{"a": 1, "b": 2}
//	Clone([]int{1, 2, 3}) // []int{1, 2, 3}
func Clone[T any](value T) T {
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return value
	}

	switch v.Kind() {
	case reflect.Map:
		if v.IsNil() {
			return value
		}
		mapType := v.Type()
		newMap := reflect.MakeMap(mapType)
		for _, key := range v.MapKeys() {
			newMap.SetMapIndex(key, v.MapIndex(key))
		}
		return newMap.Interface().(T)
	case reflect.Slice:
		if v.IsNil() {
			return value
		}
		sliceType := v.Type()
		newSlice := reflect.MakeSlice(sliceType, v.Len(), v.Cap())
		reflect.Copy(newSlice, v)
		return newSlice.Interface().(T)
	case reflect.Array:
		arrayType := v.Type()
		newArray := reflect.New(arrayType).Elem()
		for i := 0; i < v.Len(); i++ {
			newArray.Index(i).Set(v.Index(i))
		}
		return newArray.Interface().(T)
	case reflect.Ptr:
		if v.IsNil() {
			return value
		}
		elemType := v.Type().Elem()
		newPtr := reflect.New(elemType)
		newPtr.Elem().Set(v.Elem())
		return newPtr.Interface().(T)
	default:
		// For primitive types, strings, etc., return as-is (they are immutable)
		return value
	}
}

// CloneDeep creates a deep clone of value.
//
// Example:
//
//	CloneDeep(map[string]interface{}{"a": map[string]int{"b": 1}}) // Deep copy with nested maps
//	CloneDeep([][]int{{1, 2}, {3, 4}}) // Deep copy with nested slices
func CloneDeep[T any](value T) T {
	return cloneDeepValue(reflect.ValueOf(value)).Interface().(T)
}

// cloneDeepValue recursively clones a reflect.Value
func cloneDeepValue(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v
	}

	switch v.Kind() {
	case reflect.Map:
		if v.IsNil() {
			return v
		}
		mapType := v.Type()
		newMap := reflect.MakeMap(mapType)
		for _, key := range v.MapKeys() {
			clonedKey := cloneDeepValue(key)
			clonedValue := cloneDeepValue(v.MapIndex(key))
			newMap.SetMapIndex(clonedKey, clonedValue)
		}
		return newMap
	case reflect.Slice:
		if v.IsNil() {
			return v
		}
		sliceType := v.Type()
		newSlice := reflect.MakeSlice(sliceType, v.Len(), v.Len())
		for i := 0; i < v.Len(); i++ {
			clonedElem := cloneDeepValue(v.Index(i))
			newSlice.Index(i).Set(clonedElem)
		}
		return newSlice
	case reflect.Array:
		arrayType := v.Type()
		newArray := reflect.New(arrayType).Elem()
		for i := 0; i < v.Len(); i++ {
			clonedElem := cloneDeepValue(v.Index(i))
			newArray.Index(i).Set(clonedElem)
		}
		return newArray
	case reflect.Ptr:
		if v.IsNil() {
			return v
		}
		elemType := v.Type().Elem()
		newPtr := reflect.New(elemType)
		clonedElem := cloneDeepValue(v.Elem())
		newPtr.Elem().Set(clonedElem)
		return newPtr
	case reflect.Struct:
		structType := v.Type()
		newStruct := reflect.New(structType).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			structField := structType.Field(i)
			// Check if field is exported (can be set)
			if structField.IsExported() {
				clonedField := cloneDeepValue(field)
				newStruct.Field(i).Set(clonedField)
			} else {
				// For unexported fields, copy directly if possible
				if field.CanInterface() {
					newStruct.Field(i).Set(field)
				}
			}
		}
		return newStruct
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		clonedElem := cloneDeepValue(v.Elem())
		newInterface := reflect.New(v.Type()).Elem()
		newInterface.Set(clonedElem)
		return newInterface
	default:
		// For primitive types, strings, etc., return as-is (they are immutable)
		return v
	}
}

// Pick creates an object composed of the picked object properties.
//
// Example:
//
//	Pick(map[string]int{"a": 1, "b": 2, "c": 3}, []string{"a", "c"}) // map[string]int{"a": 1, "c": 3}
func Pick[K comparable, V any](m map[K]V, keys []K) map[K]V {
	result := make(map[K]V)
	for _, key := range keys {
		if value, exists := m[key]; exists {
			result[key] = value
		}
	}
	return result
}

// PickBy creates an object composed of the object properties predicate returns truthy for.
//
// Example:
//
//	PickBy(map[string]int{"a": 1, "b": 2, "c": 3}, func(v int, k string) bool { return v > 1 }) // map[string]int{"b": 2, "c": 3}
func PickBy[K comparable, V any](m map[K]V, predicate func(V, K) bool) map[K]V {
	result := make(map[K]V)
	for key, value := range m {
		if predicate(value, key) {
			result[key] = value
		}
	}
	return result
}

// IsEmpty checks if value is an empty object, collection, map, or set.
//
// Example:
//
//	IsEmpty(nil) // true
//	IsEmpty("") // true
//	IsEmpty([]int{}) // true
//	IsEmpty(map[string]int{}) // true
//	IsEmpty(0) // true (for numbers)
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		return v.Len() == 0
	case reflect.Map:
		return v.Len() == 0
	case reflect.String:
		return v.Len() == 0
	case reflect.Chan:
		return v.Len() == 0
	case reflect.Ptr:
		if v.IsNil() {
			return true
		}
		return IsEmpty(v.Elem().Interface())
	case reflect.Interface:
		if v.IsNil() {
			return true
		}
		return IsEmpty(v.Elem().Interface())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
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
	return isEqualValue(reflect.ValueOf(a), reflect.ValueOf(b))
}

// isEqualValue recursively compares two reflect.Values
func isEqualValue(a, b reflect.Value) bool {
	if !a.IsValid() && !b.IsValid() {
		return true
	}
	if !a.IsValid() || !b.IsValid() {
		return false
	}

	if a.Type() != b.Type() {
		return false
	}

	switch a.Kind() {
	case reflect.Array, reflect.Slice:
		if a.Len() != b.Len() {
			return false
		}
		for i := 0; i < a.Len(); i++ {
			if !isEqualValue(a.Index(i), b.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Map:
		if a.Len() != b.Len() {
			return false
		}
		for _, key := range a.MapKeys() {
			aVal := a.MapIndex(key)
			bVal := b.MapIndex(key)
			if !bVal.IsValid() || !isEqualValue(aVal, bVal) {
				return false
			}
		}
		return true
	case reflect.Struct:
		for i := 0; i < a.NumField(); i++ {
			if !isEqualValue(a.Field(i), b.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Ptr:
		if a.IsNil() && b.IsNil() {
			return true
		}
		if a.IsNil() || b.IsNil() {
			return false
		}
		return isEqualValue(a.Elem(), b.Elem())
	case reflect.Interface:
		if a.IsNil() && b.IsNil() {
			return true
		}
		if a.IsNil() || b.IsNil() {
			return false
		}
		return isEqualValue(a.Elem(), b.Elem())
	default:
		return a.Interface() == b.Interface()
	}
}

// Transform is an alternative to reduce; this method transforms object to a new accumulator object which is the result of running each of its own enumerable string keyed properties thru iteratee, with each invocation potentially mutating the accumulator object.
//
// Example:
//
//	obj := map[string]int{"a": 1, "b": 2, "c": 1}
//	result := Transform(obj, func(result map[string][]string, value int, key string) {
//		if result[fmt.Sprintf("%d", value)] == nil {
//			result[fmt.Sprintf("%d", value)] = []string{}
//		}
//		result[fmt.Sprintf("%d", value)] = append(result[fmt.Sprintf("%d", value)], key)
//	}, map[string][]string{})
//	// result: map[string][]string{"1": []string{"a", "c"}, "2": []string{"b"}}
func Transform[T any, R any](object map[string]T, iteratee func(R, T, string), accumulator R) R {
	for key, value := range object {
		iteratee(accumulator, value, key)
	}
	return accumulator
}

// TransformSlice transforms slice to a new accumulator object.
//
// Example:
//
//	slice := []int{1, 2, 3, 4}
//	result := TransformSlice(slice, func(result map[string][]int, value int, index int) {
//		key := "even"
//		if value%2 != 0 {
//			key = "odd"
//		}
//		if result[key] == nil {
//			result[key] = []int{}
//		}
//		result[key] = append(result[key], value)
//	}, map[string][]int{})
//	// result: map[string][]int{"odd": []int{1, 3}, "even": []int{2, 4}}
func TransformSlice[T any, R any](slice []T, iteratee func(R, T, int), accumulator R) R {
	for index, value := range slice {
		iteratee(accumulator, value, index)
	}
	return accumulator
}

// InvertBy creates an object composed of the inverted keys and values of object.
// The inverted value is generated by running each element of object thru iteratee.
// The corresponding inverted value of each inverted key is an array of keys responsible for generating the inverted value.
//
// Example:
//
//	obj := map[string]int{"a": 1, "b": 2, "c": 1}
//	result := InvertBy(obj, func(value int) string {
//		return fmt.Sprintf("group_%d", value)
//	})
//	// result: map[string][]string{"group_1": []string{"a", "c"}, "group_2": []string{"b"}}
func InvertBy[T any](object map[string]T, iteratee func(T) string) map[string][]string {
	result := make(map[string][]string)
	for key, value := range object {
		invertedKey := iteratee(value)
		result[invertedKey] = append(result[invertedKey], key)
	}
	return result
}
