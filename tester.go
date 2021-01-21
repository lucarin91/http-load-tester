package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type requestChan chan struct{}
type resultChan chan Report

func WithRequests(u string, w uint64, n uint64) (Report, error) {
	if n < w {
		return Report{}, fmt.Errorf("number of requests cannot be less then worker")
	}

	req := numChan(n, w)
	rep := spawnAndWait(u, w, req)

	return rep, nil
}

func WithDuration(u string, w uint64, d time.Duration) (Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	req := cancelChan(ctx, w)

	rep := spawnAndWait(u, w, req)

	return rep, nil
}

func spawnAndWait(u string, w uint64, req requestChan) Report {
	res := make([]resultChan, 0, w)
	for i := w; i > 0; i-- {
		fmt.Println("M: start worker")
		res = append(res, worker(u, req))
	}
	var rep Report
	for r := range merge(res...) {
		rep.Merge(r)
	}
	fmt.Println("M: done")
	return rep
}

func worker(u string, req requestChan) resultChan {
	res := make(resultChan)
	go func() {
		var rep Report
		for range req {
			r, err := http.Get(u)
			if err != nil {
				rep.Fail++
			} else {
				rep.Success++
				r.Body.Close()
			}
		}
		fmt.Println("W: done")
		res <- rep
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

func numChan(n, cap uint64) requestChan {
	req := make(requestChan, cap)
	go func() {
		for i := n; i > 0; i-- {
			req <- struct{}{}
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
