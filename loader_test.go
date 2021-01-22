package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestNum(t *testing.T) {
	var count uint64
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count++
	}))
	defer ts.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	count = 0
	rep, err := WithRequests(ctx, ts.URL, 1, 10)
	if err != nil {
		t.Errorf("get %v, want no error", err)
	}
	if 10 != count || count != rep.Requests {
		t.Errorf("get %v, want %v", rep.Requests, count)
	}

	count = 0
	rep, err = WithRequests(ctx, ts.URL, 5, 10)
	if err != nil {
		t.Errorf("get %v, want no error", err)
	}
	if 10 != count || count != rep.Requests {
		t.Errorf("get %v, want %v", rep.Requests, count)
	}
}
