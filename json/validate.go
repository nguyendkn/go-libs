package json

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

// ValidationError represents a JSON validation error
type ValidationError struct {
	Line   int
	Column int
	Offset int
	Reason string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("JSON validation error at line %d, column %d (offset %d): %s",
		e.Line, e.Column, e.Offset, e.Reason)
}

// ValidationResult contains the result of JSON validation
type ValidationResult struct {
	Valid  bool
	Errors []*ValidationError
}

// Validate performs comprehensive JSON validation
func Validate(data []byte) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: make([]*ValidationError, 0),
	}

	if len(data) == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, &ValidationError{
			Line:   1,
			Column: 1,
			Offset: 0,
			Reason: "empty JSON",
		})
		return result
	}

	// First, use Go's built-in validator
	if !json.Valid(data) {
		result.Valid = false

		// Try to parse to get more detailed error
		var temp interface{}
		if err := json.Unmarshal(data, &temp); err != nil {
			line, col, offset := findErrorPosition(data, err.Error())
			result.Errors = append(result.Errors, &ValidationError{
				Line:   line,
				Column: col,
				Offset: offset,
				Reason: err.Error(),
			})
		}
		return result
	}

	// Additional custom validations
	errors := validateStructure(data)
	if len(errors) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, errors...)
	}

	return result
}

// ValidateString validates a JSON string
func ValidateString(s string) *ValidationResult {
	return Validate([]byte(s))
}

// ValidateValue validates a JSON Value
func (v *Value) Validate() *ValidationResult {
	if v == nil || v.data == nil {
		return &ValidationResult{Valid: true, Errors: nil}
	}

	data, err := json.Marshal(v.data)
	if err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []*ValidationError{{
				Line:   1,
				Column: 1,
				Offset: 0,
				Reason: fmt.Sprintf("failed to marshal value: %v", err),
			}},
		}
	}

	return Validate(data)
}

// findErrorPosition attempts to find the line and column of a JSON error
func findErrorPosition(data []byte, errorMsg string) (line, col, offset int) {
	line = 1
	col = 1
	offset = 0

	// Try to extract offset from error message if possible
	// This is a best-effort approach as Go's JSON error messages vary

	for i, b := range data {
		if b == '\n' {
			line++
			col = 1
		} else {
			col++
		}

		// Simple heuristic: if we find common error indicators, stop here
		if strings.Contains(errorMsg, "invalid character") && i < len(data)-1 {
			// Look for the character mentioned in the error
			if strings.Contains(errorMsg, fmt.Sprintf("'%c'", b)) {
				offset = i
				break
			}
		}
	}

	return line, col, offset
}

// validateStructure performs additional structural validation
func validateStructure(data []byte) []*ValidationError {
	var errors []*ValidationError

	// Check for common issues
	str := string(data)

	// Check for trailing commas (not allowed in JSON)
	if strings.Contains(str, ",}") || strings.Contains(str, ",]") {
		errors = append(errors, &ValidationError{
			Line:   1,
			Column: 1,
			Offset: 0,
			Reason: "trailing comma detected",
		})
	}

	// Check for unescaped control characters
	for i, r := range str {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			line, col := getLineColumn(data, i)
			errors = append(errors, &ValidationError{
				Line:   line,
				Column: col,
				Offset: i,
				Reason: fmt.Sprintf("unescaped control character U+%04X", r),
			})
		}
	}

	return errors
}

