// Package string provides utility functions for working with strings.
// All functions are thread-safe and designed for high performance.
package string

import (
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

// Deburr converts Latin-1 Supplement and Latin Extended-A letters to basic Latin letters and removes combining diacritical marks.
//
// Example:
//
//	Deburr("déjà vu") // "deja vu"
//	Deburr("café") // "cafe"
//	Deburr("naïve") // "naive"
func Deburr(s string) string {
	if s == "" {
		return ""
	}

	// Map of common diacritical characters to their base forms
	deburMap := map[rune]rune{
		// Latin-1 Supplement
		'À': 'A', 'Á': 'A', 'Â': 'A', 'Ã': 'A', 'Ä': 'A', 'Å': 'A',
		'à': 'a', 'á': 'a', 'â': 'a', 'ã': 'a', 'ä': 'a', 'å': 'a',
		'Ç': 'C', 'ç': 'c',
		'È': 'E', 'É': 'E', 'Ê': 'E', 'Ë': 'E',
		'è': 'e', 'é': 'e', 'ê': 'e', 'ë': 'e',
		'Ì': 'I', 'Í': 'I', 'Î': 'I', 'Ï': 'I',
		'ì': 'i', 'í': 'i', 'î': 'i', 'ï': 'i',
		'Ñ': 'N', 'ñ': 'n',
		'Ò': 'O', 'Ó': 'O', 'Ô': 'O', 'Õ': 'O', 'Ö': 'O',
		'ò': 'o', 'ó': 'o', 'ô': 'o', 'õ': 'o', 'ö': 'o',
		'Ù': 'U', 'Ú': 'U', 'Û': 'U', 'Ü': 'U',
		'ù': 'u', 'ú': 'u', 'û': 'u', 'ü': 'u',
		'Ý': 'Y', 'ý': 'y', 'ÿ': 'y',
		// Additional common characters
		'Æ': 'A', 'æ': 'a',
		'Œ': 'O', 'œ': 'o',
		'ß': 's',
	}

	var result strings.Builder
	for _, r := range s {
		if replacement, exists := deburMap[r]; exists {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// Escape converts the characters "&", "<", ">", '"', and "'" in string to their corresponding HTML entities.
//
// Example:
//
//	Escape("fred, barney, & pebbles") // "fred, barney, &amp; pebbles"
//	Escape("<script>alert('xss')</script>") // "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"
func Escape(s string) string {
	if s == "" {
		return ""
	}

	// HTML entity mappings
	escapeMap := map[rune]string{
		'&':  "&amp;",
		'<':  "&lt;",
		'>':  "&gt;",
		'"':  "&quot;",
		'\'': "&#39;",
	}

	var result strings.Builder
	for _, r := range s {
		if entity, exists := escapeMap[r]; exists {
			result.WriteString(entity)
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// Unescape converts HTML entities "&amp;", "&lt;", "&gt;", "&quot;", and "&#39;" in string to their corresponding characters.
//
// Example:
//
//	Unescape("fred, barney, &amp; pebbles") // "fred, barney, & pebbles"
//	Unescape("&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;") // "<script>alert('xss')</script>"
func Unescape(s string) string {
	if s == "" {
		return ""
	}

	// HTML entity mappings (reverse of escape)
	unescapeMap := map[string]string{
		"&amp;":  "&",
		"&lt;":   "<",
		"&gt;":   ">",
		"&quot;": "\"",
		"&#39;":  "'",
		"&#x27;": "'", // Alternative single quote encoding
	}

	result := s
	for entity, char := range unescapeMap {
		result = strings.ReplaceAll(result, entity, char)
	}

	return result
}

// EscapeRegExp escapes the RegExp special characters "^", "$", "\", ".", "*", "+", "?", "(", ")", "[", "]", "{", "}", and "|" in string.
//
// Example:
//
//	EscapeRegExp("[lodash](https://lodash.com/)") // "\\[lodash\\]\\(https://lodash\\.com/\\)"
//	EscapeRegExp("$100.00") // "\\$100\\.00"
func EscapeRegExp(s string) string {
	if s == "" {
		return ""
	}

	// RegExp special characters that need escaping
	regexpChars := map[rune]bool{
		'^': true, '$': true, '\\': true, '.': true, '*': true, '+': true,
		'?': true, '(': true, ')': true, '[': true, ']': true, '{': true,
		'}': true, '|': true,
	}

	var result strings.Builder
	for _, r := range s {
		if regexpChars[r] {
			result.WriteRune('\\')
		}
		result.WriteRune(r)
	}

	return result.String()
}

// LowerCase converts string, as space separated words, to lower case.
//
// Example:
//
//	LowerCase("--Foo-Bar--") // "foo bar"
//	LowerCase("fooBar") // "foo bar"
//	LowerCase("__FOO_BAR__") // "foo bar"
func LowerCase(s string) string {
	words := extractWords(s)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	for i, word := range words {
		if i > 0 {
			result.WriteRune(' ')
		}
		result.WriteString(strings.ToLower(word))
	}

	return result.String()
}

// UpperCase converts string, as space separated words, to upper case.
//
// Example:
//
//	UpperCase("--Foo-Bar--") // "FOO BAR"
//	UpperCase("fooBar") // "FOO BAR"
//	UpperCase("__FOO_BAR__") // "FOO BAR"
func UpperCase(s string) string {
	words := extractWords(s)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	for i, word := range words {
		if i > 0 {
			result.WriteRune(' ')
		}
		result.WriteString(strings.ToUpper(word))
	}

	return result.String()
}

// extractWords extracts words from a string, handling camelCase, snake_case, kebab-case, etc.
func extractWords(s string) []string {
	if s == "" {
		return []string{}
	}

	var words []string
	var currentWord strings.Builder

	runes := []rune(s)
	for i, r := range runes {
		if isWordSeparator(r) {
			// Skip separators, but finalize current word if any
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		} else if i > 0 && isWordBoundary(runes[i-1], r) {
			// Word boundary detected (e.g., camelCase transition)
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			currentWord.WriteRune(r)
		} else {
			currentWord.WriteRune(r)
		}
	}

	// Add the last word if any
	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	// Filter out empty words and non-alphabetic words
	var result []string
	for _, word := range words {
		if word != "" && hasAlphaNumeric(word) {
			result = append(result, word)
		}
	}

	return result
}

// hasAlphaNumeric checks if string contains at least one alphanumeric character
func hasAlphaNumeric(s string) bool {
	for _, r := range s {
		if isLetter(r) || isDigit(r) {
			return true
		}
	}
	return false
}

// isWordSeparator checks if a rune is a word separator
func isWordSeparator(r rune) bool {
	return r == ' ' || r == '-' || r == '_' || r == '.' || r == '/' || r == '\\' ||
		r == ',' || r == ';' || r == ':' || r == '!' || r == '?' || r == '@' ||
		r == '#' || r == '$' || r == '%' || r == '^' || r == '&' || r == '*' ||
		r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' ||
		r == '|' || r == '=' || r == '+' || r == '<' || r == '>' || r == '~' ||
		r == '`' || r == '"' || r == '\''
}

// isWordBoundary checks if there's a word boundary between two runes
func isWordBoundary(prev, curr rune) bool {
	// Transition from lowercase to uppercase (camelCase)
	if isLower(prev) && isUpper(curr) {
		return true
	}
	// Transition from letter to digit or digit to letter
	if (isLetter(prev) && isDigit(curr)) || (isDigit(prev) && isLetter(curr)) {
		return true
	}
	return false
}

// Helper functions for character classification
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isUpper(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

// ToLower converts string to lower case.
//
// Example:
//
//	ToLower("HELLO WORLD") // "hello world"
//	ToLower("FooBar") // "foobar"
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper converts string to upper case.
//
// Example:
//
//	ToUpper("hello world") // "HELLO WORLD"
//	ToUpper("FooBar") // "FOOBAR"
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// StartCase converts string to start case.
//
// Example:
//
//	StartCase("--foo-bar--") // "Foo Bar"
//	StartCase("fooBar") // "Foo Bar"
//	StartCase("__FOO_BAR__") // "FOO BAR"
func StartCase(s string) string {
	words := extractWords(s)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	for i, word := range words {
		if i > 0 {
			result.WriteRune(' ')
		}
		// Capitalize first letter, lowercase the rest
		if len(word) > 0 {
			runes := []rune(word)
			result.WriteRune(unicode.ToUpper(runes[0]))
			for _, r := range runes[1:] {
				result.WriteRune(unicode.ToLower(r))
			}
		}
	}

	return result.String()
}

// ParseInt converts string to an integer of the specified radix. If radix is undefined or 0, a radix of 10 is used unless the value is a hexadecimal, in which case a radix of 16 is used.
//
// Example:
//
//	ParseInt("08") // 8
//	ParseInt("10", 2) // 2 (binary)
//	ParseInt("ff", 16) // 255 (hexadecimal)
//	ParseInt("0x10") // 16 (auto-detect hex)
func ParseInt(str string, radix ...int) int64 {
	if str == "" {
		return 0
	}

	// Trim whitespace
	str = strings.TrimSpace(str)
	if str == "" {
		return 0
	}

	// Determine radix
	base := 10
	if len(radix) > 0 && radix[0] != 0 {
		base = radix[0]
		if base < 2 || base > 36 {
			return 0 // Invalid radix
		}
	}

	// Handle sign
	negative := false
	if strings.HasPrefix(str, "-") {
		negative = true
		str = str[1:]
	} else if strings.HasPrefix(str, "+") {
		str = str[1:]
	}

	// Auto-detect hexadecimal if radix is 0 or 16
	if (base == 10 || base == 16) && (strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X")) {
		base = 16
		str = str[2:]
	}

	// Parse the number
	var result int64
	for _, char := range str {
		var digit int

		if char >= '0' && char <= '9' {
			digit = int(char - '0')
		} else if char >= 'a' && char <= 'z' {
			digit = int(char - 'a' + 10)
		} else if char >= 'A' && char <= 'Z' {
			digit = int(char - 'A' + 10)
		} else {
			// Invalid character, stop parsing
			break
		}

		if digit >= base {
			// Invalid digit for this base, stop parsing
			break
		}

		result = result*int64(base) + int64(digit)
	}

	if negative {
		result = -result
	}

	return result
}

// Replace replaces matches for pattern in string with replacement.
//
// Example:
//
//	Replace("Hi Fred", "Fred", "Barney") // "Hi Barney"
//	Replace("hello world", "world", "Go") // "hello Go"
func Replace(str, pattern, replacement string) string {
	if str == "" || pattern == "" {
		return str
	}

	return strings.Replace(str, pattern, replacement, 1) // Replace only first occurrence
}

// ReplaceAll replaces all matches for pattern in string with replacement.
//
// Example:
//
//	ReplaceAll("hello hello", "hello", "hi") // "hi hi"
//	ReplaceAll("foo bar foo", "foo", "baz") // "baz bar baz"
func ReplaceAll(str, pattern, replacement string) string {
	if str == "" || pattern == "" {
		return str
	}

	return strings.ReplaceAll(str, pattern, replacement)
}
