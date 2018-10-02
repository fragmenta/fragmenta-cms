package adapters

// FIXME: Sqlite drivers are broken compiling with cgo at present
// therefore we don't use this adapter

import (
	"database/sql"
	"fmt"
	// Unfortunately can't cross compile with sqlite support enabled -
	// see https://github.com/mattn/go-sqlite3/issues/106
	// For now for we just turn off sqlite as we don't use it in production...
	// pure go version of sqlite, or ditch sqlite and find some other pure go simple db
	// would be nice not to require a db at all for very simple usage
	//_ "github.com/mattn/go-sqlite3"
)

// SqliteAdapter conforms to the query.Database interface
type SqliteAdapter struct {
	*Adapter
	options map[string]string
	sqlDB   *sql.DB
	debug   bool
}

// Open this database
func (db *SqliteAdapter) Open(opts map[string]string) error {

	db.debug = false
	db.options = map[string]string{
		"adapter": "sqlite3",
		"db":      "./tests/query_test.sqlite",
	}

	if opts["debug"] == "true" {
		db.debug = true
	}

	for k, v := range opts {
		db.options[k] = v
	}

	var err error
	db.sqlDB, err = sql.Open(db.options["adapter"], db.options["db"])
	if err != nil {
		return err
	}

	if db.sqlDB != nil && db.debug {
		fmt.Printf("Database %s opened using %s\n", db.options["db"], db.options["adapter"])
	}

	// Call ping on the db to check it does actually exist!
	err = db.sqlDB.Ping()
	if err != nil {
		return err
	}

	return err

}

// Close the database
func (db *SqliteAdapter) Close() error {
	if db.sqlDB != nil {
		return db.sqlDB.Close()
	}
	return nil
}

// SQLDB returns the internal db.sqlDB pointer
func (db *SqliteAdapter) SQLDB() *sql.DB {
	return db.sqlDB
}

// Query execute Query SQL - NB caller must call use defer rows.Close() with rows returned
func (db *SqliteAdapter) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.performQuery(db.sqlDB, db.debug, query, args...)
}

// Exec - use this for non-select statements
func (db *SqliteAdapter) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.performExec(db.sqlDB, db.debug, query, args...)
}

// Insert a record with params and return the id - psql behaves differently
func (db *SqliteAdapter) Insert(query string, args ...interface{}) (id int64, err error) {

	// Execute the sql using db
	result, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil

}
