package json

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// ConversionOptions provides options for type conversion
type ConversionOptions struct {
	StrictMode     bool
	TimeFormat     string
	NullAsZero     bool
	EmptyAsZero    bool
	TruncateFloats bool
}

// DefaultConversionOptions returns default conversion options
func DefaultConversionOptions() *ConversionOptions {
	return &ConversionOptions{
		StrictMode:     false,
		TimeFormat:     time.RFC3339,
		NullAsZero:     false,
		EmptyAsZero:    false,
		TruncateFloats: false,
	}
}

// SafeConvert safely converts a JSON value to the specified Go type
func (v *Value) SafeConvert(targetType reflect.Type, opts *ConversionOptions) (interface{}, error) {
	if opts == nil {
		opts = DefaultConversionOptions()
	}
	
	if v == nil || v.data == nil {
		if opts.NullAsZero {
			return reflect.Zero(targetType).Interface(), nil
		}
		return nil, ErrNilValue
	}
	
	return v.convertToType(targetType, opts)
}

// convertToType converts the value to the specified type
func (v *Value) convertToType(targetType reflect.Type, opts *ConversionOptions) (interface{}, error) {
	sourceValue := reflect.ValueOf(v.data)
	
	// Handle nil interface
	if !sourceValue.IsValid() {
		if opts.NullAsZero {
			return reflect.Zero(targetType).Interface(), nil
		}
		return nil, ErrNilValue
	}
	
	// Direct type match
	if sourceValue.Type() == targetType {
		return v.data, nil
	}
	
	// Handle pointer types
	if targetType.Kind() == reflect.Ptr {
		elemType := targetType.Elem()
		converted, err := v.convertToType(elemType, opts)
		if err != nil {
			return nil, err
		}
		
		result := reflect.New(elemType)
		result.Elem().Set(reflect.ValueOf(converted))
		return result.Interface(), nil
	}
	
	// Handle conversion based on target type
	switch targetType.Kind() {
	case reflect.String:
		return v.convertToString(opts)
	case reflect.Bool:
		return v.convertToBool(opts)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.convertToInt(targetType, opts)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.convertToUint(targetType, opts)
	case reflect.Float32, reflect.Float64:
		return v.convertToFloat(targetType, opts)
	case reflect.Slice:
		return v.convertToSlice(targetType, opts)
	case reflect.Map:
		return v.convertToMap(targetType, opts)
	case reflect.Struct:
		return v.convertToStruct(targetType, opts)
	case reflect.Interface:
		return v.data, nil
	default:
		return nil, fmt.Errorf("%w: unsupported target type %s", ErrTypeConversion, targetType.String())
	}
}

// convertToString converts to string
func (v *Value) convertToString(opts *ConversionOptions) (string, error) {
	switch val := v.data.(type) {
	case string:
		return val, nil
	case bool:
		return strconv.FormatBool(val), nil
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64), nil
	case int:
		return strconv.Itoa(val), nil
	case nil:
		if opts.NullAsZero {
			return "", nil
		}
		return "", ErrNilValue
	default:
		if opts.StrictMode {
			return "", fmt.Errorf("%w: cannot convert %T to string", ErrTypeConversion, val)
		}
		// Try JSON marshaling as fallback
		data, err := json.Marshal(val)
		if err != nil {
			return "", fmt.Errorf("%w: cannot convert %T to string", ErrTypeConversion, val)
		}
		return string(data), nil
	}
}

// convertToBool converts to bool
func (v *Value) convertToBool(opts *ConversionOptions) (bool, error) {
	switch val := v.data.(type) {
	case bool:
		return val, nil
	case string:
		if val == "" && opts.EmptyAsZero {
			return false, nil
		}
		b, err := strconv.ParseBool(val)
		if err != nil && !opts.StrictMode {
			// Try common string representations
			switch val {
			case "yes", "y", "1", "true", "on":
				return true, nil
			case "no", "n", "0", "false", "off", "":
				return false, nil
			default:
				return false, fmt.Errorf("%w: cannot convert string '%s' to bool", ErrTypeConversion, val)
			}
		}
		return b, err
	case float64:
		return val != 0, nil
	case int:
		return val != 0, nil
	case nil:
		if opts.NullAsZero {
			return false, nil
		}
		return false, ErrNilValue
	default:
		if opts.StrictMode {
			return false, fmt.Errorf("%w: cannot convert %T to bool", ErrTypeConversion, val)
		}
		return false, nil
	}
}

