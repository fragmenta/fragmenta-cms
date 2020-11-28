package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// Deprecated - use server/log pkg instead to log

// Logger interface for a logger - deprecated for 2.0
type Logger interface {
	Printf(format string, args ...interface{})
}

// Logf logs the message with the given arguments to our internal logger
func (s *Server) Logf(format string, v ...interface{}) {
	s.Logger.Printf(format, v...)
}

// Log logs the message to our internal logger
func (s *Server) Log(message string) {
	s.Logf(message)
}

// Fatalf the message with the given arguments to our internal logger, and then exits with status 1
func (s *Server) Fatalf(format string, v ...interface{}) {
	s.Logger.Printf(format, v...)

	// Now exit
	os.Exit(1)
}

// Fatal logs the message, and then exits with status 1
func (s *Server) Fatal(format string) {
	s.Fatalf(format)
}

// Timef logs a time since starting, when used with defer at the start of a function to time
// Usage: defer s.Timef("Completed %s in %s",time.Now(),args...)
func (s *Server) Timef(format string, start time.Time, v ...interface{}) {
	end := time.Since(start).String()
	var args []interface{}
	args = append(args, end)
	args = append(args, v...)
	s.Logf(format, args...)
}

// Deprecated - this config parsing (mostly internal to the server anyway)
// has been moved to a new file and will be removed in 2.0
// Instead of passing config to setup and thence to handlers via router context,
// apps/handlers should use server/config to access it if required.

// Mode returns the mode (production or development)
func (s *Server) Mode() string {
	if s.production {
		return "Production"
	}
	return "Development"
}

// SetProduction sets the mode manually to SetProduction
func (s *Server) SetProduction(value bool) {
	s.production = value
}

// Production tells the caller if this server is in production mode or not?
func (s *Server) Production() bool {
	return s.production
}

// Configuration returns the map of configuration keys to values
func (s *Server) Configuration() map[string]string {
	if s.production {
		return s.configProduction
	}
	return s.configDevelopment

}

// Config returns a specific configuration value or "" if no value
func (s *Server) Config(key string) string {
	return s.Configuration()[key]
}

// ConfigInt returns the current configuration value as int64, or 0 if no value
func (s *Server) ConfigInt(key string) int64 {
	v := s.Config(key)
	if v != "" {
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i
		}
	}
	return 0
}

// ConfigBool returns the current configuration value as bool (yes=true, no=false), or false if no value
func (s *Server) ConfigBool(key string) bool {
	v := s.Config(key)
	return (v == "yes")
}

// configPath returns our expected config file path
func (s *Server) configPath() string {
	return "secrets/fragmenta.json"
}

// Read our config file and set up the server accordingly
func (s *Server) readConfig() error {

	path := s.configPath()

	// Read the config json file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error opening config %s %v", path, err)
	}

	var data map[string]map[string]string
	err = json.Unmarshal(file, &data)
	if err != nil {
		return fmt.Errorf("Error reading config %s %v", path, err)
	}

	s.configDevelopment = data["development"]
	s.configProduction = data["production"]
	s.configTest = data["test"]

	// Update our port from the config port if we have it
	portString := s.Config("port")
	if portString != "" {
		s.port, err = strconv.Atoi(portString)
		if err != nil {
			return fmt.Errorf("Error reading port %s", err)
		}
	}

	return nil
}

// Deprecated - the server relies on config for lots of settings
// the port can be changed in config instead for development easily.
// This flag will be removed in 2.0

// readArguments reads command line arguments
func (s *Server) readArguments() error {

	var p int
	flag.IntVar(&p, "p", p, "Port")
	flag.Parse()

	if p > 0 {
		s.port = p
	}

	return nil
}
