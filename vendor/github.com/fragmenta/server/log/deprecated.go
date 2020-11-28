package log

// This file contains MultiWriter logging to a file and stdout,
// it has now been deprecated in favour of structured logging found in log.go
// It will be removed in Fragmenta 2.0

import (
	"io"
	stdlog "log"
	"os"
	"strings"
)

// We accept hashtags like #error in log messages, which can be used to filter messages
// These tags can also be used to indicate the level of the message
// NB if filter is set we only output messages containg filter string

// Logger conforms with the server.Logger interface
type Logger struct {
	log    *stdlog.Logger
	Filter string
}

// New creates a new Logger which writes to a file and to stderr
func New(path string, production bool) *Logger {
	var logWriter io.Writer
	stdlog.SetFlags(stdlog.Llongfile)
	// doubleWriter writes to stdErr and to a file
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0640)
	if err != nil {
		logWriter = io.MultiWriter(os.Stderr)
	} else {

		// Do not write to Stderr in production
		if production {
			logWriter = io.MultiWriter(logFile)
		} else {
			logWriter = io.MultiWriter(os.Stderr, logFile)
		}

	}

	// By default logger logs to console and a file
	l := stdlog.New(logWriter, "", stdlog.Ldate|stdlog.Ltime)
	if l == nil {
		stdlog.Printf("Error setting up log at path %s", path)
	}

	logger := &Logger{
		log:    l,
		Filter: "",
	}

	logger.Printf("#info Opened log file at %s", path)
	return logger
}

// Printf logs events selectively given our filter
func (l *Logger) Printf(format string, args ...interface{}) {

	if l.Filter == "" {
		// If we have no filter, print all
		l.writeLog(format, args...)
	} else if strings.Contains(format, l.Filter) {
		// if we have a filter, print only those messages which match it (e.g. only match #error in production)
		l.writeLog(format, args...)
	}

}

// Log events to the server log file and other output
func (l *Logger) writeLog(format string, args ...interface{}) {

	if l.log != nil {
		if strings.Contains(format, "%") {
			l.log.Printf(format, args...)
		} else {
			l.log.Print(format)
		}
	} else {
		// If we failed to create a log, just log something to stdout
		stdlog.Printf(format, args...)
	}

}