// convertToInt converts to integer types
func (v *Value) convertToInt(targetType reflect.Type, opts *ConversionOptions) (interface{}, error) {
	var intVal int64
	
	switch val := v.data.(type) {
	case float64:
		if opts.TruncateFloats {
			intVal = int64(val)
		} else {
			if val != float64(int64(val)) {
				return nil, fmt.Errorf("%w: float %f cannot be converted to int without truncation", ErrTypeConversion, val)
			}
			intVal = int64(val)
		}
	case int:
		intVal = int64(val)
	case string:
		if val == "" && opts.EmptyAsZero {
			intVal = 0
		} else {
			i, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: cannot convert string '%s' to int", ErrTypeConversion, val)
			}
			intVal = i
		}
	case bool:
		if val {
			intVal = 1
		} else {
			intVal = 0
		}
	case nil:
		if opts.NullAsZero {
			intVal = 0
		} else {
			return nil, ErrNilValue
		}
	default:
		return nil, fmt.Errorf("%w: cannot convert %T to int", ErrTypeConversion, val)
	}
	
	// Convert to specific int type
	switch targetType.Kind() {
	case reflect.Int:
		return int(intVal), nil
	case reflect.Int8:
		if intVal < -128 || intVal > 127 {
			return nil, fmt.Errorf("%w: value %d out of range for int8", ErrTypeConversion, intVal)
		}
		return int8(intVal), nil
	case reflect.Int16:
		if intVal < -32768 || intVal > 32767 {
			return nil, fmt.Errorf("%w: value %d out of range for int16", ErrTypeConversion, intVal)
		}
		return int16(intVal), nil
	case reflect.Int32:
		if intVal < -2147483648 || intVal > 2147483647 {
			return nil, fmt.Errorf("%w: value %d out of range for int32", ErrTypeConversion, intVal)
		}
		return int32(intVal), nil
	case reflect.Int64:
		return intVal, nil
	default:
		return nil, fmt.Errorf("%w: unsupported int type %s", ErrTypeConversion, targetType.String())
	}
}

// convertToUint converts to unsigned integer types
func (v *Value) convertToUint(targetType reflect.Type, opts *ConversionOptions) (interface{}, error) {
	var uintVal uint64
	
	switch val := v.data.(type) {
	case float64:
		if val < 0 {
			return nil, fmt.Errorf("%w: negative value %f cannot be converted to uint", ErrTypeConversion, val)
		}
		if opts.TruncateFloats {
			uintVal = uint64(val)
		} else {
			if val != float64(uint64(val)) {
				return nil, fmt.Errorf("%w: float %f cannot be converted to uint without truncation", ErrTypeConversion, val)
			}
			uintVal = uint64(val)
		}
	case int:
		if val < 0 {
			return nil, fmt.Errorf("%w: negative value %d cannot be converted to uint", ErrTypeConversion, val)
		}
		uintVal = uint64(val)
	case string:
		if val == "" && opts.EmptyAsZero {
			uintVal = 0
		} else {
			u, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: cannot convert string '%s' to uint", ErrTypeConversion, val)
			}
			uintVal = u
		}
	case bool:
		if val {
			uintVal = 1
		} else {
			uintVal = 0
		}
	case nil:
		if opts.NullAsZero {
			uintVal = 0
		} else {
			return nil, ErrNilValue
		}
	default:
		return nil, fmt.Errorf("%w: cannot convert %T to uint", ErrTypeConversion, val)
	}
	
	// Convert to specific uint type
	switch targetType.Kind() {
	case reflect.Uint:
		return uint(uintVal), nil
	case reflect.Uint8:
		if uintVal > 255 {
			return nil, fmt.Errorf("%w: value %d out of range for uint8", ErrTypeConversion, uintVal)
		}
		return uint8(uintVal), nil
	case reflect.Uint16:
		if uintVal > 65535 {
			return nil, fmt.Errorf("%w: value %d out of range for uint16", ErrTypeConversion, uintVal)
		}
		return uint16(uintVal), nil
	case reflect.Uint32:
		if uintVal > 4294967295 {
			return nil, fmt.Errorf("%w: value %d out of range for uint32", ErrTypeConversion, uintVal)
		}
		return uint32(uintVal), nil
	case reflect.Uint64:
		return uintVal, nil
	default:
		return nil, fmt.Errorf("%w: unsupported uint type %s", ErrTypeConversion, targetType.String())
	}
}

