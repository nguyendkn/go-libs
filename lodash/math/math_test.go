package math

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        6,
			b:        4,
			expected: 10,
		},
		{
			name:     "negative numbers",
			a:        -3,
			b:        -2,
			expected: -5,
		},
		{
			name:     "mixed signs",
			a:        5,
			b:        -3,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        6,
			b:        4,
			expected: 2,
		},
		{
			name:     "negative result",
			a:        3,
			b:        5,
			expected: -2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{
			name:     "positive numbers",
			a:        6,
			b:        4,
			expected: 24,
		},
		{
			name:     "with zero",
			a:        5,
			b:        0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Multiply(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Multiply() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{
			name:     "basic division",
			a:        6,
			b:        4,
			expected: 1.5,
		},
		{
			name:     "exact division",
			a:        10,
			b:        2,
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Divide(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Divide() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name        string
		numbers     []int
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "basic max",
			numbers:     []int{4, 2, 8, 6},
			expectedVal: 8,
			expectedOk:  true,
		},
		{
			name:        "negative numbers",
			numbers:     []int{-1, -5, -3},
			expectedVal: -1,
			expectedOk:  true,
		},
		{
			name:        "empty slice",
			numbers:     []int{},
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := Max(tt.numbers)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("Max() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name        string
		numbers     []int
		expectedVal int
		expectedOk  bool
	}{
		{
			name:        "basic min",
			numbers:     []int{4, 2, 8, 6},
			expectedVal: 2,
			expectedOk:  true,
		},
		{
			name:        "negative numbers",
			numbers:     []int{-1, -5, -3},
			expectedVal: -5,
			expectedOk:  true,
		},
		{
			name:        "empty slice",
			numbers:     []int{},
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := Min(tt.numbers)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("Min() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{
			name:     "basic sum",
			numbers:  []int{4, 2, 8, 6},
			expected: 20,
		},
		{
			name:     "empty slice",
			numbers:  []int{},
			expected: 0,
		},
		{
			name:     "negative numbers",
			numbers:  []int{-1, -2, -3},
			expected: -6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sum(tt.numbers)
			if result != tt.expected {
				t.Errorf("Sum() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMean(t *testing.T) {
	tests := []struct {
		name        string
		numbers     []int
		expectedVal float64
		expectedOk  bool
	}{
		{
			name:        "basic mean",
			numbers:     []int{4, 2, 8, 6},
			expectedVal: 5.0,
			expectedOk:  true,
		},
		{
			name:        "empty slice",
			numbers:     []int{},
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := Mean(tt.numbers)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("Mean() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestCeil(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "positive decimal",
			input:    4.006,
			expected: 5.0,
		},
		{
			name:     "negative decimal",
			input:    -4.006,
			expected: -4.0,
		},
		{
			name:     "integer",
			input:    6,
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Ceil(tt.input)
			if result != tt.expected {
				t.Errorf("Ceil() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFloor(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "positive decimal",
			input:    4.006,
			expected: 4.0,
		},
		{
			name:     "negative decimal",
			input:    -4.006,
			expected: -5.0,
		},
		{
			name:     "integer",
			input:    4,
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Floor(tt.input)
			if result != tt.expected {
				t.Errorf("Floor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "round down",
			input:    4.006,
			expected: 4.0,
		},
		{
			name:     "round up",
			input:    4.6,
			expected: 5.0,
		},
		{
			name:     "negative round up",
			input:    -4.6,
			expected: -5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Round(tt.input)
			if result != tt.expected {
				t.Errorf("Round() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "negative number",
			input:    -5,
			expected: 5,
		},
		{
			name:     "positive number",
			input:    5,
			expected: 5,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Abs(tt.input)
			if result != tt.expected {
				t.Errorf("Abs() = %v, want %v", result, tt.expected)
			}
		})
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
			name:     "below range",
			number:   -10,
			lower:    -5,
			upper:    5,
			expected: -5,
		},
		{
			name:     "above range",
			number:   10,
			lower:    -5,
			upper:    5,
			expected: 5,
		},
		{
			name:     "within range",
			number:   3,
			lower:    -5,
			upper:    5,
			expected: 3,
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

func TestMaxBy(t *testing.T) {
	tests := []struct {
		name        string
		slice       []string
		iteratee    func(string) int
		expectedVal string
		expectedOk  bool
	}{
		{
			name:        "max by length",
			slice:       []string{"a", "bb", "ccc"},
			iteratee:    func(s string) int { return len(s) },
			expectedVal: "ccc",
			expectedOk:  true,
		},
		{
			name:        "empty slice",
			slice:       []string{},
			iteratee:    func(s string) int { return len(s) },
			expectedVal: "",
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := MaxBy(tt.slice, tt.iteratee)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("MaxBy() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestMinBy(t *testing.T) {
	tests := []struct {
		name        string
		slice       []string
		iteratee    func(string) int
		expectedVal string
		expectedOk  bool
	}{
		{
			name:        "min by length",
			slice:       []string{"a", "bb", "ccc"},
			iteratee:    func(s string) int { return len(s) },
			expectedVal: "a",
			expectedOk:  true,
		},
		{
			name:        "empty slice",
			slice:       []string{},
			iteratee:    func(s string) int { return len(s) },
			expectedVal: "",
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := MinBy(tt.slice, tt.iteratee)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("MinBy() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestSumBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		iteratee func(string) int
		expected int
	}{
		{
			name:     "sum by length",
			slice:    []string{"a", "bb", "ccc"},
			iteratee: func(s string) int { return len(s) },
			expected: 6,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			iteratee: func(s string) int { return len(s) },
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SumBy(tt.slice, tt.iteratee)
			if result != tt.expected {
				t.Errorf("SumBy() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMeanBy(t *testing.T) {
	tests := []struct {
		name        string
		slice       []string
		iteratee    func(string) int
		expectedVal float64
		expectedOk  bool
	}{
		{
			name:        "mean by length",
			slice:       []string{"a", "bb", "ccc"},
			iteratee:    func(s string) int { return len(s) },
			expectedVal: 2.0,
			expectedOk:  true,
		},
		{
			name:        "empty slice",
			slice:       []string{},
			iteratee:    func(s string) int { return len(s) },
			expectedVal: 0,
			expectedOk:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := MeanBy(tt.slice, tt.iteratee)
			if val != tt.expectedVal || ok != tt.expectedOk {
				t.Errorf("MeanBy() = (%v, %v), want (%v, %v)", val, ok, tt.expectedVal, tt.expectedOk)
			}
		})
	}
}

func TestPow(t *testing.T) {
	tests := []struct {
		name     string
		base     float64
		exponent float64
		expected float64
	}{
		{
			name:     "basic power",
			base:     2,
			exponent: 3,
			expected: 8,
		},
		{
			name:     "square root",
			base:     4,
			exponent: 0.5,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pow(tt.base, tt.exponent)
			if result != tt.expected {
				t.Errorf("Pow() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSqrt(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{
			name:     "perfect square",
			input:    9,
			expected: 3,
		},
		{
			name:     "non-perfect square",
			input:    2,
			expected: math.Sqrt(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sqrt(tt.input)
			if result != tt.expected {
				t.Errorf("Sqrt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsNaN(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected bool
	}{
		{
			name:     "NaN",
			input:    math.NaN(),
			expected: true,
		},
		{
			name:     "normal number",
			input:    1.0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNaN(tt.input)
			if result != tt.expected {
				t.Errorf("IsNaN() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsInf(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		sign     int
		expected bool
	}{
		{
			name:     "positive infinity",
			input:    math.Inf(1),
			sign:     0,
			expected: true,
		},
		{
			name:     "normal number",
			input:    1.0,
			sign:     0,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInf(tt.input, tt.sign)
			if result != tt.expected {
				t.Errorf("IsInf() = %v, want %v", result, tt.expected)
			}
		})
	}
}
