package array

import (
	"reflect"
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
