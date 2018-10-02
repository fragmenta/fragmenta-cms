package log

import (
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

const (
	// Separator is used to separate key:value pairs
	Separator = ":"
	// PrefixDate constant for date prefixes
	PrefixDate = "2006-01-02 "
	// PrefixTime  constants for time prefix
	PrefixTime = "15:04:05 "
	// PrefixDateTime  constants for date + time prefix
	PrefixDateTime = "2006-01-02:15:04:05 "
)

// NewStdErr returns a new StructuredLogger of type Default which writes to stderr. By default
// prefix is empty, and level is LevelDebug (the lowest),
// so all output is captured.
func NewStdErr(prefix string) (*Default, error) {
	d := &Default{
		Prefix: prefix, // Treated as a time format if set
		Level:  LevelInfo,
		Writer: os.Stderr,
		Color:  true,
	}
	return d, nil
}

// Default defines a default logger which simply logs to Writer,
// Writer is set to stderr, Level is LevelDebug and Prefix is empty by default.
type Default struct {

	// Prefix is used to prefix any log lines emitted.
	Prefix string

	// Level is the level above which input is ignored.
	Level int

	// Writer is the output of this logger.
	Writer io.Writer

	// Color sets whether terminal colour instructions are emitted.
	Color bool
}

// Log logs the key:value pairs given to the writer. Keys are sorted before
// output in alphabetical order to ensure consistent results.
func (d *Default) Log(values V) {
	l := d.LevelValue(values)
	if l < d.Level {
		return
	}

	// Start by writing the prefix (treated as a time format string)
	d.WriteString(time.Now().UTC().Format(d.Prefix))

	// If keys contains message, extract that first
	msg, ok := values[MessageKey].(string)
	if ok {
		d.WriteString(msg + " ")
	}
	// If keys contains duration, extract that next
	duration, ok := values[DurationKey].(time.Duration)
	if ok {
		d.WriteString("in " + duration.String() + " ")
	}

	// Now print other keys with colouring
	var prefix, suffix string
	keys := d.SortedKeys(values)
	for _, k := range keys {
		d.WriteString(k)
		d.WriteString(Separator)

		switch k {
		case IPKey:
			fallthrough
		case TraceKey:
			d.WriteString(fmt.Sprintf("%s%v%s ", TraceColor, values[k], ClearColors))
		default:
			d.WriteString(fmt.Sprintf("%v ", values[k]))
		}

	}

	if d.Color {
		prefix = d.LevelColor(l)
		suffix = ClearColors
	}
	d.WriteString(fmt.Sprintf("%s#%v%s ", prefix, d.LevelName(l), suffix))

	d.WriteString("\n")

}

// WriteString writes the string to the Writer.
func (d *Default) WriteString(s string) {
	d.Writer.Write([]byte(s))
}

// LevelValue extracts the Level from values (if present) or returns 0 if not.
func (d *Default) LevelValue(values V) int {
	l, ok := values[LevelKey].(int)
	if ok {
		return l
	}
	return 0
}

// LevelName returns the human-readable name for this level.
func (d *Default) LevelName(l int) string {
	return LevelNames[l]
}

// LevelColor returns the human-readable colour for this level.
func (d *Default) LevelColor(l int) string {
	return LevelColors[l]
}

// SortedKeys returns an array of keys for a map sorted in alpha order,
// this means we get a predictable order for the map entries when we print.
// The special keys level and message are ommitted.
func (d *Default) SortedKeys(values V) []string {
	var keys []string
	for k := range values {
		// Ignore these special keys
		if k == DurationKey || k == MessageKey || k == LevelKey {
			continue
		}
		// Append the sorted key
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
