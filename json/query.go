package json

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Query represents a JSON query
type Query struct {
	path       string
	filters    []Filter
	projection []string
}

// Filter represents a query filter
type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

// NewQuery creates a new JSON query
func NewQuery(path string) *Query {
	return &Query{
		path:       path,
		filters:    make([]Filter, 0),
		projection: make([]string, 0),
	}
}

// Where adds a filter to the query
func (q *Query) Where(field, operator string, value interface{}) *Query {
	q.filters = append(q.filters, Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return q
}

// Select adds fields to project in the result
func (q *Query) Select(fields ...string) *Query {
	q.projection = append(q.projection, fields...)
	return q
}

// Execute executes the query on a JSON value
func (q *Query) Execute(v *Value) ([]*Value, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	// Get the base value using path
	var baseValue *Value
	var err error

	if q.path == "" || q.path == "." {
		baseValue = v
	} else {
		baseValue, err = v.GetPath(q.path)
		if err != nil {
			return nil, fmt.Errorf("failed to get path '%s': %w", q.path, err)
		}
	}

	// If base value is an array, process each element
	if baseValue.IsArray() {
		arr, err := baseValue.GetArray()
		if err != nil {
			return nil, err
		}

		var results []*Value
		for _, item := range arr {
			if q.matchesFilters(item) {
				projected := q.applyProjection(item)
				results = append(results, projected)
			}
		}
		return results, nil
	}

	// If base value is an object, check if it matches filters
	if baseValue.IsObject() {
		if q.matchesFilters(baseValue) {
			projected := q.applyProjection(baseValue)
			return []*Value{projected}, nil
		}
		return []*Value{}, nil
	}

	// For primitive values, return as-is if no filters
	if len(q.filters) == 0 {
		return []*Value{baseValue}, nil
	}

	return []*Value{}, nil
}

// matchesFilters checks if a value matches all filters
func (q *Query) matchesFilters(v *Value) bool {
	for _, filter := range q.filters {
		if !q.matchesFilter(v, filter) {
			return false
		}
	}
	return true
}

// matchesFilter checks if a value matches a single filter
func (q *Query) matchesFilter(v *Value, filter Filter) bool {
	fieldValue, err := v.GetPath(filter.Field)
	if err != nil {
		return false
	}

	switch filter.Operator {
	case "=", "==", "eq":
		return q.compareValues(fieldValue.Interface(), filter.Value) == 0
	case "!=", "ne":
		return q.compareValues(fieldValue.Interface(), filter.Value) != 0
	case ">", "gt":
		return q.compareValues(fieldValue.Interface(), filter.Value) > 0
	case ">=", "gte":
		return q.compareValues(fieldValue.Interface(), filter.Value) >= 0
	case "<", "lt":
		return q.compareValues(fieldValue.Interface(), filter.Value) < 0
	case "<=", "lte":
		return q.compareValues(fieldValue.Interface(), filter.Value) <= 0
	case "contains":
		return q.contains(fieldValue, filter.Value)
	case "startswith":
		return q.startsWith(fieldValue, filter.Value)
	case "endswith":
		return q.endsWith(fieldValue, filter.Value)
	case "regex":
		return q.matchesRegex(fieldValue, filter.Value)
	case "in":
		return q.in(fieldValue, filter.Value)
	case "exists":
		return !fieldValue.IsNull()
	default:
		return false
	}
}

// compareValues compares two values
func (q *Query) compareValues(a, b interface{}) int {
	// Convert to comparable types
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)

	// Try numeric comparison first
	if aNum, aErr := strconv.ParseFloat(aStr, 64); aErr == nil {
		if bNum, bErr := strconv.ParseFloat(bStr, 64); bErr == nil {
			if aNum < bNum {
				return -1
			} else if aNum > bNum {
				return 1
			}
			return 0
		}
	}

	// String comparison
	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}
	return 0
}

// contains checks if a value contains another value
func (q *Query) contains(v *Value, target interface{}) bool {
	if v.IsString() {
		str, _ := v.GetString()
		targetStr := fmt.Sprintf("%v", target)
		return strings.Contains(str, targetStr)
	}

	if v.IsArray() {
		arr, _ := v.GetArray()
		for _, item := range arr {
			if q.compareValues(item.Interface(), target) == 0 {
				return true
			}
		}
	}

	return false
}

// startsWith checks if a string value starts with a prefix
func (q *Query) startsWith(v *Value, prefix interface{}) bool {
	if !v.IsString() {
		return false
	}

	str, _ := v.GetString()
	prefixStr := fmt.Sprintf("%v", prefix)
	return strings.HasPrefix(str, prefixStr)
}

// endsWith checks if a string value ends with a suffix
func (q *Query) endsWith(v *Value, suffix interface{}) bool {
	if !v.IsString() {
		return false
	}

	str, _ := v.GetString()
	suffixStr := fmt.Sprintf("%v", suffix)
	return strings.HasSuffix(str, suffixStr)
}

// matchesRegex checks if a string value matches a regex pattern
func (q *Query) matchesRegex(v *Value, pattern interface{}) bool {
	if !v.IsString() {
		return false
	}

	str, _ := v.GetString()
	patternStr := fmt.Sprintf("%v", pattern)

	regex, err := regexp.Compile(patternStr)
	if err != nil {
		return false
	}

	return regex.MatchString(str)
}

