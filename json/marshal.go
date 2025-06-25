package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Marshaler interface for custom JSON marshaling
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

// Unmarshaler interface for custom JSON unmarshaling
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

// MarshalOptions provides options for JSON marshaling
type MarshalOptions struct {
	Indent          string
	Prefix          string
	EscapeHTML      bool
	SortKeys        bool
	DisallowUnknown bool
}

// DefaultMarshalOptions returns default marshaling options
func DefaultMarshalOptions() *MarshalOptions {
	return &MarshalOptions{
		Indent:          "",
		Prefix:          "",
		EscapeHTML:      true,
		SortKeys:        false,
		DisallowUnknown: false,
	}
}

// MarshalWithOptions marshals a value with custom options
func MarshalWithOptions(v interface{}, opts *MarshalOptions) ([]byte, error) {
	if opts == nil {
		opts = DefaultMarshalOptions()
	}
	
	// Handle custom marshaler
	if marshaler, ok := v.(Marshaler); ok {
		return marshaler.MarshalJSON()
	}
	
	// Handle time.Time specially
	if t, ok := v.(time.Time); ok {
		return json.Marshal(t.Format(time.RFC3339))
	}
	
	// Handle *time.Time specially
	if t, ok := v.(*time.Time); ok && t != nil {
		return json.Marshal(t.Format(time.RFC3339))
	}
	
	if opts.Indent != "" {
		return json.MarshalIndent(v, opts.Prefix, opts.Indent)
	}
	
	return json.Marshal(v)
}

// UnmarshalWithOptions unmarshals JSON with custom options
func UnmarshalWithOptions(data []byte, v interface{}, opts *MarshalOptions) error {
	if opts == nil {
		opts = DefaultMarshalOptions()
	}
	
	// Handle custom unmarshaler
	if unmarshaler, ok := v.(Unmarshaler); ok {
		return unmarshaler.UnmarshalJSON(data)
	}
	
	decoder := json.NewDecoder(nil)
	if opts.DisallowUnknown {
		decoder.DisallowUnknownFields()
	}
	
	return json.Unmarshal(data, v)
}

// MarshalStruct marshals a struct to JSON with field validation
func MarshalStruct(v interface{}) (*Value, error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return &Value{data: nil}, nil
		}
		rv = rv.Elem()
	}
	
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("%w: value is not a struct", ErrTypeConversion)
	}
	
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %w", err)
	}
	
	return ParseBytes(data)
}

// UnmarshalStruct unmarshals JSON to a struct with validation
func UnmarshalStruct(data []byte, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("%w: target must be a non-nil pointer", ErrTypeConversion)
	}
	
	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("%w: target must be a pointer to struct", ErrTypeConversion)
	}
	
	return json.Unmarshal(data, v)
}

// MarshalMap marshals a map to JSON
func MarshalMap(m map[string]interface{}) (*Value, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal map: %w", err)
	}
	
	return ParseBytes(data)
}

// UnmarshalMap unmarshals JSON to a map
func UnmarshalMap(data []byte) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
	}
	
	return m, nil
}

// MarshalSlice marshals a slice to JSON
func MarshalSlice(s []interface{}) (*Value, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal slice: %w", err)
	}
	
	return ParseBytes(data)
}

// UnmarshalSlice unmarshals JSON to a slice
func UnmarshalSlice(data []byte) ([]interface{}, error) {
	var s []interface{}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to slice: %w", err)
	}
	
	return s, nil
}

// ToJSON converts any Go value to a JSON Value
func ToJSON(v interface{}) (*Value, error) {
	switch val := v.(type) {
	case *Value:
		return val, nil
	case nil:
		return &Value{data: nil}, nil
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return &Value{data: val}, nil
	case []interface{}:
		return &Value{data: val}, nil
	case map[string]interface{}:
		return &Value{data: val}, nil
	default:
		// Use JSON marshal/unmarshal for complex types
		data, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to JSON: %w", err)
		}
		
		return ParseBytes(data)
	}
}

// FromJSON converts a JSON Value to a Go value of the specified type
func FromJSON(v *Value, target interface{}) error {
	if v == nil || v.data == nil {
		return ErrNilValue
	}
	
	return v.UnmarshalTo(target)
}

// Set sets a value in the JSON (for objects and arrays)
func (v *Value) Set(key interface{}, value interface{}) error {
	if v == nil {
		return ErrNilValue
	}
	
	switch k := key.(type) {
	case string:
		return v.SetKey(k, value)
	case int:
		return v.SetIndex(k, value)
	default:
		return fmt.Errorf("%w: key must be string or int", ErrInvalidPath)
	}
}

// SetKey sets a value by key in a JSON object
func (v *Value) SetKey(key string, value interface{}) error {
	if v == nil {
		return ErrNilValue
	}
	
	// Initialize as object if nil
	if v.data == nil {
		v.data = make(map[string]interface{})
	}
	
	obj, ok := v.data.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%w: value is not an object", ErrTypeConversion)
	}
	
	obj[key] = value
	return nil
}

// SetIndex sets a value by index in a JSON array
func (v *Value) SetIndex(index int, value interface{}) error {
	if v == nil {
		return ErrNilValue
	}
	
	// Initialize as array if nil
	if v.data == nil {
		v.data = make([]interface{}, 0)
	}
	
	arr, ok := v.data.([]interface{})
	if !ok {
		return fmt.Errorf("%w: value is not an array", ErrTypeConversion)
	}
	
	// Extend array if necessary
	for len(arr) <= index {
		arr = append(arr, nil)
	}
	
	arr[index] = value
	v.data = arr
	return nil
}

// Append appends a value to a JSON array
func (v *Value) Append(value interface{}) error {
	if v == nil {
		return ErrNilValue
	}
	
	// Initialize as array if nil
	if v.data == nil {
		v.data = make([]interface{}, 0)
	}
	
	arr, ok := v.data.([]interface{})
	if !ok {
		return fmt.Errorf("%w: value is not an array", ErrTypeConversion)
	}
	
	arr = append(arr, value)
	v.data = arr
	return nil
}

// Remove removes a key from a JSON object or index from array
func (v *Value) Remove(key interface{}) error {
	if v == nil || v.data == nil {
		return ErrNilValue
	}
	
	switch k := key.(type) {
	case string:
		obj, ok := v.data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: value is not an object", ErrTypeConversion)
		}
		delete(obj, k)
		return nil
		
	case int:
		arr, ok := v.data.([]interface{})
		if !ok {
			return fmt.Errorf("%w: value is not an array", ErrTypeConversion)
		}
		if k < 0 || k >= len(arr) {
			return fmt.Errorf("%w: index %d out of range", ErrIndexOutOfRange, k)
		}
		// Remove element at index
		copy(arr[k:], arr[k+1:])
		arr = arr[:len(arr)-1]
		v.data = arr
		return nil
		
	default:
		return fmt.Errorf("%w: key must be string or int", ErrInvalidPath)
	}
}
