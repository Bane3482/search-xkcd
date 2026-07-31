package utils

import (
	"fmt"
	"net/http"
	"regexp"

	"go.uber.org/zap"
)

type (
	HelloHandler struct{}
)

var (
	HelloPing = regexp.MustCompile(`^/ping$`)
	HelloUser = regexp.MustCompile(`^/hello$`)
)

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	switch {
	case HelloPing.MatchString(r.URL.Path):
		logger.Info("ping request", zap.String("URI", r.URL.RequestURI()))

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("pong\n")); err != nil {
			return
		}
		return
	case HelloUser.MatchString(r.URL.Path):
		logger.Info("hello request", zap.String("URI", r.URL.RequestURI()))

		name := r.URL.Query().Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write([]byte("empty name\n")); err != nil {
				logger.Error("write empty message", zap.Error(err))
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, "Hello, %s!\n", name); err != nil {
			logger.Error("response writer", zap.Error(err))
			return
		}
		return
	default:
		logger.Info("bad request")
		w.WriteHeader(http.StatusBadRequest)
	}
}
