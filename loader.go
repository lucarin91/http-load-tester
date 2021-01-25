package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type requestChan chan struct{}
type resultChan chan Result

func WithRequests(ctx context.Context, u string, w uint64, n uint64) (Report, error) {
	if n < w {
		return Report{}, fmt.Errorf("number of requests cannot be less then worker")
	}
	g, ctx := errgroup.WithContext(ctx)

	req := make(requestChan, w)
	g.Go(numChan(ctx, req, n))

	rep, err := spawnAndWait(ctx, g, u, w, req)
	if err != nil {
		return rep, err
	}
	return rep, nil
}

func WithDuration(ctx context.Context, u string, w uint64, d time.Duration) (Report, error) {
	g, ctx := errgroup.WithContext(ctx)
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	req := make(requestChan, w)
	g.Go(cancelChan(ctx, req))

	rep, err := spawnAndWait(ctx, g, u, w, req)
	if err != nil {
		return rep, err
	}
	return rep, nil
}

func spawnAndWait(ctx context.Context, g *errgroup.Group, u string, w uint64, req requestChan) (Report, error) {
	res := make(resultChan, w+1)

	stat := NewStatistics()
	fmt.Println("M: start workers")
	for i := w; i > 0; i-- {
		g.Go(worker(ctx, req, res, u))
	}

	go func() {
		g.Wait()
		close(res)
	}()

	fmt.Println("M: wait workers")
	for r := range res {
		stat.Add(r)
	}

	fmt.Println("M: check errors and finalize")
	err := g.Wait()
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		return Report{}, fmt.Errorf("errgroup: %v", err)
	}
	rep, err := stat.Finalize()
	if err != nil {
		return rep, err
	}

	fmt.Println("M: done")
	return rep, nil
}

func worker(ctx context.Context, in requestChan, out resultChan, u string) func() error {
	return func() error {
		for range in {
			start := time.Now()
			r, err := http.Get(u)
			if err != nil {
				// fmt.Println("W: error")
				return err
			}
			r.Body.Close()
			d := time.Since(start)
			select {
			case out <- Result{dur: d, code: r.StatusCode}:
			case <-ctx.Done():
				// fmt.Println("W: ctx done")
				return ctx.Err()
			}
		}
		// fmt.Println("W: done")
		return nil
	}
}

func cancelChan(ctx context.Context, out requestChan) func() error {
	return func() error {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				// fmt.Println("S: ctx done")
				return ctx.Err()
			case out <- struct{}{}:
			}
		}
	}
}

func numChan(ctx context.Context, out requestChan, n uint64) func() error {
	return func() error {
		defer close(out)
		for i := n; i > 0; i-- {
			select {
			case <-ctx.Done():
				// fmt.Println("S: ctx done")
				return ctx.Err()
			case out <- struct{}{}:
			}
		}
		// fmt.Println("S: done")
		return nil
	}
}
