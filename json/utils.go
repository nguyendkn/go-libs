package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// Marshal converts a Go value to JSON bytes
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// MarshalIndent converts a Go value to indented JSON bytes
func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}

// Unmarshal parses JSON bytes into a Go value
func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// UnmarshalToValue parses JSON bytes into a Value
func UnmarshalToValue(data []byte) (*Value, error) {
	return ParseBytes(data)
}

// Clone creates a deep copy of the JSON value
func (v *Value) Clone() *Value {
	if v == nil || v.data == nil {
		return &Value{data: nil}
	}

	// Use JSON marshal/unmarshal for deep copy
	data, err := json.Marshal(v.data)
	if err != nil {
		return &Value{data: nil}
	}

	var cloned interface{}
	if err := json.Unmarshal(data, &cloned); err != nil {
		return &Value{data: nil}
	}

	return &Value{data: cloned}
}

// Equal compares two JSON values for equality
func (v *Value) Equal(other *Value) bool {
	if v == nil && other == nil {
		return true
	}
	if v == nil || other == nil {
		return false
	}

	return reflect.DeepEqual(v.data, other.data)
}

// Merge merges another JSON object into this one
func (v *Value) Merge(other *Value) error {
	if v == nil {
		return ErrNilValue
	}
	if other == nil || other.data == nil {
		return nil
	}

	// Both must be objects
	obj1, ok1 := v.data.(map[string]interface{})
	if !ok1 {
		return fmt.Errorf("%w: target is not an object", ErrTypeConversion)
	}

	obj2, ok2 := other.data.(map[string]interface{})
	if !ok2 {
		return fmt.Errorf("%w: source is not an object", ErrTypeConversion)
	}

	// Merge recursively
	for key, val := range obj2 {
		if existing, exists := obj1[key]; exists {
			// If both are objects, merge recursively
			if existingObj, ok := existing.(map[string]interface{}); ok {
				if valObj, ok := val.(map[string]interface{}); ok {
					existingValue := &Value{data: existingObj}
					valValue := &Value{data: valObj}
					if err := existingValue.Merge(valValue); err != nil {
						return err
					}
					continue
				}
			}
		}
		// Otherwise, overwrite
		obj1[key] = val
	}

	return nil
}

// Keys returns all keys of a JSON object
func (v *Value) Keys() ([]string, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}

	obj, ok := v.data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an object", ErrTypeConversion)
	}

	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}

	return keys, nil
}

// Values returns all values of a JSON object or array
func (v *Value) Values() ([]*Value, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}

	switch val := v.data.(type) {
	case map[string]interface{}:
		values := make([]*Value, 0, len(val))
		for _, v := range val {
			values = append(values, &Value{data: v})
		}
		return values, nil
	case []interface{}:
		values := make([]*Value, len(val))
		for i, v := range val {
			values[i] = &Value{data: v}
		}
		return values, nil
	default:
		return nil, fmt.Errorf("%w: value is not an object or array", ErrTypeConversion)
	}
}

// ToMap converts the JSON value to a map[string]interface{}
func (v *Value) ToMap() (map[string]interface{}, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}

	obj, ok := v.data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an object", ErrTypeConversion)
	}

	// Create a copy
	result := make(map[string]interface{})
	for k, v := range obj {
		result[k] = v
	}

	return result, nil
}

// ToSlice converts the JSON value to a []interface{}
func (v *Value) ToSlice() ([]interface{}, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}

	arr, ok := v.data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an array", ErrTypeConversion)
	}

	// Create a copy
	result := make([]interface{}, len(arr))
	copy(result, arr)

	return result, nil
}

// UnmarshalTo unmarshals the JSON value into a Go struct
func (v *Value) UnmarshalTo(target interface{}) error {
	if v == nil || v.data == nil {
		return ErrNilValue
	}

	data, err := json.Marshal(v.data)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return json.Unmarshal(data, target)
}

// FromStruct creates a JSON value from a Go struct
func FromStruct(v interface{}) (*Value, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %w", err)
	}

	return ParseBytes(data)
}

// Compact removes all whitespace from JSON
func Compact(data []byte) []byte {
	var buf bytes.Buffer
	if err := json.Compact(&buf, data); err != nil {
		return data
	}
	return buf.Bytes()
}

// CompactString removes all whitespace from JSON string
func CompactString(s string) string {
	return string(Compact([]byte(s)))
}

// Indent adds indentation to JSON
func Indent(data []byte, prefix, indent string) []byte {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, prefix, indent); err != nil {
		return data
	}
	return buf.Bytes()
}

// IndentString adds indentation to JSON string
func IndentString(s, prefix, indent string) string {
	return string(Indent([]byte(s), prefix, indent))
}

// Size returns the size of the JSON value in bytes
func (v *Value) Size() int {
	if v == nil || v.data == nil {
		return 4 // "null"
	}

	data, err := json.Marshal(v.data)
	if err != nil {
		return 4 // "null"
	}

	return len(data)
}
