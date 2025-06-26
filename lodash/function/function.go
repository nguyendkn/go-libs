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

// Curry2 creates a curried version of a function that takes 2 arguments.
//
// Example:
//
//	add := func(a, b int) int { return a + b }
//	curriedAdd := Curry2(add)
//	addFive := curriedAdd(5)
//	result := addFive(3) // 8
func Curry2[T1, T2, R any](fn func(T1, T2) R) func(T1) func(T2) R {
	return func(arg1 T1) func(T2) R {
		return func(arg2 T2) R {
			return fn(arg1, arg2)
		}
	}
}

// Curry3 creates a curried version of a function that takes 3 arguments.
//
// Example:
//
//	add3 := func(a, b, c int) int { return a + b + c }
//	curriedAdd3 := Curry3(add3)
//	result := curriedAdd3(1)(2)(3) // 6
func Curry3[T1, T2, T3, R any](fn func(T1, T2, T3) R) func(T1) func(T2) func(T3) R {
	return func(arg1 T1) func(T2) func(T3) R {
		return func(arg2 T2) func(T3) R {
			return func(arg3 T3) R {
				return fn(arg1, arg2, arg3)
			}
		}
	}
}

// Curry4 creates a curried version of a function that takes 4 arguments.
//
// Example:
//
//	add4 := func(a, b, c, d int) int { return a + b + c + d }
//	curriedAdd4 := Curry4(add4)
//	result := curriedAdd4(1)(2)(3)(4) // 10
func Curry4[T1, T2, T3, T4, R any](fn func(T1, T2, T3, T4) R) func(T1) func(T2) func(T3) func(T4) R {
	return func(arg1 T1) func(T2) func(T3) func(T4) R {
		return func(arg2 T2) func(T3) func(T4) R {
			return func(arg3 T3) func(T4) R {
				return func(arg4 T4) R {
					return fn(arg1, arg2, arg3, arg4)
				}
			}
		}
	}
}

// Partial2 creates a function that invokes func with partials prepended to the arguments it receives.
//
// Example:
//
//	greet := func(greeting, name string) string { return greeting + " " + name }
//	sayHello := Partial2(greet, "Hello")
//	result := sayHello("World") // "Hello World"
func Partial2[T1, T2, R any](fn func(T1, T2) R, arg1 T1) func(T2) R {
	return func(arg2 T2) R {
		return fn(arg1, arg2)
	}
}

// Partial3 creates a function that invokes func with partials prepended to the arguments it receives.
//
// Example:
//
//	add3 := func(a, b, c int) int { return a + b + c }
//	addToFive := Partial3(add3, 2, 3)
//	result := addToFive(4) // 9
func Partial3[T1, T2, T3, R any](fn func(T1, T2, T3) R, arg1 T1, arg2 T2) func(T3) R {
	return func(arg3 T3) R {
		return fn(arg1, arg2, arg3)
	}
}

// Partial4 creates a function that invokes func with partials prepended to the arguments it receives.
//
// Example:
//
//	add4 := func(a, b, c, d int) int { return a + b + c + d }
//	addToSix := Partial4(add4, 1, 2, 3)
//	result := addToSix(4) // 10
func Partial4[T1, T2, T3, T4, R any](fn func(T1, T2, T3, T4) R, arg1 T1, arg2 T2, arg3 T3) func(T4) R {
	return func(arg4 T4) R {
		return fn(arg1, arg2, arg3, arg4)
	}
}

// Flip creates a function that invokes func with arguments flipped.
//
// Example:
//
//	divide := func(a, b float64) float64 { return a / b }
//	flippedDivide := Flip2(divide)
//	result := flippedDivide(2, 10) // 10 / 2 = 5
func Flip2[T1, T2, R any](fn func(T1, T2) R) func(T2, T1) R {
	return func(arg2 T2, arg1 T1) R {
		return fn(arg1, arg2)
	}
}

// Flip3 creates a function that invokes func with arguments flipped.
//
// Example:
//
//	subtract := func(a, b, c int) int { return a - b - c }
//	flippedSubtract := Flip3(subtract)
//	result := flippedSubtract(1, 2, 10) // 10 - 2 - 1 = 7
func Flip3[T1, T2, T3, R any](fn func(T1, T2, T3) R) func(T3, T2, T1) R {
	return func(arg3 T3, arg2 T2, arg1 T1) R {
		return fn(arg1, arg2, arg3)
	}
}

// Rearg2 creates a function that invokes func with arguments arranged according to the specified indexes.
//
// Example:
//
//	greet := func(greeting, name string) string { return greeting + " " + name }
//	reordered := Rearg2(greet, 1, 0) // Swap arguments
//	result := reordered("World", "Hello") // "Hello World"
func Rearg2[T1, T2, R any](fn func(T1, T2) R, index1, index2 int) func(T1, T2) R {
	return func(arg1 T1, arg2 T2) R {
		args := []interface{}{arg1, arg2}
		if index1 == 0 && index2 == 1 {
			return fn(args[0].(T1), args[1].(T2))
		} else if index1 == 1 && index2 == 0 {
			return fn(args[1].(T1), args[0].(T2))
		}
		// Default case
		return fn(arg1, arg2)
	}
}

// Ary creates a function that invokes func with up to n arguments, ignoring any additional arguments.
//
// Example:
//
//	sum := func(args ...int) int {
//		total := 0
//		for _, v := range args { total += v }
//		return total
//	}
//	sumTwo := Ary(sum, 2)
//	result := sumTwo(1, 2, 3, 4, 5) // Only uses first 2 arguments: 1 + 2 = 3
func Ary(fn func(...interface{}) interface{}, n int) func(...interface{}) interface{} {
	return func(args ...interface{}) interface{} {
		if len(args) > n {
			args = args[:n]
		}
		return fn(args...)
	}
}

// Unary creates a function that accepts up to one argument, ignoring any additional arguments.
//
// Example:
//
//	parseIntFunc := func(s string, base int) int {
//		// Simplified parseInt
//		if s == "10" { return 10 }
//		return 0
//	}
//	unaryParseInt := Unary2(parseIntFunc)
//	result := unaryParseInt("10") // Only uses first argument
func Unary2[T1, T2, R any](fn func(T1, T2) R, defaultArg2 T2) func(T1) R {
	return func(arg1 T1) R {
		return fn(arg1, defaultArg2)
	}
}
