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
			err := init.db.Drop(ctx)

			if err != nil {
				init.log.Error("initialize index delete", "error", err)
				break
			}

			keywords, err := init.db.Keywords(ctx)

			if err != nil {
				init.log.Error("initialize index get keywords", "error", err)
				break
			}

			for _, keyword := range keywords {
				ids, err := init.db.Search(ctx, keyword)

				if err != nil {
					init.log.Error("initialize index search", "error", err, "keyword", keyword)
					continue
				}

				if err = init.db.Insert(ctx, keyword, ids); err != nil {
					init.log.Error("initialize index insert", "error", err, "keywords", keyword, "IDs", ids)
				}
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
