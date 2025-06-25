package json

import (
	"fmt"
	"strconv"
	"strings"
)

// GetPath extracts a value using a JSON path (e.g., "user.name", "items[0].id")
func (v *Value) GetPath(path string) (*Value, error) {
	if v == nil || v.data == nil {
		return nil, ErrNilValue
	}

	if path == "" {
		return v, nil
	}

	parts, err := parsePath(path)
	if err != nil {
		return nil, err
	}

	current := v
	for _, part := range parts {
		switch p := part.(type) {
		case string:
			current, err = current.GetByKey(p)
			if err != nil {
				return nil, fmt.Errorf("path '%s': %w", path, err)
			}
		case int:
			current, err = current.GetByIndex(p)
			if err != nil {
				return nil, fmt.Errorf("path '%s': %w", path, err)
			}
		default:
			return nil, fmt.Errorf("%w: invalid path part type", ErrInvalidPath)
		}
	}

	return current, nil
}

// SetPath sets a value using a JSON path
func (v *Value) SetPath(path string, value interface{}) error {
	if v == nil {
		return ErrNilValue
	}

	if path == "" {
		v.data = value
		return nil
	}

	parts, err := parsePath(path)
	if err != nil {
		return err
	}

	return v.setPathRecursive(parts, value)
}

// DeletePath deletes a value at the specified path
func (v *Value) DeletePath(path string) error {
	if v == nil || v.data == nil {
		return ErrNilValue
	}

	if path == "" {
		return fmt.Errorf("%w: cannot delete root", ErrInvalidPath)
	}

	parts, err := parsePath(path)
	if err != nil {
		return err
	}

	if len(parts) == 0 {
		return fmt.Errorf("%w: empty path", ErrInvalidPath)
	}

	return v.deletePathRecursive(parts, v.data)
}

// deletePathRecursive recursively deletes a path
func (v *Value) deletePathRecursive(parts []interface{}, current interface{}) error {
	if len(parts) == 1 {
		// Delete the final key/index
		lastPart := parts[0]
		switch p := lastPart.(type) {
		case string:
			obj, ok := current.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%w: parent is not an object", ErrTypeConversion)
			}
			delete(obj, p)
		case int:
			arr, ok := current.([]interface{})
			if !ok {
				return fmt.Errorf("%w: parent is not an array", ErrTypeConversion)
			}
			if p < 0 || p >= len(arr) {
				return fmt.Errorf("%w: index %d out of range", ErrIndexOutOfRange, p)
			}
			// Remove element at index
			newArr := make([]interface{}, 0, len(arr)-1)
			newArr = append(newArr, arr[:p]...)
			newArr = append(newArr, arr[p+1:]...)

			// Update the parent reference
			if current == v.data {
				v.data = newArr
			} else {
				// This case should be handled by the parent call
				return fmt.Errorf("cannot update array reference")
			}
		}
		return nil
	}

	// Navigate to next level
	part := parts[0]
	remaining := parts[1:]

	switch p := part.(type) {
	case string:
		obj, ok := current.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: value is not an object", ErrTypeConversion)
		}
		next, exists := obj[p]
		if !exists {
			return fmt.Errorf("%w: key '%s' not found", ErrKeyNotFound, p)
		}

		// Special handling for array deletion at next level
		if len(remaining) == 1 {
			if idx, ok := remaining[0].(int); ok {
				if arr, ok := next.([]interface{}); ok {
					if idx >= 0 && idx < len(arr) {
						newArr := make([]interface{}, 0, len(arr)-1)
						newArr = append(newArr, arr[:idx]...)
						newArr = append(newArr, arr[idx+1:]...)
						obj[p] = newArr
						return nil
					}
				}
			}
		}

		return v.deletePathRecursive(remaining, next)
	case int:
		arr, ok := current.([]interface{})
		if !ok {
			return fmt.Errorf("%w: value is not an array", ErrTypeConversion)
		}
		if p < 0 || p >= len(arr) {
			return fmt.Errorf("%w: index %d out of range", ErrIndexOutOfRange, p)
		}
		return v.deletePathRecursive(remaining, arr[p])
	default:
		return fmt.Errorf("%w: invalid path part type", ErrInvalidPath)
	}
}

// parsePath parses a JSON path string into parts
func parsePath(path string) ([]interface{}, error) {
	if path == "" {
		return nil, nil
	}

	var parts []interface{}
	var current strings.Builder
	var inBrackets bool

	for _, r := range path {
		switch r {
		case '.':
			if inBrackets {
				current.WriteRune(r)
			} else {
				if current.Len() > 0 {
					parts = append(parts, current.String())
					current.Reset()
				}
			}
		case '[':
			if inBrackets {
				return nil, fmt.Errorf("%w: nested brackets not allowed", ErrInvalidPath)
			}
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
			inBrackets = true
		case ']':
			if !inBrackets {
				return nil, fmt.Errorf("%w: unexpected ']'", ErrInvalidPath)
			}
			if current.Len() == 0 {
				return nil, fmt.Errorf("%w: empty brackets", ErrInvalidPath)
			}

			// Try to parse as integer
			if idx, err := strconv.Atoi(current.String()); err == nil {
				parts = append(parts, idx)
			} else {
				// Treat as string key
				parts = append(parts, current.String())
			}
			current.Reset()
			inBrackets = false
		default:
			current.WriteRune(r)
		}
	}

	if inBrackets {
		return nil, fmt.Errorf("%w: unclosed bracket", ErrInvalidPath)
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts, nil
}

// setPathRecursive recursively sets a value at the specified path
func (v *Value) setPathRecursive(parts []interface{}, value interface{}) error {
	if len(parts) == 0 {
		v.data = value
		return nil
	}

	part := parts[0]
	remaining := parts[1:]

	switch p := part.(type) {
	case string:
		// Ensure current value is an object
		if v.data == nil {
			v.data = make(map[string]interface{})
		}

		obj, ok := v.data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: cannot set key on non-object", ErrTypeConversion)
		}

		if len(remaining) == 0 {
			obj[p] = value
		} else {
			// Get or create nested value
			nested, exists := obj[p]
			if !exists {
				nested = nil
			}
			nestedValue := &Value{data: nested}
			if err := nestedValue.setPathRecursive(remaining, value); err != nil {
				return err
			}
			obj[p] = nestedValue.data
		}

	case int:
		// Ensure current value is an array
		if v.data == nil {
			v.data = make([]interface{}, 0)
		}

		arr, ok := v.data.([]interface{})
		if !ok {
			return fmt.Errorf("%w: cannot set index on non-array", ErrTypeConversion)
		}

		// Extend array if necessary
		for len(arr) <= p {
			arr = append(arr, nil)
		}
		v.data = arr

		if len(remaining) == 0 {
			arr[p] = value
		} else {
			// Get or create nested value
			nestedValue := &Value{data: arr[p]}
			if err := nestedValue.setPathRecursive(remaining, value); err != nil {
				return err
			}
			arr[p] = nestedValue.data
		}
	}

	return nil
}

// PathExists checks if a path exists in the JSON
func (v *Value) PathExists(path string) bool {
	_, err := v.GetPath(path)
	return err == nil
}
