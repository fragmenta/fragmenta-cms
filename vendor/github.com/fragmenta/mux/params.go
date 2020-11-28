package mux

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Params returns a new set of params parsed from the request.
func Params(r *http.Request) (*RequestParams, error) {
	return ParamsWithMux(mux, r)
}

// ParamsWithMux returns params for a given mux and request
func ParamsWithMux(m *Mux, r *http.Request) (*RequestParams, error) {
	params := &RequestParams{
		Values: make(url.Values, 0),
		Files:  make(map[string][]*multipart.FileHeader, 0),
	}

	// Find the route for request
	route := mux.Match(r)
	if route == nil {
		return nil, errors.New("mux: could not find route for request")
	}

	// Parse the request path params first
	urlParams := route.Parse(r.URL.Path)
	for k, v := range urlParams {
		params.Set(k, []string{v})
	}

	// Add query string params from request
	queryParams := r.URL.Query()
	for k, v := range queryParams {
		params.Add(k, v)
	}

	// If the body is empty, return now without error
	if r.Body == nil {
		return params, nil
	}

	// Parse based on content type
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			return nil, err
		}
		for k, v := range r.Form {
			params.Add(k, v)
		}

	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		err := r.ParseMultipartForm(20 << 20) // 20MB
		if err != nil {
			return nil, err
		}

		// Add the form values
		for k, v := range r.MultipartForm.Value {
			params.Add(k, v)
		}

		// Add the form files
		for k, v := range r.MultipartForm.File {
			params.Files[k] = v
		}
	}

	return params, nil
}

// ParamsJSON returns a new set of params parsed from the request (json included, for testing).
// This is a temporary method for testing json parsing, we should add this capability to Params()
func ParamsJSON(r *http.Request) (*RequestParams, error) {

	params := &RequestParams{
		Values: make(url.Values, 0),
		Files:  make(map[string][]*multipart.FileHeader, 0),
	}

	// Find the route for request
	route := mux.Match(r)
	if route == nil {
		return nil, errors.New("mux: could not find route for request")
	}

	// Parse the request path params first
	urlParams := route.Parse(r.URL.Path)
	for k, v := range urlParams {
		params.Set(k, []string{v})
	}

	// Add query string params from request
	queryParams := r.URL.Query()
	for k, v := range queryParams {
		params.Add(k, v)
	}

	// If the body is empty, return now without error
	if r.Body == nil {
		return params, nil
	}

	// Parse based on content type
	contentType := r.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			return nil, err
		}
		for k, v := range r.Form {
			params.Add(k, v)
		}

	} else if strings.HasPrefix(contentType, "multipart/form-data") {
		err := r.ParseMultipartForm(20 << 20) // 20MB
		if err != nil {
			return nil, err
		}

		// Add the form values
		for k, v := range r.MultipartForm.Value {
			params.Add(k, v)
		}

		// Add the form files
		for k, v := range r.MultipartForm.File {
			params.Files[k] = v
		}
	} else if strings.HasPrefix(contentType, "application/json") {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return params, err
		}

		// If no body provided, return straight away
		if len(body) == 0 {
			return params, err
		}

		rawData := map[string]interface{}{}
		if err := json.Unmarshal(body, &rawData); err != nil {
			return params, err
		}

		// Convert the underlying type of the JSON values
		for k, v := range rawData {
			setParam(params, k, v)
		}
	}

	return params, nil
}

// setParam converts the underlying type of JSON values to strings that we can add
// to the given params. Also handles arrays of strings, floats, ints, and booleans.
func setParam(params *RequestParams, k string, v interface{}) {
	switch v.(type) {
	case int64:
		params.Values.Add(k, fmt.Sprint(v))
	case float64:
		// JSON doesn't handle integers so we get floats
		// interpret as integer if we don't lose information
		f := v.(float64)
		if f == math.Trunc(f) {
			// Add as integer
			params.Values.Add(k, fmt.Sprint(int64(f)))
		} else {
			// Add as float
			params.Values.Add(k, fmt.Sprint(f))
		}
	case string, bool:
		params.Values.Add(k, fmt.Sprint(v))
	case []interface{}:

		// If this is an array, iterate and call setParam recursively to find the underlying
		// type of each element in the array
		for _, i := range v.([]interface{}) {
			setParam(params, k, i)
		}
	default:
		params.Values.Add(k, "")
	}
}

