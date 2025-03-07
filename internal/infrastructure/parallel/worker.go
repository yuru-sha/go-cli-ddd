package parallel

import (
	"context"
	"sync"

	"golang.org/x/sync/errgroup"
)

// WorkerPool は並列処理を行うワーカープールを表します
type WorkerPool struct {
	concurrency int
}

// NewWorkerPool は指定された並列数のワーカープールを作成します
func NewWorkerPool(concurrency int) *WorkerPool {
	if concurrency <= 0 {
		concurrency = 1
	}
	return &WorkerPool{
		concurrency: concurrency,
	}
}

// Process は指定されたタスクを並列に処理します
func (wp *WorkerPool) Process(ctx context.Context, tasks []func() error) error {
	g, ctx := errgroup.WithContext(ctx)

	// セマフォとして使用するチャネル
	sem := make(chan struct{}, wp.concurrency)

	for _, task := range tasks {
		task := task // ゴルーチン内で使用するためにローカル変数にコピー

		select {
		case <-ctx.Done():
			return ctx.Err()
		case sem <- struct{}{}: // セマフォを取得
			g.Go(func() error {
				defer func() { <-sem }() // セマフォを解放
				return task()
			})
		}
	}

	return g.Wait()
}

// ProcessWithResults は指定されたタスクを並列に処理し、結果を収集します
func ProcessWithResults[T any](wp *WorkerPool, ctx context.Context, tasks []func() (T, error)) ([]T, error) {
	g, ctx := errgroup.WithContext(ctx)

	// セマフォとして使用するチャネル
	sem := make(chan struct{}, wp.concurrency)

	// 結果を格納するスライス
	results := make([]T, len(tasks))

	// ミューテックス
	var mu sync.Mutex

	for i, task := range tasks {
		i, task := i, task // ゴルーチン内で使用するためにローカル変数にコピー

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case sem <- struct{}{}: // セマフォを取得
			g.Go(func() error {
				defer func() { <-sem }() // セマフォを解放

				result, err := task()
				if err != nil {
					return err
				}

				mu.Lock()
				results[i] = result
				mu.Unlock()

				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return results, nil
}
