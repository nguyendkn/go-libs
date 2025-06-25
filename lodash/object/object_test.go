package object

import (
	"reflect"
	"sort"
	"testing"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []string
	}{
		{
			name:     "basic keys",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Keys(tt.input)
			sort.Strings(result)
			sort.Strings(tt.expected)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Keys() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValues(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		expected []int
	}{
		{
			name:     "basic values",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Values(tt.input)
			sort.Ints(result)
			sort.Ints(tt.expected)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Values() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHas(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		key      string
		expected bool
	}{
		{
			name:     "key exists",
			input:    map[string]int{"a": 1, "b": 2},
			key:      "a",
			expected: true,
		},
		{
			name:     "key does not exist",
			input:    map[string]int{"a": 1, "b": 2},
			key:      "c",
			expected: false,
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			key:      "a",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Has(tt.input, tt.key)
			if result != tt.expected {
				t.Errorf("Has() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name         string
		obj          interface{}
		path         string
		defaultValue interface{}
		expected     interface{}
	}{
		{
			name: "nested map access",
			obj: map[string]interface{}{
				"a": map[string]interface{}{
					"b": 2,
				},
			},
			path:         "a.b",
			defaultValue: 0,
			expected:     2,
		},
		{
			name: "path not found",
			obj: map[string]interface{}{
				"a": 1,
			},
			path:         "a.b",
			defaultValue: 0,
			expected:     0,
		},
		{
			name:         "nil object",
			obj:          nil,
			path:         "a.b",
			defaultValue: 0,
			expected:     0,
		},
		{
			name: "simple access",
			obj: map[string]interface{}{
				"a": 1,
			},
			path:         "a",
			defaultValue: 0,
			expected:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Get(tt.obj, tt.path, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name     string
		obj      map[string]interface{}
		path     string
		value    interface{}
		expected map[string]interface{}
		success  bool
	}{
		{
			name:  "set nested value",
			obj:   make(map[string]interface{}),
			path:  "a.b",
			value: 2,
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": 2,
				},
			},
			success: true,
		},
		{
			name: "set simple value",
			obj:  make(map[string]interface{}),
			path: "a",
			value: 1,
			expected: map[string]interface{}{
				"a": 1,
			},
			success: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success := Set(tt.obj, tt.path, tt.value)
			if success != tt.success {
				t.Errorf("Set() success = %v, want %v", success, tt.success)
			}
			if success && !reflect.DeepEqual(tt.obj, tt.expected) {
				t.Errorf("Set() result = %v, want %v", tt.obj, tt.expected)
			}
		})
	}
}

func TestPick(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		keys     []string
		expected map[string]int
	}{
		{
			name:     "pick existing keys",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			keys:     []string{"a", "c"},
			expected: map[string]int{"a": 1, "c": 3},
		},
		{
			name:     "pick non-existing keys",
			input:    map[string]int{"a": 1, "b": 2},
			keys:     []string{"c", "d"},
			expected: map[string]int{},
		},
		{
			name:     "pick mixed keys",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			keys:     []string{"a", "d"},
			expected: map[string]int{"a": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pick(tt.input, tt.keys)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Pick() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestOmit(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		keys     []string
		expected map[string]int
	}{
		{
			name:     "omit existing keys",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			keys:     []string{"a", "c"},
			expected: map[string]int{"b": 2},
		},
		{
			name:     "omit non-existing keys",
			input:    map[string]int{"a": 1, "b": 2},
			keys:     []string{"c", "d"},
			expected: map[string]int{"a": 1, "b": 2},
		},
		{
			name:     "omit all keys",
			input:    map[string]int{"a": 1, "b": 2},
			keys:     []string{"a", "b"},
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Omit(tt.input, tt.keys)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Omit() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	tests := []struct {
		name     string
		dest     map[string]interface{}
		sources  []map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "basic merge",
			dest: map[string]interface{}{"a": 1},
			sources: []map[string]interface{}{
				{"b": 2},
				{"c": 3},
			},
			expected: map[string]interface{}{"a": 1, "b": 2, "c": 3},
		},
		{
			name: "nested merge",
			dest: map[string]interface{}{
				"a": map[string]interface{}{"x": 1},
			},
			sources: []map[string]interface{}{
				{
					"a": map[string]interface{}{"y": 2},
					"b": 3,
				},
			},
			expected: map[string]interface{}{
				"a": map[string]interface{}{"x": 1, "y": 2},
				"b": 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Merge(tt.dest, tt.sources...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Merge() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAssign(t *testing.T) {
	tests := []struct {
		name     string
		dest     map[string]int
		sources  []map[string]int
		expected map[string]int
	}{
		{
			name: "basic assign",
			dest: map[string]int{"a": 1},
			sources: []map[string]int{
				{"b": 2},
				{"c": 3},
			},
			expected: map[string]int{"a": 1, "b": 2, "c": 3},
		},
		{
			name: "overwrite values",
			dest: map[string]int{"a": 1},
			sources: []map[string]int{
				{"a": 2, "b": 3},
			},
			expected: map[string]int{"a": 2, "b": 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Assign(tt.dest, tt.sources...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Assign() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInvert(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name:     "basic invert",
			input:    map[string]string{"a": "1", "b": "2"},
			expected: map[string]string{"1": "a", "2": "b"},
		},
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Invert(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Invert() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMapKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		iteratee func(string) string
		expected map[string]int
	}{
		{
			name:     "append suffix",
			input:    map[string]int{"a": 1, "b": 2},
			iteratee: func(k string) string { return k + "1" },
			expected: map[string]int{"a1": 1, "b1": 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapKeys(tt.input, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MapKeys() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMapValues(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]int
		iteratee func(int) int
		expected map[string]int
	}{
		{
			name:     "double values",
			input:    map[string]int{"a": 1, "b": 2},
			iteratee: func(v int) int { return v * 2 },
			expected: map[string]int{"a": 2, "b": 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapValues(tt.input, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MapValues() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: true,
		},
		{
			name:     "non-empty map",
			input:    map[string]int{"a": 1},
			expected: false,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: true,
		},
		{
			name:     "non-empty slice",
			input:    []int{1},
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "non-empty string",
			input:    "hello",
			expected: false,
		},
		{
			name:     "nil",
			input:    nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsEmpty(tt.input)
			if result != tt.expected {
				t.Errorf("IsEmpty() = %v, want %v", result, tt.expected)
			}
		})
	}
}
