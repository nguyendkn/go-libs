// Package string provides utility functions for working with strings.
// All functions are thread-safe and designed for high performance.
package string

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// CamelCase converts string to camel case.
//
// Example:
//
//	CamelCase("foo bar") // "fooBar"
//	CamelCase("--foo-bar--") // "fooBar"
//	CamelCase("__FOO_BAR__") // "fooBar"
func CamelCase(s string) string {
	words := extractWords(s)
	if len(words) == 0 {
		return ""
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		result += Capitalize(strings.ToLower(words[i]))
	}
	return result
}

// KebabCase converts string to kebab case.
//
// Example:
//
//	KebabCase("Foo Bar") // "foo-bar"
//	KebabCase("fooBar") // "foo-bar"
//	KebabCase("__FOO_BAR__") // "foo-bar"
func KebabCase(s string) string {
	words := extractWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "-")
}

// SnakeCase converts string to snake case.
//
// Example:
//
//	SnakeCase("Foo Bar") // "foo_bar"
//	SnakeCase("fooBar") // "foo_bar"
//	SnakeCase("--FOO-BAR--") // "foo_bar"
func SnakeCase(s string) string {
	words := extractWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return strings.Join(words, "_")
}

// PascalCase converts string to pascal case.
//
// Example:
//
//	PascalCase("foo bar") // "FooBar"
//	PascalCase("--foo-bar--") // "FooBar"
//	PascalCase("__FOO_BAR__") // "FooBar"
func PascalCase(s string) string {
	words := extractWords(s)
	var result strings.Builder
	for _, word := range words {
		result.WriteString(Capitalize(strings.ToLower(word)))
	}
	return result.String()
}

