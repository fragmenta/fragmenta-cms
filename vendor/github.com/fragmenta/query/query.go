// Package query lets you build and execute SQL chainable queries against a database of your choice, and defer execution of SQL until you wish to extract a count or array of models.

// NB in order to allow cross-compilation, we exlude sqlite drivers by default
// uncomment them to allow use of sqlite

package query

import (
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// FIXME - this package global should in theory be protected by a mutex, even if it is only for debugging

// Debug sets whether we output debug statements for SQL
var Debug bool

func init() {
	Debug = false // default to false
}

// Result holds the results of a query as map[string]interface{}
type Result map[string]interface{}

// Func is a function which applies effects to queries
type Func func(q *Query) *Query

// Query provides all the chainable relational query builder methods
type Query struct {

	// Database - database name and primary key, set with New()
	tablename  string
	primarykey string

	// SQL - Private fields used to store sql before building sql query
	sql    string
	sel    string
	join   string
	where  string
	group  string
	having string
	order  string
	offset string
	limit  string

	// Extra args to be substituted in the *where* clause
	args []interface{}
}

// New builds a new Query, given the table and primary key
func New(t string, pk string) *Query {

	// If we have no db, return nil
	if database == nil {
		return nil
	}

	q := &Query{
		tablename:  t,
		primarykey: pk,
	}

	return q
}

// Exec the given sql and args against the database directly
// Returning sql.Result (NB not rows)
func Exec(sql string, args ...interface{}) (sql.Result, error) {
	results, err := database.Exec(sql, args...)
	return results, err
}

// Rows executes the given sql and args against the database directly
// Returning sql.Rows
func Rows(sql string, args ...interface{}) (*sql.Rows, error) {
	results, err := database.Query(sql, args...)
	return results, err
}

// Copy returns a new copy of this query which can be mutated without affecting the original
func (q *Query) Copy() *Query {
	return &Query{
		tablename:  q.tablename,
		primarykey: q.primarykey,
		sql:        q.sql,
		sel:        q.sel,
		join:       q.join,
		where:      q.where,
		group:      q.group,
		having:     q.having,
		order:      q.order,
		offset:     q.offset,
		limit:      q.limit,
		args:       q.args,
	}
}

// TODO: These should instead be something like query.New("table_name").Join(a,b).Insert() and just have one multiple function?

// InsertJoin inserts a join clause on the query
func (q *Query) InsertJoin(a int64, b int64) error {
	return q.InsertJoins([]int64{a}, []int64{b})
}

// InsertJoins using an array of ids (more general version of above)
// This inserts joins for every possible relation between the ids
func (q *Query) InsertJoins(a []int64, b []int64) error {

	// Make sure we have some data
	if len(a) == 0 || len(b) == 0 {
		return fmt.Errorf("Null data for joins insert %s", q.table())
	}

	// Check for null entries in start of data - this is not a good idea.
	//	if a[0] == 0 || b[0]  == 0 {
	//		return fmt.Errorf("Zero data for joins insert %s", q.table())
	//	}

	values := ""
	for _, av := range a {
		for _, bv := range b {
			// NB no zero values allowed, we simply ignore zero values
			if av != 0 && bv != 0 {
				values += fmt.Sprintf("(%d,%d),", av, bv)
			}

		}
	}

	values = strings.TrimRight(values, ",")

	sql := fmt.Sprintf("INSERT into %s VALUES %s;", q.table(), values)

	if Debug {
		fmt.Printf("JOINS SQL:%s\n", sql)
	}

	_, err := database.Exec(sql)
	return err
}

// UpdateJoins updates the given joins, using the given id to clear joins first
func (q *Query) UpdateJoins(id int64, a []int64, b []int64) error {

	if Debug {
		fmt.Printf("SetJoins %s %s=%d: %v %v \n", q.table(), q.pk(), id, a, b)
	}

	// First delete any existing joins
	err := q.Where(fmt.Sprintf("%s=?", q.pk()), id).Delete()
	if err != nil {
		return err
	}

	// Now join all a's with all b's by generating joins for each possible combination

	// Make sure we have data in both cases, otherwise do not attempt insert any joins
	if len(a) > 0 && len(b) > 0 {
		// Now insert all new ids - NB the order of arguments here MUST match the order in the table
		err = q.InsertJoins(a, b)
		if err != nil {
			return err
		}
	}

	return nil
}

// Insert inserts a record in the database
func (q *Query) Insert(params map[string]string) (int64, error) {

	// Insert and retrieve ID in one step from db
	sql := q.insertSQL(params)

	if Debug {
		fmt.Printf("INSERT SQL:%s %v\n", sql, valuesFromParams(params))
	}

	id, err := database.Insert(sql, valuesFromParams(params)...)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// insertSQL sets the insert sql for update statements, turn params into sql i.e. "col"=?
// NB we always use parameterized queries, never string values.
func (q *Query) insertSQL(params map[string]string) string {
	var cols, vals []string

	for i, k := range sortedParamKeys(params) {
		cols = append(cols, database.QuoteField(k))
		vals = append(vals, database.Placeholder(i+1))
	}
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s) %s;", q.table(), strings.Join(cols, ","), strings.Join(vals, ","), database.InsertSQL(q.pk()))

	return query
}

// Update one model specified in this query - the column names MUST be verified in the model
func (q *Query) Update(params map[string]string) error {
	// We should check the query has a where limitation to avoid updating all?
	// pq unfortunately does not accept limit(1) here
	return q.UpdateAll(params)
}

// Delete one model specified in this relation
func (q *Query) Delete() error {
	// We should check the query has a where limitation?
	return q.DeleteAll()
}

// UpdateAll updates all models specified in this relation
func (q *Query) UpdateAll(params map[string]string) error {
	// Create sql for update from ALL params
	q.Select(fmt.Sprintf("UPDATE %s SET %s", q.table(), querySQL(params)))

	// Execute, after PREpending params to args
	// in an update statement, the where comes at the end
	q.args = append(valuesFromParams(params), q.args...)

	if Debug {
		fmt.Printf("UPDATE SQL:%s\n%v\n", q.QueryString(), valuesFromParams(params))
	}

	_, err := q.Result()

	return err
}

// DeleteAll delets *all* models specified in this relation
func (q *Query) DeleteAll() error {

	q.Select(fmt.Sprintf("DELETE FROM %s", q.table()))

	if Debug {
		fmt.Printf("DELETE SQL:%s <= %v\n", q.QueryString(), q.args)
	}

	// Execute
	_, err := q.Result()

	return err
}

// Count fetches a count of model objects (executes SQL).
func (q *Query) Count() (int64, error) {

	// In order to get consistent results, we use the same query builder
	// but reset select to simple count select

	// Store the previous select and set
	s := q.sel
	countSelect := fmt.Sprintf("SELECT COUNT(%s) FROM %s", q.pk(), q.table())
	q.Select(countSelect)

	// Store the previous order (minus order by) and set to empty
	// Order must be blank on count because of limited select
	o := strings.Replace(q.order, "ORDER BY ", "", 1)
	q.order = ""

	// Fetch count from db for our sql with count select and no order set
	var count int64
	rows, err := q.Rows()
	if err != nil {
		return 0, fmt.Errorf("Error querying database for count: %s\nQuery:%s", err, q.QueryString())
	}

	// We expect just one row, with one column (count)
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}

	// Reset select after getting count query
	q.Select(s)
	q.Order(o)
	q.reset()

	return count, err
}

