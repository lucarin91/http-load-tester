package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type requestChan chan struct{}
type resultChan chan Result

func WithRequests(ctx context.Context, u string, w uint64, n uint64) (Report, error) {
	if n < w {
		return Report{}, fmt.Errorf("number of requests cannot be less then worker")
	}

	req := numChan(ctx, n, w)
	rep, err := spawnAndWait(u, w, req)
	if err != nil {
		return rep, err
	}
	return rep, nil
}

func WithDuration(ctx context.Context, u string, w uint64, d time.Duration) (Report, error) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	req := cancelChan(ctx, w)
	rep, err := spawnAndWait(u, w, req)
	if err != nil {
		return rep, err
	}
	return rep, nil
}

func spawnAndWait(u string, w uint64, req requestChan) (Report, error) {
	res := make([]resultChan, 0, w)
	for i := w; i > 0; i-- {
		fmt.Println("M: start worker")
		res = append(res, worker(u, req))
	}
	stat := NewStatistics()
	for r := range merge(res...) {
		stat.Add(r)
	}
	rep, err := stat.Finalize()
	if err != nil {
		return rep, err
	}
	fmt.Println("M: done")
	return rep, nil
}

func worker(u string, req requestChan) resultChan {
	res := make(resultChan)
	go func() {
		for range req {
			start := time.Now()
			r, err := http.Get(u)
			if err != nil {
				// TODO: manage error with errGroup
			}
			r.Body.Close()
			d := time.Since(start)
			res <- Result{dur: d, code: r.StatusCode}
		}
		fmt.Println("W: done")
		close(res)
	}()
	return res
}

func cancelChan(ctx context.Context, cap uint64) requestChan {
	req := make(requestChan, cap)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(req)
				return
			default:
				req <- struct{}{}
			}
		}
	}()
	return req
}

func numChan(ctx context.Context, n, cap uint64) requestChan {
	req := make(requestChan, cap)
	go func() {
		for i := n; i > 0; i-- {
			select {
			case <-ctx.Done():
				close(req)
				return
			default:
				req <- struct{}{}
			}
		}
		close(req)
	}()
	return req
}

func merge(cs ...resultChan) resultChan {
	var wg sync.WaitGroup
	out := make(resultChan)

	wg.Add(len(cs))
	for _, c := range cs {
		go func(c resultChan) {
			for e := range c {
				out <- e
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
