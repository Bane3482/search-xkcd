package core

import (
	"context"
	"time"
)

type Words interface {
	Norm(ctx context.Context, phrase string) ([]string, error)
}

type DB interface {
	Search(ctx context.Context, keyword string) ([]int, error)
	ISearch(ctx context.Context, keyword string) ([]int, error)
	Get(ctx context.Context, id int) (ComicsInfo, error)
}

type Updater interface {
	Get(ctx context.Context) ([]*Comics, error)
}

type Searcher interface {
	Search(ctx context.Context, phrase string, limit int) ([]ComicsInfo, error)
	ISearch(ctx context.Context, phrase string, limit int) ([]ComicsInfo, error)
}

type Initor interface {
	Init(ctx context.Context, waitTime time.Duration) error
}
