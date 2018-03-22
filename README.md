# MHub API Utilities

## Quick Start
Install MHub API Utilities by running the following in your terminal:
```
go get -u github.com/TheRedBricks/mhub-api-utilities
```

## Logger Middleware
Logger middleware to extend HTTP servers. Logger will print method, URL, HTTP Status and time taken in your terminal.

### With Goji
```go
package main

import (
	"io"
	"net/http"

	"github.com/TheRedBricks/mhub-api-utilities/logger"

	"goji.io"
	"goji.io/pat"
)

func main() {
	mux := goji.NewMux()

	mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "MHub API")
	})

	mux.Use(logger.Middleware)

	http.ListenAndServe(":8000", mux)
}
```
