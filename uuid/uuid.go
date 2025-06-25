package uuid

import (
	"encoding/json"
)

// UUID represents a UUID as a 16-byte byte array
type UUID struct {
	bytes [UUIDSize]byte
}

// OfInner creates a UUID object from the internal representation, a 16-byte byte array
// containing the binary UUID representation in the big-endian byte order.
//
// This method does NOT copy the argument, and thus the created object
// holds the reference to the underlying buffer.
func OfInner(bytes [UUIDSize]byte) UUID {
	return UUID{bytes: bytes}
}

// OfInnerSlice creates a UUID object from a byte slice.
// Returns an error if the length is not 16 bytes.
func OfInnerSlice(bytes []byte) (UUID, error) {
	if len(bytes) != UUIDSize {
		return UUID{}, ErrInvalidLength
	}

	var arr [UUIDSize]byte
	copy(arr[:], bytes)
	return UUID{bytes: arr}, nil
}

// FromFieldsV7 builds a UUID from UUIDv7 field values.
//
// Parameters:
//   - unixTsMs: A 48-bit unix_ts_ms field value
//   - randA: A 12-bit rand_a field value
//   - randBHi: The higher 30 bits of 62-bit rand_b field value
//   - randBLo: The lower 32 bits of 62-bit rand_b field value
//
// Returns an error if any field value is out of the specified range.
func FromFieldsV7(unixTsMs, randA, randBHi, randBLo uint64) (UUID, error) {
	if unixTsMs > 0xffff_ffff_ffff ||
		randA > 0xfff ||
		randBHi > 0x3fff_ffff ||
		randBLo > 0xffff_ffff {
		return UUID{}, ErrInvalidFieldValue
	}

	var bytes [UUIDSize]byte

	// Set timestamp (48 bits)
	bytes[0] = byte(unixTsMs >> 40)
	bytes[1] = byte(unixTsMs >> 32)
	bytes[2] = byte(unixTsMs >> 24)
	bytes[3] = byte(unixTsMs >> 16)
	bytes[4] = byte(unixTsMs >> 8)
	bytes[5] = byte(unixTsMs)

	// Set version (4 bits) and rand_a (12 bits)
	bytes[6] = Version7Mask | byte(randA>>8)
	bytes[7] = byte(randA)

	// Set variant (2 bits) and rand_b_hi (30 bits)
	// Variant10Mask = 0x80 (10xxxxxx), we need 6 bits from randBHi for the remaining bits
	bytes[8] = Variant10Mask | byte((randBHi>>24)&0x3f)
	bytes[9] = byte(randBHi >> 16)
	bytes[10] = byte(randBHi >> 8)
	bytes[11] = byte(randBHi)

	// Set rand_b_lo (32 bits)
	bytes[12] = byte(randBLo >> 24)
	bytes[13] = byte(randBLo >> 16)
	bytes[14] = byte(randBLo >> 8)
	bytes[15] = byte(randBLo)

	return UUID{bytes: bytes}, nil
}

// Bytes returns a copy of the UUID's byte array
func (u UUID) Bytes() [UUIDSize]byte {
	return u.bytes
}

// String returns the 8-4-4-4-12 canonical hexadecimal string representation
// (e.g., "0189dcd5-5311-7d40-8db0-9496a2eef37b")
func (u UUID) String() string {
	var buf [36]byte
	encodeHex(buf[:], u.bytes[:])
	return string(buf[:])
}

// StringNoDash returns the 32-digit hexadecimal string representation without hyphens
// (e.g., "0189dcd553117d408db09496a2eef37b")
func (u UUID) StringNoDash() string {
	var buf [32]byte
	encodeHexNoDash(buf[:], u.bytes[:])
	return string(buf[:])
}

// Hex returns the 32-digit hexadecimal representation without hyphens
// (alias for StringNoDash)
func (u UUID) Hex() string {
	return u.StringNoDash()
}

// MarshalJSON implements the json.Marshaler interface
func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (u *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := Parse(s)
	if err != nil {
		return err
	}

	*u = parsed
	return nil
}

// GetVariant reports the variant field value of the UUID
func (u UUID) GetVariant() Variant {
	b := u.bytes[8]
	n := b >> 4

	if n <= 0b0111 {
		if u.isNil() {
			return VarNil
		}
		return Var0
	} else if n <= 0b1011 {
		return Var10
	} else if n <= 0b1101 {
		return Var110
	} else {
		if u.isMax() {
			return VarMax
		}
		return VarReserved
	}
}

// GetVersion returns the version field value of the UUID or 0 if the UUID does
// not have the variant field value of VAR_10
func (u UUID) GetVersion() int {
	if u.GetVariant() == Var10 {
		return int(u.bytes[6] >> 4)
	}
	return 0
}

// Clone creates a copy of the UUID
func (u UUID) Clone() UUID {
	return UUID{bytes: u.bytes}
}

// Equals returns true if this UUID is equivalent to other
func (u UUID) Equals(other UUID) bool {
	return u.bytes == other.bytes
}

// CompareTo returns a negative integer, zero, or positive integer if this UUID is less
// than, equal to, or greater than other, respectively
func (u UUID) CompareTo(other UUID) int {
	for i := 0; i < UUIDSize; i++ {
		if u.bytes[i] < other.bytes[i] {
			return -1
		} else if u.bytes[i] > other.bytes[i] {
			return 1
		}
	}
	return 0
}

// isNil checks if this is the nil UUID
func (u UUID) isNil() bool {
	for _, b := range u.bytes {
		if b != 0 {
			return false
		}
	}
	return true
}

// isMax checks if this is the max UUID
func (u UUID) isMax() bool {
	for _, b := range u.bytes {
		if b != 0xff {
			return false
		}
	}
	return true
}

// encodeHex encodes src into dst as hexadecimal with hyphens
func encodeHex(dst []byte, src []byte) {
	const hextable = "0123456789abcdef"

	j := 0
	for i, v := range src {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		j += 2

		// Add hyphens at positions 8, 12, 16, 20
		if i == 3 || i == 5 || i == 7 || i == 9 {
			dst[j] = '-'
			j++
		}
	}
}

// encodeHexNoDash encodes src into dst as hexadecimal without hyphens
func encodeHexNoDash(dst []byte, src []byte) {
	const hextable = "0123456789abcdef"

	for i, v := range src {
		dst[i*2] = hextable[v>>4]
		dst[i*2+1] = hextable[v&0x0f]
	}
}
