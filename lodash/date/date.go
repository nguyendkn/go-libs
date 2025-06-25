// Package date provides utility functions for working with dates.
// All functions are thread-safe and designed for high performance.
package date

import (
	"time"
)

// Now gets the timestamp of the number of milliseconds that have elapsed since the Unix epoch.
//
// Example:
//
//	Now() // 1640995200000 (example timestamp)
func Now() int64 {
	return time.Now().UnixMilli()
}

// ToDate converts value to a Date.
//
// Example:
//
//	ToDate(1640995200000) // time.Time corresponding to the timestamp
//	ToDate("2022-01-01T00:00:00Z") // parsed time.Time
func ToDate(value any) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case int64:
		return time.UnixMilli(v), nil
	case int:
		return time.UnixMilli(int64(v)), nil
	case string:
		// Try common date formats
		formats := []string{
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02T15:04:05Z",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, v); err == nil {
				return t, nil
			}
		}
		return time.Time{}, &time.ParseError{Layout: "various formats", Value: v}
	default:
		return time.Time{}, &time.ParseError{Layout: "unsupported type", Value: ""}
	}
}

// IsDate checks if value is classified as a Date object.
//
// Example:
//
//	IsDate(time.Now()) // true
//	IsDate("2022-01-01") // false
func IsDate(value any) bool {
	_, ok := value.(time.Time)
	return ok
}

// IsValid checks if the given time is valid (not zero time).
//
// Example:
//
//	IsValid(time.Now()) // true
//	IsValid(time.Time{}) // false
func IsValid(t time.Time) bool {
	return !t.IsZero()
}

// Format formats a time according to the given layout.
//
// Example:
//
//	Format(time.Now(), "2006-01-02") // "2022-01-01"
//	Format(time.Now(), "15:04:05") // "12:30:45"
func Format(t time.Time, layout string) string {
	return t.Format(layout)
}

// Add adds the duration to the time.
//
// Example:
//
//	Add(time.Now(), time.Hour) // time one hour from now
//	Add(time.Now(), -time.Hour) // time one hour ago
func Add(t time.Time, d time.Duration) time.Time {
	return t.Add(d)
}

// Sub returns the duration t-u.
//
// Example:
//
//	Sub(time.Now(), time.Now().Add(-time.Hour)) // time.Hour
func Sub(t, u time.Time) time.Duration {
	return t.Sub(u)
}

// Before reports whether the time instant t is before u.
//
// Example:
//
//	Before(time.Now(), time.Now().Add(time.Hour)) // true
func Before(t, u time.Time) bool {
	return t.Before(u)
}

// After reports whether the time instant t is after u.
//
// Example:
//
//	After(time.Now(), time.Now().Add(-time.Hour)) // true
func After(t, u time.Time) bool {
	return t.After(u)
}

// Equal reports whether t and u represent the same time instant.
//
// Example:
//
//	Equal(time.Now(), time.Now()) // false (different nanoseconds)
//	Equal(t, t) // true
func Equal(t, u time.Time) bool {
	return t.Equal(u)
}

// StartOfDay returns the start of the day for the given time.
//
// Example:
//
//	StartOfDay(time.Date(2022, 1, 1, 15, 30, 45, 0, time.UTC)) // 2022-01-01 00:00:00
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay returns the end of the day for the given time.
//
// Example:
//
//	EndOfDay(time.Date(2022, 1, 1, 15, 30, 45, 0, time.UTC)) // 2022-01-01 23:59:59.999999999
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// StartOfWeek returns the start of the week (Monday) for the given time.
//
// Example:
//
//	StartOfWeek(time.Date(2022, 1, 5, 15, 30, 45, 0, time.UTC)) // 2022-01-03 00:00:00 (Monday)
func StartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	days := weekday - 1 // Days since Monday
	return StartOfDay(t.AddDate(0, 0, -days))
}

// EndOfWeek returns the end of the week (Sunday) for the given time.
//
// Example:
//
//	EndOfWeek(time.Date(2022, 1, 5, 15, 30, 45, 0, time.UTC)) // 2022-01-09 23:59:59.999999999 (Sunday)
func EndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	days := 7 - weekday // Days until Sunday
	return EndOfDay(t.AddDate(0, 0, days))
}

// StartOfMonth returns the start of the month for the given time.
//
// Example:
//
//	StartOfMonth(time.Date(2022, 1, 15, 15, 30, 45, 0, time.UTC)) // 2022-01-01 00:00:00
func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// EndOfMonth returns the end of the month for the given time.
//
// Example:
//
//	EndOfMonth(time.Date(2022, 1, 15, 15, 30, 45, 0, time.UTC)) // 2022-01-31 23:59:59.999999999
func EndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location())
	return EndOfDay(firstOfNextMonth.AddDate(0, 0, -1))
}

// StartOfYear returns the start of the year for the given time.
//
// Example:
//
//	StartOfYear(time.Date(2022, 6, 15, 15, 30, 45, 0, time.UTC)) // 2022-01-01 00:00:00
func StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

// EndOfYear returns the end of the year for the given time.
//
// Example:
//
//	EndOfYear(time.Date(2022, 6, 15, 15, 30, 45, 0, time.UTC)) // 2022-12-31 23:59:59.999999999
func EndOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return EndOfDay(time.Date(year, 12, 31, 0, 0, 0, 0, t.Location()))
}

// DaysInMonth returns the number of days in the month for the given time.
//
// Example:
//
//	DaysInMonth(time.Date(2022, 2, 15, 0, 0, 0, 0, time.UTC)) // 28
//	DaysInMonth(time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC)) // 29 (leap year)
func DaysInMonth(t time.Time) int {
	year, month, _ := t.Date()
	firstOfNextMonth := time.Date(year, month+1, 1, 0, 0, 0, 0, t.Location())
	lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
	return lastOfMonth.Day()
}

// IsLeapYear checks if the given year is a leap year.
//
// Example:
//
//	IsLeapYear(2020) // true
//	IsLeapYear(2021) // false
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
