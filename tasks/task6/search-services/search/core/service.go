package core

import (
	"cmp"
	"context"
	"log/slog"
	"maps"
	"slices"
)

type Service struct {
	log   *slog.Logger
	db    DB
	words Words
}

func NewService(
	log *slog.Logger,
	db DB,
	words Words,
) (*Service, error) {
	return &Service{
		log:   log,
		db:    db,
		words: words,
	}, nil
}

type Task func(context.Context) error

func (s *Service) Search(ctx context.Context, phrase string, limit int) ([]ComicsInfo, error) {
	return s.search(s.db.Search, ctx, phrase, limit)
}

func (s *Service) ISearch(ctx context.Context, phrase string, limit int) ([]ComicsInfo, error) {
	return s.search(s.db.ISearch, ctx, phrase, limit)
}

type searcher func(context.Context, string) ([]int, error)

func (s *Service) search(searcher searcher, ctx context.Context, phrase string, limit int) ([]ComicsInfo, error) {
	keywords, err := s.words.Norm(ctx, phrase)

	if err != nil {
		s.log.Error("search service norm", "error", err)
		return nil, err
	}

	count := make(map[int]int)

	for _, keyword := range keywords {
		IDs, err := searcher(ctx, keyword)

		if err != nil {
			s.log.Error("search service search", "error", err)
			return nil, err
		}

		for _, id := range IDs {
			count[id]++
		}
	}

	relevant := slices.SortedFunc(maps.Keys(count), func(i, j int) int {
		return cmp.Compare(count[j], count[i])
	})

	limit = min(len(relevant), limit)

	result := make([]ComicsInfo, limit)

	for i := 0; i < limit; i++ {
		info, err := s.db.Get(ctx, relevant[i])

		if err != nil {
			s.log.Error("search service get", "error", err)
			return nil, err
		}

		result[i] = info
	}

	return result, nil
}
