// Package config offers utilities for parsing a json config file.
// Values are read as strings, and can be fetched with Get, GetInt or GetBool.
// The caller is expected to parse them for more complex types.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

const (
	// DefaultPath is where our config is normally found for fragmenta apps.
	DefaultPath = "secrets/fragmenta.json"
)

// Config modes are set when creating a new config
const (
	ModeDevelopment = iota
	ModeProduction
	ModeTest
)

// Current is the current configuration object for
var Current *Config

// Config represents a set of key/value pairs for each mode of the app,
// production, development and test. Which set of values is used
// is set by Mode.
type Config struct {
	Mode    int
	configs []map[string]string
}

// New returns a new config, which defaults to development
func New() *Config {
	return &Config{
		Mode:    ModeDevelopment,
		configs: make([]map[string]string, 3),
	}
}

// Load our json config file from the path
func (c *Config) Load(path string) error {

	// Read the config json file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error opening config %s %v", path, err)
	}

	var data map[string]map[string]string
	err = json.Unmarshal(file, &data)
	if err != nil {
		return fmt.Errorf("error reading config %s %v", path, err)
	}

	if len(data) < 3 {
		return fmt.Errorf("error reading config - not enough configs, got :%d expected 3", len(data))
	}

	c.configs[ModeDevelopment] = data["development"]
	c.configs[ModeProduction] = data["production"]
	c.configs[ModeTest] = data["test"]

	return nil
}

// Production returns true if current config is production.
func (c *Config) Production() bool {
	return c.Mode == ModeProduction
}

// Configuration returns all the configuration key/values for a given mode.
func (c *Config) Configuration(m int) map[string]string {
	return c.configs[c.Mode]
}

// Get returns a specific value or "" if no value
func (c *Config) Get(key string) string {
	if c == nil {
		return ""
	}
	return c.configs[c.Mode][key]
}

// GetInt returns the current configuration value as int64, or 0 if no value
func (c *Config) GetInt(key string) int64 {
	v := c.Get(key)
	if v != "" {
		i, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return i
		}
	}
	return 0
}

// GetBool returns the current configuration value as bool
// (yes=true, no=false), or false if no value
func (c *Config) GetBool(key string) bool {
	v := c.Get(key)
	return (v == "yes")
}

// Config (Get) returns a specific value or "" if no value
// For compatability with older server config, we wrap this function
// Deprecated
func (c *Config) Config(key string) string {
	return c.Get(key)
}

// These convenience functions wrap the Current pkg global

// Production returns true if current config is production.
func Production() bool {
	return Current.Production()
}

// Configuration returns all the configuration key/values for a given mode.
func Configuration(m int) map[string]string {
	return Current.Configuration(m)
}

// Get returns a specific value or "" if no value
func Get(key string) string {
	return Current.Get(key)
}

// GetInt returns the current configuration value as int64, or 0 if no value
func GetInt(key string) int64 {
	return Current.GetInt(key)
}

// GetBool returns the current configuration value as bool
// (yes=true, no=false), or false if no value
func GetBool(key string) bool {
	return Current.GetBool(key)
}