// in checks if a value is in a list
func (q *Query) in(v *Value, list interface{}) bool {
	// Convert list to slice if it's not already
	var items []interface{}

	switch l := list.(type) {
	case []interface{}:
		items = l
	case []string:
		for _, item := range l {
			items = append(items, item)
		}
	case []int:
		for _, item := range l {
			items = append(items, item)
		}
	case []float64:
		for _, item := range l {
			items = append(items, item)
		}
	default:
		return false
	}

	for _, item := range items {
		if q.compareValues(v.Interface(), item) == 0 {
			return true
		}
	}

	return false
}

// applyProjection applies field projection to a value
func (q *Query) applyProjection(v *Value) *Value {
	if len(q.projection) == 0 {
		return v
	}

	if !v.IsObject() {
		return v
	}

	obj, err := v.GetObject()
	if err != nil {
		return v
	}

	projected := make(map[string]interface{})
	for _, field := range q.projection {
		if val, exists := obj[field]; exists {
			projected[field] = val.Interface()
		}
	}

	return &Value{data: projected}
}

// Find finds all values matching a path pattern
func (v *Value) Find(pattern string) ([]*Value, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	var results []*Value
	err := v.findRecursive(pattern, "", &results)
	return results, err
}

// findRecursive recursively finds values matching a pattern
func (v *Value) findRecursive(pattern, currentPath string, results *[]*Value) error {
	// Simple pattern matching - supports wildcards
	if matchesPattern(currentPath, pattern) {
		*results = append(*results, v)
	}

	// Recurse into objects and arrays
	if v.IsObject() {
		obj, err := v.GetObject()
		if err != nil {
			return err
		}

		for key, val := range obj {
			newPath := currentPath
			if newPath != "" {
				newPath += "."
			}
			newPath += key

			if err := val.findRecursive(pattern, newPath, results); err != nil {
				return err
			}
		}
	} else if v.IsArray() {
		arr, err := v.GetArray()
		if err != nil {
			return err
		}

		for i, val := range arr {
			newPath := fmt.Sprintf("%s[%d]", currentPath, i)
			if err := val.findRecursive(pattern, newPath, results); err != nil {
				return err
			}
		}
	}

	return nil
}

// matchesPattern checks if a path matches a pattern (supports * wildcard)
func matchesPattern(path, pattern string) bool {
	// Simple wildcard matching
	if pattern == "*" {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return path == pattern
	}

	// Handle array wildcard patterns like "items[*].name"
	if strings.Contains(pattern, "[*]") {
		// Replace [*] with [\\d+] for regex matching
		regexPattern := strings.ReplaceAll(pattern, "[*]", "\\[\\d+\\]")
		regexPattern = strings.ReplaceAll(regexPattern, "*", ".*")
		regex, err := regexp.Compile("^" + regexPattern + "$")
		if err != nil {
			return false
		}
		return regex.MatchString(path)
	}

	// Handle patterns with array indices like "items[0].*"
	if strings.Contains(pattern, "[") && strings.Contains(pattern, "]") {
		// Escape brackets for regex
		regexPattern := strings.ReplaceAll(pattern, "[", "\\[")
		regexPattern = strings.ReplaceAll(regexPattern, "]", "\\]")
		regexPattern = strings.ReplaceAll(regexPattern, "*", ".*")
		regex, err := regexp.Compile("^" + regexPattern + "$")
		if err != nil {
			return false
		}
		return regex.MatchString(path)
	}

	// Convert pattern to regex
	regexPattern := strings.ReplaceAll(pattern, "*", ".*")
	regex, err := regexp.Compile("^" + regexPattern + "$")
	if err != nil {
		return false
	}

	return regex.MatchString(path)
}

// Extract extracts values from multiple paths
func (v *Value) Extract(paths ...string) (map[string]*Value, error) {
	if v == nil {
		return nil, ErrNilValue
	}

	result := make(map[string]*Value)

	for _, path := range paths {
		val, err := v.GetPath(path)
		if err != nil {
			// Continue with other paths even if one fails
			result[path] = &Value{data: nil}
		} else {
			result[path] = val
		}
	}

	return result, nil
}

// Transform applies a transformation function to all values matching a pattern
func (v *Value) Transform(pattern string, transformer func(*Value) *Value) error {
	if v == nil {
		return ErrNilValue
	}

	return v.transformRecursive(pattern, "", transformer)
}

// transformRecursive recursively transforms values matching a pattern
func (v *Value) transformRecursive(pattern, currentPath string, transformer func(*Value) *Value) error {
	// Apply transformation if pattern matches
	if matchesPattern(currentPath, pattern) {
		transformed := transformer(v)
		if transformed != nil {
			v.data = transformed.data
		}
	}

	// Recurse into objects and arrays
	if v.IsObject() {
		obj, ok := v.data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%w: value is not an object", ErrTypeConversion)
		}

		for key, val := range obj {
			newPath := currentPath
			if newPath != "" {
				newPath += "."
			}
			newPath += key

			valValue := &Value{data: val}
			if err := valValue.transformRecursive(pattern, newPath, transformer); err != nil {
				return err
			}
			// Update the object with potentially transformed value
			obj[key] = valValue.data
		}
	} else if v.IsArray() {
		arr, ok := v.data.([]interface{})
		if !ok {
			return fmt.Errorf("%w: value is not an array", ErrTypeConversion)
		}

		for i, val := range arr {
			newPath := fmt.Sprintf("%s[%d]", currentPath, i)
			valValue := &Value{data: val}
			if err := valValue.transformRecursive(pattern, newPath, transformer); err != nil {
				return err
			}
			// Update the array with potentially transformed value
			arr[i] = valValue.data
		}
	}

	return nil
}
