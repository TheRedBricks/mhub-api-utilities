package trackrequest_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/TheRedBricks/mhub-api-utilities/trackrequest"
)

func testHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ok" {
			io.WriteString(w, "OK")
			return
		}
		if r.URL.Path == "/user" {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			w.Write(bodyBytes)
			return
		}
		if r.URL.Path == "/not_found" {
			http.Error(w, "Page Not Found", http.StatusNotFound)
			return
		}
		if r.URL.Path == "/error" {
			http.Error(w, "An error occurred", http.StatusInternalServerError)
			return
		}
		panic("URL Not handled")
	}
	return http.HandlerFunc(fn)
}

const ip = "192.168.0.1"

func TestTrackRequest(t *testing.T) {
	type expects struct {
		recordLog      *trackrequest.RequestLog
		responseStatus int
		responseBody   string
	}
	tests := map[string]struct {
		description string
		method      string
		url         string
		body        string
		expects     *expects
	}{
		"GET /ok": {
			method: "GET",
			url:    "/ok",
			body:   "",
			expects: &expects{
				recordLog: &trackrequest.RequestLog{
					IdentityID: "jo@nne.my",
					Method:     "GET",
					URL:        "/ok",
					IP:         ip,
					Body:       []byte(""),
					Headers: map[string][]string{
						"User-Agent":      []string{"test-client"},
						"X-Forwarded-For": []string{ip},
					},
					Cookies: map[string][]string{
						"cookie_name":         []string{"cookie_value_1", "cookie_value_2"},
						"another_cookie_name": []string{"cookie_value_=strange"},
					},
				},
				responseStatus: http.StatusOK,
				responseBody:   "OK",
			},
		},
		"POST /user": {
			method: "POST",
			url:    "/user",
			body:   "{\"name\":\"Joanne\",\"email\":\"jo@nne.my\"}",
			expects: &expects{
				recordLog: &trackrequest.RequestLog{
					IdentityID: "jo@nne.my",
					Method:     "POST",
					URL:        "/user",
					IP:         ip,
					Body:       []byte("{\"name\":\"Joanne\",\"email\":\"jo@nne.my\"}"),
					Headers: map[string][]string{
						"User-Agent":      []string{"test-client"},
						"X-Forwarded-For": []string{ip},
					},
					Cookies: map[string][]string{
						"cookie_name":         []string{"cookie_value_1", "cookie_value_2"},
						"another_cookie_name": []string{"cookie_value_=strange"},
					},
				},
				responseStatus: http.StatusOK,
				responseBody:   "{\"name\":\"Joanne\",\"email\":\"jo@nne.my\"}",
			},
		},
		"PUT /user": {
			method: "PUT",
			url:    "/user",
			body:   "{\"name\":\"Joanne\",\"email\":\"jo@nne.my\"}",
			expects: &expects{
				recordLog: &trackrequest.RequestLog{
					IdentityID: "jo@nne.my",
					Method:     "PUT",
					URL:        "/user",
					IP:         ip,
					Body:       []byte("{\"name\":\"Joanne\",\"email\":\"jo@nne.my\"}"),
					Headers: map[string][]string{
						"User-Agent":      []string{"test-client"},
						"X-Forwarded-For": []string{ip},
					},
					Cookies: map[string][]string{
						"cookie_name":         []string{"cookie_value_1", "cookie_value_2"},
						"another_cookie_name": []string{"cookie_value_=strange"},
					},
				},
				responseStatus: http.StatusOK,
				responseBody:   "{\"name\":\"Joanne\",\"email\":\"jo@nne.my\"}",
			},
		},
		"GET /error": {
			method: "GET",
			url:    "/error",
			body:   "",
			expects: &expects{
				recordLog: &trackrequest.RequestLog{
					IdentityID: "jo@nne.my",
					Method:     "GET",
					URL:        "/error",
					IP:         ip,
					Body:       []byte(""),
					Headers: map[string][]string{
						"User-Agent":      []string{"test-client"},
						"X-Forwarded-For": []string{ip},
					},
					Cookies: map[string][]string{
						"cookie_name":         []string{"cookie_value_1", "cookie_value_2"},
						"another_cookie_name": []string{"cookie_value_=strange"},
					},
				},
				responseStatus: http.StatusInternalServerError,
				responseBody:   "An error occurred",
			},
		},
		"GET /forbidden": {
			method: "GET",
			url:    "/forbidden",
			body:   "",
			expects: &expects{
				recordLog: &trackrequest.RequestLog{
					IdentityID: "jo@nne.my",
					Method:     "GET",
					URL:        "/forbidden",
					IP:         ip,
					Body:       []byte(""),
					Headers: map[string][]string{
						"User-Agent":      []string{"test-client"},
						"X-Forwarded-For": []string{ip},
					},
					Cookies: map[string][]string{
						"cookie_name":         []string{"cookie_value_1", "cookie_value_2"},
						"another_cookie_name": []string{"cookie_value_=strange"},
					},
				},
				responseStatus: http.StatusForbidden,
				responseBody:   "Forbidden",
			},
		},
	}

	tr := trackrequest.NewManager(&trackrequest.Manager{})
	tr.OnRequest = func(log *trackrequest.RequestLog) {
		test := tests[log.Method+" "+log.URL]
		assert.Equal(t, "", log.IdentityID)
		assert.Equal(t, test.expects.recordLog.Method, log.Method)
		assert.Equal(t, test.expects.recordLog.URL, log.URL)
		assert.Equal(t, test.expects.recordLog.IP, log.IP)
		assert.Equal(t, test.expects.recordLog.Body, log.Body)
		assert.Equal(t, time.Duration(0), log.TimeTaken)
		for key, expectedHeader := range test.expects.recordLog.Headers {
			logHeader := log.Headers[key]
			assert.Subset(t, expectedHeader, logHeader)
		}
		for key, expectedCookie := range test.expects.recordLog.Cookies {
			cookie := log.Cookies[key]
			assert.Subset(t, expectedCookie, cookie)
		}
	}
	tr.OnRequestComplete = func(log *trackrequest.RequestLog) {
		test := tests[log.Method+" "+log.URL]
		assert.Equal(t, test.expects.recordLog.IdentityID, log.IdentityID)
		assert.Equal(t, test.expects.recordLog.Method, log.Method)
		assert.Equal(t, test.expects.recordLog.URL, log.URL)
		assert.Equal(t, test.expects.recordLog.IP, log.IP)
		assert.Equal(t, test.expects.recordLog.Body, log.Body)
		assert.True(t, log.TimeTaken > time.Duration(0))
		for key, expectedHeader := range test.expects.recordLog.Headers {
			logHeader := log.Headers[key]
			assert.Subset(t, expectedHeader, logHeader)
		}
		for key, expectedCookie := range test.expects.recordLog.Cookies {
			cookie := log.Cookies[key]
			assert.Subset(t, expectedCookie, cookie)
		}
	}
	tr.OnError = func(err error) {
		assert.NoError(t, err)
	}

	secureMiddleware := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tr.Identify = func(log trackrequest.RequestLog) string {
				return "jo@nne.my"
			}
			if r.URL.Path == "/forbidden" {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}

	ts := httptest.NewServer(tr.Middleware(secureMiddleware(testHandler())))
	defer ts.Close()

	for _, tc := range tests {
		// setup URL
		var u bytes.Buffer
		u.WriteString(string(ts.URL))
		u.WriteString(tc.url)

		// setup client
		client := &http.Client{}
		req, err := http.NewRequest(tc.method, u.String(), bytes.NewReader([]byte(tc.body)))
		assert.NoError(t, err)

		// set headers
		req.Header.Add("User-Agent", "test-client")
		req.Header.Add("X-Forwarded-For", ip)

		// set cookies
		req.AddCookie(&http.Cookie{
			Name:  "cookie_name",
			Value: "cookie_value_1",
		})
		req.AddCookie(&http.Cookie{
			Name:  "cookie_name",
			Value: "cookie_value_2",
		})
		req.AddCookie(&http.Cookie{
			Name:  "another_cookie_name",
			Value: "cookie_value_=strange",
		})

		// send request
		res, err := client.Do(req)
		assert.NoError(t, err)

		if res != nil {
			defer res.Body.Close()
		}

		bodyBytes, err := ioutil.ReadAll(res.Body)
		assert.Equal(t, tc.expects.responseStatus, res.StatusCode)
		assert.Equal(t, tc.expects.responseBody, strings.TrimSpace(string(bodyBytes[:])))
	}
}
