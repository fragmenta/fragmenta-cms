Query [![GoDoc](https://godoc.org/github.com/fragmenta/query?status.svg)](https://godoc.org/github.com/fragmenta/query) [![Go Report Card](https://goreportcard.com/badge/github.com/fragmenta/query)](https://goreportcard.com/report/github.com/fragmenta/query)
=====



Query lets you build SQL queries with chainable methods, and defer execution of SQL until you wish to extract a count or array of models. It will probably remain limited in scope - it is not intended to be a full ORM with strict mapping between db tables and structs, but a tool for querying the database with minimum friction, and performing CRUD operations linked to models; simplifying your use of SQL to store model data without getting in the way. Full or partial SQL queries are of course also available, and full control over sql. Model creation and column are delegated to the model, to avoid dictating any particular model structure or interface, however a suggested interface is given (see below and in tests), which makes usage painless in your handlers without any boilerplate.

Supported databases: PostgreSQL, SQLite, MySQL. Bug fixes, suggestions and contributions welcome. 

Usage
=====


```go

// In your app - open a database with options
options := map[string]string{"adapter":"postgres","db":"query_test"}
err := query.OpenDatabase(options)
defer query.CloseDatabase()

...

// In your model
type Page struct {
	ID			int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	MyField	    myStruct
	...
	// Models can have any structure, any PK, here an int is used
}

// Normally you'd define helpers on your model class to load rows from the database
// Query does not attempt to read data into columns with reflection or tags - 
// that is left to your model so you can read as little or as much as you want from queries

func Find(ID int64) (*Page, error) {
	result, err := PagesQuery().Where("id=?", ID).FirstResult()
	if err != nil {
		return nil, err
	}
	return NewWithColumns(result), nil
}

func FindAll(q *Query) ([]*Page, error) {
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	var models []*Page
	for _, r := range results {
		m := NewWithColumns(r)
		models = append(models, m)
	}

	return models, nil
}

...

// In your handlers, construct queries and ask your models for the data

// Find a simple model by id
page, err := pages.Find(1)

// Start querying the database using chained finders
q := page.Query().Where("id IN (?,?)",4,5).Order("id desc").Limit(30)

// Build up chains depending on other app logic, still no db requests
if shouldRestrict {
	q.Where("id > ?",3).OrWhere("keywords ~* ?","Page")
}

// Pass the relation around, until you are ready to retrieve models from the db
results, err := pages.FindAll(q)
```

What it does
============

* Builds chainable queries including where, orwhere,group,having,order,limit,offset or plain sql
* Allows any Primary Key/Table name or model fields (query.New lets you define this)
* Allows Delete and Update operations on queried records, without creating objects
* Defers SQL requests until full query is built and results requested
* Provide helpers and return results for join ids, counts, single rows, or multiple rows


What it doesn't do
==================

* Attempt to read your models with reflection or struct tags
* Require changes to your structs like tagging fields or specific fields
* Cause problems with untagged fields, embedding, and fields not in the database
* Provide hooks after/before update etc - your models are fully in charge of queries and their lifecycle



Tests
==================

All 3 databases supported have a test suite - to run the tests, create a database called query_test in mysql and psql then run go test at the root of the package. The sqlite tests are disabled by default because enabling them prohibits cross compilation, which is useful if you don't want to install go on your server but just upload a binary compiled locally. 

```bash
go test
```



Versions
==================

- 1.0 - First version with interfaces and chainable finders
- 1.0.1 - Updated to quote table names and fields, for use of reserved words, bug fix for mysql concurrency
- 1.3 - updated API, now shifted instantiation to models instead, to avoid use of reflection
- 1.3.1 - Fixed bugs in Mysql import, updated tests
