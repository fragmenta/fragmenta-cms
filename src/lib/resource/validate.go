package resource

import (
	"time"
)

// Methods for validating params passed from the database row as interface{} types
// perhaps this should be a sub-package for clarity?

// ValidateFloat returns the float value of param or 0.0
func ValidateFloat(param interface{}) float64 {
	var v float64
	if param != nil {
		switch param.(type) {
		case float64:
			v = param.(float64)
		case float32:
			v = float64(param.(float32))
		case int:
			v = float64(param.(int))
		case int64:
			v = float64(param.(int64))
		}
	}
	return v
}

// ValidateBoolean returns the bool value of param or false
func ValidateBoolean(param interface{}) bool {
	var v bool
	if param != nil {
		switch param.(type) {
		case bool:
			v = param.(bool)
		}
	}
	return v
}

// ValidateInt returns the int value of param or 0
func ValidateInt(param interface{}) int64 {
	var v int64
	if param != nil {
		switch param.(type) {
		case int64:
			v = param.(int64)
		case float64:
			v = int64(param.(float64))
		case int:
			v = int64(param.(int))
		}
	}
	return v
}

// ValidateString returns the string value of param or ""
func ValidateString(param interface{}) string {
	var v string
	if param != nil {
		switch param.(type) {
		case string:
			v = param.(string)
		}
	}
	return v
}

// ValidateTime returns the time value of param or the zero value of time.Time
func ValidateTime(param interface{}) time.Time {
	var v time.Time
	if param != nil {
		switch param.(type) {
		case time.Time:
			v = param.(time.Time)
		}
	}
	return v
}
