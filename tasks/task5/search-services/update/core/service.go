package core

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"sync"
	"sync/atomic"
)

type Service struct {
	log            *slog.Logger
	db             DB
	xkcd           XKCD
	words          Words
	isRuningUpdate atomic.Bool
	concurrency    int
}

func NewService(
	log *slog.Logger, db DB, xkcd XKCD, words Words, concurrency int,
) (*Service, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("wrong concurrency specified: %d", concurrency)
	}
	return &Service{
		log:         log,
		db:          db,
		xkcd:        xkcd,
		words:       words,
		concurrency: concurrency,
	}, nil
}

func (s *Service) Update(ctx context.Context) error {
	if !s.isRuningUpdate.CompareAndSwap(false, true) {
		return ErrAlreadyExists
	}

	defer s.isRuningUpdate.Store(false)

	nums, err := s.db.IDs(ctx)

	if err != nil {
		return err
	}

	maxNum, err := s.xkcd.LastID(ctx)

	if err != nil {
		return err
	}

	sort.Slice(nums, func(i, j int) bool {
		return nums[i] < nums[j]
	})

	taskChan, errChan := workerPool(ctx, s.concurrency)

	ptr := 0

	run := true

	for cur := range maxNum {
		if !run {
			break
		}
		if ptr < len(nums) && nums[ptr] == cur+1 {
			ptr++
			continue
		}
		s.log.Info("send task", "cur", cur+1)
		select {
		case <-ctx.Done():
			run = false
		case errNew, ok := <-errChan:
			if ok {
				err = errNew
			}
			run = false
		case taskChan <- Task(func(ctx context.Context) error {
			s.log.Info("start task", "cur", cur+1)
			var comics Comics

			var err error

			defer func() {
				if err != nil {
					s.log.Error("end task with error", "err", err)
				}
			}()

			if cur+1 == 404 {
				comics = Comics{
					ID:    404,
					URL:   "invalid_url.com",
					Words: []string{"not-found"},
				}
			} else {
				var info XKCDInfo

				info, err = s.xkcd.Get(ctx, cur+1)

				if err != nil {
					return err
				}

				var words []string

				words, err = s.words.Norm(ctx, fmt.Sprintf("%s %s %s %s", info.Title, info.SafeTitle, info.Alt, info.Transcript))

				if err != nil {
					return err
				}

				comics = Comics{
					ID:    info.ID,
					URL:   info.URL,
					Words: words,
				}
			}

			s.log.Error("done task", "cur", cur+1)

			return s.db.Add(ctx, comics)
		}):
			s.log.Info("sended succesfully")
		}
	}

	close(taskChan)

	if err != nil {
		return err
	}

	if err, ok := <-errChan; ok {
		return err
	}

	return nil
}

func (s *Service) Stats(ctx context.Context) (ServiceStats, error) {
	dbStats, err := s.db.Stats(ctx)

	if err != nil {
		return ServiceStats{}, err
	}

	count, err := s.xkcd.LastID(ctx)

	if err != nil {
		return ServiceStats{}, err
	}

	return ServiceStats{
		DBStats:     dbStats,
		ComicsTotal: count,
	}, nil
}

func (s *Service) Status(ctx context.Context) ServiceStatus {
	if s.isRuningUpdate.Load() {
		return StatusRunning
	}
	return StatusIdle
}

func (s *Service) Drop(ctx context.Context) error {
	return s.db.Drop(ctx)
}

type Task func(context.Context) error

func workerPool(ctx context.Context, workers int) (chan<- Task, <-chan error) {
	wg := &sync.WaitGroup{}

	taskChan := make(chan Task)

	errChan := make(chan error)

	end, cancel := context.WithCancel(context.Background())

	for range workers {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-end.Done():
					return
				case task, ok := <-taskChan:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case <-end.Done():
						return
					default:
						err := task(ctx)
						if err != nil {
							select {
							case errChan <- err:
							default:
							}
							cancel()
							return
						}
					}
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(errChan)
		cancel()
	}()

	return taskChan, errChan
}