// RequestParams parses all params in a request and stores them in Values
// this includes:
// path params (from route)
// query params (from request)
// body params (from form request bodies)
type RequestParams struct {
	Values url.Values
	Files  map[string][]*multipart.FileHeader
}

// Map returns a flattened map of params with only one entry for each key,
// rather than the array of values Request params allow.
func (p *RequestParams) Map() map[string]string {
	flat := make(map[string]string)

	for k, v := range p.Values {
		flat[k] = v[0]
	}

	return flat
}

// Set sets this key to these values, removing any other entries.
func (p *RequestParams) Set(key string, values []string) {
	p.Values[key] = values
}

// SetString sets this key to this single string value, removing any other entries.
func (p *RequestParams) SetString(key string, v string) {
	p.Set(key, []string{v})
}

// SetInt sets this key to this single string value, removing any other entries.
func (p *RequestParams) SetInt(key string, v int64) {
	p.Set(key, []string{fmt.Sprintf("%d", v)})
}

// Add appends these values to this key, without removing any other entries.
func (p *RequestParams) Add(key string, values []string) {
	p.Values[key] = append(p.Values[key], values...)
}

// Delete all values associated with the key.
func (p *RequestParams) Delete(key string) {
	delete(p.Values, key)
}

// Exists returns true if this key exists in Values
func (p *RequestParams) Exists(key string) bool {
	_, ok := p.Values[key]
	return ok
}

// Get returns the first value for this key or a blank string if no entry.
func (p *RequestParams) Get(key string) string {
	v, ok := p.Values[key]
	if !ok {
		return ""
	}
	return v[0]
}

// GetStrings returns all string values associated with this key
// if there are no values associated an empty array is returned
func (p *RequestParams) GetStrings(key string) []string {
	return p.Values[key]
}

// GetDate returns the first value associated with a given key as a time,
//  using the given time format.
func (p *RequestParams) GetDate(key string, format string) (time.Time, error) {
	v := p.Get(key)
	return time.Parse(format, v)
}

// GetInt returns the first value associated with the given key as an integer.
// If there is no value or a parse error, it returns 0
// If the string contains non-numeric characters, it is truncated from
// the first non-numeric character.
func (p *RequestParams) GetInt(key string) int64 {
	var i int64
	v := p.Get(key)
	// We truncate the string at the first non-numeric character
	v = v[0 : strings.LastIndexAny(v, "0123456789")+1]
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

// GetInts returns all values associated with the key as an array of integers.
func (p *RequestParams) GetInts(key string) []int64 {
	ints := []int64{}

	for _, v := range p.Values[key] {
		vi, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			vi = 0
		}
		ints = append(ints, vi)
	}

	return ints
}

// GetUniqueInts returns all unique non-zero int values
// associated with the given key as an array of integers
func (p *RequestParams) GetUniqueInts(key string) []int64 {
	ints := []int64{}

	for _, v := range p.Values[key] {
		if string(v) == "" {
			continue // ignore blank ints
		}
		vi, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			vi = 0
		}

		// Do not insert 0, or duplicate entries
		if vi > 0 && !contains(ints, vi) {
			ints = append(ints, vi)
		}
	}

	return ints
}

// GetIntsString returns all values associated with the key
// as a comma separated string.
func (p *RequestParams) GetIntsString(key string) string {
	ints := ""

	for _, v := range p.Values[key] {
		if "" == string(v) {
			continue // ignore blank ints
		}

		if len(ints) > 0 {
			ints += "," + string(v)
		} else {
			ints += string(v)
		}

	}

	return ints
}

// GetFloat returns the first value associated with the key as an integer.
// If there is no value or a parse error, it returns 0.0
func (p *RequestParams) GetFloat(key string) float64 {
	var value float64
	v := p.Get(key)
	// Remove percent signs from float values
	v = strings.Replace(v, "%", "", -1)
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0.0
	}
	return value
}

// GetFloats returns all values associated with the key as an array of floats.
func (p *RequestParams) GetFloats(key string) []float64 {
	var values []float64
	for _, v := range p.Values[key] {
		// Remove percent signs from float values
		v = strings.Replace(v, "%", "", -1)
		value, err := strconv.ParseFloat(v, 64)
		if err != nil {
			value = 0.0
		}
		values = append(values, value)
	}
	return values
}

// contains returns true if this array of ints contains the given int
func contains(list []int64, item int64) bool {
	for _, b := range list {
		if b == item {
			return true
		}
	}
	return false
}
