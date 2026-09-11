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

type TaskFunc func(context.Context) error

type Task struct {
	ID   int
	Func TaskFunc
}

type TaskResult struct {
	ID    int
	Error error
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

	taskChan, resultChan := s.workerPool(ctx, s.concurrency)

	go func() {
		ptr := 0
		for cur := range maxNum {
			if ptr < len(nums) && nums[ptr] == cur+1 {
				ptr++
				continue
			}

			task := Task{
				ID:   cur + 1,
				Func: s.newTaskFunc(cur + 1),
			}

			select {
			case <-ctx.Done():
				return
			case taskChan <- task:
				s.log.Info("task sent", "comics_id", cur+1)
			}
		}
		close(taskChan)
	}()

	for result := range resultChan {
		if result.Error == nil {
			s.log.Info("task finished successfully", "comics_id", result.ID)
		} else {
			s.log.Error("task finished with error", "comics_id", result.ID, "error", result.Error)
		}
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

func (s *Service) newTaskFunc(id int) TaskFunc {
	return func(ctx context.Context) error {
		var comics Comics

		if id == 404 {
			comics = Comics{
				ID:    404,
				URL:   "invalid_url.com",
				Words: []string{"not-found"},
			}
		} else {
			info, err := s.xkcd.Get(ctx, id)

			if err != nil {
				return err
			}

			words, err := s.words.Norm(ctx, fmt.Sprintf("%s %s %s %s", info.Title, info.SafeTitle, info.Alt, info.Transcript))

			if err != nil {
				return err
			}

			comics = Comics{
				ID:    info.ID,
				URL:   info.URL,
				Words: words,
			}
		}

		return s.db.Add(ctx, comics)
	}
}

func (s *Service) workerPool(ctx context.Context, workers int) (chan<- Task, <-chan TaskResult) {
	wg := &sync.WaitGroup{}

	taskChan := make(chan Task, s.concurrency)
	resultChan := make(chan TaskResult, s.concurrency)

	for range workers {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case task, ok := <-taskChan:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case resultChan <- TaskResult{ID: task.ID, Error: task.Func(ctx)}:
					}
				}
			}
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	return taskChan, resultChan
}
