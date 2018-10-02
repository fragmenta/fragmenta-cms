package log

import (
	"errors"
	"os"
)

// File logs to a local file for all messages at or above Level.
type File struct {
	Default // File embeds default
	Path    string
}

const (
	// FileFlags serts the flags for OpenFile on the log file
	FileFlags = os.O_WRONLY | os.O_APPEND | os.O_CREATE

	// FilePermissions serts the perms for OpenFile on the log file
	FilePermissions = 0640
)

// NewFile creates a new file logger for the given path at Level Info.
func NewFile(path string) (*File, error) {
	if path == "" {
		return nil, errors.New("log: null file path for file log")
	}

	f := &File{
		Default: Default{
			Prefix: PrefixDateTime,
			Level:  LevelInfo,
			Writer: nil,
			Color:  false,
		},
	}

	// Set the writer to the given file
	logFile, err := os.OpenFile(path, FileFlags, FilePermissions)
	if err != nil {
		return nil, err
	}

	f.Writer = logFile

	return f, nil
}
