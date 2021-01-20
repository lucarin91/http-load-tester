package main

import "fmt"

type Report struct {
	Fail    uint64
	Success uint64
}

func (r1 *Report) Merge(r2 Report) {
	r1.Fail += r2.Fail
	r1.Success += r2.Success
}

func (r Report) String() string {
	return fmt.Sprintf("Fail: %v, Success: %v", r.Fail, r.Success)
}
