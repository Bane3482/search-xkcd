package core

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sort"
)

type Service struct {
	log        *slog.Logger
	words      Words
	update     Updater
	concurrent int
}

func NewService(
	log *slog.Logger,
	words Words,
	update Updater,
	concurrent int,
) (*Service, error) {
	if concurrent < 1 {
		return nil, fmt.Errorf("wrong concurrency specified: %d", concurrent)
	}
	return &Service{
		log:        log,
		words:      words,
		update:     update,
		concurrent: concurrent,
	}, nil
}

type Task func(context.Context) error

func (s *Service) Search(ctx context.Context, phrase string, limit int) ([]*Comics, error) {
	norm, err := s.words.Norm(ctx, phrase)

	if err != nil {
		s.log.Error("service norm", "error", err)
		return nil, err
	}

	comics, err := s.update.Get(ctx)

	if err != nil {
		s.log.Error("service get", "error", err)
		return nil, err
	}

	sort.Slice(comics, func(i, j int) bool {
		return comics[i].URL < comics[j].URL
	})

	result := make([]*Comics, 0)

	for _, c := range comics {
		if slices.ContainsFunc(c.Words, func(el string) bool {
			return slices.Contains(norm, el)
		}) {
			result = append(result, c)
			if len(result) == limit {
				break
			}
		}
	}

	return result, nil
}
