# Fragmenta Multiplexer (mux)

Fragmenta mux is a replacement the standard http.ServeMux which offers a few additional features and improved efficiency. Features are very similar to gorilla/mux but with a few additions, and it is compatible with the standard http.Handler interface or handlers returning error.  

It offers the following features:

* Named paramaters including regexp matches for params (e.g. {id:\d+} to match id only to one or more numerals)
* Delayed param parsing (url,query,form) with utility functions for extracting Int, Bool, Float params. 
* Routes are evaluated strictly in order - add important routes first and catch-alls at the end 
* Zero allocations when matching means low-memory use and responses as fast as httprouter for static routes
* A cache in front of route matching speeds up responses (under 100ns/op in a simple static case)
* Low memory usage (even with cache) 
* Accepts either the standard http.Handler interface or mux.Handler (same but with error return)
* Add middleware http.HandlerFunc for chaining standard Go middleware for auth, logging etc.

It does not offer:

* Nested routes or groups 


## Install 

Perform the usual incantation: 

```sh
  go get -u github.com/fragmenta/mux
```

## Usage 

Usage is as you'd expect if you've used the stdlib mux or gorilla mux. You can use the mux.Add/Get/Post to add handlers which return an error, or mux.AddHandler to add a stdlib http.HandlerFunc.

```go

func main() {
  m := mux.New()
  m.Get(`/`,homeHandler)
  m.Get(`/users`,users.HandleIndex)
  m.Post(`/users`,users.HandleCreate)
  m.Post(`/users/{id:\d+}/update`,users.HandleUpdate)
  http.Handle("/", r)
}


```

## Errors

Because of the handler signature returning errors, you can set an ErrorHandler which is called if an error occurs inside one of your handlers, and a FileHandler which is called for serving files if no route is found. This makes handling errors more elegant, instead of this:


```go

if err != nil {
  log.Printf("error occured:%s",err)
  // .. do something to handle and display to user
  return 
}

```

you can do this in your handlers: 

```go

if err != nil {
  return err
}

```

and display errors in a consistent way using your ErrorHandler function (you can also return a custom error type from handlers as fragmenta does to send more information than just error).


## Params

Parsing of params is delayed until you require them in your handler - no parsing is done until that point. When you do require them, just parse params as follows, and a full params object will be available with a map of all params from urls, and form bodies. Multipart file forms are parsed automatically and the files made available for use. 

```go

// Parse  params (any url, query and form params)
params,err := mux.Params(request)
if err != nil {
  return err
}

params.Values["key"][4]
params.Get("my_query_key")
params.GetInt("user_id")
params.GetFloat("float")
params.GetBool("bool")
params.GetDate("published_at","2017-01-02")

for _,fh := range params.Files {
  
}

```

## Benchmarks 

Speed isn't everything (see the list of features above), but it is important the router doesn't slow down request times, particularly if you have a lot of urls to match. For benchmarks against a few popular routers, see https://github.com/kennygrant/routebench

Performance is adequate:

```

BenchmarkStatic/stdlib_mux-4         	    1000	   1946545 ns/op	   20619 B/op	     537 allocs/op
BenchmarkStatic/gorilla_mux-4        	    1000	   1846382 ns/op	  115648 B/op	    1578 allocs/op
BenchmarkStatic/fragmenta_mux-4      	  100000	     13969 ns/op	       0 B/op	       0 allocs/op
BenchmarkStatic/httprouter_mux-4     	  100000	     16240 ns/op	       0 B/op	       0 allocs/op

BenchmarkGithubFuzz/stdlib_mux-4               	     300	   4592686 ns/op	   35767 B/op	     902 allocs/op
BenchmarkGithubFuzz/gorilla_mux-4              	     100	  12931693 ns/op	  246784 B/op	    2590 allocs/op
BenchmarkGithubFuzz/fragmenta_mux-4            	    5000	    324911 ns/op	    7617 B/op	     136 allocs/op
BenchmarkGithubFuzz/httprouter_mux-4           	   10000	    101702 ns/op	   23791 B/op	     296 allocs/op


```


