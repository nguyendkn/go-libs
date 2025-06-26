# Date Package

Comprehensive date and time manipulation utilities for Go, providing essential functions for working with time.Time objects. This package offers 20 high-performance, thread-safe functions for date operations, formatting, and calculations.

## Features

- **🚀 High Performance**: Optimized algorithms with minimal allocations
- **🔒 Thread Safe**: All functions are safe for concurrent use
- **🎯 Type Safe**: Strong typing with Go's time.Time
- **📦 Zero Dependencies**: No external dependencies beyond Go standard library
- **✅ Well Tested**: Comprehensive test coverage with edge cases
- **🌍 Timezone Aware**: Proper timezone handling

## Installation

```bash
go get github.com/nguyendkn/go-libs/lodash/date
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    "github.com/nguyendkn/go-libs/lodash/date"
)

func main() {
    now := time.Now()

    // Get current timestamp
    timestamp := date.Now()
    fmt.Println(timestamp) // 1640995200000

    // Start and end of day
    startOfDay := date.StartOfDay(now)
    endOfDay := date.EndOfDay(now)
    fmt.Printf("Day: %v to %v\n", startOfDay, endOfDay)

    // Format date
    formatted := date.Format(now, "2006-01-02 15:04:05")
    fmt.Println(formatted) // 2022-01-01 12:30:45
}
```

## Core Functions

### 🕐 **Time Creation & Conversion**
- **`Now`** - Get current timestamp in milliseconds
- **`ToDate`** - Convert various types to time.Time
- **`IsDate`** - Check if value is a time.Time
- **`IsValid`** - Check if time is valid (not zero)

### 📅 **Date Boundaries**
- **`StartOfDay`** - Get start of day (00:00:00)
- **`EndOfDay`** - Get end of day (23:59:59.999999999)
- **`StartOfWeek`** - Get start of week (Monday 00:00:00)
- **`EndOfWeek`** - Get end of week (Sunday 23:59:59.999999999)
- **`StartOfMonth`** - Get start of month (1st day 00:00:00)
- **`EndOfMonth`** - Get end of month (last day 23:59:59.999999999)
- **`StartOfYear`** - Get start of year (Jan 1st 00:00:00)
- **`EndOfYear`** - Get end of year (Dec 31st 23:59:59.999999999)

### ⏰ **Time Operations**
- **`Add`** - Add duration to time
- **`Sub`** - Subtract time from time
- **`Before`** - Check if time is before another
- **`After`** - Check if time is after another
- **`Equal`** - Check if times are equal

### 🔧 **Utility Functions**
- **`Format`** - Format time with layout
- **`DaysInMonth`** - Get number of days in month
- **`IsLeapYear`** - Check if year is leap year

## Detailed Examples

### Working with Date Boundaries
```go
now := time.Date(2022, 6, 15, 14, 30, 45, 0, time.UTC)

// Day boundaries
startDay := date.StartOfDay(now)    // 2022-06-15 00:00:00
endDay := date.EndOfDay(now)        // 2022-06-15 23:59:59.999999999

// Week boundaries (Monday to Sunday)
startWeek := date.StartOfWeek(now)  // 2022-06-13 00:00:00 (Monday)
endWeek := date.EndOfWeek(now)      // 2022-06-19 23:59:59.999999999 (Sunday)

// Month boundaries
startMonth := date.StartOfMonth(now) // 2022-06-01 00:00:00
endMonth := date.EndOfMonth(now)     // 2022-06-30 23:59:59.999999999

// Year boundaries
startYear := date.StartOfYear(now)   // 2022-01-01 00:00:00
endYear := date.EndOfYear(now)       // 2022-12-31 23:59:59.999999999
```

### Time Calculations
```go
now := time.Now()
future := now.Add(2 * time.Hour)
past := now.Add(-30 * time.Minute)

// Time comparisons
isBefore := date.Before(past, now)     // true
isAfter := date.After(future, now)     // true
isEqual := date.Equal(now, now)        // true

// Duration calculations
duration := date.Sub(future, now)     // 2h0m0s
fmt.Printf("Duration: %v\n", duration)

// Add/subtract time
tomorrow := date.Add(now, 24*time.Hour)
yesterday := date.Add(now, -24*time.Hour)
```

### Date Conversion and Validation
```go
// Convert various types to time.Time
timestamp := int64(1640995200000)
dateFromTimestamp, err := date.ToDate(timestamp)
if err == nil {
    fmt.Println("Converted:", dateFromTimestamp)
}

// Parse string to time
dateStr := "2022-01-01T12:00:00Z"
dateFromString, err := date.ToDate(dateStr)
if err == nil {
    fmt.Println("Parsed:", dateFromString)
}

// Validation
isValidTime := date.IsValid(time.Now())        // true
isZeroTime := date.IsValid(time.Time{})        // false
isTimeType := date.IsDate(time.Now())          // true
isNotTimeType := date.IsDate("2022-01-01")     // false
```

### Calendar Operations
```go
// Month information
feb2020 := time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC)
feb2021 := time.Date(2021, 2, 15, 0, 0, 0, 0, time.UTC)

daysInFeb2020 := date.DaysInMonth(feb2020)  // 29 (leap year)
daysInFeb2021 := date.DaysInMonth(feb2021)  // 28

// Leap year check
isLeap2020 := date.IsLeapYear(2020)  // true
isLeap2021 := date.IsLeapYear(2021)  // false

// Formatting
formatted := date.Format(time.Now(), "Monday, January 2, 2006")
fmt.Println(formatted) // "Wednesday, June 15, 2022"
```

## Performance Notes

- **Memory Efficient**: Functions avoid unnecessary allocations
- **Timezone Aware**: All operations preserve timezone information
- **Standard Library**: Built on Go's robust time package
- **Concurrent Safe**: All functions are thread-safe

## Error Handling

- Functions handle invalid inputs gracefully
- ToDate returns appropriate errors for unparseable values
- Zero time values are handled consistently

## Thread Safety

All functions are thread-safe and can be used concurrently without additional synchronization.

## Contributing

See the main [lodash README](../README.md) for contribution guidelines.

## License

This package is part of the go-libs project and follows the same license terms.