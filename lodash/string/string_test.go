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

func TestDeburr(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "French characters",
			input:    "déjà vu",
			expected: "deja vu",
		},
		{
			name:     "cafe with accent",
			input:    "café",
			expected: "cafe",
		},
		{
			name:     "naive with diaeresis",
			input:    "naïve",
			expected: "naive",
		},
		{
			name:     "Spanish characters",
			input:    "niño",
			expected: "nino",
		},
		{
			name:     "German characters",
			input:    "Müller",
			expected: "Muller",
		},
		{
			name:     "mixed case with accents",
			input:    "CAFÉ and café",
			expected: "CAFE and cafe",
		},
		{
			name:     "no accents",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "special characters",
			input:    "Æsop's Œuvre",
			expected: "Asop's Ouvre",
		},
		{
			name:     "German eszett",
			input:    "Straße",
			expected: "Strase",
		},
		{
			name:     "comprehensive test",
			input:    "Àlphà Bètà Gàmmà",
			expected: "Alpha Beta Gamma",
		},
		{
			name:     "numbers and symbols unchanged",
			input:    "test123!@#",
			expected: "test123!@#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Deburr(tt.input)
			if result != tt.expected {
				t.Errorf("Deburr() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ampersand",
			input:    "fred, barney, & pebbles",
			expected: "fred, barney, &amp; pebbles",
		},
		{
			name:     "script tag",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "double quotes",
			input:    `He said "Hello"`,
			expected: "He said &quot;Hello&quot;",
		},
		{
			name:     "single quotes",
			input:    "It's a test",
			expected: "It&#39;s a test",
		},
		{
			name:     "all special characters",
			input:    `<>&"'`,
			expected: "&lt;&gt;&amp;&quot;&#39;",
		},
		{
			name:     "no special characters",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "mixed content",
			input:    `<div class="test">Hello & "world"</div>`,
			expected: "&lt;div class=&quot;test&quot;&gt;Hello &amp; &quot;world&quot;&lt;/div&gt;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Escape(tt.input)
			if result != tt.expected {
				t.Errorf("Escape() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUnescape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ampersand",
			input:    "fred, barney, &amp; pebbles",
			expected: "fred, barney, & pebbles",
		},
		{
			name:     "script tag",
			input:    "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
			expected: "<script>alert('xss')</script>",
		},
		{
			name:     "double quotes",
			input:    "He said &quot;Hello&quot;",
			expected: `He said "Hello"`,
		},
		{
			name:     "single quotes",
			input:    "It&#39;s a test",
			expected: "It's a test",
		},
		{
			name:     "alternative single quote",
			input:    "It&#x27;s a test",
			expected: "It's a test",
		},
		{
			name:     "all special characters",
			input:    "&lt;&gt;&amp;&quot;&#39;",
			expected: `<>&"'`,
		},
		{
			name:     "no entities",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "mixed content",
			input:    "&lt;div class=&quot;test&quot;&gt;Hello &amp; &quot;world&quot;&lt;/div&gt;",
			expected: `<div class="test">Hello & "world"</div>`,
		},
		{
			name:     "partial entities (should not change)",
			input:    "&am; &l; &g;",
			expected: "&am; &l; &g;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Unescape(tt.input)
			if result != tt.expected {
				t.Errorf("Unescape() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestEscapeRegExp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lodash URL",
			input:    "[lodash](https://lodash.com/)",
			expected: "\\[lodash\\]\\(https://lodash\\.com/\\)",
		},
		{
			name:     "price with dollar and dot",
			input:    "$100.00",
			expected: "\\$100\\.00",
		},
		{
			name:     "all special characters",
			input:    "^$\\.\\*+?()[]{}|",
			expected: "\\^\\$\\\\\\.\\\\\\*\\+\\?\\(\\)\\[\\]\\{\\}\\|",
		},
		{
			name:     "caret and dollar",
			input:    "^start$",
			expected: "\\^start\\$",
		},
		{
			name:     "parentheses and brackets",
			input:    "(test)[array]{object}",
			expected: "\\(test\\)\\[array\\]\\{object\\}",
		},
		{
			name:     "wildcard and plus",
			input:    "*.txt+",
			expected: "\\*\\.txt\\+",
		},
		{
			name:     "question mark and pipe",
			input:    "test?|backup",
			expected: "test\\?\\|backup",
		},
		{
			name:     "no special characters",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "backslash",
			input:    "path\\to\\file",
			expected: "path\\\\to\\\\file",
		},
		{
			name:     "mixed content",
			input:    "Hello (world) $100.00!",
			expected: "Hello \\(world\\) \\$100\\.00!",
		},
		{
			name:     "email pattern",
			input:    "user@domain.com",
			expected: "user@domain\\.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeRegExp(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeRegExp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLowerCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "kebab case",
			input:    "--Foo-Bar--",
			expected: "foo bar",
		},
		{
			name:     "camel case",
			input:    "fooBar",
			expected: "foo bar",
		},
		{
			name:     "snake case",
			input:    "__FOO_BAR__",
			expected: "foo bar",
		},
		{
			name:     "pascal case",
			input:    "FooBar",
			expected: "foo bar",
		},
		{
			name:     "mixed separators",
			input:    "foo-bar_baz.qux",
			expected: "foo bar baz qux",
		},
		{
			name:     "with numbers",
			input:    "foo2Bar3",
			expected: "foo 2 bar 3",
		},
		{
			name:     "all uppercase",
			input:    "HELLO WORLD",
			expected: "hello world",
		},
		{
			name:     "already lowercase",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "Hello",
			expected: "hello",
		},
		{
			name:     "XML parser example",
			input:    "XMLParser",
			expected: "xmlparser", // Current implementation doesn't split consecutive uppercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LowerCase(tt.input)
			if result != tt.expected {
				t.Errorf("LowerCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUpperCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "kebab case",
			input:    "--Foo-Bar--",
			expected: "FOO BAR",
		},
		{
			name:     "camel case",
			input:    "fooBar",
			expected: "FOO BAR",
		},
		{
			name:     "snake case",
			input:    "__FOO_BAR__",
			expected: "FOO BAR",
		},
		{
			name:     "pascal case",
			input:    "FooBar",
			expected: "FOO BAR",
		},
		{
			name:     "mixed separators",
			input:    "foo-bar_baz.qux",
			expected: "FOO BAR BAZ QUX",
		},
		{
			name:     "with numbers",
			input:    "foo2Bar3",
			expected: "FOO 2 BAR 3",
		},
		{
			name:     "all lowercase",
			input:    "hello world",
			expected: "HELLO WORLD",
		},
		{
			name:     "already uppercase",
			input:    "HELLO WORLD",
			expected: "HELLO WORLD",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "HELLO",
		},
		{
			name:     "XML parser example",
			input:    "XMLParser",
			expected: "XMLPARSER", // Current implementation doesn't split consecutive uppercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpperCase(tt.input)
			if result != tt.expected {
				t.Errorf("UpperCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "all uppercase",
			input:    "HELLO WORLD",
			expected: "hello world",
		},
		{
			name:     "mixed case",
			input:    "FooBar",
			expected: "foobar",
		},
		{
			name:     "already lowercase",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "with numbers and symbols",
			input:    "Hello123!@#",
			expected: "hello123!@#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToLower(tt.input)
			if result != tt.expected {
				t.Errorf("ToLower() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "all lowercase",
			input:    "hello world",
			expected: "HELLO WORLD",
		},
		{
			name:     "mixed case",
			input:    "FooBar",
			expected: "FOOBAR",
		},
		{
			name:     "already uppercase",
			input:    "HELLO WORLD",
			expected: "HELLO WORLD",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "with numbers and symbols",
			input:    "hello123!@#",
			expected: "HELLO123!@#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToUpper(tt.input)
			if result != tt.expected {
				t.Errorf("ToUpper() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStartCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "kebab case",
			input:    "--foo-bar--",
			expected: "Foo Bar",
		},
		{
			name:     "camel case",
			input:    "fooBar",
			expected: "Foo Bar",
		},
		{
			name:     "snake case",
			input:    "__FOO_BAR__",
			expected: "Foo Bar",
		},
		{
			name:     "pascal case",
			input:    "FooBar",
			expected: "Foo Bar",
		},
		{
			name:     "mixed separators",
			input:    "foo-bar_baz.qux",
			expected: "Foo Bar Baz Qux",
		},
		{
			name:     "with numbers",
			input:    "foo2Bar3",
			expected: "Foo 2 Bar 3",
		},
		{
			name:     "all uppercase",
			input:    "HELLO WORLD",
			expected: "Hello World",
		},
		{
			name:     "all lowercase",
			input:    "hello world",
			expected: "Hello World",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single word",
			input:    "hello",
			expected: "Hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StartCase(tt.input)
			if result != tt.expected {
				t.Errorf("StartCase() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		radix    []int
		expected int64
	}{
		{
			name:     "decimal number",
			str:      "42",
			radix:    []int{},
			expected: 42,
		},
		{
			name:     "decimal with radix 10",
			str:      "42",
			radix:    []int{10},
			expected: 42,
		},
		{
			name:     "binary number",
			str:      "1010",
			radix:    []int{2},
			expected: 10,
		},
		{
			name:     "hexadecimal number",
			str:      "ff",
			radix:    []int{16},
			expected: 255,
		},
		{
			name:     "hexadecimal with 0x prefix",
			str:      "0x10",
			radix:    []int{},
			expected: 16,
		},
		{
			name:     "octal number",
			str:      "77",
			radix:    []int{8},
			expected: 63,
		},
		{
			name:     "negative number",
			str:      "-42",
			radix:    []int{},
			expected: -42,
		},
		{
			name:     "positive sign",
			str:      "+42",
			radix:    []int{},
			expected: 42,
		},
		{
			name:     "with whitespace",
			str:      "  42  ",
			radix:    []int{},
			expected: 42,
		},
		{
			name:     "invalid characters stop parsing",
			str:      "42abc",
			radix:    []int{},
			expected: 42,
		},
		{
			name:     "empty string",
			str:      "",
			radix:    []int{},
			expected: 0,
		},
		{
			name:     "invalid radix",
			str:      "42",
			radix:    []int{1},
			expected: 0,
		},
		{
			name:     "base 36",
			str:      "zz",
			radix:    []int{36},
			expected: 1295, // 35*36 + 35
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseInt(tt.str, tt.radix...)
			if result != tt.expected {
				t.Errorf("ParseInt() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		name        string
		str         string
		pattern     string
		replacement string
		expected    string
	}{
		{
			name:        "basic replacement",
			str:         "Hi Fred",
			pattern:     "Fred",
			replacement: "Barney",
			expected:    "Hi Barney",
		},
		{
			name:        "replace word",
			str:         "hello world",
			pattern:     "world",
			replacement: "Go",
			expected:    "hello Go",
		},
		{
			name:        "first occurrence only",
			str:         "hello hello",
			pattern:     "hello",
			replacement: "hi",
			expected:    "hi hello",
		},
		{
			name:        "no match",
			str:         "hello world",
			pattern:     "foo",
			replacement: "bar",
			expected:    "hello world",
		},
		{
			name:        "empty string",
			str:         "",
			pattern:     "test",
			replacement: "replace",
			expected:    "",
		},
		{
			name:        "empty pattern",
			str:         "hello",
			pattern:     "",
			replacement: "test",
			expected:    "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Replace(tt.str, tt.pattern, tt.replacement)
			if result != tt.expected {
				t.Errorf("Replace() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestReplaceAll(t *testing.T) {
	tests := []struct {
		name        string
		str         string
		pattern     string
		replacement string
		expected    string
	}{
		{
			name:        "replace all occurrences",
			str:         "hello hello",
			pattern:     "hello",
			replacement: "hi",
			expected:    "hi hi",
		},
		{
			name:        "multiple replacements",
			str:         "foo bar foo",
			pattern:     "foo",
			replacement: "baz",
			expected:    "baz bar baz",
		},
		{
			name:        "no match",
			str:         "hello world",
			pattern:     "foo",
			replacement: "bar",
			expected:    "hello world",
		},
		{
			name:        "empty string",
			str:         "",
			pattern:     "test",
			replacement: "replace",
			expected:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceAll(tt.str, tt.pattern, tt.replacement)
			if result != tt.expected {
				t.Errorf("ReplaceAll() = %v, want %v", result, tt.expected)
			}
		})
	}
}
