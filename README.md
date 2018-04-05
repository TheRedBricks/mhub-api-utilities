# MHub API Utilities

## Quick Start
Install MHub API Utilities by running the following in your terminal:
```
go get -u github.com/TheRedBricks/mhub-api-utilities
```

## Logger Middleware
Middleware to extend HTTP servers. Logger will print method, URL, HTTP Status and time taken in your terminal.

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
	mux.Use(logger.Middleware)

	mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "MHub API")
	})

	http.ListenAndServe(":8000", mux)
}
```

## Track Request Middleware
Middleware to extend HTTP servers. Use tracker to keep track of requests coming through the server. Method, URL, IP, Headers, Body & Cookies are returned on request.

![Terminal Screenshot](https://user-images.githubusercontent.com/1572333/38349283-e579f4ec-38d8-11e8-9be0-f317d1cf3f8a.png "Terminal Screenshot")

### With Goji
```go
package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/TheRedBricks/mhub-api-utilities/trackrequest"

	"goji.io"
	"goji.io/pat"
)

func main() {
	mux := goji.NewMux()

	tr := trackrequest.NewManager(&trackrequest.Manager{})
	tr.OnRequest = func(log *trackrequest.RequestLog) {
		// save log to DB or display into terminal
		// onRequest is triggered as soon as requests comes in
		fmt.Printf("-- onRequest --------- %+v\n", log)
	}
	tr.OnRequestComplete = func(log *trackrequest.RequestLog) {
		// save log to DB or display into terminal
		// onRequestComplete is triggered at the end
		// onRequestComplete includes TimeTaken and IdentityID
		fmt.Printf("-- onRequestComplete - %+v\n", log)
	}
	tr.OnError = func(err error) {
		// handle err optional
		fmt.Println(err)
	}
	mux.Use(tr.Middleware)

	// other middleware...
	mux.Use(func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// attach user to log if logged in
			tr.Identify = func(log trackrequest.RequestLog) string {
				return "jo@nne.my"
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})

	mux.HandleFunc(pat.Get("/"), func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "MHub API")
	})

	http.ListenAndServe(":8000", mux)
}
```
