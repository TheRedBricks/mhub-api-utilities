package trackrequest

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// onRequest receives a RequestLog
type onRequest func(*RequestLog)

// onError receives error
type onError func(error)

// identify request log with user
type identify func(RequestLog) string

// Manager structure
type Manager struct {
	OnRequest         onRequest
	OnRequestComplete onRequest
	OnError           onError

	Identify identify
}

// RequestLog structure
type RequestLog struct {
	IdentityID string
	Method     string
	URL        string
	IP         string
	Headers    map[string][]string
	Cookies    map[string][]string
	Body       []byte
	CreatedAt  *time.Time
	TimeTaken  time.Duration
}

// NewManager sets up new Manager
func NewManager(manager *Manager) *Manager {
	return manager
}

// Middleware to track requests
func (manager *Manager) Middleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log := &RequestLog{
			Method:    r.Method,
			URL:       r.URL.String(),
			IP:        r.Header.Get("X-Forwarded-For"),
			Headers:   make(map[string][]string),
			Cookies:   make(map[string][]string),
			CreatedAt: &start,
		}

		// store headers
		for name, value := range r.Header {
			log.Headers[name] = value
		}

		// store cookies
		for _, value := range r.Cookies() {
			c := strings.Split(value.String(), "=")
			if len(c) >= 2 {
				log.Cookies[c[0]] = append(log.Cookies[c[0]], strings.Join(c[1:], "="))
			}
		}

		// store body
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			if manager.OnError != nil {
				manager.OnError(err)
			}
		} else {
			log.Body = bodyBytes

			// reattach body
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// trigger OnRequest when present
		if manager.OnRequest != nil {
			manager.OnRequest(log)
		}

		h.ServeHTTP(w, r)

		// identify request to the user
		if manager.Identify != nil {
			identityID := manager.Identify(*log)
			if identityID != "" {
				log.IdentityID = identityID
			}
		}

		// trigger completed call
		if manager.OnRequestComplete != nil {
			// keep track of time taken
			log.TimeTaken = time.Since(start)

			manager.OnRequestComplete(log)
		}
	}
	return http.HandlerFunc(fn)
}
