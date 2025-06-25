// Package json provides a comprehensive JSON library for Go with zero external dependencies.
// It offers fast parsing, flexible serialization, validation, pretty printing, and query capabilities.
package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Common errors
var (
	ErrInvalidJSON     = errors.New("invalid JSON format")
	ErrInvalidPath     = errors.New("invalid JSON path")
	ErrTypeConversion  = errors.New("type conversion error")
	ErrNilValue        = errors.New("nil value")
	ErrIndexOutOfRange = errors.New("index out of range")
	ErrKeyNotFound     = errors.New("key not found")
)

// Value represents a JSON value that can be of any type
type Value struct {
	data interface{}
}

// New creates a new JSON Value from any Go value
func New(v interface{}) *Value {
	// Convert Go types to JSON-compatible types
	data, err := json.Marshal(v)
	if err != nil {
		return &Value{data: v}
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return &Value{data: v}
	}

	return &Value{data: jsonData}
}

// Parse parses JSON from a string
func Parse(s string) (*Value, error) {
	return ParseBytes([]byte(s))
}

// ParseBytes parses JSON from a byte slice
func ParseBytes(data []byte) (*Value, error) {
	if len(data) == 0 {
		return nil, ErrInvalidJSON
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return &Value{data: v}, nil
}

// ParseReader parses JSON from an io.Reader
func ParseReader(r io.Reader) (*Value, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	return ParseBytes(data)
}

// IsValid checks if a string contains valid JSON
func IsValid(s string) bool {
	return json.Valid([]byte(s))
}

// IsValidBytes checks if a byte slice contains valid JSON
func IsValidBytes(data []byte) bool {
	return json.Valid(data)
}

// String returns the JSON string representation
func (v *Value) String() string {
	if v == nil || v.data == nil {
		return "null"
	}

	data, err := json.Marshal(v.data)
	if err != nil {
		return "null"
	}

	return string(data)
}

// PrettyString returns a pretty-printed JSON string with default indentation
func (v *Value) PrettyString() string {
	return v.PrettyStringIndent("  ")
}

// PrettyStringIndent returns a pretty-printed JSON string with custom indentation
func (v *Value) PrettyStringIndent(indent string) string {
	if v == nil || v.data == nil {
		return "null"
	}

	data, err := json.MarshalIndent(v.data, "", indent)
	if err != nil {
		return "null"
	}

	return string(data)
}

// Bytes returns the JSON byte representation
func (v *Value) Bytes() []byte {
	if v == nil || v.data == nil {
		return []byte("null")
	}

	data, err := json.Marshal(v.data)
	if err != nil {
		return []byte("null")
	}

	return data
}

// Interface returns the underlying interface{} value
func (v *Value) Interface() interface{} {
	if v == nil {
		return nil
	}
	return v.data
}

// Type returns the reflect.Type of the underlying value
func (v *Value) Type() reflect.Type {
	if v == nil || v.data == nil {
		return nil
	}
	return reflect.TypeOf(v.data)
}

// IsNull checks if the value is null
func (v *Value) IsNull() bool {
	return v == nil || v.data == nil
}

// IsObject checks if the value is a JSON object
func (v *Value) IsObject() bool {
	if v == nil || v.data == nil {
		return false
	}
	_, ok := v.data.(map[string]interface{})
	return ok
}

// IsArray checks if the value is a JSON array
func (v *Value) IsArray() bool {
	if v == nil || v.data == nil {
		return false
	}
	_, ok := v.data.([]interface{})
	return ok
}

// IsString checks if the value is a string
func (v *Value) IsString() bool {
	if v == nil || v.data == nil {
		return false
	}
	_, ok := v.data.(string)
	return ok
}

// IsNumber checks if the value is a number
func (v *Value) IsNumber() bool {
	if v == nil || v.data == nil {
		return false
	}
	switch v.data.(type) {
	case float64, int, int64, float32:
		return true
	default:
		return false
	}
}

// IsBool checks if the value is a boolean
func (v *Value) IsBool() bool {
	if v == nil || v.data == nil {
		return false
	}
	_, ok := v.data.(bool)
	return ok
}
