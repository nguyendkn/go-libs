package main

import (
	"fmt"
	"time"

	"github.com/nguyendkn/go-libs/lodash/array"
	"github.com/nguyendkn/go-libs/lodash/collection"
	"github.com/nguyendkn/go-libs/lodash/date"
	"github.com/nguyendkn/go-libs/lodash/function"
	"github.com/nguyendkn/go-libs/lodash/lang"
	"github.com/nguyendkn/go-libs/lodash/math"
	"github.com/nguyendkn/go-libs/lodash/object"
	str "github.com/nguyendkn/go-libs/lodash/string"
	"github.com/nguyendkn/go-libs/lodash/util"
)

func main() {
	fmt.Println("=== Go-Lodash Demo ===")

	// Array operations
	fmt.Println("🔢 Array Operations:")

	// Chunk
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
	chunks := array.Chunk(numbers, 3)
	fmt.Printf("Chunk([1,2,3,4,5,6,7,8], 3): %v\n", chunks)

	// Compact
	mixed := []interface{}{0, 1, false, 2, "", 3, nil, 4}
	compact := array.Compact(mixed)
	fmt.Printf("Compact([0,1,false,2,'',3,nil,4]): %v\n", compact)

	// Uniq
	duplicates := []int{1, 2, 2, 3, 3, 3, 4}
	unique := array.Uniq(duplicates)
	fmt.Printf("Uniq([1,2,2,3,3,3,4]): %v\n", unique)

	fmt.Println()

	// Collection operations
	fmt.Println("📚 Collection Operations:")

	// Filter
	filtered := collection.Filter(numbers, func(x int) bool { return x%2 == 0 })
	fmt.Printf("Filter(evens): %v\n", filtered)

	// Map
	doubled := collection.Map([]int{1, 2, 3}, func(x int) int { return x * 2 })
	fmt.Printf("Map(double): %v\n", doubled)

	// Reduce
	sum := collection.Reduce(numbers, func(acc, x int) int { return acc + x }, 0)
	fmt.Printf("Reduce(sum): %d\n", sum)

	// GroupBy
	words := []string{"one", "two", "three", "four", "five"}
	grouped := collection.GroupBy(words, func(s string) int { return len(s) })
	fmt.Printf("GroupBy(length): %v\n", grouped)

	fmt.Println()

	// String operations
	fmt.Println("📝 String Operations:")

	text := "hello world"
	fmt.Printf("Original: '%s'\n", text)
	fmt.Printf("CamelCase: '%s'\n", str.CamelCase(text))
	fmt.Printf("KebabCase: '%s'\n", str.KebabCase(text))
	fmt.Printf("SnakeCase: '%s'\n", str.SnakeCase(text))
	fmt.Printf("PascalCase: '%s'\n", str.PascalCase(text))

	longText := "This is a very long text that needs to be truncated"
	fmt.Printf("Truncate(30): '%s'\n", str.Truncate(longText, 30))

	fmt.Printf("Words('%s'): %v\n", "hello, world & universe", str.Words("hello, world & universe"))

	fmt.Println()

	// Object operations
	fmt.Println("🗂️  Object Operations:")

	data := map[string]interface{}{
		"name":    "John",
		"age":     30,
		"city":    "New York",
		"country": "USA",
	}

	keys := object.Keys(data)
	fmt.Printf("Keys: %v\n", keys)

	picked := object.Pick(data, []string{"name", "age"})
	fmt.Printf("Pick(name, age): %v\n", picked)

	omitted := object.Omit(data, []string{"city", "country"})
	fmt.Printf("Omit(city, country): %v\n", omitted)

	fmt.Println()

	// Math operations
	fmt.Println("🧮 Math Operations:")

	nums := []int{10, 5, 8, 3, 12, 7}

	max, _ := math.Max(nums)
	min, _ := math.Min(nums)
	sum = math.Sum(nums)
	mean, _ := math.Mean(nums)

	fmt.Printf("Numbers: %v\n", nums)
	fmt.Printf("Max: %d, Min: %d\n", max, min)
	fmt.Printf("Sum: %d, Mean: %.2f\n", sum, mean)

	clamped := math.Clamp(15, 5, 10)
	fmt.Printf("Clamp(15, 5, 10): %d\n", clamped)

	inRange := math.InRange(7, 5, 10)
	fmt.Printf("InRange(7, 5, 10): %t\n", inRange)

	fmt.Println()

	// Function operations
	fmt.Println("⚡ Function Operations:")

	// Debounce example
	counter := 0
	debounced := function.Debounce(func() {
		counter++
		fmt.Printf("Debounced function called! Count: %d\n", counter)
	}, 100*time.Millisecond)

	fmt.Println("Calling debounced function multiple times...")
	debounced()
	debounced()
	debounced()

	time.Sleep(150 * time.Millisecond) // Wait for debounce

	// Memoize example
	expensiveFunc := function.Memoize(func(n int) int {
		fmt.Printf("Computing expensive operation for %d...\n", n)
		time.Sleep(10 * time.Millisecond) // Simulate expensive operation
		return n * n
	})

	fmt.Println("First call to memoized function:")
	result1 := expensiveFunc(5)
	fmt.Printf("Result: %d\n", result1)

	fmt.Println("Second call to memoized function (should use cache):")
	result2 := expensiveFunc(5)
	fmt.Printf("Result: %d\n", result2)

	// Once example
	initFunc := function.Once(func() string {
		fmt.Println("Initialization function called!")
		return "initialized"
	})

	fmt.Println("Calling once function multiple times:")
	fmt.Printf("First call: %s\n", initFunc())
	fmt.Printf("Second call: %s\n", initFunc())
	fmt.Printf("Third call: %s\n", initFunc())

	fmt.Println()

	// Lang operations
	fmt.Println("🔍 Lang Operations:")

	// Type checking
	fmt.Printf("IsArray([]int{1,2,3}): %t\n", lang.IsArray([]int{1, 2, 3}))
	fmt.Printf("IsString('hello'): %t\n", lang.IsString("hello"))
	fmt.Printf("IsNumber(42): %t\n", lang.IsNumber(42))

	// Cloning
	original := []int{1, 2, 3}
	cloned := lang.Clone(original).([]int)
	fmt.Printf("Original: %v, Cloned: %v\n", original, cloned)

	// Type conversion
	converted := lang.ToString(42)
	fmt.Printf("ToString(42): '%s'\n", converted)

	fmt.Println()

	// Util operations
	fmt.Println("🛠️  Util Operations:")

	// Range generation
	range1 := util.Range(5)
	fmt.Printf("Range(5): %v\n", range1)

	range2 := util.Range(2, 8, 2)
	fmt.Printf("Range(2, 8, 2): %v\n", range2)

	// Times
	squares := util.Times(5, func(i int) int { return i * i })
	fmt.Printf("Times(5, square): %v\n", squares)

	// Unique ID
	id1 := util.UniqueId()
	id2 := util.UniqueId("user_")
	fmt.Printf("UniqueId(): %s, UniqueId('user_'): %s\n", id1, id2)

	// Default values
	defaulted := util.DefaultTo(0, 42)
	fmt.Printf("DefaultTo(0, 42): %d\n", defaulted)

	fmt.Println()

	// Date operations
	fmt.Println("📅 Date Operations:")

	now := time.Now()
	fmt.Printf("Current time: %s\n", now.Format("2006-01-02 15:04:05"))

	// Date utilities
	startOfDay := date.StartOfDay(now)
	endOfDay := date.EndOfDay(now)
	fmt.Printf("Start of day: %s\n", date.Format(startOfDay, "2006-01-02 15:04:05"))
	fmt.Printf("End of day: %s\n", date.Format(endOfDay, "2006-01-02 15:04:05"))

	// Date arithmetic
	tomorrow := date.Add(now, 24*time.Hour)
	fmt.Printf("Tomorrow: %s\n", date.Format(tomorrow, "2006-01-02"))

	// Date validation
	fmt.Printf("IsDate(time.Now()): %t\n", date.IsDate(now))
	fmt.Printf("IsLeapYear(2024): %t\n", date.IsLeapYear(2024))

	// Timestamp
	timestamp := date.Now()
	fmt.Printf("Current timestamp: %d\n", timestamp)

	fmt.Println()
	fmt.Println("✅ Demo completed!")
}
