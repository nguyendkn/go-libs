package array

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		size     int
		expected [][]int
	}{
		{
			name:     "basic chunking",
			slice:    []int{1, 2, 3, 4, 5},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:     "exact division",
			slice:    []int{1, 2, 3, 4},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}},
		},
		{
			name:     "size larger than slice",
			slice:    []int{1, 2},
			size:     5,
			expected: [][]int{{1, 2}},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			size:     2,
			expected: [][]int{},
		},
		{
			name:     "size zero",
			slice:    []int{1, 2, 3},
			size:     0,
			expected: [][]int{},
		},
		{
			name:     "size negative",
			slice:    []int{1, 2, 3},
			size:     -1,
			expected: [][]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Chunk(tt.slice, tt.size)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Chunk() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCompact(t *testing.T) {
	tests := []struct {
		name     string
		slice    []interface{}
		expected []interface{}
	}{
		{
			name:     "mixed falsey values",
			slice:    []interface{}{0, 1, false, 2, "", 3, nil},
			expected: []interface{}{1, 2, 3},
		},
		{
			name:     "all truthy",
			slice:    []interface{}{1, 2, 3, "hello"},
			expected: []interface{}{1, 2, 3, "hello"},
		},
		{
			name:     "all falsey",
			slice:    []interface{}{0, false, "", nil},
			expected: []interface{}{},
		},
		{
			name:     "empty slice",
			slice:    []interface{}{},
			expected: []interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Compact(tt.slice)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Compact() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConcat(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		others   [][]int
		expected []int
	}{
		{
			name:     "basic concat",
			slice:    []int{1},
			others:   [][]int{{2, 3}, {4}},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			others:   [][]int{{1, 2}},
			expected: []int{1, 2},
		},
		{
			name:     "no others",
			slice:    []int{1, 2},
			others:   [][]int{},
			expected: []int{1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Concat(tt.slice, tt.others...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Concat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		others   [][]int
		expected []int
	}{
		{
			name:     "basic difference",
			slice:    []int{2, 1},
			others:   [][]int{{2, 3}},
			expected: []int{1},
		},
		{
			name:     "multiple others",
			slice:    []int{1, 2, 3, 4},
			others:   [][]int{{2}, {3, 4}},
			expected: []int{1},
		},
		{
			name:     "no difference",
			slice:    []int{1, 2},
			others:   [][]int{{3, 4}},
			expected: []int{1, 2},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			others:   [][]int{{1, 2}},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Difference(tt.slice, tt.others...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Difference() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDrop(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		expected []int
	}{
		{
			name:     "drop one",
			slice:    []int{1, 2, 3},
			n:        1,
			expected: []int{2, 3},
		},
		{
			name:     "drop multiple",
			slice:    []int{1, 2, 3},
			n:        2,
			expected: []int{3},
		},
		{
			name:     "drop more than length",
			slice:    []int{1, 2, 3},
			n:        5,
			expected: []int{},
		},
		{
			name:     "drop zero",
			slice:    []int{1, 2, 3},
			n:        0,
			expected: []int{1, 2, 3},
		},
		{
			name:     "drop negative",
			slice:    []int{1, 2, 3},
			n:        -1,
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Drop(tt.slice, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Drop() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDropRight(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		expected []int
	}{
		{
			name:     "drop right one",
			slice:    []int{1, 2, 3},
			n:        1,
			expected: []int{1, 2},
		},
		{
			name:     "drop right multiple",
			slice:    []int{1, 2, 3},
			n:        2,
			expected: []int{1},
		},
		{
			name:     "drop right more than length",
			slice:    []int{1, 2, 3},
			n:        5,
			expected: []int{},
		},
		{
			name:     "drop right zero",
			slice:    []int{1, 2, 3},
			n:        0,
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DropRight(tt.slice, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("DropRight() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFill(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		value    int
		start    int
		end      int
		expected []int
	}{
		{
			name:     "basic fill",
			slice:    []int{1, 2, 3, 4},
			value:    0,
			start:    1,
			end:      3,
			expected: []int{1, 0, 0, 4},
		},
		{
			name:     "fill all",
			slice:    []int{1, 2, 3},
			value:    9,
			start:    0,
			end:      3,
			expected: []int{9, 9, 9},
		},
		{
			name:     "start greater than end",
			slice:    []int{1, 2, 3},
			value:    0,
			start:    2,
			end:      1,
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]int, len(tt.slice))
			copy(slice, tt.slice)
			Fill(slice, tt.value, tt.start, tt.end)
			if !reflect.DeepEqual(slice, tt.expected) {
				t.Errorf("Fill() = %v, want %v", slice, tt.expected)
			}
		})
	}
}

func TestHead(t *testing.T) {
	tests := []struct {
		name        string
		slice       []int
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "non-empty slice",
			slice:       []int{1, 2, 3},
			expectedVal: 1,
			expectedOk:  true,
		},
		{
			name:        "empty slice",
			slice:       []int{},
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := Head(tt.slice)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("Head() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestLast(t *testing.T) {
	tests := []struct {
		name        string
		slice       []int
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "non-empty slice",
			slice:       []int{1, 2, 3},
			expectedVal: 3,
			expectedOk:  true,
		},
		{
			name:        "empty slice",
			slice:       []int{},
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := Last(tt.slice)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("Last() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{
			name:     "basic reverse",
			slice:    []int{1, 2, 3},
			expected: []int{3, 2, 1},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			slice:    []int{1},
			expected: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]int, len(tt.slice))
			copy(slice, tt.slice)
			Reverse(slice)
			if !reflect.DeepEqual(slice, tt.expected) {
				t.Errorf("Reverse() = %v, want %v", slice, tt.expected)
			}
		})
	}
}

func TestUniq(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{
			name:     "basic uniq",
			slice:    []int{2, 1, 2},
			expected: []int{2, 1},
		},
		{
			name:     "no duplicates",
			slice:    []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "all duplicates",
			slice:    []int{1, 1, 1},
			expected: []int{1},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uniq(tt.slice)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Uniq() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		slice    [][]int
		expected []int
	}{
		{
			name:     "basic flatten",
			slice:    [][]int{{1, 2}, {3, 4}},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "empty slices",
			slice:    [][]int{{}, {1, 2}, {}},
			expected: []int{1, 2},
		},
		{
			name:     "empty input",
			slice:    [][]int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.slice)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Flatten() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFlattenDeep(t *testing.T) {
	tests := []struct {
		name     string
		slice    interface{}
		expected []interface{}
	}{
		{
			name:     "nested arrays",
			slice:    []interface{}{1, []interface{}{2, []interface{}{3, 4}}},
			expected: []interface{}{1, 2, 3, 4},
		},
		{
			name:     "deeply nested",
			slice:    []interface{}{[]interface{}{[]interface{}{1, 2}}, []interface{}{[]interface{}{3, 4}}},
			expected: []interface{}{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenDeep(tt.slice)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FlattenDeep() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		value    int
		expected int
	}{
		{
			name:     "found at beginning",
			slice:    []int{1, 2, 3, 2},
			value:    1,
			expected: 0,
		},
		{
			name:     "found in middle",
			slice:    []int{1, 2, 3, 2},
			value:    2,
			expected: 1,
		},
		{
			name:     "not found",
			slice:    []int{1, 2, 3},
			value:    4,
			expected: -1,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			value:    1,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IndexOf(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("IndexOf() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLastIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		value    int
		expected int
	}{
		{
			name:     "found at end",
			slice:    []int{1, 2, 3, 2},
			value:    2,
			expected: 3,
		},
		{
			name:     "found once",
			slice:    []int{1, 2, 3},
			value:    2,
			expected: 1,
		},
		{
			name:     "not found",
			slice:    []int{1, 2, 3},
			value:    4,
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LastIndexOf(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("LastIndexOf() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		name      string
		slice     []interface{}
		separator string
		expected  string
	}{
		{
			name:      "string slice",
			slice:     []interface{}{"a", "b", "c"},
			separator: "~",
			expected:  "a~b~c",
		},
		{
			name:      "int slice",
			slice:     []interface{}{1, 2, 3},
			separator: "-",
			expected:  "1-2-3",
		},
		{
			name:      "empty slice",
			slice:     []interface{}{},
			separator: ",",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Join(tt.slice, tt.separator)
			if result != tt.expected {
				t.Errorf("Join() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		start    int
		end      int
		expected []int
	}{
		{
			name:     "basic slice",
			slice:    []int{1, 2, 3, 4},
			start:    1,
			end:      3,
			expected: []int{2, 3},
		},
		{
			name:     "negative end",
			slice:    []int{1, 2, 3, 4},
			start:    2,
			end:      -1,
			expected: []int{3},
		},
		{
			name:     "out of bounds",
			slice:    []int{1, 2, 3},
			start:    1,
			end:      10,
			expected: []int{2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Slice(tt.slice, tt.start, tt.end)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Slice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTake(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		expected []int
	}{
		{
			name:     "take some",
			slice:    []int{1, 2, 3},
			n:        2,
			expected: []int{1, 2},
		},
		{
			name:     "take more than length",
			slice:    []int{1, 2, 3},
			n:        5,
			expected: []int{1, 2, 3},
		},
		{
			name:     "take zero",
			slice:    []int{1, 2, 3},
			n:        0,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Take(tt.slice, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Take() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTakeRight(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		expected []int
	}{
		{
			name:     "take right some",
			slice:    []int{1, 2, 3},
			n:        2,
			expected: []int{2, 3},
		},
		{
			name:     "take right more than length",
			slice:    []int{1, 2, 3},
			n:        5,
			expected: []int{1, 2, 3},
		},
		{
			name:     "take right zero",
			slice:    []int{1, 2, 3},
			n:        0,
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TakeRight(tt.slice, tt.n)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TakeRight() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWithout(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		values   []int
		expected []int
	}{
		{
			name:     "remove some values",
			slice:    []int{2, 1, 2, 3},
			values:   []int{1, 2},
			expected: []int{3},
		},
		{
			name:     "remove no values",
			slice:    []int{1, 2, 3},
			values:   []int{4, 5},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			values:   []int{1, 2},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Without(tt.slice, tt.values...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Without() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFlattenDepth(t *testing.T) {
	tests := []struct {
		name     string
		slice    interface{}
		depth    int
		expected []interface{}
	}{
		{
			name:     "depth 1 - 2D array",
			slice:    [][]int{{1, 2}, {3, 4}},
			depth:    1,
			expected: []interface{}{1, 2, 3, 4},
		},
		{
			name:     "depth 2 - 3D array",
			slice:    [][][]int{{{1, 2}}, {{3, 4}}},
			depth:    2,
			expected: []interface{}{1, 2, 3, 4},
		},
		{
			name:     "depth 1 - 3D array (partial flatten)",
			slice:    [][][]int{{{1, 2}}, {{3, 4}}},
			depth:    1,
			expected: []interface{}{[]int{1, 2}, []int{3, 4}},
		},
		{
			name:     "depth 0 - no flattening",
			slice:    [][]int{{1, 2}, {3, 4}},
			depth:    0,
			expected: []interface{}{[][]int{{1, 2}, {3, 4}}},
		},
		{
			name:     "negative depth",
			slice:    [][]int{{1, 2}, {3, 4}},
			depth:    -1,
			expected: []interface{}{[][]int{{1, 2}, {3, 4}}},
		},
		{
			name:     "mixed types",
			slice:    []interface{}{[]int{1, 2}, "hello", []string{"a", "b"}},
			depth:    1,
			expected: []interface{}{1, 2, "hello", "a", "b"},
		},
		{
			name:     "empty slice",
			slice:    [][]int{},
			depth:    1,
			expected: []interface{}{},
		},
		{
			name:     "non-slice input",
			slice:    42,
			depth:    1,
			expected: []interface{}{42},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FlattenDepth(tt.slice, tt.depth)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FlattenDepth() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFromPairs(t *testing.T) {
	tests := []struct {
		name     string
		pairs    [][2]interface{}
		expected map[interface{}]interface{}
	}{
		{
			name:     "string keys with mixed values",
			pairs:    [][2]interface{}{{"a", 1}, {"b", "hello"}, {"c", true}},
			expected: map[interface{}]interface{}{"a": 1, "b": "hello", "c": true},
		},
		{
			name:     "numeric keys",
			pairs:    [][2]interface{}{{1, "one"}, {2, "two"}, {3, "three"}},
			expected: map[interface{}]interface{}{1: "one", 2: "two", 3: "three"},
		},
		{
			name:     "empty pairs",
			pairs:    [][2]interface{}{},
			expected: map[interface{}]interface{}{},
		},
		{
			name:     "duplicate keys - last wins",
			pairs:    [][2]interface{}{{"a", 1}, {"b", 2}, {"a", 3}},
			expected: map[interface{}]interface{}{"a": 3, "b": 2},
		},
		{
			name:     "mixed key types",
			pairs:    [][2]interface{}{{"string", 1}, {42, "number"}, {true, "boolean"}},
			expected: map[interface{}]interface{}{"string": 1, 42: "number", true: "boolean"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromPairs(tt.pairs)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FromPairs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFromPairsString(t *testing.T) {
	tests := []struct {
		name     string
		pairs    [][2]interface{}
		expected map[string]interface{}
	}{
		{
			name:     "string keys with mixed values",
			pairs:    [][2]interface{}{{"a", 1}, {"b", "hello"}, {"c", true}},
			expected: map[string]interface{}{"a": 1, "b": "hello", "c": true},
		},
		{
			name:     "empty pairs",
			pairs:    [][2]interface{}{},
			expected: map[string]interface{}{},
		},
		{
			name:     "duplicate keys - last wins",
			pairs:    [][2]interface{}{{"a", 1}, {"b", 2}, {"a", 3}},
			expected: map[string]interface{}{"a": 3, "b": 2},
		},
		{
			name:     "mixed key types - only strings accepted",
			pairs:    [][2]interface{}{{"string", 1}, {42, "number"}, {"valid", "value"}},
			expected: map[string]interface{}{"string": 1, "valid": "value"},
		},
		{
			name:     "non-string keys ignored",
			pairs:    [][2]interface{}{{1, "one"}, {true, "boolean"}, {"valid", "accepted"}},
			expected: map[string]interface{}{"valid": "accepted"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FromPairsString(tt.pairs)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("FromPairsString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntersectionBy(t *testing.T) {
	tests := []struct {
		name     string
		iteratee func(int) int
		slices   [][]int
		expected []int
	}{
		{
			name:     "identity function",
			iteratee: func(x int) int { return x },
			slices:   [][]int{{2, 1}, {2, 3}},
			expected: []int{2},
		},
		{
			name: "absolute value function",
			iteratee: func(x int) int {
				if x < 0 {
					return -x
				}
				return x
			},
			slices:   [][]int{{-2, 1}, {2, 3}},
			expected: []int{-2},
		},
		{
			name:     "modulo function - has intersection",
			iteratee: func(x int) int { return x % 3 },
			slices:   [][]int{{1, 4, 7}, {4, 5, 8}, {7, 6, 10}},
			expected: []int{1}, // 1%3=1, 4%3=1, 7%3=1 from first slice; 4%3=1 from second; 7%3=1, 10%3=1 from third
		},
		{
			name:     "empty slices",
			iteratee: func(x int) int { return x },
			slices:   [][]int{},
			expected: []int{},
		},
		{
			name:     "single slice",
			iteratee: func(x int) int { return x },
			slices:   [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "no intersection",
			iteratee: func(x int) int { return x },
			slices:   [][]int{{1, 2}, {3, 4}},
			expected: []int{}, // no common values
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntersectionBy(tt.iteratee, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IntersectionBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntersectionByString(t *testing.T) {
	tests := []struct {
		name     string
		iteratee func(string) int
		slices   [][]string
		expected []string
	}{
		{
			name:     "by length",
			iteratee: func(s string) int { return len(s) },
			slices:   [][]string{{"a", "bb", "ccc"}, {"dd", "e", "fff"}},
			expected: []string{"a", "bb", "ccc"}, // lengths 1,2,3 all appear in both slices
		},
		{
			name: "by first character",
			iteratee: func(s string) int {
				if len(s) == 0 {
					return 0
				}
				return int(s[0])
			},
			slices:   [][]string{{"apple", "banana"}, {"avocado", "cherry"}},
			expected: []string{"apple"},
		},
		{
			name:     "empty strings",
			iteratee: func(s string) int { return len(s) },
			slices:   [][]string{{"", "a"}, {"", "b"}},
			expected: []string{"", "a"}, // length 0 and 1 both appear in both slices, order from first slice
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntersectionBy(tt.iteratee, tt.slices...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IntersectionBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntersectionWith(t *testing.T) {
	tests := []struct {
		name       string
		comparator func(int, int) bool
		slices     [][]int
		expected   []int
	}{
		{
			name:       "equality comparator",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{{2, 1}, {2, 3}},
			expected:   []int{2},
		},
		{
			name:       "absolute value comparator",
			comparator: func(a, b int) bool { return abs(a) == abs(b) },
			slices:     [][]int{{-2, 1}, {2, 3}},
			expected:   []int{-2},
		},
		{
			name:       "modulo comparator",
			comparator: func(a, b int) bool { return a%3 == b%3 },
			slices:     [][]int{{1, 4, 7}, {4, 5, 8}, {7, 6, 10}},
			expected:   []int{1}, // 1%3=1, matches with 4%3=1 and 7%3=1, 10%3=1
		},
		{
			name:       "empty slices",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{},
			expected:   []int{},
		},
		{
			name:       "single slice",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{{1, 2, 2, 3}},
			expected:   []int{1, 2, 3}, // unique elements
		},
		{
			name:       "no intersection",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{{1, 2}, {3, 4}},
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntersectionWith(tt.comparator, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IntersectionWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntersectionWithString(t *testing.T) {
	tests := []struct {
		name       string
		comparator func(string, string) bool
		slices     [][]string
		expected   []string
	}{
		{
			name:       "by length",
			comparator: func(a, b string) bool { return len(a) == len(b) },
			slices:     [][]string{{"a", "bb", "ccc"}, {"dd", "e", "fff"}},
			expected:   []string{"a", "bb", "ccc"}, // all lengths have matches
		},
		{
			name:       "case insensitive",
			comparator: func(a, b string) bool { return strings.EqualFold(a, b) },
			slices:     [][]string{{"Apple", "banana"}, {"APPLE", "cherry"}},
			expected:   []string{"Apple"}, // "Apple" matches "APPLE"
		},
		{
			name: "by first character",
			comparator: func(a, b string) bool {
				if len(a) == 0 || len(b) == 0 {
					return len(a) == len(b)
				}
				return a[0] == b[0]
			},
			slices:   [][]string{{"apple", "banana"}, {"avocado", "cherry"}},
			expected: []string{"apple"}, // "apple" and "avocado" both start with 'a'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntersectionWith(tt.comparator, tt.slices...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("IntersectionWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func TestPullAll(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		values   []int
		expected []int
	}{
		{
			name:     "remove multiple values",
			slice:    []int{1, 2, 3, 1, 2, 3},
			values:   []int{2, 3},
			expected: []int{1, 1},
		},
		{
			name:     "remove single value",
			slice:    []int{1, 2, 3, 4, 5},
			values:   []int{3},
			expected: []int{1, 2, 4, 5},
		},
		{
			name:     "remove non-existent values",
			slice:    []int{1, 2, 3},
			values:   []int{4, 5},
			expected: []int{1, 2, 3},
		},
		{
			name:     "remove all values",
			slice:    []int{1, 2, 3},
			values:   []int{1, 2, 3},
			expected: []int{},
		},
		{
			name:     "empty values",
			slice:    []int{1, 2, 3},
			values:   []int{},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			values:   []int{1, 2},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]int, len(tt.slice))
			copy(slice, tt.slice)
			PullAll(&slice, tt.values)
			if !reflect.DeepEqual(slice, tt.expected) {
				t.Errorf("PullAll() = %v, want %v", slice, tt.expected)
			}
		})
	}
}

func TestPullAllBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		values   []int
		iteratee func(int) int
		expected []int
	}{
		{
			name:     "remove by modulo",
			slice:    []int{1, 2, 3, 4, 5, 6},
			values:   []int{2, 4}, // even numbers
			iteratee: func(x int) int { return x % 2 },
			expected: []int{1, 3, 5}, // removes all even numbers
		},
		{
			name:   "remove by absolute value",
			slice:  []int{-1, 2, -3, 4, -5},
			values: []int{1, 3, 5}, // absolute values 1, 3, and 5
			iteratee: func(x int) int {
				if x < 0 {
					return -x
				}
				return x
			},
			expected: []int{2, 4}, // removes -1, -3, -5
		},
		{
			name:     "identity function",
			slice:    []int{1, 2, 3, 4, 5},
			values:   []int{2, 4},
			iteratee: func(x int) int { return x },
			expected: []int{1, 3, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]int, len(tt.slice))
			copy(slice, tt.slice)
			PullAllBy(&slice, tt.values, tt.iteratee)
			if !reflect.DeepEqual(slice, tt.expected) {
				t.Errorf("PullAllBy() = %v, want %v", slice, tt.expected)
			}
		})
	}
}

func TestPullAllWith(t *testing.T) {
	tests := []struct {
		name       string
		slice      []int
		values     []int
		comparator func(int, int) bool
		expected   []int
	}{
		{
			name:       "equality comparator",
			slice:      []int{1, 2, 3, 4, 5},
			values:     []int{2, 4},
			comparator: func(a, b int) bool { return a == b },
			expected:   []int{1, 3, 5},
		},
		{
			name:       "same parity comparator",
			slice:      []int{1, 2, 3, 4, 5, 6},
			values:     []int{2}, // even number
			comparator: func(a, b int) bool { return a%2 == b%2 },
			expected:   []int{1, 3, 5}, // removes all even numbers
		},
		{
			name:       "absolute value comparator",
			slice:      []int{-1, 2, -3, 4, -5},
			values:     []int{1, 3, 5}, // absolute values 1, 3, and 5
			comparator: func(a, b int) bool { return abs(a) == abs(b) },
			expected:   []int{2, 4}, // removes -1, -3, -5
		},
		{
			name:       "no matches",
			slice:      []int{1, 2, 3},
			values:     []int{4, 5},
			comparator: func(a, b int) bool { return a == b },
			expected:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]int, len(tt.slice))
			copy(slice, tt.slice)
			PullAllWith(&slice, tt.values, tt.comparator)
			if !reflect.DeepEqual(slice, tt.expected) {
				t.Errorf("PullAllWith() = %v, want %v", slice, tt.expected)
			}
		})
	}
}

func TestPullAt(t *testing.T) {
	tests := []struct {
		name            string
		slice           []string
		indexes         []int
		expectedSlice   []string
		expectedRemoved []string
	}{
		{
			name:            "remove at specific indexes",
			slice:           []string{"a", "b", "c", "d", "e"},
			indexes:         []int{1, 3},
			expectedSlice:   []string{"a", "c", "e"},
			expectedRemoved: []string{"b", "d"},
		},
		{
			name:            "remove at negative indexes",
			slice:           []string{"a", "b", "c", "d"},
			indexes:         []int{-1, -3},
			expectedSlice:   []string{"a", "c"},
			expectedRemoved: []string{"b", "d"},
		},
		{
			name:            "remove at out of bounds indexes",
			slice:           []string{"a", "b", "c"},
			indexes:         []int{1, 5, -5},
			expectedSlice:   []string{"a", "c"},
			expectedRemoved: []string{"b"},
		},
		{
			name:            "remove all elements",
			slice:           []string{"a", "b", "c"},
			indexes:         []int{0, 1, 2},
			expectedSlice:   []string{},
			expectedRemoved: []string{"a", "b", "c"},
		},
		{
			name:            "no indexes provided",
			slice:           []string{"a", "b", "c"},
			indexes:         []int{},
			expectedSlice:   []string{"a", "b", "c"},
			expectedRemoved: []string{},
		},
		{
			name:            "empty slice",
			slice:           []string{},
			indexes:         []int{0, 1},
			expectedSlice:   []string{},
			expectedRemoved: []string{},
		},
		{
			name:            "duplicate indexes",
			slice:           []string{"a", "b", "c", "d"},
			indexes:         []int{1, 1, 2},
			expectedSlice:   []string{"a", "d"},
			expectedRemoved: []string{"b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slice := make([]string, len(tt.slice))
			copy(slice, tt.slice)
			removed := PullAt(&slice, tt.indexes...)

			// Handle empty slice comparison
			if len(slice) == 0 && len(tt.expectedSlice) == 0 {
				// Both are empty, check removed
			} else if !reflect.DeepEqual(slice, tt.expectedSlice) {
				t.Errorf("PullAt() slice = %v, want %v", slice, tt.expectedSlice)
			}

			if len(removed) == 0 && len(tt.expectedRemoved) == 0 {
				// Both are empty, test passes
			} else if !reflect.DeepEqual(removed, tt.expectedRemoved) {
				t.Errorf("PullAt() removed = %v, want %v", removed, tt.expectedRemoved)
			}
		})
	}
}

func TestTail(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{
			name:     "normal slice",
			slice:    []int{1, 2, 3, 4, 5},
			expected: []int{2, 3, 4, 5},
		},
		{
			name:     "two elements",
			slice:    []int{1, 2},
			expected: []int{2},
		},
		{
			name:     "single element",
			slice:    []int{1},
			expected: []int{},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Tail(tt.slice)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Tail() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniqBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		iteratee func(int) int
		expected []int
	}{
		{
			name:     "identity function",
			slice:    []int{2, 1, 2, 3, 1},
			iteratee: func(x int) int { return x },
			expected: []int{2, 1, 3},
		},
		{
			name:     "modulo function",
			slice:    []int{1, 2, 3, 4, 5, 6},
			iteratee: func(x int) int { return x % 3 },
			expected: []int{1, 2, 3}, // 1%3=1, 2%3=2, 3%3=0, 4%3=1(dup), 5%3=2(dup), 6%3=0(dup)
		},
		{
			name:  "absolute value function",
			slice: []int{-1, 1, -2, 2, -3},
			iteratee: func(x int) int {
				if x < 0 {
					return -x
				}
				return x
			},
			expected: []int{-1, -2, -3}, // abs(-1)=1, abs(1)=1(dup), abs(-2)=2, abs(2)=2(dup), abs(-3)=3
		},
		{
			name:     "empty slice",
			slice:    []int{},
			iteratee: func(x int) int { return x },
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqBy(tt.slice, tt.iteratee)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UniqBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniqByString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		iteratee func(string) int
		expected []string
	}{
		{
			name:     "by length",
			slice:    []string{"a", "bb", "c", "dd", "eee"},
			iteratee: func(s string) int { return len(s) },
			expected: []string{"a", "bb", "eee"}, // lengths: 1, 2, 1(dup), 2(dup), 3
		},
		{
			name:  "by first character",
			slice: []string{"apple", "banana", "avocado", "cherry"},
			iteratee: func(s string) int {
				if len(s) == 0 {
					return 0
				}
				return int(s[0])
			},
			expected: []string{"apple", "banana", "cherry"}, // 'a', 'b', 'a'(dup), 'c'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqBy(tt.slice, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UniqBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniqWith(t *testing.T) {
	tests := []struct {
		name       string
		slice      []int
		comparator func(int, int) bool
		expected   []int
	}{
		{
			name:       "equality comparator",
			slice:      []int{1, 2, 2, 3, 1},
			comparator: func(a, b int) bool { return a == b },
			expected:   []int{1, 2, 3},
		},
		{
			name:       "same parity comparator",
			slice:      []int{1, 2, 3, 4, 5, 6},
			comparator: func(a, b int) bool { return a%2 == b%2 },
			expected:   []int{1, 2}, // 1(odd), 2(even), 3(odd-dup), 4(even-dup), 5(odd-dup), 6(even-dup)
		},
		{
			name:       "absolute value comparator",
			slice:      []int{-1, 1, -2, 2, -3},
			comparator: func(a, b int) bool { return abs(a) == abs(b) },
			expected:   []int{-1, -2, -3}, // abs(-1)=1, abs(1)=1(dup), abs(-2)=2, abs(2)=2(dup), abs(-3)=3
		},
		{
			name:       "empty slice",
			slice:      []int{},
			comparator: func(a, b int) bool { return a == b },
			expected:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqWith(tt.slice, tt.comparator)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UniqWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniqWithString(t *testing.T) {
	tests := []struct {
		name       string
		slice      []string
		comparator func(string, string) bool
		expected   []string
	}{
		{
			name:       "case insensitive",
			slice:      []string{"a", "A", "b", "B", "c"},
			comparator: func(a, b string) bool { return strings.EqualFold(a, b) },
			expected:   []string{"a", "b", "c"},
		},
		{
			name:       "by length",
			slice:      []string{"a", "bb", "c", "dd", "eee"},
			comparator: func(a, b string) bool { return len(a) == len(b) },
			expected:   []string{"a", "bb", "eee"}, // lengths: 1, 2, 1(dup), 2(dup), 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqWith(tt.slice, tt.comparator)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UniqWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedIndex(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		value    int
		expected int
	}{
		{
			name:     "insert in middle",
			slice:    []int{30, 50},
			value:    40,
			expected: 1,
		},
		{
			name:     "insert at beginning",
			slice:    []int{30, 50},
			value:    20,
			expected: 0,
		},
		{
			name:     "insert at end",
			slice:    []int{30, 50},
			value:    60,
			expected: 2,
		},
		{
			name:     "duplicate values - first position",
			slice:    []int{4, 5, 5, 5, 6},
			value:    5,
			expected: 1,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			value:    5,
			expected: 0,
		},
		{
			name:     "single element - before",
			slice:    []int{5},
			value:    3,
			expected: 0,
		},
		{
			name:     "single element - after",
			slice:    []int{5},
			value:    7,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedIndex(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("SortedIndex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedIndexBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		value    string
		iteratee func(string) int
		expected int
	}{
		{
			name:     "by length",
			slice:    []string{"a", "bb", "ccc"},
			value:    "dd",
			iteratee: func(s string) int { return len(s) },
			expected: 1, // length 2 should be inserted at index 1 (before "bb")
		},
		{
			name:     "by length - at beginning",
			slice:    []string{"bb", "ccc", "dddd"},
			value:    "a",
			iteratee: func(s string) int { return len(s) },
			expected: 0, // length 1 should be inserted at index 0
		},
		{
			name:     "by length - at end",
			slice:    []string{"a", "bb", "ccc"},
			value:    "eeee",
			iteratee: func(s string) int { return len(s) },
			expected: 3, // length 4 should be inserted at index 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedIndexBy(tt.slice, tt.value, tt.iteratee)
			if result != tt.expected {
				t.Errorf("SortedIndexBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		value    int
		expected int
	}{
		{
			name:     "found - first occurrence",
			slice:    []int{4, 5, 5, 5, 6},
			value:    5,
			expected: 1,
		},
		{
			name:     "found - single occurrence",
			slice:    []int{1, 2, 3, 4, 5},
			value:    3,
			expected: 2,
		},
		{
			name:     "not found",
			slice:    []int{4, 5, 5, 5, 6},
			value:    3,
			expected: -1,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			value:    5,
			expected: -1,
		},
		{
			name:     "found at beginning",
			slice:    []int{1, 2, 3},
			value:    1,
			expected: 0,
		},
		{
			name:     "found at end",
			slice:    []int{1, 2, 3},
			value:    3,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedIndexOf(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("SortedIndexOf() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedLastIndex(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		value    int
		expected int
	}{
		{
			name:     "duplicate values - last position",
			slice:    []int{4, 5, 5, 5, 6},
			value:    5,
			expected: 4,
		},
		{
			name:     "insert in middle",
			slice:    []int{30, 50},
			value:    40,
			expected: 1,
		},
		{
			name:     "insert at beginning",
			slice:    []int{30, 50},
			value:    20,
			expected: 0,
		},
		{
			name:     "insert at end",
			slice:    []int{30, 50},
			value:    60,
			expected: 2,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			value:    5,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedLastIndex(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("SortedLastIndex() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedLastIndexBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		value    string
		iteratee func(string) int
		expected int
	}{
		{
			name:     "by length - duplicate lengths",
			slice:    []string{"a", "b", "cc", "dd", "eee"},
			value:    "ff",
			iteratee: func(s string) int { return len(s) },
			expected: 4, // length 2 should be inserted after all length 2 elements
		},
		{
			name:     "by length - at end",
			slice:    []string{"a", "bb", "ccc"},
			value:    "dddd",
			iteratee: func(s string) int { return len(s) },
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedLastIndexBy(tt.slice, tt.value, tt.iteratee)
			if result != tt.expected {
				t.Errorf("SortedLastIndexBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedLastIndexOf(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		value    int
		expected int
	}{
		{
			name:     "found - last occurrence",
			slice:    []int{4, 5, 5, 5, 6},
			value:    5,
			expected: 3,
		},
		{
			name:     "found - single occurrence",
			slice:    []int{1, 2, 3, 4, 5},
			value:    3,
			expected: 2,
		},
		{
			name:     "not found",
			slice:    []int{4, 5, 5, 5, 6},
			value:    3,
			expected: -1,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			value:    5,
			expected: -1,
		},
		{
			name:     "found at beginning",
			slice:    []int{1, 2, 3},
			value:    1,
			expected: 0,
		},
		{
			name:     "found at end",
			slice:    []int{1, 2, 3},
			value:    3,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedLastIndexOf(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("SortedLastIndexOf() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedUniq(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected []int
	}{
		{
			name:     "basic sorted uniq",
			slice:    []int{1, 1, 2, 2, 3},
			expected: []int{1, 2, 3},
		},
		{
			name:     "no duplicates",
			slice:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "all duplicates",
			slice:    []int{1, 1, 1, 1},
			expected: []int{1},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			expected: []int{},
		},
		{
			name:     "single element",
			slice:    []int{5},
			expected: []int{5},
		},
		{
			name:     "consecutive duplicates",
			slice:    []int{1, 2, 2, 2, 3, 3, 4},
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedUniq(tt.slice)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortedUniq() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedUniqBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		iteratee func(int) int
		expected []int
	}{
		{
			name:     "identity function",
			slice:    []int{1, 1, 2, 2, 3},
			iteratee: func(x int) int { return x },
			expected: []int{1, 2, 3},
		},
		{
			name:     "modulo function",
			slice:    []int{3, 1, 4, 2, 5}, // sorted by x%3: [0,1,1,2,2]
			iteratee: func(x int) int { return x % 3 },
			expected: []int{3, 1, 2}, // first occurrence of each modulo value: 0, 1, 2
		},
		{
			name:  "absolute value function",
			slice: []int{-2, -1, 1, 2}, // sorted by abs: [-2,-1,1,2] -> [2,1,1,2]
			iteratee: func(x int) int {
				if x < 0 {
					return -x
				}
				return x
			},
			expected: []int{-2, -1, 2}, // first occurrence of each absolute value
		},
		{
			name:     "empty slice",
			slice:    []int{},
			iteratee: func(x int) int { return x },
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedUniqBy(tt.slice, tt.iteratee)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortedUniqBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSortedUniqByString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		iteratee func(string) string
		expected []string
	}{
		{
			name:     "case insensitive",
			slice:    []string{"a", "A", "b", "B", "c"},
			iteratee: func(s string) string { return strings.ToLower(s) },
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "by length",
			slice:    []string{"a", "b", "cc", "dd", "eee"},
			iteratee: func(s string) string { return fmt.Sprintf("%d", len(s)) },
			expected: []string{"a", "cc", "eee"}, // lengths: 1, 2, 3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortedUniqBy(tt.slice, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortedUniqBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestZipWith(t *testing.T) {
	tests := []struct {
		name     string
		iteratee func(...int) int
		slices   [][]int
		expected []int
	}{
		{
			name:     "sum two slices",
			iteratee: func(args ...int) int { return args[0] + args[1] },
			slices:   [][]int{{1, 2}, {3, 4}},
			expected: []int{4, 6},
		},
		{
			name:     "sum three slices",
			iteratee: func(args ...int) int { return args[0] + args[1] + args[2] },
			slices:   [][]int{{1, 2}, {3, 4}, {5, 6}},
			expected: []int{9, 12},
		},
		{
			name:     "different lengths - use minimum",
			iteratee: func(args ...int) int { return args[0] + args[1] },
			slices:   [][]int{{1, 2, 3}, {4, 5}},
			expected: []int{5, 7},
		},
		{
			name:     "empty slices",
			iteratee: func(args ...int) int { return args[0] + args[1] },
			slices:   [][]int{},
			expected: []int{},
		},
		{
			name:     "single slice",
			iteratee: func(args ...int) int { return args[0] * 2 },
			slices:   [][]int{{1, 2, 3}},
			expected: []int{2, 4, 6},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZipWith(tt.iteratee, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ZipWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestZipObject(t *testing.T) {
	tests := []struct {
		name     string
		keys     []string
		values   []int
		expected map[string]int
	}{
		{
			name:     "basic zip object",
			keys:     []string{"a", "b", "c"},
			values:   []int{1, 2, 3},
			expected: map[string]int{"a": 1, "b": 2, "c": 3},
		},
		{
			name:     "more keys than values",
			keys:     []string{"a", "b", "c"},
			values:   []int{1, 2},
			expected: map[string]int{"a": 1, "b": 2},
		},
		{
			name:     "more values than keys",
			keys:     []string{"a", "b"},
			values:   []int{1, 2, 3},
			expected: map[string]int{"a": 1, "b": 2},
		},
		{
			name:     "empty keys and values",
			keys:     []string{},
			values:   []int{},
			expected: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZipObject(tt.keys, tt.values)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ZipObject() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestZipObjectDeep(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		values   []interface{}
		expected map[string]interface{}
	}{
		{
			name:   "simple paths",
			paths:  []string{"a", "b"},
			values: []interface{}{1, 2},
			expected: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
		},
		{
			name:   "nested paths",
			paths:  []string{"a.b", "c.d"},
			values: []interface{}{1, 2},
			expected: map[string]interface{}{
				"a": map[string]interface{}{"b": 1},
				"c": map[string]interface{}{"d": 2},
			},
		},
		{
			name:   "deep nested paths",
			paths:  []string{"a.b.c", "d"},
			values: []interface{}{1, 2},
			expected: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{"c": 1},
				},
				"d": 2,
			},
		},
		{
			name:     "empty paths and values",
			paths:    []string{},
			values:   []interface{}{},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ZipObjectDeep(tt.paths, tt.values)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ZipObjectDeep() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "take while less than 4",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x < 4 },
			expected:  []int{1, 2, 3},
		},
		{
			name:      "take while even",
			slice:     []int{2, 4, 6, 1, 8},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "take all - predicate always true",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{1, 2, 3, 4, 5},
		},
		{
			name:      "take none - predicate false from start",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 10 },
			expected:  []int{},
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{},
		},
		{
			name:      "single element - true",
			slice:     []int{5},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{5},
		},
		{
			name:      "single element - false",
			slice:     []int{5},
			predicate: func(x int) bool { return x < 0 },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TakeWhile(tt.slice, tt.predicate)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TakeWhile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTakeWhileString(t *testing.T) {
	tests := []struct {
		name      string
		slice     []string
		predicate func(string) bool
		expected  []string
	}{
		{
			name:      "take while alphabetic",
			slice:     []string{"a", "b", "c", "1", "d"},
			predicate: func(s string) bool { return s >= "a" && s <= "z" },
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "take while short strings",
			slice:     []string{"a", "bb", "c", "dddd", "e"},
			predicate: func(s string) bool { return len(s) <= 2 },
			expected:  []string{"a", "bb", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TakeWhile(tt.slice, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TakeWhile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTakeRightWhile(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "take right while greater than 2",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 2 },
			expected:  []int{3, 4, 5},
		},
		{
			name:      "take right while even",
			slice:     []int{1, 3, 2, 4, 6},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "take all - predicate always true",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{1, 2, 3, 4, 5},
		},
		{
			name:      "take none - predicate false from end",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 10 },
			expected:  []int{},
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{},
		},
		{
			name:      "single element - true",
			slice:     []int{5},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{5},
		},
		{
			name:      "single element - false",
			slice:     []int{5},
			predicate: func(x int) bool { return x < 0 },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TakeRightWhile(tt.slice, tt.predicate)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TakeRightWhile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTakeRightWhileString(t *testing.T) {
	tests := []struct {
		name      string
		slice     []string
		predicate func(string) bool
		expected  []string
	}{
		{
			name:      "take right while numeric",
			slice:     []string{"a", "1", "2", "3"},
			predicate: func(s string) bool { return s >= "0" && s <= "9" },
			expected:  []string{"1", "2", "3"},
		},
		{
			name:      "take right while short strings",
			slice:     []string{"hello", "a", "bb", "c"},
			predicate: func(s string) bool { return len(s) <= 2 },
			expected:  []string{"a", "bb", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TakeRightWhile(tt.slice, tt.predicate)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("TakeRightWhile() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnionBy(t *testing.T) {
	tests := []struct {
		name     string
		iteratee func(int) int
		slices   [][]int
		expected []int
	}{
		{
			name:     "identity function",
			iteratee: func(x int) int { return x },
			slices:   [][]int{{2, 1}, {2, 3}},
			expected: []int{2, 1, 3},
		},
		{
			name: "absolute value function",
			iteratee: func(x int) int {
				if x < 0 {
					return -x
				}
				return x
			},
			slices:   [][]int{{-2, 1}, {2, -3, 4}},
			expected: []int{-2, 1, -3, 4}, // abs(-2)=2, abs(1)=1, abs(2)=2(dup), abs(-3)=3, abs(4)=4
		},
		{
			name:     "modulo function",
			iteratee: func(x int) int { return x % 3 },
			slices:   [][]int{{1, 4}, {2, 5}, {3, 6}},
			expected: []int{1, 2, 3}, // 1%3=1, 4%3=1(dup), 2%3=2, 5%3=2(dup), 3%3=0, 6%3=0(dup)
		},
		{
			name:     "empty slices",
			iteratee: func(x int) int { return x },
			slices:   [][]int{},
			expected: []int{},
		},
		{
			name:     "single slice",
			iteratee: func(x int) int { return x },
			slices:   [][]int{{1, 2, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "multiple slices with overlaps",
			iteratee: func(x int) int { return x },
			slices:   [][]int{{1, 2}, {2, 3}, {3, 4}},
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnionBy(tt.iteratee, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UnionBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnionByString(t *testing.T) {
	tests := []struct {
		name     string
		iteratee func(string) int
		slices   [][]string
		expected []string
	}{
		{
			name:     "by length",
			iteratee: func(s string) int { return len(s) },
			slices:   [][]string{{"a", "bb"}, {"ccc", "d"}},
			expected: []string{"a", "bb", "ccc"}, // lengths: 1, 2, 3, 1(dup)
		},
		{
			name: "by first character",
			iteratee: func(s string) int {
				if len(s) == 0 {
					return 0
				}
				return int(s[0])
			},
			slices:   [][]string{{"apple", "banana"}, {"avocado", "cherry"}},
			expected: []string{"apple", "banana", "cherry"}, // 'a', 'b', 'a'(dup), 'c'
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnionBy(tt.iteratee, tt.slices...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UnionBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnionWith(t *testing.T) {
	tests := []struct {
		name       string
		comparator func(int, int) bool
		slices     [][]int
		expected   []int
	}{
		{
			name:       "equality comparator",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{{2, 1}, {2, 3}},
			expected:   []int{2, 1, 3},
		},
		{
			name:       "same parity comparator",
			comparator: func(a, b int) bool { return a%2 == b%2 },
			slices:     [][]int{{1, 2}, {3, 4}},
			expected:   []int{1, 2}, // 1(odd), 2(even), 3(odd-dup), 4(even-dup)
		},
		{
			name:       "absolute value comparator",
			comparator: func(a, b int) bool { return abs(a) == abs(b) },
			slices:     [][]int{{-1, 2}, {1, -3}},
			expected:   []int{-1, 2, -3}, // abs(-1)=1, abs(2)=2, abs(1)=1(dup), abs(-3)=3
		},
		{
			name:       "empty slices",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{},
			expected:   []int{},
		},
		{
			name:       "single slice",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{{1, 2, 2, 3}},
			expected:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnionWith(tt.comparator, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UnionWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnionWithString(t *testing.T) {
	tests := []struct {
		name       string
		comparator func(string, string) bool
		slices     [][]string
		expected   []string
	}{
		{
			name:       "case insensitive",
			comparator: func(a, b string) bool { return strings.EqualFold(a, b) },
			slices:     [][]string{{"a", "B"}, {"A", "c"}},
			expected:   []string{"a", "B", "c"}, // "a", "B", "A"(dup), "c"
		},
		{
			name:       "by length",
			comparator: func(a, b string) bool { return len(a) == len(b) },
			slices:     [][]string{{"a", "bb"}, {"ccc", "d"}},
			expected:   []string{"a", "bb", "ccc"}, // lengths: 1, 2, 3, 1(dup)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnionWith(tt.comparator, tt.slices...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UnionWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnzip(t *testing.T) {
	tests := []struct {
		name     string
		slice    [][]int
		expected [][]int
	}{
		{
			name:     "basic unzip",
			slice:    [][]int{{1, 4}, {2, 5}, {3, 6}},
			expected: [][]int{{1, 2, 3}, {4, 5, 6}},
		},
		{
			name:     "unequal inner slice lengths",
			slice:    [][]int{{1, 4, 7}, {2, 5}, {3}},
			expected: [][]int{{1, 2, 3}, {4, 5}, {7}},
		},
		{
			name:     "single inner slice",
			slice:    [][]int{{1, 2, 3}},
			expected: [][]int{{1}, {2}, {3}},
		},
		{
			name:     "empty slice",
			slice:    [][]int{},
			expected: [][]int{},
		},
		{
			name:     "empty inner slices",
			slice:    [][]int{{}, {}, {}},
			expected: [][]int{},
		},
		{
			name:     "single element inner slices",
			slice:    [][]int{{1}, {2}, {3}},
			expected: [][]int{{1, 2, 3}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unzip(tt.slice)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Unzip() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnzipString(t *testing.T) {
	tests := []struct {
		name     string
		slice    [][]string
		expected [][]string
	}{
		{
			name:     "basic string unzip",
			slice:    [][]string{{"a", "d"}, {"b", "e"}, {"c", "f"}},
			expected: [][]string{{"a", "b", "c"}, {"d", "e", "f"}},
		},
		{
			name:     "mixed length strings",
			slice:    [][]string{{"hello", "world"}, {"foo"}, {"bar", "baz", "qux"}},
			expected: [][]string{{"hello", "foo", "bar"}, {"world", "baz"}, {"qux"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unzip(tt.slice)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Unzip() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnzipWith(t *testing.T) {
	tests := []struct {
		name     string
		slice    [][]int
		iteratee func(...int) int
		expected []int
	}{
		{
			name:  "sum function",
			slice: [][]int{{1, 4}, {2, 5}, {3, 6}},
			iteratee: func(args ...int) int {
				sum := 0
				for _, v := range args {
					sum += v
				}
				return sum
			},
			expected: []int{6, 15}, // [1+2+3, 4+5+6]
		},
		{
			name:  "max function",
			slice: [][]int{{1, 4}, {2, 5}, {3, 6}},
			iteratee: func(args ...int) int {
				if len(args) == 0 {
					return 0
				}
				max := args[0]
				for _, v := range args[1:] {
					if v > max {
						max = v
					}
				}
				return max
			},
			expected: []int{3, 6}, // [max(1,2,3), max(4,5,6)]
		},
		{
			name:  "unequal lengths",
			slice: [][]int{{1, 4, 7}, {2, 5}, {3}},
			iteratee: func(args ...int) int {
				sum := 0
				for _, v := range args {
					sum += v
				}
				return sum
			},
			expected: []int{6, 9, 7}, // [1+2+3, 4+5, 7]
		},
		{
			name:     "empty slice",
			slice:    [][]int{},
			iteratee: func(args ...int) int { return 0 },
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UnzipWith(tt.slice, tt.iteratee)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UnzipWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestXor(t *testing.T) {
	tests := []struct {
		name     string
		slices   [][]int
		expected []int
	}{
		{
			name:     "basic xor",
			slices:   [][]int{{2, 1}, {2, 3}},
			expected: []int{1, 3},
		},
		{
			name:     "three slices",
			slices:   [][]int{{1, 2}, {2, 3}, {3, 4}},
			expected: []int{1, 4},
		},
		{
			name:     "no common elements",
			slices:   [][]int{{1, 2}, {3, 4}},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "all common elements",
			slices:   [][]int{{1, 2}, {1, 2}},
			expected: []int{},
		},
		{
			name:     "single slice",
			slices:   [][]int{{1, 2, 3}},
			expected: []int{1, 2, 3},
		},
		{
			name:     "empty slices",
			slices:   [][]int{},
			expected: []int{},
		},
		{
			name:     "duplicates within slice",
			slices:   [][]int{{1, 1, 2}, {2, 3, 3}},
			expected: []int{1, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Xor(tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			// Sort both slices for comparison since order might vary
			sort.Ints(result)
			sort.Ints(tt.expected)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Xor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestXorBy(t *testing.T) {
	tests := []struct {
		name     string
		iteratee func(int) int
		slices   [][]int
		expected []int
	}{
		{
			name:     "identity function",
			iteratee: func(x int) int { return x },
			slices:   [][]int{{2, 1}, {2, 3}},
			expected: []int{1, 3},
		},
		{
			name:     "modulo function",
			iteratee: func(x int) int { return x % 3 },
			slices:   [][]int{{1, 4}, {2, 5}}, // 1%3=1, 4%3=1(same), 2%3=2, 5%3=2(same)
			expected: []int{1, 2},             // criteria 1 and 2 each appear in only one slice
		},
		{
			name: "absolute value function",
			iteratee: func(x int) int {
				if x < 0 {
					return -x
				}
				return x
			},
			slices:   [][]int{{-1, 2}, {1, -3}}, // abs(-1)=1, abs(2)=2, abs(1)=1(same), abs(-3)=3
			expected: []int{2, -3},              // abs(2)=2 and abs(-3)=3 are unique
		},
		{
			name:     "empty slices",
			iteratee: func(x int) int { return x },
			slices:   [][]int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := XorBy(tt.iteratee, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			// Sort both slices for comparison since order might vary
			sort.Ints(result)
			sort.Ints(tt.expected)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("XorBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestXorByString(t *testing.T) {
	tests := []struct {
		name     string
		iteratee func(string) int
		slices   [][]string
		expected []string
	}{
		{
			name:     "by length",
			iteratee: func(s string) int { return len(s) },
			slices:   [][]string{{"a", "bb"}, {"cc", "d"}}, // lengths: [1,2], [2,1] -> both have 1 and 2
			expected: []string{},                           // no unique lengths
		},
		{
			name:     "by length - unique",
			iteratee: func(s string) int { return len(s) },
			slices:   [][]string{{"a", "bb"}, {"ccc", "d"}}, // lengths: [1,2], [3,1] -> 2 and 3 are unique
			expected: []string{"bb", "ccc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := XorBy(tt.iteratee, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			// Sort both slices for comparison since order might vary
			sort.Strings(result)
			sort.Strings(tt.expected)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("XorBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestXorWith(t *testing.T) {
	tests := []struct {
		name       string
		comparator func(int, int) bool
		slices     [][]int
		expected   []int
	}{
		{
			name:       "equality comparator",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{{2, 1}, {2, 3}},
			expected:   []int{1, 3},
		},
		{
			name:       "same parity comparator",
			comparator: func(a, b int) bool { return a%2 == b%2 },
			slices:     [][]int{{1, 2}, {3, 4}}, // [odd,even], [odd,even] -> no unique parities
			expected:   []int{},
		},
		{
			name:       "absolute value comparator",
			comparator: func(a, b int) bool { return abs(a) == abs(b) },
			slices:     [][]int{{-1, 2}, {1, -3}}, // abs(-1)=1, abs(2)=2, abs(1)=1(same), abs(-3)=3
			expected:   []int{2, -3},              // abs(2)=2 and abs(-3)=3 are unique
		},
		{
			name:       "empty slices",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{},
			expected:   []int{},
		},
		{
			name:       "single slice",
			comparator: func(a, b int) bool { return a == b },
			slices:     [][]int{{1, 2, 3}},
			expected:   []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := XorWith(tt.comparator, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			// Sort both slices for comparison since order might vary
			sort.Ints(result)
			sort.Ints(tt.expected)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("XorWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestXorWithString(t *testing.T) {
	tests := []struct {
		name       string
		comparator func(string, string) bool
		slices     [][]string
		expected   []string
	}{
		{
			name:       "case insensitive",
			comparator: func(a, b string) bool { return strings.EqualFold(a, b) },
			slices:     [][]string{{"a", "B"}, {"A", "c"}}, // "a"/"A" appear in both, "B" and "c" are unique
			expected:   []string{"B", "c"},
		},
		{
			name:       "by length",
			comparator: func(a, b string) bool { return len(a) == len(b) },
			slices:     [][]string{{"a", "bb"}, {"ccc", "d"}}, // lengths: [1,2], [3,1] -> 2 and 3 are unique
			expected:   []string{"bb", "ccc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := XorWith(tt.comparator, tt.slices...)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			// Sort both slices for comparison since order might vary
			sort.Strings(result)
			sort.Strings(tt.expected)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("XorWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}
