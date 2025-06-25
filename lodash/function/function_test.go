package function

import (
	"sync"
	"testing"
	"time"
)

func TestDebounce(t *testing.T) {
	var called bool
	var mutex sync.Mutex

	fn := func() {
		mutex.Lock()
		called = true
		mutex.Unlock()
	}

	debounced := Debounce(fn, 50*time.Millisecond)

	// Call multiple times quickly
	debounced()
	debounced()
	debounced()

	// Should not be called yet
	mutex.Lock()
	if called {
		t.Error("Function should not be called immediately")
	}
	mutex.Unlock()

	// Wait for debounce period
	time.Sleep(60 * time.Millisecond)

	mutex.Lock()
	if !called {
		t.Error("Function should be called after debounce period")
	}
	mutex.Unlock()
}

func TestThrottle(t *testing.T) {
	var callCount int
	var mutex sync.Mutex

	fn := func() {
		mutex.Lock()
		callCount++
		mutex.Unlock()
	}

	throttled := Throttle(fn, 50*time.Millisecond)

	// First call should execute immediately
	throttled()
	
	mutex.Lock()
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
	mutex.Unlock()

	// Subsequent calls within throttle period should be ignored
	throttled()
	throttled()

	mutex.Lock()
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
	mutex.Unlock()

	// Wait for throttle period and call again
	time.Sleep(60 * time.Millisecond)
	throttled()

	mutex.Lock()
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
	mutex.Unlock()
}

