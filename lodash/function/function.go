// Package function provides utility functions for working with functions.
// All functions are thread-safe and designed for high performance.
package function

import (
	"sync"
	"time"
)

// Debounce creates a debounced function that delays invoking func until after wait milliseconds
// have elapsed since the last time the debounced function was invoked.
//
// Example:
//
//	debounced := Debounce(func() { fmt.Println("Hello") }, 100*time.Millisecond)
//	debounced() // Will be called after 100ms if no other calls are made
func Debounce(fn func(), wait time.Duration) func() {
	var timer *time.Timer
	var mutex sync.Mutex

	return func() {
		mutex.Lock()
		defer mutex.Unlock()

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(wait, fn)
	}
}

// DebounceWithArgs creates a debounced function that delays invoking func until after wait milliseconds
// have elapsed since the last time the debounced function was invoked. Supports arguments.
//
// Example:
//
//	debounced := DebounceWithArgs(func(args ...interface{}) {
//		fmt.Println("Hello", args[0])
//	}, 100*time.Millisecond)
//	debounced("World") // Will be called after 100ms with "World"
func DebounceWithArgs(fn func(...interface{}), wait time.Duration) func(...interface{}) {
	var timer *time.Timer
	var mutex sync.Mutex

	return func(args ...interface{}) {
		mutex.Lock()
		defer mutex.Unlock()

		if timer != nil {
			timer.Stop()
		}

		timer = time.AfterFunc(wait, func() {
			fn(args...)
		})
	}
}

// Throttle creates a throttled function that only invokes func at most once per every wait milliseconds.
//
// Example:
//
//	throttled := Throttle(func() { fmt.Println("Hello") }, 100*time.Millisecond)
//	throttled() // Will be called immediately
//	throttled() // Will be ignored if called within 100ms
func Throttle(fn func(), wait time.Duration) func() {
	var lastCall time.Time
	var mutex sync.Mutex

	return func() {
		mutex.Lock()
		defer mutex.Unlock()

		now := time.Now()
		if now.Sub(lastCall) >= wait {
			lastCall = now
			fn()
		}
	}
}

// ThrottleWithArgs creates a throttled function that only invokes func at most once per every wait milliseconds.
// Supports arguments.
//
// Example:
//
//	throttled := ThrottleWithArgs(func(args ...interface{}) {
//		fmt.Println("Hello", args[0])
//	}, 100*time.Millisecond)
//	throttled("World") // Will be called immediately with "World"
func ThrottleWithArgs(fn func(...interface{}), wait time.Duration) func(...interface{}) {
	var lastCall time.Time
	var mutex sync.Mutex

	return func(args ...interface{}) {
		mutex.Lock()
		defer mutex.Unlock()

		now := time.Now()
		if now.Sub(lastCall) >= wait {
			lastCall = now
			fn(args...)
		}
	}
}

// Once creates a function that is restricted to invoking func once.
// Repeat calls to the function return the value of the first invocation.
//
// Example:
//
//	initialize := Once(func() int {
//		fmt.Println("Initializing...")
//		return 42
//	})
//	result1 := initialize() // Prints "Initializing..." and returns 42
//	result2 := initialize() // Returns 42 without printing
func Once[T any](fn func() T) func() T {
	var once sync.Once
	var result T

	return func() T {
		once.Do(func() {
			result = fn()
		})
		return result
	}
}

// OnceVoid creates a function that is restricted to invoking func once.
// For functions that don't return values.
//
// Example:
//
//	initialize := OnceVoid(func() { fmt.Println("Initializing...") })
//	initialize() // Prints "Initializing..."
//	initialize() // Does nothing
func OnceVoid(fn func()) func() {
	var once sync.Once

	return func() {
		once.Do(fn)
	}
}

// Memoize creates a function that memoizes the result of func.
// If resolver is provided, it determines the cache key for storing the result based on the arguments.
//
// Example:
//
//	fibonacci := Memoize(func(n int) int {
//		if n <= 1 { return n }
//		return fibonacci(n-1) + fibonacci(n-2)
//	})
//	result := fibonacci(10) // Computed and cached
//	result2 := fibonacci(10) // Retrieved from cache
func Memoize[K comparable, V any](fn func(K) V) func(K) V {
	cache := make(map[K]V)
	var mutex sync.RWMutex

	return func(key K) V {
		mutex.RLock()
		if value, exists := cache[key]; exists {
			mutex.RUnlock()
			return value
		}
		mutex.RUnlock()

		mutex.Lock()
		defer mutex.Unlock()

		// Double-check in case another goroutine computed it
		if value, exists := cache[key]; exists {
			return value
		}

		result := fn(key)
		cache[key] = result
		return result
	}
}