// Result executes the query against the database, returning sql.Result, and error (no rows)
// (Executes SQL)
func (q *Query) Result() (sql.Result, error) {
	results, err := database.Exec(q.QueryString(), q.args...)
	return results, err
}

// Rows executes the query against the database, and return the sql rows result for this query
// (Executes SQL)
func (q *Query) Rows() (*sql.Rows, error) {
	results, err := database.Query(q.QueryString(), q.args...)
	return results, err
}

// FirstResult executes the SQL and returrns the first result
func (q *Query) FirstResult() (Result, error) {

	// Set a limit on the query
	q.Limit(1)

	// Fetch all results (1)
	results, err := q.Results()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("No results found for Query:%s", q.QueryString())
	}

	// Return the first result
	return results[0], nil
}

// ResultInt64 returns the first result from a query stored in the column named col as an int64.
func (q *Query) ResultInt64(c string) (int64, error) {
	result, err := q.FirstResult()
	if err != nil || result[c] == nil {
		return 0, err
	}
	var i int64
	switch result[c].(type) {
	case int64:
		i = result[c].(int64)
	case int:
		i = int64(result[c].(int))
	case float64:
		i = int64(result[c].(float64))
	case string:
		f, err := strconv.ParseFloat(result[c].(string), 64)
		if err != nil {
			return i, err
		}
		i = int64(f)
	}

	return i, nil
}