func TestOnce(t *testing.T) {
	var callCount int

	fn := func() int {
		callCount++
		return 42
	}

	onceFn := Once(fn)

	// First call
	result1 := onceFn()
	if result1 != 42 {
		t.Errorf("Expected 42, got %d", result1)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second call should return same result without calling original function
	result2 := onceFn()
	if result2 != 42 {
		t.Errorf("Expected 42, got %d", result2)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

func TestOnceVoid(t *testing.T) {
	var callCount int

	fn := func() {
		callCount++
	}

	onceFn := OnceVoid(fn)

	// First call
	onceFn()
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second call should not execute
	onceFn()
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}
}

func TestMemoize(t *testing.T) {
	var callCount int

	fn := func(n int) int {
		callCount++
		return n * 2
	}

	memoized := Memoize(fn)

	// First call with argument 5
	result1 := memoized(5)
	if result1 != 10 {
		t.Errorf("Expected 10, got %d", result1)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second call with same argument should use cache
	result2 := memoized(5)
	if result2 != 10 {
		t.Errorf("Expected 10, got %d", result2)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Call with different argument should execute function
	result3 := memoized(3)
	if result3 != 6 {
		t.Errorf("Expected 6, got %d", result3)
	}
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestAfter(t *testing.T) {
	var callCount int

	fn := func() {
		callCount++
	}

	afterThree := After(3, fn)

	// First two calls should not execute
	afterThree()
	afterThree()
	if callCount != 0 {
		t.Errorf("Expected 0 calls, got %d", callCount)
	}

	// Third call should execute
	afterThree()
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Fourth call should also execute
	afterThree()
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestBefore(t *testing.T) {
	var callCount int

	fn := func() int {
		callCount++
		return callCount * 10
	}

	beforeThree := Before(3, fn)

	// First call
	result1 := beforeThree()
	if result1 != 10 {
		t.Errorf("Expected 10, got %d", result1)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second call
	result2 := beforeThree()
	if result2 != 20 {
		t.Errorf("Expected 20, got %d", result2)
	}
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}

	// Third call should return last result without calling function
	result3 := beforeThree()
	if result3 != 20 {
		t.Errorf("Expected 20, got %d", result3)
	}
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestNegate(t *testing.T) {
	isEven := func(n int) bool {
		return n%2 == 0
	}

	isOdd := Negate(isEven)

	if !isOdd(3) {
		t.Error("Expected true for odd number")
	}

	if isOdd(4) {
		t.Error("Expected false for even number")
	}
}

func TestCompose(t *testing.T) {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }

	// Compose applies functions from right to left
	addOneThenDouble := Compose(double, addOne)
	result := addOneThenDouble(3) // (3 + 1) * 2 = 8

	if result != 8 {
		t.Errorf("Expected 8, got %d", result)
	}
}

func TestPipe(t *testing.T) {
	addOne := func(x int) int { return x + 1 }
	double := func(x int) int { return x * 2 }

	// Pipe applies functions from left to right
	addOneThenDouble := Pipe(addOne, double)
	result := addOneThenDouble(3) // (3 + 1) * 2 = 8

	if result != 8 {
		t.Errorf("Expected 8, got %d", result)
	}
}

func TestDebounceWithArgs(t *testing.T) {
	var lastArg interface{}
	var mutex sync.Mutex

	fn := func(args ...interface{}) {
		mutex.Lock()
		if len(args) > 0 {
			lastArg = args[0]
		}
		mutex.Unlock()
	}

	debounced := DebounceWithArgs(fn, 50*time.Millisecond)

	// Call with different arguments
	debounced("first")
	debounced("second")
	debounced("third")

	// Wait for debounce period
	time.Sleep(60 * time.Millisecond)

	mutex.Lock()
	if lastArg != "third" {
		t.Errorf("Expected 'third', got %v", lastArg)
	}
	mutex.Unlock()
}

func TestThrottleWithArgs(t *testing.T) {
	var lastArg interface{}
	var callCount int
	var mutex sync.Mutex

	fn := func(args ...interface{}) {
		mutex.Lock()
		callCount++
		if len(args) > 0 {
			lastArg = args[0]
		}
		mutex.Unlock()
	}

	throttled := ThrottleWithArgs(fn, 50*time.Millisecond)

	// First call should execute immediately
	throttled("first")
	
	mutex.Lock()
	if callCount != 1 || lastArg != "first" {
		t.Errorf("Expected 1 call with 'first', got %d calls with %v", callCount, lastArg)
	}
	mutex.Unlock()

	// Subsequent calls within throttle period should be ignored
	throttled("second")
	throttled("third")

	mutex.Lock()
	if callCount != 1 || lastArg != "first" {
		t.Errorf("Expected 1 call with 'first', got %d calls with %v", callCount, lastArg)
	}
	mutex.Unlock()
}

func TestDelay(t *testing.T) {
	var called bool
	var mutex sync.Mutex

	fn := func() {
		mutex.Lock()
		called = true
		mutex.Unlock()
	}

	start := time.Now()
	Delay(fn, 50*time.Millisecond)

	// Should not be called immediately
	mutex.Lock()
	if called {
		t.Error("Function should not be called immediately")
	}
	mutex.Unlock()

	// Wait for delay
	time.Sleep(60 * time.Millisecond)

	mutex.Lock()
	if !called {
		t.Error("Function should be called after delay")
	}
	mutex.Unlock()

	elapsed := time.Since(start)
	if elapsed < 50*time.Millisecond {
		t.Errorf("Function called too early: %v", elapsed)
	}
}

func TestDelayWithArgs(t *testing.T) {
	var arg interface{}
	var mutex sync.Mutex

	fn := func(args ...interface{}) {
		mutex.Lock()
		if len(args) > 0 {
			arg = args[0]
		}
		mutex.Unlock()
	}

	DelayWithArgs(fn, 50*time.Millisecond, "test")

	// Wait for delay
	time.Sleep(60 * time.Millisecond)

	mutex.Lock()
	if arg != "test" {
		t.Errorf("Expected 'test', got %v", arg)
	}
	mutex.Unlock()
}

func TestMemoizeWithResolver(t *testing.T) {
	var callCount int

	type Args struct {
		A, B int
	}

	fn := func(args Args) int {
		callCount++
		return args.A + args.B
	}

	resolver := func(args Args) string {
		return string(rune(args.A)) + string(rune(args.B))
	}

	memoized := MemoizeWithResolver(fn, resolver)

	// First call
	result1 := memoized(Args{1, 2})
	if result1 != 3 {
		t.Errorf("Expected 3, got %d", result1)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second call with same arguments should use cache
	result2 := memoized(Args{1, 2})
	if result2 != 3 {
		t.Errorf("Expected 3, got %d", result2)
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Call with different arguments
	result3 := memoized(Args{2, 3})
	if result3 != 5 {
		t.Errorf("Expected 5, got %d", result3)
	}
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}
