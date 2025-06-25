package date

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	before := time.Now().UnixMilli()
	now := Now()
	after := time.Now().UnixMilli()
	
	if now < before || now > after {
		t.Errorf("Now() = %v, should be between %v and %v", now, before, after)
	}
}

func TestToDate(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
	}{
		{
			name:      "time.Time",
			value:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			expectErr: false,
		},
		{
			name:      "int64 timestamp",
			value:     int64(1640995200000),
			expectErr: false,
		},
		{
			name:      "int timestamp",
			value:     1640995200000,
			expectErr: false,
		},
		{
			name:      "RFC3339 string",
			value:     "2022-01-01T00:00:00Z",
			expectErr: false,
		},
		{
			name:      "date string",
			value:     "2022-01-01",
			expectErr: false,
		},
		{
			name:      "invalid string",
			value:     "invalid-date",
			expectErr: true,
		},
		{
			name:      "unsupported type",
			value:     []int{1, 2, 3},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ToDate(tt.value)
			if tt.expectErr {
				if err == nil {
					t.Errorf("ToDate() should return error for %v", tt.value)
				}
			} else {
				if err != nil {
					t.Errorf("ToDate() should not return error for %v, got %v", tt.value, err)
				}
				if result.IsZero() {
					t.Errorf("ToDate() should not return zero time for valid input %v", tt.value)
				}
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
			value:    "2022-01-01",
			expected: false,
		},
		{
			name:     "int",
			value:    1640995200000,
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

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected bool
	}{
		{
			name:     "valid time",
			time:     time.Now(),
			expected: true,
		},
		{
			name:     "zero time",
			time:     time.Time{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValid(tt.time)
			if result != tt.expected {
				t.Errorf("IsValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormat(t *testing.T) {
	testTime := time.Date(2022, 1, 1, 12, 30, 45, 0, time.UTC)
	
	tests := []struct {
		name     string
		layout   string
		expected string
	}{
		{
			name:     "date format",
			layout:   "2006-01-02",
			expected: "2022-01-01",
		},
		{
			name:     "time format",
			layout:   "15:04:05",
			expected: "12:30:45",
		},
		{
			name:     "RFC3339 format",
			layout:   time.RFC3339,
			expected: "2022-01-01T12:30:45Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Format(testTime, tt.layout)
			if result != tt.expected {
				t.Errorf("Format() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	baseTime := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	
	result := Add(baseTime, time.Hour)
	expected := time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC)
	
	if !result.Equal(expected) {
		t.Errorf("Add() = %v, want %v", result, expected)
	}
}

func TestSub(t *testing.T) {
	t1 := time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC)
	t2 := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	
	result := Sub(t1, t2)
	expected := time.Hour
	
	if result != expected {
		t.Errorf("Sub() = %v, want %v", result, expected)
	}
}

func TestBefore(t *testing.T) {
	t1 := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC)
	
	if !Before(t1, t2) {
		t.Error("Before() should return true when first time is before second")
	}
	
	if Before(t2, t1) {
		t.Error("Before() should return false when first time is after second")
	}
}

func TestAfter(t *testing.T) {
	t1 := time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC)
	t2 := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	
	if !After(t1, t2) {
		t.Error("After() should return true when first time is after second")
	}
	
	if After(t2, t1) {
		t.Error("After() should return false when first time is before second")
	}
}

func TestEqual(t *testing.T) {
	t1 := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	t2 := time.Date(2022, 1, 1, 12, 0, 0, 0, time.UTC)
	t3 := time.Date(2022, 1, 1, 13, 0, 0, 0, time.UTC)
	
	if !Equal(t1, t2) {
		t.Error("Equal() should return true for equal times")
	}
	
	if Equal(t1, t3) {
		t.Error("Equal() should return false for different times")
	}
}

func TestStartOfDay(t *testing.T) {
	input := time.Date(2022, 1, 1, 15, 30, 45, 123456789, time.UTC)
	result := StartOfDay(input)
	expected := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	
	if !result.Equal(expected) {
		t.Errorf("StartOfDay() = %v, want %v", result, expected)
	}
}

func TestEndOfDay(t *testing.T) {
	input := time.Date(2022, 1, 1, 15, 30, 45, 123456789, time.UTC)
	result := EndOfDay(input)
	expected := time.Date(2022, 1, 1, 23, 59, 59, 999999999, time.UTC)
	
	if !result.Equal(expected) {
		t.Errorf("EndOfDay() = %v, want %v", result, expected)
	}
}

func TestStartOfMonth(t *testing.T) {
	input := time.Date(2022, 6, 15, 15, 30, 45, 0, time.UTC)
	result := StartOfMonth(input)
	expected := time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC)
	
	if !result.Equal(expected) {
		t.Errorf("StartOfMonth() = %v, want %v", result, expected)
	}
}

func TestEndOfMonth(t *testing.T) {
	input := time.Date(2022, 6, 15, 15, 30, 45, 0, time.UTC)
	result := EndOfMonth(input)
	expected := time.Date(2022, 6, 30, 23, 59, 59, 999999999, time.UTC)
	
	if !result.Equal(expected) {
		t.Errorf("EndOfMonth() = %v, want %v", result, expected)
	}
}

func TestStartOfYear(t *testing.T) {
	input := time.Date(2022, 6, 15, 15, 30, 45, 0, time.UTC)
	result := StartOfYear(input)
	expected := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	
	if !result.Equal(expected) {
		t.Errorf("StartOfYear() = %v, want %v", result, expected)
	}
}

func TestEndOfYear(t *testing.T) {
	input := time.Date(2022, 6, 15, 15, 30, 45, 0, time.UTC)
	result := EndOfYear(input)
	expected := time.Date(2022, 12, 31, 23, 59, 59, 999999999, time.UTC)
	
	if !result.Equal(expected) {
		t.Errorf("EndOfYear() = %v, want %v", result, expected)
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected int
	}{
		{
			name:     "January",
			time:     time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		{
			name:     "February non-leap year",
			time:     time.Date(2022, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: 28,
		},
		{
			name:     "February leap year",
			time:     time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: 29,
		},
		{
			name:     "April",
			time:     time.Date(2022, 4, 15, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DaysInMonth(tt.time)
			if result != tt.expected {
				t.Errorf("DaysInMonth() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		expected bool
	}{
		{
			name:     "leap year divisible by 4",
			year:     2020,
			expected: true,
		},
		{
			name:     "non-leap year",
			year:     2021,
			expected: false,
		},
		{
			name:     "century non-leap year",
			year:     1900,
			expected: false,
		},
		{
			name:     "century leap year",
			year:     2000,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsLeapYear(tt.year)
			if result != tt.expected {
				t.Errorf("IsLeapYear() = %v, want %v", result, tt.expected)
			}
		})
	}
}
