package object

import (
	"fmt"
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
			name:  "set simple value",
			obj:   make(map[string]interface{}),
			path:  "a",
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

func TestClone(t *testing.T) {
	t.Run("clone map", func(t *testing.T) {
		original := map[string]int{"a": 1, "b": 2}
		cloned := Clone(original)

		// Should be equal but different instances
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone() = %v, want %v", cloned, original)
		}

		// Modify original to ensure they're separate
		original["c"] = 3
		if len(cloned) != 2 {
			t.Errorf("Clone() should create separate instance, cloned map was affected")
		}
	})

	t.Run("clone slice", func(t *testing.T) {
		original := []int{1, 2, 3}
		cloned := Clone(original)

		// Should be equal but different instances
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone() = %v, want %v", cloned, original)
		}

		// Modify original to ensure they're separate
		original[0] = 99
		if cloned[0] != 1 {
			t.Errorf("Clone() should create separate instance, cloned slice was affected")
		}
	})

	t.Run("clone array", func(t *testing.T) {
		original := [3]int{1, 2, 3}
		cloned := Clone(original)

		// Should be equal but different instances
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone() = %v, want %v", cloned, original)
		}

		// Modify original to ensure they're separate
		original[0] = 99
		if cloned[0] != 1 {
			t.Errorf("Clone() should create separate instance, cloned array was affected")
		}
	})

	t.Run("clone nil map", func(t *testing.T) {
		var original map[string]int
		cloned := Clone(original)

		if cloned != nil {
			t.Errorf("Clone() of nil map should be nil, got %v", cloned)
		}
	})

	t.Run("clone primitive", func(t *testing.T) {
		original := 42
		cloned := Clone(original)

		if cloned != original {
			t.Errorf("Clone() = %v, want %v", cloned, original)
		}
	})
}

