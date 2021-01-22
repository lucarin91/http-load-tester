package main

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"
)

const host = "localhost:8080"
const path = "/pZ3hLHse"

func TestRequestNum(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = host
	defer srv.Shutdown(ctx)
	var count uint64
	http.HandleFunc("/pZ3hLHse", func(w http.ResponseWriter, r *http.Request) {
		count++
	})
	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("http server: %v", err)
		}
	}()
	time.Sleep(1 * time.Second)

	count = 0
	rep, err := WithRequests(ctx, "http://"+host+path, 1, 10)
	if err != nil {
		t.Errorf("get %v, want no error", err)
	}
	if 10 != count || count != rep.Requests {
		t.Errorf("get %v, want %v", rep.Requests, count)
	}

	count = 0
	rep, err = WithRequests(ctx, "http://"+host+path, 5, 10)
	if err != nil {
		t.Errorf("get %v, want no error", err)
	}
	if 10 != count || count != rep.Requests {
		t.Errorf("get %v, want %v", rep.Requests, count)
	}
}
