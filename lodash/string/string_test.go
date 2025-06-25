package string

import (
	"reflect"
	"testing"
)

func TestCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic camel case",
			input:    "foo bar",
			expected: "fooBar",
		},
		{
			name:     "with dashes",
			input:    "--foo-bar--",
			expected: "fooBar",
		},
		{
			name:     "with underscores",
			input:    "__FOO_BAR__",
			expected: "fooBar",
		},
		{
			name:     "already camel case",
			input:    "fooBar",
			expected: "fooBar",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CamelCase(tt.input)
			if result != tt.expected {
				t.Errorf("CamelCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestKebabCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic kebab case",
			input:    "Foo Bar",
			expected: "foo-bar",
		},
		{
			name:     "from camel case",
			input:    "fooBar",
			expected: "foo-bar",
		},
		{
			name:     "with underscores",
			input:    "__FOO_BAR__",
			expected: "foo-bar",
		},
		{
			name:     "already kebab case",
			input:    "foo-bar",
			expected: "foo-bar",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KebabCase(tt.input)
			if result != tt.expected {
				t.Errorf("KebabCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic snake case",
			input:    "Foo Bar",
			expected: "foo_bar",
		},
		{
			name:     "from camel case",
			input:    "fooBar",
			expected: "foo_bar",
		},
		{
			name:     "with dashes",
			input:    "--FOO-BAR--",
			expected: "foo_bar",
		},
		{
			name:     "already snake case",
			input:    "foo_bar",
			expected: "foo_bar",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SnakeCase(tt.input)
			if result != tt.expected {
				t.Errorf("SnakeCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic pascal case",
			input:    "foo bar",
			expected: "FooBar",
		},
		{
			name:     "with dashes",
			input:    "--foo-bar--",
			expected: "FooBar",
		},
		{
			name:     "with underscores",
			input:    "__FOO_BAR__",
			expected: "FooBar",
		},
		{
			name:     "already pascal case",
			input:    "FooBar",
			expected: "FooBar",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PascalCase(tt.input)
			if result != tt.expected {
				t.Errorf("PascalCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "all caps",
			input:    "FRED",
			expected: "Fred",
		},
		{
			name:     "mixed case",
			input:    "fRED",
			expected: "Fred",
		},
		{
			name:     "lowercase",
			input:    "fred",
			expected: "Fred",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Capitalize(tt.input)
			if result != tt.expected {
				t.Errorf("Capitalize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLowerFirst(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "capitalize first",
			input:    "Fred",
			expected: "fred",
		},
		{
			name:     "all caps",
			input:    "FRED",
			expected: "fRED",
		},
		{
			name:     "already lowercase",
			input:    "fred",
			expected: "fred",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LowerFirst(tt.input)
			if result != tt.expected {
				t.Errorf("LowerFirst() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUpperFirst(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase first",
			input:    "fred",
			expected: "Fred",
		},
		{
			name:     "already uppercase",
			input:    "FRED",
			expected: "FRED",
		},
		{
			name:     "mixed case",
			input:    "fRED",
			expected: "FRED",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpperFirst(tt.input)
			if result != tt.expected {
				t.Errorf("UpperFirst() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "spaces",
			input:    "  abc  ",
			expected: "abc",
		},
		{
			name:     "tabs and newlines",
			input:    "\t\nabc\t\n",
			expected: "abc",
		},
		{
			name:     "no whitespace",
			input:    "abc",
			expected: "abc",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Trim(tt.input)
			if result != tt.expected {
				t.Errorf("Trim() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPad(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		chars    string
		expected string
	}{
		{
			name:     "basic padding",
			input:    "abc",
			length:   8,
			chars:    "_-",
			expected: "_-abc_-_",
		},
		{
			name:     "single char padding",
			input:    "abc",
			length:   6,
			chars:    "_",
			expected: "_abc__",
		},
		{
			name:     "no padding needed",
			input:    "abc",
			length:   3,
			chars:    "_",
			expected: "abc",
		},
		{
			name:     "shorter than input",
			input:    "abcdef",
			length:   3,
			chars:    "_",
			expected: "abcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Pad(tt.input, tt.length, tt.chars)
			if result != tt.expected {
				t.Errorf("Pad() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{
			name:     "repeat asterisk",
			input:    "*",
			n:        3,
			expected: "***",
		},
		{
			name:     "repeat string",
			input:    "abc",
			n:        2,
			expected: "abcabc",
		},
		{
			name:     "zero repetitions",
			input:    "abc",
			n:        0,
			expected: "",
		},
		{
			name:     "negative repetitions",
			input:    "abc",
			n:        -1,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Repeat(tt.input, tt.n)
			if result != tt.expected {
				t.Errorf("Repeat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStartsWith(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		target   string
		expected bool
	}{
		{
			name:     "starts with",
			input:    "abc",
			target:   "a",
			expected: true,
		},
		{
			name:     "does not start with",
			input:    "abc",
			target:   "b",
			expected: false,
		},
		{
			name:     "empty target",
			input:    "abc",
			target:   "",
			expected: true,
		},
		{
			name:     "empty input",
			input:    "",
			target:   "a",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StartsWith(tt.input, tt.target)
			if result != tt.expected {
				t.Errorf("StartsWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		target   string
		expected bool
	}{
		{
			name:     "ends with",
			input:    "abc",
			target:   "c",
			expected: true,
		},
		{
			name:     "does not end with",
			input:    "abc",
			target:   "b",
			expected: false,
		},
		{
			name:     "empty target",
			input:    "abc",
			target:   "",
			expected: true,
		},
		{
			name:     "empty input",
			input:    "",
			target:   "a",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EndsWith(tt.input, tt.target)
			if result != tt.expected {
				t.Errorf("EndsWith() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		separator string
		expected  []string
	}{
		{
			name:      "split by dash",
			input:     "a-b-c",
			separator: "-",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "split by empty string",
			input:     "abc",
			separator: "",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "no separator found",
			input:     "abc",
			separator: "-",
			expected:  []string{"abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Split(tt.input, tt.separator)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Split() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestWords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "punctuation separated",
			input:    "fred, barney, & pebbles",
			expected: []string{"fred", "barney", "pebbles"},
		},
		{
			name:     "camel case",
			input:    "camelCase",
			expected: []string{"camel", "Case"},
		},
		{
			name:     "snake case",
			input:    "snake_case",
			expected: []string{"snake", "case"},
		},
		{
			name:     "kebab case",
			input:    "kebab-case",
			expected: []string{"kebab", "case"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Words(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Words() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		length   int
		omission []string
		expected string
	}{
		{
			name:     "basic truncate",
			input:    "hi-diddly-ho there, neighborino",
			length:   24,
			expected: "hi-diddly-ho there, n...",
		},
		{
			name:     "custom omission",
			input:    "hi-diddly-ho there, neighborino",
			length:   24,
			omission: []string{" [...]"},
			expected: "hi-diddly-ho there [...]",
		},
		{
			name:     "no truncation needed",
			input:    "short",
			length:   10,
			expected: "short",
		},
		{
			name:     "zero length",
			input:    "test",
			length:   0,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result string
			if len(tt.omission) > 0 {
				result = Truncate(tt.input, tt.length, tt.omission[0])
			} else {
				result = Truncate(tt.input, tt.length)
			}
			if result != tt.expected {
				t.Errorf("Truncate() = %v, want %v", result, tt.expected)
			}
		})
	}
}
