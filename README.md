# MHub API Utilities

## Quick Start
Install MHub API Utilities by running the following in your terminal:
```
go get -u github.com/TheRedBricks/mhub-api-utilities
```

## Logger Middleware
Logger middleware to extend HTTP servers. Logger will print method, URL, HTTP Status and time taken in your terminal.

![Terminal Screenshot](https://user-images.githubusercontent.com/1572333/37758694-0e1f0df8-2dec-11e8-920e-e30dcb0160f2.png "Terminal Screenshot")

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
