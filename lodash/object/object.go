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

// IsEmpty checks if value is an empty object, collection, map, or set.
//
// Example:
//
//	IsEmpty(map[string]int{}) // true
//	IsEmpty(map[string]int{"a": 1}) // false
//	IsEmpty([]int{}) // true
//	IsEmpty([]int{1}) // false
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return true
		}
		return IsEmpty(v.Elem().Interface())
	default:
		return false
	}
}