// MemoizeWithResolver creates a function that memoizes the result of func with a custom resolver.
//
// Example:
//
//	add := MemoizeWithResolver(
//		func(a, b int) int { return a + b },
//		func(a, b int) string { return fmt.Sprintf("%d,%d", a, b) },
//	)
//	result := add(1, 2) // Computed and cached with key "1,2"
func MemoizeWithResolver[T any, K comparable, V any](fn func(T) V, resolver func(T) K) func(T) V {
	cache := make(map[K]V)
	var mutex sync.RWMutex

	return func(args T) V {
		key := resolver(args)

		mutex.RLock()
		if value, exists := cache[key]; exists {
			mutex.RUnlock()
			return value
		}
		mutex.RUnlock()

		mutex.Lock()
		defer mutex.Unlock()

		// Double-check in case another goroutine computed it
		if value, exists := cache[key]; exists {
			return value
		}

		result := fn(args)
		cache[key] = result
		return result
	}
}

// Delay invokes func after wait milliseconds.
//
// Example:
//
//	Delay(func() { fmt.Println("Hello") }, 1*time.Second)
//	// Prints "Hello" after 1 second
func Delay(fn func(), wait time.Duration) {
	time.AfterFunc(wait, fn)
}

// DelayWithArgs invokes func after wait milliseconds with arguments.
//
// Example:
//
//	DelayWithArgs(func(args ...interface{}) {
//		fmt.Println("Hello", args[0])
//	}, 1*time.Second, "World")
//	// Prints "Hello World" after 1 second
func DelayWithArgs(fn func(...interface{}), wait time.Duration, args ...interface{}) {
	time.AfterFunc(wait, func() {
		fn(args...)
	})
}

// Defer invokes func on the next tick of the event loop (using goroutine).
//
// Example:
//
//	Defer(func() { fmt.Println("Deferred") })
//	fmt.Println("Immediate")
//	// Prints "Immediate" then "Deferred"
func Defer(fn func()) {
	go fn()
}

// DeferWithArgs invokes func on the next tick with arguments.
//
// Example:
//
//	DeferWithArgs(func(args ...interface{}) {
//		fmt.Println("Deferred", args[0])
//	}, "Hello")
func DeferWithArgs(fn func(...interface{}), args ...interface{}) {
	go func() {
		fn(args...)
	}()
}

// After creates a function that invokes func once it's called n or more times.
//
// Example:
//
//	afterThree := After(3, func() { fmt.Println("Called 3 times!") })
//	afterThree() // Nothing happens
//	afterThree() // Nothing happens
//	afterThree() // Prints "Called 3 times!"
//	afterThree() // Prints "Called 3 times!" again
func After(n int, fn func()) func() {
	var count int
	var mutex sync.Mutex

	return func() {
		mutex.Lock()
		defer mutex.Unlock()

		count++
		if count >= n {
			fn()
		}
	}
}

// Before creates a function that invokes func while it's called less than n times.
// Subsequent calls to the created function return the result of the last func invocation.
//
// Example:
//
//	beforeThree := Before(3, func() int {
//		fmt.Println("Called")
//		return 42
//	})
//	result1 := beforeThree() // Prints "Called", returns 42
//	result2 := beforeThree() // Prints "Called", returns 42
//	result3 := beforeThree() // Returns 42 (last result)
func Before[T any](n int, fn func() T) func() T {
	var count int
	var result T
	var mutex sync.Mutex

	return func() T {
		mutex.Lock()
		defer mutex.Unlock()

		if count < n-1 {
			result = fn()
			count++
		}
		return result
	}
}

// Negate creates a function that negates the result of the predicate func.
//
// Example:
//
//	isEven := func(n int) bool { return n%2 == 0 }
//	isOdd := Negate(isEven)
//	fmt.Println(isOdd(3)) // true
//	fmt.Println(isOdd(4)) // false
func Negate[T any](predicate func(T) bool) func(T) bool {
	return func(value T) bool {
		return !predicate(value)
	}
}

// Compose creates a function that is the composition of the provided functions,
// where each successive invocation is supplied the return value of the previous.
//
// Example:
//
//	addOne := func(x int) int { return x + 1 }
//	double := func(x int) int { return x * 2 }
//	addOneThenDouble := Compose(double, addOne)
//	result := addOneThenDouble(3) // (3 + 1) * 2 = 8
func Compose[T any](fns ...func(T) T) func(T) T {
	return func(value T) T {
		result := value
		for i := len(fns) - 1; i >= 0; i-- {
			result = fns[i](result)
		}
		return result
	}
}

// Pipe creates a function that is the composition of the provided functions,
// where each successive invocation is supplied the return value of the previous.
// This is the reverse of Compose.
//
// Example:
//
//	addOne := func(x int) int { return x + 1 }
//	double := func(x int) int { return x * 2 }
//	addOneThenDouble := Pipe(addOne, double)
//	result := addOneThenDouble(3) // (3 + 1) * 2 = 8
func Pipe[T any](fns ...func(T) T) func(T) T {
	return func(value T) T {
		result := value
		for _, fn := range fns {
			result = fn(result)
		}
		return result
	}
}
