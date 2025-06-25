package lang

import (
	"reflect"
	"testing"
	"time"
)

func TestIsArray(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "slice",
			value:    []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "array",
			value:    [3]int{1, 2, 3},
			expected: true,
		},
		{
			name:     "string",
			value:    "hello",
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsArray(tt.value)
			if result != tt.expected {
				t.Errorf("IsArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsBoolean(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "true",
			value:    true,
			expected: true,
		},
		{
			name:     "false",
			value:    false,
			expected: true,
		},
		{
			name:     "number",
			value:    1,
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBoolean(tt.value)
			if result != tt.expected {
				t.Errorf("IsBoolean() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsDate(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "time.Time",
			value:    time.Now(),
			expected: true,
		},
		{
			name:     "string",
			value:    "2023-01-01",
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDate(tt.value)
			if result != tt.expected {
				t.Errorf("IsDate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "empty slice",
			value:    []int{},
			expected: true,
		},
		{
			name:     "non-empty slice",
			value:    []int{1},
			expected: false,
		},
		{
			name:     "empty string",
			value:    "",
			expected: true,
		},
		{
			name:     "non-empty string",
			value:    "hello",
			expected: false,
		},
		{
			name:     "zero number",
			value:    0,
			expected: true,
		},
		{
			name:     "non-zero number",
			value:    42,
			expected: false,
		},
		{
			name:     "false boolean",
			value:    false,
			expected: true,
		},
		{
			name:     "true boolean",
			value:    true,
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmpty(tt.value)
			if result != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected bool
	}{
		{
			name:     "equal slices",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "different slices",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 4},
			expected: false,
		},
		{
			name:     "equal strings",
			a:        "hello",
			b:        "hello",
			expected: true,
		},
		{
			name:     "different strings",
			a:        "hello",
			b:        "world",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEqual(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("IsEqual() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsFunction(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "function",
			value:    func() {},
			expected: true,
		},
		{
			name:     "string",
			value:    "hello",
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsFunction(tt.value)
			if result != tt.expected {
				t.Errorf("IsFunction() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNumber(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "int",
			value:    42,
			expected: true,
		},
		{
			name:     "float64",
			value:    3.14,
			expected: true,
		},
		{
			name:     "string",
			value:    "42",
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumber(tt.value)
			if result != tt.expected {
				t.Errorf("IsNumber() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsObject(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "map",
			value:    map[string]int{"a": 1},
			expected: true,
		},
		{
			name:     "struct",
			value:    struct{}{},
			expected: true,
		},
		{
			name:     "slice",
			value:    []int{1, 2},
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsObject(tt.value)
			if result != tt.expected {
				t.Errorf("IsObject() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsString(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "string",
			value:    "hello",
			expected: true,
		},
		{
			name:     "number",
			value:    42,
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsString(tt.value)
			if result != tt.expected {
				t.Errorf("IsString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClone(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		checkFn  func(original, cloned interface{}) bool
	}{
		{
			name:  "slice",
			value: []int{1, 2, 3},
			checkFn: func(original, cloned interface{}) bool {
				orig := original.([]int)
				clone := cloned.([]int)
				// Check if values are equal but slices are different
				return reflect.DeepEqual(orig, clone) && 
					   reflect.ValueOf(orig).Pointer() != reflect.ValueOf(clone).Pointer()
			},
		},
		{
			name:  "map",
			value: map[string]int{"a": 1, "b": 2},
			checkFn: func(original, cloned interface{}) bool {
				orig := original.(map[string]int)
				clone := cloned.(map[string]int)
				return reflect.DeepEqual(orig, clone) && 
					   reflect.ValueOf(orig).Pointer() != reflect.ValueOf(clone).Pointer()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := Clone(tt.value)
			if !tt.checkFn(tt.value, cloned) {
				t.Errorf("Clone() failed for %v", tt.value)
			}
		})
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{
			name:     "string",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "int",
			value:    42,
			expected: "42",
		},
		{
			name:     "bool true",
			value:    true,
			expected: "true",
		},
		{
			name:     "bool false",
			value:    false,
			expected: "false",
		},
		{
			name:     "nil",
			value:    nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.value)
			if result != tt.expected {
				t.Errorf("ToString() = %v, want %v", result, tt.expected)
			}
		})
	}
}
