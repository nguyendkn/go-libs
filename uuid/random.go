package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
)

// RandomGenerator defines the interface for random number generators
type RandomGenerator interface {
	NextUint32() uint32
}

// CryptoRandom implements RandomGenerator using crypto/rand
type CryptoRandom struct {
	mu     sync.Mutex
	buffer []byte
	cursor int
}

// NewCryptoRandom creates a new cryptographically secure random number generator
func NewCryptoRandom() *CryptoRandom {
	return &CryptoRandom{
		buffer: make([]byte, 32), // Buffer for 8 uint32 values
		cursor: len(make([]byte, 32)), // Force initial fill
	}
}

// NextUint32 returns a cryptographically secure random uint32
func (r *CryptoRandom) NextUint32() uint32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Refill buffer if needed
	if r.cursor >= len(r.buffer) {
		if _, err := rand.Read(r.buffer); err != nil {
			panic(err) // Should never happen with crypto/rand
		}
		r.cursor = 0
	}
	
	// Read next uint32 from buffer
	val := binary.BigEndian.Uint32(r.buffer[r.cursor:])
	r.cursor += 4
	return val
}

// BufferedCryptoRandom wraps crypto.getRandomValues() to enable buffering
// This uses a small buffer by default to avoid both unbearable throughput 
// decline in some environments and the waste of time and space for unused values.
type BufferedCryptoRandom struct {
	buffer [8]uint32
	cursor int
	mu     sync.Mutex
}

// NewBufferedCryptoRandom creates a new buffered crypto random generator
func NewBufferedCryptoRandom() *BufferedCryptoRandom {
	return &BufferedCryptoRandom{
		cursor: len([8]uint32{}), // Force initial fill
	}
}

// NextUint32 returns a cryptographically secure random uint32 with buffering
func (r *BufferedCryptoRandom) NextUint32() uint32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.cursor >= len(r.buffer) {
		// Fill buffer with random bytes
		var bytes [32]byte // 8 * 4 bytes
		if _, err := rand.Read(bytes[:]); err != nil {
			panic(err)
		}
		
		// Convert to uint32 array
		for i := 0; i < 8; i++ {
			r.buffer[i] = binary.BigEndian.Uint32(bytes[i*4:])
		}
		r.cursor = 0
	}
	
	val := r.buffer[r.cursor]
	r.cursor++
	return val
}

// getDefaultRandom returns the default random number generator available
func getDefaultRandom() RandomGenerator {
	return NewBufferedCryptoRandom()
}

// SecureRandom provides a global instance of secure random generator
var SecureRandom RandomGenerator = getDefaultRandom()

// SetRandomGenerator allows setting a custom random generator for testing
func SetRandomGenerator(rng RandomGenerator) {
	SecureRandom = rng
}

// TestRandom is a deterministic random generator for testing purposes
type TestRandom struct {
	values []uint32
	index  int
}

// NewTestRandom creates a new test random generator with predefined values
func NewTestRandom(values ...uint32) *TestRandom {
	return &TestRandom{
		values: values,
		index:  0,
	}
}

// NextUint32 returns the next predefined value, cycling through the array
func (r *TestRandom) NextUint32() uint32 {
	if len(r.values) == 0 {
		return 0
	}
	
	val := r.values[r.index%len(r.values)]
	r.index++
	return val
}

// Reset resets the test random generator to the beginning
func (r *TestRandom) Reset() {
	r.index = 0
}