// Capitalize converts the first character of string to upper case and the remaining to lower case.
//
// Example:
//
//	Capitalize("FRED") // "Fred"
//	Capitalize("fRED") // "Fred"
func Capitalize(s string) string {
	if s == "" {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + strings.ToLower(s[size:])
}

// LowerFirst converts the first character of string to lower case.
//
// Example:
//
//	LowerFirst("Fred") // "fred"
//	LowerFirst("FRED") // "fRED"
func LowerFirst(s string) string {
	if s == "" {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[size:]
}

// UpperFirst converts the first character of string to upper case.
//
// Example:
//
//	UpperFirst("fred") // "Fred"
//	UpperFirst("FRED") // "FRED"
func UpperFirst(s string) string {
	if s == "" {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[size:]
}

// Trim removes leading and trailing whitespace from string.
//
// Example:
//
//	Trim("  abc  ") // "abc"
//	Trim("\t\nabc\t\n") // "abc"
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// TrimStart removes leading whitespace from string.
//
// Example:
//
//	TrimStart("  abc  ") // "abc  "
func TrimStart(s string) string {
	return strings.TrimLeftFunc(s, unicode.IsSpace)
}

// TrimEnd removes trailing whitespace from string.
//
// Example:
//
//	TrimEnd("  abc  ") // "  abc"
func TrimEnd(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// Pad pads string on the left and right sides if it's shorter than length.
// Padding characters are truncated if they can't be evenly divided by length.
//
// Example:
//
//	Pad("abc", 8, "_-") // "_-abc_-_"
//	Pad("abc", 6, "_") // "_abc__"
func Pad(s string, length int, chars string) string {
	if chars == "" {
		chars = " "
	}

	strLen := utf8.RuneCountInString(s)
	if strLen >= length {
		return s
	}

	padLen := length - strLen
	leftPad := padLen / 2
	rightPad := padLen - leftPad

	return PadStart(PadEnd(s, strLen+rightPad, chars), length, chars)
}

// PadStart pads string on the left side if it's shorter than length.
//
// Example:
//
//	PadStart("abc", 6, "_-") // "_-_abc"
//	PadStart("abc", 6, "_") // "___abc"
func PadStart(s string, length int, chars string) string {
	if chars == "" {
		chars = " "
	}

	strLen := utf8.RuneCountInString(s)
	if strLen >= length {
		return s
	}

	padLen := length - strLen
	pad := strings.Repeat(chars, (padLen/len(chars))+1)
	padRunes := []rune(pad)

	return string(padRunes[:padLen]) + s
}

// PadEnd pads string on the right side if it's shorter than length.
//
// Example:
//
//	PadEnd("abc", 6, "_-") // "abc_-_"
//	PadEnd("abc", 6, "_") // "abc___"
func PadEnd(s string, length int, chars string) string {
	if chars == "" {
		chars = " "
	}

	strLen := utf8.RuneCountInString(s)
	if strLen >= length {
		return s
	}

	padLen := length - strLen
	pad := strings.Repeat(chars, (padLen/len(chars))+1)
	padRunes := []rune(pad)

	return s + string(padRunes[:padLen])
}

// Repeat repeats the given string n times.
//
// Example:
//
//	Repeat("*", 3) // "***"
//	Repeat("abc", 2) // "abcabc"
func Repeat(s string, n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(s, n)
}

// StartsWith checks if string starts with target.
//
// Example:
//
//	StartsWith("abc", "a") // true
//	StartsWith("abc", "b") // false
func StartsWith(s, target string) bool {
	return strings.HasPrefix(s, target)
}

// EndsWith checks if string ends with target.
//
// Example:
//
//	EndsWith("abc", "c") // true
//	EndsWith("abc", "b") // false
func EndsWith(s, target string) bool {
	return strings.HasSuffix(s, target)
}

// Split splits string by separator.
//
// Example:
//
//	Split("a-b-c", "-") // []string{"a", "b", "c"}
//	Split("a-b-c", "") // []string{"a", "-", "b", "-", "c"}
func Split(s, separator string) []string {
	if separator == "" {
		return strings.Split(s, "")
	}
	return strings.Split(s, separator)
}

// Words splits string into an array of its words.
//
// Example:
//
//	Words("fred, barney, & pebbles") // []string{"fred", "barney", "pebbles"}
//	Words("camelCase") // []string{"camel", "Case"}
func Words(s string) []string {
	return extractWords(s)
}

// extractWords extracts words from a string using various delimiters and case changes
func extractWords(s string) []string {
	if s == "" {
		return []string{}
	}

	// Replace common delimiters and punctuation with spaces
	re := regexp.MustCompile(`[_\-\s,&]+`)
	s = re.ReplaceAllString(s, " ")

	// Split on case changes (camelCase, PascalCase)
	re = regexp.MustCompile(`([a-z])([A-Z])`)
	s = re.ReplaceAllString(s, "$1 $2")

	// Split on number boundaries
	re = regexp.MustCompile(`([a-zA-Z])(\d)`)
	s = re.ReplaceAllString(s, "$1 $2")
	re = regexp.MustCompile(`(\d)([a-zA-Z])`)
	s = re.ReplaceAllString(s, "$1 $2")

	// Clean up and split
	words := strings.Fields(s)
	var result []string
	for _, word := range words {
		// Only include alphabetic words
		if word != "" && regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(word) {
			result = append(result, word)
		}
	}

	return result
}

// Truncate truncates string if it's longer than the given maximum string length.
//
// Example:
//
//	Truncate("hi-diddly-ho there, neighborino", 24) // "hi-diddly-ho there, n..."
//	Truncate("hi-diddly-ho there, neighborino", 24, " [...]") // "hi-diddly-ho there[...]"
func Truncate(s string, length int, omission ...string) string {
	if length <= 0 {
		return ""
	}

	omit := "..."
	if len(omission) > 0 {
		omit = omission[0]
	}

	runes := []rune(s)
	if len(runes) <= length {
		return s
	}

	omitRunes := []rune(omit)
	if len(omitRunes) >= length {
		return string(omitRunes[:length])
	}

	// Calculate the position to truncate, considering the omission length
	truncatePos := max(0, length-len(omitRunes))

	return string(runes[:truncatePos]) + omit
}
