package json

import (
	"fmt"
	"strconv"
)

// GetString extracts a string value from the JSON
func (v *Value) GetString() (string, error) {
	if v == nil || v.data == nil {
		return "", ErrNilValue
	}
	
	switch val := v.data.(type) {
	case string:
		return val, nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(val), nil
	case bool:
		return strconv.FormatBool(val), nil
	default:
		return "", fmt.Errorf("%w: cannot convert %T to string", ErrTypeConversion, val)
	}
}

// GetInt extracts an integer value from the JSON
func (v *Value) GetInt() (int, error) {
	if v == nil || v.data == nil {
		return 0, ErrNilValue
	}
	
	switch val := v.data.(type) {
	case float64:
		return int(val), nil
	case int:
		return val, nil
	case string:
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, fmt.Errorf("%w: cannot convert string '%s' to int", ErrTypeConversion, val)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("%w: cannot convert %T to int", ErrTypeConversion, val)
	}
}

// GetInt64 extracts an int64 value from the JSON
func (v *Value) GetInt64() (int64, error) {
	if v == nil || v.data == nil {
		return 0, ErrNilValue
	}
	
	switch val := v.data.(type) {
	case float64:
		return int64(val), nil
	case int:
		return int64(val), nil
	case int64:
		return val, nil
	case string:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%w: cannot convert string '%s' to int64", ErrTypeConversion, val)
		}
		return i, nil
	default:
		return 0, fmt.Errorf("%w: cannot convert %T to int64", ErrTypeConversion, val)
	}
}

// GetFloat64 extracts a float64 value from the JSON
func (v *Value) GetFloat64() (float64, error) {
	if v == nil || v.data == nil {
		return 0, ErrNilValue
	}
	
	switch val := v.data.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case string:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, fmt.Errorf("%w: cannot convert string '%s' to float64", ErrTypeConversion, val)
		}
		return f, nil
	default:
		return 0, fmt.Errorf("%w: cannot convert %T to float64", ErrTypeConversion, val)
	}
}

// GetBool extracts a boolean value from the JSON
func (v *Value) GetBool() (bool, error) {
	if v == nil || v.data == nil {
		return false, ErrNilValue
	}
	
	switch val := v.data.(type) {
	case bool:
		return val, nil
	case string:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return false, fmt.Errorf("%w: cannot convert string '%s' to bool", ErrTypeConversion, val)
		}
		return b, nil
	case float64:
		return val != 0, nil
	case int:
		return val != 0, nil
	default:
		return false, fmt.Errorf("%w: cannot convert %T to bool", ErrTypeConversion, val)
	}
}

// GetArray extracts an array value from the JSON
func (v *Value) GetArray() ([]*Value, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}
	
	arr, ok := v.data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an array", ErrTypeConversion)
	}
	
	result := make([]*Value, len(arr))
	for i, item := range arr {
		result[i] = &Value{data: item}
	}
	
	return result, nil
}

// GetObject extracts an object value from the JSON
func (v *Value) GetObject() (map[string]*Value, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}
	
	obj, ok := v.data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an object", ErrTypeConversion)
	}
	
	result := make(map[string]*Value)
	for key, val := range obj {
		result[key] = &Value{data: val}
	}
	
	return result, nil
}

// Get extracts a value by key (for objects) or index (for arrays)
func (v *Value) Get(key interface{}) (*Value, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}
	
	switch k := key.(type) {
	case string:
		return v.GetByKey(k)
	case int:
		return v.GetByIndex(k)
	default:
		return nil, fmt.Errorf("%w: key must be string or int", ErrInvalidPath)
	}
}

// GetByKey extracts a value by key from a JSON object
func (v *Value) GetByKey(key string) (*Value, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}
	
	obj, ok := v.data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an object", ErrTypeConversion)
	}
	
	val, exists := obj[key]
	if !exists {
		return nil, fmt.Errorf("%w: key '%s' not found", ErrKeyNotFound, key)
	}
	
	return &Value{data: val}, nil
}

// GetByIndex extracts a value by index from a JSON array
func (v *Value) GetByIndex(index int) (*Value, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}
	
	arr, ok := v.data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an array", ErrTypeConversion)
	}
	
	if index < 0 || index >= len(arr) {
		return nil, fmt.Errorf("%w: index %d out of range [0, %d)", ErrIndexOutOfRange, index, len(arr))
	}
	
	return &Value{data: arr[index]}, nil
}

// Has checks if a key exists in a JSON object
func (v *Value) Has(key string) bool {
	if v == nil || v.data == nil {
		return false
	}
	
	obj, ok := v.data.(map[string]interface{})
	if !ok {
		return false
	}
	
	_, exists := obj[key]
	return exists
}

// Len returns the length of an array or object
func (v *Value) Len() int {
	if v == nil || v.data == nil {
		return 0
	}
	
	switch val := v.data.(type) {
	case []interface{}:
		return len(val)
	case map[string]interface{}:
		return len(val)
	case string:
		return len(val)
	default:
		return 0
	}
}