// ResultFloat64 returns the first result from a query stored in the column named col as a float64.
func (q *Query) ResultFloat64(c string) (float64, error) {
	result, err := q.FirstResult()
	if err != nil || result[c] == nil {
		return 0, err
	}
	var f float64
	switch result[c].(type) {
	case float64:
		f = result[c].(float64)
	case int:
		f = float64(result[c].(int))
	case int64:
		f = float64(result[c].(int))
	case string:
		f, err = strconv.ParseFloat(result[c].(string), 64)
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

// Results returns an array of results
func (q *Query) Results() ([]Result, error) {

	// Make an empty result set map
	var results []Result

	// Fetch rows from db for our sql
	rows, err := q.Rows()

	if err != nil {
		return results, fmt.Errorf("Error querying database for rows: %s\nQUERY:%s", err, q)
	}

	// Close rows before returning
	defer rows.Close()

	// Fetch the columns from the database
	cols, err := rows.Columns()
	if err != nil {
		return results, fmt.Errorf("Error fetching columns: %s\nQUERY:%s\nCOLS:%s", err, q, cols)
	}

	// For each row, construct an entry in results with a map of column string keys to values
	for rows.Next() {
		result, err := scanRow(cols, rows)
		if err != nil {
			return results, fmt.Errorf("Error fetching row: %s\nQUERY:%s\nCOLS:%s", err, q, cols)
		}
		results = append(results, result)
	}

	return results, nil
}

// ResultIDs returns an array of ids as the result of a query
// FIXME - this should really use the query primary key, not "id" hardcoded
func (q *Query) ResultIDs() []int64 {
	var ids []int64
	if Debug {
		fmt.Printf("#info ResultIDs:%s\n", q.DebugString())
	}
	results, err := q.Results()
	if err != nil {
		return ids
	}

	for _, r := range results {
		if r["id"] != nil {
			ids = append(ids, r["id"].(int64))
		}
	}

	return ids
}

// ResultIDSets returns a map from a values to arrays of b values, the order of a,b is respected not the table key order
func (q *Query) ResultIDSets(a, b string) map[int64][]int64 {
	idSets := make(map[int64][]int64, 0)

	results, err := q.Results()
	if err != nil {
		return idSets
	}

	for _, r := range results {
		if r[a] != nil && r[b] != nil {
			av := r[a].(int64)
			bv := r[b].(int64)
			idSets[av] = append(idSets[av], bv)
		}
	}
	if Debug {
		fmt.Printf("#info ResultIDSets:%s\n", q.DebugString())
	}
	return idSets
}

// QueryString builds a query string to use for results
func (q *Query) QueryString() string {

	if q.sql == "" {

		// if we have arguments override the selector
		if q.sel == "" {
			// Note q.table() etc perform quoting on field names
			q.sel = fmt.Sprintf("SELECT %s.* FROM %s", q.table(), q.table())
		}

		q.sql = fmt.Sprintf("%s %s %s %s %s %s %s %s", q.sel, q.join, q.where, q.group, q.having, q.order, q.offset, q.limit)
		q.sql = strings.TrimRight(q.sql, " ")
		q.sql = strings.Replace(q.sql, "  ", " ", -1)
		q.sql = strings.Replace(q.sql, "   ", " ", -1)

		// Replace ? with whatever placeholder db prefers
		q.replaceArgPlaceholders()

		q.sql = q.sql + ";"
	}

	return q.sql
}

// CHAINABLE FINDERS

// Apply the Func to this query, and return the modified Query
// This allows chainable finders from other packages
// e.g. q.Apply(status.Published) where status.Published is a Func
func (q *Query) Apply(f Func) *Query {
	return f(q)
}

// Conditions applies a series of query funcs to a query
func (q *Query) Conditions(funcs ...Func) *Query {
	for _, f := range funcs {
		q = f(q)
	}
	return q
}

// SQL defines sql manually and overrides all other setters
// Completely replaces all stored sql
func (q *Query) SQL(sql string) *Query {
	q.sql = sql
	q.reset()
	return q
}

// Limit sets the sql LIMIT with an int
func (q *Query) Limit(limit int) *Query {
	q.limit = fmt.Sprintf("LIMIT %d", limit)
	q.reset()
	return q
}

// Offset sets the sql OFFSET with an int
func (q *Query) Offset(offset int) *Query {
	q.offset = fmt.Sprintf("OFFSET %d", offset)
	q.reset()
	return q
}

// Where defines a WHERE clause on SQL - Additional calls add WHERE () AND () clauses
func (q *Query) Where(sql string, args ...interface{}) *Query {

	if len(q.where) > 0 {
		q.where = fmt.Sprintf("%s AND (%s)", q.where, sql)
	} else {
		q.where = fmt.Sprintf("WHERE (%s)", sql)
	}

	// NB this assumes that args are only supplied for where clauses
	// this may be an incorrect assumption!
	if args != nil {
		if q.args == nil {
			q.args = args
		} else {
			q.args = append(q.args, args...)
		}
	}

	q.reset()
	return q
}

// OrWhere defines a where clause on SQL - Additional calls add WHERE () OR () clauses
func (q *Query) OrWhere(sql string, args ...interface{}) *Query {

	if len(q.where) > 0 {
		q.where = fmt.Sprintf("%s OR (%s)", q.where, sql)
	} else {
		q.where = fmt.Sprintf("WHERE (%s)", sql)
	}

	if args != nil {
		if q.args == nil {
			q.args = args
		} else {
			q.args = append(q.args, args...)
		}
	}

	q.reset()
	return q
}

// WhereIn adds a Where clause which selects records IN() the given array
// If IDs is an empty array, the query limit is set to 0
func (q *Query) WhereIn(col string, IDs []int64) *Query {
	// Return no results, so that when chaining callers
	// don't have to check for empty arrays
	if len(IDs) == 0 {
		q.Limit(0)
		q.reset()
		return q
	}

	in := ""
	for _, ID := range IDs {
		in = fmt.Sprintf("%s%d,", in, ID)
	}
	in = strings.TrimRight(in, ",")
	sql := fmt.Sprintf("%s IN (%s)", col, in)

	if len(q.where) > 0 {
		q.where = fmt.Sprintf("%s AND (%s)", q.where, sql)
	} else {
		q.where = fmt.Sprintf("WHERE (%s)", sql)
	}

	q.reset()
	return q
}

// Define a join clause on SQL - we create an inner join like this:
// INNER JOIN extras_seasons ON extras.id = extra_id
// q.Select("SELECT units.* FROM units INNER JOIN sites ON units.site_id = sites.id")

// rails join example
// INNER JOIN "posts_tags" ON "posts_tags"."tag_id" = "tags"."id" WHERE "posts_tags"."post_id" = 111

// Join adds an inner join to the query
func (q *Query) Join(otherModel string) *Query {
	modelTable := q.tablename

	tables := []string{
		modelTable,
		ToPlural(otherModel),
	}
	sort.Strings(tables)
	joinTable := fmt.Sprintf("%s_%s", tables[0], tables[1])

	sql := fmt.Sprintf("INNER JOIN %s ON %s.id = %s.%s_id", database.QuoteField(joinTable), database.QuoteField(modelTable), database.QuoteField(joinTable), ToSingular(modelTable))

	if len(q.join) > 0 {
		q.join = fmt.Sprintf("%s %s", q.join, sql)
	} else {
		q.join = fmt.Sprintf("%s", sql)
	}

	q.reset()
	return q
}

// Order defines ORDER BY sql
func (q *Query) Order(sql string) *Query {
	if sql == "" {
		q.order = ""
	} else {
		q.order = fmt.Sprintf("ORDER BY %s", sql)
	}
	q.reset()

	return q
}

// Group defines GROUP BY sql
func (q *Query) Group(sql string) *Query {
	if sql == "" {
		q.group = ""
	} else {
		q.group = fmt.Sprintf("GROUP BY %s", sql)
	}
	q.reset()
	return q
}

// Having defines HAVING sql
func (q *Query) Having(sql string) *Query {
	if sql == "" {
		q.having = ""
	} else {
		q.having = fmt.Sprintf("HAVING %s", sql)
	}
	q.reset()
	return q
}

// Select defines SELECT  sql
func (q *Query) Select(sql string) *Query {
	q.sel = sql
	q.reset()
	return q
}

// DebugString returns a query representation string useful for debugging
func (q *Query) DebugString() string {
	return fmt.Sprintf("--\nQuery-SQL:%s\nARGS:%s\n--", q.QueryString(), q.argString())
}

// Clear sql/query caches
func (q *Query) reset() {
	// Perhaps later clear cached compiled representation of query too

	// clear stored sql
	q.sql = ""
}

// Return an arg string (for debugging)
func (q *Query) argString() string {
	output := "-"

	for _, a := range q.args {
		output = output + fmt.Sprintf("'%s',", q.argToString(a))
	}
	output = strings.TrimRight(output, ",")
	output = output + ""

	return output
}

// Convert arguments to string - used only for debug argument strings
// Not to be exported or used to try to escape strings...
func (q *Query) argToString(arg interface{}) string {
	switch arg.(type) {
	case string:
		return arg.(string)
	case []byte:
		return string(arg.([]byte))
	case int, int8, int16, int32, uint, uint8, uint16, uint32:
		return fmt.Sprintf("%d", arg)
	case int64, uint64:
		return fmt.Sprintf("%d", arg)
	case float32, float64:
		return fmt.Sprintf("%f", arg)
	case bool:
		return fmt.Sprintf("%d", arg)
	default:
		return fmt.Sprintf("%v", arg)

	}

}

// Ask model for primary key name to use
func (q *Query) pk() string {
	return database.QuoteField(q.primarykey)
}

// Ask model for table name to use
func (q *Query) table() string {
	return database.QuoteField(q.tablename)
}

// Replace ? with whatever database prefers (psql uses numbered args)
func (q *Query) replaceArgPlaceholders() {
	// Match ? and replace with argument placeholder from database
	for i := range q.args {
		q.sql = strings.Replace(q.sql, "?", database.Placeholder(i+1), 1)
	}
}

// Sorts the param names given - map iteration order is explicitly random in Go
// but we need params in a defined order to avoid unexpected results.
func sortedParamKeys(params map[string]string) []string {
	sortedKeys := make([]string, len(params))
	i := 0
	for k := range params {
		sortedKeys[i] = k
		i++
	}
	sort.Strings(sortedKeys)

	return sortedKeys
}

// Generate a set of values for the params in order
func valuesFromParams(params map[string]string) []interface{} {

	// NB DO NOT DEPEND ON PARAMS ORDER - see note on SortedParamKeys
	var values []interface{}
	for _, key := range sortedParamKeys(params) {
		values = append(values, params[key])
	}
	return values
}

// Used for update statements, turn params into sql i.e. "col"=?
func querySQL(params map[string]string) string {
	var output []string
	for _, k := range sortedParamKeys(params) {
		output = append(output, fmt.Sprintf("%s=?", database.QuoteField(k)))
	}
	return strings.Join(output, ",")
}

func scanRow(cols []string, rows *sql.Rows) (Result, error) {

	// We return a map[string]interface{} for each row scanned
	result := Result{}

	values := make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		var col interface{}
		values[i] = &col
	}

	// Scan results into these interfaces
	err := rows.Scan(values...)
	if err != nil {
		return nil, fmt.Errorf("Error scanning row: %s", err)
	}

	// Make a string => interface map and hand off to caller
	// We fix up a few types which the pq driver returns as less handy equivalents
	// We enforce usage of int64 at all times as all our records use int64
	for i := 0; i < len(cols); i++ {
		v := *values[i].(*interface{})
		if values[i] != nil {
			switch v.(type) {
			default:
				result[cols[i]] = v
			case bool:
				result[cols[i]] = v.(bool)
			case int:
				result[cols[i]] = int64(v.(int))
			case []byte: // text cols are given as bytes
				result[cols[i]] = string(v.([]byte))
			case int64:
				result[cols[i]] = v.(int64)
			}
		}

	}

	return result, nil
}
