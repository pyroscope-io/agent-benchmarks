/*
   Runner for the the benchmark suite.
   It takes a sequence of paths to the actual benchmarks to run.

   Each benchmark is just a directory with a Dockerfile that speficies
   how to build the image of the program to benchmark.

   The program to benchmark should handle the PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING
   environment variable:
   - when the envvar is set, the benchmarked program should enable profiling.
   - when the envvar is not set, the benchmarked program should not enable profiling.

   The runner currently executes three different versions of the benchmarked program:
   - A program with profiling not enabled. This is the baseline.
   - A program with profiling enabled with a fast ingester available.
   - A program with profiling enabled with a slow ingester available.
   - A program with profiling enabled with an ingester unavailable.
   It will measure the time it takes to run all of them and compare the last two with the baseline.
   The benchmarked programs are run several times to have more samples and generate more reliable results.

   The runner takes care of the whole setup and teardown of each benchmark, including:
   - Building the docker images and containres for the ingester and benchmarked program.
   - Creating a network and connecting/disconnecting the ingester and benchmarked program.
   - Removing the containers, images and network when no longer needed.
*/
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"golang.org/x/perf/benchstat"
)

const (
	id           = "pyroscope-agent-benchmark"
	ingesterPath = "ingester"
	outputPath   = "results.txt"
	n            = 5
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Panicf("Unable to create docker client: %s", err)
	}
	var noprof, prof bytes.Buffer
	r := newRunner(cli)
	for _, agent := range os.Args[1:] {
		run(ctx, r, agent, &noprof, &prof)
	}
	c := &benchstat.Collection{
		Alpha:      0.05,
		AddGeoMean: false,
		DeltaTest:  benchstat.UTest,
	}
	c.AddConfig("no profiling", noprof.Bytes())
	c.AddConfig("profiling", prof.Bytes())
	f, err := os.Create(outputPath)
	if err != nil {
		log.Panicf("Unable to open results file: %s", err)
	}
	defer f.Close()
	benchstat.FormatText(f, c.Tables())
}

func run(ctx context.Context, r *runner, name string, noprof, prof *bytes.Buffer) {
	agent := agentName(name)
	log.Printf("Running %s benchmark", agent)
	if err := r.buildImage(ctx, ingesterPath, id+"/ingester"); err != nil {
		log.Panicf("Unable to create ingester image: %s", err)
	}

	if err := r.buildImage(ctx, name, id+"/benchmarked"); err != nil {
		log.Panicf("Unable to create benchmarked image: %s", err)
	}

	if err := r.createNetwork(ctx); err != nil {
		log.Panicf("Unable to create network: %s", err)
	}
	defer r.removeNetwork(ctx)

	// Fast ingester
	if err := r.createIngester(ctx, "fast"); err != nil {
		log.Panicf("Unable to create ingester container: %s", err)
	}
	defer r.removeIngester(ctx)

	if err := r.connectIngester(ctx); err != nil {
		log.Panicf("Unable to connect the ingester to the network: %s", err)
	}

	if err := r.startIngester(ctx); err != nil {
		log.Panicf("Unable to start ingester: %s", err)
	}

	defer r.removeBenchmarked(ctx)

	var fast []time.Duration
	for i := 0; i < n; i++ {
		if err := r.createBenchmarked(ctx, true); err != nil {
			log.Panicf("Unable to create benchmarked container: %s", err)
		}

		if err := r.connectBenchmarked(ctx); err != nil {
			log.Panicf("Unable to connect the benchmarked to the network: %s", err)
		}

		log.Printf(">>> profiled with fast ingester benchmark %d/%d", i+1, n)
		t0 := time.Now()
		if err := r.runBenchmarked(ctx); err != nil {
			log.Panicf("Failed to run benchmarked container: %s", err)
		}
		fast = append(fast, time.Since(t0))

		r.removeBenchmarked(ctx)
	}

	// Slow ingester
	r.removeIngester(ctx)
	if err := r.createIngester(ctx, "slow"); err != nil {
		log.Panicf("Unable to create ingester container: %s", err)
	}

	if err := r.connectIngester(ctx); err != nil {
		log.Panicf("Unable to connect the ingester to the network: %s", err)
	}

	if err := r.startIngester(ctx); err != nil {
		log.Panicf("Unable to start ingester: %s", err)
	}

	var slow []time.Duration
	for i := 0; i < n; i++ {
		if err := r.createBenchmarked(ctx, true); err != nil {
			log.Panicf("Unable to create benchmarked container: %s", err)
		}

		if err := r.connectBenchmarked(ctx); err != nil {
			log.Panicf("Unable to connect the benchmarked to the network: %s", err)
		}

		log.Printf(">>> profiled with slow ingester benchmark %d/%d", i+1, n)
		t0 := time.Now()
		if err := r.runBenchmarked(ctx); err != nil {
			log.Panicf("Failed to run benchmarked container: %s", err)
		}
		slow = append(slow, time.Since(t0))

		r.removeBenchmarked(ctx)
	}

	// Remove the container and network, not needed anymore.
	r.removeIngester(ctx)
	r.removeNetwork(ctx)

	var noserver []time.Duration
	for i := 0; i < n; i++ {
		if err := r.createBenchmarked(ctx, true); err != nil {
			log.Panicf("Unable to create benchmarked container: %s", err)
		}

		log.Printf(">>> profiled without ingester benchmark %d/%d", i+1, n)
		t0 := time.Now()
		if err := r.runBenchmarked(ctx); err != nil {
			log.Panicf("Failed to run benchmarked container: %s", err)
		}
		noserver = append(noserver, time.Since(t0))

		r.removeBenchmarked(ctx)
	}

	var base []time.Duration
	for i := 0; i < n; i++ {
		if err := r.createBenchmarked(ctx, false); err != nil {
			log.Panicf("Unable to create benchmarked container: %s", err)
		}

		log.Printf(">>> non-profiled (baseline) benchmark %d/%d", i+1, n)
		t0 := time.Now()
		if err := r.runBenchmarked(ctx); err != nil {
			log.Panicf("Failed to run benchmarked container: %s", err)
		}
		base = append(base, time.Since(t0))

		r.removeBenchmarked(ctx)
	}

	for _, r := range base {
		noprof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "fast"), r))
		noprof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "slow"), r))
		noprof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "noserver"), r))
	}
	for _, r := range fast {
		prof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "fast"), r))
	}
	for _, r := range slow {
		prof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "slow"), r))
	}
	for _, r := range noserver {
		prof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "noserver"), r))
	}
}

func agentName(name string) string {
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, "/", "")
	return name
}

func benchmarkName(name, bench string) string {
	return "Benchmark" + name + "-" + bench
}
