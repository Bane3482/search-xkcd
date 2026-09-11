package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/VictoriaMetrics/metrics"
	"yadro.com/course/api/core"
)

func encodeReply(w io.Writer, reply any) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(reply); err != nil {
		return fmt.Errorf("could not encode comics: %v", err)
	}
	return nil
}

func NewMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.WritePrometheus(w, true)

		w.WriteHeader(http.StatusOK)
	}
}

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}

func NewPingHandler(log *slog.Logger, pingers map[string]core.Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := PingResponse{
			Replies: make(map[string]string),
		}

		for name, pinger := range pingers {
			err := pinger.Ping(r.Context())

			if err != nil {
				reply.Replies[name] = "unavailable"
				log.Error("one of services is not available", "service", name, "error", err)
			} else {
				reply.Replies[name] = "ok"
			}
		}

		if err := encodeReply(w, reply); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

type Authenticator interface {
	Login(user, password string) (string, error)
}

type AuthInfo struct {
	User     string `json:"name"`
	Password string `json:"password"`
}

func NewLoginHandler(log *slog.Logger, auth Authenticator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		info := AuthInfo{}

		if err := decoder.Decode(&info); err != nil {
			log.Error("login decode info", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		token, err := auth.Login(info.User, info.Password)

		if err != nil {
			log.Error("auth login", "error", err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		w.WriteHeader(http.StatusOK)

		if _, err = fmt.Fprintf(w, "%s", token); err != nil {
			log.Error("write token", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NewUpdateHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := updater.Update(r.Context()); err != nil {
			if errors.Is(err, core.ErrAlreadyExists) {
				log.Info("update handle status accepted")
				w.WriteHeader(http.StatusAccepted)
			}
			log.Error("update handle update", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type UpdateStats struct {
	WordsTotal    int `json:"words_total"`
	WordsUnique   int `json:"words_unique"`
	ComicsFetched int `json:"comics_fetched"`
	ComicsTotal   int `json:"comics_total"`
}

func NewUpdateStatsHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stats, err := updater.Stats(context.Background())

		if err != nil {
			log.Error("stats handle get stats", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reply := UpdateStats{
			WordsTotal:    stats.WordsTotal,
			WordsUnique:   stats.WordsUnique,
			ComicsFetched: stats.ComicsFetched,
			ComicsTotal:   stats.ComicsTotal,
		}

		if err := encodeReply(w, reply); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

type UpdateStatus struct {
	Status string `json:"status"`
}

func NewUpdateStatusHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := updater.Status(context.Background())

		if err != nil {
			log.Error("status handle get status", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reply := UpdateStatus{
			Status: string(status),
		}

		if err := encodeReply(w, reply); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

func NewDropHandler(log *slog.Logger, updater core.Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := updater.Drop(r.Context()); err != nil {
			log.Error("drop handle drop", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type Comics struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type ComicsReply struct {
	Comics []Comics `json:"comics"`
	Total  int      `json:"total"`
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

		comics, err := searcher.Search(r.Context(), phrase, limit)

		if err != nil {
			log.Error("search handler search", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reply := ComicsReply{
			Comics: make([]Comics, 0, len(comics)),
			Total:  len(comics),
		}

		for _, c := range comics {
			reply.Comics = append(reply.Comics, Comics{ID: c.ID, URL: c.URL})
		}

		if err := encodeReply(w, reply); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}

func NewSearchIndexHandler(log *slog.Logger, searcher core.Searcher) http.HandlerFunc {
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

		comics, err := searcher.SearchIndex(r.Context(), phrase, limit)

		if err != nil {
			log.Error("search handler search", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		reply := ComicsReply{
			Comics: make([]Comics, 0, len(comics)),
			Total:  len(comics),
		}

		for _, c := range comics {
			reply.Comics = append(reply.Comics, Comics{ID: c.ID, URL: c.URL})
		}

		if err := encodeReply(w, reply); err != nil {
			log.Error("cannot encode reply", "error", err)
		}
	}
}
