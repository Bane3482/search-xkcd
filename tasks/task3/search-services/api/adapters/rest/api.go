package rest

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"yadro.com/course/api/core"
)

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/ping" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		replies := core.PingResponse{
			Replies: make(map[string]string),
		}

		for name, pinger := range pingers {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err := pinger.Ping(ctx)

			if err != nil {
				replies.Replies[name] = "unavailable"
			} else {
				replies.Replies[name] = "ok"
			}
		}

		body, err := json.Marshal(replies)

		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if _, err = w.Write(body); err != nil {
			log.Error("write response", slog.Any("error", err))
		}
	}
}

func NewWordsHandler(log *slog.Logger, normalizer core.Normalizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phrase := r.URL.Query().Get("phrase")
		if r.Method != http.MethodGet || phrase == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := normalizer.Norm(context.Background(), phrase)

		if err != nil {
			slog.Error("normalize", slog.Any("error", err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := core.WordsResponse{
			Words: result,
			Total: len(result),
		}

		body, err := json.Marshal(resp)

		if err != nil {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if _, err = w.Write(body); err != nil {
			log.Error("write response", slog.Any("error", err))
		}
	}
}
