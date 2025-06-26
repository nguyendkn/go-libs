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

func TestCurry2(t *testing.T) {
	add := func(a, b int) int {
		return a + b
	}

	curriedAdd := Curry2(add)
	addFive := curriedAdd(5)

	result := addFive(3)
	if result != 8 {
		t.Errorf("Expected 8, got %d", result)
	}

	// Test with different values
	addTen := curriedAdd(10)
	result2 := addTen(7)
	if result2 != 17 {
		t.Errorf("Expected 17, got %d", result2)
	}
}

func TestCurry3(t *testing.T) {
	add3 := func(a, b, c int) int {
		return a + b + c
	}

	curriedAdd3 := Curry3(add3)

	// Test full currying
	result := curriedAdd3(1)(2)(3)
	if result != 6 {
		t.Errorf("Expected 6, got %d", result)
	}

	// Test partial application
	addToThree := curriedAdd3(1)(2)
	result2 := addToThree(4)
	if result2 != 7 {
		t.Errorf("Expected 7, got %d", result2)
	}
}

func TestCurry4(t *testing.T) {
	add4 := func(a, b, c, d int) int {
		return a + b + c + d
	}

	curriedAdd4 := Curry4(add4)

	// Test full currying
	result := curriedAdd4(1)(2)(3)(4)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}

	// Test partial application
	addToSix := curriedAdd4(1)(2)(3)
	result2 := addToSix(5)
	if result2 != 11 {
		t.Errorf("Expected 11, got %d", result2)
	}
}

func TestPartial2(t *testing.T) {
	greet := func(greeting, name string) string {
		return greeting + " " + name
	}

	sayHello := Partial2(greet, "Hello")

	result := sayHello("World")
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result)
	}

	result2 := sayHello("Go")
	if result2 != "Hello Go" {
		t.Errorf("Expected 'Hello Go', got '%s'", result2)
	}
}

func TestPartial3(t *testing.T) {
	add3 := func(a, b, c int) int {
		return a + b + c
	}

	addToFive := Partial3(add3, 2, 3)

	result := addToFive(4)
	if result != 9 {
		t.Errorf("Expected 9, got %d", result)
	}

	result2 := addToFive(10)
	if result2 != 15 {
		t.Errorf("Expected 15, got %d", result2)
	}
}

func TestPartial4(t *testing.T) {
	add4 := func(a, b, c, d int) int {
		return a + b + c + d
	}

	addToSix := Partial4(add4, 1, 2, 3)

	result := addToSix(4)
	if result != 10 {
		t.Errorf("Expected 10, got %d", result)
	}

	result2 := addToSix(0)
	if result2 != 6 {
		t.Errorf("Expected 6, got %d", result2)
	}
}

func TestFlip2(t *testing.T) {
	divide := func(a, b float64) float64 {
		return a / b
	}

	flippedDivide := Flip2(divide)

	result := flippedDivide(2, 10) // 10 / 2 = 5
	if result != 5.0 {
		t.Errorf("Expected 5.0, got %f", result)
	}

	// Test with different values
	result2 := flippedDivide(4, 20) // 20 / 4 = 5
	if result2 != 5.0 {
		t.Errorf("Expected 5.0, got %f", result2)
	}
}

func TestFlip3(t *testing.T) {
	subtract := func(a, b, c int) int {
		return a - b - c
	}

	flippedSubtract := Flip3(subtract)

	result := flippedSubtract(1, 2, 10) // 10 - 2 - 1 = 7
	if result != 7 {
		t.Errorf("Expected 7, got %d", result)
	}

	// Test with different values
	result2 := flippedSubtract(3, 5, 20) // 20 - 5 - 3 = 12
	if result2 != 12 {
		t.Errorf("Expected 12, got %d", result2)
	}
}

func TestRearg2(t *testing.T) {
	greet := func(greeting, name string) string {
		return greeting + " " + name
	}

	// Test swapping arguments (index 1, 0)
	reordered := Rearg2(greet, 1, 0)
	result := reordered("World", "Hello") // Should be "Hello World"
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result)
	}

	// Test normal order (index 0, 1)
	normal := Rearg2(greet, 0, 1)
	result2 := normal("Hello", "World") // Should be "Hello World"
	if result2 != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result2)
	}
}

func TestAry(t *testing.T) {
	sum := func(args ...interface{}) interface{} {
		total := 0
		for _, v := range args {
			if num, ok := v.(int); ok {
				total += num
			}
		}
		return total
	}

	sumTwo := Ary(sum, 2)

	result := sumTwo(1, 2, 3, 4, 5) // Only uses first 2 arguments: 1 + 2 = 3
	if result != 3 {
		t.Errorf("Expected 3, got %v", result)
	}

	// Test with fewer arguments than limit
	result2 := sumTwo(10) // Only 1 argument: 10
	if result2 != 10 {
		t.Errorf("Expected 10, got %v", result2)
	}
}

func TestUnary2(t *testing.T) {
	concat := func(a, b string) string {
		return a + b
	}

	unaryConcat := Unary2(concat, " World")

	result := unaryConcat("Hello") // "Hello" + " World" = "Hello World"
	if result != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", result)
	}

	result2 := unaryConcat("Hi") // "Hi" + " World" = "Hi World"
	if result2 != "Hi World" {
		t.Errorf("Expected 'Hi World', got '%s'", result2)
	}
}
