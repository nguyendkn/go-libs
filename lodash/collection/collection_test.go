package collection

import (
	"reflect"
	"strconv"
	"testing"
)

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "filter even numbers",
			slice:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{2, 4, 6},
		},
		{
			name:      "filter greater than 3",
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(x int) bool { return x > 3 },
			expected:  []int{4, 5},
		},
		{
			name:      "no matches",
			slice:     []int{1, 3, 5},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{},
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.slice, tt.predicate)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Filter() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMap(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		mapper   func(int) string
		expected []string
	}{
		{
			name:     "convert to string",
			slice:    []int{1, 2, 3},
			mapper:   func(x int) string { return strconv.Itoa(x) },
			expected: []string{"1", "2", "3"},
		},
		{
			name:     "double and convert",
			slice:    []int{1, 2, 3},
			mapper:   func(x int) string { return strconv.Itoa(x * 2) },
			expected: []string{"2", "4", "6"},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			mapper:   func(x int) string { return strconv.Itoa(x) },
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Map(tt.slice, tt.mapper)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Map() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReduce(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		reducer  func(int, int) int
		initial  int
		expected int
	}{
		{
			name:     "sum",
			slice:    []int{1, 2, 3, 4},
			reducer:  func(acc, x int) int { return acc + x },
			initial:  0,
			expected: 10,
		},
		{
			name:     "product",
			slice:    []int{2, 3, 4},
			reducer:  func(acc, x int) int { return acc * x },
			initial:  1,
			expected: 24,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			reducer:  func(acc, x int) int { return acc + x },
			initial:  5,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reduce(tt.slice, tt.reducer, tt.initial)
			if result != tt.expected {
				t.Errorf("Reduce() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name        string
		slice       []int
		predicate   func(int) bool
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "find first even",
			slice:       []int{1, 3, 4, 6},
			predicate:   func(x int) bool { return x%2 == 0 },
			expectedVal: 4,
			expectedOk:  true,
		},
		{
			name:        "not found",
			slice:       []int{1, 3, 5},
			predicate:   func(x int) bool { return x%2 == 0 },
			expectedVal: 0,
			expectedOk:  false,
		},
		{
			name:        "empty slice",
			slice:       []int{},
			predicate:   func(x int) bool { return x > 0 },
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := Find(tt.slice, tt.predicate)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("Find() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestFindIndex(t *testing.T) {
	tests := []struct {
		name        string
		slice       []int
		predicate   func(int) bool
		expectedIdx int
		expectedOk  bool
	}{
		{
			name:        "find first even index",
			slice:       []int{1, 3, 4, 6},
			predicate:   func(x int) bool { return x%2 == 0 },
			expectedIdx: 2,
			expectedOk:  true,
		},
		{
			name:        "not found",
			slice:       []int{1, 3, 5},
			predicate:   func(x int) bool { return x%2 == 0 },
			expectedIdx: -1,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx, ok := FindIndex(tt.slice, tt.predicate)
			if idx != tt.expectedIdx || ok != tt.expectedOk {
				t.Errorf("FindIndex() = (%v, %v), want (%v, %v)", idx, ok, tt.expectedIdx, tt.expectedOk)
			}
		})
	}
}

func TestEvery(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "all even",
			slice:     []int{2, 4, 6},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  true,
		},
		{
			name:      "not all even",
			slice:     []int{1, 2, 4},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  false,
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Every(tt.slice, tt.predicate)
			if result != tt.expected {
				t.Errorf("Every() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSome(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		expected  bool
	}{
		{
			name:      "some even",
			slice:     []int{1, 2, 3},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  true,
		},
		{
			name:      "none even",
			slice:     []int{1, 3, 5},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  false,
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Some(tt.slice, tt.predicate)
			if result != tt.expected {
				t.Errorf("Some() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGroupBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		keyFunc  func(string) int
		expected map[int][]string
	}{
		{
			name:    "group by length",
			slice:   []string{"one", "two", "three", "four"},
			keyFunc: func(s string) int { return len(s) },
			expected: map[int][]string{
				3: {"one", "two"},
				4: {"four"},
				5: {"three"},
			},
		},
		{
			name:     "empty slice",
			slice:    []string{},
			keyFunc:  func(s string) int { return len(s) },
			expected: map[int][]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GroupBy(tt.slice, tt.keyFunc)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("GroupBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCountBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		keyFunc  func(string) int
		expected map[int]int
	}{
		{
			name:    "count by length",
			slice:   []string{"one", "two", "three", "four"},
			keyFunc: func(s string) int { return len(s) },
			expected: map[int]int{
				3: 2,
				4: 1,
				5: 1,
			},
		},
		{
			name:     "empty slice",
			slice:    []string{},
			keyFunc:  func(s string) int { return len(s) },
			expected: map[int]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountBy(tt.slice, tt.keyFunc)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CountBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	tests := []struct {
		name           string
		slice          []int
		predicate      func(int) bool
		expectedTruthy []int
		expectedFalsy  []int
	}{
		{
			name:           "partition even/odd",
			slice:          []int{1, 2, 3, 4, 5, 6},
			predicate:      func(x int) bool { return x%2 == 0 },
			expectedTruthy: []int{2, 4, 6},
			expectedFalsy:  []int{1, 3, 5},
		},
		{
			name:           "all truthy",
			slice:          []int{2, 4, 6},
			predicate:      func(x int) bool { return x%2 == 0 },
			expectedTruthy: []int{2, 4, 6},
			expectedFalsy:  []int{},
		},
		{
			name:           "empty slice",
			slice:          []int{},
			predicate:      func(x int) bool { return x%2 == 0 },
			expectedTruthy: []int{},
			expectedFalsy:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			truthy, falsy := Partition(tt.slice, tt.predicate)
			// Handle empty slice comparison for truthy
			if len(truthy) == 0 && len(tt.expectedTruthy) == 0 {
				// Both truthy slices are empty, continue to check falsy
			} else if !reflect.DeepEqual(truthy, tt.expectedTruthy) {
				t.Errorf("Partition() truthy = %v, want %v", truthy, tt.expectedTruthy)
			}
			// Handle empty slice comparison for falsy
			if len(falsy) == 0 && len(tt.expectedFalsy) == 0 {
				// Both falsy slices are empty, test passes
			} else if !reflect.DeepEqual(falsy, tt.expectedFalsy) {
				t.Errorf("Partition() falsy = %v, want %v", falsy, tt.expectedFalsy)
			}
		})
	}
}

func TestSize(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected int
	}{
		{
			name:     "non-empty slice",
			slice:    []int{1, 2, 3},
			expected: 3,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Size(tt.slice)
			if result != tt.expected {
				t.Errorf("Size() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		value    int
		expected bool
	}{
		{
			name:     "value exists",
			slice:    []int{1, 2, 3},
			value:    2,
			expected: true,
		},
		{
			name:     "value does not exist",
			slice:    []int{1, 2, 3},
			value:    4,
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []int{},
			value:    1,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Includes(tt.slice, tt.value)
			if result != tt.expected {
				t.Errorf("Includes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReject(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		expected  []int
	}{
		{
			name:      "reject even numbers",
			slice:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{1, 3, 5},
		},
		{
			name:      "reject none",
			slice:     []int{1, 3, 5},
			predicate: func(x int) bool { return x%2 == 0 },
			expected:  []int{1, 3, 5},
		},
		{
			name:      "empty slice",
			slice:     []int{},
			predicate: func(x int) bool { return x > 0 },
			expected:  []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Reject(tt.slice, tt.predicate)
			// Handle empty slice comparison
			if len(result) == 0 && len(tt.expected) == 0 {
				return // Both are empty, test passes
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Reject() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFindLast(t *testing.T) {
	tests := []struct {
		name        string
		slice       []int
		predicate   func(int) bool
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "find last even",
			slice:       []int{1, 2, 3, 4, 5, 6},
			predicate:   func(x int) bool { return x%2 == 0 },
			expectedVal: 6,
			expectedOk:  true,
		},
		{
			name:        "not found",
			slice:       []int{1, 3, 5},
			predicate:   func(x int) bool { return x%2 == 0 },
			expectedVal: 0,
			expectedOk:  false,
		},
		{
			name:        "empty slice",
			slice:       []int{},
			predicate:   func(x int) bool { return x > 0 },
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := FindLast(tt.slice, tt.predicate)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("FindLast() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}
