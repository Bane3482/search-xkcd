package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func WithMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, code: http.StatusOK}

		startTime := time.Now()
		next.ServeHTTP(rw, r)

		metrics.GetOrCreateHistogram(
			fmt.Sprintf(`http_request_duration_seconds{url=%q,status="%d"}`, r.URL.Path, rw.code),
		).UpdateDuration(startTime)
	})
}
