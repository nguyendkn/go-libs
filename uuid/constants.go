package uuid

import "errors"

// Constants for UUID formatting
const (
	// HexDigits represents the hexadecimal digits used in UUID string representation
	HexDigits = "0123456789abcdef"
	
	// UUIDSize is the size of a UUID in bytes
	UUIDSize = 16
	
	// MaxCounter is the maximum value for the counter in UUIDv7
	MaxCounter = 0x3ff_ffff_ffff
	
	// DefaultRollbackAllowance is the default rollback allowance in milliseconds
	DefaultRollbackAllowance = 10_000
)

// UUID variant types
type Variant string

const (
	VarNil      Variant = "NIL"
	VarMax      Variant = "MAX"
	Var0        Variant = "VAR_0"
	Var10       Variant = "VAR_10"
	Var110      Variant = "VAR_110"
	VarReserved Variant = "VAR_RESERVED"
)

// Common errors
var (
	ErrInvalidLength     = errors.New("not 128-bit length")
	ErrInvalidFieldValue = errors.New("invalid field value")
	ErrInvalidUUIDString = errors.New("could not parse UUID string")
	ErrInvalidTimestamp  = errors.New("unixTsMs must be a 48-bit positive integer")
	ErrInvalidRollback   = errors.New("rollbackAllowance out of reasonable range")
	ErrNoSecureRNG       = errors.New("no cryptographically strong RNG available")
)

// UUID string format lengths
const (
	UUIDStringLength32  = 32  // 32-digit hex without hyphens
	UUIDStringLength36  = 36  // 8-4-4-4-12 hyphenated format
	UUIDStringLength38  = 38  // hyphenated format with braces
	UUIDStringLength45  = 45  // RFC 9562 URN format
)

// Bit manipulation helpers
const (
	// Version 7 identifier
	Version7Mask = 0x70
	
	// Version 4 identifier  
	Version4Mask = 0x40
	
	// Variant 10 identifier
	Variant10Mask = 0x80
)
