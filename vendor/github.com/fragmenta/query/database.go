package query

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fragmenta/query/adapters"
)

// database is the package global db  - this reference is not exported outside the package.
var database adapters.Database

// OpenDatabase opens the database with the given options
func OpenDatabase(opts map[string]string) error {

	// If we already have a db, return it
	if database != nil {
		return fmt.Errorf("query: database already open - %s", database)
	}

	// Assign the db global in query package
	switch opts["adapter"] {
	case "sqlite3":
		database = &adapters.SqliteAdapter{}
	case "mysql":
		database = &adapters.MysqlAdapter{}
	case "postgres":
		database = &adapters.PostgresqlAdapter{}
	default:
		database = nil // fail
	}

	if database == nil {
		return fmt.Errorf("query: database adapter not recognised - %s", opts)
	}

	// Ask the db adapter to open
	return database.Open(opts)
}

// CloseDatabase closes the database opened by OpenDatabase
func CloseDatabase() error {
	var err error
	if database != nil {
		err = database.Close()
		database = nil
	}

	return err
}

// SetMaxOpenConns sets the maximum number of open connections
func SetMaxOpenConns(max int) {
	database.SQLDB().SetMaxOpenConns(max)
}

// QuerySQL executes the given sql Query against our database, with arbitrary args
func QuerySQL(query string, args ...interface{}) (*sql.Rows, error) {
	if database == nil {
		return nil, fmt.Errorf("query: QuerySQL called with nil database")
	}
	results, err := database.Query(query, args...)
	return results, err
}

// ExecSQL executes the given sql against our database with arbitrary args
// NB returns sql.Result - not to be used when rows expected
func ExecSQL(query string, args ...interface{}) (sql.Result, error) {
	if database == nil {
		return nil, fmt.Errorf("query: ExecSQL called with nil database")
	}
	results, err := database.Exec(query, args...)
	return results, err
}

// TimeString returns a string formatted as a time for this db
// if the database is nil, an empty string is returned.
func TimeString(t time.Time) string {
	if database != nil {
		return database.TimeString(t)
	}
	return ""
}
