package main

import (
	"fmt"
	"math"
	"time"
)

type Statistics struct {
	n     uint64
	max   time.Duration
	min   time.Duration
	start time.Time
	total time.Duration
	codes map[int]uint64
}
type Report struct {
	Requests  uint64
	Slowest   time.Duration
	Fastest   time.Duration
	Average   time.Duration
	ReqPerSec float64
	Codes     map[int]uint64
}

type Result struct {
	dur  time.Duration
	code int
}

func NewStatistics() Statistics {
	return Statistics{
		start: time.Now(),
		min:   time.Duration(math.MaxInt64),
		codes: make(map[int]uint64),
	}
}

func (s *Statistics) Add(res Result) {
	s.total += res.dur
	s.n++
	s.max = time.Duration(math.Max(float64(s.max), float64(res.dur)))
	s.min = time.Duration(math.Min(float64(s.min), float64(res.dur)))
	s.codes[res.code]++
}

func (s *Statistics) Finalize() (Report, error) {
	if s.n == 0 || s.total == 0 {
		return Report{}, fmt.Errorf("finalize a not inizialized Report")
	}
	return Report{
		Requests:  s.n,
		Slowest:   s.max,
		Fastest:   s.min,
		Average:   time.Duration(uint64(s.total) / s.n),
		ReqPerSec: float64(s.n) / time.Since(s.start).Seconds(),
		Codes:     s.codes,
	}, nil
}
