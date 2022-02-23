package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pyroscope-io/client/pyroscope"
)

func fib(n int64) int64 {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func run() {
	fmt.Println(fib(50))
}

func main() {
	if os.Getenv("PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING") != "" {
		cfg := pyroscope.Config{
			ApplicationName: "fibonacci-go-all-push",
			ServerAddress:   "http://ingester:4040",
		}
		p, err := pyroscope.Start(cfg)
		if err != nil {
			log.Fatal(err)
		}
		run()
		if err := p.Stop(); err != nil {
			log.Fatal(err)
		}
	} else {
		run()
	}
}
