package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type requestChan chan struct{}
type resultChan chan Report

func WithRequests(u string, w uint64, n uint64) (Report, error) {
	if n < w {
		return Report{}, fmt.Errorf("number of requests cannot be less then worker")
	}
	req := make(requestChan, n)
	res := make(resultChan)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := w; i > 0; i-- {
		fmt.Println("M: start worker")
		requestWorker(ctx, u, req, res)
	}

	var rep Report
	for i := w; i > 0; i-- {
		rep.Merge(<-res)
	}
	fmt.Println("M: done")
	return rep, nil
}

func WithDuration(u string, w uint64, d time.Duration) (Report, error) {
	res := make(resultChan)
	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	for i := w; i > 0; i-- {
		fmt.Println("M: start worker")
		durationWorker(ctx, u, res)
	}

	var rep Report
	for i := w; i > 0; i-- {
		rep.Merge(<-res)
	}
	fmt.Println("M: done")
	return rep, nil
}

func requestWorker(ctx context.Context, u string, req requestChan, res resultChan) {
	go func() {
		var rep Report
		for {
			select {
			case <-ctx.Done():
			default:
				fmt.Println("W: done")
				res <- rep
				return
			case req <- struct{}{}:
				workerRequest(u, &rep)
			}
		}
	}()
}

func durationWorker(ctx context.Context, u string, res resultChan) {
	go func() {
		var rep Report
		for {
			select {
			case <-ctx.Done():
				fmt.Println("W: done")
				res <- rep
				return
			default:
				workerRequest(u, &rep)
			}
		}
	}()
}

func workerRequest(u string, rep *Report) {
	r, err := http.Get(u)
	if err != nil {
		rep.Fail++
	} else {
		rep.Success++
		r.Body.Close()
	}
}
