package util

import (
	"fmt"
	"reflect"
	"testing"
)

func TestIdentity(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{
			name:  "int",
			value: 42,
		},
		{
			name:  "string",
			value: "hello",
		},
		{
			name:  "slice",
			value: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Identity(tt.value)
			if !reflect.DeepEqual(result, tt.value) {
				t.Errorf("Identity() = %v, want %v", result, tt.value)
			}
		})
	}
}

func TestConstant(t *testing.T) {
	getValue := Constant(42)

	result1 := getValue()
	result2 := getValue()

	if result1 != 42 || result2 != 42 {
		t.Errorf("Constant() function should always return 42, got %v and %v", result1, result2)
	}
}

func TestNoop(t *testing.T) {
	// Just test that it doesn't panic
	Noop()
}

func TestRange(t *testing.T) {
	tests := []struct {
		name     string
		args     []int
		expected []int
	}{
		{
			name:     "single argument",
			args:     []int{4},
			expected: []int{0, 1, 2, 3},
		},
		{
			name:     "start and end",
			args:     []int{1, 5},
			expected: []int{1, 2, 3, 4},
		},
		{
			name:     "start, end, and step",
			args:     []int{0, 20, 5},
			expected: []int{0, 5, 10, 15},
		},
		{
			name:     "negative step",
			args:     []int{10, 0, -2},
			expected: []int{10, 8, 6, 4, 2},
		},
		{
			name:     "zero step",
			args:     []int{0, 10, 0},
			expected: []int{},
		},
		{
			name:     "empty args",
			args:     []int{},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Range(tt.args...)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Range() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTimes(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		iteratee func(int) int
		expected []int
	}{
		{
			name:     "multiply by 2",
			n:        3,
			iteratee: func(i int) int { return i * 2 },
			expected: []int{0, 2, 4},
		},
		{
			name:     "zero times",
			n:        0,
			iteratee: func(i int) int { return i },
			expected: []int{},
		},
		{
			name:     "negative times",
			n:        -1,
			iteratee: func(i int) int { return i },
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Times(tt.n, tt.iteratee)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Times() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniqueId(t *testing.T) {
	// Test without prefix
	id1 := UniqueId()
	id2 := UniqueId()

	if id1 == id2 {
		t.Errorf("UniqueId() should generate unique IDs, got %s and %s", id1, id2)
	}

	// Test with prefix
	prefixedId := UniqueId("test_")
	if len(prefixedId) <= 5 { // "test_" + at least one digit
		t.Errorf("UniqueId() with prefix should be longer than prefix, got %s", prefixedId)
	}
}

func TestDefaultTo(t *testing.T) {
	tests := []struct {
		name         string
		value        int
		defaultValue int
		expected     int
	}{
		{
			name:         "non-zero value",
			value:        42,
			defaultValue: 10,
			expected:     42,
		},
		{
			name:         "zero value",
			value:        0,
			defaultValue: 10,
			expected:     10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultTo(tt.value, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("DefaultTo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDefaultToAny(t *testing.T) {
	tests := []struct {
		name         string
		value        interface{}
		defaultValue interface{}
		expected     interface{}
	}{
		{
			name:         "non-nil value",
			value:        42,
			defaultValue: "default",
			expected:     42,
		},
		{
			name:         "nil value",
			value:        nil,
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DefaultToAny(tt.value, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("DefaultToAny() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAttempt(t *testing.T) {
	// Test successful function
	result, err := Attempt(func() (int, error) {
		return 42, nil
	})

	if err != nil {
		t.Errorf("Attempt() should not return error for successful function, got %v", err)
	}
	if result != 42 {
		t.Errorf("Attempt() = %v, want 42", result)
	}

	// Test function with error
	_, err = Attempt(func() (int, error) {
		return 0, fmt.Errorf("test error")
	})

	if err == nil {
		t.Error("Attempt() should return error when function fails")
	}
}

func TestFlow(t *testing.T) {
	add := func(x int) int { return x + 1 }
	multiply := func(x int) int { return x * 2 }

	composed := Flow(add, multiply)
	result := composed(3) // (3 + 1) * 2 = 8

	if result != 8 {
		t.Errorf("Flow() = %v, want 8", result)
	}
}

func TestFlowRight(t *testing.T) {
	add := func(x int) int { return x + 1 }
	multiply := func(x int) int { return x * 2 }

	composed := FlowRight(add, multiply)
	result := composed(3) // (3 * 2) + 1 = 7

	if result != 7 {
		t.Errorf("FlowRight() = %v, want 7", result)
	}
}

func TestStubFunctions(t *testing.T) {
	if !reflect.DeepEqual(StubArray(), []interface{}{}) {
		t.Error("StubArray() should return empty slice")
	}

	if StubFalse() != false {
		t.Error("StubFalse() should return false")
	}

	if !reflect.DeepEqual(StubObject(), map[string]interface{}{}) {
		t.Error("StubObject() should return empty map")
	}

	if StubString() != "" {
		t.Error("StubString() should return empty string")
	}

	if StubTrue() != true {
		t.Error("StubTrue() should return true")
	}
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		lower    int
		upper    int
		expected int
	}{
		{
			name:     "within bounds",
			number:   10,
			lower:    5,
			upper:    15,
			expected: 10,
		},
		{
			name:     "below lower bound",
			number:   3,
			lower:    5,
			upper:    15,
			expected: 5,
		},
		{
			name:     "above upper bound",
			number:   20,
			lower:    5,
			upper:    15,
			expected: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Clamp(tt.number, tt.lower, tt.upper)
			if result != tt.expected {
				t.Errorf("Clamp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInRange(t *testing.T) {
	tests := []struct {
		name     string
		number   int
		args     []int
		expected bool
	}{
		{
			name:     "in range with start and end",
			number:   3,
			args:     []int{2, 4},
			expected: true,
		},
		{
			name:     "in range with only end",
			number:   4,
			args:     []int{8},
			expected: true,
		},
		{
			name:     "not in range",
			number:   4,
			args:     []int{2},
			expected: false,
		},
		{
			name:     "equal to end",
			number:   2,
			args:     []int{2},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InRange(tt.number, tt.args...)
			if result != tt.expected {
				t.Errorf("InRange() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToPath(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected []string
	}{
		{
			name:     "dot notation",
			str:      "a.b.c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "bracket notation",
			str:      "a[0].b",
			expected: []string{"a", "0", "b"},
		},
		{
			name:     "mixed notation",
			str:      "a.b[1].c",
			expected: []string{"a", "b", "1", "c"},
		},
		{
			name:     "empty string",
			str:      "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPath(tt.str)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ToPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}