// convertToFloat converts to float types
func (v *Value) convertToFloat(targetType reflect.Type, opts *ConversionOptions) (interface{}, error) {
	var floatVal float64
	
	switch val := v.data.(type) {
	case float64:
		floatVal = val
	case int:
		floatVal = float64(val)
	case string:
		if val == "" && opts.EmptyAsZero {
			floatVal = 0
		} else {
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return nil, fmt.Errorf("%w: cannot convert string '%s' to float", ErrTypeConversion, val)
			}
			floatVal = f
		}
	case bool:
		if val {
			floatVal = 1.0
		} else {
			floatVal = 0.0
		}
	case nil:
		if opts.NullAsZero {
			floatVal = 0.0
		} else {
			return nil, ErrNilValue
		}
	default:
		return nil, fmt.Errorf("%w: cannot convert %T to float", ErrTypeConversion, val)
	}
	
	// Convert to specific float type
	switch targetType.Kind() {
	case reflect.Float32:
		return float32(floatVal), nil
	case reflect.Float64:
		return floatVal, nil
	default:
		return nil, fmt.Errorf("%w: unsupported float type %s", ErrTypeConversion, targetType.String())
	}
}

// convertToSlice converts to slice types
func (v *Value) convertToSlice(targetType reflect.Type, opts *ConversionOptions) (interface{}, error) {
	arr, ok := v.data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an array", ErrTypeConversion)
	}
	
	elemType := targetType.Elem()
	result := reflect.MakeSlice(targetType, len(arr), len(arr))
	
	for i, item := range arr {
		itemValue := &Value{data: item}
		converted, err := itemValue.convertToType(elemType, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to convert array element %d: %w", i, err)
		}
		result.Index(i).Set(reflect.ValueOf(converted))
	}
	
	return result.Interface(), nil
}

// convertToMap converts to map types
func (v *Value) convertToMap(targetType reflect.Type, opts *ConversionOptions) (interface{}, error) {
	obj, ok := v.data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%w: value is not an object", ErrTypeConversion)
	}
	
	keyType := targetType.Key()
	valueType := targetType.Elem()
	
	// Only support string keys for now
	if keyType.Kind() != reflect.String {
		return nil, fmt.Errorf("%w: only string keys are supported for maps", ErrTypeConversion)
	}
	
	result := reflect.MakeMap(targetType)
	
	for key, val := range obj {
		valValue := &Value{data: val}
		converted, err := valValue.convertToType(valueType, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to convert map value for key '%s': %w", key, err)
		}
		result.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(converted))
	}
	
	return result.Interface(), nil
}

// convertToStruct converts to struct types
func (v *Value) convertToStruct(targetType reflect.Type, opts *ConversionOptions) (interface{}, error) {
	// Handle time.Time specially
	if targetType == reflect.TypeOf(time.Time{}) {
		return v.convertToTime(opts)
	}
	
	// Use JSON marshaling for general struct conversion
	data, err := json.Marshal(v.data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal for struct conversion: %w", err)
	}
	
	result := reflect.New(targetType).Interface()
	if err := json.Unmarshal(data, result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to struct: %w", err)
	}
	
	return reflect.ValueOf(result).Elem().Interface(), nil
}

// convertToTime converts to time.Time
func (v *Value) convertToTime(opts *ConversionOptions) (time.Time, error) {
	switch val := v.data.(type) {
	case string:
		t, err := time.Parse(opts.TimeFormat, val)
		if err != nil {
			// Try common time formats
			formats := []string{
				time.RFC3339,
				time.RFC3339Nano,
				"2006-01-02T15:04:05Z",
				"2006-01-02 15:04:05",
				"2006-01-02",
			}
			
			for _, format := range formats {
				if t, err := time.Parse(format, val); err == nil {
					return t, nil
				}
			}
			return time.Time{}, fmt.Errorf("%w: cannot parse time string '%s'", ErrTypeConversion, val)
		}
		return t, nil
	case float64:
		// Assume Unix timestamp
		return time.Unix(int64(val), 0), nil
	case int:
		// Assume Unix timestamp
		return time.Unix(int64(val), 0), nil
	case nil:
		if opts.NullAsZero {
			return time.Time{}, nil
		}
		return time.Time{}, ErrNilValue
	default:
		return time.Time{}, fmt.Errorf("%w: cannot convert %T to time.Time", ErrTypeConversion, val)
	}
}

// ConvertTo is a convenience method for type conversion
func (v *Value) ConvertTo(target interface{}, opts *ConversionOptions) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return fmt.Errorf("%w: target must be a non-nil pointer", ErrTypeConversion)
	}
	
	targetType := targetValue.Elem().Type()
	converted, err := v.SafeConvert(targetType, opts)
	if err != nil {
		return err
	}
	
	targetValue.Elem().Set(reflect.ValueOf(converted))
	return nil
}
