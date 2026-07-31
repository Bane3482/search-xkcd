package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"yadro.com/course/api/core"
)

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			log.Error("ping handle marshal", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err = w.Write(body); err != nil {
			log.Error("ping handle write response", slog.Any("error", err))
		}
	}
}

func NewUpdateHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		if err := updater.Update(ctx); err != nil {
			if errors.Is(err, core.ErrAlreadyExists) {
				log.Info("update handle status accepted")
				w.WriteHeader(http.StatusAccepted)
			} else {
				log.Error("update handle update", "error", err)
				w.WriteHeader(http.StatusBadRequest)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewUpdateStatsHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := updater.Stats(context.Background())

		if err != nil {
			log.Error("stats handle get stats", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(resp)

		if err != nil {
			log.Error("stats handle marshal", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err = w.Write(body); err != nil {
			log.Error("stats handle write response", "error", err)
		}
	}
}

func NewUpdateStatusHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := updater.Status(context.Background())

		if err != nil {
			log.Error("status handle get status", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp := core.UpdateStatus{
			Status: status,
		}

		body, err := json.Marshal(resp)

		if err != nil {
			log.Error("status handle marshal", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, err = w.Write(body); err != nil {
			log.Error("status handle write response", "error", err)
		}
	}
}

func NewDropHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := updater.Drop(context.Background()); err != nil {
			log.Error("drop handle drop", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func NewSearchHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		phrase := query.Get("phrase")

		var limit int

		limitStr := query.Get("limit")

		if limitStr != "" {
			var err error
			if limit, err = strconv.Atoi(limitStr); err != nil {
				log.Error("search handler parse limit", "error", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if limit < 1 {
				log.Error("search handler wrong limit", "error", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			limit = 10
		}

		reply, err := searcher.Search(context.Background(), phrase, limit)

		if err != nil {
			log.Error("search handler search", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var body []byte

		body, err = json.Marshal(&reply)

		if err != nil {
			log.Error("search handler marshal", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Write(body)
	}
}
