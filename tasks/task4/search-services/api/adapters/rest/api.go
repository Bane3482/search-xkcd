package rest

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"yadro.com/course/api/core"
)

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
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
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

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
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := updater.Drop(context.Background()); err != nil {
			log.Error("drop handle drop", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
