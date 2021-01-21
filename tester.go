package main

import (
	"fmt"
	"net/http"
	"time"
)

type Report struct {
	Fail    uint64
	Success uint64
}

func (r Report) String() string {
	return fmt.Sprintf("Fail: %v, Success: %v", r.Fail, r.Success)
}

func WithRequests(u string, w uint64, n uint64) (Report, error) {
	// TODO: parallelize with 'w' worker
	var rep Report
	if n < w {
		return rep, fmt.Errorf("number of requests cannot be less then worker")
	}
	for i := n; i > 0; i-- {
		r, err := http.Get(u)
		if err != nil {
			rep.Fail++
		} else {
			rep.Success++
			r.Body.Close()
		}
	}
	return rep, nil
}

func WithDuration(u string, w uint64, d time.Duration) (Report, error) {
	// TODO: parallelize with 'w' worker
	var rep Report
	after := time.After(d)
	for {
		select {
		case <-after:
			return rep, nil
		default:
			r, err := http.Get(u)
			if err != nil {
				rep.Fail++
			} else {
				rep.Success++
				r.Body.Close()
			}
		}
	}
}
