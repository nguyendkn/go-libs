package lang

import (
	"fmt"
	"reflect"
	"regexp"
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
		name    string
		value   interface{}
		checkFn func(original, cloned interface{}) bool
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

func TestIsError(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "error",
			value:    fmt.Errorf("test error"),
			expected: true,
		},
		{
			name:     "string",
			value:    "error string",
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
		{
			name:     "number",
			value:    42,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsError(tt.value)
			if result != tt.expected {
				t.Errorf("IsError() = %v, want %v", result, tt.expected)
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
			name:     "function with params",
			value:    func(int) string { return "" },
			expected: true,
		},
		{
			name:     "string",
			value:    "not a function",
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
		{
			name:     "number",
			value:    42,
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

func TestIsInteger(t *testing.T) {
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
			name:     "int64",
			value:    int64(42),
			expected: true,
		},
		{
			name:     "float with integer value",
			value:    42.0,
			expected: true,
		},
		{
			name:     "float with decimal",
			value:    3.14,
			expected: false,
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
			result := IsInteger(tt.value)
			if result != tt.expected {
				t.Errorf("IsInteger() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsFloat(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "float32",
			value:    float32(3.14),
			expected: true,
		},
		{
			name:     "float64",
			value:    3.14,
			expected: true,
		},
		{
			name:     "int",
			value:    42,
			expected: false,
		},
		{
			name:     "string",
			value:    "3.14",
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
			result := IsFloat(tt.value)
			if result != tt.expected {
				t.Errorf("IsFloat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsPlainObject(t *testing.T) {
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
			value:    struct{ Name string }{Name: "test"},
			expected: true,
		},
		{
			name:     "pointer to struct",
			value:    &struct{ Name string }{Name: "test"},
			expected: true,
		},
		{
			name:     "time.Time (excluded)",
			value:    time.Now(),
			expected: false,
		},
		{
			name:     "slice",
			value:    []int{1, 2, 3},
			expected: false,
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
			result := IsPlainObject(tt.value)
			if result != tt.expected {
				t.Errorf("IsPlainObject() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsMap(t *testing.T) {
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
			name:     "empty map",
			value:    map[string]int{},
			expected: true,
		},
		{
			name:     "struct",
			value:    struct{}{},
			expected: false,
		},
		{
			name:     "slice",
			value:    []int{1, 2, 3},
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
			result := IsMap(tt.value)
			if result != tt.expected {
				t.Errorf("IsMap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToNumber(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected float64
	}{
		{
			name:     "string number",
			value:    "42",
			expected: 42.0,
		},
		{
			name:     "string float",
			value:    "3.14",
			expected: 3.14,
		},
		{
			name:     "int",
			value:    42,
			expected: 42.0,
		},
		{
			name:     "float",
			value:    3.14,
			expected: 3.14,
		},
		{
			name:     "true",
			value:    true,
			expected: 1.0,
		},
		{
			name:     "false",
			value:    false,
			expected: 0.0,
		},
		{
			name:     "empty string",
			value:    "",
			expected: 0.0,
		},
		{
			name:     "invalid string",
			value:    "abc",
			expected: 0.0,
		},
		{
			name:     "nil",
			value:    nil,
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToNumber(tt.value)
			if result != tt.expected {
				t.Errorf("ToNumber() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToInteger(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected int64
	}{
		{
			name:     "string number",
			value:    "42",
			expected: 42,
		},
		{
			name:     "string float",
			value:    "3.14",
			expected: 3,
		},
		{
			name:     "int",
			value:    42,
			expected: 42,
		},
		{
			name:     "float",
			value:    3.14,
			expected: 3,
		},
		{
			name:     "true",
			value:    true,
			expected: 1,
		},
		{
			name:     "false",
			value:    false,
			expected: 0,
		},
		{
			name:     "nil",
			value:    nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToInteger(tt.value)
			if result != tt.expected {
				t.Errorf("ToInteger() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsRegExp(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "regexp",
			value:    regexp.MustCompile("test"),
			expected: true,
		},
		{
			name:     "regexp with pattern",
			value:    regexp.MustCompile(`\d+`),
			expected: true,
		},
		{
			name:     "string",
			value:    "test",
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
		{
			name:     "number",
			value:    42,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRegExp(tt.value)
			if result != tt.expected {
				t.Errorf("IsRegExp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsSymbol(t *testing.T) {
	// Define a custom Symbol type for testing
	type Symbol string
	type CustomStruct struct {
		name string
	}
	type TestSymbol struct {
		value string
	}

	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "custom Symbol type",
			value:    Symbol("test"),
			expected: true,
		},
		{
			name:     "struct with Symbol suffix",
			value:    TestSymbol{value: "test"},
			expected: true,
		},
		{
			name:     "struct without Symbol",
			value:    CustomStruct{name: "test"},
			expected: false,
		},
		{
			name:     "string",
			value:    "test",
			expected: false,
		},
		{
			name:     "nil",
			value:    nil,
			expected: false,
		},
		{
			name:     "number",
			value:    42,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSymbol(tt.value)
			if result != tt.expected {
				t.Errorf("IsSymbol() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsArrayBuffer(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{
			name:     "byte slice",
			value:    []byte{1, 2, 3},
			expected: true,
		},
		{
			name:     "empty byte slice",
			value:    []byte{},
			expected: true,
		},
		{
			name:     "int slice",
			value:    []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "string",
			value:    "test",
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
			result := IsArrayBuffer(tt.value)
			if result != tt.expected {
				t.Errorf("IsArrayBuffer() = %v, want %v", result, tt.expected)
			}
		})
	}
}
