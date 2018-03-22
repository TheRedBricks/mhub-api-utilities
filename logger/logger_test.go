package logger_test

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/TheRedBricks/mhub-api-utilities/logger"
	"github.com/stretchr/testify/assert"
)

func testHandler() http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ok" {
			io.WriteString(w, "OK")
			return
		}
		if r.URL.Path == "/not_found" {
			http.Error(w, "Page Not Found", http.StatusNotFound)
			return
		}
		if r.URL.Path == "/error" {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		panic("URL Not handled")
	}
	return http.HandlerFunc(fn)
}

func TestLogger(t *testing.T) {
	timeRegex := "\\d{4}\\/\\d{2}\\/\\d{2}\\s\\d{2}:\\d{2}:\\d{2}"
	tests := []struct {
		description string
		url         string
		expectedLog string
	}{
		{
			description: "Test OK Log",
			url:         "/ok",
			expectedLog: "^" + timeRegex + " GET \\/ok 200\\s.*\\n$",
		},
		{
			description: "Test Not Found Log",
			url:         "/not_found",
			expectedLog: "^" + timeRegex + " GET \\/not_found 404\\s.*\\n$",
		},
		{
			description: "Test Internal Server Error Log",
			url:         "/error",
			expectedLog: "^" + timeRegex + " GET \\/error 500\\s.*\\n$",
		},
	}

	ts := httptest.NewServer(logger.Middleware(testHandler()))
	defer ts.Close()

	for _, tc := range tests {
		var logBuf bytes.Buffer
		log.SetOutput(&logBuf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		var u bytes.Buffer
		u.WriteString(string(ts.URL))
		u.WriteString(tc.url)

		res, err := http.Get(u.String())
		assert.NoError(t, err)

		if res != nil {
			defer res.Body.Close()
		}

		assert.Regexp(t, tc.expectedLog, logBuf.String(), tc.description)
	}
}
