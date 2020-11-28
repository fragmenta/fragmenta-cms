# server
A wrapper for the net/http server offering a few other features:

* Config loading from a json config file (for use in setup/handlers)
* Levelled, Structured logging using a flexible set of loggers - easy to add your own custom loggers 
* Optional logging middleware for requests
* Scheduling of tasks at specific times of day and intervals
* Uses the new request context in the Go stdlib
* Add tracing to requests so that you can follow a request id through handlers
* Requires Go 1.8

Example usage:

```go

  // Redirect all :80 traffic to our canonical url on :443
	server.StartRedirectAll(80, config.Get("root_url"))

	// If in production, serve over tls with autocerts from let's encrypt
	err = server.StartTLSAutocert(config.Get("autocert_email"), config.Get("autocert_domains"))
	if err != nil {
			server.Fatalf("Error starting server %s", err)
	}

```

## Config 

The config package offers access to json config files containing dev/production/test configs. 

Example usage:

```go

  // Load config
  config.Load(path)

  // Get a key from our current config (dev/prod/test)
  config.Get("mykey")

```

## Logging

The logging package offers structured, levelled logging which can be configured to send to a file, stdout, and/or other services like an influxdb server with additional plugin loggers. You can add as many loggers which log events as you want, and because logging is structured, each logger can decide which information to act on. Example log output to sdtout is below (real colouring is nicer):

Example usage:

```go
  
  // Set up a stderr logger with time prefix
	logger, err := log.NewStdErr(log.PrefixDateTime)
	if err != nil {
		return err
	}
  
  // Add to the list of loggers receiving events
	log.Add(logger)

```

```bash
2017-01-16:00:37:05 Starting server port:3000 #info 
2017-01-16:00:37:05 Finished loading assets in 109.483µs #info 
2017-01-16:00:37:05 Finished loading templates in 3.184977ms #info 
2017-01-16:00:37:05 Finished opening database in 6.387409ms db:mydb user:myuser #info 
2017-01-16:00:37:05 Finished loading server in 9.99619ms #info 
2017-01-16:00:37:06 <- Request ip:[::1]:64913 len:0 method:GET trace:07466847-28899DB4 url:/ #info 
2017-01-16:00:37:06  in handler using request context trace:07466847-28899DB4 #info 
2017-01-16:00:37:06 -> Response in 3.005292ms trace:07466847-28899DB4 url:/ #info 
2017-01-16:00:37:07 <- Request ip:[::1]:64913 len:0 method:GET trace:A0E55A1B-012DA648 url:/ #info 
2017-01-16:00:37:07  in handler using request context trace:A0E55A1B-012DA648 #info 
2017-01-16:00:37:07 -> Response in 3.32221ms trace:A0E55A1B-012DA648 url:/ #info 
```

## Scheduling

A scheduling facility so that you can schedule actions (like sending a tweet) on app startup. 

```go 
  
   schedule.At(func(){}, context, time, repeatDuration)

```
