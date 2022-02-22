/*
   It uses a name of a benchmark to run.
   The name of the benchmark is actually a directory which contains a Dockerfile.
   The Dockerfile must contain instructions to build and image of the benchmark to run.
   The benchmark will be a program that will accept the following environment option:
   PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING
   When the envvar is set, the benchmarked program should enable profiling.
   When the envvar is not set, the benchmarked program should disable profiling.

   The runner will do the following:
   - Build the ingestor image.
   - Build the benchmarked image.
   - Create a network.
   - Create the ingestor container.
   - Run the ingestor container.
   - Create a benchmarked container.
   - Run the benchmarked container with profiling enabled
   - When the benchmarked container finishes, destroy the benchmarked container
   - Drop the network.
   - Create a benchmarked container.
   - Run the benchmarked container with profiling enabled.
   - When the benchmarked container finishes, destroy the benchmarked container.
   - Create a benchmarked container.
   - Run the benchmarked container with profiling disabled.
   - When the benchmarked container finishes, destroy the benchmarked container.
   - Give the results of the three runs and compare them.
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

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/perf/benchstat"
)

const (
	id           = "pyroscope-agent-benchmark"
	ingestorPath = "../ingestor"
	n            = 5
)

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
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
	benchstat.FormatText(os.Stdout, c.Tables())
}

func run(ctx context.Context, r *runner, name string, noprof, prof *bytes.Buffer) {
	agent := agentName(name)
	log.Printf("Running %s benchmark", agent)
	if err := r.buildImage(ctx, ingestorPath, id+"/ingestor"); err != nil {
		log.Fatal(err)
	}

	if err := r.buildImage(ctx, name, id+"/benchmarked"); err != nil {
		log.Fatal(err)
	}

	if err := r.createNetwork(ctx); err != nil {
		log.Fatal(err)
	}
	defer r.removeNetwork(ctx)

	if err := r.createIngestor(ctx); err != nil {
		log.Fatal(err)
	}
	defer r.removeIngestor(ctx)

	if err := r.connectIngestor(ctx); err != nil {
		log.Fatal(err)
	}

	if err := r.startIngestor(ctx); err != nil {
		log.Fatal(err)
	}

	var r0 []time.Duration
	for i := 0; i < n; i++ {
		if err := r.createBenchmarked(ctx, true); err != nil {
			log.Fatal(err)
		}
		defer r.removeBenchmarked(ctx)

		if err := r.connectBenchmarked(ctx); err != nil {
			log.Fatal(err)
		}

		log.Printf(">>> profiled with server benchmark %d/%d", i+1, n)
		t0 := time.Now()
		if err := r.startBenchmarked(ctx); err != nil {
			log.Fatal(err)
		}

		statusCh, errCh := r.cli.ContainerWait(ctx, r.benchmarkedID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				log.Fatal(err)
			}
		case <-statusCh:
		}
		r0 = append(r0, time.Since(t0))
	}

	// Delete the container, drop the network and start again.
	r.removeIngestor(ctx)
	r.removeNetwork(ctx)

	var r1 []time.Duration
	for i := 0; i < n; i++ {
		r.removeBenchmarked(ctx)
		if err := r.createBenchmarked(ctx, true); err != nil {
			log.Fatal(err)
		}

		log.Printf(">>> profiled without server benchmark %d/%d", i+1, n)
		t0 := time.Now()
		if err := r.startBenchmarked(ctx); err != nil {
			log.Fatal(err)
		}

		statusCh, errCh := r.cli.ContainerWait(ctx, r.benchmarkedID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}
		r1 = append(r1, time.Since(t0))
	}

	var r2 []time.Duration
	for i := 0; i < n; i++ {
		// Delete the container, start again without profiling.
		r.removeBenchmarked(ctx)
		if err := r.createBenchmarked(ctx, false); err != nil {
			log.Fatal(err)
		}

		log.Printf(">>> non-profiled benchmark %d/%d", i+1, n)
		t0 := time.Now()
		if err := r.startBenchmarked(ctx); err != nil {
			log.Fatal(err)
		}

		statusCh, errCh := r.cli.ContainerWait(ctx, r.benchmarkedID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				log.Fatal(err)
			}
		case <-statusCh:
		}
		r2 = append(r2, time.Since(t0))
	}

	for _, r := range r2 {
		noprof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "fast"), r))
		noprof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "noserver"), r))
	}
	for _, r := range r0 {
		prof.WriteString(fmt.Sprintf("%s 1 %d ns/op\n", benchmarkName(agent, "fast"), r))
	}
	for _, r := range r1 {
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
