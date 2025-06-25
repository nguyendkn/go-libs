package uuid

import (
	"sync"
	"time"
)

// V7Generator encapsulates the monotonic counter state.
//
// This struct provides APIs to utilize a separate counter state from that of the
// default generator used by UUIDv7() and UUIDv7Obj(). In addition to
// the default Generate method, this struct has GenerateOrAbort
// that is useful to absolutely guarantee the monotonically increasing order of
// generated UUIDs.
type V7Generator struct {
	mu        sync.Mutex
	timestamp uint64
	counter   uint64
	random    RandomGenerator
}

// NewV7Generator creates a generator object with the default random number generator,
// or with the specified one if passed as an argument. The specified random
// number generator should be cryptographically strong and securely seeded.
func NewV7Generator(rng ...RandomGenerator) *V7Generator {
	var random RandomGenerator
	if len(rng) > 0 && rng[0] != nil {
		random = rng[0]
	} else {
		random = getDefaultRandom()
	}
	
	return &V7Generator{
		random: random,
	}
}

// Generate generates a new UUIDv7 object from the current timestamp, or resets the
// generator upon significant timestamp rollback.
//
// This method returns a monotonically increasing UUID by reusing the previous
// timestamp even if the up-to-date timestamp is smaller than the immediately
// preceding UUID's. However, when such a clock rollback is considered
// significant (i.e., by more than ten seconds), this method resets the
// generator and returns a new UUID based on the given timestamp, breaking the
// increasing order of UUIDs.
func (g *V7Generator) Generate() UUID {
	return g.GenerateOrResetCore(uint64(time.Now().UnixMilli()), DefaultRollbackAllowance)
}

// GenerateOrAbort generates a new UUIDv7 object from the current timestamp, or returns
// an error upon significant timestamp rollback.
//
// This method returns a monotonically increasing UUID by reusing the previous
// timestamp even if the up-to-date timestamp is smaller than the immediately
// preceding UUID's. However, when such a clock rollback is considered
// significant (i.e., by more than ten seconds), this method aborts and
// returns an error immediately.
func (g *V7Generator) GenerateOrAbort() (UUID, error) {
	return g.GenerateOrAbortCore(uint64(time.Now().UnixMilli()), DefaultRollbackAllowance)
}

// GenerateOrResetCore generates a new UUIDv7 object from the unixTsMs passed, or resets the
// generator upon significant timestamp rollback.
//
// This method is equivalent to Generate except that it takes a custom
// timestamp and clock rollback allowance.
//
// Parameters:
//   - unixTsMs: Unix timestamp in milliseconds
//   - rollbackAllowance: The amount of unixTsMs rollback that is considered significant.
//     A suggested value is 10_000 (milliseconds).
//
// Returns an error if unixTsMs is not a 48-bit positive integer.
func (g *V7Generator) GenerateOrResetCore(unixTsMs, rollbackAllowance uint64) UUID {
	uuid, err := g.GenerateOrAbortCore(unixTsMs, rollbackAllowance)
	if err != nil {
		// Reset state and resume
		g.mu.Lock()
		g.timestamp = 0
		g.mu.Unlock()
		
		uuid, err = g.GenerateOrAbortCore(unixTsMs, rollbackAllowance)
		if err != nil {
			panic(err) // Should not happen after reset
		}
	}
	return uuid
}

// GenerateOrAbortCore generates a new UUIDv7 object from the unixTsMs passed, or returns
// an error upon significant timestamp rollback.
//
// This method is equivalent to GenerateOrAbort except that it takes a
// custom timestamp and clock rollback allowance.
//
// Parameters:
//   - unixTsMs: Unix timestamp in milliseconds
//   - rollbackAllowance: The amount of unixTsMs rollback that is considered significant.
//     A suggested value is 10_000 (milliseconds).
//
// Returns an error if unixTsMs is not a 48-bit positive integer.
func (g *V7Generator) GenerateOrAbortCore(unixTsMs, rollbackAllowance uint64) (UUID, error) {
	if unixTsMs < 1 || unixTsMs > 0xffff_ffff_ffff {
		return UUID{}, ErrInvalidTimestamp
	}
	if rollbackAllowance > 0xffff_ffff_ffff {
		return UUID{}, ErrInvalidRollback
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	if unixTsMs > g.timestamp {
		g.timestamp = unixTsMs
		g.resetCounter()
	} else if unixTsMs+rollbackAllowance >= g.timestamp {
		// Go on with previous timestamp if new one is not much smaller
		g.counter++
		if g.counter > MaxCounter {
			// Increment timestamp at counter overflow
			g.timestamp++
			g.resetCounter()
		}
	} else {
		// Abort if clock went backwards to unbearable extent
		return UUID{}, ErrInvalidTimestamp
	}

	// Split counter into rand_a (12 bits) and rand_b_hi (30 bits)
	randA := g.counter >> 30
	randBHi := g.counter & (1<<30 - 1)
	randBLo := uint64(g.random.NextUint32())

	return FromFieldsV7(g.timestamp, randA, randBHi, randBLo)
}

// resetCounter initializes the counter at a 42-bit random integer
func (g *V7Generator) resetCounter() {
	// Generate 42-bit random counter
	hi := uint64(g.random.NextUint32()) & 0x3ff // 10 bits
	lo := uint64(g.random.NextUint32())         // 32 bits
	g.counter = hi<<32 | lo
}

// GenerateV4 generates a new UUIDv4 object utilizing the random number generator inside.
func (g *V7Generator) GenerateV4() UUID {
	var bytes [UUIDSize]byte
	
	// Fill with random bytes
	for i := 0; i < UUIDSize; i += 4 {
		val := g.random.NextUint32()
		bytes[i] = byte(val >> 24)
		bytes[i+1] = byte(val >> 16)
		bytes[i+2] = byte(val >> 8)
		bytes[i+3] = byte(val)
	}
	
	// Set version (4) and variant (10)
	bytes[6] = Version4Mask | (bytes[6] & 0x0f)
	bytes[8] = Variant10Mask | (bytes[8] & 0x3f)
	
	return UUID{bytes: bytes}
}

// Default generator instance
var defaultGenerator *V7Generator
var defaultGeneratorOnce sync.Once

// getDefaultGenerator returns the default V7Generator instance
func getDefaultGenerator() *V7Generator {
	defaultGeneratorOnce.Do(func() {
		defaultGenerator = NewV7Generator()
	})
	return defaultGenerator
}
