package uuid

import "fmt"

// UUIDv7 generates a UUIDv7 string.
//
// Parameters:
//   - noDash: if true, returns 32-digit hex without hyphens; if false, returns 8-4-4-4-12 format
//
// Returns the canonical hexadecimal string representation.
func UUIDv7(noDash ...bool) string {
	uuid := UUIDv7Obj()
	if len(noDash) > 0 && noDash[0] {
		return uuid.StringNoDash()
	}
	return uuid.String()
}

// UUIDv7PrimaryKey generates a UUIDv7 string wrapped in single quotes for database usage.
//
// Parameters:
//   - noDash: if true, returns 32-digit hex without hyphens; if false, returns 8-4-4-4-12 format
//
// Returns the UUID string wrapped in single quotes.
func UUIDv7PrimaryKey(noDash ...bool) string {
	uuid := UUIDv7Obj()
	var s string
	if len(noDash) > 0 && noDash[0] {
		s = uuid.StringNoDash()
	} else {
		s = uuid.String()
	}
	return fmt.Sprintf("'%s'", s)
}

// UUIDv7Obj generates a UUIDv7 object using the default generator.
func UUIDv7Obj() UUID {
	return getDefaultGenerator().Generate()
}

// UUIDv4 generates a UUIDv4 string.
//
// Returns the 8-4-4-4-12 canonical hexadecimal string representation.
func UUIDv4() string {
	return UUIDv4Obj().String()
}

// UUIDv4Obj generates a UUIDv4 object using the default generator.
func UUIDv4Obj() UUID {
	return getDefaultGenerator().GenerateV4()
}

// New is an alias for UUIDv7Obj for compatibility
func New() UUID {
	return UUIDv7Obj()
}

// NewString is an alias for UUIDv7 for compatibility
func NewString() string {
	return UUIDv7()
}

// NewV4 is an alias for UUIDv4Obj for compatibility
func NewV4() UUID {
	return UUIDv4Obj()
}

// NewV4String is an alias for UUIDv4 for compatibility
func NewV4String() string {
	return UUIDv4()
}

// NewV7 is an alias for UUIDv7Obj for compatibility
func NewV7() UUID {
	return UUIDv7Obj()
}

// NewV7String is an alias for UUIDv7 for compatibility
func NewV7String() string {
	return UUIDv7()
}

// Must wraps a UUID generation function and panics if an error occurs.
// This is useful for generating UUIDs in variable declarations.
func Must(uuid UUID, err error) UUID {
	if err != nil {
		panic(err)
	}
	return uuid
}

// FromString is an alias for Parse for compatibility
func FromString(s string) (UUID, error) {
	return Parse(s)
}

// MustFromString is an alias for MustParse for compatibility
func MustFromString(s string) UUID {
	return MustParse(s)
}

// FromBytes creates a UUID from a byte slice
func FromBytes(b []byte) (UUID, error) {
	return OfInnerSlice(b)
}

// MustFromBytes creates a UUID from a byte slice, panicking on error
func MustFromBytes(b []byte) UUID {
	uuid, err := FromBytes(b)
	if err != nil {
		panic(err)
	}
	return uuid
}
