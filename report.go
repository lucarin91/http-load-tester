package main

import (
	"fmt"
	"math"
	"time"
)

type Report struct {
	Fail      uint64
	Success   uint64
	Slowest   time.Duration
	Fastest   time.Duration
	Total     time.Duration
	Average   time.Duration
	ReqPerSec float64
}

func NewReport() Report {
	return Report{
		Fail:      0,
		Success:   0,
		Slowest:   time.Duration(0),
		Fastest:   time.Duration(math.MaxInt64),
		Total:     time.Duration(0),
		Average:   time.Duration(0),
		ReqPerSec: 0.0,
	}
}

func (r1 *Report) Merge(r2 Report) {
	r1.Fail += r2.Fail
	r1.Success += r2.Success
	r1.Slowest = time.Duration(math.Max(float64(r1.Slowest), float64(r2.Slowest)))
	r1.Fastest = time.Duration(math.Min(float64(r1.Fastest), float64(r2.Fastest)))
	r1.Total += r2.Total
	// Average and ReqPerSec cannot be combine
}

func (r *Report) AddTime(d time.Duration) {
	r.Total += d
	r.Slowest = time.Duration(math.Max(float64(r.Slowest), float64(d)))
	r.Fastest = time.Duration(math.Min(float64(r.Fastest), float64(d)))
}

func (r *Report) Finalize() error {
	if r.Success == 0 || r.Total == 0 {
		return fmt.Errorf("finalize a not inizialized Report")
	}
	r.Average = time.Duration(uint64(r.Total) / r.Success)
	r.ReqPerSec = float64(r.Success) / r.Total.Seconds()
	return nil
}
