package uuid

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	// Regex patterns for different UUID formats
	hex32Pattern = regexp.MustCompile(`^[0-9a-fA-F]{32}$`)
	hex36Pattern = regexp.MustCompile(`^([0-9a-fA-F]{8})-([0-9a-fA-F]{4})-([0-9a-fA-F]{4})-([0-9a-fA-F]{4})-([0-9a-fA-F]{12})$`)
	hex38Pattern = regexp.MustCompile(`^\{([0-9a-fA-F]{8})-([0-9a-fA-F]{4})-([0-9a-fA-F]{4})-([0-9a-fA-F]{4})-([0-9a-fA-F]{12})\}$`)
	urnPattern   = regexp.MustCompile(`^urn:uuid:([0-9a-fA-F]{8})-([0-9a-fA-F]{4})-([0-9a-fA-F]{4})-([0-9a-fA-F]{4})-([0-9a-fA-F]{12})$`)
)

// Parse builds a UUID from a string representation.
//
// This function accepts the following formats:
//
//   - 32-digit hexadecimal format without hyphens: `0189dcd553117d408db09496a2eef37b`
//   - 8-4-4-4-12 hyphenated format: `0189dcd5-5311-7d40-8db0-9496a2eef37b`
//   - Hyphenated format with surrounding braces: `{0189dcd5-5311-7d40-8db0-9496a2eef37b}`
//   - RFC 9562 URN format: `urn:uuid:0189dcd5-5311-7d40-8db0-9496a2eef37b`
//
// Leading and trailing whitespaces represent an error.
//
// Returns an error if the argument could not be parsed as a valid UUID string.
func Parse(s string) (UUID, error) {
	var hex string
	
	switch len(s) {
	case UUIDStringLength32:
		// 32-digit hex without hyphens
		if !hex32Pattern.MatchString(s) {
			return UUID{}, ErrInvalidUUIDString
		}
		hex = strings.ToLower(s)
		
	case UUIDStringLength36:
		// 8-4-4-4-12 hyphenated format
		matches := hex36Pattern.FindStringSubmatch(s)
		if matches == nil {
			return UUID{}, ErrInvalidUUIDString
		}
		hex = strings.ToLower(strings.Join(matches[1:], ""))
		
	case UUIDStringLength38:
		// Hyphenated format with braces
		matches := hex38Pattern.FindStringSubmatch(s)
		if matches == nil {
			return UUID{}, ErrInvalidUUIDString
		}
		hex = strings.ToLower(strings.Join(matches[1:], ""))
		
	case UUIDStringLength45:
		// RFC 9562 URN format
		matches := urnPattern.FindStringSubmatch(s)
		if matches == nil {
			return UUID{}, ErrInvalidUUIDString
		}
		hex = strings.ToLower(strings.Join(matches[1:], ""))
		
	default:
		return UUID{}, ErrInvalidUUIDString
	}
	
	return parseHex(hex)
}

// parseHex converts a 32-character hex string to UUID bytes
func parseHex(hex string) (UUID, error) {
	if len(hex) != 32 {
		return UUID{}, ErrInvalidUUIDString
	}
	
	var bytes [UUIDSize]byte
	
	// Parse 4 bytes at a time for efficiency
	for i := 0; i < UUIDSize; i += 4 {
		// Parse 8 hex characters (4 bytes) at once
		chunk := hex[2*i : 2*i+8]
		val, err := strconv.ParseUint(chunk, 16, 32)
		if err != nil {
			return UUID{}, ErrInvalidUUIDString
		}
		
		// Convert to big-endian bytes
		bytes[i] = byte(val >> 24)
		bytes[i+1] = byte(val >> 16)
		bytes[i+2] = byte(val >> 8)
		bytes[i+3] = byte(val)
	}
	
	return UUID{bytes: bytes}, nil
}

// MustParse is like Parse but panics if the string cannot be parsed.
// It simplifies safe initialization of UUID constants.
func MustParse(s string) UUID {
	uuid, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return uuid
}

// IsValid reports whether s is a valid UUID string in any supported format.
func IsValid(s string) bool {
	_, err := Parse(s)
	return err == nil
}

// Nil returns the nil UUID (all zeros)
func Nil() UUID {
	return UUID{}
}

// Max returns the max UUID (all 0xFF)
func Max() UUID {
	var bytes [UUIDSize]byte
	for i := range bytes {
		bytes[i] = 0xff
	}
	return UUID{bytes: bytes}
}
