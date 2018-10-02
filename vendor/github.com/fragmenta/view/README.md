# view
Package view provides template registration, rendering, and helper functions for golang views

### Usage 

Load templates on app startup:

```Go 
	err := view.LoadTemplates()
	if err != nil {
		server.Fatalf("Error reading templates %s", err)
	}
```

Render a template 

```Go 
    // Set up the view
	view := view.New(context)
    // Add a key to the view
    view.AddKey("page", page)
    // Optionally set template, layout or other attributes
    view.Template("src/pages/views/home.html.got")
    // Render the view
    return view.Render()
```


Public subpackages:

* helpers - utilities for handling files
