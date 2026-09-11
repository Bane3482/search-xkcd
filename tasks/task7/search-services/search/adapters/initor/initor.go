package initor

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"yadro.com/course/search/adapters/db"
)

type Initor struct {
	db  *db.DB
	log *slog.Logger
}

func New(log *slog.Logger, db *db.DB) *Initor {
	return &Initor{
		db:  db,
		log: log,
	}
}

func (init *Initor) work(ctx context.Context, waitTime time.Duration) {
	ticker := time.NewTicker(waitTime)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := init.db.RebuildIndex(ctx)
			if err != nil {
				init.log.Error("initialize index", "error", err)
			}
		}
	}
}

func (init *Initor) Init(ctx context.Context, waitTime time.Duration) *sync.WaitGroup {
	wg := &sync.WaitGroup{}

	wg.Go(func() {
		init.work(ctx, waitTime)
	})

	return wg
}
