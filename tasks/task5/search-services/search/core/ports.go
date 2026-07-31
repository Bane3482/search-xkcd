package core

import "context"

type Words interface {
	Norm(ctx context.Context, phrase string) ([]string, error)
}

type Updater interface {
	Get(ctx context.Context) ([]*Comics, error)
}

type Searcher interface {
	Search(ctx context.Context, phrase string, limit int) ([]*Comics, error)
}
