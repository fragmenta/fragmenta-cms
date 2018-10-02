# Assets
Assets provides asset compilation, concatenation and fingerprinting. Asset details are stored in a file at secrets/assets.json by default. 

### Usage

Use the assets package to organize assets however you like within your src folder, and output them in compressed form in your public/assets folder when you come to deploy your app. 

```Go 
  // Load asset details from json file on each run
  err := appAssets.Load()
  if err != nil {
    // If no assets loaded, compile for the first time (produces files in public/assets)
    err := appAssets.Compile("src", "public")
    if err != nil {
      server.Fatalf("#error compiling assets %s", err)
    }
  }
```

// Use the asset helpers to generate fingerprinted assets (either one fingerprinted file in production or a list of all files in development) - this is similar to the Rails asset pipeline. 
```Go 
  view.Helpers["style"] = appAssets.StyleLink
  view.Helpers["script"] = appAssets.ScriptLink
```