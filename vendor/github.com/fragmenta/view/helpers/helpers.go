package helpers

import (
	"fmt"
	got "html/template"
	"strings"
	"time"
)

// ARRAYS

// Array takes a set of interface pointers as variadic args, and returns a single array
func Array(args ...interface{}) []interface{} {
	return []interface{}{args}
}

// CommaSeparatedArray returns the values as a comma separated string
func CommaSeparatedArray(args []string) string {
	result := ""
	for _, v := range args {
		if len(result) > 0 {
			result = fmt.Sprintf("%s,%s", result, v)
		} else {
			result = v
		}

	}
	return result
}

// MAPS

// Empty returns an empty map[string]interface{} for use as a context
func Empty() map[string]interface{} {
	return map[string]interface{}{}
}

// Map sets a map key and return the map
func Map(m map[string]interface{}, k string, v interface{}) map[string]interface{} {
	m[k] = v
	return m
}

// Set a map key and return an empty string
func Set(m map[string]interface{}, k string, v interface{}) string {
	m[k] = v
	return "" // Render nothing, we want no side effects
}

// SetIf sets a map key if the given condition is true
func SetIf(m map[string]interface{}, k string, v interface{}, t bool) string {
	if t {
		m[k] = v
	} else {
		m[k] = ""
	}
	return "" // Render nothing, we want no side effects
}

// Append all args to an array, and return that array
func Append(m []interface{}, args ...interface{}) []interface{} {
	for _, v := range args {
		m = append(m, v)
	}
	return m
}

// CreateMap - given a set of interface pointers as variadic args, generate and return a map to the values
// This is currently unused as we just use simpler Map add above to add to context
func CreateMap(args ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, 0)

	key := ""
	for _, v := range args {
		if len(key) == 0 {
			key = string(v.(string))
		} else {
			m[key] = v
		}
	}

	return m
}

// Contains returns true if this array of ints contains the given int
func Contains(list []int64, item int64) bool {
	for _, b := range list {
		if b == item {
			return true
		}
	}
	return false
}

// Blank returns true if a string is empty
func Blank(s string) bool {
	return len(s) == 0
}

// Exists returns true if this string has a length greater than 0
func Exists(s string) bool {
	return len(s) > 0
}

// Time returns a formatted time string given a time and optional format
func Time(time time.Time, formats ...string) got.HTML {
	layout := "Jan 2, 2006 at 15:04"
	if len(formats) > 0 {
		layout = formats[0]
	}
	value := fmt.Sprintf(time.Format(layout))
	return got.HTML(Escape(value))
}

// Ago returns a time string reporting distance from the current date
// of form 5 hours ago
func Ago(t time.Time, formats ...string) string {
	duration := time.Since(t)
	absDuration := duration
	if duration < 0 {
		absDuration = -duration
	}
	hours := absDuration / time.Hour

	// Use ago only for past dates
	ago := " ago"
	if duration < 0 {
		ago = ""
	}

	switch {
	case absDuration < time.Minute:
		return fmt.Sprintf("%d seconds%s", duration/time.Second, ago)
	case absDuration < time.Hour:
		return fmt.Sprintf("%d minutes%s", duration/time.Minute, ago)
	case absDuration < time.Hour*24:
		unit := "hour"
		if hours > 1 {
			unit = "hours"
		}
		return fmt.Sprintf("%d %s%s", hours, unit, ago)
	default:
		unit := "day"
		if hours > 48 {
			unit = "days"
		}
		return fmt.Sprintf("%d %s%s", hours/24, unit, ago)
	}
}

// Date returns a formatted date string given a time and optional format
// Date format layouts are for the date 2006-01-02
func Date(t time.Time, formats ...string) got.HTML {

	//layout := "2006-01-02" // Jan 2, 2006
	layout := "Jan 2, 2006"
	if len(formats) > 0 {
		layout = formats[0]
	}
	value := fmt.Sprintf(t.Format(layout))
	return got.HTML(Escape(value))
}

// UTCDate returns a formatted date string in 2006-01-02
func UTCDate(t time.Time) got.HTML {
	return Date(t.UTC(), "2006-01-02")
}

// UTCTime returns a formatted date string in 2006-01-02
func UTCTime(t time.Time) got.HTML {
	return Time(t.UTC(), "2006-01-02T15:04:00:00.000Z")
}

// JSONTime returns a formatted date string with format
// suitable for using in a json file
func JSONTime(t time.Time) got.HTML {
	return Time(t.UTC(), "2006-01-02T15:04:05Z07:00")
}

// UTCNow returns a formatted date string in 2006-01-02
func UTCNow() got.HTML {
	return Date(time.Now().UTC(), "2006-01-02")
}

// YearNow returns a formatted date string for the current year
func YearNow() got.HTML {
	return Date(time.Now().UTC(), "2006")
}

// Truncate text to a given length
func Truncate(s string, l int64) string {
	return s
}

// CSV escape (replace , with ,,)
func CSV(s got.HTML) string {
	return strings.Replace(string(s), ",", ",,", -1)
}

// JSON escapes a string for use in a json template (html template)
func JSON(t string) got.HTML {
	// Escape mandatory characters
	t = strings.Replace(t, "\r", " ", -1)
	t = strings.Replace(t, "\n", " ", -1)
	t = strings.Replace(t, "\t", " ", -1)
	t = strings.Replace(t, "\\", "\\\\", -1)
	t = strings.Replace(t, "\"", "\\\"", -1)
	// Because we use html/template escape as temlate.HTML
	return got.HTML(t)
}
