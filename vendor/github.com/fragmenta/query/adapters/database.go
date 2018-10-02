// Package adapters offers adapters for pouplar databases
package adapters

import (
	"database/sql"
	"fmt"
	"time"
)

// Database provides an interface for database adapters to conform to
type Database interface {

	// Open and close
	Open(opts map[string]string) error
	Close() error
	SQLDB() *sql.DB

	// Execute queries with or without returned rows
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)

	// Insert a record, returning id
	Insert(sql string, args ...interface{}) (id int64, err error)

	// Return extra SQL for insert statement (see psql)
	InsertSQL(pk string) string

	// A format string for the arg placeholder
	Placeholder(i int) string

	// Quote Table and Column names
	QuoteField(name string) string

	// Convert a time to a string
	TimeString(t time.Time) string

	// Convert a string to a time
	ParseTime(s string) (time.Time, error)
}

// Adapter is a struct defining a few functions used by all adapters
type Adapter struct {
	queries map[string]interface{}
}

// ReplaceArgPlaceholder does no replacements by default, and use default ? placeholder for args
// psql requires a different placeholder numericall labelled
func (db *Adapter) ReplaceArgPlaceholder(sql string, args []interface{}) string {
	return sql
}

// Placeholder is the argument placeholder for this adapter
func (db *Adapter) Placeholder(i int) string {
	return "?"
}

// TimeString - given a time, return the standard string representation
func (db *Adapter) TimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000 -0700")
}

// ParseTime - given a string, return a time object built from it
func (db *Adapter) ParseTime(s string) (time.Time, error) {

	// Deal with broken mysql dates - deal better with this?
	if s == "0000-00-00 00:00:00" {
		return time.Now(), nil
	}

	// Try to choose the right format for date string
	format := "2006-01-02 15:04:05"
	if len(s) > len(format) {
		format = "2006-01-02 15:04:05.000"
	}
	if len(s) > len(format) {
		format = "2006-01-02 15:04:05.000 -0700"
	}

	t, err := time.Parse(format, s)
	if err != nil {
		fmt.Println("Unhandled field type:", s, "\n", err)
	}

	return t, err
}

// QuoteField quotes a table name or column name
func (db *Adapter) QuoteField(name string) string {
	return fmt.Sprintf(`"%s"`, name)
}

// InsertSQL provides extra SQL for end of insert statement (RETURNING for psql)
func (db *Adapter) InsertSQL(pk string) string {
	return ""
}

// performQuery executes Query SQL on the given sqlDB and return the rows.
// NB caller must call use defer rows.Close() with rows returned
func (db *Adapter) performQuery(sqlDB *sql.DB, debug bool, query string, args ...interface{}) (*sql.Rows, error) {

	if sqlDB == nil {
		return nil, fmt.Errorf("No database available.")
	}

	if debug {
		fmt.Println("QUERY:", query, "ARGS", args)
	}

	// This should be cached, perhaps hold a map in memory of queries strings and compiled queries?
	// use queries map to store this
	stmt, err := sqlDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)

	if err != nil {
		return nil, err
	}

	// Caller is responsible for closing rows with defer rows.Close()
	return rows, err
}

// performExec executes Query SQL on the given sqlDB with no rows returned, just result
func (db *Adapter) performExec(sqlDB *sql.DB, debug bool, query string, args ...interface{}) (sql.Result, error) {

	if sqlDB == nil {
		return nil, fmt.Errorf("No database available.")
	}

	if debug {
		fmt.Println("QUERY:", query, "ARGS", args)
	}

	stmt, err := sqlDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)

	if err != nil {
		return result, err
	}

	// Caller is responsible for closing rows with defer rows.Close()
	return result, err
}
