package main

import (
	"fmt"
	"testing"
	"time"
)

func TestFinalize(t *testing.T) {
	var hashtests = []struct {
		inSuccess    uint64
		inTotal      string
		outAverage   string
		outReqPerSec float64
	}{
		{10, "1s", "100ms", 10},
		{10, "10s", "1s", 1},
		{10, "100s", "10s", 0.1},
	}
	for _, tt := range hashtests {
		name := fmt.Sprintf("%v-%v", tt.inSuccess, tt.inTotal)
		t.Run(name, func(t *testing.T) {
			r := NewReport()
			r.Success = tt.inSuccess
			r.Total, _ = time.ParseDuration(tt.inTotal)
			err := r.Finalize()
			if err != nil {
				t.Errorf("got %v, want %q", err, "nill")
			}
			if r.Average.String() != tt.outAverage {
				t.Errorf("got %q, want %q", r.Average.String(), tt.outAverage)
			}
			if r.ReqPerSec != tt.outReqPerSec {
				t.Errorf("got %v, want %v", r.ReqPerSec, tt.outReqPerSec)
			}
		})
	}
}
