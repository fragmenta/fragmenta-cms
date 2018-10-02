// Package log provides a structured, levelled logger interface
// for use in server handlers, which handles multiple output streams.
// A typical use might be to log everything to stderr, but to add another
// logger to send important data off to
// The Default logger simply logs to stderr, a local File logger is available,
// and data can be extracted and sent elsewhere by additional loggers
// (for example page hits to a stats service).
//
// Usage:
// logger,err := log.NewStdErr()
// log.Add(logger)
// log.Error(log.V{"key":value,"key":value})
//
package log

import (
	"os"
	"time"
)

const (
	// LevelKey is the key for setting level
	LevelKey = "level"
	// MessageKey is the key for a message
	MessageKey = "msg"
	// DurationKey is used by the Time function
	DurationKey = "duration"
	// ErrorKey is used for errors
	ErrorKey = "error"
	// IPKey is used for IP addresses (for colouring)
	IPKey = "ip"
	// URLKey is used for identifying URLs (for filtering)
	URLKey = "url"
	// TraceKey is used for trace ids emitted in middleware
	TraceKey = "trace"
)

// Debug sends the key/value map at level Debug to all registered (log)gers.
func Debug(values map[string]interface{}) {
	values[LevelKey] = LevelDebug
	Log(values)
}

// Info sends the key/value map at level Info to all registered loggers.
func Info(values map[string]interface{}) {
	values[LevelKey] = LevelInfo
	Log(values)
}

// Error sends the key/value map at level Error to all registered loggers.
func Error(values map[string]interface{}) {
	values[LevelKey] = LevelError
	Log(values)
}

// Fatal sends the key/value map at level Fatal to all registered loggers,
// no other action is taken.
func Fatal(values map[string]interface{}) {
	values[LevelKey] = LevelFatal
	Log(values)
}

// Time sends the key/value map to all registered loggers with an additional duration, start and end params set.
func Time(start time.Time, values map[string]interface{}) {
	values[DurationKey] = time.Now().UTC().Sub(start)
	Log(values)
}

// Log sends the key/value map to all registered loggers. If level is not set,
// it defaults to LevelInfo.
func Log(values map[string]interface{}) {
	_, ok := values[LevelKey]
	if !ok {
		values[LevelKey] = LevelInfo
	}

	for _, l := range loggers {
		l.Log(values)
	}
}

// Add adds the given logger to the list of outputs,
// it should not be called from other goroutines.
func Add(l StructuredLogger) {
	loggers = append(loggers, l)
}

// Valid levels for logging.
const (
	LevelNone = iota
	LevelDebug
	LevelInfo
	LevelError
	LevelFatal
)

var (
	// LevelNames is a list of human-readable for levels.
	LevelNames = []string{"none", "debug", "info", "error", "fatal"}

	// NoColor determines if a terminal is colourable or not
	NoColor = os.Getenv("TERM") == "dumb"

	// LevelColors is a list of human-friendly terminal colors for levels.
	LevelColors = []string{"\033[0m", "\033[34m", "\033[32m", "\033[33m", "\033[31m"}

	// TraceColor sets a for IP addresses or request id
	TraceColor = "\033[33m"

	// ClearColors clears all formatting
	ClearColors = "\033[0m"
)

// This variable stores multiple loggers, which may decide whether
// to print or not depending on the level and/or message content.
// They may log to a file, stderr, or over the network, and different
// destinations may all log the same messages.
var loggers []StructuredLogger

// StructuredLogger defines an interface for loggers
// which may be added with Add() to the list of outputs.
type StructuredLogger interface {
	Log(V)
}

// Values is a map of structured key value pairs
// usage: log.Warn(log.Values{"user":1,"foo":"bar"})
type Values map[string]interface{}

// V is a shorthand for values
type V map[string]interface{}
