package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

type cmdArgs struct {
	workers  uint64
	requests uint64
	duration string
	url      string
}

func parseArgs() (cmdArgs, error) {
	var args cmdArgs
	flag.Uint64Var(&args.workers, "w", 50, "number of workers to run concurrently")
	flag.Uint64Var(&args.requests, "n", 200, "number of requests to run")
	flag.StringVar(&args.duration, "z", "", "duration of application to send requests")
	flag.Parse()

	if flag.NArg() > 1 {
		return args, fmt.Errorf("too many arguments")
	}
	if flag.NArg() < 1 {
		return args, fmt.Errorf("url not found")
	}
	args.url = flag.Arg(0)

	return args, nil
}

func main() {
	args, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// TODO: handle CTRL-C signal to graceful shoutdown and print report

	if err := run(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args cmdArgs) error {
	rep, err := func() (Report, error) {
		if "" != args.duration {
			d, err := time.ParseDuration(args.duration)
			if err != nil {
				return Report{}, err
			}
			return WithDuration(args.url, args.workers, d)
		}
		return WithRequests(args.url, args.workers, args.requests)
	}()
	if err != nil {
		return err
	}
	fmt.Printf(`
Summary:
  Requests:     %d
  Slowest:      %s
  Fastest:      %s
  Average:      %s
  Requests/sec: %.2f
`, rep.Requests, rep.Slowest, rep.Fastest, rep.Average, rep.ReqPerSec)
	return nil
}