func TestCloneDeep(t *testing.T) {
	t.Run("clone deep nested map", func(t *testing.T) {
		original := map[string]interface{}{
			"a": 1,
			"b": map[string]interface{}{
				"c": 2,
				"d": map[string]interface{}{
					"e": 3,
				},
			},
		}
		cloned := CloneDeep(original)

		// Should be equal but different instances
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("CloneDeep() = %v, want %v", cloned, original)
		}

		// Modify nested value in original
		original["b"].(map[string]interface{})["c"] = 99

		// Cloned should not be affected
		if cloned["b"].(map[string]interface{})["c"] != 2 {
			t.Errorf("CloneDeep() should create deep copy, nested value was affected")
		}
	})

	t.Run("clone deep nested slice", func(t *testing.T) {
		original := [][]int{{1, 2}, {3, 4}}
		cloned := CloneDeep(original)

		// Should be equal but different instances
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("CloneDeep() = %v, want %v", cloned, original)
		}

		// Modify nested value in original
		original[0][0] = 99

		// Cloned should not be affected
		if cloned[0][0] != 1 {
			t.Errorf("CloneDeep() should create deep copy, nested slice was affected")
		}
	})

	t.Run("clone deep struct", func(t *testing.T) {
		type Inner struct {
			Value int
		}
		type Outer struct {
			Name  string
			Inner Inner
		}

		original := Outer{
			Name:  "test",
			Inner: Inner{Value: 42},
		}
		cloned := CloneDeep(original)

		// Should be equal but different instances
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("CloneDeep() = %v, want %v", cloned, original)
		}

		// Modify original
		original.Inner.Value = 99

		// Cloned should not be affected
		if cloned.Inner.Value != 42 {
			t.Errorf("CloneDeep() should create deep copy, nested struct was affected")
		}
	})

	t.Run("clone deep with pointer", func(t *testing.T) {
		value := 42
		original := &value
		cloned := CloneDeep(original)

		// Should point to different memory locations
		if original == cloned {
			t.Errorf("CloneDeep() should create new pointer instance")
		}

		// But values should be equal
		if *original != *cloned {
			t.Errorf("CloneDeep() = %v, want %v", *cloned, *original)
		}

		// Modify original value
		*original = 99

		// Cloned should not be affected
		if *cloned != 42 {
			t.Errorf("CloneDeep() should create deep copy, pointer value was affected")
		}
	})

	t.Run("clone deep nil values", func(t *testing.T) {
		var original map[string]interface{}
		cloned := CloneDeep(original)

		if cloned != nil {
			t.Errorf("CloneDeep() of nil map should be nil, got %v", cloned)
		}
	})
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
			name:     "pick mixed existing and non-existing keys",
			input:    map[string]int{"a": 1, "b": 2, "c": 3},
			keys:     []string{"a", "d", "c"},
			expected: map[string]int{"a": 1, "c": 3},
		},
		{
			name:     "pick all keys",
			input:    map[string]int{"a": 1, "b": 2},
			keys:     []string{"a", "b"},
			expected: map[string]int{"a": 1, "b": 2},
		},
		{
			name:     "pick no keys",
			input:    map[string]int{"a": 1, "b": 2},
			keys:     []string{},
			expected: map[string]int{},
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			keys:     []string{"a", "b"},
			expected: map[string]int{},
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

func TestPickBy(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]int
		predicate func(int, string) bool
		expected  map[string]int
	}{
		{
			name:      "pick by value greater than 1",
			input:     map[string]int{"a": 1, "b": 2, "c": 3},
			predicate: func(v int, k string) bool { return v > 1 },
			expected:  map[string]int{"b": 2, "c": 3},
		},
		{
			name:      "pick by key length",
			input:     map[string]int{"a": 1, "bb": 2, "ccc": 3},
			predicate: func(v int, k string) bool { return len(k) > 1 },
			expected:  map[string]int{"bb": 2, "ccc": 3},
		},
		{
			name:      "pick even values",
			input:     map[string]int{"a": 1, "b": 2, "c": 3, "d": 4},
			predicate: func(v int, k string) bool { return v%2 == 0 },
			expected:  map[string]int{"b": 2, "d": 4},
		},
		{
			name:      "pick none",
			input:     map[string]int{"a": 1, "b": 2, "c": 3},
			predicate: func(v int, k string) bool { return v > 10 },
			expected:  map[string]int{},
		},
		{
			name:      "pick all",
			input:     map[string]int{"a": 1, "b": 2},
			predicate: func(v int, k string) bool { return true },
			expected:  map[string]int{"a": 1, "b": 2},
		},
		{
			name:      "empty map",
			input:     map[string]int{},
			predicate: func(v int, k string) bool { return true },
			expected:  map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PickBy(tt.input, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("PickBy() = %v, want %v", result, tt.expected)
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
			name:     "nil value",
			input:    nil,
			expected: true,
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
			name:     "empty slice",
			input:    []int{},
			expected: true,
		},
		{
			name:     "non-empty slice",
			input:    []int{1, 2, 3},
			expected: false,
		},
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
			name:     "zero int",
			input:    0,
			expected: true,
		},
		{
			name:     "non-zero int",
			input:    42,
			expected: false,
		},
		{
			name:     "false bool",
			input:    false,
			expected: true,
		},
		{
			name:     "true bool",
			input:    true,
			expected: false,
		},
		{
			name:     "zero float",
			input:    0.0,
			expected: true,
		},
		{
			name:     "non-zero float",
			input:    3.14,
			expected: false,
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

func TestIsEqual(t *testing.T) {
	tests := []struct {
		name     string
		a        interface{}
		b        interface{}
		expected bool
	}{
		{
			name:     "equal integers",
			a:        42,
			b:        42,
			expected: true,
		},
		{
			name:     "different integers",
			a:        42,
			b:        43,
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
			name:     "different slice lengths",
			a:        []int{1, 2},
			b:        []int{1, 2, 3},
			expected: false,
		},
		{
			name:     "equal maps",
			a:        map[string]int{"a": 1, "b": 2},
			b:        map[string]int{"a": 1, "b": 2},
			expected: true,
		},
		{
			name:     "different maps",
			a:        map[string]int{"a": 1, "b": 2},
			b:        map[string]int{"a": 1, "b": 3},
			expected: false,
		},
		{
			name:     "different map keys",
			a:        map[string]int{"a": 1},
			b:        map[string]int{"b": 1},
			expected: false,
		},
		{
			name:     "both nil",
			a:        nil,
			b:        nil,
			expected: true,
		},
		{
			name:     "one nil",
			a:        nil,
			b:        42,
			expected: false,
		},
		{
			name:     "nested equal structures",
			a:        [][]int{{1, 2}, {3, 4}},
			b:        [][]int{{1, 2}, {3, 4}},
			expected: true,
		},
		{
			name:     "nested different structures",
			a:        [][]int{{1, 2}, {3, 4}},
			b:        [][]int{{1, 2}, {3, 5}},
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

func TestTransform(t *testing.T) {
	// Test grouping by value
	obj := map[string]int{"a": 1, "b": 2, "c": 1}
	result := Transform(obj, func(result map[string][]string, value int, key string) {
		valueStr := fmt.Sprintf("%d", value)
		if result[valueStr] == nil {
			result[valueStr] = []string{}
		}
		result[valueStr] = append(result[valueStr], key)
	}, map[string][]string{})

	expected := map[string][]string{
		"1": {"a", "c"},
		"2": {"b"},
	}

	// Check that all expected keys exist and have correct values
	for key, expectedValues := range expected {
		actualValues, exists := result[key]
		if !exists {
			t.Errorf("Transform() missing key %s", key)
			continue
		}

		if len(actualValues) != len(expectedValues) {
			t.Errorf("Transform() key %s has %d values, want %d", key, len(actualValues), len(expectedValues))
			continue
		}

		// Check that all expected values are present (order may vary)
		for _, expectedValue := range expectedValues {
			found := false
			for _, actualValue := range actualValues {
				if actualValue == expectedValue {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Transform() key %s missing value %s", key, expectedValue)
			}
		}
	}
}

func TestTransformSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	result := TransformSlice(slice, func(result map[string][]int, value int, index int) {
		key := "even"
		if value%2 != 0 {
			key = "odd"
		}
		if result[key] == nil {
			result[key] = []int{}
		}
		result[key] = append(result[key], value)
	}, map[string][]int{})

	expected := map[string][]int{
		"odd":  {1, 3},
		"even": {2, 4},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("TransformSlice() = %v, want %v", result, expected)
	}
}

func TestInvertBy(t *testing.T) {
	obj := map[string]int{"a": 1, "b": 2, "c": 1}

	result := InvertBy(obj, func(value int) string {
		return fmt.Sprintf("group_%d", value)
	})

	expected := map[string][]string{
		"group_1": {"a", "c"},
		"group_2": {"b"},
	}

	// Check that all expected keys exist and have correct values
	for key, expectedValues := range expected {
		actualValues, exists := result[key]
		if !exists {
			t.Errorf("InvertBy() missing key %s", key)
			continue
		}

		if len(actualValues) != len(expectedValues) {
			t.Errorf("InvertBy() key %s has %d values, want %d", key, len(actualValues), len(expectedValues))
			continue
		}

		// Check that all expected values are present (order may vary)
		for _, expectedValue := range expectedValues {
			found := false
			for _, actualValue := range actualValues {
				if actualValue == expectedValue {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("InvertBy() key %s missing value %s", key, expectedValue)
			}
		}
	}

	// Test empty object
	emptyResult := InvertBy(map[string]int{}, func(value int) string {
		return fmt.Sprintf("group_%d", value)
	})

	if len(emptyResult) != 0 {
		t.Errorf("InvertBy() on empty object should return empty map, got %v", emptyResult)
	}
}