// getLineColumn calculates line and column from byte offset
func getLineColumn(data []byte, offset int) (line, col int) {
	line = 1
	col = 1

	for i := 0; i < offset && i < len(data); i++ {
		if data[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}

	return line, col
}

// Format formats JSON with default indentation
func Format(data []byte) ([]byte, error) {
	return FormatIndent(data, "", "  ")
}

// FormatIndent formats JSON with custom indentation
func FormatIndent(data []byte, prefix, indent string) ([]byte, error) {
	if !json.Valid(data) {
		return nil, ErrInvalidJSON
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return json.MarshalIndent(v, prefix, indent)
}

// FormatString formats a JSON string
func FormatString(s string) (string, error) {
	formatted, err := Format([]byte(s))
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

// FormatStringIndent formats a JSON string with custom indentation
func FormatStringIndent(s, prefix, indent string) (string, error) {
	formatted, err := FormatIndent([]byte(s), prefix, indent)
	if err != nil {
		return "", err
	}
	return string(formatted), nil
}

// Minify removes all unnecessary whitespace from JSON
func Minify(data []byte) ([]byte, error) {
	if !json.Valid(data) {
		return nil, ErrInvalidJSON
	}

	var buf bytes.Buffer
	if err := json.Compact(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to minify JSON: %w", err)
	}

	return buf.Bytes(), nil
}

// MinifyString removes all unnecessary whitespace from a JSON string
func MinifyString(s string) (string, error) {
	minified, err := Minify([]byte(s))
	if err != nil {
		return "", err
	}
	return string(minified), nil
}

// ValidateSchema validates JSON against a simple schema
type Schema struct {
	Type       string             `json:"type"`
	Properties map[string]*Schema `json:"properties,omitempty"`
	Items      *Schema            `json:"items,omitempty"`
	Required   []string           `json:"required,omitempty"`
	Enum       []interface{}      `json:"enum,omitempty"`
	Minimum    *float64           `json:"minimum,omitempty"`
	Maximum    *float64           `json:"maximum,omitempty"`
	MinLength  *int               `json:"minLength,omitempty"`
	MaxLength  *int               `json:"maxLength,omitempty"`
	Pattern    string             `json:"pattern,omitempty"`
}

// ValidateSchema validates a JSON value against a schema
func (v *Value) ValidateSchema(schema *Schema) *ValidationResult {
	result := &ValidationResult{
		Valid:  true,
		Errors: make([]*ValidationError, 0),
	}

	if v == nil || v.data == nil {
		if schema.Type != "null" {
			result.Valid = false
			result.Errors = append(result.Errors, &ValidationError{
				Line:   1,
				Column: 1,
				Offset: 0,
				Reason: fmt.Sprintf("expected type %s, got null", schema.Type),
			})
		}
		return result
	}

	errors := v.validateAgainstSchema(schema, "")
	if len(errors) > 0 {
		result.Valid = false
		result.Errors = append(result.Errors, errors...)
	}

	return result
}

// validateAgainstSchema recursively validates against schema
func (v *Value) validateAgainstSchema(schema *Schema, path string) []*ValidationError {
	var errors []*ValidationError

	// Type validation
	actualType := v.getJSONType()
	if schema.Type != "" && schema.Type != actualType {
		errors = append(errors, &ValidationError{
			Line:   1,
			Column: 1,
			Offset: 0,
			Reason: fmt.Sprintf("at path '%s': expected type %s, got %s", path, schema.Type, actualType),
		})
		return errors
	}

	// Enum validation
	if len(schema.Enum) > 0 {
		found := false
		for _, enumVal := range schema.Enum {
			if v.data == enumVal {
				found = true
				break
			}
		}
		if !found {
			errors = append(errors, &ValidationError{
				Line:   1,
				Column: 1,
				Offset: 0,
				Reason: fmt.Sprintf("at path '%s': value not in enum", path),
			})
		}
	}

	// Type-specific validations
	switch actualType {
	case "object":
		if schema.Properties != nil {
			obj, _ := v.data.(map[string]interface{})

			// Check required properties
			for _, required := range schema.Required {
				if _, exists := obj[required]; !exists {
					errors = append(errors, &ValidationError{
						Line:   1,
						Column: 1,
						Offset: 0,
						Reason: fmt.Sprintf("at path '%s': missing required property '%s'", path, required),
					})
				}
			}

			// Validate properties
			for key, val := range obj {
				if propSchema, exists := schema.Properties[key]; exists {
					propPath := path + "." + key
					if path == "" {
						propPath = key
					}
					propValue := &Value{data: val}
					propErrors := propValue.validateAgainstSchema(propSchema, propPath)
					errors = append(errors, propErrors...)
				}
			}
		}

	case "array":
		if schema.Items != nil {
			arr, _ := v.data.([]interface{})
			for i, item := range arr {
				itemPath := fmt.Sprintf("%s[%d]", path, i)
				if path == "" {
					itemPath = fmt.Sprintf("[%d]", i)
				}
				itemValue := &Value{data: item}
				itemErrors := itemValue.validateAgainstSchema(schema.Items, itemPath)
				errors = append(errors, itemErrors...)
			}
		}

	case "string":
		str, _ := v.data.(string)
		if schema.MinLength != nil && len(str) < *schema.MinLength {
			errors = append(errors, &ValidationError{
				Line:   1,
				Column: 1,
				Offset: 0,
				Reason: fmt.Sprintf("at path '%s': string too short", path),
			})
		}
		if schema.MaxLength != nil && len(str) > *schema.MaxLength {
			errors = append(errors, &ValidationError{
				Line:   1,
				Column: 1,
				Offset: 0,
				Reason: fmt.Sprintf("at path '%s': string too long", path),
			})
		}

	case "number":
		num, _ := v.GetFloat64()
		if schema.Minimum != nil && num < *schema.Minimum {
			errors = append(errors, &ValidationError{
				Line:   1,
				Column: 1,
				Offset: 0,
				Reason: fmt.Sprintf("at path '%s': number below minimum", path),
			})
		}
		if schema.Maximum != nil && num > *schema.Maximum {
			errors = append(errors, &ValidationError{
				Line:   1,
				Column: 1,
				Offset: 0,
				Reason: fmt.Sprintf("at path '%s': number above maximum", path),
			})
		}
	}

	return errors
}

// getJSONType returns the JSON type of the value
func (v *Value) getJSONType() string {
	if v == nil || v.data == nil {
		return "null"
	}

	switch v.data.(type) {
	case bool:
		return "boolean"
	case float64, int, int64, float32:
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "unknown"
	}
}
