package logger

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
)

type loggedResponse struct {
	w      http.ResponseWriter
	r      *http.Request
	url    string
	status int
}

func (w *loggedResponse) Flush() {
	if wf, ok := w.w.(http.Flusher); ok {
		wf.Flush()
	}
}

func (w *loggedResponse) Header() http.Header         { return w.w.Header() }
func (w *loggedResponse) Write(d []byte) (int, error) { return w.w.Write(d) }

func (w *loggedResponse) WriteHeader(status int) {
	w.status = status
	w.w.WriteHeader(status)
}

// Middleware logger for http mux
func Middleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggedResponse{w, r, r.URL.String(), 200}
		h.ServeHTTP(lrw, r)

		status := strconv.Itoa(lrw.status)
		colorStatus := color.RedString(status)
		if lrw.status < 300 {
			colorStatus = color.GreenString(status)
		} else if lrw.status < 400 {
			colorStatus = color.CyanString(status)
		}
		log.Printf("%s %s %s %s", r.Method, r.URL, colorStatus, time.Since(start))
	}
	return http.HandlerFunc(fn)
}
